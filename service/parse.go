package service

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	cssservice "github.com/accretional/proto-css/service"
	svgpb "github.com/accretional/proto-svg/proto/pb/svg"
)

// foreignParsers maps a proto package name to a parser for that grammar's text.
// At a seam (a field whose message belongs to another grammar) the SVG parser
// captures the seam's text and hands it to that grammar's own parser, so the
// <style> body becomes a real css.CssStyleSheet subtree, not an opaque string.
var foreignParsers = map[string]func(input, typeName string) (proto.Message, error){
	"css": cssservice.ParseAs,
}

// Parse parses SVG text into a typed SvgDocument (the grammar's start symbol).
// Pure reflection over the generated svg schema plus the prefix/separator
// tables — the proto schema IS the grammar. Embedded grammars (CSS) are
// delegated at the package boundary. Inverse of Render.
func Parse(input string) (*svgpb.SvgDocument, error) {
	msg, err := ParseAs(input, "svg.SvgDocument")
	if err != nil {
		return nil, err
	}
	doc, ok := msg.(*svgpb.SvgDocument)
	if !ok {
		return nil, fmt.Errorf("internal: parsed %T", msg)
	}
	return doc, nil
}

// ParseAs parses SVG text against an arbitrary svg.* message type (the start
// symbol), e.g. "svg.SVGCircleElement". Same engine as Parse, rooted wherever
// the caller asks.
func ParseAs(input, typeName string) (proto.Message, error) {
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(typeName))
	if err != nil {
		return nil, fmt.Errorf("%s not registered: %w", typeName, err)
	}
	msg := mt.New()

	p := &parser{prefix: svgpb.MessagePrefix, sep: svgpb.FieldSeparator}
	pos, err := p.parseMsg(input, p.skipWS(input, 0), msg, nil)
	if err != nil {
		return nil, err
	}
	pos = p.skipWS(input, pos)
	if pos < len(input) {
		return nil, fmt.Errorf("unconsumed input at %d: %q", pos, snippet(input[pos:]))
	}
	return msg.Interface(), nil
}

type parser struct {
	prefix map[string][]string
	sep    map[string]string
	depth  int
	steps  int
}

// maxParseDepth bounds recursion; maxParseSteps bounds total work. SVG content
// models are recursive (containers nest arbitrarily), so a naive longest-match
// descent can recurse or backtrack without bound. These caps make it fail fast.
const (
	maxParseDepth = 400
	maxParseSteps = 5_000_000
)

func (p *parser) parseMsg(input string, pos int, msg protoreflect.Message, outerStops []string) (int, error) {
	if p.depth > maxParseDepth {
		return pos, fmt.Errorf("max parse depth exceeded")
	}
	if p.steps++; p.steps > maxParseSteps {
		return pos, fmt.Errorf("parse step budget exceeded")
	}
	p.depth++
	defer func() { p.depth-- }()

	md := msg.Descriptor()
	fqn := "." + string(md.FullName())

	// 1. Consume the message's leading prefix tokens (markup/keywords).
	if pfx, ok := p.prefix[fqn]; ok {
		for _, tok := range pfx {
			pos = p.skipWS(input, pos)
			if !strings.HasPrefix(input[pos:], tok) {
				return pos, fmt.Errorf("expected %q for %s at %d", tok, md.Name(), pos)
			}
			pos += len(tok)
		}
	}

	// 2. Scalar leaf — caller captures the text via parseScalar.
	if isScalar(md) {
		return pos, nil
	}

	// 3. Walk fields in declaration order.
	fields := md.Fields()
	handledOneofs := map[int]bool{}
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		stops := p.fieldStops(md, i, outerStops)

		if od := fd.ContainingOneof(); od != nil {
			if handledOneofs[od.Index()] {
				continue
			}
			handledOneofs[od.Index()] = true
			pos = p.parseOneof(input, pos, msg, od, stops)
			continue
		}
		if fd.IsList() {
			pos = p.parseRepeated(input, pos, msg, fd, stops)
			continue
		}
		if np, err := p.parseSingular(input, pos, msg, fd, stops); err == nil {
			pos = np
		}
	}
	return pos, nil
}

