// gen is the SVG generation ENGINE: it compiles the SVG EBNF grammar (in
// memory, via the gluon v2 metaparser → compiler chain), walks the resulting
// proto message graph, and emits valid SVG markup.
//
// Pipeline (mirrors proto-css gen, NOT genproto — no scalarize/prune/strip):
//
//	lang/*.ebnf  → strip (* *) comments → WrapString → ParseEBNF
//	             → GrammarToAST → Compile(Options{OnMessage, OnField})
//	             → FileDescriptorProto  → (reflection) → SVG markup
//
// The compile here deliberately does NOT run CollapseCommaList, NameSequence,
// scalarizeLeaves, pruneUnreachable, emptyKeywordRules, or StripKeywords. It
// keeps the markup terminals as keyword-message fields (captured in `kw`) so the
// renderer can re-emit `<rect`, `="`, `>`, `</rect>`, etc. directly, and it
// leaves the leaf *Type messages in place (intercepted by reps).
//
// The defining change vs proto-css: REPEATED fields are RENDERED (children and
// attributes ARE repetitions in SVG), not dropped. See render.go.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/accretional/gluon/v2/compiler"
	metaparser "github.com/accretional/gluon/v2/metaparser"
	pb "github.com/accretional/gluon/v2/pb"
	"google.golang.org/protobuf/types/descriptorpb"
)

// ebnfComment matches an EBNF (* ... *) comment (multi-line). gluon v2's
// ParseEBNF silently empties a rule body if a comment appears inside the rule's
// RHS, so we scrub comments before parsing (same workaround as genproto).
var ebnfComment = regexp.MustCompile(`(?s)\(\*.*?\*\)`)

// rootEBNF is concatenated first so SvgDocument is the conceptual root.
const rootEBNF = "svg.ebnf"

func main() {
	langDir := flag.String("lang", "lang", "directory of EBNF files")
	out := flag.String("out", "chrome-testing/generated", "output directory")
	debugTag := flag.String("debug", "", "dump in-memory structure for an element open-tag (e.g. \"<rect\") or \"ALL\"")
	flag.Parse()

	byFQN, kw, optional := compileGrammar(*langDir)
	r := newRenderer(byFQN, kw, optional)

	if *debugTag != "" {
		if *debugTag == "ALL" {
			dumpAllElements(byFQN, kw)
		} else if strings.HasPrefix(*debugTag, ".svg.") {
			dumpMsg(byFQN, kw, *debugTag, 0, map[string]int{})
		} else {
			dumpElement(byFQN, kw, *debugTag)
		}
		return
	}

	if err := os.MkdirAll(*out, 0o755); err != nil {
		fail("mkdir %s: %v", *out, err)
	}

	// (a) Whole-document render from the root production.
	doc := r.RenderDocument()
	docPath := filepath.Join(*out, "sample-document.svg")
	writeFile(docPath, doc)

	// (b) A minimal SVG containing one <rect> rendered with several attributes.
	rectFQN := findElementFQN(byFQN, "<rect")
	if rectFQN == "" {
		fail("could not locate the SVGRectElement message in the compiled grammar")
	}
	rect := r.RenderElementWithAttrs(rectFQN, 5)
	rectPath := filepath.Join(*out, "sample-rect.svg")
	writeFile(rectPath, rect)

	fmt.Printf("=== %s ===\n%s\n\n", rectPath, rect)
	fmt.Printf("=== %s (%d bytes) ===\n", docPath, len(doc))
	printHead(doc, 50)

	// Well-formedness check over both outputs.
	okDoc, msgDoc := checkWellFormed(doc)
	okRect, msgRect := checkWellFormed(rect)
	fmt.Printf("\nwell-formedness: sample-rect.svg: %s\n", status(okRect, msgRect))
	fmt.Printf("well-formedness: sample-document.svg: %s\n", status(okDoc, msgDoc))
	if !okDoc || !okRect {
		os.Exit(1)
	}

	// (c) The all-value-paths gallery: enumerate every element's every attribute's
	// every value, inject each into the element's blueprint, and emit per-element
	// gallery pages + index + values.json + manifest.tsv.
	runGalleryPass(byFQN, kw, r)
}

// compileGrammar reads every lang/*.ebnf (svg.ebnf first, then the rest sorted),
// strips (* *) comments, and runs WrapString → ParseEBNF → GrammarToAST →
// Compile with hooks that capture keyword literals and EBNF-optional fields. It
// returns the message graph indexed by FQN, the keyword-literal map, and the
// optional-field set.
func compileGrammar(langDir string) (map[string]*descriptorpb.DescriptorProto, map[string]string, map[string]bool) {
	files := collectEBNFFiles(langDir)
	var sb strings.Builder
	for _, name := range files {
		data, err := os.ReadFile(filepath.Join(langDir, name))
		if err != nil {
			fail("read %s: %v", name, err)
		}
		sb.Write(data)
		sb.WriteByte('\n')
	}
	src := ebnfComment.ReplaceAllString(sb.String(), " ")

	doc := metaparser.WrapString(src)
	doc.Name = "svg"

	gd, err := metaparser.ParseEBNF(doc)
	if err != nil {
		fail("ParseEBNF: %v", err)
	}
	ast, err := compiler.GrammarToAST(gd)
	if err != nil {
		fail("GrammarToAST: %v", err)
	}
	ast.Language = "svg"

	kw := map[string]string{}
	optional := map[string]bool{}
	fdp, err := compiler.Compile(ast, compiler.Options{
		Package: "svg",
		OnMessage: func(fqn string, node *pb.ASTNode) {
			if node.GetKind() == compiler.KindTerminal {
				kw[fqn] = node.GetValue()
			}
		},
		OnField: func(parentFQN, fieldName string, node *pb.ASTNode) {
			if node.GetKind() == compiler.KindOptional {
				optional[parentFQN+"/"+fieldName] = true
			}
		},
	})
	if err != nil {
		fail("Compile: %v", err)
	}

	// Index every message (top-level + nested) by fully-qualified name.
	byFQN := map[string]*descriptorpb.DescriptorProto{}
	var index func(m *descriptorpb.DescriptorProto, fqn string)
	index = func(m *descriptorpb.DescriptorProto, fqn string) {
		byFQN[fqn] = m
		for _, n := range m.GetNestedType() {
			index(n, fqn+"."+n.GetName())
		}
	}
	for _, m := range fdp.GetMessageType() {
		index(m, ".svg."+m.GetName())
	}
	globalKW = kw
	return byFQN, kw, optional
}

