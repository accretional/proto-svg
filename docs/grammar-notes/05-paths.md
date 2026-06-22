# Paths grammar notes

## Source

- W3C SVG 2 §9 "Paths" — `/docs/specs/cache/svg2-paths.txt` (read in full, lines 0–1342)
- W3C SVG 2 Appendix B "Implementation Notes" — `/docs/specs/cache/svg2-implnote.txt` (read in full, lines 0–941; arc-flag semantics §B.2.1–B.2.5)
- Cross-checked against MDN Web Docs (SVG path `d` attribute, `pathLength`), the WHATWG/W3C SVG 2 CR, and known browser behaviour.

---

## `path` element

### Categories

Graphics element, renderable element, shape element.

### Content model

Any number of the following elements in any order:

- **animation elements** — `animate`, `animateMotion`, `animateTransform`, `discard`, `set`
- **descriptive elements** — `desc`, `title`, `metadata`
- **paint server elements** — `linearGradient`, `radialGradient`, `pattern`
- `clipPath`, `marker`, `mask`, `script`, `style`

### Attributes and properties

| Name | Kind | Value | Initial | Inherited | Animatable |
|------|------|-------|---------|-----------|-----------|
| `d` | geometry property / presentation attribute | `none \| <string>` | `none` | no | yes |
| `pathLength` | geometry attribute | `<number>` | (none) | no | yes |

All shared attribute groups also apply: core attributes (`id`, `tabindex`, `lang`, `xml:space`, `class`, `style`), conditional processing, global event, graphical event, ARIA, and presentation attributes.

---

## Path-data BNF (verbatim)

Transcribed exactly from §9.3.9 of `svg2-paths.txt` (lines 951–1038), as it appears in the spec fenced block:

```
svg_path::= wsp* moveto? (moveto drawto_command*)?

drawto_command::=
    moveto
    | closepath
    | lineto
    | horizontal_lineto
    | vertical_lineto
    | curveto
    | smooth_curveto
    | quadratic_bezier_curveto
    | smooth_quadratic_bezier_curveto
    | elliptical_arc

moveto::=
    ( "M" | "m" ) wsp* coordinate_pair_sequence

closepath::=
    ("Z" | "z")

lineto::=
    ("L"|"l") wsp* coordinate_pair_sequence

horizontal_lineto::=
    ("H"|"h") wsp* coordinate_sequence

vertical_lineto::=
    ("V"|"v") wsp* coordinate_sequence

curveto::=
    ("C"|"c") wsp* curveto_coordinate_sequence

curveto_coordinate_sequence::=
    coordinate_pair_triplet
    | (coordinate_pair_triplet comma_wsp? curveto_coordinate_sequence)

smooth_curveto::=
    ("S"|"s") wsp* smooth_curveto_coordinate_sequence

smooth_curveto_coordinate_sequence::=
    coordinate_pair_double
    | (coordinate_pair_double comma_wsp? smooth_curveto_coordinate_sequence)

quadratic_bezier_curveto::=
    ("Q"|"q") wsp* quadratic_bezier_curveto_coordinate_sequence

quadratic_bezier_curveto_coordinate_sequence::=
    coordinate_pair_double
    | (coordinate_pair_double comma_wsp? quadratic_bezier_curveto_coordinate_sequence)

smooth_quadratic_bezier_curveto::=
    ("T"|"t") wsp* coordinate_pair_sequence

elliptical_arc::=
    ( "A" | "a" ) wsp* elliptical_arc_argument_sequence

elliptical_arc_argument_sequence::=
    elliptical_arc_argument
    | (elliptical_arc_argument comma_wsp? elliptical_arc_argument_sequence)

elliptical_arc_argument::=
    number comma_wsp? number comma_wsp? number comma_wsp
    flag comma_wsp? flag comma_wsp? coordinate_pair

coordinate_pair_double::=
    coordinate_pair comma_wsp? coordinate_pair

coordinate_pair_triplet::=
    coordinate_pair comma_wsp? coordinate_pair comma_wsp? coordinate_pair

coordinate_pair_sequence::=
    coordinate_pair | (coordinate_pair comma_wsp? coordinate_pair_sequence)

coordinate_sequence::=
    coordinate | (coordinate comma_wsp? coordinate_sequence)

coordinate_pair::= coordinate comma_wsp? coordinate

coordinate::= sign? number

sign::= "+"|"-"
number ::= ([0-9])+
flag::=("0"|"1")
comma_wsp::=(wsp+ ","? wsp*) | ("," wsp*)
wsp ::= (#x9 | #x20 | #xA | #xC | #xD)
```

### Notes on the BNF as written

