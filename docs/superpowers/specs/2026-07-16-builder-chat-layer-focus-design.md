# Builder Chat Layer Focus — Design

**Date:** 2026-07-16

## Goal

Let an operator point an AI-builder request at one layer-tree node, giving the
agent its real identity and structure before the operator writes a request.

## UX

Each `TreeRow` gets a compact robot button visible on hover. Clicking it does
not select, drag, or mutate the node. It opens/focuses the Agent tab and fills
the composer with:

```
Focus on node "<name>" (ID: "<id>").
```

The operator adds their instruction and presses Send.

## Data Flow

1. `TreeRow` emits `ask-agent(node)`.
2. `NodeTreePanel` forwards the event.
3. The builder view/store holds one runtime focus request.
4. `ChatPanel` consumes it, switches/focuses its composer, and sends the full
   focus-node snapshot with the normal tree and node catalog.
5. Backend includes the snapshot in the system prompt. The agent targets that
   node or descendants unless the operator explicitly requests broader changes.

The full tree remains in the request so IDs and surrounding structure remain
available.

## Files

- `TreeRow.vue`: hover robot button + event.
- `NodeTreePanel.vue`: event forwarding.
- Builder shared state/injection: one runtime focus-node ref.
- `RightPanel.vue` / `ChatPanel.vue`: switch Agent tab, prefill and send focus.
- `builder_chat_controller.go`: request field and focused-node prompt section.

## Safety

- Focus is advisory; user text can explicitly broaden scope.
- The full, real node snapshot is sent; the agent must not invent IDs.
- No automatic mutation arises from pointing alone.
- Existing Apply and Auto-apply controls retain their behavior.

## Scope

Single-node focus only. Multi-select focus is out of scope.
