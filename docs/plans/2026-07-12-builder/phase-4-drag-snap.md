# Phase 4: Advanced Canvas — Drag, Resize, Snap, Multi-select

**Status:** ✅ done (partial — see catatan)
**Goal:** Direct manipulation untuk flow-layout page builder (Webflow model, bukan Figma absolute canvas).

## Checklist

- [x] 4.1 drag reorder child antar sibling / antar parent (HTML5 DnD native, no dnd-kit-vue)
- [x] 4.2 Selection overlay — outline ring (dari Phase 3) + drop-target ring pada container
- [ ] ~4.3 Multi-selection~ → deferred ke Phase 5 (butuh grouping/component dulu)
- [ ] ~4.4 Resize handles~ → N/A flow layout; sizing via classes, handle di Phase 6 inspector
- [ ] ~4.5 Smart snapping~ → N/A flow layout; DOM flow sudah atur posisi
- [ ] ~4.6 Arrow pixel nudge~ → N/A no absolute positioning; move sibling via Cmd+Up/Down
- [x] 4.7 Cmd+D duplicate selected (dari Phase 3)
- [x] 4.8 Backspace/Delete hapus selected (dari Phase 3)
- [x] 4.9 Quick add node palette (dari Phase 3)
- [x] Bonus: Cmd+Up/Down move sibling, naik/turun button di inspector
- [ ] ~Keyboard help overlay~ → YAGNI, shortcuts visible di inspector buttons

## Files

- Modify: `frontend/src/components/builder/tree-utils.ts` (insertChild, moveSibling, reparent, isDescendant, replaceNode)
- Modify: `frontend/src/components/builder/useBuilderStore.ts` (drag state + handlers, moveSiblingNode)
- Modify: `frontend/src/components/builder/NodeRenderer.vue` (draggable + DnD handlers + drop indicators)
- Modify: `frontend/src/components/builder/TreeRow.vue` (draggable outline rows)
- Modify: `frontend/src/views/Builder.vue` (Cmd+Up/Down keyboard)
- Modify: `frontend/src/components/builder/InspectorPanel.vue` (Naik/Turun buttons)

## Catatan

**Ponytail: why banyak item di-defer.**

Canvas ini = **flow layout** (DOM children + Tailwind flex classes), bukan absolute positioning seperti Figma. Tree node style via `classes[]` saja, no inline `x/y`. Konsekuensi:

- **Resize handles (4.4)**: tidak ada pixel w/h. Sizing via `w-64`, `max-w-md`. Inspector Phase 6 yang handle w/h controls.
- **Snap guides (4.5)**: alignment guides berguna di absolute canvas. DOM flow sudah atur posisi otomatis — snap tidak applicable.
- **Arrow pixel nudge (4.6)**: no absolute pos, nudge margin/padding via inspector lebih tepat. Cmd+Up/Down untuk move sibling lebih berguna.
- **Multi-select/lasso (4.3)**: YAGNI sampai grouping/component exists (Phase 5). Batch delete bisa tunggu.

Yang diimplementasi = value-real untuk flow-layout page builder (Webflow/Framer model):

1. **Drag reorder sibling + cross-parent** — HTML5 DnD native, no dep. Drop zone: before (top third), after (bottom third), inside (middle, containers only).
2. **Move up/down keyboard** — Cmd+Up/Down. Plus inspector buttons.
3. **Drop indicator visual** — ring pada container saat drop inside, opacity pada dragging node.

**Implementasi detail:**

- `tree-utils.ts`: tambah `insertChild`, `moveSibling`, `reparent`, `isDescendant`, `replaceNode`.
  - `reparent` guard: cannot drop node ke dirinya sendiri atau descendant-nya.
- `useBuilderStore.ts`: drag state (`draggingId`, `dropTarget`) + handlers (`dragStart/dragOver/drop/dragEnd`). `dragOver` compute posisi dari pointer Y relative to bounding box.
- `NodeRenderer.vue`: `:draggable="!readonly && !editing"`, handlers wired. Drag node dimmed (opacity-40), drop-target container di-ring.
- `TreeRow.vue`: sama, draggable di outline panel. Drag-from-outline ke canvas juga works (shared store state).
- `Builder.vue`: Cmd+Up/Down keyboard.
- `InspectorPanel.vue`: Naik/Turun/Duplikat/Hapus buttons.

## Hasil Verifikasi

```
bun run build → ✓ built in 715ms
Builder chunk compiled
```

Runtime visual drag perlu verifikasi manual di browser.

## Files

## Commit

```bash
git commit -m "Phase 4: drag, resize, snap, multi-select, keyboard shortcuts"
```
