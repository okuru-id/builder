# Phase 3: Frontend Builder — Canvas Editor

**Status:** ✅ done
**Goal:** Route `/builder/:id`. Recursive `NodeRenderer`. Selection + inline edit. Autosave debounced.

## Checklist

- [x] TypeScript types `frontend/src/types/page-builder.ts`
- [x] Router: tambah `/builder/:id` (top-level, fullscreen, no sidebar)
- [x] `Builder.vue` layout: toolbar top, canvas center, panels kiri (tree) + kanan (inspector)
- [x] `NodeRenderer.vue` recursive (shared canvas + codegen preview via `readonly` prop)
- [x] `Canvas.vue` viewport (breakpoint width: desktop/tablet 768/mobile 390)
- [x] `Toolbar.vue`: page name (editable), breakpoint switch, dirty/saving indicator, publish button
- [x] `InspectorPanel.vue`: props editor (text, level, src, alt, href, classes textarea, delete/duplicate)
- [x] `NodeTreePanel.vue` + `TreeRow.vue`: outline tree sidebar + add-node palette
- [x] Selection state (click node → outline ring; click canvas bg → deselect)
- [x] Inline text edit (`contenteditable` dblclick → edit → blur/enter/esc save)
- [x] Autosave debounced 1.5s (watch tree → PUT), `saveNow()` for publish flush
- [x] Add node palette (frame, section, text, heading, image, button, link)
- [x] Keyboard: Delete/Backspace hapus, Cmd/Ctrl+D duplikat (skip saat typing)
- [x] Landing pages list view (`LandingPages.vue`) + sidebar entry
- [x] `bun run build` OK

## NodeRenderer.vue (core)

```vue
<script setup lang="ts">
import type { Node } from '@/types/page-builder'
defineProps<{ node: Node; depth?: number }>()
defineEmits<{ select: [id: string] }>()
</script>

<template>
  <component :is="el(node.type)"
    :class="node.classes"
    v-bind="attrs(node)"
    @click.stop="$emit('select', node.id)">
    <template v-if="isLeaf(node)">
      {{ node.props.text }}
    </template>
    <NodeRenderer v-else v-for="c in node.children" :key="c.id"
      :node="c" :depth="(depth ?? 0)+1" @select="$emit('select', $event)" />
  </component>
</template>
```

## Autosave

```typescript
import { useDebounceFn } from '@vueuse/core'
const save = useDebounceFn(() => api.put(`/landing-pages/${id}`, { tree }), 2000)
watch(tree, save, { deep: true })
```

## Hasil Verifikasi

```
bun run build         → ✓ built in 780ms
GET /admin/builder/:id → 200 (SPA shell, 467 bytes)
Builder-*.js chunk     → ada di assets/
API login/list/show    → 200 (builder dependencies OK)
```

Runtime visual (canvas render, selection, inline edit) butuh verifikasi manual di browser.

## Catatan

- State via composable `useBuilderStore` + provide/inject (BUILDER_KEY). Ponytail: no Pinia, no new dep.
- `TreeRow.vue` SFC terpisah dengan `defineOptions({ name: 'TreeRow' })` untuk rekursi — lebih bersih dari `h()` recursive yang berisiko typing error.
- `NodeRenderer.vue` punya prop `readonly` agar bisa dipakai codegen preview / iframe snapshot tanpa interaction layer (Phase 7).
- `ponytail:` zoom/pan canvas deferred — breakpoint width switch cukup untuk MVP. Tambah zoom saat butuh (Phase 4).
- `ponytail:` undo/redo belum — butuh history stack. Tambah di Phase 4 bareng dengan drag/snap.
- `ponytail:` inspector classes = textarea mentah. Style panel UI controls nyata di Phase 6.
- `ponytail:` drag reorder belum — dnd-kit-vue Phase 4.

## Files

- Create: `frontend/src/types/page-builder.ts`
- Create: `frontend/src/components/builder/{tree-utils.ts,useBuilderStore.ts,injection.ts,NodeRenderer.vue,TreeRow.vue,Canvas.vue,Toolbar.vue,InspectorPanel.vue,NodeTreePanel.vue}`
- Create: `frontend/src/views/{Builder.vue,LandingPages.vue}`
- Modify: `frontend/src/router/index.ts` (route `/builder/:id`, `/pages`)
- Modify: `frontend/src/components/AppSidebar.vue` (entry Pages)

## Commit

```bash
git add frontend/src/types/ frontend/src/views/Builder.vue frontend/src/components/builder/
git add frontend/src/router/index.ts
git commit -m "Phase 3: canvas builder frontend — NodeRenderer, selection, inline edit, autosave"
```
