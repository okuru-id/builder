// Package services contains application business logic.
package services

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"sort"
	"strings"
)

//go:embed icons.json
var iconsJSON []byte

// iconPaths holds the parsed Tabler icon map (name → list of [tag, attrs]).
// Loaded once at init from icons.json, generated from frontend icon-map.ts.
var iconPaths map[string][][]any

func init() {
	_ = json.Unmarshal(iconsJSON, &iconPaths)
}

// LandingCodegen converts a landing-page tree (map[string]any JSON shape) into
// a standalone HTML document with Tailwind CDN. Deterministic: same tree → same bytes.
//
// Tree shape: { "root": Node } where Node =
//   { id, type, name, props: {}, classes: [], children: [] }
//
// ponytail: codegen stays pure (no DB). Component instance resolution happens
// upstream in services.ResolveComponentInstances before Generate is called.
type LandingCodegen struct{}

func NewLandingCodegen() *LandingCodegen { return &LandingCodegen{} }

// Generate produces a full HTML document from a tree with 2-space indent.
// title = page name used in <title>.
func (g *LandingCodegen) Generate(tree map[string]any, title string) string {
	body := g.renderNodeDepth(asNode(tree["root"]), 2)
	return g.document(title, body)
}

// GenerateFragment renders only the body inner HTML (no wrapper).
// Useful for diffs, tests, and embedded previews.
func (g *LandingCodegen) GenerateFragment(tree map[string]any) string {
	return g.renderNodeDepth(asNode(tree["root"]), 0)
}

// ── Safelist ────────────────────────────────────────────────────

// GenerateSafelist walks the tree and collects every distinct Tailwind utility
// class (class="...") from every node. Returns sorted, deduplicated slice.
// Useful for Tailwind JIT safelist configuration so no utility is purged.
func (g *LandingCodegen) GenerateSafelist(tree map[string]any) []string {
	set := make(map[string]struct{})
	g.collectClasses(asNode(tree["root"]), set)
	out := make([]string, 0, len(set))
	for c := range set {
		out = append(out, c)
	}
	sort.Strings(out)
	return out
}

func (g *LandingCodegen) collectClasses(n nodeMap, out map[string]struct{}) {
	if n == nil {
		return
	}
	if raw, ok := n["classes"]; ok {
		for _, s := range toStringSlice(raw) {
			out[s] = struct{}{}
		}
	}
	for _, c := range n.children() {
		g.collectClasses(asNode(c), out)
	}
}

// ── internal renderers (depth-aware) ────────────────────────────

func (g *LandingCodegen) renderNodeDepth(n nodeMap, depth int) string {
	if n == nil {
		return ""
	}
	if n.hidden() {
		return ""
	}
	t := n.typeStr()
	switch t {
	case "text":
		return g.renderLeafDepth("span", n, depth)
	case "heading":
		return g.renderHeadingDepth(n, depth)
	case "image":
		return g.renderSelfClosingDepth("img", n, depth, "src", "alt")
	case "button":
		return g.renderLeafDepth("button", n, depth)
	case "link":
		return g.renderLinkDepth(n, depth)
	case "section":
		return g.renderContainerDepth("section", n, depth)
	case "frame", "":
		return g.renderContainerDepth("div", n, depth)
	case "divider":
		return g.renderSelfClosingDepth("hr", n, depth)
	case "form":
		return g.renderFormDepth(n, depth)
	case "input":
		return g.renderInputDepth(n, depth)
	case "icon":
		return g.renderIconDepth(n, depth)
	default:
		return g.renderContainerDepth("div", n, depth)
	}
}

func (g *LandingCodegen) renderContainerDepth(tag string, n nodeMap, depth int) string {
	indent := strings.Repeat("  ", depth)
	var b strings.Builder
	b.WriteString(indent + "<" + tag + g.attrStr(n) + ">\n")
	for _, c := range n.children() {
		b.WriteString(g.renderNodeDepth(asNode(c), depth+1))
	}
	b.WriteString(indent + "</" + tag + ">\n")
	return b.String()
}

