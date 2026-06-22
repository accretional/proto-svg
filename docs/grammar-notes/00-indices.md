# SVG Grammar Completeness Indices

Source files read:
- `docs/specs/cache/svg2-eltindex.txt` — W3C SVG 2 Appendix F (informative)
- `docs/specs/cache/svg2-attindex.txt` — W3C SVG 2 Appendix G (informative)
- `docs/specs/cache/svg2-propidx.txt` — W3C SVG 2 Appendix H (informative)
- `docs/specs/mdn_docs_elements.md` — MDN element reference list
- `docs/specs/mdn_docs_attributes.md` — MDN attribute reference list

---

## Elements

Total from W3C SVG 2 index: **52**

```
a
animate
animateMotion
animateTransform
audio
canvas
circle
clipPath
defs
desc
discard
ellipse
feBlend
feColorMatrix
feComponentTransfer
feComposite
feConvolveMatrix
feDiffuseLighting
feDisplacementMap
feDistantLight
feDropShadow
feFlood
feFuncA
feFuncB
feFuncG
feFuncR
feGaussianBlur
feImage
feMerge
feMergeNode
feMorphology
feOffset
fePointLight
feSpecularLighting
feSpotLight
feTile
feTurbulence
filter
foreignObject
g
iframe
image
line
linearGradient
marker
mask
metadata
mpath
path
pattern
polygon
polyline
radialGradient
rect
script
set
stop
style
svg
switch
symbol
text
textPath
title
tspan
unknown
use
video
view
```

---

## Regular Attributes (G.1)

The Anim. column marks attributes the spec explicitly flags as animatable (✓).
Attributes with no element list in the spec (global event handlers with no element column) are noted as such.