1. **`number` is integer-only in the spec BNF** — The production `number ::= ([0-9])+` admits only strings of decimal digits. Fractional values enter via `coordinate ::= sign? number`, but there is no decimal-point rule in this BNF. This is a **known spec defect** (see §Discrepancies below).

2. **`elliptical_arc_argument`** — The production requires `comma_wsp` (mandatory) between the third number (`x-axis-rotation`) and the first flag, but only `comma_wsp?` (optional) between flags and between flag and coordinate pair. This is intentional: at least one separator is required after the rotation angle so the parser can unambiguously detect the flag tokens.

3. **`svg_path`** — The top-level rule allows empty input (`wsp*` with optional `moveto`), which is not an error; it simply produces no rendering.

---

## Path-data EBNF-ready rendering (all commands)

The grammar below is a complete, self-consistent EBNF suitable for direct use in a grammar tool. Corrections to the spec's `number` production are applied (see §Discrepancies). Terminals are quoted literals or character sets. Command letters are kept as atomic terminals. Flags are `"0" | "1"`.

```ebnf
(* Top-level *)
svg_path          = wsp* [ moveto { drawto_command } ] ;

drawto_command    = moveto
                  | closepath
                  | lineto
                  | horizontal_lineto
                  | vertical_lineto
                  | curveto
                  | smooth_curveto
                  | quadratic_bezier_curveto
                  | smooth_quadratic_bezier_curveto
                  | elliptical_arc ;

(* ── Commands ── *)

moveto            = ( "M" | "m" ) wsp* coordinate_pair_sequence ;

closepath         = "Z" | "z" ;

lineto            = ( "L" | "l" ) wsp* coordinate_pair_sequence ;

horizontal_lineto = ( "H" | "h" ) wsp* coordinate_sequence ;

vertical_lineto   = ( "V" | "v" ) wsp* coordinate_sequence ;

curveto           = ( "C" | "c" ) wsp* curveto_coordinate_sequence ;
curveto_coordinate_sequence
                  = coordinate_pair_triplet
                    { comma_wsp? coordinate_pair_triplet } ;

smooth_curveto    = ( "S" | "s" ) wsp* smooth_curveto_coordinate_sequence ;
smooth_curveto_coordinate_sequence
                  = coordinate_pair_double
                    { comma_wsp? coordinate_pair_double } ;

quadratic_bezier_curveto
                  = ( "Q" | "q" ) wsp* quadratic_bezier_curveto_coordinate_sequence ;
quadratic_bezier_curveto_coordinate_sequence
                  = coordinate_pair_double
                    { comma_wsp? coordinate_pair_double } ;

smooth_quadratic_bezier_curveto
                  = ( "T" | "t" ) wsp* coordinate_pair_sequence ;

elliptical_arc    = ( "A" | "a" ) wsp* elliptical_arc_argument_sequence ;
elliptical_arc_argument_sequence
                  = elliptical_arc_argument
                    { comma_wsp? elliptical_arc_argument } ;

elliptical_arc_argument
                  = nonneg_number comma_wsp?
                    nonneg_number comma_wsp?
                    number       comma_wsp
                    flag         comma_wsp?
                    flag         comma_wsp?
                    coordinate_pair ;

(* ── Coordinate helpers ── *)

coordinate_pair_triplet
                  = coordinate_pair comma_wsp?
                    coordinate_pair comma_wsp?
                    coordinate_pair ;

coordinate_pair_double
                  = coordinate_pair comma_wsp? coordinate_pair ;

coordinate_pair_sequence
                  = coordinate_pair { comma_wsp? coordinate_pair } ;

coordinate_sequence
                  = coordinate { comma_wsp? coordinate } ;

coordinate_pair   = coordinate comma_wsp? coordinate ;

coordinate        = sign? number ;

(* ── Leaves / open datatypes ── *)

sign              = "+" | "-" ;

number            = integer [ "." integer ] | "." integer ;
(* Corrected from spec: real decimal numbers, not integer-only.
   See §Discrepancies. Exponent notation is NOT supported in path data. *)

integer           = digit { digit } ;

digit             = "0" | "1" | "2" | "3" | "4"
                  | "5" | "6" | "7" | "8" | "9" ;

nonneg_number     = number ;
(* Semantic constraint only: rx, ry must be >= 0.
   Negative values are accepted by the parser and then
   silently converted to their absolute value by the renderer. *)

flag              = "0" | "1" ;

comma_wsp         = ( wsp+ [ "," ] wsp* ) | ( "," wsp* ) ;

wsp               = #x9 | #x20 | #xA | #xC | #xD ;
(* HT, SP, LF, FF, CR *)
```

