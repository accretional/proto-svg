# W3C SVG 2 — Attribute & Property Defining Sections

Links to the spec sections that define attribute / presentation-property **value syntax** — the per-attribute data the grammar's `*Attr` / property productions need. Base: `https://www.w3.org/TR/SVG2/`

Authoritative full lists:
- **Appendix G, Attribute Index** — https://www.w3.org/TR/SVG2/attindex.html (regular: `#RegularAttributes`, presentation: `#PresentationAttributes`)
- **Appendix H, Property Index** — https://www.w3.org/TR/SVG2/propidx.html

## Common / core attributes (struct.html)

| Attribute(s) | Link |
| --- | --- |
| Common attributes (overview) | https://www.w3.org/TR/SVG2/struct.html#CommonAttributes |
| `id` | https://www.w3.org/TR/SVG2/struct.html#Core.attrib |
| `lang`, `xml:lang` | https://www.w3.org/TR/SVG2/struct.html#LangSpaceAttrs |
| `xml:space` | https://www.w3.org/TR/SVG2/struct.html#WhitespaceProcessingXMLSpaceAttribute |
| `tabindex` | https://www.w3.org/TR/SVG2/struct.html#tabindexattribute |
| `data-*` | https://www.w3.org/TR/SVG2/struct.html#DataAttributes |

## Conditional processing (struct.html)

| Attribute(s) | Link |
| --- | --- |
| Conditional processing (overview) | https://www.w3.org/TR/SVG2/struct.html#ConditionalProcessing |
| `requiredExtensions` | https://www.w3.org/TR/SVG2/struct.html#ConditionalProcessingRequiredExtensionsAttribute |
| `systemLanguage` | https://www.w3.org/TR/SVG2/struct.html#ConditionalProcessingSystemLanguageAttribute |

## WAI-ARIA (struct.html)

| Attribute(s) | Link |
| --- | --- |
| WAI-ARIA (overview) | https://www.w3.org/TR/SVG2/struct.html#WAIARIAAttributes |
| `role` | https://www.w3.org/TR/SVG2/struct.html#roleattribute |
| `aria-*` (state & property attrs) | https://www.w3.org/TR/SVG2/struct.html#ARIAStateandPropertyAttributes |

## Styling attributes (styling.html)

| Attribute(s) | Link |
| --- | --- |
| `class`, `style` | https://www.w3.org/TR/SVG2/styling.html#ElementSpecificStyling |
| Presentation attributes (how attrs map to properties) | https://www.w3.org/TR/SVG2/styling.html#PresentationAttributes |

## Geometry properties (geometry.html)

| Property | Link |
| --- | --- |
| `cx` | https://www.w3.org/TR/SVG2/geometry.html#CX |
| `cy` | https://www.w3.org/TR/SVG2/geometry.html#CY |
| `r` | https://www.w3.org/TR/SVG2/geometry.html#R |
| `rx` | https://www.w3.org/TR/SVG2/geometry.html#RX |
| `ry` | https://www.w3.org/TR/SVG2/geometry.html#RY |
| `x` | https://www.w3.org/TR/SVG2/geometry.html#X |
| `y` | https://www.w3.org/TR/SVG2/geometry.html#Y |
| `width`, `height` | https://www.w3.org/TR/SVG2/geometry.html#Sizing |

## Coordinate systems & transforms (coords.html)

| Attribute / property | Link |
| --- | --- |
| `transform` | https://www.w3.org/TR/SVG2/coords.html#TransformProperty |
| `viewBox` | https://www.w3.org/TR/SVG2/coords.html#ViewBoxAttribute |
| `preserveAspectRatio` | https://www.w3.org/TR/SVG2/coords.html#PreserveAspectRatioAttribute |
| Units | https://www.w3.org/TR/SVG2/coords.html#Units |
| `vector-effect` | https://www.w3.org/TR/SVG2/coords.html#VectorEffects |

## Rendering-control properties (render.html)

| Property | Link |
| --- | --- |
| `display`, `visibility` | https://www.w3.org/TR/SVG2/render.html#VisibilityControl |
| `opacity` | https://www.w3.org/TR/SVG2/render.html#ObjectAndGroupOpacityProperties |
| `overflow` | https://www.w3.org/TR/SVG2/render.html#OverflowAndClipProperties |

## Painting — fill, stroke, markers, hints (painting.html)

| Property | Link |
| --- | --- |
| `color` | https://www.w3.org/TR/SVG2/painting.html#ColorProperty |
| `fill` | https://www.w3.org/TR/SVG2/painting.html#SpecifyingFillPaint |
| `fill-rule` | https://www.w3.org/TR/SVG2/painting.html#WindingRule |
| `fill-opacity` | https://www.w3.org/TR/SVG2/painting.html#FillOpacity |
| `stroke` | https://www.w3.org/TR/SVG2/painting.html#SpecifyingStrokePaint |
| `stroke-opacity` | https://www.w3.org/TR/SVG2/painting.html#StrokeOpacity |
| `stroke-width` | https://www.w3.org/TR/SVG2/painting.html#StrokeWidth |
| `stroke-linecap` | https://www.w3.org/TR/SVG2/painting.html#LineCaps |
| `stroke-linejoin`, `stroke-miterlimit` | https://www.w3.org/TR/SVG2/painting.html#LineJoin |
| `stroke-dasharray`, `stroke-dashoffset` | https://www.w3.org/TR/SVG2/painting.html#StrokeDashing |
| `marker-start`, `marker-mid`, `marker-end` | https://www.w3.org/TR/SVG2/painting.html#VertexMarkerProperties |
| `marker` (shorthand) | https://www.w3.org/TR/SVG2/painting.html#MarkerShorthand |
| `paint-order` | https://www.w3.org/TR/SVG2/painting.html#PaintOrder |
| `color-interpolation` | https://www.w3.org/TR/SVG2/painting.html#ColorInterpolation |
| `color-rendering` | https://www.w3.org/TR/SVG2/painting.html#ColorRendering |
| `shape-rendering` | https://www.w3.org/TR/SVG2/painting.html#ShapeRendering |
| `text-rendering` | https://www.w3.org/TR/SVG2/painting.html#TextRendering |
| `image-rendering` | https://www.w3.org/TR/SVG2/painting.html#ImageRendering |
| `will-change` | https://www.w3.org/TR/SVG2/painting.html#WillChange |

