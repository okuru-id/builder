package services

import (
	"strings"
	"testing"
)

func TestLandingCodegen_EmptyFrame(t *testing.T) {
	tree := map[string]any{
		"root": map[string]any{
			"type":    "frame",
			"classes": []any{"min-h-screen"},
			"children": []any{},
		},
	}
	want := `<div class="min-h-screen"></div>`
	got := NewLandingCodegen().GenerateFragment(tree)
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLandingCodegen_TextAndHeading(t *testing.T) {
	tree := map[string]any{
		"root": map[string]any{
			"type": "frame",
			"children": []any{
				map[string]any{
					"type":    "heading",
					"classes": []any{"text-4xl"},
					"props":   map[string]any{"text": "Hello", "level": float64(1)},
				},
				map[string]any{
					"type":    "text",
					"classes": []any{"text-neutral-600"},
					"props":   map[string]any{"text": "World"},
				},
			},
		},
	}
	got := NewLandingCodegen().GenerateFragment(tree)
	if !strings.Contains(got, `<h1 class="text-4xl">Hello</h1>`) {
		t.Errorf("missing heading, got %q", got)
	}
	if !strings.Contains(got, `<span class="text-neutral-600">World</span>`) {
		t.Errorf("missing text, got %q", got)
	}
}

func TestLandingCodegen_ImageAttrs(t *testing.T) {
	tree := map[string]any{
		"root": map[string]any{
			"type": "frame",
			"children": []any{
				map[string]any{
					"type":    "image",
					"classes": []any{"w-full"},
					"props":   map[string]any{"src": "/a.png", "alt": "Alt"},
				},
			},
		},
	}
	got := NewLandingCodegen().GenerateFragment(tree)
	// class must come before src, src before alt (deterministic order)
	want := `<img class="w-full" src="/a.png" alt="Alt" />`
	if !strings.Contains(got, want) {
		t.Fatalf("got %q, want substring %q", got, want)
	}
}

func TestLandingCodegen_NestedFrame(t *testing.T) {
	tree := map[string]any{
		"root": map[string]any{
			"type":    "section",
			"classes": []any{"p-4"},
			"children": []any{
				map[string]any{
					"type":    "frame",
					"classes": []any{"flex"},
					"children": []any{
						map[string]any{
							"type":  "button",
							"props": map[string]any{"text": "Go"},
							"classes": []any{"bg-blue-500"},
						},
					},
				},
			},
		},
	}
	got := NewLandingCodegen().GenerateFragment(tree)
	want := `<section class="p-4"><div class="flex"><button class="bg-blue-500">Go</button></div></section>`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLandingCodegen_Deterministic(t *testing.T) {
	tree := map[string]any{
		"root": map[string]any{
			"type": "frame",
			"children": []any{
				map[string]any{"type": "text", "props": map[string]any{"text": "x"}},
			},
		},
	}
	g := NewLandingCodegen()
	a := g.Generate(tree, "T")
	b := g.Generate(tree, "T")
	if a != b {
		t.Fatal("non-deterministic output")
	}
}

func TestLandingCodegen_LinkWithChildren(t *testing.T) {
	tree := map[string]any{
		"root": map[string]any{
			"type": "link",
			"props": map[string]any{"href": "https://okuru.id"},
			"children": []any{
				map[string]any{"type": "text", "props": map[string]any{"text": "Visit"}},
			},
		},
	}
	got := NewLandingCodegen().GenerateFragment(tree)
	want := `<a href="https://okuru.id"><span>Visit</span></a>`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLandingCodegen_HTMLEscapesText(t *testing.T) {
	tree := map[string]any{
		"root": map[string]any{
			"type":  "text",
			"props": map[string]any{"text": "<script>alert(1)</script>"},
		},
	}
	got := NewLandingCodegen().GenerateFragment(tree)
	if strings.Contains(got, "<script>") {
		t.Fatalf("XSS not escaped: %q", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Fatalf("expected escaped, got %q", got)
	}
}

func TestLandingCodegen_FullDocumentHasTailwind(t *testing.T) {
	tree := map[string]any{
		"root": map[string]any{"type": "frame", "children": []any{}},
	}
	got := NewLandingCodegen().Generate(tree, "Home")
	if !strings.HasPrefix(got, "<!DOCTYPE html>") {
		t.Error("missing doctype")
	}
	if !strings.Contains(got, "cdn.tailwindcss.com") {
		t.Error("missing tailwind CDN")
	}
	if !strings.Contains(got, "<title>Home</title>") {
		t.Error("missing title")
	}
}