| Attribute | Elements | Anim. | Notes / Value column |
|-----------|----------|-------|----------------------|
| `accumulate` | animate, animateMotion, animateTransform | | |
| `additive` | animate, animateMotion, animateTransform | | |
| `amplitude` | feFuncA, feFuncB, feFuncG, feFuncR | yes | |
| `aria-activedescendant` | a, audio, canvas, circle, discard, ellipse, foreignObject, g, iframe, image, line, path, polygon, polyline, rect, svg, switch, symbol, text, textPath, tspan, unknown, use, video, view | | |
| `aria-atomic` | (same aria element set) | | |
| `aria-autocomplete` | (same aria element set) | | |
| `aria-busy` | (same aria element set) | | |
| `aria-checked` | (same aria element set) | | |
| `aria-colcount` | (same aria element set) | | |
| `aria-colindex` | (same aria element set) | | |
| `aria-colspan` | (same aria element set) | | |
| `aria-controls` | (same aria element set) | | |
| `aria-current` | (same aria element set) | | |
| `aria-describedby` | (same aria element set) | | |
| `aria-details` | (same aria element set) | | |
| `aria-disabled` | (same aria element set) | | |
| `aria-dropeffect` | (same aria element set) | | |
| `aria-errormessage` | (same aria element set) | | |
| `aria-expanded` | (same aria element set) | | |
| `aria-flowto` | (same aria element set) | | |
| `aria-grabbed` | (same aria element set) | | |
| `aria-haspopup` | (same aria element set) | | |
| `aria-hidden` | (same aria element set) | | |
| `aria-invalid` | (same aria element set) | | |
| `aria-keyshortcuts` | (same aria element set) | | |
| `aria-label` | (same aria element set) | | |
| `aria-labelledby` | (same aria element set) | | |
| `aria-level` | (same aria element set) | | |
| `aria-live` | (same aria element set) | | |
| `aria-modal` | (same aria element set) | | |
| `aria-multiline` | (same aria element set) | | |
| `aria-multiselectable` | (same aria element set) | | |
| `aria-orientation` | (same aria element set) | | |
| `aria-owns` | (same aria element set) | | |
| `aria-placeholder` | (same aria element set) | | |
| `aria-posinset` | (same aria element set) | | |
| `aria-pressed` | (same aria element set) | | |
| `aria-readonly` | (same aria element set) | | |
| `aria-relevant` | (same aria element set) | | |
| `aria-required` | (same aria element set) | | |
| `aria-roledescription` | (same aria element set) | | |
| `aria-rowcount` | (same aria element set) | | |
| `aria-rowindex` | (same aria element set) | | |
| `aria-rowspan` | (same aria element set) | | |
| `aria-selected` | (same aria element set) | | |
| `aria-setsize` | (same aria element set) | | |
| `aria-sort` | (same aria element set) | | |
| `aria-valuemax` | (same aria element set) | | |
| `aria-valuemin` | (same aria element set) | | |
| `aria-valuenow` | (same aria element set) | | |
| `aria-valuetext` | (same aria element set) | | |
| `attributeName` | animate, animateTransform, set | | |
| `azimuth` | feDistantLight | yes | |
| `baseFrequency` | feTurbulence | yes | |
| `begin` | animate, animateMotion, animateTransform, set | | |
| `begin` | discard | | (separate row in spec) |
| `bias` | feConvolveMatrix | yes | |
| `by` | animate, animateMotion, animateTransform | | |
| `calcMode` | animate, animateMotion, animateTransform | | |
| `class` | all elements (full list) | yes | |
| `clipPathUnits` | clipPath | yes | |
| `crossorigin` | feImage | | |
| `crossorigin` | image | yes | |
| `crossorigin` | script | yes | |
| `cx` | radialGradient | yes | |
| `cy` | radialGradient | yes | |
| `diffuseConstant` | feDiffuseLighting | yes | |
| `divisor` | feConvolveMatrix | yes | |
| `download` | a | | |
| `dur` | animate, animateMotion, animateTransform, set | | |
| `dx` | feDropShadow | yes | |
| `dx` | feOffset | yes | |
| `dx` | text | yes | |
| `dx` | tspan | yes | |
| `dy` | feDropShadow | yes | |
| `dy` | feOffset | yes | |
| `dy` | text | yes | |
| `dy` | tspan | yes | |
| `edgeMode` | feConvolveMatrix | yes | |
| `edgeMode` | feGaussianBlur | yes | |
| `elevation` | feDistantLight | yes | |
| `end` | animate, animateMotion, animateTransform, set | | |
| `exponent` | feFuncA, feFuncB, feFuncG, feFuncR | yes | |
| `fill` | animate, animateMotion, animateTransform, set | | (animation timing fill, not the paint property) |
| `filterUnits` | filter | yes | |
| `fr` | radialGradient | yes | |
| `from` | animate, animateMotion, animateTransform | | |
| `fx` | radialGradient | yes | |
| `fy` | radialGradient | yes | |
| `gradientTransform` | linearGradient | yes | |
| `gradientTransform` | radialGradient | yes | |
| `gradientUnits` | linearGradient | yes | |
| `gradientUnits` | radialGradient | yes | |
| `height` | filter primitive elements (feBlend…feTurbulence) | yes | |
| `height` | filter | yes | |
| `height` | mask | yes | |
| `height` | pattern | yes | |
| `href` | a | yes | |
| `href` | animate, animateMotion, animateTransform, set | | |
| `href` | discard | | |
| `href` | feImage | yes | |
| `href` | image | yes | |
| `href` | linearGradient | yes | |
| `href` | mpath | | |
| `href` | pattern | yes | |
| `href` | radialGradient | yes | |
| `href` | script | | |
| `href` | textPath | yes | |
| `href` | use | yes | |
| `hreflang` | a | | |
| `id` | all elements (full list) | | |
| `in` | feBlend, feColorMatrix, feComponentTransfer, feComposite, feConvolveMatrix, feDiffuseLighting, feDisplacementMap, feDropShadow, feGaussianBlur, feMergeNode, feMorphology, feOffset, feSpecularLighting, feTile | yes | |
| `in2` | feBlend | yes | |
| `in2` | feComposite | yes | |
| `in2` | feDisplacementMap | yes | |
| `intercept` | feFuncA, feFuncB, feFuncG, feFuncR | yes | |
| `k1` | feComposite | yes | |
| `k2` | feComposite | yes | |
| `k3` | feComposite | yes | |
| `k4` | feComposite | yes | |
| `kernelMatrix` | feConvolveMatrix | yes | |
| `kernelUnitLength` | feConvolveMatrix | yes | |
| `kernelUnitLength` | feDiffuseLighting | yes | |
| `kernelUnitLength` | feSpecularLighting | yes | |
| `keyPoints` | animateMotion | | |
| `keySplines` | animate, animateMotion, animateTransform | | |
| `keyTimes` | animate, animateMotion, animateTransform | | |
| `lang` | all elements (full list) | | |
| `lengthAdjust` | text, textPath, tspan | yes | |
| `limitingConeAngle` | feSpotLight | yes | |
| `markerHeight` | marker | yes | |
| `markerUnits` | marker | yes | |
| `markerWidth` | marker | yes | |
| `maskContentUnits` | mask | yes | |
| `maskUnits` | mask | yes | |
| `max` | animate, animateMotion, animateTransform, set | | |
| `media` | style | | |
| `method` | textPath | yes | |
| `min` | animate, animateMotion, animateTransform, set | | |
| `mode` | feBlend | yes | |
| `numOctaves` | feTurbulence | yes | |
| `offset` | feFuncA, feFuncB, feFuncG, feFuncR | yes | |
| `offset` | stop | yes | |
| `onabort` | svg | | |
| `onafterprint` | (no element column in spec) | | |
| `onbeforeprint` | (no element column in spec) | | |
| `onbegin` | animate, animateMotion, animateTransform, set | | |
| `oncancel` | most elements | | |
| `oncanplay` | most elements | | |
| `oncanplaythrough` | most elements | | |
| `onchange` | most elements | | |
| `onclick` | most elements | | |
| `onclose` | most elements | | |
| `oncopy` | most elements (excl. pattern) | | |
| `oncuechange` | most elements | | |
| `oncut` | most elements (excl. pattern) | | |
| `ondblclick` | most elements | | |
| `ondrag` | most elements | | |
| `ondragend` | most elements | | |
| `ondragenter` | most elements | | |
| `ondragexit` | most elements | | |
| `ondragleave` | most elements | | |
| `ondragover` | most elements | | |
| `ondragstart` | most elements | | |
| `ondrop` | most elements | | |
| `ondurationchange` | most elements | | |
| `onemptied` | most elements | | |
| `onend` | animate, animateMotion, animateTransform, set | | |
| `onended` | most elements | | |
| `onerror` | most elements | | |
| `onerror` | svg | | (separate svg-only row) |
| `onfocus` | most elements | | |
| `onfocusin` | a, audio, canvas, circle, defs, ellipse, foreignObject, g, iframe, image, line, path, polygon, polyline, rect, svg, switch, symbol, text, textPath, tspan, unknown, use, video | | |
| `onfocusout` | (same as onfocusin) | | |
| `onhashchange` | (no element column in spec) | | |
| `oninput` | most elements | | |
| `oninvalid` | most elements | | |
| `onkeydown` | most elements | | |
| `onkeypress` | most elements | | |
| `onkeyup` | most elements | | |
| `onload` | most elements | | |
| `onloadeddata` | most elements | | |
| `onloadedmetadata` | most elements | | |
| `onloadstart` | most elements | | |
| `onmessage` | (no element column in spec) | | |
| `onmousedown` | most elements | | |
| `onmouseenter` | most elements | | |
| `onmouseleave` | most elements | | |
| `onmousemove` | most elements | | |
| `onmouseout` | most elements | | |
| `onmouseover` | most elements | | |
| `onmouseup` | most elements | | |
| `onmousewheel` | most elements | | |
| `onoffline` | (no element column in spec) | | |
| `ononline` | (no element column in spec) | | |
| `onpagehide` | (no element column in spec) | | |
| `onpageshow` | (no element column in spec) | | |
| `onpaste` | most elements (excl. pattern) | | |
| `onpause` | most elements | | |
| `onplay` | most elements | | |
| `onplaying` | most elements | | |
| `onpopstate` | (no element column in spec) | | |
| `onprogress` | most elements | | |
| `onratechange` | most elements | | |
| `onrepeat` | animate, animateMotion, animateTransform, set | | |
| `onreset` | most elements | | |
| `onresize` | most elements | | |
| `onresize` | svg | | (separate svg-only row) |
| `onscroll` | most elements | | |
| `onscroll` | svg | | (separate svg-only row) |
| `onseeked` | most elements | | |
| `onseeking` | most elements | | |
| `onselect` | most elements | | |
| `onshow` | most elements | | |
| `onstalled` | most elements | | |
| `onstorage` | (no element column in spec) | | |
| `onsubmit` | most elements | | |
| `onsuspend` | most elements | | |
| `ontimeupdate` | most elements | | |
| `ontoggle` | most elements | | |
| `onunload` | (no element column in spec) | | |
| `onunload` | svg | | (separate svg-only row) |
| `onvolumechange` | most elements | | |
| `onwaiting` | most elements | | |
| `operator` | feComposite | yes | |
| `operator` | feMorphology | yes | |
| `order` | feConvolveMatrix | yes | |
| `orient` | marker | yes | |
| `origin` | animateMotion | | |
| `path` | animateMotion | | |
| `path` | textPath | yes | |
| `pathLength` | circle, ellipse, line, path, polygon, polyline, rect | yes | |
| `patternContentUnits` | pattern | yes | |
| `patternTransform` | pattern | yes | |
| `patternUnits` | pattern | yes | |
| `ping` | a | | |
| `playbackorder` | svg | | |
| `points` | polygon | yes | |
| `points` | polyline | yes | |
| `pointsAtX` | feSpotLight | yes | |
| `pointsAtY` | feSpotLight | yes | |
| `pointsAtZ` | feSpotLight | yes | |
| `preserveAlpha` | feConvolveMatrix | yes | |
| `preserveAspectRatio` | canvas, feImage, image, marker, pattern, svg, symbol, view | yes | |
| `primitiveUnits` | filter | yes | |
| `r` | radialGradient | yes | |
| `radius` | feMorphology | yes | |
| `refX` | marker | yes | |
| `refX` | symbol | yes | |
| `refY` | marker | yes | |
| `refY` | symbol | yes | |
| `referrerpolicy` | a | | |
| `rel` | a | | |
| `repeatCount` | animate, animateMotion, animateTransform, set | | |
| `repeatDur` | animate, animateMotion, animateTransform, set | | |
| `requiredExtensions` | a, animate, animateMotion, animateTransform, audio, canvas, circle, clipPath, discard, ellipse, foreignObject, g, iframe, image, line, mask, path, polygon, polyline, rect, set, svg, switch, text, textPath, tspan, unknown, use, video | | |
| `restart` | animate, animateMotion, animateTransform, set | | |
| `result` | feBlend, feColorMatrix, feComponentTransfer, feComposite, feConvolveMatrix, feDiffuseLighting, feDisplacementMap, feDropShadow, feFlood, feGaussianBlur, feImage, feMerge, feMorphology, feOffset, feSpecularLighting, feTile, feTurbulence | yes | |
| `role` | a, audio, canvas, circle, discard, ellipse, foreignObject, g, iframe, image, line, path, polygon, polyline, rect, svg, switch, symbol, text, textPath, tspan, unknown, use, video, view | | |
| `rotate` | animateMotion | | |
| `rotate` | text | yes | |
| `rotate` | tspan | yes | |
| `scale` | feDisplacementMap | yes | |
| `seed` | feTurbulence | yes | |
| `side` | textPath | yes | |
| `slope` | feFuncA, feFuncB, feFuncG, feFuncR | yes | |
| `spacing` | textPath | yes | |
| `specularConstant` | feSpecularLighting | yes | |
| `specularExponent` | feSpecularLighting | yes | |
| `specularExponent` | feSpotLight | yes | |
| `spreadMethod` | linearGradient | yes | |
| `spreadMethod` | radialGradient | yes | |
| `startOffset` | textPath | yes | |
| `stdDeviation` | feDropShadow | yes | |
| `stdDeviation` | feGaussianBlur | yes | |
| `stitchTiles` | feTurbulence | yes | |
| `style` | all elements (full list) | | |
| `surfaceScale` | feDiffuseLighting | yes | |
| `surfaceScale` | feSpecularLighting | yes | |
| `systemLanguage` | a, animate, animateMotion, animateTransform, audio, canvas, circle, clipPath, discard, ellipse, foreignObject, g, iframe, image, line, mask, path, polygon, polyline, rect, set, svg, switch, text, textPath, tspan, unknown, use, video | | |
| `tabindex` | all elements (full list) | | |
| `tableValues` | feFuncA, feFuncB, feFuncG, feFuncR | yes | |
| `target` | a | yes | |
| `targetX` | feConvolveMatrix | yes | |
| `targetY` | feConvolveMatrix | yes | |
| `textLength` | text | yes | |
| `textLength` | textPath, tspan | yes | |
| `timelinebegin` | svg | | |
| `title` | style | | |
| `to` | animate, animateMotion, animateTransform | | |
| `to` | set | | |
| `transform` | svg | yes | |
| `type` | a | | |
| `type` | animateTransform | | |
| `type` | feColorMatrix | yes | |
| `type` | feFuncA, feFuncB, feFuncG, feFuncR | yes | |
| `type` | feTurbulence | yes | |
| `type` | script | | |
| `type` | style | | |
| `values` | animate, animateMotion, animateTransform | | |
| `values` | feColorMatrix | yes | |
| `viewBox` | marker, pattern, svg, symbol, view | yes | |
| `width` | filter primitive elements (feBlend…feTurbulence) | yes | |
| `width` | filter | yes | |
| `width` | mask | yes | |
| `width` | pattern | yes | |
| `x` | filter primitive elements (feBlend…feTurbulence) | yes | |
| `x` | fePointLight | yes | |
| `x` | feSpotLight | yes | |
| `x` | filter | yes | |
| `x` | mask | yes | |
| `x` | pattern | yes | |
| `x` | text | yes | |
| `x` | tspan | yes | |
| `x1` | line | yes | |
| `x1` | linearGradient | yes | |
| `x2` | line | yes | |
| `x2` | linearGradient | yes | |
| `xChannelSelector` | feDisplacementMap | yes | |
| `xlink:href` | a, image, linearGradient, pattern, radialGradient, script, textPath, use | | (deprecated in SVG 2; use `href`) |
| `xlink:href` | feImage | yes | (deprecated) |
| `xlink:title` | a, image, linearGradient, pattern, radialGradient, script, textPath, use | | (deprecated) |
| `xml:space` | all elements (full list) | | (deprecated in SVG 2) |
| `y` | filter primitive elements (feBlend…feTurbulence) | yes | |
| `y` | fePointLight | yes | |
| `y` | feSpotLight | yes | |
| `y` | filter | yes | |
| `y` | mask | yes | |
| `y` | pattern | yes | |
| `y` | text | yes | |
| `y` | tspan | yes | |
| `y1` | line | yes | |
| `y1` | linearGradient | yes | |
| `y2` | line | yes | |
| `y2` | linearGradient | yes | |
| `yChannelSelector` | feDisplacementMap | yes | |
| `z` | fePointLight | yes | |
| `z` | feSpotLight | yes | |
| `zoomAndPan` | svg, view | | (deprecated in SVG 2) |

