# W3C SVG 2 — Spec Index

Master chapter index for the W3C SVG 2 specification. Base: `https://www.w3.org/TR/SVG2/`

This is the authoritative source for the SVG grammar (element content models, attribute/property value syntax, path-data BNF, transform syntax). Use this file as the map; the grammar-bearing sections are enumerated in the companion docs below.

- [`w3c_svg2_elements.md`](w3c_svg2_elements.md) — per-element defining sections (content model + element-specific attributes)
- [`w3c_svg2_attributes.md`](w3c_svg2_attributes.md) — attribute & presentation-property defining sections
- [`w3c_svg2_grammar.md`](w3c_svg2_grammar.md) — value-syntax sections: data types, path-data BNF, transform/viewBox/preserveAspectRatio, indices

## Chapters

| # | Chapter | Link | Grammar relevance |
| --- | --- | --- | --- |
| 1 | Introduction | https://www.w3.org/TR/SVG2/intro.html | — |
| 2 | Conformance Criteria | https://www.w3.org/TR/SVG2/conform.html | — |
| 3 | Rendering Model | https://www.w3.org/TR/SVG2/render.html | `display`, `visibility`, `opacity`, `overflow` properties |
| 4 | Basic Data Types and Interfaces | https://www.w3.org/TR/SVG2/types.html | **Core** — data types & attribute syntax |
| 5 | Document Structure | https://www.w3.org/TR/SVG2/struct.html | **Core** — svg/g/defs/symbol/use/switch/desc/title/metadata; common attrs; conditional processing; ARIA |
| 6 | Styling | https://www.w3.org/TR/SVG2/styling.html | `style` element; class/style attrs; presentation attributes |
| 7 | Geometry Properties | https://www.w3.org/TR/SVG2/geometry.html | **Core** — cx, cy, r, rx, ry, x, y, width, height |
| 8 | Coordinate Systems, Transformations and Units | https://www.w3.org/TR/SVG2/coords.html | **Core** — transform, viewBox, preserveAspectRatio, units, vector-effect |
| 9 | Paths | https://www.w3.org/TR/SVG2/paths.html | **Core** — path element + path-data BNF |
| 10 | Basic Shapes | https://www.w3.org/TR/SVG2/shapes.html | **Core** — rect/circle/ellipse/line/polyline/polygon; points syntax |
| 11 | Text | https://www.w3.org/TR/SVG2/text.html | **Core** — text/tspan/textPath; text properties |
| 12 | Embedded Content | https://www.w3.org/TR/SVG2/embedded.html | **Core** — image, foreignObject |
| 13 | Painting: Filling, Stroking and Marker Symbols | https://www.w3.org/TR/SVG2/painting.html | **Core** — fill/stroke/marker properties; rendering hints |
| 14 | Paint Servers: Gradients and Patterns | https://www.w3.org/TR/SVG2/pservers.html | **Core** — linearGradient/radialGradient/stop/pattern |
| 15 | Scripting and Interactivity | https://www.w3.org/TR/SVG2/interact.html | `script` element; event attributes; `pointer-events` |
| 16 | Linking | https://www.w3.org/TR/SVG2/linking.html | **Core** — a/view; href & xlink attributes |

## Appendices

| Appendix | Title | Link | Grammar relevance |
| --- | --- | --- | --- |
| A | IDL Definitions | https://www.w3.org/TR/SVG2/idl.html | DOM only (not markup grammar) |
| B | Implementation Notes | https://www.w3.org/TR/SVG2/implnote.html | Elliptical-arc parameter conversion |
| C | Accessibility Support | https://www.w3.org/TR/SVG2/access.html | — |
| D | Animating SVG Documents | https://www.w3.org/TR/SVG2/animate.html | Pointer to animation (see companion specs) |
| E | References | https://www.w3.org/TR/SVG2/refs.html | — |
| **F** | **Element Index** | https://www.w3.org/TR/SVG2/eltindex.html | **Completeness backbone — every element** |
| **G** | **Attribute Index** | https://www.w3.org/TR/SVG2/attindex.html | **Completeness backbone — every attribute** |
| **H** | **Property Index** | https://www.w3.org/TR/SVG2/propidx.html | **Completeness backbone — every property** |
| I | IDL Index | https://www.w3.org/TR/SVG2/idlindex.html | DOM only |
| J | Media Type Registration (image/svg+xml) | https://www.w3.org/TR/SVG2/mimereg.html | — |
| K | Changes from SVG 1.1 | https://www.w3.org/TR/SVG2/changes.html | Useful for SVG1.1-vs-2 deltas |

### Attribute Index subsections (Appendix G)

| Section | Link |
| --- | --- |
| G.1 Regular attributes | https://www.w3.org/TR/SVG2/attindex.html#RegularAttributes |
| G.2 Presentation attributes | https://www.w3.org/TR/SVG2/attindex.html#PresentationAttributes |

## ⚠️ Coverage gap — SVG 2 core does NOT define these modules

The SVG 2 spec above is modularized. Filters, clipping/masking, and declarative animation live in **separate W3C specs**. A grammar that covers "all of SVG" (matching the MDN element/attribute lists) MUST also pull from:

| Module | Elements / attributes it defines | Spec link |
| --- | --- | --- |
| Filter Effects | `filter`, `feBlend`, `feColorMatrix`, `feComponentTransfer`, `feComposite`, `feConvolveMatrix`, `feDiffuseLighting`, `feDisplacementMap`, `feDistantLight`, `feDropShadow`, `feFlood`, `feFuncA/B/G/R`, `feGaussianBlur`, `feImage`, `feMerge`, `feMergeNode`, `feMorphology`, `feOffset`, `fePointLight`, `feSpecularLighting`, `feSpotLight`, `feTile`, `feTurbulence`; `filter`, `flood-color`, `flood-opacity`, `lighting-color`, `color-interpolation-filters`, and all `fe*` primitive attributes | https://www.w3.org/TR/filter-effects-1/ |
| CSS Masking | `clipPath`, `mask`; `clip-path`, `clip-rule`, `clipPathUnits`, `mask`, `mask-type`, `maskUnits`, `maskContentUnits`, `clip` | https://www.w3.org/TR/css-masking-1/ |
| SVG Animation / SMIL | `animate`, `set`, `animateMotion`, `animateTransform`, `mpath`, `discard`; `attributeName`, `begin`, `dur`, `end`, `from`/`to`/`by`/`values`, `keyTimes`, `keySplines`, `calcMode`, `repeatCount`, `repeatDur`, `restart`, `additive`, `accumulate`, `keyPoints`, `rotate`, `path`, `origin` | https://www.w3.org/TR/SVG11/animate.html (SMIL animation; SVG 2 defers to it) + https://svgwg.org/specs/animations/ |

The three SVG 2 appendix indices (F/G/H) only list what SVG 2 core defines, so they will **not** include the filter/masking/animation names above — those come from the indices/sections of the companion specs.
