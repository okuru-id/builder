# Phase 8: Two-Way Sync (Code → Tree)

**Status:** ⬜ todo
**Goal:** Parse `.vue`/HTML kembali ke tree JSON. Developer edit code, canvas reflect. Round-trip stable.

## Checklist

- [ ] 8.1 Parse Vue SFC via `@vue/compiler-sfc` → extract template AST
- [ ] 8.2 Convert AST nodes (v-for, v-if, v-bind) → builder node tree
- [ ] 8.3 Parse Tailwind classes kembali ke inspector props (reverse mapping)
- [ ] 8.4 Round-trip test: generate → parse → generate = identical
- [ ] 8.5 Import dialog — upload `.vue`/`.html` → create page dari parse
- [ ] 8.6 Diff viewer saat parse berbeda dari current tree

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
