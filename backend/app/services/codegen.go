// Package services contains application business logic.
package services

import (
	"fmt"
	"html"
	"sort"
	"strings"
)

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

// ── attribute builder ───────────────────────────────────────────

// attrStr builds the attribute string. class always first, then named attrs in
// the order given (deterministic). Only non-empty values are emitted.
func (g *LandingCodegen) attrStr(n nodeMap, extra ...string) string {
	var b strings.Builder
	if cls := n.classStr(); cls != "" {
		fmt.Fprintf(&b, ` class=%q`, cls)
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
<script src="https://cdn.tailwindcss.com"></script>
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
	raw, ok := n["classes"]
	if !ok {
		return ""
	}
	parts := toStringSlice(raw)
	return strings.Join(dedup(parts), " ")
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
