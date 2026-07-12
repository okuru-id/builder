# Phase 3: Frontend Builder — Canvas Editor

**Status:** ⬜ todo
**Goal:** Route `/builder/:id`. Recursive `NodeRenderer`. Selection + inline edit. Autosave debounced.

## Checklist

- [ ] TypeScript types `frontend/src/types/page-builder.ts`
- [ ] Router: tambah `/builder/:id` (top-level, no sidebar, fullscreen)
- [ ] `Builder.vue` layout: toolbar top, canvas center, panels kanan
- [ ] `NodeRenderer.vue` recursive (shared canvas + codegen preview)
- [ ] `Canvas.vue` viewport (zoom/pan, responsive breakpoints)
- [ ] `Toolbar.vue`: publish button, undo/redo, breakpoint switch
- [ ] `InspectorPanel.vue`: props editor (text, src, href, alt, classes)
- [ ] `NodeTreePanel.vue`: outline tree sidebar
- [ ] Selection state (click node → highlight border)
- [ ] Inline text edit (`contenteditable` dblclick → edit → blur save)
- [ ] Autosave debounced 2s (watch tree → PUT)
- [ ] Add node palette (frame, text, image, button)
- [ ] `bun run build` OK

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

## Files

- Create: `frontend/src/types/page-builder.ts`
- Create: `frontend/src/views/Builder.vue`
- Create: `frontend/src/components/builder/{NodeRenderer,Toolbar,InspectorPanel,NodeTreePanel,Canvas}.vue`
- Modify: `frontend/src/router/index.ts`

## Commit

```bash
git add frontend/src/types/ frontend/src/views/Builder.vue frontend/src/components/builder/
git add frontend/src/router/index.ts
git commit -m "Phase 3: canvas builder frontend — NodeRenderer, selection, inline edit, autosave"
```
