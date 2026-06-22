// Command genproto compiles the SVG EBNF grammar into a protobuf schema and the
// grammar-derived tables the reflection renderer/parser consult.
//
// Pipeline (all gluon v2, grammar-agnostic except scalarizeLeaves):
//
//	lang/*.ebnf
//	  → WrapString → ParseEBNF            → GrammarDescriptor
//	  → GrammarToAST                       → ASTDescriptor
//	  → CollapseCommaList                  (X (SEP X)* → repeated + separator)
//	  → NameSequence                       (name anon sequences after leading keywords)
//	  → scalarizeLeaves                    (atomic leaf types → string value=1)
//	  → Compile (prefix pass, OnMessage/OnField collects MessagePrefix+FieldSeparator)
//	  → StripKeywords                      (leading keyword terminals → prefix map)
//	  → Compile                            → FileDescriptorProto
//
// Outputs:
//
//	proto/svg.proto              bundled .proto (human-readable)
//	proto/svg.fdset              serialized FileDescriptorSet
//	proto/pb/svg/prefix_map.go   MessagePrefix: FQN → leading terminal tokens
//	proto/pb/svg/separator_map.go FieldSeparator: parentFQN.field → list separator
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/accretional/gluon/v2/compiler"
	metaparser "github.com/accretional/gluon/v2/metaparser"
	pb "github.com/accretional/gluon/v2/pb"
	"github.com/accretional/merge/descriptor"
)

// ebnfComment matches an EBNF (* ... *) comment, including multi-line, used to
// scrub comments out of the parser input (see the note in main).
var ebnfComment = regexp.MustCompile(`(?s)\(\*.*?\*\)`)

// rootEBNF is the grammar file that must be concatenated FIRST so its root
// production (SvgDocument) is the conceptual root of the schema.
const rootEBNF = "svg.ebnf"

// collectEBNFFiles returns the EBNF source files to concatenate, in order:
// svg.ebnf first (so SvgDocument is the root), then every other *.ebnf in the
// directory sorted lexically. This is tolerant of however many module files
// currently exist (today: svg.ebnf, datatype.ebnf) so the pipeline keeps
// working as more modules land.
func collectEBNFFiles(langDir string) ([]string, error) {
	entries, err := os.ReadDir(langDir)
	if err != nil {
		return nil, err
	}
	var rest []string
	haveRoot := false
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) != ".ebnf" {
			continue
		}
		if name == rootEBNF {
			haveRoot = true
			continue
		}
		rest = append(rest, name)
	}
	sort.Strings(rest)
	if !haveRoot {
		return nil, fmt.Errorf("root grammar %s not found in %s", rootEBNF, langDir)
	}
	return append([]string{rootEBNF}, rest...), nil
}

