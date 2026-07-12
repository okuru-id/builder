# Phase 5: Component / Instance System

**Status:** ⬜ todo
**Goal:** Reusable component master. Instance pasang dari palette. Edit master → update semua instance.

## Checklist

- [ ] 5.1 CRUD `/admin/api/landing-components` (backend: index, store, show, update, destroy)
- [ ] 5.2 "Create Component" action — save selected node(s) sebagai master
- [ ] 5.3 Component palette (left drawer) — list komponen tersedia
- [ ] 5.4 Instance rendering — NodeRenderer resolve `componentId` → load master tree + merge `instanceOverrides`
- [ ] 5.5 Edit master → flag instance dirty, badge "Update from master"
- [ ] 5.6 Break link — convert instance ke standalone copy (deep clone, hapus `componentId`)
- [ ] 5.7 Component registry sync (pull saat builder load)

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