## Text properties (text.html)

| Property | Link |
| --- | --- |
| `text-anchor` | https://www.w3.org/TR/SVG2/text.html#TextAnchoringProperties |
| `glyph-orientation-horizontal` | https://www.w3.org/TR/SVG2/text.html#GlyphOrientationHorizontalProperty |
| `glyph-orientation-vertical` | https://www.w3.org/TR/SVG2/text.html#GlyphOrientationVerticalProperty |
| `kerning` | https://www.w3.org/TR/SVG2/text.html#KerningProperty |
| `font-variant` | https://www.w3.org/TR/SVG2/text.html#FontVariantProperty |
| `line-height` | https://www.w3.org/TR/SVG2/text.html#LineHeightProperty |
| `writing-mode` | https://www.w3.org/TR/SVG2/text.html#WritingModeProperty |
| `direction` | https://www.w3.org/TR/SVG2/text.html#DirectionProperty |
| `dominant-baseline` | https://www.w3.org/TR/SVG2/text.html#DominantBaselineProperty |
| `alignment-baseline` | https://www.w3.org/TR/SVG2/text.html#AlignmentBaselineProperty |
| `baseline-shift` | https://www.w3.org/TR/SVG2/text.html#BaselineShiftProperty |
| `letter-spacing` | https://www.w3.org/TR/SVG2/text.html#LetterSpacingProperty |
| `word-spacing` | https://www.w3.org/TR/SVG2/text.html#WordSpacingProperty |
| `text-overflow` | https://www.w3.org/TR/SVG2/text.html#TextOverflowProperty |
| `inline-size` | https://www.w3.org/TR/SVG2/text.html#InlineSize |
| `shape-inside` | https://www.w3.org/TR/SVG2/text.html#TextShapeInside |
| `shape-subtract` | https://www.w3.org/TR/SVG2/text.html#TextShapeSubtract |
| `shape-image-threshold` | https://www.w3.org/TR/SVG2/text.html#TextShapeImageThreshold |
| `shape-margin` | https://www.w3.org/TR/SVG2/text.html#TextShapeMargin |
| `shape-padding` | https://www.w3.org/TR/SVG2/text.html#TextShapePadding |
| `white-space` | https://www.w3.org/TR/SVG2/text.html#TextWhiteSpace |
| `text-decoration-fill`, `text-decoration-stroke` | https://www.w3.org/TR/SVG2/text.html#TextDecorationFillStroke |
| `text`/`tspan` positioning attrs (`x`,`y`,`dx`,`dy`,`rotate`,`textLength`,`lengthAdjust`) | https://www.w3.org/TR/SVG2/text.html#TSpanAttributes |
| `textPath` attrs (`href`,`startOffset`,`method`,`spacing`,`side`,`path`) | https://www.w3.org/TR/SVG2/text.html#TextPathAttributes |

## Paint-server attributes (pservers.html)

| Attribute group | Link |
| --- | --- |
| linearGradient attrs (`x1`,`y1`,`x2`,`y2`,`gradientUnits`,`gradientTransform`,`spreadMethod`,`href`) | https://www.w3.org/TR/SVG2/pservers.html#LinearGradientAttributes |
| radialGradient attrs (`cx`,`cy`,`r`,`fx`,`fy`,`fr`, …) | https://www.w3.org/TR/SVG2/pservers.html#RadialGradientAttributes |
| gradient stop attrs (`offset`) | https://www.w3.org/TR/SVG2/pservers.html#GradientStopAttributes |
| `stop-color`, `stop-opacity` | https://www.w3.org/TR/SVG2/pservers.html#StopColorProperties |
| pattern attrs (`patternUnits`,`patternContentUnits`,`patternTransform`,`href`, geometry) | https://www.w3.org/TR/SVG2/pservers.html#PatternElementAttributes |

## Paths (paths.html)

| Attribute / property | Link |
| --- | --- |
| `d` | https://www.w3.org/TR/SVG2/paths.html#TheDProperty |
| `pathLength` | https://www.w3.org/TR/SVG2/paths.html#PathLengthAttribute |

## Interactivity / events (interact.html)

| Attribute(s) | Link |
| --- | --- |
| `pointer-events` | https://www.w3.org/TR/SVG2/interact.html#PointerEventsProp |
| Event attributes (`onclick`, `onload`, …) | https://www.w3.org/TR/SVG2/interact.html#EventAttributes |
| Animation event attributes (`onbegin`, `onend`, `onrepeat`) | https://www.w3.org/TR/SVG2/interact.html#AnimationEvents |

## Linking / URL attributes (linking.html)

| Attribute(s) | Link |
| --- | --- |
| URL reference attributes (`href`) | https://www.w3.org/TR/SVG2/linking.html#linkRefAttrs |
| Deprecated XLink attributes (`xlink:href`, `xlink:title`, …) | https://www.w3.org/TR/SVG2/linking.html#XLinkRefAttrs |
| Syntactic forms: URL and `<url>` | https://www.w3.org/TR/SVG2/linking.html#URLforms |