func main() {
	langDir := flag.String("lang", "lang", "directory of EBNF files")
	bundledOut := flag.String("bundled", "proto/svg.proto", "bundled .proto output")
	fdsetOut := flag.String("fdset", "proto/svg.fdset", "output FileDescriptorSet binary")
	prefixMapOut := flag.String("prefix-map", "proto/pb/svg/prefix_map.go", "generated MessagePrefix map")
	separatorMapOut := flag.String("separator-map", "proto/pb/svg/separator_map.go", "generated FieldSeparator map")
	pkgName := flag.String("package", "svg", "proto package name")
	goPkg := flag.String("go-package", "github.com/accretional/proto-svg/proto/pb/svg;svgpb", "go_package option")
	flag.Parse()

	// 1. Read + concatenate the grammar.
	ebnfFiles, err := collectEBNFFiles(*langDir)
	if err != nil {
		log.Fatalf("collect ebnf files: %v", err)
	}
	var sb strings.Builder
	for _, name := range ebnfFiles {
		data, err := os.ReadFile(filepath.Join(*langDir, name))
		if err != nil {
			log.Fatalf("read %s: %v", name, err)
		}
		sb.Write(data)
		sb.WriteByte('\n')
	}

	// Strip EBNF comments before parsing. gluon v2's ParseEBNF has a bug
	// where a comment anywhere inside a rule's right-hand side (leading or
	// interspersed between alternation arms) silently empties the entire rule
	// body. The SVG grammar comments heavily inside alternations and value
	// rules, so feeding the raw source would lose those rules and orphan
	// everything they reference. Each comment is replaced by a space so it
	// can't fuse adjacent tokens. The grammar files keep their comments; only
	// the parser input is cleaned.
	src := ebnfComment.ReplaceAllString(sb.String(), " ")

	doc := metaparser.WrapString(src)
	doc.Name = "svg"

	gd, err := metaparser.ParseEBNF(doc)
	if err != nil {
		log.Fatalf("ParseEBNF: %v", err)
	}

	ast, err := compiler.GrammarToAST(gd)
	if err != nil {
		log.Fatalf("GrammarToAST: %v", err)
	}
	ast.Language = *pkgName

	// 2. AST transforms.
	ast.Root = compiler.CollapseCommaList(ast.Root)
	ast.Root = compiler.NameSequence(ast.Root)
	ast.Root = scalarizeLeaves(ast.Root)
	var prunedRules []string
	ast.Root, prunedRules = pruneUnreachable(ast.Root, "SvgDocument")
	fmt.Printf("pruned %d unreachable rules\n", len(prunedRules))

	// 3. Prefix pass: collect leading-terminal prefixes and list separators
	// before StripKeywords discards the leading terminals.
	prefixes := map[string][]string{}
	separators := map[string]string{}
	if _, err := compiler.Compile(ast, compiler.Options{
		Package: *pkgName,
		OnMessage: func(fqn string, node *pb.ASTNode) {
			if node.GetKind() == compiler.KindTerminal {
				prefixes[fqn] = []string{node.GetValue()}
				return
			}
			// Pure-keyword rule (row = "row", colon_symbol = ":"): record its
			// literal so emptyKeywordRules can blank the body and the renderer
			// still re-emits the keyword from the prefix map.
			if lit, ok := bodyTerminal(node); ok {
				prefixes[fqn] = []string{lit}
				return
			}
			var kids []*pb.ASTNode
			switch node.GetKind() {
			case compiler.KindSequence:
				kids = node.GetChildren()
			case compiler.KindRule:
				body := node.GetChildren()
				if len(body) == 1 && body[0].GetKind() == compiler.KindSequence {
					kids = body[0].GetChildren()
				}
			}
			var toks []string
			for _, c := range kids {
				if c.GetKind() != compiler.KindTerminal {
					break
				}
				toks = append(toks, c.GetValue())
			}
			if len(toks) > 0 {
				prefixes[fqn] = toks
			}
		},
		OnField: func(parent, name string, node *pb.ASTNode) {
			if node.GetKind() == compiler.KindRepetition && node.GetValue() != "" {
				separators[parent+"."+name] = node.GetValue()
			}
		},
	}); err != nil {
		log.Fatalf("compiler.Compile (prefix pass): %v", err)
	}

	// 4. Blank pure-keyword rule bodies → empty marker messages (literals now
	// captured in prefixes), strip leading keyword terminals, then compile.
	ast.Root = emptyKeywordRules(ast.Root)
	ast.Root = compiler.StripKeywords(ast.Root)

	fdp, err := compiler.Compile(ast, compiler.Options{
		Package:   *pkgName,
		GoPackage: *goPkg,
		FileName:  filepath.Base(*bundledOut),
	})
	if err != nil {
		log.Fatalf("compiler.Compile: %v", err)
	}
	dups := dedupeMessages(fdp)
	dangling := dropDanglingFields(fdp)
	emptyOneofs := pruneEmptyOneofs(fdp)
	renamed := uniquifyFields(fdp)
	fmt.Printf("compiled %d messages from %d rules\n", len(fdp.GetMessageType()), len(gd.GetRules()))
	if len(dups) > 0 {
		fmt.Printf("note: deduped %d colliding message name(s): %v\n", len(dups), dups)
	}
	if len(dangling) > 0 {
		fmt.Printf("note: dropped %d field(s) referencing undefined grammar rules (dangling)\n", len(dangling))
	}
	if len(emptyOneofs) > 0 {
		fmt.Printf("note: dropped %d empty oneof(s) left by dangling-field removal: %v\n", len(emptyOneofs), emptyOneofs)
	}
	if renamed > 0 {
		fmt.Printf("note: uniquified %d duplicate field name(s) (repeated inline terminals)\n", renamed)
	}
	if deps := fdp.GetDependency(); len(deps) > 0 {
		fmt.Printf("note: svg.proto imports %v\n", deps)
	}

	// Write the audit reports (pruned rules + dangling refs) for inspection.
	writeReport("proto/GENPROTO_PRUNED.txt", "rules unreachable from SvgDocument (pruned)", prunedRules)
	writeReport("proto/GENPROTO_DANGLING.txt", "fields dropped — type references an undefined grammar rule", dangling)

	// 5. Write outputs.
	if err := os.MkdirAll(filepath.Dir(*bundledOut), 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(*prefixMapOut), 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}

	set := &descriptorpb.FileDescriptorSet{File: []*descriptorpb.FileDescriptorProto{fdp}}
	blob, err := proto.Marshal(set)
	if err != nil {
		log.Fatalf("marshal fdset: %v", err)
	}
	if err := os.WriteFile(*fdsetOut, blob, 0o644); err != nil {
		log.Fatalf("write %s: %v", *fdsetOut, err)
	}
	fmt.Printf("wrote %s (%d bytes)\n", *fdsetOut, len(blob))

	protoSrc, err := descriptor.ToString(fdp)
	if err != nil {
		log.Fatalf("descriptor.ToString: %v", err)
	}
	if err := os.WriteFile(*bundledOut, []byte(protoSrc), 0o644); err != nil {
		log.Fatalf("write %s: %v", *bundledOut, err)
	}
	fmt.Printf("wrote %s\n", *bundledOut)

	if err := os.WriteFile(*prefixMapOut, []byte(formatPrefixMap(prefixes)), 0o644); err != nil {
		log.Fatalf("write %s: %v", *prefixMapOut, err)
	}
	fmt.Printf("wrote %s (%d entries)\n", *prefixMapOut, len(prefixes))

	if err := os.WriteFile(*separatorMapOut, []byte(formatSeparatorMap(separators)), 0o644); err != nil {
		log.Fatalf("write %s: %v", *separatorMapOut, err)
	}
	fmt.Printf("wrote %s (%d entries)\n", *separatorMapOut, len(separators))
}