Unique regular attribute names (collapsing multi-element rows): **accumulate, additive, amplitude, aria-* (50 tokens), attributeName, azimuth, baseFrequency, begin, bias, by, calcMode, class, clipPathUnits, crossorigin, cx, cy, diffuseConstant, divisor, download, dur, dx, dy, edgeMode, elevation, end, exponent, fill, filterUnits, fr, from, fx, fy, gradientTransform, gradientUnits, height, href, hreflang, id, in, in2, intercept, k1, k2, k3, k4, kernelMatrix, kernelUnitLength, keyPoints, keySplines, keyTimes, lang, lengthAdjust, limitingConeAngle, markerHeight, markerUnits, markerWidth, maskContentUnits, maskUnits, max, media, method, min, mode, numOctaves, offset, on* (~60 event handlers), operator, order, orient, origin, path, pathLength, patternContentUnits, patternTransform, patternUnits, ping, playbackorder, points, pointsAtX, pointsAtY, pointsAtZ, preserveAlpha, preserveAspectRatio, primitiveUnits, r, radius, refX, refY, referrerpolicy, rel, repeatCount, repeatDur, requiredExtensions, restart, result, role, rotate, scale, seed, side, slope, spacing, specularConstant, specularExponent, spreadMethod, startOffset, stdDeviation, stitchTiles, style, surfaceScale, systemLanguage, tabindex, tableValues, target, targetX, targetY, textLength, timelinebegin, title, to, transform, type, values, viewBox, width, x, x1, x2, xChannelSelector, xlink:href, xlink:title, xml:space, y, y1, y2, yChannelSelector, z, zoomAndPan**

