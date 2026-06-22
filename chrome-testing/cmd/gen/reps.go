package main

// Representative values substitute ONLY the atomic, open-ended leaf *Type rules
// of the SVG value grammar — the scalars and tokens that bottom out in infinite
// character ranges (numbers, lengths, colors, idents, IRIs, free text). Every
// COMPOSITE above a leaf (PaintType, TransformList, SvgPath, Points, ViewBox,
// PreserveAspectRatio, every element, every attribute group) is walked
// structurally by the renderer and bottoms out at these leaves.
//
// Keys are the COMPILED simple proto message names (gluon's PascalCase of the
// EBNF rule name). Most leaf rules keep their spelling (LengthType, ColorType,
// …) but gluon lowercases a capital run that follows a digit run, so the
// `Bcp47Type` rule compiles to the message `Bcp47type` — verified against
// proto/svg.proto. The renderer intercepts a message by simpleName(fqn) BEFORE
// walking it, so the (un-scalarized) leaf body is never traversed in gen.
//
// CRITICAL: every sample is the VALUE TEXT ONLY — the text that sits INSIDE the
// attribute quotes. The surrounding ` name="` and `"` come from the grammar's
// markup terminals (emitted via the kw map), so reps must be quote-safe: NO
// embedded double-quote characters (`"`). A `"` here would break the
// `name="value"` well-formedness invariant.
var reps = map[string][]string{
	// ── numeric leaves ──────────────────────────────────────────────────────
	"NumberType":  {"0", "1", "-1", "3.14", "0.5", "2"},
	"IntegerType": {"0", "1", "-1", "100", "3", "10"},

	// ── length / percentage / coordinate ───────────────────────────────────
	"LengthType":           {"10", "24px", "2em", "1.5rem", "50%", "12pt"},
	"PercentageType":       {"50%", "100%", "25%", "0%", "80%"},
	"LengthPercentageType": {"10", "24px", "2em", "50%", "100%"},
	"CoordinateType":       {"0", "10", "-5.5", "100px", "50%", "2.5em"},

	// ── angle / time ────────────────────────────────────────────────────────
	"AngleType": {"0deg", "45deg", "90deg", "180deg", "0.25turn"},
	"TimeType":  {"0s", "0.3s", "1s", "200ms", "2.5s"},

	// ── color / paint references ────────────────────────────────────────────
	"ColorType": {"#e94560", "red", "currentColor", "#16c79a", "#0f3460"},

	// ── IRI / URL references ────────────────────────────────────────────────
	"IriType": {"#ref1", "#grad1", "#marker1", "#clip1"},
	"UrlType": {"url(#id)", "url(#grad1)", "url(#pattern1)"},

	// ── strings / identifiers / names ───────────────────────────────────────
	"StringType":      {"label", "Aa", "sample", "specimen"},
	"CustomIdentType": {"blur1", "result1", "myFilter", "shadow", "out"},
	"XmlNameType":     {"circle1", "grad-a", "myId", "node3", "r1"},
	"Bcp47type":       {"en", "fr-CA", "de", "zh-Hans", "pt-BR"},

	// ── number / length lists (rendered inside one attribute value) ─────────
	"ListOfNumbersType":        {"1 2 3 4", "0 1 0 0", "1 0.5", "0 0 1 1 0", "2 4 6"},
	"ListOfLengthsType":        {"10px 20px", "5% 10%", "1em 2em 3em", "0 4px"},
	"NumberOptionalNumberType": {"2", "2 3", "0.5", "3 3", "1 2"},
	"DasharrayType":            {"5 3", "10 5 2", "4", "5% 10%", "8 4 2 4"},

	// ── SMIL event symbols (begin/end value) ────────────────────────────────
	"EventSymbolType": {"click", "mouseenter", "focus", "beginEvent", "endEvent", "load"},

	// ── raw character data (text content between tags; NO markup, NO quotes) ─
	"CharacterDataType": {"Example", "Hello SVG", "Aa", "Sample Text"},
}
