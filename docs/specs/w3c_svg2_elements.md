# W3C SVG 2 — Element Defining Sections

Each link points to the spec section that defines an element's **content model** and **element-specific attributes** — the per-element data the grammar's `*Element` productions need. Base: `https://www.w3.org/TR/SVG2/`

Authoritative full list: **Appendix F, Element Index** — https://www.w3.org/TR/SVG2/eltindex.html

## Document structure (struct.html)

| Element | Link |
| --- | --- |
| `<svg>` | https://www.w3.org/TR/SVG2/struct.html#SVGElement |
| `<g>` | https://www.w3.org/TR/SVG2/struct.html#GElement |
| `<defs>` | https://www.w3.org/TR/SVG2/struct.html#DefsElement |
| `<symbol>` | https://www.w3.org/TR/SVG2/struct.html#SymbolElement |
| `<use>` | https://www.w3.org/TR/SVG2/struct.html#UseElement |
| `<switch>` | https://www.w3.org/TR/SVG2/struct.html#SwitchElement |
| `<desc>` | https://www.w3.org/TR/SVG2/struct.html#DescriptionAndTitleElements |
| `<title>` | https://www.w3.org/TR/SVG2/struct.html#DescriptionAndTitleElements |
| `<metadata>` | https://www.w3.org/TR/SVG2/struct.html#MetadataElement |

## Styling (styling.html)

| Element | Link |
| --- | --- |
| `<style>` | https://www.w3.org/TR/SVG2/styling.html#StyleElement |

## Paths (paths.html)

| Element | Link |
| --- | --- |
| `<path>` | https://www.w3.org/TR/SVG2/paths.html#PathElement |

## Basic shapes (shapes.html)

| Element | Link |
| --- | --- |
| `<rect>` | https://www.w3.org/TR/SVG2/shapes.html#RectElement |
| `<circle>` | https://www.w3.org/TR/SVG2/shapes.html#CircleElement |
| `<ellipse>` | https://www.w3.org/TR/SVG2/shapes.html#EllipseElement |
| `<line>` | https://www.w3.org/TR/SVG2/shapes.html#LineElement |
| `<polyline>` | https://www.w3.org/TR/SVG2/shapes.html#PolylineElement |
| `<polygon>` | https://www.w3.org/TR/SVG2/shapes.html#PolygonElement |

## Text (text.html)

| Element | Link |
| --- | --- |
| `<text>` | https://www.w3.org/TR/SVG2/text.html#TextElement |
| `<tspan>` | https://www.w3.org/TR/SVG2/text.html#TextElement |
| `<textPath>` | https://www.w3.org/TR/SVG2/text.html#TextPathElement |

## Embedded content (embedded.html)

| Element | Link |
| --- | --- |
| `<image>` | https://www.w3.org/TR/SVG2/embedded.html#ImageElement |
| `<foreignObject>` | https://www.w3.org/TR/SVG2/embedded.html#ForeignObjectElement |

## Painting — markers (painting.html)

| Element | Link |
| --- | --- |
| `<marker>` | https://www.w3.org/TR/SVG2/painting.html#MarkerElement |

## Paint servers (pservers.html)

| Element | Link |
| --- | --- |
| `<linearGradient>` | https://www.w3.org/TR/SVG2/pservers.html#LinearGradients |
| `<radialGradient>` | https://www.w3.org/TR/SVG2/pservers.html#RadialGradients |
| `<stop>` | https://www.w3.org/TR/SVG2/pservers.html#GradientStops |
| `<pattern>` | https://www.w3.org/TR/SVG2/pservers.html#Patterns |

## Scripting (interact.html)

| Element | Link |
| --- | --- |
| `<script>` | https://www.w3.org/TR/SVG2/interact.html#ScriptElement |

## Linking (linking.html)

| Element | Link |
| --- | --- |
| `<a>` | https://www.w3.org/TR/SVG2/linking.html#Links |
| `<view>` | https://www.w3.org/TR/SVG2/linking.html#ViewElement |

## ⚠️ Elements NOT in SVG 2 core (see companion specs)

These are referenced by MDN but defined elsewhere — needed for a complete grammar:

| Element(s) | Spec |
| --- | --- |
| `filter`, `feBlend`, `feColorMatrix`, `feComponentTransfer`, `feComposite`, `feConvolveMatrix`, `feDiffuseLighting`, `feDisplacementMap`, `feDistantLight`, `feDropShadow`, `feFlood`, `feFuncA/B/G/R`, `feGaussianBlur`, `feImage`, `feMerge`, `feMergeNode`, `feMorphology`, `feOffset`, `fePointLight`, `feSpecularLighting`, `feSpotLight`, `feTile`, `feTurbulence` | https://www.w3.org/TR/filter-effects-1/ |
| `clipPath`, `mask` | https://www.w3.org/TR/css-masking-1/ |
| `animate`, `set`, `animateMotion`, `animateTransform`, `mpath` | https://www.w3.org/TR/SVG11/animate.html · https://svgwg.org/specs/animations/ |
