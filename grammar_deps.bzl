"""Cross-grammar build dependencies for the SVG grammar.

The SVG grammar embeds CSS at the <style> element seam without physically
merging their EBNF sources. Evaluated by gluon's v2/builddep loader; two
pipelines consume the one record:

  * lang/cmd/genproto (compile/emit): the externalize pass retypes the <style>
    content field from .svg.CssStyleSheet to the imported .css.CssStyleSheet,
    drops the local opaque placeholder, and adds `import "css.proto"` to
    svg.proto.

  * gluon Metaparser.ResolveDependencies RPC (parse/compose): streams the CSS
    EBNF DocumentDescriptors so a merged grammar can descend into the CSS inside
    an SVG <style>.

Paths are filesystem-relative to this file's directory (the proto-svg repo
root, sibling-checkout layout).
"""

GRAMMAR_DEPS = [
    struct(
        name = "css",
        grammar_srcs = "../proto-css/lang",
        proto = "css.proto",
        proto_include = "../proto-css/proto",
        proto_package = "css",
        go_package = "github.com/accretional/proto-css/proto/pb/css",
        external_rules = [
            # <style> content — the full stylesheet.
            "CssStyleSheet",
            # Presentation-attribute color values (fill / stroke / color /
            # stop-color / flood-color / lighting-color). SVG's own structured
            # ColorType is dropped in favour of proto-css's richer <color>.
            "ColorType",
        ],
    ),
]
