# Landing Page Builder — Phase Tracker

Figma-like canvas editor untuk okuru.id. Tree JSON source of truth, direct manipulation, codegen ke HTML+Tailwind.

## Status

| Phase | Nama | Status | Commit | File |
|-------|------|--------|--------|------|
| 0 | Cleanup legacy | ✅ done | `998200c` | [phase-0-cleanup.md](phase-0-cleanup.md) |
| 1 | Data model + API | ✅ done | `d8bc2f2` | [phase-1-data-model.md](phase-1-data-model.md) |
| 2 | Codegen + publish | ✅ done | _(this commit)_ | [phase-2-codegen.md](phase-2-codegen.md) |
| 3 | Canvas editor | ✅ done | _(this commit)_ | [phase-3-canvas-editor.md](phase-3-canvas-editor.md) |
| 4 | Drag/snap/multi-select | ✅ done (partial) | _(this commit)_ | [phase-4-drag-snap.md](phase-4-drag-snap.md) |
| 5 | Component system | ⬜ todo | — | [phase-5-components.md](phase-5-components.md) |
| 6 | Style panel + token | ⬜ todo | — | [phase-6-style-panel.md](phase-6-style-panel.md) |
| 7 | Codegen export | ⬜ todo | — | [phase-7-codegen-export.md](phase-7-codegen-export.md) |
| 8 | Two-way sync | ⬜ todo | — | [phase-8-two-way-sync.md](phase-8-two-way-sync.md) |

## Arsitektur

- **DB:** Postgres JSONB (`landing_pages.tree`) = source of truth, bukan HTML string.
- **Canvas:** Vue 3 recursive `<NodeRenderer>` render tree → DOM.
- **Style:** Tailwind utility `classes[]` per node. No inline `style:`.
- **Publish:** tree → codegen → static HTML (cache). Autosave tree saat edit.
- **Output:** pure Tailwind utilities, 0 import `@/components/ui/*`.

## Tech Stack

Go 1.24 + Goravel + GORM · PostgreSQL · Vue 3 + shadcn-vue + reka-ui + dnd-kit-vue · Tailwind v4 · Bun.

## Node Schema (source of truth)

```typescript
interface Node {
  id: string
  type: 'frame' | 'text' | 'image' | 'button' | 'link' | 'section' | 'component'
  name: string
  props: Record<string, any>  // text, src, href, alt, ...
  classes: string[]           // Tailwind utilities
  children: Node[]
  componentId?: number
  instanceOverrides?: Partial<Node>
}
```

## Anti-Pattern Checklist

- ✅ Tree JSON source of truth, bukan HTML string
- ✅ Direct manipulation canvas, bukan form-only
- ✅ Section library = optional starting point, bukan satu-satunya cara
- ✅ Token + Tailwind utilities, bukan inline style
- ✅ Publish pipeline terpisah dari canvas edit
- ✅ Output pure Tailwind, 0 import `@/components/ui/*`

## Referensi

- Master plan (legacy, merged ke file-file per phase): `docs/plans/2026-07-12-landing-page-builder.md`
- Spec asli: `docs/plans/2026-07-05-landing-page-builder-design.md`
