// Package service renders typed SVG messages to SVG text and parses SVG text
// back, both purely by reflection over the grammar-derived schema in svg.proto
// plus the generated tables (svgpb.MessagePrefix, svgpb.FieldSeparator). There
// is no per-element logic.
//
// The SVG grammar embeds the CSS grammar at the <style> seam (its content is a
// css.CssStyleSheet). This renderer/parser is grammar-agnostic: at a field
// whose message belongs to another grammar's proto package it delegates to that
// grammar's own renderer/parser (proto-css's service). No CSS knowledge lives
// here — only the dispatch at the package boundary.
package service

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	cssservice "github.com/accretional/proto-css/service"
	svgpb "github.com/accretional/proto-svg/proto/pb/svg"
)

// localPackage is the proto package of this grammar's own messages. Messages in
// any other package are handed to a foreign renderer (see foreignRenderers).
const localPackage = "svg"

// foreignRenderers maps a proto package name to a renderer for messages in that
// package. The SVG renderer delegates at grammar seams: a css.CssStyleSheet
// embedded in <style> is handed to proto-css's renderer, which knows CSS
// spacing.
var foreignRenderers = map[string]func(proto.Message) (string, error){
	"css": cssservice.Render,
}

// Render serializes any typed SVG message back into SVG text. It walks the
// message via reflection: it emits each message's leading terminal tokens
// recorded in svgpb.MessagePrefix (the markup the compiler's StripKeywords moved
// out of the schema, e.g. "<rect", ` x="`, ">"), recurses into populated message
// fields in proto declaration order, emits scalar string leaves verbatim,
// interleaves svgpb.FieldSeparator between repeated elements, and delegates any
// foreign-package subtree (css.*) to that grammar's renderer.
//
// Like HTML, the SVG grammar bakes its spacing into the terminals, so tokens are
// concatenated directly.
func Render(msg proto.Message) (string, error) {
	var toks []string
	if err := renderMessage(msg.ProtoReflect(), &toks); err != nil {
		return "", err
	}
	return strings.Join(toks, ""), nil
}

func renderMessage(m protoreflect.Message, toks *[]string) error {
	fqn := "." + string(m.Descriptor().FullName())
	if prefix, ok := svgpb.MessagePrefix[fqn]; ok {
		*toks = append(*toks, prefix...)
	}

	fds := m.Descriptor().Fields()
	for i := 0; i < fds.Len(); i++ {
		fd := fds.Get(i)
		if !m.Has(fd) {
			continue
		}
		val := m.Get(fd)
		if fd.IsList() {
			list := val.List()
			sep := svgpb.FieldSeparator[fqn+"."+string(fd.Name())]
			for j := 0; j < list.Len(); j++ {
				if j > 0 && sep != "" {
					*toks = append(*toks, sep)
				}
				if err := renderValue(fd, list.Get(j), toks); err != nil {
					return err
				}
			}
			continue
		}
		if err := renderValue(fd, val, toks); err != nil {
			return err
		}
	}
	return nil
}

func renderValue(fd protoreflect.FieldDescriptor, val protoreflect.Value, toks *[]string) error {
	switch fd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		sub := val.Message()
		pkg := string(sub.Descriptor().ParentFile().Package())
		if pkg != localPackage {
			render, ok := foreignRenderers[pkg]
			if !ok {
				return fmt.Errorf("no renderer registered for package %q (%s)", pkg, sub.Descriptor().FullName())
			}
			out, err := render(sub.Interface())
			if err != nil {
				return fmt.Errorf("render %s subtree %s: %w", pkg, sub.Descriptor().FullName(), err)
			}
			*toks = append(*toks, out)
			return nil
		}
		return renderMessage(sub, toks)
	case protoreflect.StringKind:
		if s := val.String(); s != "" {
			*toks = append(*toks, s)
		}
		return nil
	default:
		return fmt.Errorf("unsupported field kind %v at %s", fd.Kind(), fd.FullName())
	}
}
