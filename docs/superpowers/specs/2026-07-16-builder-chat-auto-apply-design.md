# Builder Chat Auto-Apply — Design

**Date:** 2026-07-16

## Goal

Allow an operator to approve an AI-builder session once, then automatically
apply every valid `action:*` block in later assistant replies until they turn
the mode off or refresh the page.

## UX

- Add an `Auto-apply` toggle in the chat composer.
- Enabling it opens one confirmation dialog warning that all future actions,
  including `delete` and `move`, apply automatically.
- While active, show a prominent destructive-style banner: "Auto-apply active
  — changes including delete are applied automatically." A Stop button disables
  the mode.
- Refresh always disables the mode. It is deliberately not persisted.

## Flow

1. User enables toggle and confirms once.
2. A normal chat request streams its answer.
3. When streaming ends, parse action blocks in source order and call the
   existing `applyAction` for each.
4. Existing store mutations keep their own undo snapshots and autosave path.
5. Invalid payloads, unknown IDs, or failed operations are skipped with a toast;
   later actions still run.
6. Disabling auto-apply affects only future replies; already-completed actions
   stay applied and remain undoable.

## Implementation

Only `frontend/src/components/builder/ChatPanel.vue` changes.

- Add runtime `ref(false)` for auto-apply state.
- Use existing shadcn-vue `ConfirmDialog` integration for the one-time warning.
- Add `applyAll(msg)` helper. It parses `partsOf(msg)`, filters action parts,
  then applies them sequentially after the SSE stream completes.
- Keep the existing manual Apply buttons. They are disabled/marked Applied after
  auto-apply through the existing `appliedFlags` map.
- Do not alter backend protocol, storage schema, or action behavior.

## Safety

- Explicit one-time confirmation states that deletion is included.
- Visual active-state warning and Stop control remain visible.
- Session-only state limits accidental persistence.
- Existing undo/redo can recover every applied action.

## Validation

- Enable auto-apply, ask for add/classes/text/delete/move actions; verify each
  is applied after reply end and can be undone.
- Disable or refresh; verify later actions show manual Apply buttons.
- Feed malformed action JSON; verify it is skipped without blocking valid later
  actions.