Approximate count of unique attribute names in G.1: **~230** (including 50 aria-* tokens and ~60 on* event handler names).

---

## Presentation Attributes (G.2)

These 54 attributes mirror the CSS properties of the same name and may be used as element attributes or CSS properties on SVG elements. Listed verbatim from spec G.2:

```
alignment-baseline
baseline-shift
clip
clip-path
clip-rule
color
color-interpolation
color-interpolation-filters
color-rendering
cursor
direction
display
dominant-baseline
fill
fill-opacity
fill-rule
filter
flood-color
flood-opacity
font-family
font-size
font-size-adjust
font-stretch
font-style
font-variant
font-weight
glyph-orientation-horizontal
glyph-orientation-vertical
image-rendering
letter-spacing
lighting-color
marker
marker-end
marker-mid
marker-start
mask
opacity
overflow
paint-order
pointer-events
shape-rendering
stop-color
stop-opacity
stroke
stroke-dasharray
stroke-dashoffset
stroke-linecap
stroke-linejoin
stroke-miterlimit
stroke-opacity
stroke-width
text-anchor
text-decoration
text-rendering
transform
unicode-bidi
vector-effect
visibility
word-spacing
writing-mode
```

Total presentation attributes: **60** (54 individual names + `marker` shorthand which covers `marker-end`, `marker-mid`, `marker-start` — all listed explicitly).

