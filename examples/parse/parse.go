// parse is the inverse of examples/render: it reads SVG markup and rebuilds the
// svg.SvgDocument proto (emitted as a textproto), by recursive descent over the
// proto descriptors. The descriptor tree IS the grammar — so the same two maps
// render emits with are read here in reverse:
//
//   - MessagePrefix[fqn]  = the literal terminals a production begins with
//                           ("<rect", ` x="`, ">", "</rect>", `"`). To parse a
//                           message we must first match (consume) these.
//   - a message's fields  = the sequence to match, in field-number order.
//   - a oneof             = an alternation: try each arm; the one whose literals
//                           match the input wins (first viable arm).
//   - a repeated field    = { X }: parse elements until one fails to match,
//                           consuming FieldSeparator between them.
//   - a {string value}    = a free terminal: consume input up to the next known
//     leaf                  literal (the "follow" — usually the closing `"`).
//
// There is NO SVG knowledge here. The parser only matches literals, recurses
// through fields, picks oneof arms, and slices free terminals at the next
// sentinel. Feeding examples/render's own output back through this and
// re-rendering reproduces the byte-identical SVG (the faithfulness guarantee).
//
//   go run ./examples/parse examples/parse/flamegraph.svg                 # -> stdout
//   go run ./examples/parse examples/parse/flamegraph.svg out.textproto
package main

import (
	"fmt"
	"os"
	"strings"

	svgpb "github.com/accretional/proto-svg/proto/pb/svg"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

const (
	fdsetPath = "proto/svg.fdset"
	rootMsg   = "svg.SvgDocument"
)

// Furthest-failure tracking, for a useful error when the input doesn't match:
// the deepest byte offset any prefix match reached before failing.
var (
	src          string
	furthest     int
	furthestWhat string
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: parse <in.svg> [out.textproto]")
		os.Exit(2)
	}
	desc, err := loadRootDescriptor(fdsetPath, rootMsg)
	if err != nil {
		fail(err)
	}
	srcBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		fail(err)
	}
	src = string(srcBytes)

	m, rest, ok := parseMsg(desc, src, "")
	if !ok {
		fail(fmt.Errorf("does not match the SVG grammar at byte %d (%s)\n  near: %q",
			furthest, furthestWhat, src[furthest:min(furthest+60, len(src))]))
	}
	if strings.TrimSpace(rest) != "" {
		fail(fmt.Errorf("trailing input not consumed (%d bytes): %.60q…", len(rest), rest))
	}

	out, err := (prototext.MarshalOptions{Multiline: true, Indent: "  "}).Marshal(m.Interface())
	if err != nil {
		fail(err)
	}
	if len(os.Args) >= 3 {
		if err := os.WriteFile(os.Args[2], out, 0o644); err != nil {
			fail(err)
		}
		fmt.Fprintf(os.Stderr, "parsed %d bytes SVG -> %s (%d bytes textproto)\n", len(src), os.Args[2], len(out))
		return
	}
	os.Stdout.Write(out)
}

// parseMsg attempts to parse one message (described by md) from the front of in.
// follow is the literal that terminates this message's trailing free terminal,
// when its last field is a {string value} leaf (it bubbles down from the parent).
// Returns the populated message, the unconsumed remainder, and whether it matched.
func parseMsg(md protoreflect.MessageDescriptor, in, follow string) (protoreflect.Message, string, bool) {
	rest := in
	pfx := svgpb.MessagePrefix["."+string(md.FullName())]
	// LEXICAL: whitespace before a tag boundary ("<rect", "</g>") is insignificant
	// in real SVG (indentation, newlines) but absent from the structural grammar.
	// Skip it before matching a tag literal. (Not before ` attr="` / space / `>`
	// literals, whose own leading space is significant, nor before free-terminal
	// leaves, which have no prefix and so preserve their whitespace.)
	if len(pfx) > 0 && strings.HasPrefix(pfx[0], "<") {
		rest = skipWS(rest)
	}
	for _, tok := range pfx {
		if !strings.HasPrefix(rest, tok) {
			if p := len(src) - len(rest); p > furthest {
				furthest, furthestWhat = p, fmt.Sprintf("expected %q for %s", tok, md.FullName())
			}
			return nil, in, false
		}
		rest = rest[len(tok):]
	}

	m := dynamicpb.NewMessage(md)
	fields := md.Fields()
	for i := 0; i < fields.Len(); {
		fd := fields.Get(i)

		// A oneof: try its arms, then skip past all its member fields.
		if oo := fd.ContainingOneof(); oo != nil {
			end := afterOneof(fields, oo, i)
			armFd, armMsg, r2, ok := parseOneof(oo, rest, followAt(fields, end, follow))
			if !ok {
				return nil, in, false
			}
			m.Set(armFd, protoreflect.ValueOfMessage(armMsg))
			rest = r2
			i = end
			continue
		}

		sub := followAt(fields, i+1, follow)
		switch {
		case fd.IsList():
			rest = parseList(m, fd, rest, sub)
		case fd.Kind() == protoreflect.MessageKind:
			// LEXICAL: self-closing shorthand. When the next structural token is
			// the element's ">" terminator but the input reads "/>", the element
			// was written self-closed. Synthesise the canonical long form: set the
			// ">" and the trailing "</tag>" keyword, leave content empty.
			if lit, ok := firstLit(fd.Message()); ok && lit == ">" {
				if ws := skipWS(rest); strings.HasPrefix(ws, "/>") {
					m.Set(fd, protoreflect.ValueOfMessage(dynamicpb.NewMessage(fd.Message())))
					last := fields.Get(fields.Len() - 1)
					m.Set(last, protoreflect.ValueOfMessage(dynamicpb.NewMessage(last.Message())))
					rest = ws[2:]
					i = fields.Len()
					continue
				}
			}
			child, r2, ok := parseMsg(fd.Message(), rest, sub)
			if !ok {
				return nil, in, false
			}
			m.Set(fd, protoreflect.ValueOfMessage(child))
			rest = r2
		case fd.Kind() == protoreflect.StringKind:
			val, r2 := consumeUntil(rest, sub)
			m.Set(fd, protoreflect.ValueOfString(val))
			rest = r2
		default:
			return nil, in, false
		}
		i++
	}
	return m, rest, true
}