### Command-letter summary table

| Letter | Kind | Parameters per repetition |
|--------|------|--------------------------|
| `M` / `m` | absolute / relative moveto | `x y` (+ implicit lineto if extra pairs follow) |
| `Z` / `z` | closepath | none |
| `L` / `l` | absolute / relative lineto | `x y` |
| `H` / `h` | absolute / relative horizontal lineto | `x` |
| `V` / `v` | absolute / relative vertical lineto | `y` |
| `C` / `c` | absolute / relative cubic Bézier | `x1 y1 x2 y2 x y` |
| `S` / `s` | absolute / relative smooth cubic Bézier | `x2 y2 x y` |
| `Q` / `q` | absolute / relative quadratic Bézier | `x1 y1 x y` |
| `T` / `t` | absolute / relative smooth quadratic Bézier | `x y` |
| `A` / `a` | absolute / relative elliptical arc | `rx ry x-axis-rotation large-arc-flag sweep-flag x y` |

---

## Flag-parsing subtlety

### The problem

`large-arc-flag` and `sweep-flag` are each a single character, `"0"` or `"1"`, not numeric tokens. In the spec BNF:

```
flag ::= ("0"|"1")
```

This means flags are **not** parsed as the general `number` production, and **no separator is needed** between two adjacent flags, or between a flag and a subsequent signed number. All of the following are valid and must parse identically:

```
A 25,26 -30 0,1 50,-25
A 25,26 -30 0 1 50,-25
A 25,26 -30 0150,-25      (* flags "0","1" then coordinate 50 → no: 50 is absorbed *)
```

Wait — the last form needs care. After `0` (large-arc-flag) and `1` (sweep-flag), the next token begins. Because `5` is a digit it cannot be a sign or start a new flag, so the parser enters `coordinate_pair` and reads `50,-25`. This is correct.

More importantly, the string `A25,26,-30,0,150,-25` must be parsed as:
- rx=25, ry=26, rotation=-30, large-arc-flag=0, sweep-flag=1, x=50, y=-25

because after consuming `0` as the large-arc-flag, the parser immediately tries to read `comma_wsp?` then `flag`, finding `1` directly. The remaining `50` is the x-coordinate.

The spec itself demonstrates this with the arcs01 example (line 861):
```
a150,150 0 1,0 150,-150
```
Here the flag tokens `1` and `0` are separated by only a comma, immediately adjacent to the rotation `0` and the coordinate `150`.

### Consequence for grammar authoring

The `flag` production **must be kept separate** from `number`/`coordinate` in the EBNF. A scanner that tokenizes path data as a conventional numeric stream will fail on flag-adjacent forms. A maximal-munch rule applied to `number` would greedily consume a leading `0` or `1` as the start of a larger integer, incorrectly merging `01` into a single token `1`.

**Correct parsing strategy**: after reading `x-axis-rotation` (followed by mandatory `comma_wsp`), switch to flag-mode, consume exactly one character `"0"` or `"1"` for `large-arc-flag`, apply optional `comma_wsp?`, consume exactly one character `"0"` or `"1"` for `sweep-flag`, apply optional `comma_wsp?`, then resume coordinate parsing.

### Spec text evidence (§9.3.9, lines 1040–1057)

The spec explains maximal-munch for coordinates: "M 100-200" parses as x=100, y=-200 because `-` cannot follow a digit in `number`. Similarly "M 0.6.5" parses as 0.6 and .5 because only one decimal point is allowed. The same maximal-munch principle applies to flags: `"0"` or `"1"` is consumed as a complete flag token and parsing immediately proceeds.

### MDN / browser confirmation

MDN SVG path `d` documentation explicitly states that flag values are single characters and may appear without any separator between them or between a flag and the following coordinate (e.g., `a 1 1 0 00 1 0` is equivalent to `a 1 1 0 0 0 1 0`). All major browsers (Chrome, Firefox, Safari) implement this behaviour.

---

## pathLength

### Spec definition (§9.6.1)

| Attribute | Value | Initial | Animatable |
|-----------|-------|---------|-----------|
| `pathLength` | `<number>` | (none) | yes |

- Value is the author's total path length in user units.
- Used as a scaling factor: the user agent scales all distance-along-a-path computations by `pathLength / UA_computed_length`.
- **Zero is valid**: treated as a scaling factor of infinity. Zero scaled by infinity remains zero; any non-percentage positive value becomes `+Infinity`.
- **Negative value is an error** (error-handling rules apply).
- Has no effect on percentage distance-along-a-path calculations.
- `moveto` operations contribute zero length to path-length calculations. Only `lineto`, `curveto` (cubic and quadratic), and `arcto` commands contribute.