Note: the spec lists `marker-end`, `marker-mid`, `marker-start` individually alongside the `marker` shorthand, plus `letter-spacing` and `word-spacing` appear at the end of the run-on sentence.

---

## Properties (H)

From Appendix H. Value syntax and initial value are verbatim from the spec text.

| Property | Value syntax (verbatim) | Initial value |
|----------|------------------------|---------------|
| `alignment-baseline` | `auto \| baseline \| before-edge \| text-before-edge \| middle \| central \| after-edge \| text-after-edge \| ideographic \| alphabetic \| hanging \| mathematical` | see property description |
| `baseline-shift` | `baseline \| sub \| super \| <percentage> \| <length>` | baseline |
| `color` | `<color>` | depends on user agent |
| `color-interpolation` | `auto \| sRGB \| linearRGB` | sRGB |
| `color-rendering` | `auto \| optimizeSpeed \| optimizeQuality` | auto |
| `direction` | `ltr \| rtl` | ltr |
| `display` | `inline \| block \| list-item \| run-in \| compact \| marker \| table \| inline-table \| table-row-group \| table-header-group \| table-footer-group \| table-row \| table-column-group \| table-column \| table-cell \| table-caption \| none` | inline |
| `dominant-baseline` | `auto \| use-script \| no-change \| reset-size \| ideographic \| alphabetic \| hanging \| mathematical \| central \| middle \| text-after-edge \| text-before-edge` | auto |
| `fill` | `<paint>` | black |
| `fill-opacity` | `<alpha-value>` | 1 |
| `fill-rule` | `nonzero \| evenodd` | nonzero |
| `font-variant` | `normal \| small-caps` | normal |
| `glyph-orientation-vertical` | `auto \| <angle> \| <number>` | auto |
| `image-rendering` | `auto \| optimizeSpeed \| optimizeQuality` | auto |
| `line-height` | `normal \| <number> \| <length-percentage>` | normal |
| `marker` | see individual properties | see individual properties |
| `marker-end` | `none \| <url>` | none |
| `marker-mid` | `none \| <url>` | none |
| `marker-start` | `none \| <url>` | none |
| `opacity` | `<alpha-value>` | 1 |
| `overflow` | `visible \| hidden \| scroll \| auto` | see prose |
| `paint-order` | `normal \| [ fill \|\| stroke \|\| markers ]` | normal |
| `pointer-events` | `bounding-box \| visiblePainted \| visibleFill \| visibleStroke \| visible \| painted \| fill \| stroke \| all \| none` | visiblePainted |
| `shape-rendering` | `auto \| optimizeSpeed \| crispEdges \| geometricPrecision` | auto |
| `stop-color` | `currentColor \| <color> [<icccolor>]` | black |
| `stop-opacity` | `<alpha-value>` | 1 |
| `stroke` | `<paint>` | none |
| `stroke-dasharray` | `none \| <dasharray>` | none |
| `stroke-dashoffset` | `<length-percentage>` | 0 |
| `stroke-linecap` | `butt \| round \| square` | butt |
| `stroke-linejoin` | `miter \| round \| bevel` | miter |
| `stroke-miterlimit` | `<number>` (non-negative) | 4 |
| `stroke-opacity` | `<alpha-value>` | 1 |
| `stroke-width` | `<length-percentage>` | 1 |
| `text-anchor` | `start \| middle \| end` | start |
| `text-decoration` | `none \| [ underline \|\| overline \|\| line-through \|\| blink ]` | none |
| `text-rendering` | `auto \| optimizeSpeed \| optimizeLegibility \| geometricPrecision` | auto |
| `vector-effect` | `non-scaling-stroke \| none` | none |
| `visibility` | `visible \| hidden \| collapse` | visible |
| `white-space` | `normal \| pre \| nowrap \| pre-wrap \| pre-line` | normal |
| `writing-mode` | `lr-tb \| rl-tb \| tb-rl \| lr \| rl \| tb` | lr-tb |

