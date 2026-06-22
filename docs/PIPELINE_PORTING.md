# proto-css → proto-svg Pipeline Porting Guide

This is the sole, exhaustive reference for replicating the proto-css
EBNF → protobuf → example-generation pipeline in `proto-svg`. Every symbol
name, signature, file path, and behavior below was read directly from the
reference repos and verified by building/running. SVG-specific porting risks
are called out in §6.

Reference repos (read-only; do not modify):
- `/Volumes/wd_office_3/Projects/proto-css`   — the pipeline being ported
- `/Volumes/wd_office_3/Projects/gluon`        — the grammar/compiler library
- `/Volumes/wd_office_3/Projects/proto-merge`  — `descriptor.ToString` (proto printer)

---

## 0. The v2 import path puzzle (resolved)

proto-css imports `github.com/accretional/gluon/v2/{compiler,metaparser,pb}`.
Its `go.mod` has:

```
module github.com/accretional/proto-css
go 1.25.5
replace github.com/accretional/gluon => ../gluon
replace github.com/accretional/merge => ../proto-merge
require (
    github.com/accretional/gluon v0.0.0-20260430002741-82b4ddbe67a4
    github.com/accretional/merge v0.0.0-00010101000000-000000000000
    ...
)
```

The local gluon module declares `module github.com/accretional/gluon` (NO
`/v2` in the module path, and there is **no separate go.mod under
`../gluon/v2/`**). So how does `.../gluon/v2/compiler` resolve?

**Answer: `/v2` is just a *subdirectory* of the single gluon module, not a
semantic-major-version submodule.** Because the module path is plain
`github.com/accretional/gluon`, Go treats `github.com/accretional/gluon/v2`
as the package directory `v2/` inside that module. Verified:

```
$ go list -f '{{.Dir}}' github.com/accretional/gluon/v2/compiler
/Volumes/wd_office_3/Projects/gluon/v2/compiler
$ go list -m all | grep gluon
github.com/accretional/gluon v0.0.0-20260430002741-82b4ddbe67a4 => ../gluon
$ go build ./...      # succeeds with no output
```