// skipWS drops leading XML-insignificant whitespace.
func skipWS(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	return s[i:]
}

// parseOneof tries each arm in declaration order and returns the first that
// parses. The grammar gives productions distinct leading literals, so at most
// one prefixed arm matches; a free-terminal arm (a {string value} with no
// literal) acts as a catch-all and is chosen positionally — which still
// re-renders to identical markup, since every such arm emits its value verbatim.
func parseOneof(oo protoreflect.OneofDescriptor, in, follow string) (protoreflect.FieldDescriptor, protoreflect.Message, string, bool) {
	arms := oo.Fields()
	for i := 0; i < arms.Len(); i++ {
		fd := arms.Get(i)
		if fd.Kind() != protoreflect.MessageKind {
			continue
		}
		if child, rest, ok := parseMsg(fd.Message(), in, follow); ok {
			return fd, child, rest, true
		}
	}
	return nil, nil, in, false
}

// parseList parses a repeated field: { X (SEP X)* }. Elements are parsed until
// one fails to match (the list's natural terminator — e.g. ">" after attributes,
// "</tag>" after content). A non-empty FieldSeparator is consumed between
// elements, but only when another element actually follows it.
func parseList(m protoreflect.Message, fd protoreflect.FieldDescriptor, in, follow string) string {
	sep := svgpb.FieldSeparator["."+string(fd.FullName())]
	lv := m.Mutable(fd).List()
	rest := in
	for {
		child, r2, ok := parseMsg(fd.Message(), rest, follow)
		// Stop on failure OR no progress: a free-terminal arm (e.g. empty text
		// content) can "match" zero input at the list's terminator, which would
		// otherwise loop forever appending empty messages. No progress = no more
		// elements. (r2 is always a suffix of rest, so equal length ⟺ zero width.)
		if !ok || len(r2) >= len(rest) {
			break
		}
		lv.Append(protoreflect.ValueOfMessage(child))
		rest = r2
		if sep == "" {
			continue
		}
		if !strings.HasPrefix(rest, sep) {
			break
		}
		afterSep := rest[len(sep):]
		if _, _, ok := parseMsg(fd.Message(), afterSep, follow); !ok {
			break // the separator was trailing / belongs to the follow
		}
		rest = afterSep
	}
	return rest
}

// followAt returns the literal that a free terminal ending at field index idx
// runs up to: the first literal of the field AT idx (the next sibling — usually
// a closing-quote/space/`>` keyword), or the inherited follow if idx is past the
// end or that field has no determinable leading literal.
func followAt(fields protoreflect.FieldDescriptors, idx int, inherited string) string {
	if idx < fields.Len() {
		nf := fields.Get(idx)
		if nf.Kind() == protoreflect.MessageKind && !nf.IsList() && nf.ContainingOneof() == nil {
			if lit, ok := firstLit(nf.Message()); ok {
				return lit
			}
		}
	}
	return inherited
}

// firstLit is the literal terminal a message necessarily begins with: its own
// prefix, or (recursively) that of its first non-list, non-oneof field.
func firstLit(md protoreflect.MessageDescriptor) (string, bool) {
	if p := svgpb.MessagePrefix["."+string(md.FullName())]; len(p) > 0 {
		return p[0], true
	}
	fields := md.Fields()
	if fields.Len() == 0 {
		return "", false
	}
	f0 := fields.Get(0)
	if f0.Kind() == protoreflect.MessageKind && !f0.IsList() && f0.ContainingOneof() == nil {
		return firstLit(f0.Message())
	}
	return "", false
}

// afterOneof returns the index just past the last contiguous member of oo.
func afterOneof(fields protoreflect.FieldDescriptors, oo protoreflect.OneofDescriptor, i int) int {
	j := i
	for j < fields.Len() && fields.Get(j).ContainingOneof() == oo {
		j++
	}
	return j
}

// consumeUntil splits s at the first occurrence of term, returning the part
// before it (the free terminal value) and the remainder starting at term.
func consumeUntil(s, term string) (string, string) {
	if term == "" {
		return s, ""
	}
	if idx := strings.Index(s, term); idx >= 0 {
		return s[:idx], s[idx:]
	}
	return s, ""
}

func loadRootDescriptor(path, name string) (protoreflect.MessageDescriptor, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var fds descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(b, &fds); err != nil {
		return nil, fmt.Errorf("%s is not a FileDescriptorSet: %w", path, err)
	}
	files, err := protodesc.NewFiles(&fds)
	if err != nil {
		return nil, err
	}
	d, err := files.FindDescriptorByName(protoreflect.FullName(name))
	if err != nil {
		return nil, fmt.Errorf("root message %s: %w", name, err)
	}
	md, ok := d.(protoreflect.MessageDescriptor)
	if !ok {
		return nil, fmt.Errorf("%s is not a message", name)
	}
	return md, nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "parse:", err)
	os.Exit(1)
}
