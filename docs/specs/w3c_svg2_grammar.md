# W3C SVG 2 — Grammar / Value-Syntax Sources

The focal file for authoring the EBNF grammar: sections that contain **formal value syntax** (data types, BNF grammars, list/coordinate formats) rather than prose. Base: `https://www.w3.org/TR/SVG2/`

## Basic data types & attribute syntax (types.html)

The foundation for `datatype.ebnf` — defines `<number>`, `<integer>`, `<length>`, `<percentage>`, `<angle>`, `<color>`, `<paint>`, `<coordinate>`, list types, etc.

| Section | Link |
| --- | --- |
| Definitions (basic data types) | https://www.w3.org/TR/SVG2/types.html#definitions |
| Attribute syntax | https://www.w3.org/TR/SVG2/types.html#syntax |
| Real number precision | https://www.w3.org/TR/SVG2/types.html#Precision |
| Range clamping | https://www.w3.org/TR/SVG2/types.html#RangeClamping |

## Path data grammar (paths.html) — KEY

The complete BNF for the `d` attribute. Drives `path.ebnf`.

| Section | Link |
| --- | --- |
| Path data (overview) | https://www.w3.org/TR/SVG2/paths.html#PathData |
| `d` property | https://www.w3.org/TR/SVG2/paths.html#TheDProperty |
| moveto commands | https://www.w3.org/TR/SVG2/paths.html#PathDataMovetoCommands |
| closepath command | https://www.w3.org/TR/SVG2/paths.html#PathDataClosePathCommand |
| lineto commands | https://www.w3.org/TR/SVG2/paths.html#PathDataLinetoCommands |
| cubic Bézier commands | https://www.w3.org/TR/SVG2/paths.html#PathDataCubicBezierCommands |
| quadratic Bézier commands | https://www.w3.org/TR/SVG2/paths.html#PathDataQuadraticBezierCommands |
| elliptical arc commands | https://www.w3.org/TR/SVG2/paths.html#PathDataEllipticalArcCommands |
| **The grammar for path data (BNF)** | https://www.w3.org/TR/SVG2/paths.html#PathDataBNF |

## Transforms, viewBox, preserveAspectRatio, units (coords.html) — KEY

Drives `transform.ebnf` and the viewBox/preserveAspectRatio/points mini-grammars.

| Section | Link |
| --- | --- |
| `transform` property (transform-list syntax) | https://www.w3.org/TR/SVG2/coords.html#TransformProperty |
| `viewBox` attribute (4-number syntax) | https://www.w3.org/TR/SVG2/coords.html#ViewBoxAttribute |
| `preserveAspectRatio` attribute (align + meetOrSlice) | https://www.w3.org/TR/SVG2/coords.html#PreserveAspectRatioAttribute |
| Units | https://www.w3.org/TR/SVG2/coords.html#Units |

## Points list (shapes.html)

The `points` attribute syntax for `polyline`/`polygon`.

| Section | Link |
| --- | --- |
| Basic shapes — introduction & definitions | https://www.w3.org/TR/SVG2/shapes.html#Introduction |

## Elliptical-arc parameter conversion (implnote.html)

Reference for arc-flag semantics (informative, but clarifies value ranges).

| Section | Link |
| --- | --- |
| Elliptical arc endpoint syntax | https://www.w3.org/TR/SVG2/implnote.html#ArcSyntax |
| Out-of-range radii correction | https://www.w3.org/TR/SVG2/implnote.html#ArcCorrectionOutOfRangeRadii |

## Presentation-attribute ↔ property duality (styling.html)

Explains that presentation attributes share CSS property value syntax (strong reuse point with proto-css).

| Section | Link |
| --- | --- |
| Presentation attributes | https://www.w3.org/TR/SVG2/styling.html#PresentationAttributes |

## Completeness indices (authoritative name lists)

| Index | Link |
| --- | --- |
| Element Index (Appendix F) | https://www.w3.org/TR/SVG2/eltindex.html |
| Attribute Index — regular (Appendix G.1) | https://www.w3.org/TR/SVG2/attindex.html#RegularAttributes |
| Attribute Index — presentation (Appendix G.2) | https://www.w3.org/TR/SVG2/attindex.html#PresentationAttributes |
| Property Index (Appendix H) | https://www.w3.org/TR/SVG2/propidx.html |

## Companion-spec grammar sources (filters, masking, animation)

Not in SVG 2 core — required for full grammar coverage:

| Module | Grammar-bearing source |
| --- | --- |
| Filter Effects (filter & `fe*` primitives + attrs) | https://www.w3.org/TR/filter-effects-1/ |
| CSS Masking (`clipPath`, `mask`, clip/mask attrs) | https://www.w3.org/TR/css-masking-1/ |
| SMIL Animation (`animate`/`set`/… + timing/value attrs) | https://www.w3.org/TR/SVG11/animate.html · https://svgwg.org/specs/animations/ |
