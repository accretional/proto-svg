// render is the "SVG API" renderer: it takes a textproto describing an
// svg.SvgDocument (the data a service feeds in) and emits the .svg markup —
// with ZERO knowledge of SVG. It only knows how to walk a protobuf message:
//
//  1. load the grammar-compiled schema (proto/svg.fdset),
//  2. unmarshal the textproto into a dynamic SvgDocument message,
//  3. walk that message by proto reflection, and for every message emit its
//     MessagePrefix (the markup terminals gluon lifted out of the grammar —
//     "<rect", ` x="`, ">", "</rect>", …), recurse into child messages in
//     field-number order, write string leaves (the free terminal values)
//     verbatim, and interleave FieldSeparator between repeated items.
//
// That is the whole program. There is no flame-graph logic, no per-element
// special-casing, no layout — the grammar (via the schema + the prefix /
// separator maps) is the single source of truth for what the SVG looks like.
//
//	go run ./examples examples/flamegraph.textproto > out.svg
//	go run ./examples examples/flamegraph.textproto flamegraph.svg
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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: render <doc.textproto> [out.svg]")
		os.Exit(2)
	}
	in := os.Args[1]

	// 1. Load the grammar-compiled schema descriptors.
	desc, err := loadRootDescriptor(fdsetPath, rootMsg)
	if err != nil {
		fail(err)
	}

	// 2. Unmarshal the textproto into a dynamic SvgDocument. (No type resolver
	// is needed: every nested type is a plain message reachable from the
	// SvgDocument descriptor — there are no Any / extension fields.)
	raw, err := os.ReadFile(in)
	if err != nil {
		fail(err)
	}
	doc := dynamicpb.NewMessage(desc)
	if err := prototext.Unmarshal(raw, doc); err != nil {
		fail(fmt.Errorf("parse %s: %w", in, err))
	}

	// 3. Walk the message graph → SVG markup.
	var sb strings.Builder
	render(doc.ProtoReflect(), &sb)

	out := sb.String()
	if len(os.Args) >= 3 {
		if err := os.WriteFile(os.Args[2], []byte(out), 0o644); err != nil {
			fail(err)
		}
		fmt.Fprintf(os.Stderr, "wrote %d bytes -> %s\n", len(out), os.Args[2])
		return
	}
	fmt.Print(out)
}

// render emits one message: its prefix tokens, then each populated field in
// field-number order. This is the entire SVG-emission logic.
func render(m protoreflect.Message, sb *strings.Builder) {
	fqn := "." + string(m.Descriptor().FullName())
	for _, tok := range svgpb.MessagePrefix[fqn] {
		sb.WriteString(tok)
	}
	fields := m.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if !m.Has(fd) {
			continue
		}
		val := m.Get(fd)
		if fd.IsList() {
			sep := svgpb.FieldSeparator["."+string(fd.FullName())]
			list := val.List()
			for j := 0; j < list.Len(); j++ {
				if j > 0 {
					sb.WriteString(sep)
				}
				renderValue(fd, list.Get(j), sb)
			}
			continue
		}
		renderValue(fd, val, sb)
	}
}

// renderValue writes one field value: recurse into messages, write the free
// terminal (string leaf) verbatim.
func renderValue(fd protoreflect.FieldDescriptor, v protoreflect.Value, sb *strings.Builder) {
	if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
		render(v.Message(), sb)
		return
	}
	sb.WriteString(v.String())
}

// loadRootDescriptor reads a FileDescriptorSet and returns the descriptor for
// the named root message plus a type resolver over the whole set.
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
	fmt.Fprintln(os.Stderr, "render:", err)
	os.Exit(1)
}