// collectEBNFFiles returns svg.ebnf first, then every other *.ebnf sorted.
func collectEBNFFiles(langDir string) []string {
	entries, err := os.ReadDir(langDir)
	if err != nil {
		fail("read dir %s: %v", langDir, err)
	}
	var rest []string
	haveRoot := false
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".ebnf" {
			continue
		}
		if e.Name() == rootEBNF {
			haveRoot = true
			continue
		}
		rest = append(rest, e.Name())
	}
	sort.Strings(rest)
	if !haveRoot {
		fail("root grammar %s not found in %s", rootEBNF, langDir)
	}
	return append([]string{rootEBNF}, rest...)
}

// findElementFQN returns the FQN of the element message whose leading markup
// terminal is the given open tag (e.g. "<rect"). An element compiles to a
// sequence message whose first field references the open-tag keyword message, so
// we scan for the message whose first field's keyword literal equals openTag.
func findElementFQN(byFQN map[string]*descriptorpb.DescriptorProto, openTag string) string {
	for fqn, m := range byFQN {
		fields := m.GetField()
		if len(fields) == 0 {
			continue
		}
		if globalKW[fields[0].GetTypeName()] == openTag {
			return fqn
		}
	}
	return ""
}

// globalKW mirrors the keyword map so findElementFQN can resolve tag literals.
var globalKW map[string]string

// writeFile writes content to path (creating parents), or fails.
func writeFile(path, content string) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		fail("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fail("write %s: %v", path, err)
	}
}

func printHead(s string, n int) {
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = lines[:n]
	}
	for _, ln := range lines {
		fmt.Println(ln)
	}
}

func status(ok bool, msg string) string {
	if ok {
		return "PASS"
	}
	return "FAIL — " + msg
}

// tagRe matches an opening start-tag name (`<rect`), but not a close tag
// (`</rect>`) or processing instruction (`<?xml`).
var tagRe = regexp.MustCompile(`<([A-Za-z][A-Za-z0-9-]*)`)

// checkWellFormed scans the markup and asserts two invariants:
//  1. tag balance: every start tag `<name` has a matching `</name>` later, with
//     correct nesting (a stack — close tags must match the most recent open).
//  2. attribute quoting: every `="` is closed by a later `"` before the next
//     `="` or `>` (i.e. attribute values are balanced double-quoted strings).
//
// It returns (ok, message). On failure, message names the offending tag/quote.
func checkWellFormed(svg string) (bool, string) {
	// (1) Tag balance via an explicit stack.
	var stack []string
	i := 0
	n := len(svg)
	for i < n {
		c := svg[i]
		if c != '<' {
			i++
			continue
		}
		// Skip processing instructions / declarations: `<?...?>`, `<!...>`.
		if i+1 < n && (svg[i+1] == '?' || svg[i+1] == '!') {
			i++
			continue
		}
		// Close tag `</name>`.
		if i+1 < n && svg[i+1] == '/' {
			j := i + 2
			start := j
			for j < n && svg[j] != '>' {
				j++
			}
			name := svg[start:j]
			if len(stack) == 0 {
				return false, "close tag </" + name + "> with no open tag"
			}
			top := stack[len(stack)-1]
			if top != name {
				return false, "close tag </" + name + "> does not match open <" + top + ">"
			}
			stack = stack[:len(stack)-1]
			i = j + 1
			continue
		}
		// Start tag `<name ...>`; could be self-closing `.../>`.
		m := tagRe.FindStringSubmatchIndex(svg[i:])
		if m == nil {
			i++
			continue
		}
		name := svg[i+m[2] : i+m[3]]
		// Find the end of this tag (`>`), tracking whether it self-closes.
		j := i + m[1]
		selfClose := false
		for j < n && svg[j] != '>' {
			if svg[j] == '/' && j+1 < n && svg[j+1] == '>' {
				selfClose = true
			}
			j++
		}
		if !selfClose {
			stack = append(stack, name)
		}
		i = j + 1
	}
	if len(stack) > 0 {
		return false, "unclosed start tag <" + stack[len(stack)-1] + ">"
	}

	// (2) Attribute quoting: every `="` must be closed by a `"` before the
	// enclosing tag ends. Walk attribute-opening tokens.
	for k := 0; k+1 < len(svg); k++ {
		if svg[k] == '=' && svg[k+1] == '"' {
			// find the closing quote
			closed := false
			for q := k + 2; q < len(svg); q++ {
				if svg[q] == '"' {
					closed = true
					break
				}
				// a `>` or `<` before the close quote means an unbalanced value
				if svg[q] == '<' {
					break
				}
			}
			if !closed {
				ctx := svg[k:min(k+24, len(svg))]
				return false, `unclosed attribute quote near ` + ctx
			}
		}
	}
	return true, ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "gen: "+format+"\n", args...)
	os.Exit(1)
}