func (g *LandingCodegen) renderLeafDepth(tag string, n nodeMap, depth int) string {
	indent := strings.Repeat("  ", depth)
	return fmt.Sprintf("%s<%s%s>%s</%s>\n",
		indent, tag, g.attrStr(n), html.EscapeString(n.propStr("text")), tag)
}

func (g *LandingCodegen) renderHeadingDepth(n nodeMap, depth int) string {
	lvl := n.propInt("level")
	if lvl < 1 || lvl > 6 {
		lvl = 2
	}
	tag := fmt.Sprintf("h%d", lvl)
	return g.renderLeafDepth(tag, n, depth)
}

func (g *LandingCodegen) renderSelfClosingDepth(tag string, n nodeMap, depth int, attrs ...string) string {
	indent := strings.Repeat("  ", depth)
	return indent + "<" + tag + g.attrStr(n, attrs...) + " />\n"
}

func (g *LandingCodegen) renderLinkDepth(n nodeMap, depth int) string {
	indent := strings.Repeat("  ", depth)
	var inner string
	kids := n.children()
	if len(kids) > 0 {
		// Links with children render as a container block
		var b strings.Builder
		b.WriteString(indent + "<a" + g.attrStr(n, "href") + ">\n")
		for _, c := range kids {
			b.WriteString(g.renderNodeDepth(asNode(c), depth+1))
		}
		b.WriteString(indent + "</a>\n")
		return b.String()
	}
	inner = html.EscapeString(n.propStr("text"))
	return fmt.Sprintf("%s<a%s>%s</a>\n", indent, g.attrStr(n, "href"), inner)
}

// ── form, input, icon renderers ─────────────────────────────────

func (g *LandingCodegen) renderFormDepth(n nodeMap, depth int) string {
	indent := strings.Repeat("  ", depth)
	au := g.attrStr(n)
	action := n.propStr("action")
	method := n.propStr("method")
	if method == "" {
		method = "POST"
	}
	var b strings.Builder
	if action != "" {
		b.WriteString(fmt.Sprintf(` action=%q`, action))
	}
	b.WriteString(fmt.Sprintf(` method=%q`, method))
	var inner strings.Builder
	for _, c := range n.children() {
		inner.WriteString(g.renderNodeDepth(asNode(c), depth+1))
	}
	return fmt.Sprintf("%s<form%s%s>\n%s%s</form>\n", indent, au, b.String(), inner.String(), indent)
}

func (g *LandingCodegen) renderInputDepth(n nodeMap, depth int) string {
	indent := strings.Repeat("  ", depth)
	label := n.propStr("label")
	inputType := n.propStr("inputType")
	if inputType == "" {
		inputType = "text"
	}
	placeholder := n.propStr("placeholder")
	req := n.propStr("required")
	var b strings.Builder
	b.WriteString(indent + "<div" + g.attrStr(n) + ">\n")
	if label != "" {
		b.WriteString(fmt.Sprintf("%s  <label class=\"text-sm font-medium\">%s</label>\n", indent, html.EscapeString(label)))
	}
	b.WriteString(indent + "  <input")
	b.WriteString(fmt.Sprintf(` type=%q`, inputType))
	if placeholder != "" {
		b.WriteString(fmt.Sprintf(` placeholder=%q`, html.EscapeString(placeholder)))
	}
	if req == "true" || req == "1" {
		b.WriteString(` required`)
	}
	b.WriteString(` class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm"`)
	b.WriteString(" />\n")
	b.WriteString(indent + "</div>\n")
	return b.String()
}

