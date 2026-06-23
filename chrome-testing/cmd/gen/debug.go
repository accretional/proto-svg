package main

import (
	"fmt"
	"sort"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
)

// dumpElement prints the in-memory field structure of an element message and a
// couple of its attribute arms, so the enumerator can be written against the
// real (gen-pipeline) FDP shape rather than the on-disk genproto schema.
func dumpElement(byFQN map[string]*descriptorpb.DescriptorProto, kw map[string]string, openTag string) {
	fqn := findElementFQN(byFQN, openTag)
	fmt.Printf("ELEMENT %q -> %s\n", openTag, fqn)
	m := byFQN[fqn]
	if m == nil {
		fmt.Println("  <nil>")
		return
	}
	dumpMsg(byFQN, kw, fqn, 0, map[string]int{})
}

func dumpMsg(byFQN map[string]*descriptorpb.DescriptorProto, kw map[string]string, fqn string, depth int, seen map[string]int) {
	ind := strings.Repeat("  ", depth)
	if lit, ok := kw[fqn]; ok {
		fmt.Printf("%sKW %s = %q\n", ind, simpleName(fqn), lit)
		return
	}
	if samples, ok := reps[simpleName(fqn)]; ok {
		fmt.Printf("%sLEAF %s = %v\n", ind, simpleName(fqn), samples)
		return
	}
	m := byFQN[fqn]
	if m == nil {
		fmt.Printf("%sUNRESOLVED %s\n", ind, fqn)
		return
	}
	if seen[fqn] >= 1 || depth > 9 { // depth bumped for structured-value inspection
		fmt.Printf("%s... %s\n", ind, simpleName(fqn))
		return
	}
	seen[fqn]++
	defer func() { seen[fqn]-- }()
	kind := "SEQ"
	if len(m.GetOneofDecl()) > 0 {
		kind = "ONEOF"
	}
	fmt.Printf("%s%s %s {\n", ind, kind, simpleName(fqn))
	for _, f := range m.GetField() {
		rep := ""
		if f.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			rep = "repeated "
		}
		fmt.Printf("%s  %s%s -> %s\n", ind, rep, f.GetName(), simpleName(f.GetTypeName()))
		dumpMsg(byFQN, kw, f.GetTypeName(), depth+2, seen)
	}
	fmt.Printf("%s}\n", ind)
}

// dumpAllElements lists every element message FQN and its open tag.
func dumpAllElements(byFQN map[string]*descriptorpb.DescriptorProto, kw map[string]string) {
	type el struct{ tag, fqn string }
	var els []el
	for fqn, m := range byFQN {
		fields := m.GetField()
		if len(fields) == 0 {
			continue
		}
		lit := kw[fields[0].GetTypeName()]
		if strings.HasPrefix(lit, "<") && !strings.HasPrefix(lit, "</") {
			els = append(els, el{lit, fqn})
		}
	}
	sort.Slice(els, func(i, j int) bool { return els[i].tag < els[j].tag })
	for _, e := range els {
		fmt.Printf("%-22s %s\n", e.tag, e.fqn)
	}
	fmt.Printf("TOTAL %d element messages\n", len(els))
}
