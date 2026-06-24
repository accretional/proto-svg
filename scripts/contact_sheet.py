#!/usr/bin/env python3
"""contact_sheet.py — tile each element's per-preset viewer screenshots (and the
first frame of each animation GIF) into one labeled grid image, for visual
validation. Pure stdlib + PIL; reads the already-captured screenshots, no browser.

Input:
  chrome-testing/gallery/catalogue.json          (elements[].presets[].name)
  chrome-testing/screenshots/gallery/<tag>/NN-<slug>.png   (static preset)
  chrome-testing/screenshots/gallery/<tag>/NN-<slug>.gif   (temporal preset)

Output:
  chrome-testing/screenshots/contact/<tag>.png   one labeled grid per element

Usage: contact_sheet.py [tag ...]   (default: all in catalogue.json)
"""
import json, os, re, sys
from PIL import Image, ImageDraw, ImageFont

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
CAT = os.path.join(ROOT, "chrome-testing/gallery/catalogue.json")
SHOTS = os.path.join(ROOT, "chrome-testing/screenshots/gallery")
OUT = os.path.join(ROOT, "chrome-testing/screenshots/contact")

CELL, LABEL_H, PAD, COLS = 150, 34, 8, 7
BG, CARD, FG, ACCENT = (10, 13, 12), (15, 19, 17), (205, 214, 210), (78, 227, 154)

slug_re = re.compile(r"[^a-z0-9]+")


def slugify(s):
    s = slug_re.sub("-", s.lower()).strip("-")
    if len(s) > 40:
        s = s[:40].strip("-")
    return s or "v"


def font(sz):
    for p in ("/System/Library/Fonts/Menlo.ttc",
              "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf"):
        if os.path.exists(p):
            try:
                return ImageFont.truetype(p, sz)
            except Exception:
                pass
    return ImageFont.load_default()


F, FT = font(11), font(15)


def load_thumb(path):
    try:
        im = Image.open(path)
        if getattr(im, "is_animated", False):
            im.seek(0)
        im = im.convert("RGB")
        im.thumbnail((CELL, CELL))
        return im
    except Exception:
        return None


def wrap(draw, text, fnt, maxw):
    out, line = [], ""
    for ch in text:
        if draw.textlength(line + ch, font=fnt) <= maxw:
            line += ch
        else:
            out.append(line); line = ch
        if len(out) >= 2:
            break
    if line and len(out) < 2:
        out.append(line)
    return out[:2]


def sheet(el):
    tag = el["tag"]
    d = os.path.join(SHOTS, tag)
    cells = []
    for i, p in enumerate(el.get("presets", [])):
        base = "%02d-%s" % (i, slugify(p["name"]))
        png, gif = os.path.join(d, base + ".png"), os.path.join(d, base + ".gif")
        temporal = os.path.exists(gif)
        path = gif if temporal else png
        cells.append((load_thumb(path), ("▶ " if temporal else "") + p["name"]))
    if not cells:
        return None
    n = len(cells)
    cols = min(COLS, n)
    rows = (n + cols - 1) // cols
    cw, ch = CELL + PAD * 2, CELL + LABEL_H + PAD * 2
    img = Image.new("RGB", (cols * cw, rows * ch + 40), BG)
    dr = ImageDraw.Draw(img)
    dr.text((PAD, 8), "<%s>  (%s) - %d presets" % (tag, el.get("name", ""), n), fill=ACCENT, font=FT)
    for i, (thumb, label) in enumerate(cells):
        r, c = divmod(i, cols)
        x, y = c * cw + PAD, r * ch + 40 + PAD
        dr.rectangle([x, y, x + CELL, y + CELL + LABEL_H], fill=CARD, outline=(38, 48, 42))
        if thumb:
            img.paste(thumb, (x + (CELL - thumb.width) // 2, y + (CELL - thumb.height) // 2))
        else:
            dr.text((x + 8, y + CELL // 2), "(missing)", fill=(200, 80, 80), font=F)
        ty = y + CELL + 3
        for ln in wrap(dr, label, F, CELL - 6):
            dr.text((x + 4, ty), ln, fill=FG, font=F); ty += 13
    return img


def main():
    cat = json.load(open(CAT))
    by_tag = {e["tag"]: e for e in cat["elements"]}
    os.makedirs(OUT, exist_ok=True)
    targets = sys.argv[1:] or [e["tag"] for e in cat["elements"]]
    made = 0
    for tag in targets:
        el = by_tag.get(tag)
        if not el:
            continue
        img = sheet(el)
        if img:
            img.save(os.path.join(OUT, tag + ".png")); made += 1
            print("  %s: %d presets" % (tag, len(el.get("presets", []))))
    print("wrote %d contact sheets to %s" % (made, OUT))


if __name__ == "__main__":
    main()