func (g *LandingCodegen) renderIconDepth(n nodeMap, depth int) string {
	indent := strings.Repeat("  ", depth)
	props, _ := n["props"].(map[string]any)
	name, _ := props["icon"].(string)
	variant, _ := props["iconVariant"].(string)
	if variant != "outline" && variant != "filled" {
		variant = "outline"
	}
	lookupKey := name
	if variant == "filled" {
		lookupKey = name + "Filled"
	}
	fillAttr, strokeAttr := "none", "currentColor"
	if variant == "filled" {
		fillAttr, strokeAttr = "currentColor", "none"
	}
	var svg string
	segs, ok := iconPaths[lookupKey]
	if !ok || len(segs) == 0 {
		// Fallback to outline variant if filled missing.
		segs, ok = iconPaths[name]
	}
	if ok && len(segs) > 0 {
		var b strings.Builder
		b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24" stroke-width="2" stroke="` + strokeAttr + `" fill="` + fillAttr + `" stroke-linecap="round" stroke-linejoin="round">`)
		for _, seg := range segs {
			if len(seg) < 2 {
				continue
			}
			tag, _ := seg[0].(string)
			attrs, _ := seg[1].(map[string]any)
			b.WriteString("<" + tag)
			for k, v := range attrs {
				if k == "key" {
					continue
				}
				b.WriteString(fmt.Sprintf(` %s=%q`, k, fmt.Sprint(v)))
			}
			b.WriteString(" />")
		}
		b.WriteString("</svg>")
		svg = b.String()
	} else {
		// Unknown/missing icon: render its name as text, or a placeholder dot.
		if name != "" {
			svg = html.EscapeString(name)
		} else {
			svg = `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="4" fill="currentColor" /></svg>`
		}
	}
	return fmt.Sprintf("%s<span%s>%s</span>\n", indent, g.attrStr(n), svg)
}

// ── attribute builder ───────────────────────────────────────────

// attrStr builds the attribute string. class always first, then the
// data-bp-hide visibility hook (if any), then named attrs in the order given
// (deterministic). Only non-empty values are emitted.
func (g *LandingCodegen) attrStr(n nodeMap, extra ...string) string {
	var b strings.Builder
	if cls := n.classStr(); cls != "" {
		fmt.Fprintf(&b, ` class=%q`, cls)
	}
	if bp := n.hiddenOn(); bp != "" {
		fmt.Fprintf(&b, ` data-bp-hide=%q`, bp)
	}
	for _, name := range extra {
		if v := n.propStr(name); v != "" {
			fmt.Fprintf(&b, ` %s=%q`, name, v)
		}
	}
	return b.String()
}

// ── document wrapper ────────────────────────────────────────────