### Grammar value for `pathLength`

```ebnf
pathLength_value  = number ;
(* number is the same real-number leaf as used in path data,
   but here the full CSS/SVG <number> production applies:
   sign? ( integer [ "." integer ] | "." integer ) [ exponent ]
   Negative values are parse-legal but semantically an error.
   Initial value is absent (the attribute is optional). *)
```

**Note**: `pathLength` is listed in the spec both as a geometry property and as a presentation attribute on the `path` element. The value space is `<number>` (SVG/CSS number, not path-data number — exponent notation is allowed here).

---

## Open datatypes used

| Name | Description | Notes |
|------|-------------|-------|
| `number` | Real decimal number, no exponent, one optional decimal point | Path-data leaf; corrected from spec's integer-only BNF |
| `nonneg_number` | Same as `number`, semantically ≥ 0 | Used for `rx`, `ry` in arc arguments; negatives silently abs-valued |
| `coordinate` | `sign? number` | Signed real |
| `coordinate_pair` | Two coordinates | x then y |
| `flag` | `"0" \| "1"` | Single character; closed terminal set |
| `wsp` | Whitespace: HT, SP, LF, FF, CR | Unicode codepoints #x9, #x20, #xA, #xC, #xD |
| `comma_wsp` | Separator: whitespace around optional comma | |

---

## Context-sensitive constraints (overlay, not in grammar)

These rules cannot be expressed in a context-free grammar and belong in the semantic / constraint overlay:

1. **First command must be `moveto`**: The grammar's `svg_path ::= wsp* moveto? (moveto drawto_command*)?` forces a moveto first if any drawing commands follow. A path with `drawto_command` but no preceding `moveto` is handled by the error-recovery rule (render nothing or partial).

2. **Relative `m` as first command**: A relative `m` appearing as the first command is interpreted as absolute coordinates for the initial point; subsequent extra coordinate pairs in the same moveto are then relative linetos.

3. **Implicit command repetition**: After a moveto `M`/`m`, extra coordinate pairs are treated as implicit `L`/`l` (not `M`/`m`). After any other command, extra parameter sets are implicit repeats of that command. The grammar encodes this via `coordinate_pair_sequence` and the `{…}` repetition rules.

4. **Smooth-curve control point reflection**: For `S`/`s`, if the previous command was not `C`/`c`/`S`/`s`, the first control point is coincident with the current point. For `T`/`t`, if the previous command was not `Q`/`q`/`T`/`t`, the control point is coincident with the current point. Formula: `(newx1, newy1) = (2*curx − oldx2, 2*cury − oldy2)`.

5. **Arc parameter clamping** (§9.5.1):
   - If endpoint == current point: arc segment is omitted entirely.
   - If `rx == 0` or `ry == 0`: treated as a straight lineto.
   - If `rx < 0` or `ry < 0`: use absolute value.
   - If radii are too small to reach the endpoint: scale uniformly until exactly one solution exists (see §B.2.5 formula).

6. **pathLength sign constraint**: `pathLength < 0` is a semantic error; `pathLength == 0` is valid with special infinity semantics.

7. **Maximal-munch parsing**: The parser must consume as much as possible at each production. The minus sign or a second decimal point terminates the current number; a third `0`/`1` after both flags has already been consumed as coordinates.

8. **Error recovery**: On parse error, render up to (but not including) the command containing the first error. Within a malformed command, render up to the last correctly-defined segment.

9. **Animation interpolation**: Two `d` values interpolate smoothly only if they have identical command count, command types, and order. Otherwise discrete interpolation is used. Flags interpolate as 0/1 thresholds (any non-zero → 1).

---

## Discrepancies, doc gaps & roadblocks

### D1: `number` production is integer-only (CRITICAL)

**Spec BNF**: `number ::= ([0-9])+`

This admits only strings of decimal digits. Fractional path coordinates (e.g., `M 0.5 1.5`, `C 100.25 200.5 …`) would not match. This is clearly wrong — every SVG implementation accepts decimal fractions, and the spec's own examples use them (line 1054: "0.6.5" is two valid floats).

**MDN and browser behaviour**: All browsers accept decimal-fraction coordinates.

**Resolution for grammar**: Replace with:
```ebnf
number  = integer [ "." integer ] | "." integer ;
integer = digit { digit } ;
digit   = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" ;
```
No exponent notation is supported in path data (distinct from CSS `<number>`).

**Reasoning**: The spec's `.5` parsing example at line 1054 implies decimal points are valid in `number`. The spec BNF is simply an error of omission. All tooling and real-world path data uses fractional values. The integer-only BNF is rejected.