// dedupeMessages removes duplicate top-level messages by name, keeping the
// first occurrence, and returns the dropped names. The snake_case leaf atoms
// and their PascalCase wrappers (e.g. number_type / NumberType) both lower to
// the same message name; after scalarization they are identical scalar
// messages, so dropping the later one is loss-free. Any non-scalar collision
// would surface in the returned list for inspection.
func dedupeMessages(fdp *descriptorpb.FileDescriptorProto) []string {
	seen := map[string]bool{}
	var kept []*descriptorpb.DescriptorProto
	var dropped []string
	for _, m := range fdp.GetMessageType() {
		name := m.GetName()
		if seen[name] {
			dropped = append(dropped, name)
			continue
		}
		seen[name] = true
		kept = append(kept, m)
	}
	fdp.MessageType = kept
	sort.Strings(dropped)
	return dropped
}

// dropDanglingFields removes any field whose message type references a
// message that was never emitted, recursing into nested messages. These come
// from dangling nonterminal references in the grammar (a rule that references
// a name no rule defines, e.g. StrokeLinejoinProp → miter_clip, where
// `miter_clip` is undefined). Without this, descriptor serialization and
// protoc both fail to resolve the type. Returns "Parent.field -> .svg.Type"
// descriptions for the audit report.
func dropDanglingFields(fdp *descriptorpb.FileDescriptorProto) []string {
	defined := map[string]bool{}
	var indexNames func(m *descriptorpb.DescriptorProto, fqn string)
	indexNames = func(m *descriptorpb.DescriptorProto, fqn string) {
		defined[fqn] = true
		for _, n := range m.GetNestedType() {
			indexNames(n, fqn+"."+n.GetName())
		}
	}
	for _, m := range fdp.GetMessageType() {
		indexNames(m, ".svg."+m.GetName())
	}

	var dropped []string
	var clean func(m *descriptorpb.DescriptorProto, fqn string)
	clean = func(m *descriptorpb.DescriptorProto, fqn string) {
		var kept []*descriptorpb.FieldDescriptorProto
		for _, f := range m.GetField() {
			if f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE && !defined[f.GetTypeName()] {
				dropped = append(dropped, fmt.Sprintf("%s.%s -> %s", fqn, f.GetName(), f.GetTypeName()))
				continue
			}
			kept = append(kept, f)
		}
		m.Field = kept
		for _, n := range m.GetNestedType() {
			clean(n, fqn+"."+n.GetName())
		}
	}
	for _, m := range fdp.GetMessageType() {
		clean(m, ".svg."+m.GetName())
	}
	sort.Strings(dropped)
	return dropped
}

