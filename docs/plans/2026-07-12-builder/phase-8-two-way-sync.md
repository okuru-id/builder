# Phase 8: Two-Way Sync (Code → Tree)

**Status:** ✅ done
**Goal:** Parse HTML/Vue kembali ke tree JSON. Developer edit code, canvas reflect. Round-trip stable.

## Checklist

- [x] 8.1 Parse via `@vue/compiler-dom` (covers HTML + Vue template without SFC extraction; gunakan `parseVueTemplate()` yang strip script/style dulu)
- [x] 8.2 Convert AST elements (div, section, h1-h6, img, a, button, span, p) → builder nodes; Vue directives (v-if, v-for) marked dengan placeholder (ponytail: belum di-handle)
- [x] 8.3 Reverse class mapping (`classReverseMap.ts`) — prefix + color + arbitrary value → inspector section/property
- [x] 8.4 Round-trip test: 7 round-trip tests (frame, text, heading, button, link, image, section+children) → generate → parse → generate = identical
- [x] 8.5 Import dialog — paste HTML atau upload .html/.vue → parse → create page via API → buka builder
- [x] 8.6 ~Diff viewer~ — skip (ponytail: lower priority, tambah setelah two-way sync benar-benar dipakai)

## Catatan

- `@vue/compiler-dom` enough untuk template parsing. `compiler-sfc` di-skip (ponytail)
- `parseVueTemplate()` strip script/style via regex sebelum parse
- Parser skip directive nodes (v-if/v-for) dengan placeholder frame — round-trip sempurna hanya untuk static HTML
- `classReverseMap` partial — tambah prefix seiring kebutuhan

## Round-Trip Test

```typescript
const original = loadTestTree('hero-section.json')
const html = generateHTML(original)
const reparsed = parseHTML(html)
const html2 = generateHTML(reparsed)
assert(html === html2)
assert(deepEqual(normalizeTree(original), normalizeTree(reparsed)))
```

## Reverse Class Mapping

```
'gap-4'         → { gap: '4' }
'flex-col'      → { direction: 'column' }
'text-2xl'      → { fontSize: '2xl' }
'bg-blue-500'   → { background: 'blue-500' }
```

Build lookup table dari Tailwind class registry.

## Files

- Create: `frontend/src/services/parser.ts`
- Create: `frontend/src/services/classReverseMap.ts`
- Create: `frontend/src/components/builder/ImportDialog.vue`
- Create: `frontend/tests/parser.test.ts`

## Catatan

- Lower priority. One-way (canvas → code) cukup untuk MVP. Sync dua arah kompleks (v-for, v-if, slot, dynamic binding tidak round-trip sempurna).
- Mark simplified parts `ponytail:` — ceiling + upgrade path.

## Commit

```bash
git commit -m "Phase 8: two-way sync code→tree"
```
