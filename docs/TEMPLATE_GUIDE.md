# SVG Template Authoring + Verification Guide

Hand-authored, one HTML file per SVG element at `chrome-testing/html/template/<tag>.html`
(use the real tag: `rect.html`, `linearGradient.html`, `feGaussianBlur.html`).

Each template has TWO parts:

1. A **human showcase** (the visible page): the element demonstrated with multiple
   variations, each visibly distinguishable, so a person can see what every
   attribute/value does.
2. A **machine blueprint** (a hidden block): the minimal SVG scaffold that supplies
   only the element's *neighbors*, with a slot the generator fills with
   grammar-generated element markup.

These are hand-authored. Do NOT write a script that emits templates; write each by hand.

## Showcase structure

```html
<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"><title>&lt;rect&gt;</title>
<style>
  body{margin:0;background:#1a1a2e;color:#e6e6e6;font:14px/1.4 ui-monospace,Menlo,monospace;padding:24px}
  h1{color:#16c79a;font-size:18px;margin:0 0 4px}
  p.desc{color:#9aa;margin:0 0 20px}
  .grid{display:flex;flex-wrap:wrap;gap:16px}
  .card{background:#0f1530;border:1px solid #26305a;border-radius:8px;padding:10px;width:160px}
  .card svg{display:block;background:#161c3a;border-radius:4px;width:140px;height:140px}
  .card .label{margin-top:8px;color:#f5a623;font-size:12px;word-break:break-all}
</style></head>
<body>
  <h1>&lt;rect&gt;</h1>
  <p class="desc">Rectangle. Showcasing x/y/width/height, rx/ry rounding, fill, stroke.</p>
  <div class="grid">
    <div class="card"><svg viewBox="0 0 100 100"><rect x="10" y="10" width="80" height="80" fill="#e94560"/></svg><div class="label">basic</div></div>
    <div class="card"><svg viewBox="0 0 100 100"><rect x="10" y="10" width="80" height="80" rx="20" fill="#16c79a"/></svg><div class="label">rx="20"</div></div>
    <!-- one card per attribute/value worth showing; each VISUALLY distinct -->
  </div>

  <!-- machine blueprint: the generator fills {{ELEMENT}} with grammar variants -->
  <script type="application/svg-blueprint" id="blueprint">
<svg xmlns="http://www.w3.org/2000/svg" width="120" height="120" viewBox="0 0 100 100">{{ELEMENT}}</svg>
  </script>
</body></html>
```

Rules for the showcase:
- Cover every attribute/value that has a *visible* effect, one card each, with a label
  saying what is varied. Each card must look DIFFERENT from its neighbours (use
  contrasting fills/sizes/strokes) so distinguishability is obvious at a glance.
- Keep each SVG small (≈100×100 viewBox) and self-contained.
- Bright palette on the dark cards: `#e94560 #16c79a #f5a623 #4d8bff #b06 aff #ffd166`.

## Machine blueprint (the `{{ELEMENT}}` slot)

The blueprint is ONE complete `<svg>` providing only the neighbours the element needs;
`{{ELEMENT}}` marks where the generator injects the element. The generator gives the
injected element `id="slot"` when an id is needed; reference it with `url(#slot)` /
`href="#slot"`.

- Self-rendering elements (shapes, text, image, g): `<svg ...>{{ELEMENT}}</svg>`.
- Referenced paint servers (linearGradient/radialGradient/pattern): put `{{ELEMENT}}`
  in `<defs>` and add a shape that paints with it:
  `<svg ...><defs>{{ELEMENT}}</defs><rect x="5" y="5" width="90" height="90" fill="url(#slot)"/></svg>`
- `stop`: nest inside a gradient that a rect uses:
  `<svg ...><defs><linearGradient id="slot"><stop offset="0" stop-color="#e94560"/>{{ELEMENT}}</linearGradient></defs><rect width="100" height="100" fill="url(#slot)"/></svg>`
- marker: `<defs><marker id="slot">…</marker></defs><path d="M10 50 H90" stroke="#16c79a" stroke-width="3" marker-end="url(#slot)"/>` with `{{ELEMENT}}` as the marker.
- clipPath / mask: `{{ELEMENT}}` in `<defs>`, plus a rect with `clip-path="url(#slot)"` / `mask="url(#slot)"`.
- filter + primitives: `<defs><filter id="slot">{{ELEMENT}}</filter></defs><rect … filter="url(#slot)"/>` (for a primitive, the blueprint already contains the `<filter>` wrapper).
- animation elements: `{{ELEMENT}}` as a child of an animated shape:
  `<rect width="40" height="40" fill="#4d8bff">{{ELEMENT}}</rect>`.
- structural/container (svg/g/defs/symbol/use/switch/a): a scaffold that makes the
  container visible (e.g. `use` references a defined shape).

## Manual verification protocol (REQUIRED, per template, no diff scripts)

For each template you author:

1. Screenshot it:
   `./chrome-testing/snap.sh chrome-testing/html/template/<tag>.html chrome-testing/screenshots/template/<tag>.png`
   (chromerpc is already built+cached; this just captures the PNG.)
2. **Read the PNG** and look at it. Confirm with your own eyes:
   - it actually renders (not blank, not a broken/partial SVG),
   - every variation card is visibly DISTINCT from the others,
   - the element and the property being varied are clearly demonstrated.
3. If anything is blank, broken, or two cards look identical, FIX the template and
   re-screenshot. Do not consider a template done until its screenshot visually
   confirms distinguishable, correct rendering.

Report, per element: the variations shown and a one-line confirmation that you viewed
the screenshot and it renders distinctly.
