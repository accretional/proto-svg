package main

import (
	"github.com/accretional/gluon/v2/compiler"
	pb "github.com/accretional/gluon/v2/pb"
)

// pruneUnreachable drops every rule not reachable from the given root rules by
// following nonterminal references. The full SVG grammar carries a long tail of
// low-level primitives (digit/hex-digit/ident-char character ranges) that,
// after the leaf types are scalarized, no reachable rule references anymore.
// Removing them shrinks the emitted schema to the reachable core and — because
// the surviving range-bearing rules are exactly those orphans — drops the
// unicode/utf_8.proto import, leaving svg.proto self-contained.
//
// Scalarizing a PascalCase leaf wrapper (e.g. LengthType) also severs its
// reference to the snake_case atom (length_type / number), so prune then
// removes the atom, dissolving the name collision the two would otherwise
// share.
//
// The input is not mutated; a new file node with the surviving rules (in
// original order) is returned.
func pruneUnreachable(root *pb.ASTNode, roots ...string) (*pb.ASTNode, []string) {
	byName := map[string]*pb.ASTNode{}
	for _, r := range root.GetChildren() {
		if r.GetKind() == compiler.KindRule {
			byName[r.GetValue()] = r
		}
	}

	reachable := map[string]bool{}
	var visit func(name string)
	visit = func(name string) {
		if reachable[name] {
			return
		}
		r, ok := byName[name]
		if !ok {
			return // terminal-only / external symbol
		}
		reachable[name] = true
		collectRefs(r, visit)
	}
	for _, rt := range roots {
		visit(rt)
	}

	var kids []*pb.ASTNode
	var dropped []string
	for _, r := range root.GetChildren() {
		if r.GetKind() == compiler.KindRule && !reachable[r.GetValue()] {
			dropped = append(dropped, r.GetValue())
			continue
		}
		kids = append(kids, r)
	}
	return &pb.ASTNode{Kind: root.GetKind(), Value: root.GetValue(), Children: kids}, dropped
}

// collectRefs invokes fn for every nonterminal name referenced anywhere in the
// node's subtree.
func collectRefs(node *pb.ASTNode, fn func(string)) {
	if node == nil {
		return
	}
	if node.GetKind() == compiler.KindNonterminal {
		fn(node.GetValue())
	}
	for _, c := range node.GetChildren() {
		collectRefs(c, fn)
	}
}
