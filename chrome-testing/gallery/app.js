/* svglab — vanilla SVG element explorer.
   Reimplements the SVG Lab design (no React/DC runtime). Data is driven entirely
   by catalogue.json, which the generator emits from the SVG grammar. The viewer,
   control panel and live editor are wired together client-side: a control or
   preset rebuilds the SVG markup from the current attribute values; editing the
   markup updates the viewer directly. */
(function () {
  'use strict';
  var $ = function (id) { return document.getElementById(id); };

  var DATA = null;            // catalogue.json
  var cur = null;            // current element def
  var values = {};           // current attribute values
  var code = '';             // current SVG markup (editor text = source of truth for the viewer)
  var zoom = 1, panX = 0, panY = 0, dark = true, embed = false;
  var els = [], byId = {};

  // ---- code building: apply attribute values onto the base SVG (string subst) ----
  function escRe(s) { return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'); }
  function escXml(s) { return String(s).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;'); }

  function setAttrInTag(tagStr, key, value) {
    var re = new RegExp('(\\s' + escRe(key) + '=")[^"]*(")');
    if (re.test(tagStr)) return tagStr.replace(re, '$1' + value + '$2');
    return tagStr.replace(/\s*\/?>$/, function (m) { return ' ' + key + '="' + value + '"' + m; });
  }

  function buildCode(el, vals) {
    var out = el.base;
    // Target the SHOWCASED element specifically via its data-lab marker — the base
    // may contain other same-tag elements (a <rect> inside a demo mask, a <path>
    // inside a demo marker, the faint original <use>), and a bare first-match would
    // mutate the wrong one. Fall back to first-match for markerless bases.
    var openRe = new RegExp('<' + escRe(el.tag) + '\\b[^>]*?\\bdata-lab\\b[^>]*?\\/?>');
    var m = openRe.exec(out);
    if (!m) {
      openRe = new RegExp('<' + escRe(el.tag) + '(\\s[^>]*?)?\\/?>');
      m = openRe.exec(out);
    }
    if (m) {
      var tagStr = m[0];
      // Apply EVERY value (controls AND preset-only attributes such as display /
      // visibility), not just the trimmed control set — else a preset that sets a
      // non-control attribute would be a no-op.
      Object.keys(vals).forEach(function (key) {
        if (key === '_content') return;
        var v = vals[key];
        if (v == null) return;
        tagStr = setAttrInTag(tagStr, key, String(v));
      });
      out = out.slice(0, m.index) + tagStr + out.slice(m.index + m[0].length);
    }
    var ca = el.attrs.filter(function (a) { return a.key === '_content'; })[0];
    if (ca && vals._content != null) {
      var cre = new RegExp('(<' + escRe(el.tag) + '(?:\\s[^>]*?)?>)([\\s\\S]*?)(<\\/' + escRe(el.tag) + '>)');
      out = out.replace(cre, function (_, a, b, c) { return a + escXml(vals._content) + c; });
    }
    return out;
  }

  // ---- viewer ----
  function renderViewer() {
    var v = $('viewer');
    if (!v) return;
    v.innerHTML = code;
    var svg = v.querySelector('svg');
    if (svg) { svg.setAttribute('width', '100%'); svg.setAttribute('height', '100%'); svg.style.display = 'block'; }
    v.style.transform = 'translate(' + panX + 'px,' + panY + 'px) scale(' + zoom + ')';
    var cv = $('canvas');
    if (cv) cv.style.background = dark ? '#0c0f0e' : '#f4f7f5';
  }

  // ---- syntax highlight + gutter ----
  function highlight(src) {
    var sp = function (col, t) { return '<span style="color:' + col + '">' + t + '</span>'; };
    var re = /(<\/?)([a-zA-Z][\w:.-]*)((?:\s+[\w:.-]+(?:\s*=\s*"[^"]*")?)*)(\s*\/?>)/g;
    var out = '', last = 0, m;
    while ((m = re.exec(src))) {
      out += escXml(src.slice(last, m.index));
      var attrs = m[3].replace(/([\w:.-]+)(\s*=\s*)("[^"]*")?/g, function (_, n, eq, val) {
        return sp('#6fd3e8', n) + (eq ? sp('#5b6a64', eq) : '') + (val ? sp('#e8c98a', escXml(val)) : '');
      });
      out += sp('#5b6a64', m[1]) + sp('#4ee39a', m[2]) + attrs + sp('#5b6a64', m[4]);
      last = m.index + m[0].length;
    }
    out += escXml(src.slice(last));
    return out;
  }
  function paintEditor() {
    var ta = $('editor'), hl = $('hl'), gut = $('gutter');
    if (!ta) return;
    if (ta.value !== code) ta.value = code;
    var lines = code.split('\n');
    if (hl) hl.innerHTML = highlight(code) + '\n';
    if (gut) gut.textContent = lines.map(function (_, i) { return i + 1; }).join('\n');
    var lc = $('linecount'); if (lc) lc.textContent = lines.length + ' lines';
  }

  // ---- regenerate from values (control/preset path) ----
  function regen() {
    code = buildCode(cur, values);
    renderViewer();
    paintEditor();
  }

  // ---- control panel ----
  function attrRow(a) {
    var row = document.createElement('div');
    row.style.cssText = 'display:flex;align-items:center;gap:12px;padding:9px 0;border-bottom:1px solid #11160f';
    var lab = document.createElement('label');
    lab.style.cssText = 'flex:none;width:104px;font-size:12px;color:#9aa8a0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap';
    lab.textContent = a.label;
    var wrap = document.createElement('div');
    wrap.style.cssText = 'flex:1;min-width:0';
    var val = values[a.key] != null ? values[a.key] : '';

    function commit(v) { values[a.key] = v; regen(); }

    if (a.control === 'paint') {
      var box = document.createElement('div'); box.style.cssText = 'display:flex;gap:7px;align-items:center';
      var col = document.createElement('input'); col.type = 'color';
      col.value = /^#[0-9a-fA-F]{6}$/.test(val) ? val : '#4ee39a';
      var txt = document.createElement('input'); txt.type = 'text'; txt.className = 'fin'; txt.value = val;
      txt.style.cssText += ';flex:1;min-width:0';
      col.oninput = function () { txt.value = col.value; commit(col.value); };
      txt.oninput = function () { if (/^#[0-9a-fA-F]{6}$/.test(txt.value)) col.value = txt.value; commit(txt.value); };
      box.appendChild(col); box.appendChild(txt); wrap.appendChild(box);
    } else if (a.control === 'range') {
      var disp = document.createElement('div');
      disp.style.cssText = 'display:flex;justify-content:flex-end;font-size:11px;color:#4ee39a;margin-bottom:3px';
      disp.textContent = val;
      var r = document.createElement('input'); r.type = 'range';
      r.min = a.min; r.max = a.max; r.step = a.step != null ? a.step : 1; r.value = val;
      r.oninput = function () { disp.textContent = r.value; commit(r.value); };
      wrap.appendChild(disp); wrap.appendChild(r);
    } else if (a.control === 'select') {
      var sel = document.createElement('select'); sel.className = 'fin';
      (a.options || []).forEach(function (o) {
        var opt = document.createElement('option'); opt.value = o; opt.textContent = o;
        if (o === val) opt.selected = true; sel.appendChild(opt);
      });
      sel.onchange = function () { commit(sel.value); };
      wrap.appendChild(sel);
    } else if (a.control === 'number') {
      var n = document.createElement('input'); n.type = 'number'; n.className = 'fin';
      n.style.cssText += ';width:100%'; if (a.min != null) n.min = a.min; if (a.max != null) n.max = a.max;
      if (a.step != null) n.step = a.step; n.value = val;
      n.oninput = function () { commit(n.value); };
      wrap.appendChild(n);
    } else { // text
      var t = document.createElement('input'); t.type = 'text'; t.className = 'fin';
      t.style.cssText += ';width:100%'; t.value = val;
      t.oninput = function () { commit(t.value); };
      wrap.appendChild(t);
    }
    row.appendChild(lab); row.appendChild(wrap);
    return row;
  }

  function renderControls() {
    var c = $('controls'); c.innerHTML = '';
    cur.attrs.forEach(function (a) { c.appendChild(attrRow(a)); });
    var p = $('presets'); p.innerHTML = '';
    (cur.presets || []).forEach(function (pr) {
      var b = document.createElement('button'); b.className = 'preset'; b.textContent = pr.name;
      if (pr.meaning) b.title = pr.meaning;
      b.onclick = function () { applyPreset(pr.values, pr.meaning); };
      p.appendChild(b);
    });
    var pc = $('preset-count'); if (pc) pc.textContent = (cur.presets || []).length ? '· ' + cur.presets.length : '';
  }

  // showMeaning displays the active preset's plain-language caption above the
  // preset list (created lazily so index.html needs no extra markup).
  function showMeaning(meaning) {
    var host = $('presets'); if (!host) return;
    var el = $('preset-meaning');
    if (!el) {
      el = document.createElement('div'); el.id = 'preset-meaning';
      el.style.cssText = 'font-size:12px;color:#9fb0a8;margin:2px 2px 8px;line-height:1.45;min-height:16px';
      host.parentNode.insertBefore(el, host);
    }
    el.textContent = meaning || '';
  }

  function applyPreset(vals, meaning) {
    values = Object.assign({}, cur.defaults, vals);
    renderControls();
    regen();
    showMeaning(meaning);
  }

  // ---- nav rail ----
  // namespaceIds rewrites every id (and #/url(#)/href="#" reference) in a
  // thumbnail SVG to a unique prefix, so the many rail/home thumbnails don't
  // collide on id="slot" with each other OR with the viewer's live SVG — a
  // duplicate id="slot" would make url(#slot)/href="#slot" resolve to the wrong
  // element (filters/gradients/patterns/textPath would silently break).
  var _uid = 0;
  function namespaceIds(svg, uid) {
    return svg
      .replace(/\bid="([^"]+)"/g, 'id="' + uid + '-$1"')
      .replace(/url\(#([^)]+)\)/g, 'url(#' + uid + '-$1)')
      .replace(/href="#([^"]+)"/g, 'href="#' + uid + '-$1"');
  }
  function miniSvg(el) {
    // small thumbnail: the element's default render, fit into 22px, with all ids
    // namespaced so it never steals the viewer's #slot reference.
    var c = namespaceIds(buildCode(el, el.defaults), 'm' + (_uid++));
    return c.replace(/<svg /, '<svg width="100%" height="100%" ');
  }
  function buildRail() {
    var rail = $('rail'); rail.innerHTML = '';
    DATA.groups.forEach(function (g) {
      var cards = els.filter(function (e) { return e.cat === g; });
      if (!cards.length) return;
      var wrap = document.createElement('div'); wrap.style.marginBottom = '2px';
      var btn = document.createElement('button'); btn.className = 'grp-btn';
      var open = true;
      btn.innerHTML = '<span style="width:12px;color:#4ee39a;font-size:10px">▾</span>' +
        '<span style="flex:1;font-family:\'Space Grotesk\',sans-serif;font-weight:600;font-size:12.5px;color:#cdd6d2">' + g + '</span>' +
        '<span style="font-size:10px;color:#3f4d46">' + cards.length + '</span>';
      var list = document.createElement('div'); list.style.paddingLeft = '2px';
      cards.forEach(function (el) {
        var cb = document.createElement('button'); cb.className = 'card-btn'; cb.dataset.id = el.id;
        cb.innerHTML = '<span style="width:22px;height:22px;flex:none;background:#0a0e0c;border:1px solid #1a221e;border-radius:5px;display:flex;align-items:center;justify-content:center;padding:3px;overflow:hidden">' + miniSvg(el) + '</span>' +
          '<span style="flex:1;font-size:12px;color:#cdd6d2;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">' + el.name + '</span>' +
          '<span style="font-size:10px;color:#4a584f">&lt;' + el.tag + '&gt;</span>';
        cb.onclick = function () { location.hash = '#/el/' + el.id; };
        list.appendChild(cb);
      });
      btn.onclick = function () { open = !open; list.style.display = open ? '' : 'none'; btn.firstChild.textContent = open ? '▾' : '▸'; };
      wrap.appendChild(btn); wrap.appendChild(list); rail.appendChild(wrap);
    });
  }
  function markActive(id) {
    document.querySelectorAll('.card-btn').forEach(function (b) {
      b.classList.toggle('active', b.dataset.id === id);
    });
  }

  // ---- open element / breadcrumb ----
  function openElement(id, presetVals) {
    var el = byId[id]; if (!el) return;
    cur = el; values = Object.assign({}, el.defaults, presetVals || {});
    zoom = 1; panX = 0; panY = 0;
    $('bc-cat').textContent = el.cat;
    $('bc-name').textContent = el.name;
    $('bc-tag').textContent = '<' + el.tag + '>';
    $('bc-desc').textContent = el.desc || '';
    $('b-reset').textContent = '100%';
    renderControls();
    regen();
    markActive(id);
  }

  // ---- embed mode (clean, chrome-free viewer for screenshots) ----
  function setEmbed(on) {
    embed = on;
    var d = on ? 'none' : '';
    document.querySelector('header').style.display = d;
    document.querySelectorAll('aside')[0].style.display = d;   // nav rail
    ['bcbar', 'editorpane', 'ctrlpanel', 'vieweroverlay'].forEach(function (id) { var e = $(id); if (e) e.style.display = on ? 'none' : (id === 'vieweroverlay' ? 'flex' : ''); });
    var v = $('viewer');
    if (on) { v.style.width = '440px'; v.style.height = '440px'; $('canvas').style.background = '#0c0f0e'; }
    else { v.style.width = '300px'; v.style.height = '300px'; }
  }

  // ---- home view ----
  function showPage(which) {
    $('home').style.display = which === 'home' ? 'block' : 'none';
    $('elpage').style.display = which === 'home' ? 'none' : 'flex';
  }
  function featuredIds() {
    var pref = ['rect', 'circle', 'path', 'polygon', 'linearGradient', 'radialGradient',
      'feGaussianBlur', 'feDropShadow', 'text', 'pattern', 'feTurbulence', 'animate'];
    var out = pref.filter(function (id) { return byId[id]; });
    els.forEach(function (e) { if (out.length < 12 && out.indexOf(e.id) < 0) out.push(e.id); });
    return out.slice(0, 12);
  }
  function renderHome() {
    var feat = $('featured'); if (!feat) return;
    feat.innerHTML = '';
    featuredIds().forEach(function (id) {
      var el = byId[id]; if (!el) return;
      var card = document.createElement('button');
      card.style.cssText = 'text-align:left;background:#0f1311;border:1px solid #1c2422;border-radius:12px;overflow:hidden;cursor:pointer;padding:0;font-family:inherit;transition:border-color .15s';
      card.onmouseenter = function () { card.style.borderColor = '#33463c'; };
      card.onmouseleave = function () { card.style.borderColor = '#1c2422'; };
      card.innerHTML =
        '<div style="height:118px;background:#0c0f0e;display:flex;align-items:center;justify-content:center;padding:18px;border-bottom:1px solid #161c19">' +
        '<div style="width:100%;height:100%;max-width:88px;display:flex;align-items:center;justify-content:center">' + miniSvg(el) + '</div></div>' +
        '<div style="padding:11px 13px">' +
        '<div style="font-size:12px;color:#4ee39a;font-family:\'JetBrains Mono\',monospace">&lt;' + el.tag + '&gt;</div>' +
        '<div style="font-family:\'Space Grotesk\',sans-serif;font-weight:600;font-size:14px;color:#eef4f0;margin-top:2px">' + el.name + '</div>' +
        '<div style="font-size:10.5px;color:#6b7a73;margin-top:3px">' + el.cat + '</div></div>';
      card.onclick = function () { location.hash = '#/el/' + el.id; };
      feat.appendChild(card);
    });
  }

  // ---- routing ----
  function route() {
    var h = location.hash.replace(/^#\/?/, '');
    var parts = h.split('/').filter(Boolean);
    if (parts[0] === 'embed' && parts[1]) {
      setEmbed(true);
      showPage('element');
      var el = byId[parts[1]];
      var idx = parseInt(parts[2] || '0', 10);
      var preset = el && el.presets && el.presets[idx] ? el.presets[idx].values : null;
      openElement(parts[1], preset);
      return;
    }
    setEmbed(false);
    if (parts[0] === 'el' && parts[1]) { showPage('element'); openElement(parts[1]); return; }
    showPage('home');
    markActive(null);
    renderHome();
  }

  // ---- viewer interactions ----
  function wireViewer() {
    $('b-canvas').onclick = function () { dark = !dark; $('b-canvas').textContent = dark ? 'dark bg' : 'light bg'; renderViewer(); };
    $('b-zoomin').onclick = function () { zoom = Math.min(4, +(zoom + 0.2).toFixed(2)); renderViewer(); $('b-reset').textContent = Math.round(zoom * 100) + '%'; };
    $('b-zoomout').onclick = function () { zoom = Math.max(0.3, +(zoom - 0.2).toFixed(2)); renderViewer(); $('b-reset').textContent = Math.round(zoom * 100) + '%'; };
    $('b-reset').onclick = function () { zoom = 1; panX = 0; panY = 0; renderViewer(); $('b-reset').textContent = '100%'; };
    var cv = $('canvas');
    cv.onwheel = function (e) { e.preventDefault(); zoom = Math.min(4, Math.max(0.3, +(zoom + (e.deltaY > 0 ? -0.12 : 0.12)).toFixed(2))); renderViewer(); $('b-reset').textContent = Math.round(zoom * 100) + '%'; };
    cv.onmousedown = function (e) {
      if (e.button !== 0) return;
      var sx = e.clientX, sy = e.clientY, px = panX, py = panY; cv.style.cursor = 'grabbing';
      function mv(ev) { panX = px + (ev.clientX - sx); panY = py + (ev.clientY - sy); renderViewer(); }
      function up() { document.removeEventListener('mousemove', mv); document.removeEventListener('mouseup', up); cv.style.cursor = 'grab'; }
      document.addEventListener('mousemove', mv); document.addEventListener('mouseup', up);
    };
    // editor: edits drive the viewer directly
    var ta = $('editor');
    ta.oninput = function () { code = ta.value; renderViewer(); paintEditor(); };
    ta.onscroll = function () { $('hl').scrollTop = ta.scrollTop; $('hl').scrollLeft = ta.scrollLeft; $('gutter').scrollTop = ta.scrollTop; };
    // export
    $('b-copy').onclick = function () { try { navigator.clipboard.writeText(code); } catch (e) {} var b = $('b-copy'); b.textContent = 'copied'; setTimeout(function () { b.textContent = 'copy'; }, 1200); };
    $('b-svg').onclick = function () { save(new Blob([code], { type: 'image/svg+xml' }), (cur ? cur.id : 'svg') + '.svg'); };
    $('b-png').onclick = downloadPng;
    $('home-link').onclick = function () { location.hash = '#/'; };
  }
  function save(blob, name) { var u = URL.createObjectURL(blob); var a = document.createElement('a'); a.href = u; a.download = name; document.body.appendChild(a); a.click(); a.remove(); setTimeout(function () { URL.revokeObjectURL(u); }, 1500); }
  function downloadPng() {
    var blob = new Blob([code], { type: 'image/svg+xml;charset=utf-8' }), url = URL.createObjectURL(blob), img = new Image();
    img.onload = function () {
      var S = 720, c = document.createElement('canvas'); c.width = S; c.height = S;
      var ctx = c.getContext('2d'); ctx.fillStyle = dark ? '#0e1211' : '#fff'; ctx.fillRect(0, 0, S, S);
      ctx.drawImage(img, 80, 80, S - 160, S - 160);
      c.toBlob(function (b) { save(b, (cur ? cur.id : 'svg') + '.png'); URL.revokeObjectURL(url); });
    };
    img.onerror = function () { URL.revokeObjectURL(url); };
    img.src = url;
  }

  // ---- boot ----
  fetch('./catalogue.json').then(function (r) { return r.json(); }).then(function (d) {
    DATA = d; els = d.elements || []; byId = {};
    els.forEach(function (e) { byId[e.id] = e; });
    $('el-count').textContent = els.length;
    $('grp-count').textContent = (d.groups || []).filter(function (g) { return els.some(function (e) { return e.cat === g; }); }).length;
    buildRail();
    wireViewer();
    window.addEventListener('hashchange', route);
    route();
  }).catch(function (err) {
    document.body.insertAdjacentHTML('beforeend', '<div style="position:fixed;top:60px;left:260px;color:#ef7a7a;font-size:13px">Failed to load catalogue.json: ' + err + '</div>');
  });

  // expose for the screenshot harness
  window.svglab = { open: function (id, idx) { location.hash = '#/embed/' + id + '/' + (idx || 0); } };
})();
