# Phase 4: Advanced Canvas — Drag, Resize, Snap, Multi-select

**Status:** ⬜ todo
**Goal:** Direct manipulation Figma-class: drag move/reorder, resize handles, smart snapping, multi-select, keyboard shortcuts.

## Checklist

- [ ] 4.1 dnd-kit-vue drag reorder child antar sibling / antar parent
- [ ] 4.2 Selection overlay — highlight border + resize handles
- [ ] 4.3 Multi-selection (Shift+click, lasso drag)
- [ ] 4.4 Resize handles (frame corners, edges right/bottom)
- [ ] 4.5 Smart snapping — alignment guides saat drag (top/left/center/right/bottom nearest sibling)
- [ ] 4.6 Arrow key nudge (1px, 10px Shift)
- [ ] 4.7 Cmd+D duplicate selected
- [ ] 4.8 Backspace/Delete hapus selected
- [ ] 4.9 Quick add node (palette drag to canvas / "+" edge)
- [ ] Keyboard shortcut help overlay (?)

## Files

- Modify: `frontend/src/components/builder/Canvas.vue` (overlay layer)
- Create: `frontend/src/components/builder/SelectionOverlay.vue`
- Create: `frontend/src/components/builder/SnapGuides.vue`
- Create: `frontend/src/composables/useSelection.ts`
- Create: `frontend/src/composables/useKeyboard.ts`

## Catatan

- Canvas pakai **pointer events custom** untuk free positioning, bukan dnd-kit list reorder (dnd-kit untuk reorder tree sibling).
- Snapping = nearest neighbor alignment, threshold 4px.

## Commit

```bash
git commit -m "Phase 4: drag, resize, snap, multi-select, keyboard shortcuts"
```
