# Builder Chat Agent Fixes — Design

**Date:** 2026-07-16
**Scope:** Fix 5 issues in the builder AI agent (frontend `ChatPanel.vue`, backend `builder_chat_controller.go`, plus one store helper).

## Context

The builder AI agent streams OpenAI-compatible completions over SSE and emits
fenced `action:*` blocks that render as "Apply" buttons. An audit found five
issues ranging from UX feedback to data loss.

## Issues

1. **Typing indicator never renders.** `partsOf('')` returns an empty array, so
   the `v-for` has zero iterations and the in-loop typing indicator is never
   shown. Users get no feedback while waiting for the first token.
2. **`action:classes` is destructive.** The backend prompt says "replace matched
   prefixes" (merge), but the frontend calls `patchNode(id, { classes: set })`,
   which replaces the entire class array. Sending `["flex"]` wipes out color,
   padding, typography classes. Data loss.
3. **Backend has no HTTP timeout.** `nethttp.DefaultClient` blocks forever if the
   LLM stalls — goroutine/connection leak under load.
4. **Stream cannot be cancelled.** No stop button; long generations block the UI
   until [DONE].
5. **Actions limited to add/classes/text.** No `delete` or `move/reorder`.

## Design

**#1 Typing indicator.** Move the indicator out of the `v-for` loop into a
sibling element rendered when `busy && lastAssistantMsg.content === ''`. Visible
during the wait for the first token, then replaced by streamed parts.

**#2 Classes merge.** Use the existing `cn()` helper (`src/lib/utils.ts` =
twMerge + clsx, already installed) to merge the agent's classes onto the node's
current classes:

```ts
const cur = findNode(tree.value.root, id)?.classes ?? []
store.patchNode(id, { classes: cn(cur, ...set).split(/\s+/).filter(Boolean) })
```

twMerge strips conflicting prefixes (e.g. `flex-row` overridden by `flex-col`)
and preserves the rest. The backend prompt already matches this semantics — no
prompt change.

**#3 Backend timeout.** Replace `nethttp.DefaultClient` with a client scoped to
the controller: `&nethttp.Client{Timeout: 60 * time.Second}`. Request context
already flows from the gin request via `Origin().Context()`, so client
disconnects cancel the upstream call.

**#4 Cancel stream.** Add an `AbortController` per `send()`. Show a "Stop"
button while `busy` (replaces Send). On abort, fetch rejects; keep the
partially-streamed assistant content. No backend change — the gin request
context cancels the upstream request.

**#5 `delete` + `move` actions.** New fenced formats:

```
action:delete
{ "nodeId": "<id>" }

action:move
{ "nodeId": "<id>", "parentId": "root", "index": -1 }
```

Frontend handling:
- `delete` → existing `store.removeNode(id)` (already history-aware).
- `move` → new `store.moveNode(id, parentId, index)` helper (~6 lines):
  `pushHistory()` then `tree.value = { root: reparentTree(tree.value.root, {id},
  parentId, index) }`, reusing existing `reparentTree` from `tree-utils.ts`.

Add the two formats to `builderSystemPrompt` and extend the action parser regex
to `add|classes|text|delete|move` and the `ActionPart['kind']` union.

## Components Touched

| File | Change |
|---|---|
| `frontend/src/components/builder/ChatPanel.vue` | typing indicator, cancel button + AbortController, merge classes via `cn()`, parse + apply `delete`/`move` |
| `frontend/src/components/builder/useBuilderStore.ts` | add `moveNode` helper, export it |
| `backend/app/http/controllers/admin/builder_chat_controller.go` | HTTP client timeout, extend system prompt with `delete`/`move` formats |

## Data Flow

```
User types → send() → fetch /admin/api/builder/chat (AbortController)
           → backend proxy → upstream LLM (60s client timeout)
           ← SSE tokens stream → assistantMsg.content appends
           ← [DONE] → busy=false
Cancel → abort() → fetch rejects → keep partial content
Apply action → merge/delete/move via store (pushHistory) → autosave
```

## Error Handling

- Upstream HTTP error → existing path emits `LLM error (HTTP %d): <body>` event.
- Timeout → context deadline emits guidance event + [DONE].
- Abort → frontend swallows `AbortError`, keeps partial content, no toast.
- Invalid action JSON → `payload = null`, Apply button disabled with toast.
- `move`/`delete` to unknown id → store no-op (findNode returns null),
  toast "Node not found".

## Testing

- No new test files. `backend/app/services/codegen_test.go` exists; controller
  is a thin proxy and not unit-tested today — out of scope.
- Manual checks: (a) typing indicator shows before first token; (b) `classes`
  action preserves unrelated classes; (c) Stop button cancels mid-stream;
  (d) `delete` and `move` apply and are undoable; (e) backend recovers from a
  stalled upstream within 60s.

## Out of Scope

- Rate limiting / abuse protection (add when multi-tenant).
- Streaming progress bar.
- Tool/function-calling migration (still fenced-text protocol).
