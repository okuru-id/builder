# Phase 7: One-Way Codegen (Tree → HTML/Vue Export)

**Status:** ✅ done
**Goal:** Generate clean downloadable HTML+Tailwind (atau Vue SFC) dari tree. Deterministic, idempotent. Tests.

## Checklist

- [x] 7.1 Go codegen `services/codegen.go` refine — deterministic, idempotent, indent 2sp, class dedup
- [x] 7.2 TypeScript codegen mirror (`frontend/src/composables/useCodegen.ts`)
- [x] 7.3 Download button — ExportDialog (HTML preview, copy, download, safelist export)
- [x] 7.4 Codegen tests: 12 tests (empty, deterministic×2, nested, no-shadcn, classes, dedup, self-closing, heading, link, safelist, whitespace)
- [x] 7.5 Safelist generator — `GenerateSafelist()` (Go) + `collectClasses()` (TS)
- [x] 7.6 Output pure Tailwind utilities — 0 shadcn imports (verified by test)

## Tests (critical)

```go
func TestGenerateHTML_EmptyTree(t *testing.T) { /* ... */ }
func TestGenerateHTML_Deterministic(t *testing.T) {
    html1 := GenerateHTML(tree)
    html2 := GenerateHTML(tree)
    assert.Equal(t, html1, html2)
}
func TestGenerateHTML_NestedFrame(t *testing.T) { /* ... */ }
func TestGenerateHTML_NoShadcnImports(t *testing.T) {
    // assert output contains no @/components/ui/*
}
```

## Deterministic Rules

- Atribut urut: `class` → `id` → `src` → `alt` → `href` → `type` → `value`
- Children urut sesuai tree order (no sort)
- Whitespace: indent 2 spaces, newline antar sibling
- Self-closing: `<img />`, `<input />`, `<br />`

## Files

- Modify: `backend/app/services/codegen.go`
- Create: `backend/tests/services/codegen_test.go`
- Create: `frontend/src/composables/useCodegen.ts`
- Create: `frontend/src/components/builder/ExportDialog.vue`

## Commit

```bash
git commit -m "Phase 7: one-way codegen tree→HTML/Vue + tests"
```
