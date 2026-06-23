#!/usr/bin/env python3
"""contact_sheet.py — tile each element's per-value specimen PNGs (and the first
frame of each animation GIF) into a single labeled grid image, for visual
uniqueness review. Pure stdlib + PIL; reads the already-captured screenshots, no
browser needed.

Input:
  chrome-testing/generated/specimens.json   (element -> [{label, file, temporal}])
  chrome-testing/screenshots/specimens/<el>/<NN>-<slug>.png  (static)
  chrome-testing/screenshots/specimens/<el>/<NN>-<slug>.gif  (temporal)

Output:
  chrome-testing/screenshots/contact/<el>.png   one labeled grid per element

Usage: contact_sheet.py [element ...]   (default: all in specimens.json)
"""
import json, os, sys
from PIL import Image, ImageDraw, ImageFont

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SPECS = os.path.join(ROOT, "chrome-testing/generated/specimens.json")
SHOTS = os.path.join(ROOT, "chrome-testing/screenshots/specimens")
OUT = os.path.join(ROOT, "chrome-testing/screenshots/contact")

CELL = 150        # thumbnail size
LABEL_H = 34      # label strip height
PAD = 8
COLS = 7
BG = (26, 26, 46)
CARD = (15, 21, 48)
FG = (230, 230, 230)
ACCENT = (245, 166, 35)

def font(sz):
    for p in ("/System/Library/Fonts/Menlo.ttc",
              "/System/Library/Fonts/Supplemental/Andale Mono.ttf",
              "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf"):
        if os.path.exists(p):
            try: return ImageFont.truetype(p, sz)
            except Exception: pass
    return ImageFont.load_default()

F = font(11)
FT = font(15)

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

def sheet(el, items):
    cells = []
    for it in items:
        base = os.path.splitext(os.path.basename(it["file"]))[0]
        png = os.path.join(SHOTS, el, base + (".gif" if it["temporal"] else ".png"))
        if it["temporal"] and not os.path.exists(png):
            png = os.path.join(SHOTS, el, base + ".gif")
        thumb = load_thumb(png)
        tag = ("▶ " if it["temporal"] else "") + it["label"]
        cells.append((thumb, tag))
    if not cells:
        return None
    n = len(cells)
    cols = min(COLS, n)
    rows = (n + cols - 1) // cols
    cw = CELL + PAD * 2
    ch = CELL + LABEL_H + PAD * 2
    W = cols * cw
    H = rows * ch + 40
    img = Image.new("RGB", (W, H), BG)
    d = ImageDraw.Draw(img)
    d.text((PAD, 8), f"<{el}>  ({n} value-paths)", fill=(22, 199, 154), font=FT)
    for i, (thumb, tag) in enumerate(cells):
        r, c = divmod(i, cols)
        x = c * cw + PAD
        y = r * ch + 40 + PAD
        d.rectangle([x, y, x + CELL, y + CELL + LABEL_H], fill=CARD, outline=(38, 48, 90))
        if thumb:
            img.paste(thumb, (x + (CELL - thumb.width) // 2, y + (CELL - thumb.height) // 2))
        else:
            d.text((x + 8, y + CELL // 2), "(missing)", fill=(200, 80, 80), font=F)
        ty = y + CELL + 3
        for ln in wrap(d, tag, F, CELL - 6):
            d.text((x + 4, ty), ln, fill=ACCENT, font=F)
            ty += 13
    return img

def main():
    spec = json.load(open(SPECS))
    os.makedirs(OUT, exist_ok=True)
    targets = sys.argv[1:] or sorted(spec.keys())
    made = 0
    for el in targets:
        items = spec.get(el)
        if not items:
            continue
        img = sheet(el, items)
        if img:
            img.save(os.path.join(OUT, el + ".png"))
            made += 1
            print(f"  {el}: {len(items)} cells")
    print(f"wrote {made} contact sheets to {OUT}")

if __name__ == "__main__":
    main()
