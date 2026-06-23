package main

import (
	"encoding/json"
	"fmt"
	"html"
	"sort"
	"strings"
)

// emit.go — write the generated galleries. Per element a dark-theme showcase
// page (one card per enumerated value-path: blueprint-wrapped SVG + a monospace
// attr="value" label), an index.html linking them, values.json
// (element→[{attr,value}]) and manifest.tsv.

// page bundles an element's enumerated variants for emission.
type page struct {
	tag      string
	variants []Variant
}

// galleryCSS is the showcase CSS from TEMPLATE_GUIDE (dark theme).
const galleryCSS = `body{margin:0;background:#1a1a2e;color:#e6e6e6;font:14px/1.4 ui-monospace,Menlo,monospace;padding:24px}
h1{color:#16c79a;font-size:18px;margin:0 0 4px}
p.desc{color:#9aa;margin:0 0 20px}
a{color:#4d8bff}
.grid{display:flex;flex-wrap:wrap;gap:16px}
.card{background:#0f1530;border:1px solid #26305a;border-radius:8px;padding:10px;width:160px}
.card svg{display:block;background:#161c3a;border-radius:4px;width:140px;height:140px}
.card .label{margin-top:8px;color:#f5a623;font-size:12px;word-break:break-all}
.card .attr{color:#16c79a}`

// emitPage renders one element's gallery HTML.
func emitPage(p page) string {
	var b strings.Builder
	t := html.EscapeString("<" + p.tag + ">")
	fmt.Fprintf(&b, `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>%s</title>
<style>%s</style></head><body>
<h1>%s</h1>
<p class="desc">%d enumerated value-paths walked from the grammar. <a href="index.html">&larr; index</a></p>
<div class="grid">
`, t, galleryCSS, t, len(p.variants))
	for _, v := range p.variants {
		label := fmt.Sprintf(`<span class="attr">%s</span>="%s"`,
			html.EscapeString(v.Attr), html.EscapeString(v.Value))
		fmt.Fprintf(&b, `  <div class="card">%s<div class="label">%s</div></div>
`, v.WrappedSVG, label)
	}
	b.WriteString("</div></body></html>\n")
	return b.String()
}

// emitIndex renders the index linking every element page.
func emitIndex(pages []page, totalVariants int) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>SVG grammar gallery</title>
<style>%s
ul{columns:4;list-style:none;padding:0}li{margin:2px 0}.n{color:#9aa}</style></head><body>
<h1>SVG grammar — all value-paths</h1>
<p class="desc">%d elements, %d enumerated value-path variants. Each page walks one element's grammar for every attribute and every value, rendered as an SVG.</p>
<ul>
`, galleryCSS, len(pages), totalVariants)
	for _, p := range pages {
		fmt.Fprintf(&b, `  <li><a href="%s.html">&lt;%s&gt;</a> <span class="n">%d</span></li>
`, p.tag, html.EscapeString(p.tag), len(p.variants))
	}
	b.WriteString("</ul></body></html>\n")
	return b.String()
}

// emitValuesJSON builds the element→[{attr,value}] map.
func emitValuesJSON(pages []page) string {
	type av struct {
		Attr  string `json:"attr"`
		Value string `json:"value"`
	}
	m := map[string][]av{}
	for _, p := range pages {
		for _, v := range p.variants {
			m[p.tag] = append(m[p.tag], av{Attr: v.Attr, Value: v.Value})
		}
	}
	data, _ := json.MarshalIndent(m, "", "  ")
	return string(data) + "\n"
}

// emitManifest builds the TSV: tag⇥attr⇥value⇥needsID.
func emitManifest(pages []page) string {
	var b strings.Builder
	b.WriteString("element\tattr\tvalue\tneeds_id\n")
	tags := make([]string, 0, len(pages))
	for _, p := range pages {
		tags = append(tags, p.tag)
	}
	sort.Strings(tags)
	byTag := map[string]page{}
	for _, p := range pages {
		byTag[p.tag] = p
	}
	for _, tag := range tags {
		for _, v := range byTag[tag].variants {
			needs := "0"
			if v.NeedsID {
				needs = "1"
			}
			fmt.Fprintf(&b, "%s\t%s\t%s\t%s\n",
				tag, v.Attr, strings.ReplaceAll(v.Value, "\t", " "), needs)
		}
	}
	return b.String()
}