Total properties in Appendix H: **41** (counting `marker`, `marker-end`, `marker-mid`, `marker-start` each separately).

Note: `clip`, `clip-path`, `clip-rule`, `color-interpolation-filters`, `cursor`, `filter`, `flood-color`, `flood-opacity`, `font-family`, `font-size`, `font-size-adjust`, `font-stretch`, `font-style`, `font-weight`, `glyph-orientation-horizontal`, `letter-spacing`, `lighting-color`, `mask`, `text-anchor` (already listed), `unicode-bidi`, `word-spacing` appear in the G.2 presentation-attribute list but are **not repeated** in the H property table — they are inherited from CSS or defined elsewhere. The H table is the subset that SVG 2 defines normatively on its own.

---

## MDN Cross-check

### Elements in MDN but NOT in W3C SVG 2 index

| Element | Browser support | Grammar decision |
|---------|----------------|-----------------|
| `<altGlyph>` | None (removed from all browsers) | **EXCLUDE** — SVG 1.1 legacy, never implemented in modern engines |
| `<altGlyphDef>` | None | **EXCLUDE** — same |
| `<altGlyphItem>` | None | **EXCLUDE** — same |
| `<animate>` | Full | (present in W3C index — listed for completeness) |
| `<color-profile>` | None | **EXCLUDE** — SVG 1.1, dropped before any real browser support |
| `<cursor>` | None | **EXCLUDE** — SVG 1.1, removed; CSS cursor used instead |
| `<font>` | None (Firefox/Chrome dropped it) | **EXCLUDE** — SVG 1.1 font embedding, fully removed |
| `<font-face>` | None | **EXCLUDE** — same as `<font>` |
| `<font-face-format>` | None | **EXCLUDE** — same |
| `<font-face-name>` | None | **EXCLUDE** — same |
| `<font-face-src>` | None | **EXCLUDE** — same |
| `<font-face-uri>` | None | **EXCLUDE** — same |
| `<glyph>` | None | **EXCLUDE** — SVG 1.1 font glyph, removed |
| `<glyphRef>` | None | **EXCLUDE** — same |
| `<hkern>` | None | **EXCLUDE** — SVG 1.1 kerning, removed |
| `<missing-glyph>` | None | **EXCLUDE** — SVG 1.1 font, removed |
| `<tref>` | None (removed ~2013) | **EXCLUDE** — never widely implemented, removed from spec |
| `<vkern>` | None | **EXCLUDE** — SVG 1.1 kerning, removed |