func (p *parser) parseSingular(input string, pos int, msg protoreflect.Message, fd protoreflect.FieldDescriptor, stops []string) (int, error) {
	if fd.Kind() != protoreflect.MessageKind {
		text, np := matchUntilAny(input, p.skipWS(input, pos), stops)
		text = strings.TrimSpace(text)
		if text != "" {
			msg.Set(fd, protoreflect.ValueOfString(text))
		}
		return np, nil
	}
	// Grammar seam: a field whose message belongs to another grammar's package
	// (css.CssStyleSheet) is captured as text and handed to that grammar's
	// parser, which returns the structured subtree to embed.
	if pkg := string(fd.Message().ParentFile().Package()); pkg != localPackage {
		return p.parseForeign(input, pos, msg, fd, pkg, stops)
	}
	sub := newSub(fd.Message())
	if sub == nil {
		return pos, fmt.Errorf("cannot create %s", fd.Message().FullName())
	}
	if isScalar(fd.Message()) {
		return p.parseScalar(input, pos, msg, fd, sub, stops)
	}
	np, err := p.parseMsg(input, pos, sub, stops)
	if err != nil {
		return pos, err
	}
	msg.Set(fd, protoreflect.ValueOfMessage(sub))
	return np, nil
}

// parseForeign captures the seam text up to the next stop and delegates to the
// embedded grammar's parser, embedding the returned message.
func (p *parser) parseForeign(input string, pos int, msg protoreflect.Message, fd protoreflect.FieldDescriptor, pkg string, stops []string) (int, error) {
	parse, ok := foreignParsers[pkg]
	if !ok {
		return pos, fmt.Errorf("no parser registered for package %q (%s)", pkg, fd.Message().FullName())
	}
	text, np := matchUntilAny(input, p.skipWS(input, pos), stops)
	text = strings.TrimSpace(text)
	if text == "" {
		return pos, nil
	}
	sub, err := parse(text, string(fd.Message().FullName()))
	if err != nil {
		return pos, fmt.Errorf("parse %s seam %q: %w", pkg, snippet(text), err)
	}
	msg.Set(fd, protoreflect.ValueOfMessage(sub.ProtoReflect()))
	return np, nil
}

func (p *parser) parseScalar(input string, pos int, parent protoreflect.Message, fd protoreflect.FieldDescriptor, sub protoreflect.Message, stops []string) (int, error) {
	text, np := matchUntilAny(input, p.skipWS(input, pos), stops)
	text = strings.TrimSpace(text)
	if text == "" {
		return pos, fmt.Errorf("empty scalar for %s", fd.Name())
	}
	if vfd := sub.Descriptor().Fields().ByName("value"); vfd != nil {
		sub.Set(vfd, protoreflect.ValueOfString(text))
	}
	parent.Set(fd, protoreflect.ValueOfMessage(sub))
	return np, nil
}

func (p *parser) parseOneof(input string, pos int, msg protoreflect.Message, od protoreflect.OneofDescriptor, stops []string) int {
	bestPos := pos
	var bestFD protoreflect.FieldDescriptor
	var bestMsg protoreflect.Message

	for i := 0; i < od.Fields().Len(); i++ {
		fd := od.Fields().Get(i)
		if fd.Kind() != protoreflect.MessageKind {
			continue
		}
		sub := newSub(fd.Message())
		if sub == nil {
			continue
		}
		if isScalar(fd.Message()) {
			continue
		}
		if np, err := p.parseMsg(input, pos, sub, stops); err == nil && np > bestPos {
			bestPos, bestFD, bestMsg = np, fd, sub
		}
	}
	if bestFD == nil {
		for i := 0; i < od.Fields().Len(); i++ {
			fd := od.Fields().Get(i)
			if fd.Kind() != protoreflect.MessageKind || !isScalar(fd.Message()) {
				continue
			}
			sub := newSub(fd.Message())
			text, np := matchUntilAny(input, p.skipWS(input, pos), stops)
			text = strings.TrimSpace(text)
			if text == "" {
				continue
			}
			if vfd := fd.Message().Fields().ByName("value"); vfd != nil {
				sub.Set(vfd, protoreflect.ValueOfString(text))
			}
			bestPos, bestFD, bestMsg = np, fd, sub
			break
		}
	}
	if bestFD != nil {
		msg.Set(bestFD, protoreflect.ValueOfMessage(bestMsg))
	}
	return bestPos
}