func (g *LandingCodegen) document(title, body string) string {
	const tpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>%s</title>
<link rel="preconnect" href="https://fonts.googleapis.com" />
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Poppins:wght@400;500;600;700&family=Roboto:wght@400;500;700&family=Montserrat:wght@400;500;600;700&family=Open+Sans:wght@400;500;600;700&family=Playfair+Display:wght@400;500;600;700&family=Lora:wght@400;500;600;700&family=Merriweather:wght@400;700&family=JetBrains+Mono:wght@400;500;700&family=Source+Sans+3:wght@400;500;600;700&family=Nunito:wght@400;500;600;700&family=Raleway:wght@400;500;600;700&family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap" rel="stylesheet" />
<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
<style>
/* Per-breakpoint visibility. !important wins over the node's own flex/grid/block
   utilities (classic Tailwind source-order conflict). Breakpoints match the
   builder: mobile <768, tablet 768..1023, desktop ≥1024. */
@media (max-width:767px){[data-bp-hide~="m"]{display:none!important}}
@media (min-width:768px) and (max-width:1023px){[data-bp-hide~="t"]{display:none!important}}
@media (min-width:1024px){[data-bp-hide~="d"]{display:none!important}}
/* Google Font utility classes (mirrors frontend/src/style.css). Tailwind
   browser CDN only knows font-sans/serif/mono, so custom Google Fonts use
   the gfont-* prefix. */
.gfont-inter{font-family:'Inter',sans-serif}
.gfont-poppins{font-family:'Poppins',sans-serif}
.gfont-roboto{font-family:'Roboto',sans-serif}
.gfont-montserrat{font-family:'Montserrat',sans-serif}
.gfont-opensans{font-family:'Open Sans',sans-serif}
.gfont-playfair{font-family:'Playfair Display',serif}
.gfont-lora{font-family:'Lora',serif}
.gfont-merriweather{font-family:'Merriweather',serif}
.gfont-jetbrains{font-family:'JetBrains Mono',monospace}
.gfont-sourcesans{font-family:'Source Sans 3',sans-serif}
.gfont-nunito{font-family:'Nunito',sans-serif}
.gfont-raleway{font-family:'Raleway',sans-serif}
.gfont-jakarta{font-family:'Plus Jakarta Sans',sans-serif}
</style>
</head>
<body>
%s</body>
</html>`
	return fmt.Sprintf(tpl, html.EscapeString(title), body)
}

// ── node map helpers ────────────────────────────────────────────

type nodeMap map[string]any

func asNode(v any) nodeMap {
	m, _ := v.(map[string]any)
	return m
}

func (n nodeMap) typeStr() string {
	s, _ := n["type"].(string)
	return s
}

func (n nodeMap) hidden() bool {
	b, _ := n["hidden"].(bool)
	return b
}

// hiddenOn returns the breakpoints a node should be hidden on as a single
// whitespace-separated token string (e.g. "m", "t", "m t"). Emitted as the
// `data-bp-hide` attribute and matched by a scoped <style> block with
// !important media queries. This is bulletproof against the node's own
// display utilities (flex/grid/block) — the classic Tailwind source-order
// footgun where `flex` overrides `hidden`.
func (n nodeMap) hiddenOn() string {
	raw, ok := n["hiddenOn"]
	if !ok {
		return ""
	}
	var parts []string
	for _, bp := range toStringSlice(raw) {
		switch bp {
		case "mobile":
			parts = append(parts, "m")
		case "tablet":
			parts = append(parts, "t")
		case "desktop":
			parts = append(parts, "d")
		}
	}
	return strings.Join(parts, " ")
}

func (n nodeMap) propStr(key string) string {
	if n == nil {
		return ""
	}
	v, ok := n["props"]
	if !ok {
		return ""
	}
	p, ok := v.(map[string]any)
	if !ok {
		return ""
	}
	switch t := p[key].(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	default:
		return fmt.Sprintf("%v", t)
	}
}

func (n nodeMap) propInt(key string) int {
	if n == nil {
		return 0
	}
	v, ok := n["props"]
	if !ok {
		return 0
	}
	p, ok := v.(map[string]any)
	if !ok {
		return 0
	}
	switch t := p[key].(type) {
	case float64:
		return int(t)
	case int:
		return t
	default:
		return 0
	}
}

// classStr joins deduplicated classes deterministically. Duplicates (which
// can arise from style panel interactions) are silently dropped, first
// occurrence wins (stable order in, stable order out).
func (n nodeMap) classStr() string {
	if n == nil {
		return ""
	}
	return strings.Join(dedup(toStringSlice(n["classes"])), " ")
}

func (n nodeMap) children() []any {
	if n == nil {
		return nil
	}
	c, _ := n["children"].([]any)
	return c
}

// toStringSlice accepts []string, []any (JSON shape), and returns stable slice.
func toStringSlice(raw any) []string {
	switch t := raw.(type) {
	case []string:
		return t
	case []any:
		out := make([]string, 0, len(t))
		for _, v := range t {
			if s, ok := v.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

// dedup removes duplicates while preserving order. First occurrence wins.
func dedup(in []string) []string {
	if len(in) < 2 {
		return in
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// SortedClasses is exported for tests. Returns class list sorted alphabetically.
func SortedClasses(raw any) []string {
	out := toStringSlice(raw)
	sort.Strings(out)
	return out
}
