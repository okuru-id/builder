package services_test

import (
	"strings"
	"testing"

	"okuru/app/services"
)

func tree(root map[string]any) map[string]any {
	return map[string]any{"root": root}
}

func node(typ, name string, overrides ...map[string]any) map[string]any {
	n := map[string]any{
		"id":       "n1",
		"type":     typ,
		"name":     name,
		"props":    map[string]any{},
		"classes":  []any{},
		"children": []any{},
	}
	for _, o := range overrides {
		for k, v := range o {
			n[k] = v
		}
	}
	return n
}

// ── Tests ─────────────────────────────────────────────────────────

func TestGenerateHTML_EmptyTree(t *testing.T) {
	cg := services.NewLandingCodegen()
	tr := tree(node("frame", "root"))
	html := cg.Generate(tr, "Test")
	if !strings.Contains(html, "<div") {
		t.Error("expected a <div> wrapper, got:", html)
	}
	if !strings.Contains(html, "Test") {
		t.Error("expected title 'Test', got:", html)
	}
}

func TestGenerateHTML_Deterministic(t *testing.T) {
	cg := services.NewLandingCodegen()
	tr := tree(node("frame", "root", map[string]any{
		"children": []any{
			node("text", "a", map[string]any{"props": map[string]any{"text": "Hello"}}),
			node("button", "b", map[string]any{"props": map[string]any{"text": "Click"}}),
		},
	}))
	html1 := cg.Generate(tr, "Test")
	html2 := cg.Generate(tr, "Test")
	if html1 != html2 {
		t.Error("Generate is not deterministic")
	}
}

func TestGenerateFragment_Deterministic(t *testing.T) {
	cg := services.NewLandingCodegen()
	tr := tree(node("frame", "root", map[string]any{
		"children": []any{
			node("text", "a", map[string]any{"props": map[string]any{"text": "Hello"}}),
		},
	}))
	frag1 := cg.GenerateFragment(tr)
	frag2 := cg.GenerateFragment(tr)
	if frag1 != frag2 {
		t.Error("GenerateFragment is not deterministic")
	}
}

func TestGenerateHTML_NestedFrame(t *testing.T) {
	cg := services.NewLandingCodegen()
	inner := node("frame", "inner", map[string]any{
		"children": []any{
			node("text", "t", map[string]any{"props": map[string]any{"text": "hi"}}),
		},
	})
	tr := tree(node("frame", "root", map[string]any{
		"children": []any{inner},
	}))
	html := cg.Generate(tr, "Nested")
	if !strings.Contains(html, "hi") {
		t.Error("expected nested text content, got:", html)
	}
}

func TestGenerateHTML_NoShadcnImports(t *testing.T) {
	cg := services.NewLandingCodegen()
	tr := tree(node("frame", "root", map[string]any{
		"children": []any{
			node("text", "t", map[string]any{"props": map[string]any{"text": "ok"}}),
		},
	}))
	html := cg.Generate(tr, "Test")
	if strings.Contains(html, "@/components/ui") {
		t.Error("output must not contain shadcn-vue imports")
	}
}

func TestGenerateHTML_ClassesPreserved(t *testing.T) {
	cg := services.NewLandingCodegen()
	n := node("frame", "root")
	n["classes"] = []any{"flex", "gap-4", "p-4"}
	tr := tree(n)
	html := cg.Generate(tr, "Test")
	if !strings.Contains(html, `class="flex gap-4 p-4"`) {
		t.Error("expected class string, got:", html)
	}
}

func TestGenerateHTML_ClassDedup(t *testing.T) {
	cg := services.NewLandingCodegen()
	n := node("frame", "root")
	n["classes"] = []any{"flex", "gap-4", "flex", "p-4", "gap-4"}
	tr := tree(n)
	html := cg.Generate(tr, "Dedup")
	// "flex" and "gap-4" each appear twice but should emit only once
	if !strings.Contains(html, `class="flex gap-4 p-4"`) {
		t.Error("expected deduplicated class string, got:", html)
	}
}

func TestGenerateHTML_SelfClosingImage(t *testing.T) {
	cg := services.NewLandingCodegen()
	img := node("image", "img")
	img["props"] = map[string]any{"src": "https://example.com/pic.jpg", "alt": "A pic"}
	tr := tree(node("frame", "root", map[string]any{"children": []any{img}}))
	html := cg.Generate(tr, "Img")
	if !strings.Contains(html, "<img ") || !strings.Contains(html, "/>") {
		t.Error("expected self-closing img, got:", html)
	}
	if !strings.Contains(html, `src="https://example.com/pic.jpg"`) {
		t.Error("expected src attribute")
	}
}

func TestGenerateHTML_HeadingLevels(t *testing.T) {
	cg := services.NewLandingCodegen()
	h2 := node("heading", "h2", map[string]any{"props": map[string]any{"text": "Section", "level": float64(2)}})
	tr := tree(node("frame", "root", map[string]any{"children": []any{h2}}))
	html := cg.Generate(tr, "Heading")
	if !strings.Contains(html, "<h2") {
		t.Error("expected <h2> tag, got:", html)
	}
}

func TestGenerateHTML_Link(t *testing.T) {
	cg := services.NewLandingCodegen()
	lnk := node("link", "link")
	lnk["props"] = map[string]any{"text": "Visit", "href": "https://okuru.id"}
	tr := tree(node("frame", "root", map[string]any{"children": []any{lnk}}))
	html := cg.Generate(tr, "Link")
	if !strings.Contains(html, `href="https://okuru.id"`) {
		t.Error("expected href attribute, got:", html)
	}
}

func TestGenerateSafelist(t *testing.T) {
	cg := services.NewLandingCodegen()
	n1 := node("frame", "root")
	n1["classes"] = []any{"flex", "gap-4", "p-4"}
	n2 := node("text", "t")
	n2["classes"] = []any{"text-lg", "font-bold"}
	tr := tree(n1)
	tr["root"].(map[string]any)["children"] = []any{n2}

	safelist := cg.GenerateSafelist(tr)
	expected := []string{"flex", "font-bold", "gap-4", "p-4", "text-lg"}
	if len(safelist) != len(expected) {
		t.Fatalf("expected %d classes, got %d: %v", len(expected), len(safelist), safelist)
	}
	for i, s := range expected {
		if safelist[i] != s {
			t.Errorf("position %d: expected %q, got %q", i, s, safelist[i])
		}
	}
}

func TestGenerateHTML_WhitespaceIndent(t *testing.T) {
	cg := services.NewLandingCodegen()
	tr := tree(node("frame", "root", map[string]any{
		"children": []any{
			node("text", "a", map[string]any{"props": map[string]any{"text": "Hi"}}),
		},
	}))
	html := cg.GenerateFragment(tr)
	// Should have indentation:  <div>\n    <span>Hi</span>\n  </div>
	if !strings.Contains(html, "  ") {
		t.Error("expected indentation in output, got:\n", html)
	}
	if !strings.Contains(html, "\n") {
		t.Error("expected newlines in output, got:\n", html)
	}
}