This works **only** because the module path lacks the `/v2` suffix (a real
`/v2` major-version module would need `module github.com/accretional/gluon/v2`
in its own go.mod, and Go's import-comment rules would reject importing it as a
subdir). It is effectively a directory naming convention ("v2 = the new
rewrite") that Go happens to allow. `go doc github.com/accretional/gluon/v2/compiler`
works and lists the full API.

For the SVG port: keep the identical layout — proto-svg's go.mod uses the same
two `replace` directives, and the import paths are byte-for-byte the same
(`github.com/accretional/gluon/v2/compiler`, etc.).

---

## 1. Gluon v2 public API (the pipeline surface)

All paths under `/Volumes/wd_office_3/Projects/gluon/v2/`. Package import
paths in parentheses.

### 1.1 `metaparser` package (`github.com/accretional/gluon/v2/metaparser`)

Source: `v2/metaparser/{metaparser.go,ebnf.go,compile.go,...}`.

```go
// metaparser.go
func WrapString(s string) *pb.DocumentDescriptor
// Lifts a Go string into a single-chunk DocumentDescriptor. Never errors.
// The returned doc has no Name/Uri — set doc.Name yourself if you want one.

// ebnf.go
func ParseEBNF(doc *pb.DocumentDescriptor) (*pb.GrammarDescriptor, error)
// Concatenates the doc's text chunks, runs the v1 lexkit ISO-14977 EBNF
// parser, converts to the v2 flat-Production grammar shape.

func TextOf(doc *pb.DocumentDescriptor) string
// Concatenates a DocumentDescriptor's text chunks back to a Go string.

func EBNFLexV2() *pb.LexDescriptor   // the canonical ISO-14977 lex table
```

Typical use (from `genproto/main.go` and `gen/main.go`):

```go
doc := metaparser.WrapString(src)
doc.Name = "css"                    // → "svg" in the port
gd, err := metaparser.ParseEBNF(doc)   // *pb.GrammarDescriptor
```

> **Known parser bug (must replicate the workaround):** gluon v2's `ParseEBNF`
> silently empties a rule body if an EBNF `(* ... *)` comment appears anywhere
> inside the rule's right-hand side. proto-css strips comments before parsing
> (see §2, `ebnfComment`). The SVG grammar must do the same.

### 1.2 `compiler` package (`github.com/accretional/gluon/v2/compiler`)

Source: `v2/compiler/{compiler.go,kinds.go,collapse.go,name_sequence.go,strip_keywords.go,grammar2ast.go,names.go}`.
Package doc (`go doc .../compiler`) defines the "schema-AST" kind conventions.

**Grammar → AST:**
```go
func GrammarToAST(gd *pb.GrammarDescriptor) (*pb.ASTDescriptor, error)
// Regroups the flat Production lists (with CONCATENATION/ALTERNATION delimiters)
// into sequence/alternation subtrees by EBNF precedence (alternation < concat).
// Result root is a `file` node whose children are `rule` nodes in source order.
// ASTDescriptor.Language is taken from gd.Name.
```

**AST transforms (all pure; return a NEW tree, input not mutated):**
```go
func CollapseCommaList(root *pb.ASTNode) *pb.ASTNode
// Rewrites "X (SEP X)*" → bare repetition of X, dropping the SEP terminal and
// stashing the separator literal in the repetition node's Value. Applied
// bottom-up recursively. SEP is any one-token terminal ("," ";" "."). No-op if
// the pattern is absent.

func NameSequence(root *pb.ASTNode) *pb.ASTNode
// Annotates each sequence node whose leading children are terminals with a
// PascalCase wrapper-message name (concatenation of those terminals'
// identifierized values), stored in node.Value. e.g. seq["ORDER","BY",x] → "OrderBy".
// The compiler reads Value as the preferred nested-message stem.

func StripKeywords(root *pb.ASTNode) *pb.ASTNode
// Removes the leading run of terminal children from every sequence node (their
// identity is now carried by the rule/wrapper name). seq(IF,EXISTS) → empty
// sequence (a pure marker message).
```

**Compile:**
```go
func Compile(ast *pb.ASTDescriptor, opts Options) (*descriptorpb.FileDescriptorProto, error)
// One proto message per `rule` node. Keyword terminals collected into a
// deduplicated set of empty messages appended after the rule messages.

type Options struct {
    Package   string  // proto package; defaults to ast.Language (sanitized) or "lang"
    GoPackage string  // go_package file option; empty = omit
    FileName  string  // FileDescriptorProto.name; defaults to "<package>.proto"
    OnMessage func(fqn string, node *pb.ASTNode)
    OnField   func(parentFQN, fieldName string, node *pb.ASTNode)
}
```

`OnMessage` fires per emitted message as soon as its fully-qualified name is
known; `node` is the AST node that produced it (rule node for top-level
messages, sequence/alternation node for nested wrappers, terminal node for
keyword messages). `OnField` fires per emitted field; `node` is the AST node
**before** repetition/optional/group peeling — so the caller can read wrapper
metadata (e.g. the separator stashed on a `repetition` node by
`CollapseCommaList`, or detect a `KindOptional` wrapper). Range lowering emits
two fields and fires `OnField` twice.

**Kind constants** (`v2/compiler/kinds.go`) — string-typed:
```go
const (
    KindFile        = "file"        // root; children are rules
    KindRule        = "rule"        // Value = rule name; 1 body child
    KindSequence    = "sequence"    // ordered concatenation → fields
    KindAlternation = "alternation" // choice → proto3 oneof
    KindOptional    = "optional"    // [ body ] — 1 child
    KindRepetition  = "repetition"  // { body } — 1 child → LABEL_REPEATED
    KindGroup       = "group"       // ( body ) — transparent, 1 child
    KindTerminal    = "terminal"    // quoted literal; Value = literal text
    KindNonterminal = "nonterminal" // ref to another rule; Value = name
    KindRange       = "range"       // 2 children: range_lower, range_upper
    KindRangeLower  = "range_lower"
    KindRangeUpper  = "range_upper"
    KindScalar      = "scalar"      // → proto3 string; Value = field name (default "value")
)
```

**Compile lowering rules (`compiler.go`, the part the port depends on most):**
- `rule` body that is a `sequence` → each item becomes a field.
- `rule` body that is an `alternation` → a oneof named `value`.
- `repetition` peels to `LABEL_REPEATED`; `optional` and `group` peel
  transparently (both stay `LABEL_OPTIONAL`). Peel loop is in `emitField`.
- `terminal` → reference to a deduped empty `*Keyword` message; field named
  `<ident>_keyword`.
- `scalar` → `string` field (`appendStringField`).
- `nonterminal` → message field referencing `.<pkg>.<PascalCase(name)>`.
- nested `sequence`/`alternation` → promoted to a **nested message**
  (`nestedFromSequence` / `nestedFromAlternation`), named from `node.Value`
  (set by `NameSequence`) or `Seq`/`Alt` + N.
- `range` → two `.unicode.UTF8` message fields and sets `usesUTF8` (pulls in
  `unicode/utf_8.proto` dependency). proto-css avoids this by scalarizing all
  range-bearing rules (§2).

### 1.3 `pb` package (`github.com/accretional/gluon/v2/pb`)

> **Naming note:** the task brief calls these `pb.GrammarDefinition`/`ASTNode`.
> The actual type is **`GrammarDescriptor`** (not `GrammarDefinition`). Verify
> against the real types below.

**`ASTNode`** (`v2/pb/ast.pb.go`):
```go
type ASTNode struct {
    Kind     string          // proto field 1
    Value    string          // proto field 2 (matched text / rule name / literal)
    Children []*ASTNode      // proto field 3 (source order)
    Location *SourceLocation // proto field 4
}
func (x *ASTNode) GetKind() string
func (x *ASTNode) GetValue() string
func (x *ASTNode) GetChildren() []*ASTNode
```

**`ASTDescriptor`** (`v2/pb/ast.pb.go`):
```go
type ASTDescriptor struct {
    Language string    // 1
    Version  string    // 2
    Root     *ASTNode  // 3
}
```

**`GrammarDescriptor`** (`v2/pb/grammar.pb.go`):
```go
type GrammarDescriptor struct {
    Name  string            // 1
    Lex   *LexDescriptor    // 2
    Rules []*RuleDescriptor // 3
}
type RuleDescriptor struct {
    Name        string        // 1
    Expressions []*Production // 2 (flat list; structure via delimiter/scoper positions)
}
type Production struct {
    Kind isProduction_Kind // oneof: Delimiter|Scoper|Terminal|Nonterminal|Range
}
// wrappers: Production_Terminal{Terminal string}, Production_Nonterminal{Nonterminal string},
//           Production_Range{Range *StringRange}, Production_Scoper{Scoper *ScopedProduction},
//           Production_Delimiter{Delimiter Delimiter}
```

**Enums** (`v2/pb/lex.pb.go`):
```go
type Scoper int32
const (
    Scoper_SCOPER_UNSPECIFIED Scoper = 0
    Scoper_OPTIONAL           Scoper = 1   // [ ]
    Scoper_REPETITION         Scoper = 2   // { }
    Scoper_GROUP              Scoper = 3   // ( )
    Scoper_TERMINAL           Scoper = 4   // " " ' '
    Scoper_COMMENT            Scoper = 5   // (* *)
)
type Delimiter int32
const (
    Delimiter_DELIMITER_UNSPECIFIED Delimiter = 0
    Delimiter_WHITESPACE            Delimiter = 1
    Delimiter_DEFINITION            Delimiter = 2   // =
    Delimiter_CONCATENATION         Delimiter = 3   // ,
    Delimiter_TERMINATION           Delimiter = 4   // ;
    Delimiter_ALTERNATION           Delimiter = 5   // |
)
```

The port only references `pb.ASTNode`, `pb.ASTDescriptor`, `pb.GrammarDescriptor`
directly; the Production/Scoper/Delimiter layer is internal to `GrammarToAST`
and you never touch it.

### 1.4 `merge/descriptor` package (`github.com/accretional/merge/descriptor`)

Source: `/Volumes/wd_office_3/Projects/proto-merge/descriptor/descriptor.go`.
```go
func ToString(fdp *descriptorpb.FileDescriptorProto) (string, error)
// Renders a FileDescriptorProto to human-readable .proto source. Internally:
// desc.CreateFileDescriptor(fdp) + protoprint.Printer{Compact:false}
// (jhump/protoreflect). Used by genproto to write proto/css.proto.
```

---

## 2. The `genproto` command

Path: `/Volumes/wd_office_3/Projects/proto-css/lang/cmd/genproto/{main.go,scalarize.go,prune.go,keywordize.go}`.

Compiles the CSS EBNF grammar into a protobuf schema plus the
grammar-derived lookup tables the renderer/parser consult.

### 2.1 Transform order (`main.go`)

Flags: `-lang` (dir, default `lang`), `-bundled` (`proto/css.proto`),
`-fdset` (`proto/css.fdset`), `-prefix-map` (`proto/pb/css/prefix_map.go`),
`-separator-map` (`proto/pb/css/separator_map.go`), `-package` (`css`),
`-go-package` (`github.com/accretional/proto-css/proto/pb/css;csspb`).

`ebnfFiles` (concatenation order — `css.ebnf` first so `CssStyleSheet` is the
root): `css.ebnf, primitive.ebnf, combinator.ebnf, datatype.ebnf,
functions.ebnf, pseudo-class.ebnf, pseudo-element.ebnf, selector.ebnf,
property.ebnf, atrule.ebnf`.

Pipeline (exact sequence):
```
1. Read+concat ebnfFiles → src
2. src = ebnfComment.ReplaceAllString(src, " ")   // strip (* … *) — parser bug workaround
   ebnfComment = regexp.MustCompile(`(?s)\(\*.*?\*\)`)
3. doc := metaparser.WrapString(src); doc.Name = "css"
4. gd  := metaparser.ParseEBNF(doc)
5. ast := compiler.GrammarToAST(gd); ast.Language = pkgName
6. ast.Root = compiler.CollapseCommaList(ast.Root)   // X (SEP X)* → repeated + separator
7. ast.Root = compiler.NameSequence(ast.Root)        // name anon sequences
8. ast.Root = scalarizeLeaves(ast.Root)              // LOCAL: leaf rules → string value=1
9. ast.Root, prunedRules = pruneUnreachable(ast.Root, "CssStyleSheet")
10. PREFIX PASS: compiler.Compile(ast, Options{Package, OnMessage, OnField})
       — discards the returned FDP; the hooks populate two maps:
         prefixes   map[fqn][]string   (MessagePrefix)
         separators map[parent.field]string (FieldSeparator)
11. ast.Root = emptyKeywordRules(ast.Root)   // pure-keyword rule bodies → empty marker
12. ast.Root = compiler.StripKeywords(ast.Root)
13. fdp := compiler.Compile(ast, Options{Package, GoPackage, FileName})  // FINAL
14. dedupeMessages(fdp); dropDanglingFields(fdp); uniquifyFields(fdp)
15. Write fdset, css.proto (via descriptor.ToString), prefix_map.go, separator_map.go,
    plus audit reports proto/GENPROTO_PRUNED.txt, proto/GENPROTO_DANGLING.txt
```

Verified run output:
```
pruned 230 unreachable rules
compiled 2859 messages from 2133 rules
note: deduped 8 colliding message name(s): [AngleType IntegerType LengthType …]
note: dropped 1 field(s) referencing undefined grammar rules (dangling)
note: uniquified 42 duplicate field name(s) (repeated inline terminals)
wrote proto/css.fdset (601587 bytes); wrote proto/css.proto
wrote prefix_map.go (3431 entries); wrote separator_map.go (113 entries)
```

### 2.2 `scalarizeLeaves` (`scalarize.go`)

Replaces the body of every **leaf rule** with a single `scalar` node so it
lowers to `message X { string value = 1; }`. A rule is a leaf when **either**:
- its normalized name (lowercased, underscores stripped — so `time_type` and
  `TimeType` collide) is in the curated `leafTypes` set, **or**
- its body contains a character `range` (`hasRange`) — collapsing ranges keeps
  the proto self-contained (avoids `.unicode.UTF8` + `unicode/utf_8.proto`).

`leafTypes` (CSS-specific; the port must rewrite this list for SVG): `number_type`,
`non_negative_number_type`, `integer_type`, `non_negative_integer_type`,
`positive_integer_type`, `negative_integer_type`, `percentage_type`,
`non_negative_percentage_type`, `length_type`, `non_negative_length_type`,
`time_type`, `non_negative_time_type`, `angle_type`, `flex_type`,
`non_negative_flex_type`, `frequency_type`, `resolution_type`, `ident_type`,
`custom_ident_type`, `dashed_ident_type`, `hex_color_type`, `string_type`,
`dimension_type`.

`norm(s) = strings.ToLower(strings.ReplaceAll(s, "_", ""))`. This is the only
**grammar-specific** step in an otherwise grammar-agnostic pipeline.

### 2.3 `pruneUnreachable` (`prune.go`)

```go
func pruneUnreachable(root *pb.ASTNode, roots ...string) (*pb.ASTNode, []string)
```
DFS from `roots` following `KindNonterminal` references (`collectRefs`); drops
every rule not reachable. Returns the pruned `file` node (surviving rules in
original order) and the dropped names. After scalarization the low-level
primitive ranges (digit/hex-digit char ranges) become unreachable, so pruning
shrinks ~3900 rules to the reachable core **and** removes the last
range-bearing rules — dropping the `unicode/utf_8.proto` import. For SVG the
single root is the document/root-element rule (CSS uses `"CssStyleSheet"`).

### 2.4 Keyword handling (`keywordize.go`)

```go
func bodyTerminal(rule *pb.ASTNode) (string, bool)
// True iff a rule body is effectively one terminal literal (e.g. row = "row",
// colon_symbol = ":"). Unwraps a single enclosing sequence/group.

func emptyKeywordRules(root *pb.ASTNode) *pb.ASTNode
// Blanks the body of every pure-keyword rule → empty marker message (message Row {}).
// The literal was already captured in the prefix map during the prefix pass, so the
// renderer re-emits it. Removes ~1700 redundant *Keyword wrapper messages.
// MUST run after the prefix pass and before the final Compile.
```

### 2.5 Prefix map & separator map (the load-bearing outputs)

Built in the **prefix pass** (step 10) via the `Compile` hooks:

```go
OnMessage: func(fqn string, node *pb.ASTNode) {
    // 1. bare terminal → record its literal
    if node.GetKind() == compiler.KindTerminal { prefixes[fqn] = []string{node.GetValue()}; return }
    // 2. pure-keyword rule → record its single literal (so emptyKeywordRules can blank it)
    if lit, ok := bodyTerminal(node); ok { prefixes[fqn] = []string{lit}; return }
    // 3. otherwise: collect the LEADING RUN of terminal children of the rule/sequence body
    //    (these are what StripKeywords will delete)
    var toks []string
    for _, c := range kids { if c.GetKind() != compiler.KindTerminal { break }; toks = append(toks, c.GetValue()) }
    if len(toks) > 0 { prefixes[fqn] = toks }
},
OnField: func(parent, name string, node *pb.ASTNode) {
    // CollapseCommaList stashed the list separator on the repetition node's Value
    if node.GetKind() == compiler.KindRepetition && node.GetValue() != "" {
        separators[parent+"."+name] = node.GetValue()
    }
},
```

- **`MessagePrefix` (`prefix_map.go`)**: `map[string][]string` — FQN →
  leading terminal tokens that `StripKeywords` removed. The renderer emits these
  before walking the message's fields; the parser consumes them first.
  Generated package `csspb`.
- **`FieldSeparator` (`separator_map.go`)**: `map[string]string` —
  `"parentFQN.fieldName"` → the list separator `CollapseCommaList` dropped. The
  renderer interleaves it; the parser splits on it. Generated package `csspb`.

Both are emitted as Go source by `formatPrefixMap` / `formatSeparatorMap`
(sorted keys, `// Code generated by genproto. DO NOT EDIT.` header).

### 2.6 Descriptor cleanup passes (`main.go`)

Run on the final FDP, all loss-free (renderer/parser address fields by
position):
- `dedupeMessages(fdp) []string` — drops duplicate top-level messages by name
  (the snake_case atom and PascalCase wrapper scalarize to identical messages).
- `dropDanglingFields(fdp) []string` — removes any message-typed field whose
  `TypeName` was never emitted (dangling nonterminal refs); recurses into nested
  messages. Without this, descriptor serialization and protoc fail to resolve.
- `uniquifyFields(fdp) int` — renames duplicate field names within a message
  (`name`, `name_2`, …); inline terminals can collide (e.g. the commas in
  `device-cmyk(c,m,y,k)`).

### 2.7 Outputs

| Output | Path | What |
|---|---|---|
| Bundled .proto | `proto/css.proto` | `descriptor.ToString(fdp)` |
| FileDescriptorSet | `proto/css.fdset` | `proto.Marshal(&FileDescriptorSet{File:[fdp]})` |
| Prefix map | `proto/pb/css/prefix_map.go` | `var MessagePrefix map[string][]string`, package `csspb` |
| Separator map | `proto/pb/css/separator_map.go` | `var FieldSeparator map[string]string`, package `csspb` |
| Audit | `proto/GENPROTO_PRUNED.txt`, `proto/GENPROTO_DANGLING.txt` | inspection reports |

---

## 3. The `gen` command (gallery data)

Path: `/Volumes/wd_office_3/Projects/proto-css/chrome-testing/cmd/gen/{main.go,render.go,reps.go,classify.go,emit.go}`.

Walks the **same** grammar (re-compiled in-process; it does **not** read the
fdset on disk — see note below) and renders every CSS property/value pair via
proto reflection over the compiled message graph, emitting the gallery data
files consumed by `chrome-testing/gallery`.

> **Important difference from genproto:** `gen` runs a *simpler* compile —
> `WrapString → ParseEBNF → GrammarToAST → Compile` with **no** CollapseCommaList,
> NameSequence, scalarizeLeaves, prune, StripKeywords, or keyword-emptying. It
> uses the `OnMessage`/`OnField` hooks to build two in-memory maps (`kw` =
> keyword literals, `optional` = EBNF-optional fields) instead of writing files.
> It does not load `css.fdset`; it keeps the FDP in memory. (The grammar walk is
> over `*descriptorpb.FileDescriptorProto`, indexed by FQN via reflection-style
> maps, **not** `protoreflect.MessageDescriptor`.)

### 3.1 `compileGrammar` (`main.go`)

```go
func compileGrammar(langDir string) (*descriptorpb.FileDescriptorProto, map[string]string, string)
```
- Same `ebnfFiles` order as genproto (but no comment stripping here — gen uses
  the raw source, since it only reads keyword literals/structure).
- `compiler.Compile` with:
  - `OnMessage`: if `node.GetKind()==KindTerminal` → `kw[fqn] = node.GetValue()`
  - `OnField`: if `node.GetKind()==KindOptional` → `optional[parentFQN+"/"+fieldName]=true`
- Stashes `optional` in package global `globalOptional`; returns `(fdp, kw, src)`.

### 3.2 Enumeration (`main.go`)

`enumerate(fdp, r, ruleText) []Property` walks the `Property` message's oneof
fields. For each `XExpr` field it resolves the kebab property name via
`propertyKeyword(exprMsg, r)` (renders the leading keyword nonterminal),
locates the `XProp` value rule (`.css.XProp`, or `findPropFQN` fallback that
prefers a `*Prop` field that is not `AllProp`), then `r.RenderP(propFQN)` to
produce values + provenance. A high-cap renderer (`maxVals=maxProduct=160`)
probes the true value count to flag truncation.

`Property` struct: `Name, Expr, Prop, Values []string, Syntax, Provenance,
Truncated, TrueCount, Assists []Assist, Description, Experimental, Nonstandard,
Deprecated, Warning`.

### 3.3 The Renderer algorithm (`render.go`) — the heart of the port

```go
type Renderer struct {
    byFQN    map[string]*descriptorpb.DescriptorProto // ".css.Foo" → message
    kw       map[string]string                        // FQN → keyword literal
    optional map[string]bool                          // "<msgFQN>/<field>" → EBNF-optional
    maxDepth, maxVals, maxProduct int                 // 16, 50, 50
    cur    *prov; argIdx int
}
func (r *Renderer) Render(fqn string) []string
func (r *Renderer) RenderP(fqn string) ([]string, prov)
func (r *Renderer) render(fqn string, depth int, seen map[string]int, argPos bool) []string
```

`render` dispatch order (this is the core logic — reproduce exactly):
1. **depth guard**: `depth > maxDepth` → `nil` (`prov.depthCut`).
2. **open-ended leaf** (`reps[simpleName(fqn)]` exists): return the sample set;
   set `prov.usedRep`/`addAssist`. In `argPos` (inside a function/arg list) it
   returns ONE position-cycled sample (`v[argIdx%len(v)]; argIdx++`) so channels
   vary; otherwise the full sample set.
3. **keyword message** (`kw[fqn]` exists): return `[]string{lit}` (empty `lit` → nil).
4. **unresolved** (`byFQN[fqn]==nil`): `nil` (e.g. `.unicode.UTF8`).
5. **cycle guard**: `seen[fqn] >= 2` → nil (`prov.cycleCut`); else `seen[fqn]++`/defer `--`.
6. **oneof / alternation** (`len(m.GetOneofDecl())>0`): union of every variant
   field, each variant capped at `maxVals`, whole union capped at `maxVals*2`.
   Special case: `CalcValueType` in `argPos` emits only the first (dimensional)
   variant.
7. **uniform shorthand** (`uniformSeqType(m) != ""` — a sequence whose 2+
   non-repeated fields all reference the same type, i.e. `<t>{1,n}`): render the
   single type once (`prov.uniform`).
8. **sequence** (default): a **function** sequence (`isFunctionSeq` — has a `(`
   terminal field) is rendered as ONE canonical instance with leaves collapsed
   to a single sample; a plain sequence is the bounded cartesian product
   (`joinProduct(parts, maxProduct)`), with EBNF-optional fields appended an
   extra `""` so absent variants are last.

`renderField` (`render.go`):
```go
func (r *Renderer) renderField(f *FieldDescriptorProto, depth int, seen, argPos) []string {
    if f.Type == TYPE_STRING { return nil }            // bare scalar — no grammar value
    if f.Label == LABEL_REPEATED {                     // EBNF { } — OMITTED entirely
        r.cur.repeated = true; return nil
    }
    return r.render(f.GetTypeName(), depth+1, seen, argPos)
}
```

> **⚠️ This is the single biggest SVG porting hazard. See §6.**

Helpers: `isFunctionSeq` (detects `(` literal field), `uniformSeqType`,
`joinProduct`/`glue` (CSS-aware spacing: no space before `,` `/` `)` `(`; no
space after `(`), `productOverflows`, `dedupe`/`cleanValue` (trims leading
`,`/`/`, collapses whitespace, fixes punctuation spacing).

`prov` struct: `usedRep, capped, depthCut, cycleCut, uniform, repeated bool;
assists map[string][]string`.

### 3.4 `reps.go` — representative leaf samples

`var reps map[string][]string` — substitutes ONLY the atomic, open-ended leaf
types (everything composite is walked). Keys are compiled message names
(PascalCase of the rule). Full CSS list of leaf type names and their sample
sets (the port replaces these with SVG leaf types):

| Leaf message | Samples |
|---|---|
| `LengthType` / `NonNegativeLengthType` | `24px 48px 8px 16px 64px` |
| `DimensionType` | `24px 2em 48px 0.5rem` |
| `PercentageType` / `NonNegativePercentageType` | `50% 80% 30% 65% 15% 100%` |
| `NumberType` | `1 0.5 2 0.75 0.25` |
| `NonNegativeNumberType` | `1 2 0.5 0.75 0.25` |
| `IntegerType` | `2 1 3 5 4` |
| `NegativeIntegerType` | `-1 -2 -3` |
| `NonNegativeIntegerType` | `2 1 3 5 4` |
| `PositiveIntegerType` | `2 3 1 5 4` |
| `AngleType` | `45deg 135deg 250deg 320deg 90deg` |
| `ObliqueAngleType` | `14deg -14deg 45deg -45deg 90deg` |
| `CubicBezierProgressType` | `0.25 0.5 0.75 0 1` |
| `MiterlimitType` | `4 1 10 2 1.5` |
| `InitialLetterSizeType` | `2 1.5 3 1 2.5` |
| `FilterAmountType` | `0.4 1.8 0.7 2` |
| `ScaleFactorType` | `2 0.5 1.6 0.7` |
| `FontWeightNumberType` | `900 100 400 700 300` |
| `RotateAxisVectorType` | `1 1 1 / 1 0 0 / 0 1 0 / 1 1 0` (multi-token) |
| `BorderImageSliceList`/`WidthList`/`OutsetList` | varied 1–4 value tuples |
| `BorderRadiusList` | corner tuples incl. `/` form |
| `FontShortSizeType` / `FontShortLineHeightType` | `24px 16px` / `1.4` |
| `TranslateList` / `ScaleList` | x[y[z]] tuples |
| `SpreadShadowType` / `ShadowType` | curated shadows (color+offsets+blur) |
| `ClipRectEdgesType` | `8px, 56px, 56px, 8px` |
| `QuotesPropItem` | quote pairs (multi-token) |
| `TimeType` / `NonNegativeTimeType` | `0.3s 1s 200ms` / `0.3s 1s` |
| `FrequencyType` | `440Hz 1kHz` |
| `ResolutionType` | `96dpi 2x 300dpi` |
| `FlexType` / `NonNegativeFlexType` | `1fr 2fr 3fr` / `1fr 2fr` |
| `HexColorType` | `#c5483c #2f5fd0` |
| `FeatureTagType` / `AxisTagType` / `AxisValueType` | OpenType tags / axis values |
| `StringType` | `"Specimen" "Aa"` |
| `IdentType` / `CustomIdentType` / `DashedIdentType` | `alpha beta` / `my-ident` / `--my-var` |
| `TransitionPropertyNameType` | real property names |
| `UrlType` / `MaskUrlType` | `url(assets/photo.jpg)` / mask SVG |
| `LinearGradientFn` …`BasicShapeType` | one canonical instance each |

Multi-token reps (e.g. `RotateAxisVectorType`, `QuotesPropItem`) exist
**specifically because** identical adjacent leaves collapse to one repeated
proto field; supplying the whole tuple as one rep works around that. This is a
direct consequence of the repeated-field omission in §6 and will be pervasive
in SVG.

`leafDisplay map[string]string` maps a leaf message name → CSS data-type
spelling (`"LengthType" → "<length>"`) for the provenance banner; `leafName`
falls back to the raw name.

### 3.5 `classify.go` — family taxonomy + value-type inference

- `families []Family` — 34 presentation groups (id/sigil/title/blurb/demo/
  sampleKind/group/focus). Note family `32` is `"SVG Paint & Geometry"`.
- `classifyRules []classifyRule` + `classify(name) string` — maps a property
  name to a family id (first match wins; `has`/`pre`/`exact` predicates).
- `inferValueType(values []string) string` — returns a control hint
  (`none/angle/color/length/number/function/keyword`) by regex-matching the
  rendered values (`reNumber/reLength/reAngle/reColor`). **This is purely
  presentational — NOT the "provenance" the brief mentions.**

> **Provenance** (`pure`/`assisted`/`empty`) is produced in `main.go`, not
> classify.go: `classifyProvenance(vals, p)` returns `"empty"` (no values),
> `"assisted"` (`p.usedRep` — a leaf was sampled), else `"pure"`.
> `analyzeProvenance` (in `main.go`, run with `-analyze`) reports a finer
> breakdown: pure-complete / all-keywords-collapsed / pure-truncated /
> representative-assisted / empty.

### 3.6 `emit.go` — output files

`emit(props, outDir)` groups props by family and writes three files into
`chrome-testing/generated/` (default `-out`):

- **`codex-data.jsx`** — an IIFE defining `window.CODEX`:
  ```
  window.CODEX = {
    total, stats: { pureComplete, pureTruncated, assisted, empty },
    familyCount, galleryCount, families, byId,
  }
  ```
  Each family carries `properties: [...]`, each property:
  ```
  { number, name, maturity, description, experimental, nonstandard, deprecated,
    warning, syntax, ebnf, valueType, defaultValue, related, provenance,
    truncated, trueCount, shown, assists, values: [ v(value, css), … ] }
  ```
  (`v = (value, css, label) => ({value, css, note:"", label})`). Values capped
  at `maxGalleryValues = 50`.
- **`manifest.tsv`** — one tab-separated row per property:
  `name⇥familyID⇥Prop⇥valueType⇥len(Values)⇥provenance⇥trueCount`.
- **`values.json`** — `map[propertyName][]string` (the shown values), pretty-printed.

---

## 4. Shell wiring

All scripts are in `/Volumes/wd_office_3/Projects/proto-css/`.

### chrome-testing/*

| Script | Does | Calls / ports / artifacts |
|---|---|---|
| `gen.sh` | Generate gallery data from the grammar (idempotent). Runs `docs.sh` first if `generated/descriptions.json` missing. | `go run ./chrome-testing/cmd/gen/ -lang lang -out chrome-testing/generated`. Then `dist.sh` unless `SKIP_DIST=1`. Output: `generated/codex-data.jsx`. |
| `build.sh` (CT) | Full CT pipeline. | `gen.sh` then `shoot.sh`. |
| `serve.sh` | Serve the self-contained `dist/` bundle and open it. | `dist.sh`, then `python3 -m http.server $SERVE_PORT --directory dist` (**port 8888** default, `SERVE_PORT` override). Kills stale server on the port first. |
| `dist.sh` | Assemble deployable static bundle into `dist/`. | Copies `gallery/*.jsx`, `styles.css`, `assets/`, `generated/codex-data.jsx`; vendors React/Babel (curl from unpkg, cached in `.vendor-cache`, `VENDOR_OFFLINE=1`); rewrites `index.html` paths. Output: `chrome-testing/dist/`. |
| `docs.sh` | Fetch MDN property descriptions → `generated/descriptions.json`. | `go run ./chrome-testing/cmd/mdndesc/ docs/reference/mdnproperties-reference.md`. |
| `shoot.sh` | Serve gallery + headless `chromerpc`, screenshot every property. | Builds `chromerpc`+`automate` from `$CHROMERPC_SRC` (default `~/Documents/chromerpc`) into `/tmp/chromerpc-ct2/bin`; serves CT on a free port; `go run ./chrome-testing/cmd/shoot/ -values generated/values.json -manifest generated/manifest.tsv …`; runs `automate` chunks; post-processes PDFs→PNG and frames→GIF (`go run ./chrome-testing/cmd/gifenc/`). Output: `chrome-testing/screenshots/`. `ONLY=`, `RESUME=`, `RESTART_EVERY=40`. |

### root/*

| Script | Does | Calls |
|---|---|---|
| `setup.sh` | Check prereqs (go, python3, Chrome), best-effort install deploy/proto tooling, build `chromerpc` into `/tmp/chromerpc-testing/bin`, `go mod tidy`, verify EBNF files + that `cmd/gen` and `cmd/shoot` build. | `go mod tidy`; `go build ./chrome-testing/cmd/gen/ ./chrome-testing/cmd/shoot/`. |
| `build.sh` (root) | "EBNF → HTML → screenshots → gallery". | **Calls `chrome-testing/run.sh --generated`** — see ⚠️ below. |
| `test.sh` | The only sanctioned pre-commit validator: full pipeline build, `go vet`, `go test -race`, template/screenshot/gallery checks, smoke-test (serves `generated_gallery.html` on a random port, checks HTTP 200). | **Calls `chrome-testing/run.sh`** (no flag) — see ⚠️ below. |
| `LET_IT_RIP.sh` | Full pipeline: `setup.sh` → `gen.sh` → (shoot, currently **disabled** in source — "Takes too much time") → `serve.sh`. `SKIP_SERVE=1` stops before serving. | `setup.sh`, `gen.sh`, `serve.sh` (which serves on port 8888). |

> ⚠️ **`chrome-testing/run.sh` does NOT exist.** Root `build.sh` and `test.sh`
> both invoke it (`build.sh`: `chrome-testing/run.sh --generated`; `test.sh`:
> `chrome-testing/run.sh`). Root `build.sh` even guards with
> `if [[ -x "$CHROME_TESTING/run.sh" ]]` and prints
> `"ERROR: chrome-testing/run.sh not found or not executable"` then `exit 1`.
> So in the CURRENT proto-css tree, **root `build.sh` and `test.sh` are broken**;
> the working entry points are `LET_IT_RIP.sh` / `chrome-testing/gen.sh` /
> `chrome-testing/build.sh`. The SVG port should either create a `run.sh` or
> (cleaner) point its root `build.sh`/`test.sh` at `chrome-testing/gen.sh` +
> `chrome-testing/build.sh` directly. Confirmed by `ls`: no
> `chrome-testing/run.sh` is present.

**Genproto is not wired into any shell script.** `genproto` lives in
`lang/cmd/genproto/` and is run manually (it produces the committed `proto/`
schema + maps that the *service*/parser use). The gallery `gen` command
re-derives its own in-memory FDP and does not consume `genproto`'s output. The
SVG port should add an explicit `go run ./lang/cmd/genproto/` invocation to a
script if it wants the on-disk proto regenerated.

---

## 5. `go.mod` for proto-svg

Mirror proto-css exactly, swapping the module path:

```
module github.com/accretional/proto-svg

go 1.25.5

replace github.com/accretional/gluon => ../gluon
replace github.com/accretional/merge => ../proto-merge

require (
    github.com/accretional/gluon v0.0.0-20260430002741-82b4ddbe67a4
    github.com/accretional/merge v0.0.0-00010101000000-000000000000
    google.golang.org/grpc v1.80.0
    google.golang.org/protobuf v1.36.11
)

require (
    github.com/accretional/proto-expr v0.0.0-20260416071217-9a69001c59bb // indirect
    github.com/golang/protobuf v1.5.4 // indirect
    github.com/jhump/protoreflect v1.18.0 // indirect
    github.com/jhump/protoreflect/v2 v2.0.0-beta.1 // indirect
    golang.org/x/net v0.52.0 // indirect
    golang.org/x/sys v0.42.0 // indirect
    golang.org/x/text v0.35.0 // indirect
    google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
)
```

(`go mod tidy` will normalize versions; the two `replace` directives and the
two `require`d accretional modules are the load-bearing parts. The
`proto-expr` indirect comes in transitively via gluon.)

**Import paths a genproto/gen port must use:**
```go
import (
    "github.com/accretional/gluon/v2/compiler"
    metaparser "github.com/accretional/gluon/v2/metaparser"
    pb         "github.com/accretional/gluon/v2/pb"
    "github.com/accretional/merge/descriptor"          // genproto only (ToString)
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/descriptorpb"
)
```
(`gen` uses `compiler`, `metaparser`, `pb`, `descriptorpb` but **not**
`descriptor` or `proto`.)

---

## 6. SVG-specific porting notes — repeated fields are CENTRAL

proto-css is a **flat declaration grammar**: a stylesheet is a flat list of
`property: value;` pairs, and the value grammar is a tree that bottoms out in
scalar leaves. Crucially, proto-css **deliberately omits repeated fields** at
render time, because in CSS a `{ }` list is almost always a comma/space list of
values where one representative element suffices.

SVG is an **attributed tree**: `element → attributes`, `element → children`,
`attribute → value`. Here the EBNF `{ }` repetition **IS the document
structure** — a `<g>` has repeated child elements; an element has repeated
attributes; a `<path d>` has a repeated list of path commands; `points` is a
repeated coordinate list. **Omitting repeated fields would erase the entire
document.** Every place the proto-css renderer/reps assume "flat CSS, drop the
repetitions" must be revisited. The concrete sites:

### 6.1 `render.go: renderField` — repeated fields hard-dropped (THE critical fix)
```go
if f.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
    if r.cur != nil { r.cur.repeated = true }
    return nil            // ← returns NOTHING for every { } repetition
}
```
This is the #1 change. In CSS, omitting `{ }` avoids leading separators and
combinatorial blow-up and is harmless (the singular head of a `+`/`#` list is a
separate field that still renders). In SVG, this would emit elements with **no
children and no attributes**. The port must instead emit 1..N repetitions
(e.g. render the element type once or a few times and join with the field's
`FieldSeparator`, or — for child-element lists — recurse and emit a small fixed
count). The whole notion of "the mandatory head is a separate singular field"
does not hold for SVG attribute/child lists, which are pure `{ }`.

### 6.2 `render.go: uniformSeqType` — ignores repeated fields
```go
for _, f := range m.GetField() {
    if f.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED { continue } // ← skipped
    ...
}
```
The `<t>{1,4}` shorthand collapse skips repeated fields entirely. For SVG
coordinate/point lists this collapse is both wrong (a `points` list is not a
1-value shorthand) and would interact badly once repeated fields actually
render. Re-evaluate whether uniform-collapse should apply to SVG at all.

### 6.3 `render.go: render` sequence branch — optional-absent + product cap
The plain-sequence path appends `""` for EBNF-optional fields and caps the
cartesian product at `maxProduct=50`. For an SVG element with many attributes
this product explosion is real (every attribute combination), and the "absent
variant ordered last" trick assumes a small fixed field count. The port needs a
different strategy for element-with-attributes: pick a representative attribute
*set* per element rather than the cross-product of all attribute values.

### 6.4 `render.go: isFunctionSeq` / `glue` — CSS punctuation assumptions
`isFunctionSeq` keys on a `(` terminal (CSS functions). `glue` hard-codes CSS
spacing rules (`,` `/` `)` `(` attach without spaces; space-join otherwise).
SVG syntaxes differ: path data uses command letters + space/comma-separated
numbers (`M 10 10 L 90 90`), `transform` uses `translate(…)` (function-like),
`viewBox`/`points` are space/comma number lists, attribute values are
quoted strings inside `name="value"`. The CSS `glue`/`cleanValue` punctuation
model must be replaced with SVG-aware joining (likely several, keyed by context:
path-data vs. point-list vs. attribute-string vs. transform-function).

### 6.5 `reps.go` — multi-token reps are a *workaround* for 6.1, and CSS leaf set is wrong
- The CSS leaf set (`LengthType`, `AngleType`, `HexColorType`, …) must be
  replaced with SVG leaf types: SVG `<number>`, `<coordinate>`,
  `<length>` (note SVG lengths/percentages differ from CSS), `<color>`,
  `<paint>` (`none|currentColor|<color>|url(#id)`), `<path-data>` (a whole
  mini-grammar), `<points>`, `<transform-list>`, `<list-of-numbers>`,
  `<angle>`, IRI/`url(#id)` references, `<opacity-value>`, etc.
- The existing multi-token reps (`RotateAxisVectorType: "1 1 1"`,
  `QuotesPropItem`, `TranslateList`, `BorderRadiusList`) exist **only because**
  §6.1 collapses identical adjacent leaves and omits repetitions — the author
  smuggles whole tuples through as a single string. In SVG this hack would be
  needed for *almost every* coordinate/point/path-data leaf, which is a strong
  signal that the port should fix §6.1 properly (render repetitions) rather than
  pre-bake every multi-token value as a rep string. At minimum, `d` (path data),
  `points`, `viewBox`, `transform`, `stroke-dasharray` will each need either a
  curated canonical instance (rep) or genuine repeated-field rendering.

### 6.6 `main.go: enumerate` / `propertyKeyword` / `findPropFQN` — CSS `Property`/`*Expr`/`*Prop` shape
`enumerate` is hard-wired to the CSS grammar's `Property` oneof of `XExpr`
messages, each wrapping an `XProp` value rule, with the property name being the
leading keyword. SVG has no such `property: value` declaration shape — it has
elements and attributes. The entire enumeration entry point must be rewritten to
walk the SVG root/element messages and enumerate (element, attribute, value)
triples, or (element → rendered instance). `propertyKeyword` (reads the leading
keyword to get the kebab name) and `findPropFQN` (finds the `*Prop` field) are
CSS-shape-specific and have no SVG analogue.

### 6.7 `classify.go` families & `emit.go` schema are CSS-property-centric
`families`/`classifyRules` taxonomize CSS properties; `emit.go`'s `window.CODEX`
schema is `families → properties → values`. For SVG the natural axis is
elements (and their attributes), so both the taxonomy and the emitted schema
need an SVG-shaped redesign (e.g. `elements → attributes → sample-values`, or
`elements → rendered SVG snippets`). `inferValueType`'s CSS regexes
(`reLength` with `px/em/rem/…`, `reColor`) also need SVG value vocabularies.

### 6.8 genproto leaf/keyword tuning is grammar-specific
In `genproto`, `scalarizeLeaves`'s `leafTypes` list and the `pruneUnreachable`
root (`"CssStyleSheet"`) are CSS-specific. For SVG: set the prune root to the
SVG document/root-element rule (e.g. `Svg`), and curate an SVG `leafTypes` set.
The keyword-emptying and prefix/separator machinery is grammar-agnostic and
ports unchanged — but note SVG keywords are element/attribute names and enum
keyword values, and SVG's heavy use of `{ }` means `CollapseCommaList` and the
`FieldSeparator` map will be far more important than in CSS (where 113
separators were recorded; SVG will have proportionally more, and they carry the
real list syntax).

### 6.9 Summary of the structural inversion
| Aspect | proto-css (flat) | proto-svg (tree) |
|---|---|---|
| Top structure | flat `property: value` list | nested `element{attrs, children}` |
| `{ }` repetition | mostly omitted at render | **must be rendered** (children, attrs, points, path-data) |
| Render entry | `Property` oneof → `XExpr`/`XProp` | root element → recurse elements/attrs |
| Leaf reps | CSS `<length>`/`<color>`/… | SVG `<coordinate>`/`<paint>`/`<path-data>`/… |
| Joining | CSS `glue` punctuation | SVG path/point/transform/attribute syntaxes |
| Output schema | families → properties → values | elements → attributes → values (redesign) |

The pipeline **infrastructure** (gluon v2 import resolution, the
WrapString→ParseEBNF→GrammarToAST→transforms→Compile chain, the OnMessage/OnField
prefix+separator collection, the descriptor cleanup passes, the
fdset/proto/map outputs, the dist/serve shell wiring) ports essentially
unchanged. The **grammar-specific layers** — `scalarizeLeaves.leafTypes`, the
prune root, `reps`, `render.go`'s repeated-field omission, the
`Property`/`Expr`/`Prop` enumeration, `classify`/`emit` schemas — are where all
the SVG work concentrates, and the repeated-field omission (§6.1) is the single
change without which the SVG output is structurally empty.