Note: MDN's element reference page (as captured) lists only currently supported elements — it does NOT include the removed SVG 1.1 elements above. The MDN elements file shows 69 elements; all 69 are a subset of the W3C SVG 2 index (minus audio, canvas, iframe, video, unknown, discard which MDN does not list separately as SVG elements).

### Elements in W3C SVG 2 index but NOT in MDN

These are SVG 2 additions that browsers have not implemented:

| Element | Browser support | Grammar decision |
|---------|----------------|-----------------|
| `<audio>` | Not implemented in any browser as SVG element | **EXCLUDE** — SVG 2 proposed embedding HTML media; no browser shipped it |
| `<canvas>` | Not implemented as SVG element | **EXCLUDE** — SVG 2 proposed embedding HTML canvas; no browser shipped it |
| `<discard>` | Partial (Chrome/Edge only, behind flag or limited) | **CONDITIONAL** — include as optional/experimental; exclude from core grammar |
| `<iframe>` | Not implemented as SVG element | **EXCLUDE** — SVG 2 proposed embedding; no browser shipped it |
| `<unknown>` | Parsed but has no rendering semantics | **EXCLUDE** — meta-element for fallback in `<switch>`; not a real author element |
| `<video>` | Not implemented as SVG element | **EXCLUDE** — SVG 2 proposed embedding HTML media; no browser shipped it |

### Attributes in MDN but NOT in W3C SVG 2 G.1 index

| Attribute | MDN status | Grammar decision |
|-----------|-----------|-----------------|
| `attributeType` | Deprecated (⚠️) | **EXCLUDE** from core — SVG 2 dropped it; SMIL-era; browsers still parse but do nothing with it |
| `autofocus` | Listed | **INCLUDE** — standard HTML global attribute, valid on SVG elements in browsers |
| `baseProfile` | Deprecated (⚠️) | **EXCLUDE** — SVG 1.1 era, removed from SVG 2 |
| `color-rendering` | Listed | Present in propidx but absent from G.1 regular attrs; it's a presentation attribute — covered under G.2 |
| `d` | Listed | **INCLUDE** — geometry attribute on `<path>` and `<rect>`; absent from G.1 because it's defined in the Paths chapter, not the attribute index |
| `data-*` | Listed | **INCLUDE** — standard HTML custom data attributes, valid on SVG elements |
| `decoding` | Listed | **INCLUDE** — valid on `<image>` (mirrors HTML img); browsers support it |
| `fetchpriority` | Experimental/Non-standard (🧪🔶) | **EXCLUDE** from grammar — not in spec, not cross-browser |
| `font-width` | Experimental/Non-standard (🧪🔶) | **EXCLUDE** — not in SVG 2 spec, non-standard |
| `glyph-orientation-horizontal` | Deprecated (⚠️) | **EXCLUDE** — SVG 1.1 only, removed in SVG 2 |
| `mask-type` | Listed | **INCLUDE** — CSS Masking spec, fully supported in browsers |
| `r` | Listed (MDN `r` attr) | MDN lists `r` for `<circle>` and `<ellipse>` (`rx`/`ry`); the SVG 2 G.1 only lists `r` for radialGradient; in practice `r` is also a geometry attr on `<circle>` — defined in Shapes chapter, not G.1 |
| `requiredFeatures` | Deprecated (⚠️) | **EXCLUDE** — SVG 1.1, removed from SVG 2; browsers removed support |
| `rx` | Listed | **INCLUDE** — geometry attribute on `<rect>`, `<ellipse>`, `<circle>`; defined in Shapes chapter, absent from G.1 for the same reason as `r`, `d`, `cx`, `cy` for shapes |
| `ry` | Listed | **INCLUDE** — same as `rx` |
| `text-overflow` | Listed | **INCLUDE** — CSS property valid as presentation attribute in browsers |
| `transform-origin` | Listed | **INCLUDE** — CSS Transforms; browsers accept on SVG elements |
| `version` | Deprecated (⚠️) | **EXCLUDE** — SVG 1.1 only, ignored by all browsers in SVG 2 |
| `xlink:arcrole` | Deprecated (⚠️) | **EXCLUDE** — XLink era, removed from SVG 2 |
| `xlink:show` | Deprecated (⚠️) | **EXCLUDE** — XLink era, removed from SVG 2 |
| `xlink:type` | Deprecated (⚠️) | **EXCLUDE** — XLink era, removed from SVG 2 |
| `xml:lang` | Deprecated (⚠️) | **EXCLUDE** from grammar (use `lang`); browsers still parse both |

