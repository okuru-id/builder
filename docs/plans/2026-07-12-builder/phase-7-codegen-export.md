# Phase 7: One-Way Codegen (Tree → HTML/Vue Export)

**Status:** ⬜ todo
**Goal:** Generate clean downloadable HTML+Tailwind (atau Vue SFC) dari tree. Deterministic, idempotent. Tests.

## Checklist

- [ ] 7.1 Go codegen `services/codegen.go` refine — deterministic, idempotent
- [ ] 7.2 TypeScript codegen mirror (live preview export)
- [ ] 7.3 Download button — export HTML atau Vue SFC
- [ ] 7.4 Codegen tests: input tree X → expected output (snapshot/hash)
- [ ] 7.5 Safelist generator — scan tree classes → `safelist.txt` untuk Tailwind JIT build
- [ ] 7.6 Output **pure Tailwind utilities**, 0 import `@/components/ui/*`

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
