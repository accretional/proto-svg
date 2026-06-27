# examples — the SVG API, both directions

The grammar-compiled proto schema (`proto/svg.fdset` + the `MessagePrefix` /
`FieldSeparator` maps in `proto/pb/svg/`) is a complete, self-describing SVG API.
These two tiny programs prove it carries everything needed to go **proto → SVG**
and **SVG → proto**, each by *pure proto reflection* — neither has any SVG
knowledge, no per-element code, no layout. The grammar is the only source of truth.

| dir | program | direction |
|---|---|---|
| `render/` | `render.go` | textproto → `.svg` (walk the message, emit each `MessagePrefix`, write leaves, interleave `FieldSeparator`) |
| `parse/`  | `parse.go`  | `.svg` → textproto (recursive descent over the *same* descriptors + maps, read in reverse) |

```
go run ./examples/render examples/render/flamegraph.textproto out.svg
go run ./examples/parse  examples/parse/flamegraph.svg        out.textproto
```

## render/ — proto → SVG

`flamegraph.textproto` describes a Brendan-Gregg-style flame graph as an
`svg.SvgDocument`: gradient background, a `<style>` block, a `<script>` with a
CDATA hover/zoom body, `<title>` tooltips, and colored frame `<g><rect><text>`
groups. `render.go` walks it and emits `flamegraph.svg`. It has **no flame-graph
logic** — it only knows how to walk a protobuf.

## parse/ — SVG → proto (the inverse)

`parse.go` is `render.go` run backwards: the descriptor tree *is* the grammar, so
a message's `MessagePrefix` is the literal it must start with, its fields are the
sequence to match, a `oneof` is an alternation (try the arms), a repeated field is
`{ X }`, and a `{string value}` leaf is a free terminal consumed up to the next
known literal. That's the whole parser.

- **`flamegraph.svg`** — a real, in-the-wild flame graph (~3000 frames).
- **`flamegraph.textproto`** — the `svg.SvgDocument` `parse.go` recovered from it.
- **`flamegraph_gen.svg`** — `render.go` run on that textproto. It is produced
  **entirely from the proto**; nothing crosses the proto boundary. It reproduces
  the original DOM frame-for-frame (identical `<g>/<rect>/<text>/<title>` counts).

### Faithfulness

- **Canonical round-trip is byte-exact.** Feeding `render.go`'s own output back
  through `parse.go` and re-rendering reproduces the input byte-for-byte.
- **Real-world round-trip is DOM-exact.** `flamegraph_gen.svg` differs from the
  input only in *serialization* — canonical long-form close tags instead of
  self-closing `/>`, and no insignificant indentation — never in structure,
  attributes, values, or order. It renders identically.

### Two lexical bridges (not SVG logic)

The structural grammar models canonical markup, so `parse.go` adds a thin lexical
layer for two XML serialization variants real files use — both purely textual,
with zero SVG semantics:

1. **insignificant whitespace** between tags is skipped (only before tag
   boundaries — never inside text/CDATA/CSS, which are free-terminal leaves);
2. **self-closing `/>`** is read as the canonical `></tag>`.

### The parser is also a validator

Because it is grammar-faithful, `parse.go` rejects anything the grammar doesn't
allow, pointing at the exact byte. The original `flamegraph.svg` (as emitted by
`flamegraph.pl`) used `text-anchor=""` 3074× — *invalid* SVG (the spec allows only
`start | middle | end`). The committed file is normalized to the spec default
`text-anchor="start"` (identical rendering) so it parses; an unfixed copy makes
`parse.go` report `text-anchor` at the offending offset.