// pruneEmptyOneofs removes oneof declarations that no surviving field
// references and re-indexes the remaining fields' OneofIndex values, recursing
// into nested messages. dropDanglingFields can strip every variant of an
// alternation oneof (e.g. a content-category union like AnimationElement whose
// members all live in not-yet-authored module files), leaving an empty
// OneofDecl that the descriptor printer/protoc reject ("oneof must contain at
// least one field declaration"). When the grammar is complete no oneof is ever
// emptied, so this is a no-op then; during the partial-grammar foundation stage
// it keeps the pipeline compiling. Returns "Parent.oneofName" descriptions.
func pruneEmptyOneofs(fdp *descriptorpb.FileDescriptorProto) []string {
	var dropped []string
	var walk func(m *descriptorpb.DescriptorProto, fqn string)
	walk = func(m *descriptorpb.DescriptorProto, fqn string) {
		// Count surviving field references per oneof index.
		used := make([]bool, len(m.GetOneofDecl()))
		for _, f := range m.GetField() {
			if f.OneofIndex != nil {
				if i := int(f.GetOneofIndex()); i >= 0 && i < len(used) {
					used[i] = true
				}
			}
		}
		// Build the kept oneof list and an old→new index remap.
		remap := make([]int32, len(m.GetOneofDecl()))
		var keptOneofs []*descriptorpb.OneofDescriptorProto
		for i, od := range m.GetOneofDecl() {
			if used[i] {
				remap[i] = int32(len(keptOneofs))
				keptOneofs = append(keptOneofs, od)
			} else {
				remap[i] = -1
				dropped = append(dropped, fmt.Sprintf("%s.%s", fqn, od.GetName()))
			}
		}
		if len(keptOneofs) != len(m.GetOneofDecl()) {
			m.OneofDecl = keptOneofs
			for _, f := range m.GetField() {
				if f.OneofIndex != nil {
					f.OneofIndex = proto.Int32(remap[f.GetOneofIndex()])
				}
			}
		}
		for _, n := range m.GetNestedType() {
			walk(n, fqn+"."+n.GetName())
		}
	}
	for _, m := range fdp.GetMessageType() {
		walk(m, ".svg."+m.GetName())
	}
	sort.Strings(dropped)
	return dropped
}

// uniquifyFields renames duplicate field names within each message (and nested
// messages), keeping the first and suffixing later ones (_2, _3, …). Inlining
// keyword nonterminals into inline terminals, and the compiler names a terminal
// field after its literal without deduplicating; a message with two of the same
// terminal ends up with duplicate field names that descriptor serialization and
// protoc both reject. The renderer/parser address fields by position, so
// renaming is loss-free. Returns the count of renamed fields.
func uniquifyFields(fdp *descriptorpb.FileDescriptorProto) int {
	n := 0
	var walk func(m *descriptorpb.DescriptorProto)
	walk = func(m *descriptorpb.DescriptorProto) {
		seen := map[string]int{}
		for _, f := range m.GetField() {
			base := f.GetName()
			if seen[base] == 0 {
				seen[base] = 1
				continue
			}
			seen[base]++
			f.Name = proto.String(fmt.Sprintf("%s_%d", base, seen[base]))
			n++
		}
		for _, nm := range m.GetNestedType() {
			walk(nm)
		}
	}
	for _, m := range fdp.GetMessageType() {
		walk(m)
	}
	return n
}

// writeReport writes a sorted, headed list to path for human inspection.
func writeReport(path, title string, items []string) {
	sorted := append([]string(nil), items...)
	sort.Strings(sorted)
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n# count: %d\n\n", title, len(sorted))
	for _, it := range sorted {
		b.WriteString(it)
		b.WriteByte('\n')
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		log.Printf("write %s: %v", path, err)
		return
	}
	fmt.Printf("wrote %s (%d entries)\n", path, len(sorted))
}

func formatPrefixMap(prefixes map[string][]string) string {
	keys := make([]string, 0, len(prefixes))
	for k := range prefixes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("// Code generated by genproto. DO NOT EDIT.\n\n")
	b.WriteString("package svgpb\n\n")
	b.WriteString("// MessagePrefix maps a fully-qualified proto message name to the leading\n")
	b.WriteString("// terminal tokens that StripKeywords removed from its schema. The renderer\n")
	b.WriteString("// emits these tokens before walking the message's fields; the parser\n")
	b.WriteString("// consumes them first.\n")
	b.WriteString("var MessagePrefix = map[string][]string{\n")
	for _, k := range keys {
		b.WriteString("\t")
		b.WriteString(strconv.Quote(k))
		b.WriteString(": {")
		for i, tok := range prefixes[k] {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(strconv.Quote(tok))
		}
		b.WriteString("},\n")
	}
	b.WriteString("}\n")
	return b.String()
}

func formatSeparatorMap(separators map[string]string) string {
	keys := make([]string, 0, len(separators))
	for k := range separators {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString("// Code generated by genproto. DO NOT EDIT.\n\n")
	b.WriteString("package svgpb\n\n")
	b.WriteString("// FieldSeparator maps a \"parentFQN.fieldName\" key to the literal that\n")
	b.WriteString("// CollapseCommaList dropped when it rewrote \"X (SEP X)*\" into a repeated\n")
	b.WriteString("// field. The renderer interleaves it between elements; the parser splits on it.\n")
	b.WriteString("var FieldSeparator = map[string]string{\n")
	for _, k := range keys {
		b.WriteString("\t")
		b.WriteString(strconv.Quote(k))
		b.WriteString(": ")
		b.WriteString(strconv.Quote(separators[k]))
		b.WriteString(",\n")
	}
	b.WriteString("}\n")
	return b.String()
}
