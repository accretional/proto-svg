package main

// docs.go — HAND-AUTHORED long-form documentation for every SVG element, shown in
// the gallery's About panel. Each entry is written by hand (grounded in MDN, never
// script-generated) as a one-sentence summary, one or two prose paragraphs, and a
// few "good to know" tips. Backtick `spans` render as inline code in the panel.
//
// The short one-line taglines live in catmeta.go's descTable; this is the deep
// explanation of what the element IS and what it DOES.

var docTable = map[string]elementDoc{

	// ── shapes ──────────────────────────────────────────────────────────────
	"rect": {
		Summary: "The rectangle — SVG's box primitive, positioned by its top-left corner.",
		Body: []string{
			"`<rect>` draws an axis-aligned rectangle. `x`/`y` set the top-left corner and `width`/`height` the size, all in the current user coordinate system. Give `rx`/`ry` a value to round the corners; set only one and the other mirrors it.",
			"Like every shape it is painted with `fill` and `stroke`, and can be transformed, clipped, masked and filtered. A zero width or height makes it disappear entirely.",
		},
		Tips: []string{
			"`rx`/`ry` are clamped to half the width/height — a fully rounded short side becomes a stadium/pill.",
			"Percentages on `width`/`height` resolve against the viewport, not the parent shape.",
		},
	},
	"circle": {
		Summary: "A circle defined by a centre point and a single radius.",
		Body: []string{
			"`<circle>` is the simplest curved primitive: `cx`/`cy` place the centre and `r` sets the radius. There is no separate width and height — for an oval use `<ellipse>` instead.",
			"A circle with `r=\"0\"` is not rendered. The geometry is exact, so it scales crisply at any zoom.",
		},
		Tips: []string{
			"Stroke straddles the edge — half its width sits outside the radius, half inside.",
		},
	},
	"ellipse": {
		Summary: "An ellipse with independent horizontal and vertical radii.",
		Body: []string{
			"`<ellipse>` extends the circle with two radii: `rx` along x and `ry` along y, centred at `cx`/`cy`. Equal radii give a circle.",
			"In SVG 2 `rx` or `ry` may be `auto`, in which case the other radius is used for both — handy when only one axis is known.",
		},
	},

	// ── lines & paths ───────────────────────────────────────────────────────
	"line": {
		Summary: "A single straight segment between two points.",
		Body: []string{
			"`<line>` connects (`x1`,`y1`) to (`x2`,`y2`) with one straight stroke. It has no interior, so `fill` does nothing — only stroke properties (`stroke`, `stroke-width`, `stroke-dasharray`, `stroke-linecap`) affect how it looks.",
		},
		Tips: []string{
			"For more than two points use `<polyline>`; for curves use `<path>`.",
		},
	},
	"polyline": {
		Summary: "A connected run of straight segments through a list of points.",
		Body: []string{
			"`<polyline>` draws an open chain of line segments through the coordinate pairs in its `points` list. It is not automatically closed, so by default it reads as a multi-segment line rather than a shape.",
			"It can still take a `fill`, which paints the region enclosed by the implicit closing edge between the last and first points — useful for area charts.",
		},
	},
	"polygon": {
		Summary: "A closed shape built from a list of points.",
		Body: []string{
			"`<polygon>` is like `<polyline>` but automatically closes — an edge is drawn from the last point back to the first, enclosing a fillable area. Use it for triangles, stars, arrows and any straight-edged shape.",
		},
		Tips: []string{
			"`fill-rule` (`nonzero` vs `evenodd`) decides which regions count as inside for self-intersecting stars.",
		},
	},
	"path": {
		Summary: "The universal shape — any line or curve from a `d` command string.",
		Body: []string{
			"`<path>` is the most powerful primitive. Its `d` attribute is a mini-language of commands: `M` moveto, `L` lineto, `H`/`V` horizontal/vertical, `C`/`S` cubic Béziers, `Q`/`T` quadratic Béziers, `A` elliptical arcs, and `Z` to close. Uppercase commands are absolute, lowercase relative.",
			"Every other shape can be expressed as a path. It is what drawing tools export, what `<textPath>` follows, and what `pathLength` lets you normalise for dash and animation maths.",
		},
		Tips: []string{
			"`pathLength` rescales the path's reported length so `stroke-dasharray` and motion offsets use round numbers.",
			"Keep `d` precise but compact — relative commands often produce smaller, more editable strings.",
		},
	},

	// ── containers ──────────────────────────────────────────────────────────
	"svg": {
		Summary: "The root canvas — and, when nested, a viewport with its own coordinates.",
		Body: []string{
			"`<svg>` is the outermost element of any SVG document, but it can also nest inside another `<svg>` to establish a fresh coordinate system. `viewBox` defines the internal coordinates that map onto the element's `width`/`height`, and `preserveAspectRatio` controls how that mapping fits when the aspect ratios differ.",
			"A nested `<svg>` clips to its bounds by default, making it a clean way to embed a self-contained sub-drawing or reframe content.",
		},
		Tips: []string{
			"`viewBox=\"minX minY width height\"` is the single most important attribute for responsive, scalable SVG.",
			"`preserveAspectRatio=\"xMidYMid meet\"` fits and centres; `slice` fills and crops; `none` stretches.",
		},
	},
	"g": {
		Summary: "Groups elements so they share a transform, paint and other inherited properties.",
		Body: []string{
			"`<g>` has no geometry of its own. It bundles child elements so a single `transform`, `opacity`, `fill`, `clip-path` or `filter` applies to them all at once, and so they can be referenced or animated as a unit.",
			"Groups also organise a document logically, much like a folder, without affecting the rendered result beyond the inherited attributes you set on them.",
		},
		Tips: []string{
			"`opacity` on a `<g>` is group opacity — the whole group is composited then faded, so overlaps don't double up.",
		},
	},
	"defs": {
		Summary: "A container for reusable definitions that are not drawn until referenced.",
		Body: []string{
			"`<defs>` holds elements that exist to be referenced elsewhere — gradients, patterns, filters, clip paths, markers, symbols and reusable shapes. Nothing inside `<defs>` renders on its own; it becomes visible only when something points at it by `id`, e.g. `fill=\"url(#grad)\"` or `<use href=\"#icon\">`.",
			"Grouping definitions in `<defs>` keeps a document tidy and makes intent explicit, though technically any element with an `id` can be referenced from anywhere.",
		},
		Tips: []string{
			"In this lab the `<defs>` holds a shape with `id=\"slot\"` that a `<use>` instantiates — swap the definition and the rendered shape changes.",
		},
	},
	"use": {
		Summary: "Stamps a copy of an already-defined element by reference.",
		Body: []string{
			"`<use href=\"#id\">` clones the referenced element (or `<symbol>`) into the document at the `x`/`y` offset you give. The clone is a live, shadow copy: edit the original and every `<use>` of it updates, which is the backbone of icon systems and repeated motifs.",
			"Referencing a `<symbol>` also applies that symbol's own `width`/`height`/`viewBox`, so `<use>` becomes a sized, self-contained instance.",
		},
		Tips: []string{
			"Inherited paint (e.g. `fill`) only reaches the clone for properties the original left unset — set `fill=\"currentColor\"` on the original to recolour per instance.",
		},
	},
	"symbol": {
		Summary: "A reusable template that renders only when instantiated by `<use>`, with its own viewport.",
		Body: []string{
			"`<symbol>` defines graphics meant to be reused, like `<defs>` content, but it is never drawn directly — it appears only through `<use href=\"#id\">`. Unlike `<g>`, it carries its own `viewBox`, `preserveAspectRatio`, `refX` and `refY`, giving each instance an isolated, scalable coordinate system.",
			"This makes `<symbol>` the natural building block for icon libraries: define once, size and place freely with `<use>`.",
		},
		Tips: []string{
			"`refX`/`refY` set the symbol's anchor point — the coordinate that lands on the `<use>` element's `x`/`y`.",
		},
	},
	"switch": {
		Summary: "Renders only the first child whose conditional-processing tests pass.",
		Body: []string{
			"`<switch>` evaluates each direct child in order and draws the first one whose conditional attributes all evaluate true, skipping the rest. The tested attributes are `requiredExtensions`, `systemLanguage` (and the legacy `requiredFeatures`).",
			"It is most often used for localisation — offer the same label in several languages and let `systemLanguage` pick the right one — or to provide a fallback when an extension is unsupported.",
		},
		Tips: []string{
			"Order matters: put the most specific match first and a no-condition fallback last.",
		},
	},
	"a": {
		Summary: "A hyperlink that wraps SVG content.",
		Body: []string{
			"`<a>` is SVG's anchor: wrap any shape, group or text in it and that content becomes a clickable link via `href`. It behaves like an HTML `<a>` — `target`, `download`, `rel` and friends all apply.",
			"Because it is a container, the link region is exactly the bounding geometry of whatever it contains.",
		},
		Tips: []string{
			"Add a `<title>` child for an accessible link name and hover tooltip.",
		},
	},
	"foreignObject": {
		Summary: "Embeds HTML (or any other XML namespace) inside an SVG drawing.",
		Body: []string{
			"`<foreignObject>` carves out a rectangular region (`x`,`y`,`width`,`height`) inside the SVG and hands it to another renderer — almost always HTML. It is the bridge for putting wrapping paragraphs, form controls or live web content into a vector scene.",
			"Content must declare its namespace, e.g. a `<div xmlns=\"http://www.w3.org/1999/xhtml\">`. The region participates in SVG transforms and clipping like any other element.",
		},
		Tips: []string{
			"Support in standalone SVG-as-image contexts is limited — it shines when the SVG is inline in an HTML page.",
		},
	},
	"view": {
		Summary: "A named viewport you can jump to by URL fragment.",
		Body: []string{
			"`<view id=\"name\">` predefines a `viewBox` (and optional `preserveAspectRatio`) that can be activated by referencing the SVG with that fragment, e.g. `picture.svg#name`. Loading that URL reframes the whole document to the view's box — like a saved camera position.",
			"It renders nothing inline; its only effect is the framing it applies when targeted.",
		},
		Tips: []string{
			"Great for slide-style storytelling or deep-linking to a detail within one SVG file.",
		},
	},

	// ── text ────────────────────────────────────────────────────────────────
	"text": {
		Summary: "Real, selectable, accessible type rendered as vector glyphs.",
		Body: []string{
			"`<text>` places a run of characters at (`x`,`y`), where the y is the baseline. It honours the full set of font and text properties (`font-family`, `font-size`, `text-anchor`, `letter-spacing`) and is painted with `fill`/`stroke` like any shape.",
			"`x`/`y`/`dx`/`dy`/`rotate` accept lists, letting you position or rotate individual glyphs, and `textLength` with `lengthAdjust` fits the run to an exact width.",
		},
		Tips: []string{
			"`text-anchor` (`start`/`middle`/`end`) aligns the run horizontally around the `x` position.",
			"Real text stays selectable and searchable — prefer it over outlining type to paths.",
		},
	},
	"tspan": {
		Summary: "A styled or repositioned sub-run inside a `<text>`.",
		Body: []string{
			"`<tspan>` carves a piece of text out of its surrounding `<text>` so it can be styled or moved independently — a coloured word, a superscript, a second line. Its `x`/`y` reposition absolutely while `dx`/`dy` nudge relative to the flow.",
			"Nesting `<tspan>`s lets you build rich inline formatting from a single text element.",
		},
		Tips: []string{
			"Set `dy` on a `<tspan>` and reset `x` to fake line breaks (SVG `<text>` does not wrap on its own).",
		},
	},
	"textPath": {
		Summary: "Text flowed along the shape of a referenced path.",
		Body: []string{
			"`<textPath href=\"#path\">` lays its characters along an arbitrary `<path>` (or basic shape) instead of a straight baseline, so type can curve, loop or spiral. `startOffset` shifts where the text begins along the path.",
			"`side` chooses which side of the path the text sits on, `method` picks between aligning and stretching glyphs, and `spacing` controls inter-glyph spacing.",
		},
		Tips: []string{
			"A negative or percentage `startOffset` slides the text — easy to animate for a marquee-on-a-curve effect.",
		},
	},

	// ── paint & fills ─────────────────────────────────────────────────────────
	"linearGradient": {
		Summary: "A smooth colour blend along a straight vector, used as a paint.",
		Body: []string{
			"`<linearGradient>` interpolates colour between `<stop>` children along the line from (`x1`,`y1`) to (`x2`,`y2`). Reference it with `fill=\"url(#id)\"` or `stroke=\"url(#id)\"`.",
			"`gradientUnits` decides whether those coordinates are fractions of the shape's bounding box (`objectBoundingBox`, the default) or absolute user units (`userSpaceOnUse`); `spreadMethod` controls what happens beyond the end stops.",
		},
		Tips: []string{
			"`gradientTransform` rotates or skews the gradient independently of the shape it paints.",
		},
	},
	"radialGradient": {
		Summary: "A colour blend radiating outward from a focal point.",
		Body: []string{
			"`<radialGradient>` blends its `<stop>`s from a centre (`cx`,`cy`) out to radius `r`, optionally with a separate focal point (`fx`,`fy`) for an off-centre highlight that reads as a light source or sphere.",
			"Like its linear sibling it supports `gradientUnits`, `gradientTransform` and `spreadMethod`, and is referenced through `url(#id)`.",
		},
		Tips: []string{
			"Move `fx`/`fy` away from `cx`/`cy` to suggest a glossy, directionally-lit surface.",
		},
	},
	"stop": {
		Summary: "One colour stop within a gradient.",
		Body: []string{
			"`<stop>` defines a single point in a gradient: `offset` (0–1, or a percentage) is its position along the gradient and `stop-color` its colour, with `stop-opacity` for partial transparency.",
			"A gradient needs at least two stops; their order and offsets shape the entire blend.",
		},
		Tips: []string{
			"Two stops at the same offset create a hard colour band instead of a smooth transition.",
		},
	},
	"pattern": {
		Summary: "A tile that repeats across any fill or stroke.",
		Body: []string{
			"`<pattern>` defines a small piece of artwork that is tiled to paint a shape via `fill=\"url(#id)\"`. `width`/`height` set the tile size and `x`/`y` its origin. `patternUnits` controls how those are measured and `patternContentUnits` how the tile's own contents are measured.",
			"A `viewBox` lets the tile scale independently, and `patternTransform` rotates or scales the whole tiling — perfect for hatching, polka dots, textures and repeating motifs.",
		},
		Tips: []string{
			"Default `patternUnits=\"objectBoundingBox\"` makes the tile size a fraction of the shape — switch to `userSpaceOnUse` for fixed-size tiles.",
		},
	},

	// ── clipping & masking ──────────────────────────────────────────────────
	"clipPath": {
		Summary: "Clips an element to an arbitrary shape with hard, 1-bit edges.",
		Body: []string{
			"`<clipPath>` defines a region — built from shapes, paths or text — and everything outside it is cut away when applied via `clip-path=\"url(#id)\"`. Clipping is binary: a pixel is either fully kept or fully removed, giving crisp silhouette edges.",
			"`clipPathUnits` chooses between bounding-box fractions and user-space coordinates for the clip geometry.",
		},
		Tips: []string{
			"For soft, graduated edges or partial transparency use `<mask>` instead — clipping cannot feather.",
		},
	},
	"mask": {
		Summary: "Masks an element by luminance or alpha — soft and graduated.",
		Body: []string{
			"`<mask>` composites its content against an element via `mask=\"url(#id)\"`. By default it is a luminance mask: bright mask pixels reveal, dark ones hide, and greys produce partial transparency — so gradients in the mask create smooth fades.",
			"Set `mask-type=\"alpha\"` (or the CSS property) to drive visibility from the mask's alpha channel instead. `maskUnits` and `maskContentUnits` control the coordinate systems.",
		},
		Tips: []string{
			"A black→white gradient mask is the classic way to fade an element out along an edge.",
		},
	},
	"marker": {
		Summary: "A graphic stamped onto the vertices of a path or shape.",
		Body: []string{
			"`<marker>` defines a small symbol — an arrowhead, dot or tick — that is drawn at a shape's points via the `marker-start`, `marker-mid` and `marker-end` properties. `refX`/`refY` set the anchor that lands on each vertex and `markerWidth`/`markerHeight` the marker viewport.",
			"`orient=\"auto\"` rotates the marker to follow the path direction (so arrowheads point the right way), and `markerUnits` scales it either with the stroke width or in fixed user units.",
		},
		Tips: []string{
			"`orient=\"auto-start-reverse\"` flips the start marker so a single arrow definition works at both ends.",
		},
	},

	// ── filters ───────────────────────────────────────────────────────────────
	"filter": {
		Summary: "A container for a pipeline of filter primitives applied as a post-process.",
		Body: []string{
			"`<filter>` holds an ordered chain of `fe*` primitives. Applied via `filter=\"url(#id)\"`, it takes the element's rendering (`SourceGraphic`/`SourceAlpha`) and runs it through the chain to produce effects like blur, shadow, lighting and texture.",
			"`x`/`y`/`width`/`height` define the filter region — the canvas the effect can paint into. Primitives pass results to one another via `result` names and the `in`/`in2` inputs.",
		},
		Tips: []string{
			"Effects often spill outside the element — widen the region (e.g. `x=\"-20%\" width=\"140%\"`) so blurs and shadows aren't clipped.",
		},
	},
	"feBlend": {
		Summary: "Blends two filter inputs with a Photoshop-style blend mode.",
		Body: []string{
			"`<feBlend>` combines its `in` and `in2` inputs using a `mode` such as `normal`, `multiply`, `screen`, `darken`, `lighten`, `overlay` or `color-dodge`. The modes match those in image editors, mixing the two layers' colours rather than simply stacking them.",
		},
		Tips: []string{
			"Blend modes only differ visibly when the two layers overlap and contrast — feed it distinct, overlapping inputs.",
		},
	},
	"feColorMatrix": {
		Summary: "Transforms colours through a matrix — saturate, hue-rotate or fully recolour.",
		Body: []string{
			"`<feColorMatrix>` multiplies each pixel's [R,G,B,A] by a transformation. `type=\"matrix\"` takes a full 20-value 5×4 matrix for arbitrary channel mixing; `saturate` takes one number (0 = greyscale, 1 = normal); `hueRotate` takes an angle in degrees; and `luminanceToAlpha` turns brightness into opacity.",
			"It is the go-to primitive for tinting, desaturating, channel swaps and converting a graphic into an alpha matte.",
		},
		Tips: []string{
			"`luminanceToAlpha` is a quick way to build a mask from any artwork's brightness.",
		},
	},
	"feComponentTransfer": {
		Summary: "Remaps each colour channel through its own transfer function.",
		Body: []string{
			"`<feComponentTransfer>` adjusts colour curves per channel using up to four children — `<feFuncR>`, `<feFuncG>`, `<feFuncB>`, `<feFuncA>`. Each function reshapes its channel: `identity` leaves it, `linear` applies slope+intercept, `gamma` applies amplitude·Cᵉˣᵖ+offset, and `table`/`discrete` remap through a lookup list.",
			"Together they cover brightness, contrast, colour balance, posterisation and thresholding.",
		},
		Tips: []string{
			"`discrete` with a few `tableValues` posterises a channel into hard steps.",
		},
	},
	"feComposite": {
		Summary: "Combines two inputs with Porter-Duff or arithmetic compositing.",
		Body: []string{
			"`<feComposite>` merges `in` and `in2` using an `operator`: the Porter-Duff set `over`, `in`, `out`, `atop` and `xor` controls how the two layers' coverage combines, while `lighter` adds them.",
			"`operator=\"arithmetic\"` instead computes `result = k1·i1·i2 + k2·i1 + k3·i2 + k4` per channel, a flexible formula for masking, fades and custom mixes via the `k1`–`k4` coefficients.",
		},
		Tips: []string{
			"`operator=\"in\"` is a common way to clip one filter result to the shape of another.",
		},
	},
	"feConvolveMatrix": {
		Summary: "Convolves pixels with a kernel — blur, sharpen, emboss or detect edges.",
		Body: []string{
			"`<feConvolveMatrix>` recomputes each pixel as a weighted blend of its neighbours. `order` sets the kernel size and `kernelMatrix` its weights; `divisor` normalises the result and `bias` offsets it. The chosen weights determine the effect — a centre-heavy kernel sharpens, a ring detects edges, an asymmetric one embosses.",
			"`edgeMode` decides how pixels at the boundary are sampled (`duplicate`, `wrap` or `none`).",
		},
		Tips: []string{
			"It needs textured input to show anything — a flat fill convolves to the same flat colour.",
		},
	},
	"feDiffuseLighting": {
		Summary: "Lights an alpha bump-map with matte (Lambertian) shading.",
		Body: []string{
			"`<feDiffuseLighting>` treats its input's alpha channel as a height-field and shades it as a non-shiny surface lit by a child light source (`<feDistantLight>`, `<fePointLight>` or `<feSpotLight>`). `surfaceScale` exaggerates the relief, `diffuseConstant` scales brightness and `lighting-color` tints the light.",
			"The output is opaque RGBA, usually composited back over the original artwork to add dimensional shading.",
		},
		Tips: []string{
			"Blur the alpha source first (`feGaussianBlur` on `SourceAlpha`) so there's a smooth gradient for the light to grade across.",
		},
	},
	"feDisplacementMap": {
		Summary: "Warps one input using the pixel values of another.",
		Body: []string{
			"`<feDisplacementMap>` shifts each pixel of `in` by an amount read from `in2`: `P'(x,y) = P(x + scale·(Xc−0.5), y + scale·(Yc−0.5))`. `xChannelSelector` and `yChannelSelector` pick which channels of the map drive horizontal and vertical displacement, and `scale` sets the strength.",
			"Paired with `feTurbulence` as the map it produces ripples, liquid distortion and organic wobble.",
		},
		Tips: []string{
			"Animate `scale` from 0 for a melt/ripple-in effect.",
		},
	},
	"feDistantLight": {
		Summary: "An infinitely-distant directional light for a lighting filter.",
		Body: []string{
			"`<feDistantLight>` is a light source used inside `<feDiffuseLighting>` or `<feSpecularLighting>`. Because it is infinitely far away, only its direction matters: `azimuth` is the compass angle in the xy-plane and `elevation` the angle above the surface, like sunlight.",
		},
	},
	"feDropShadow": {
		Summary: "A one-shot offset, blurred, coloured drop shadow.",
		Body: []string{
			"`<feDropShadow>` is a convenience primitive bundling offset, blur and composite into one step: it blurs the source alpha by `stdDeviation`, offsets it by `dx`/`dy`, tints it with `flood-color`/`flood-opacity`, and draws the original on top.",
			"It replaces the classic `feGaussianBlur`+`feOffset`+`feMerge` recipe with a single, readable element.",
		},
		Tips: []string{
			"Give the `<filter>` room (negative `x`/`y`, >100% width/height) or the shadow clips at the edge.",
		},
	},
	"feFlood": {
		Summary: "Fills the entire filter region with a solid colour.",
		Body: []string{
			"`<feFlood>` paints a flat rectangle of `flood-color` at `flood-opacity` across the filter region (or a sub-region you set). On its own it is just a colour block; combined with `feComposite`/`feMerge` it becomes the colour layer behind shadows, tints and glows.",
		},
	},
	"feFuncR": {
		Summary: "The red-channel transfer function inside feComponentTransfer.",
		Body: []string{
			"`<feFuncR>` reshapes the red channel within a `<feComponentTransfer>`. Its `type` chooses the curve — `identity`, `linear` (slope/intercept), `gamma` (amplitude/exponent/offset) or `table`/`discrete` (a `tableValues` lookup) — and the result is recombined with the other channels.",
		},
	},
	"feFuncG": {
		Summary: "The green-channel transfer function inside feComponentTransfer.",
		Body: []string{
			"`<feFuncG>` is the green-channel counterpart of `<feFuncR>`, applying its own transfer curve (`identity`, `linear`, `gamma`, `table` or `discrete`) to remap green within a `<feComponentTransfer>`.",
		},
	},
	"feFuncB": {
		Summary: "The blue-channel transfer function inside feComponentTransfer.",
		Body: []string{
			"`<feFuncB>` remaps the blue channel within a `<feComponentTransfer>` using its `type` curve — the same options as the other channels — letting you shift colour balance or posterise blue independently.",
		},
	},
	"feFuncA": {
		Summary: "The alpha-channel transfer function inside feComponentTransfer.",
		Body: []string{
			"`<feFuncA>` reshapes the alpha (opacity) channel within a `<feComponentTransfer>`. A `linear` or `gamma` curve fades transparency smoothly; a `discrete` table thresholds it into hard cut-outs.",
		},
	},
	"feGaussianBlur": {
		Summary: "Softens an input with a Gaussian blur of a given standard deviation.",
		Body: []string{
			"`<feGaussianBlur>` is the workhorse softening primitive. `stdDeviation` sets the blur amount, and a two-number value blurs x and y by different amounts for directional softness.",
			"Beyond plain blurring it is the first step of most shadows, glows and lighting bump-maps (blurring `SourceAlpha` into a smooth height-field).",
		},
		Tips: []string{
			"`stdDeviation=\"0\"` is a no-op; larger values spread further and need a roomier filter region.",
		},
	},
	"feImage": {
		Summary: "Loads an external image or document element as a filter input.",
		Body: []string{
			"`<feImage>` brings outside pixels into the filter chain — referenced by `href` (a raster, an SVG, or even another element by fragment) — so they can be blended, displaced or composited with the source. `preserveAspectRatio` controls how the image fits the region.",
		},
		Tips: []string{
			"Use it to supply a custom texture as the displacement map for `feDisplacementMap`.",
		},
	},
	"feMerge": {
		Summary: "Stacks several filter results into one layered output.",
		Body: []string{
			"`<feMerge>` composites multiple earlier results on top of one another in order, each named by an `<feMergeNode in=\"...\">` child. It is how the shadow recipe puts the blurred, offset shadow underneath the crisp original in a single step.",
		},
	},
	"feMergeNode": {
		Summary: "One input layer inside an feMerge.",
		Body: []string{
			"`<feMergeNode>` names a single layer to stack inside its parent `<feMerge>`. Its only attribute, `in`, points at a previous primitive's `result` (or a source keyword); the order of the nodes is the painting order, bottom to top.",
		},
	},
	"feMorphology": {
		Summary: "Fattens (dilate) or thins (erode) an input.",
		Body: []string{
			"`<feMorphology>` grows or shrinks the opaque parts of its input. `operator=\"dilate\"` spreads them outward (thicker strokes, bolder text); `operator=\"erode\"` pulls them inward (thinner, eaten away). `radius` sets how far, and a two-number radius works on x and y separately.",
		},
		Tips: []string{
			"Dilate then subtract the original to build an outline/stroke around arbitrary artwork.",
		},
	},
	"feOffset": {
		Summary: "Shifts an input by dx/dy — the foundation of drop shadows.",
		Body: []string{
			"`<feOffset>` simply translates its input by `dx`/`dy` within the filter region. Alone it just moves pixels; offsetting a blurred `SourceAlpha` and merging the original on top is the classic drop-shadow construction.",
		},
	},
	"fePointLight": {
		Summary: "A point light source at a 3-D position for a lighting filter.",
		Body: []string{
			"`<fePointLight>` is an omnidirectional light placed at (`x`,`y`,`z`) inside `<feDiffuseLighting>` or `<feSpecularLighting>`. Like a bare bulb it radiates in all directions, so its position relative to the surface — especially the `z` height — controls where highlights fall and how they spread.",
		},
	},
	"feSpecularLighting": {
		Summary: "Lights an alpha bump-map with glossy specular highlights.",
		Body: []string{
			"`<feSpecularLighting>` is the shiny counterpart to `feDiffuseLighting`: it adds bright, view-dependent highlights to an alpha height-field. `specularConstant` scales highlight intensity and `specularExponent` their tightness (higher = sharper, glassier), while `surfaceScale` and `lighting-color` behave as in diffuse lighting.",
			"Its result is typically composited additively over the artwork to give a wet or metallic sheen.",
		},
	},
	"feSpotLight": {
		Summary: "A cone-shaped spotlight for a lighting filter.",
		Body: []string{
			"`<feSpotLight>` is positioned at (`x`,`y`,`z`) and aimed at (`pointsAtX`,`pointsAtY`,`pointsAtZ`), casting a directional cone of light. `specularExponent` concentrates the beam toward its axis and `limitingConeAngle` clips it to a hard-edged circle of light.",
		},
	},
	"feTile": {
		Summary: "Tiles a filter result repeatedly across the filter region.",
		Body: []string{
			"`<feTile>` takes a small input result (defined by its sub-region) and repeats it edge-to-edge to fill the area. Combined with `feFlood` or `feImage` it builds repeating textures and step-and-repeat patterns entirely within the filter chain.",
		},
		Tips: []string{
			"The tiled input must itself be non-uniform — tiling a flat colour produces the same flat colour.",
		},
	},
	"feTurbulence": {
		Summary: "Generates Perlin noise — clouds, marble and organic texture.",
		Body: []string{
			"`<feTurbulence>` synthesises procedural noise to fill the filter region. `type=\"fractalNoise\"` gives soft, cloud-like fields; `type=\"turbulence\"` gives sharper, marbled patterns. `baseFrequency` sets the feature scale (small = large blobs), `numOctaves` adds fine detail, and `seed` makes it reproducible.",
			"It is the raw material for textures and, fed into `feDisplacementMap`, for organic distortion.",
		},
		Tips: []string{
			"`stitchTiles=\"stitch\"` makes the noise seamlessly tileable across the region edges.",
		},
	},

	// ── content & motion ──────────────────────────────────────────────────────
	"animate": {
		Summary: "Declarative SMIL animation — tweens a single attribute over time.",
		Body: []string{
			"`<animate>` smoothly changes one attribute (named by `attributeName`) of its parent element from `from` to `to` (or across a `values` list) over `dur`. `begin`, `repeatCount`, `fill` and `calcMode` control timing, looping and easing.",
			"It runs without any script, defined entirely in markup, making simple motion self-contained inside the SVG.",
		},
		Tips: []string{
			"`fill=\"freeze\"` holds the final value after the animation ends instead of snapping back.",
			"Use `keyTimes`+`keySplines` with `calcMode=\"spline\"` for custom easing curves.",
		},
	},
	"animateMotion": {
		Summary: "Moves an element along a motion path.",
		Body: []string{
			"`<animateMotion>` translates its parent along a trajectory — given inline as a `path`, or by referencing a real `<path>` through a child `<mpath>`. `rotate=\"auto\"` turns the element to face its direction of travel, so it banks through curves.",
			"`keyPoints` with `keyTimes` lets you vary speed along the path rather than moving at constant rate.",
		},
		Tips: []string{
			"Pair it with `<mpath href=\"#track\">` to reuse an existing path as the route.",
		},
	},
	"animateTransform": {
		Summary: "Animates a transform — translate, scale, rotate or skew.",
		Body: []string{
			"`<animateTransform>` animates the `transform` attribute specifically, with `type` selecting `translate`, `scale`, `rotate`, `skewX` or `skewY`. Because it targets the transform list it can spin, grow or slide an element where a plain `<animate>` cannot.",
		},
		Tips: []string{
			"For `rotate`, the `from`/`to` values can include a centre: `0 50 50` → `360 50 50` spins about (50,50).",
		},
	},
	"set": {
		Summary: "Sets an attribute to a value for a span of time — a discrete, non-interpolated change.",
		Body: []string{
			"`<set>` snaps `attributeName` to a single `to` value at `begin` and (optionally) reverts after `dur`. Unlike `<animate>` it does not tween — it is an instantaneous state change, ideal for toggling visibility, swapping a colour on interaction, or scripting show/hide.",
		},
		Tips: []string{
			"`begin=\"click\"` makes `<set>` a pure-markup interaction handler.",
		},
	},
	"mpath": {
		Summary: "References an existing path for animateMotion to follow.",
		Body: []string{
			"`<mpath href=\"#path\">` is a child of `<animateMotion>` that points it at a real `<path>` element to use as the motion route, instead of duplicating the geometry in a `path` attribute. Edit the referenced path and the motion updates with it.",
		},
	},
	"discard": {
		Summary: "Removes its target element from the document at a set time.",
		Body: []string{
			"`<discard>` deletes its target (the parent, or an `href` element) once its `begin` time is reached, permanently removing it from the DOM. It is a performance hint for long or looping animations — drop elements that will never be seen again to free resources.",
		},
		Tips: []string{
			"The removal is one-way; the element does not come back.",
		},
	},
	"image": {
		Summary: "Embeds a raster or SVG image into the drawing.",
		Body: []string{
			"`<image href=\"...\">` places an external bitmap (PNG/JPEG/GIF) or SVG into a rectangle defined by `x`/`y`/`width`/`height`. `preserveAspectRatio` controls how the image fits that box, and the image participates in SVG transforms, clipping and filters like any element.",
		},
		Tips: []string{
			"A `data:` URI embeds the image inline so the SVG stays self-contained with no extra request.",
		},
	},

	// ── descriptive / metadata / content ────────────────────────────────────
	"desc": {
		Summary: "A long-form, accessible description of its parent element.",
		Body: []string{
			"`<desc>` provides extended descriptive text for the element that contains it, read by assistive technology. It is never rendered visually — its job is semantic, conveying detail that the artwork alone cannot.",
			"Place it as the first child of any container or shape you want to describe; pair it with `<title>` for a short name plus a fuller explanation.",
		},
		Tips: []string{
			"Reference it from the parent with `aria-describedby` for the most reliable screen-reader support.",
		},
	},
	"title": {
		Summary: "An accessible name and hover tooltip for its parent element.",
		Body: []string{
			"`<title>` gives its parent a short human-readable name. Most browsers surface it as a native tooltip on hover, and assistive technology announces it as the element's accessible name — the SVG equivalent of an HTML `title`/`alt`.",
			"It renders nothing itself; being the first child of the element it names is what wires it up.",
		},
		Tips: []string{
			"Add a `<title>` to interactive elements (`<a>`, buttons) and to the root `<svg>` for a document name.",
		},
	},
	"metadata": {
		Summary: "Machine-readable metadata about the document — authorship, licensing, RDF.",
		Body: []string{
			"`<metadata>` carries data meant for machines rather than viewers: typically RDF, Dublin Core or other XML vocabularies describing the document's author, license, creation tool or keywords. It is invisible and ignored by the renderer, but available to crawlers and asset pipelines.",
			"Editors such as Inkscape write provenance here, and licensing tools read it back.",
		},
		Tips: []string{
			"Content is usually namespaced XML (e.g. an `<rdf:RDF>` block), not plain text.",
		},
	},
	"script": {
		Summary: "Embedded or referenced ECMAScript that runs in the SVG.",
		Body: []string{
			"`<script>` adds scripting to an SVG, either inline as a text body or via `href` to an external file. In a live browser it can manipulate the DOM, respond to events and drive interaction, exactly like a script in HTML.",
			"Note that a `<script>` inserted through `innerHTML` (as this lab renders previews) does not execute per the HTML spec — its code runs when the SVG is loaded as a real document or file.",
		},
		Tips: []string{
			"`type` defaults to JavaScript; a non-JS MIME type makes the browser treat the body as data, not code.",
		},
	},
	"style": {
		Summary: "Embedded CSS that styles elements in the document.",
		Body: []string{
			"`<style>` holds CSS rules that target SVG elements by tag, class or id — `fill`, `stroke`, `opacity`, transforms, even `@keyframes` animations. It centralises presentation so many elements share one rule instead of repeated presentation attributes.",
			"CSS in a `<style>` wins over presentation attributes (which act like the lowest-priority defaults), so it is the clean way to theme and animate a drawing.",
		},
		Tips: []string{
			"`@keyframes` plus an `animation` property gives script-free, CSS-driven motion.",
			"Class and tag selectors work just as in HTML — `.icon { fill: currentColor }` is a common pattern.",
		},
	},
}
