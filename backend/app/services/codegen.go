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
// ponytail: component instance resolution (componentId → master tree) deferred to
// Phase 5. For now instances render as plain frames.
type LandingCodegen struct{}

func NewLandingCodegen() *LandingCodegen { return &LandingCodegen{} }

// Generate produces a full HTML document from a tree.
// title = page name used in <title>.
func (g *LandingCodegen) Generate(tree map[string]any, title string) string {
	body := g.renderNode(asNode(tree["root"]))
	return g.document(title, body)
}

// GenerateFragment renders only the body inner HTML (no doctype/html/head wrapper).
// Useful for diffs, tests, and embedded previews.
func (g *LandingCodegen) GenerateFragment(tree map[string]any) string {
	return g.renderNode(asNode(tree["root"]))
}

// renderNode dispatches a single node to its element renderer. depth controls indent.
func (g *LandingCodegen) renderNode(n nodeMap) string {
	if n == nil {
		return ""
	}
	t := n.typeStr()
	switch t {
	case "text":
		return g.renderLeaf("span", n)
	case "heading":
		return g.renderHeading(n)
	case "image":
		return g.renderSelfClosing("img", n, "src", "alt")
	case "button":
		return g.renderLeaf("button", n)
	case "link":
		return g.renderLink(n)
	case "section":
		return g.renderContainer("section", n)
	case "frame", "":
		return g.renderContainer("div", n)
	default:
		// unknown types fall back to a neutral div
		return g.renderContainer("div", n)
	}
}

// renderContainer renders a node that has children.
func (g *LandingCodegen) renderContainer(tag string, n nodeMap) string {
	var b strings.Builder
	b.WriteString("<" + tag + g.attrStr(n) + ">")
	for _, c := range n.children() {
		b.WriteString(g.renderNode(asNode(c)))
	}
	b.WriteString("</" + tag + ">")
	return b.String()
}

// renderLeaf renders a node whose only content is props.text.
func (g *LandingCodegen) renderLeaf(tag string, n nodeMap) string {
	return fmt.Sprintf(
		"<%s%s>%s</%s>",
		tag, g.attrStr(n), html.EscapeString(n.propStr("text")), tag,
	)
}

// renderHeading maps props.level → h1..h6 (default h2).
func (g *LandingCodegen) renderHeading(n nodeMap) string {
	lvl := n.propInt("level")
	if lvl < 1 || lvl > 6 {
		lvl = 2
	}
	tag := fmt.Sprintf("h%d", lvl)
	return g.renderLeaf(tag, n)
}

// renderSelfClosing renders void elements (img, input, br).
func (g *LandingCodegen) renderSelfClosing(tag string, n nodeMap, attrs ...string) string {
	return "<" + tag + g.attrStr(n, attrs...) + " />"
}

// renderLink prefers children, falls back to props.text.
func (g *LandingCodegen) renderLink(n nodeMap) string {
	var inner string
	kids := n.children()
	if len(kids) > 0 {
		for _, c := range kids {
			inner += g.renderNode(asNode(c))
		}
	} else {
		inner = html.EscapeString(n.propStr("text"))
	}
	return fmt.Sprintf("<a%s>%s</a>", g.attrStr(n, "href"), inner)
}

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

// document wraps body HTML with a minimal HTML5 doc + Tailwind CDN.
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
%s
</body>
</html>`
	return fmt.Sprintf(tpl, html.EscapeString(title), body)
}

// --- node map helpers ---

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

// classStr joins node.classes deterministically. If classes is an []any (from
// JSON unmarshal) we sort-copy to avoid map-order drift, but preserve original
// order since a slice already is ordered.
func (n nodeMap) classStr() string {
	if n == nil {
		return ""
	}
	raw, ok := n["classes"]
	if !ok {
		return ""
	}
	parts := toStringSlice(raw)
	return strings.Join(parts, " ")
}

func (n nodeMap) children() []any {
	if n == nil {
		return nil
	}
	c, _ := n["children"].([]any)
	return c
}

// toStringSlice accepts []string, []any (JSON shape), and returns stable slice.
// ponytail: no dedup yet — codegen output matches input order. Add dedup when
// style panel emits duplicates (Phase 6).
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

// SortedClasses is exported for tests and Phase 7 deterministic comparisons.
// It returns the class list sorted alphabetically — used to detect class-set
// equality without caring about order.
func SortedClasses(raw any) []string {
	out := toStringSlice(raw)
	sort.Strings(out)
	return out
}
