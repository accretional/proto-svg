package main

import (
	"github.com/accretional/gluon/v2/compiler"
	pb "github.com/accretional/gluon/v2/pb"
)

// bodyTerminal reports whether a rule's body is effectively a single terminal
// literal (a pure keyword or punctuation rule, e.g. `meet = "meet"`,
// `colon_symbol = ":"`). It unwraps a single enclosing sequence or group.
// Returns the literal and true on match.
func bodyTerminal(rule *pb.ASTNode) (string, bool) {
	if rule.GetKind() != compiler.KindRule {
		return "", false
	}
	body := rule.GetChildren()
	if len(body) != 1 {
		return "", false
	}
	n := body[0]
	for (n.GetKind() == compiler.KindSequence || n.GetKind() == compiler.KindGroup) && len(n.GetChildren()) == 1 {
		n = n.GetChildren()[0]
	}
	if n.GetKind() == compiler.KindTerminal {
		return n.GetValue(), true
	}
	return "", false
}

// emptyKeywordRules replaces the body of every pure-keyword rule with nothing,
// so it lowers to an empty marker message (e.g. `message Meet {}`). The rule's
// literal is recorded in the prefix map during the prefix pass, so the renderer
// re-emits "meet" when it reaches a Meet, and a caller selects the keyword by
// instantiating the empty `&Meet{}` rather than building a nested keyword
// wrapper. This removes the redundant *Keyword wrapper messages the compiler
// would otherwise dedup into existence.
//
// Must run AFTER the prefix pass has recorded the literals and BEFORE the final
// Compile. The input is not mutated.
func emptyKeywordRules(root *pb.ASTNode) *pb.ASTNode {
	if root == nil {
		return nil
	}
	if _, ok := bodyTerminal(root); ok {
		return &pb.ASTNode{Kind: compiler.KindRule, Value: root.GetValue()}
	}
	kids := make([]*pb.ASTNode, 0, len(root.GetChildren()))
	for _, c := range root.GetChildren() {
		kids = append(kids, emptyKeywordRules(c))
	}
	return &pb.ASTNode{Kind: root.GetKind(), Value: root.GetValue(), Children: kids}
}