func (p *parser) parseRepeated(input string, pos int, msg protoreflect.Message, fd protoreflect.FieldDescriptor, outerStops []string) int {
	if fd.Kind() != protoreflect.MessageKind {
		return pos
	}
	list := msg.Mutable(fd).List()
	sepKey := "." + string(msg.Descriptor().FullName()) + "." + string(fd.Name())
	sep := p.sep[sepKey]

	for {
		tryPos := pos
		if list.Len() > 0 && sep != "" {
			tryPos = p.skipWS(input, tryPos)
			if !strings.HasPrefix(input[tryPos:], sep) {
				break
			}
			tryPos += len(sep)
		}
		sub := newSub(fd.Message())
		if sub == nil {
			break
		}
		np, err := p.parseMsg(input, tryPos, sub, outerStops)
		if err != nil || np <= tryPos {
			break
		}
		list.Append(protoreflect.ValueOfMessage(sub))
		pos = np
	}
	return pos
}

func (p *parser) fieldStops(md protoreflect.MessageDescriptor, fieldIdx int, outerStops []string) []string {
	stops := append([]string(nil), outerStops...)
	fields := md.Fields()
	var skipOneof protoreflect.OneofDescriptor
	if fieldIdx < fields.Len() {
		skipOneof = fields.Get(fieldIdx).ContainingOneof()
	}
	handled := map[int]bool{}
	for j := fieldIdx + 1; j < fields.Len(); j++ {
		fd := fields.Get(j)
		if od := fd.ContainingOneof(); od != nil {
			if skipOneof != nil && od.Index() == skipOneof.Index() {
				continue
			}
			if handled[od.Index()] {
				continue
			}
			handled[od.Index()] = true
			for k := 0; k < od.Fields().Len(); k++ {
				if vfd := od.Fields().Get(k); vfd.Kind() == protoreflect.MessageKind {
					if t := p.leadingTerminal(vfd.Message()); t != "" {
						stops = append(stops, t)
					}
				}
			}
		} else if fd.Kind() == protoreflect.MessageKind {
			if t := p.leadingTerminal(fd.Message()); t != "" {
				stops = append(stops, t)
			}
		}
	}
	return stops
}

// leadingTerminal returns the first prefix token reachable at the front of a
// message. SVG content models are recursive, so recursion is guarded by a
// visited set.
func (p *parser) leadingTerminal(md protoreflect.MessageDescriptor) string {
	return p.leadingTerminalSeen(md, map[protoreflect.FullName]bool{})
}

func (p *parser) leadingTerminalSeen(md protoreflect.MessageDescriptor, seen map[protoreflect.FullName]bool) string {
	fqn := "." + string(md.FullName())
	if pfx, ok := p.prefix[fqn]; ok && len(pfx) > 0 {
		return pfx[0]
	}
	if seen[md.FullName()] {
		return ""
	}
	seen[md.FullName()] = true
	if md.Fields().Len() > 0 {
		fd := md.Fields().Get(0)
		if od := fd.ContainingOneof(); od != nil {
			for i := 0; i < od.Fields().Len(); i++ {
				if vfd := od.Fields().Get(i); vfd.Kind() == protoreflect.MessageKind {
					if t := p.leadingTerminalSeen(vfd.Message(), seen); t != "" {
						return t
					}
				}
			}
			return ""
		}
		if fd.Kind() == protoreflect.MessageKind {
			return p.leadingTerminalSeen(fd.Message(), seen)
		}
	}
	return ""
}

// skipWS is a no-op: like HTML, the SVG grammar bakes whitespace into its
// terminals, and rendered SVG carries no insignificant whitespace, so matching
// is exact and rendered output round-trips.
func (p *parser) skipWS(input string, pos int) int {
	return pos
}

func matchUntilAny(input string, pos int, stops []string) (string, int) {
	end := len(input)
	for _, s := range stops {
		if s == "" {
			continue
		}
		if idx := strings.Index(input[pos:], s); idx >= 0 && pos+idx < end {
			end = pos + idx
		}
	}
	return input[pos:end], end
}

func isScalar(md protoreflect.MessageDescriptor) bool {
	fields := md.Fields()
	if fields.Len() != 1 {
		return false
	}
	fd := fields.Get(0)
	return string(fd.Name()) == "value" && fd.Kind() == protoreflect.StringKind
}

func newSub(md protoreflect.MessageDescriptor) protoreflect.Message {
	mt, err := protoregistry.GlobalTypes.FindMessageByName(md.FullName())
	if err != nil {
		return nil
	}
	return mt.New()
}

func snippet(s string) string {
	if len(s) > 48 {
		return s[:48] + "…"
	}
	return s
}