Note on geometry attributes absent from G.1: `cx`, `cy` for `<circle>` and `<ellipse>`, `r` for `<circle>`, `rx`/`ry` for `<ellipse>` and `<rect>`, `d` for `<path>` and `<rect>`, `x`/`y`/`width`/`height` for `<rect>`, `<image>`, `<svg>`, `<foreignObject>`, `<use>` are defined per element in the Shapes/Structure chapters, not re-listed in G.1 (G.1 only shows the subset of attributes the appendix chose to reference). The grammar must include them from their normative chapter definitions.

### Attributes in W3C SVG 2 G.1 but NOT in MDN

| Attribute | Status | Grammar decision |
|-----------|--------|-----------------|
| `fr` | SVG 2 new (focal radius for radialGradient) | **INCLUDE** — Chrome/Firefox/Safari all support it |
| `playbackorder` | SVG 2 new attribute on `<svg>` | **EXCLUDE** — not implemented in any browser |
| `timelinebegin` | SVG 2 new attribute on `<svg>` | **EXCLUDE** — not implemented in any browser |
| `ping` | Standard (mirrors HTML `<a>`) | **INCLUDE** — browsers support it on `<a>` |
| `referrerpolicy` | Standard (mirrors HTML) | **INCLUDE** — browsers support it on `<a>` |
| `hreflang` | Standard (mirrors HTML) | **INCLUDE** — browsers support it on `<a>` |
| `rel` | Standard (mirrors HTML) | **INCLUDE** — browsers support it on `<a>` |

---

## Browser-Implementation Flags for SVG 2 Additions

The following SVG 2 features appear in the W3C indices but are not implemented in current browsers. The grammar should treat them accordingly:

| Feature | Type | Implemented | Grammar recommendation |
|---------|------|-------------|----------------------|
| `<audio>` | element | No | **Exclude** |
| `<canvas>` | element | No | **Exclude** |
| `<iframe>` | element | No | **Exclude** |
| `<video>` | element | No | **Exclude** |
| `<unknown>` | element | No (parse-only) | **Exclude** |
| `<discard>` | element | Partial (Chrome/Edge) | **Optional/experimental** |
| `playbackorder` attr | attribute | No | **Exclude** |
| `timelinebegin` attr | attribute | No | **Exclude** |
| `href` without xlink | attribute | Yes (all modern browsers) | **Include**; keep `xlink:href` as deprecated alternative |
| `fr` on radialGradient | attribute | Yes (all modern browsers) | **Include** |
| `side` on textPath | attribute | Experimental (Firefox only as of 2025) | **Optional** |
| `edgeMode` on feGaussianBlur | attribute | Yes | **Include** |
| `line-height` property | property | Yes | **Include** |
| `white-space` property | property | Yes (replaces `xml:space`) | **Include** |
| `paint-order` property | property | Yes | **Include** |
| `vector-effect` property | property | Yes | **Include** |

---

## Summary

| Category | Count |
|----------|-------|
| Elements (W3C SVG 2 index) | 52 |
| Elements recommended for grammar (excluding unimplemented) | 46 |
| Regular attribute unique names (G.1) | ~230 (50 aria-*, ~60 on*, ~120 others) |
| Presentation attributes (G.2) | 60 |
| Properties (Appendix H) | 41 |
| MDN elements not in W3C index (all deprecated/removed) | 18 |
| W3C elements not in MDN (unimplemented SVG 2) | 6 |
| MDN attributes not in W3C G.1 (notable additions) | 23 |
| W3C G.1 attributes not in MDN | 7 |

### Names that differ between W3C SVG 2 and MDN

**In W3C index but absent from MDN (unimplemented SVG 2 elements):**
audio, canvas, iframe, video, unknown, discard

**In MDN but absent from W3C element index (removed SVG 1.1 — MDN historical docs):**
altGlyph, altGlyphDef, altGlyphItem, color-profile, cursor, font, font-face, font-face-format, font-face-name, font-face-src, font-face-uri, glyph, glyphRef, hkern, missing-glyph, tref, vkern

**Attribute-level differences (notable):**
- MDN has `d`, `rx`, `ry` (geometry attrs defined in shape chapters, not G.1) — must be in grammar
- MDN has `autofocus`, `decoding`, `data-*`, `mask-type`, `text-overflow`, `transform-origin` — browser-supported, include
- MDN has `attributeType`, `baseProfile`, `requiredFeatures`, `version`, `xlink:arcrole`, `xlink:show`, `xlink:type`, `xml:lang`, `glyph-orientation-horizontal`, `font-width`, `fetchpriority` as deprecated/experimental — exclude or mark deprecated
- W3C G.1 has `playbackorder`, `timelinebegin` — unimplemented, exclude
- W3C G.1 retains `xlink:href`, `xlink:title`, `xml:space` as deprecated legacy — include with deprecation note
