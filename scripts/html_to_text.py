#!/usr/bin/env python3
"""html_to_text.py — Convert a W3C/HTML spec page into readable plain text.

Purpose: the SVG grammar is authored by reading the W3C spec prose, the formal
value-syntax blocks (BNF in <pre>), and the per-element/attribute descriptions.
Raw spec HTML is too noisy to read directly. This converter strips markup while
PRESERVING the two things that matter for grammar work:

  1. <pre> blocks verbatim (the path-data BNF, transform syntax, etc. live here),
     fenced so the agent can see exact whitespace/grammar.
  2. heading structure (h1..h6 become markdown #), so sections like
     "Value syntax", "Formal syntax", and per-element definitions stay locatable.

Dependency-free (stdlib html.parser only). Deterministic output.

Usage: html_to_text.py <input.html> <output.txt>
"""

import sys
from html.parser import HTMLParser
from html import unescape

# Tags whose text content is dropped entirely. NOTE: <head> is deliberately NOT
# here. HTML5/Bikeshed pages (filter-effects, css-masking) legally omit </head>,
# so skipping the <head> subtree would never close and would swallow the whole
# body. Its <style>/<script> children are skipped individually anyway; <title>
# and <meta> contribute no/negligible text.
SKIP_TAGS = {"script", "style", "noscript", "svg", "object"}

# Block-level tags: emit a newline boundary around them (outside <pre>).
BLOCK_TAGS = {
    "p", "div", "section", "article", "header", "footer", "main", "aside",
    "ul", "ol", "li", "dl", "dt", "dd", "table", "thead", "tbody", "tfoot",
    "tr", "th", "td", "blockquote", "figure", "figcaption", "hr", "nav",
    "h1", "h2", "h3", "h4", "h5", "h6", "pre", "form", "fieldset",
}

HEADINGS = {"h1": 1, "h2": 2, "h3": 3, "h4": 4, "h5": 5, "h6": 6}


class TextExtractor(HTMLParser):
    def __init__(self):
        super().__init__(convert_charrefs=True)
        self.out = []          # list of text fragments
        self.skip_depth = 0    # >0 while inside a SKIP_TAGS subtree
        self.pre_depth = 0     # >0 while inside <pre>
        self.heading = None    # current heading level, or None
        self.heading_buf = []  # text collected for the current heading

    # -- helpers --------------------------------------------------------------
    def _emit(self, s):
        if self.heading is not None:
            self.heading_buf.append(s)
        else:
            self.out.append(s)

    def _newline(self):
        if self.pre_depth == 0:
            self._emit("\n")

    # -- parser callbacks -----------------------------------------------------
    def handle_starttag(self, tag, attrs):
        if tag in SKIP_TAGS:
            self.skip_depth += 1
            return
        if self.skip_depth:
            return
        if tag == "br":
            self._emit("\n")
            return
        if tag == "pre":
            self.pre_depth += 1
            self.out.append("\n\n```\n")
            return
        if tag in HEADINGS and self.pre_depth == 0:
            self.out.append("\n\n" + "#" * HEADINGS[tag] + " ")
            self.heading = HEADINGS[tag]
            self.heading_buf = []
            return
        if tag in BLOCK_TAGS:
            self._newline()

    def handle_endtag(self, tag):
        if tag in SKIP_TAGS:
            if self.skip_depth:
                self.skip_depth -= 1
            return
        if self.skip_depth:
            return
        if tag == "pre":
            if self.pre_depth:
                self.pre_depth -= 1
                self.out.append("\n```\n\n")
            return
        if tag in HEADINGS and self.heading is not None:
            text = " ".join("".join(self.heading_buf).split())
            self.out.append(text)
            self.heading = None
            self.heading_buf = []
            self.out.append("\n")
            return
        if tag in BLOCK_TAGS:
            self._newline()

    def handle_data(self, data):
        if self.skip_depth:
            return
        if self.pre_depth:
            self.out.append(data)
        else:
            self._emit(data)

    def text(self):
        # Flush any heading left open by a missing close tag.
        if self.heading is not None and self.heading_buf:
            self.out.append(" ".join("".join(self.heading_buf).split()))
            self.heading = None
            self.heading_buf = []
        raw = "".join(self.out)
        # Collapse whitespace OUTSIDE fenced code blocks; keep code verbatim.
        lines = []
        in_fence = False
        for line in raw.split("\n"):
            if line.strip() == "```":
                in_fence = not in_fence
                lines.append("```")
                continue
            if in_fence:
                lines.append(line.rstrip())
            else:
                collapsed = " ".join(line.split())
                lines.append(collapsed)
        # Squeeze runs of 3+ blank lines down to a single blank line.
        result = []
        blanks = 0
        for line in lines:
            if line == "":
                blanks += 1
                if blanks <= 1:
                    result.append(line)
            else:
                blanks = 0
                result.append(line)
        return "\n".join(result).strip() + "\n"


def main():
    if len(sys.argv) != 3:
        print("usage: html_to_text.py <input.html> <output.txt>", file=sys.stderr)
        sys.exit(2)
    with open(sys.argv[1], "r", encoding="utf-8", errors="replace") as f:
        html = f.read()
    parser = TextExtractor()
    parser.feed(html)
    with open(sys.argv[2], "w", encoding="utf-8") as f:
        f.write(parser.text())


if __name__ == "__main__":
    main()
