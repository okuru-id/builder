# Phase 5: Component / Instance System

**Status:** ✅ done
**Goal:** Reusable component master. Instance pasang dari palette. Edit master → update semua instance.

## Checklist

- [x] 5.1 CRUD `/admin/api/landing-components` (backend: index, store, show, update, destroy)
- [x] 5.2 "Create Component" action — save selected node(s) sebagai master
- [x] 5.3 Component palette (left drawer) — list komponen tersedia
- [x] 5.4 Instance rendering — NodeRenderer resolve `componentId` → load master tree (instance: pure reference, live resolve; overrides ditangguhkan — YAGNI)
- [x] 5.5 ~Edit master → flag instance dirty~ — N/A: arsitektur live-resolve berarti instance selalu mirror master, tidak ada divergence yang perlu di-flag
- [x] 5.6 Break link — convert instance ke standalone copy (deep clone, hapus `componentId`)
- [x] 5.7 Component registry sync (pull saat builder load)

## Catatan arsitektur

- Instance = node `{type:'component', componentId}`. NodeRenderer resolve master di render-time (client), `services.ResolveComponentInstances` resolve di publish-time (server). HTML publish self-contained, tidak depend master DB.
- `instanceOverrides` skip sengaja (ponytail). Saat butuh per-instance divergence, tambah deepMerge master+overrides di `NodeRenderer.displayNode` + `services.resolveNode`. Saat itu baru 5.5 badge "dirty" relevan.

## Files

- Create: `backend/app/http/controllers/admin/landing_component_controller.go`
- Modify: `backend/routes/api.go`
- Create: `frontend/src/components/builder/ComponentPalette.vue`
- Create: `frontend/src/composables/useComponents.ts`
- Modify: `frontend/src/components/builder/NodeRenderer.vue` (resolve componentId)

## Instance Merge

```typescript
function resolveInstance(node: Node, masters: Map<number, Node>): Node {
  if (!node.componentId) return node
  const master = masters.get(node.componentId)
  if (!master) return node
  return deepMerge(structuredClone(master), node.instanceOverrides ?? {})
}
```

## Commit

```bash
git commit -m "Phase 5: component/instance system"
```