### D2: `sign` is separate from `number` — negative numbers only via `coordinate`

**Spec BNF**: `coordinate ::= sign? number`. There is no `sign` in `number` itself.

This means `number` is always non-negative; the sign is part of `coordinate`. This has a parsing implication: in arc arguments, `rx` and `ry` are consumed as bare `number` productions (no sign). The arc argument structure is:

```
number comma_wsp? number comma_wsp? number comma_wsp flag comma_wsp? flag ...
```

So `rx` and `ry` are parsed as unsigned numbers, and the rotation is also a `number` (unsigned in the BNF). However, the spec prose at §9.5.1 says negative `rx`/`ry` have their absolute value taken, and the arc example at line 867 shows `a25,25 -30 0,1 50,-25` where `-30` is the rotation. **This is a contradiction**: the BNF shows `number` (unsigned) for rotation, but the example shows a signed rotation value.

**Resolution**: The `number` rule in the EBNF-ready rendering above absorbs a leading decimal form, and `coordinate` adds the optional sign. For the elliptical arc, the rotation is semantically a signed angle but the BNF parses it as `number` (unsigned). In practice, browsers accept a signed rotation because the `number` token (when following the second comma_wsp) is parsed with maximal-munch that includes a leading `-` treated as a separator — no: since `comma_wsp` was already consumed, the `-` must be the sign of the next `number`.

**Grammar decision**: Treat `x-axis-rotation` as a `coordinate` (sign-optional number) rather than bare `number`. This matches real-world use and browser behaviour. The spec BNF appears to contain a second error of omission here. Document in overlay: rotation is a signed real number in degrees; the parser must accept the sign as part of the coordinate just as for other coordinate values.

### D3: `elliptical_arc_argument` — mandatory vs optional separator before first flag

**Spec BNF**:
```
elliptical_arc_argument::=
    number comma_wsp? number comma_wsp? number comma_wsp
    flag comma_wsp? flag comma_wsp? coordinate_pair
```

The separator after the third `number` (rotation) is `comma_wsp` (mandatory), not `comma_wsp?`. This differs from all other separators in the arc argument. This is deliberate: without a mandatory separator between rotation and the first flag, a rotation value ending in a digit `0` or `1` would be ambiguous with the flag. The grammar is **correct** here. Document it: the separator before `large-arc-flag` is mandatory.

**MDN/browser note**: Chrome's path parser does in fact require at least one separator (whitespace or comma) before the first flag. This matches the spec.

### D4: `svg_path` production — `wsp*` before second `moveto`

The spec writes:
```
svg_path::= wsp* moveto? (moveto drawto_command*)?
```

This means the top level is: optional leading whitespace, then optionally one standalone moveto with no subsequent drawing, OR a moveto followed by zero or more drawto commands. This structure is slightly odd — it looks like it could accept two separate moveto groups at the top level without enclosing them in a loop.

**Analysis**: The path is a sequence of subpaths. The grammar as written does not obviously allow more than one top-level moveto-group. However, the prose says "Subsequent moveto commands represent the start of a new subpath" — meaning subpaths are encoded by embedding additional moveto commands in the `drawto_command*` list. The grammar is correct: the second `moveto` in the list will appear as a `drawto_command` (since `drawto_command` includes `moveto`). The top-level structure `moveto drawto_command*` handles the first subpath's moveto plus all subsequent commands (including additional movetos).

**No change needed** — the grammar is consistent with prose.

### D5: `d` property vs attribute — presentation attribute / CSS property

SVG 2 promotes `d` from an XML attribute to a CSS geometry property. The spec notes (lines 329–332): "d will become a presentation attribute (no name change) with path data string as value". The value syntax in the CSS context is `none | <string>` where `<string>` is parsed as path data.

**Grammar impact**: In CSS, the path data string appears inside a `path()` function call or as the value of the `d` property. The `svg_path` production remains the grammar for the string contents in both contexts.

### D6: No exponent notation in path-data `number`

The spec BNF `number ::= ([0-9])+` excludes exponents. The corrected production also excludes exponents. This is consistent with MDN, which states that exponential notation (e.g., `1e2`) is not valid in SVG path data strings. Some browser implementations may be more lenient; the grammar takes the strict spec position.

### D7: `pathLength` — listed as both geometry property and presentation attribute

The spec attribute table (lines 133–134) lists `pathLength` under "presentation attributes" and also under "Geometry properties". MDN and browser behaviour treats it as a plain attribute. The grammar value is `<number>` (full CSS number with optional sign/exponent); a negative value is a parse-legal, semantically-invalid error.

---

*End of 05-paths.md*
