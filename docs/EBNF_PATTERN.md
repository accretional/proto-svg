# SVG EBNF Structural-Markup Pattern (authoritative)

This grammar is a STRUCTURAL grammar for real SVG documents: the strings it
derives are valid `.svg` markup that a browser renders directly. The literal
markup (`<tag>`, `="`, `"`, `</tag>`, attribute spacing) is therefore part of the
grammar as terminals. This supersedes any earlier "renderer adds punctuation"
note. gluon lifts the leading markup terminals into the prefix map and the
renderer re-emits them, exactly as proto-css bakes CSS punctuation (`{ } : ;`).

## Two hard rules

1. **Enumerate every fixed keyword set in full.** If an attribute or value
   ranges over a defined set of keywords, list ALL of them as quoted terminals,
   however many (16 blend modes, ~80 ARIA roles, all of them). A scalarized leaf
   (`NumberType`, `StringType`, `ColorType`, ...) is ONLY for values that are
   genuinely arbitrary (a real number, a free-text label, author-chosen idents,
   a URL, a script body). Never use a leaf where a closed keyword set exists.

2. **Bake the markup.** Elements and attributes carry their literal SVG syntax.

## Element pattern

```ebnf
Svg = "<svg" , { SvgAttribute } , ">" , { SvgContentModel } , "</svg>" ;
```

- Open tag `"<svg"` (no trailing space). Attribute repetitions follow; each
  attribute supplies its own leading space. Then `">"`, the content repetition,
  and the close tag `"</svg>"`.
- `{ SvgContentModel }` is zero-or-more child elements, so `<g></g>` (empty) and
  `<g><rect.../></g>` both derive. Every SVG element uses the open/close form
  (all SVG elements may legally carry children such as `animate`/`desc`).
- For elements holding character data (`title`, `desc`, `text`, `style`,
  `script`, `metadata`, `foreignObject`), the content model includes the
  `CharacterDataType` leaf (raw text, no quotes/markup).

## Attribute pattern

```ebnf
SvgXAttr          = ' x="' , LengthPercentageType , '"' ;          (* open value -> leaf  *)
SvgZoomAndPanAttr = ' zoomAndPan="' , ( "disable" | "magnify" ) , '"' ;  (* closed set -> terminals *)
SvgViewBoxAttr    = ' viewBox="' , ViewBox , '"' ;                 (* structured value     *)
```

- The leading terminal is a single-quoted string ` name="` (leading SPACE,
  attribute name, `=`, opening double-quote). Single quotes let the terminal
  contain the `"`. This whole token becomes the attribute's prefix-map entry.
- The value is the enumeration / leaf / structured value type.
- The trailing terminal `'"'` closes the quote.
- Names with colons/kebab are literal: `' xml:lang="'`, `' stroke-width="'`,
  `' aria-hidden="'`. No name-mangling.

## Attribute groups and content models

- An attribute group is an alternation of attribute productions:
  `SvgAttribute = SvgXAttr | SvgYAttr | ... | CoreAttribute | AriaAttribute | ...`.
  The element body repeats it: `{ SvgAttribute }`.
- A content model is an alternation of allowed child elements, drawn from the
  content-category unions (`AnimationElement`, `ShapeElement`, ...) when the FULL
  category is allowed, or enumerated element-by-element when only a SUBSET is
  allowed. If an element admits only some members of a category, list exactly
  those members; do not reference the broad category.

## Element rule names are DOM interface names

Each element rule is named by its **DOM interface** (from the spec), so content
models read unambiguously (`ShapeElement = SVGCircleElement | SVGRectElement | …`).
The emitted tag is the literal terminal (`"<rect"`), so the rule name is just a
clear label and never affects output. (gluon lowercases capital runs in the
compiled proto message name, e.g. `SVGRectElement` → `Svgrectelement`; cosmetic
only, names stay collision-unique.) Attribute productions and helper unions keep
their short element-prefixed names (`RectXAttr`, `RectAttribute`, `RectContent`).

Canonical tag → element-rule map:

```
svg=SVGSVGElement  g=SVGGElement  defs=SVGDefsElement  symbol=SVGSymbolElement
use=SVGUseElement  switch=SVGSwitchElement  a=SVGAElement  desc=SVGDescElement
title=SVGTitleElement  metadata=SVGMetadataElement  unknown=SVGUnknownElement
rect=SVGRectElement  circle=SVGCircleElement  ellipse=SVGEllipseElement
line=SVGLineElement  polyline=SVGPolylineElement  polygon=SVGPolygonElement
path=SVGPathElement
text=SVGTextElement  tspan=SVGTSpanElement  textPath=SVGTextPathElement
image=SVGImageElement  foreignObject=SVGForeignObjectElement  view=SVGViewElement
script=SVGScriptElement  style=SVGStyleElement
linearGradient=SVGLinearGradientElement  radialGradient=SVGRadialGradientElement
stop=SVGStopElement  pattern=SVGPatternElement  marker=SVGMarkerElement
clipPath=SVGClipPathElement  mask=SVGMaskElement  filter=SVGFilterElement
feBlend=SVGFEBlendElement  feColorMatrix=SVGFEColorMatrixElement
feComponentTransfer=SVGFEComponentTransferElement  feComposite=SVGFECompositeElement
feConvolveMatrix=SVGFEConvolveMatrixElement  feDiffuseLighting=SVGFEDiffuseLightingElement
feDisplacementMap=SVGFEDisplacementMapElement  feDistantLight=SVGFEDistantLightElement
feDropShadow=SVGFEDropShadowElement  feFlood=SVGFEFloodElement
feFuncR=SVGFEFuncRElement  feFuncG=SVGFEFuncGElement  feFuncB=SVGFEFuncBElement
feFuncA=SVGFEFuncAElement  feGaussianBlur=SVGFEGaussianBlurElement
feImage=SVGFEImageElement  feMerge=SVGFEMergeElement  feMergeNode=SVGFEMergeNodeElement
feMorphology=SVGFEMorphologyElement  feOffset=SVGFEOffsetElement
fePointLight=SVGFEPointLightElement  feSpecularLighting=SVGFESpecularLightingElement
feSpotLight=SVGFESpotLightElement  feTile=SVGFETileElement  feTurbulence=SVGFETurbulenceElement
animate=SVGAnimateElement  set=SVGSetElement  animateMotion=SVGAnimateMotionElement
animateTransform=SVGAnimateTransformElement  mpath=SVGMPathElement  discard=SVGDiscardElement
```

The document root stays `SvgDocument` (`SvgDocument = SVGSVGElement ;`), the
genproto prune root.

## Value grammars carry no markup

Value type files (`datatype.ebnf`, `transform.ebnf`, `path.ebnf`) define the
content that sits INSIDE `="..."`. They never include `<`, `>`, or the attribute
quotes. `TransformList`, `SvgPath`, `Points`, `ViewBox`, `PaintType` are
structured value grammars; `*Type` leaves are scalarized.
