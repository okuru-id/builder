// Builder store: one instance per Builder.vue mount. Holds page + tree + selection,
// exposes mutation helpers, autosave, and publish. No Pinia (ponytail: avoid new dep).
import { computed, ref, shallowRef } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import type { Breakpoint, Node, NodeType, Page, TreeShape } from '@/types/page-builder'

export type { Breakpoint }
import {
  addChild,
  cloneTree,
  deleteNode,
  emptyRoot,
  findNode,
  insertChild,
  makeNode,
  moveSibling,
  reparent as reparentTree,
  replaceNode,
  updateNode,
} from '@/components/builder/tree-utils'
import { useComponents } from '@/composables/useComponents'

// Breakpoint now lives in @/types/page-builder (shared with Node.hiddenOn).

const BP_WIDTH: Record<Breakpoint, number | null> = {
  desktop: 1920, // 1080p laptop
  tablet: 768,
  mobile: 390,
}

export function useBuilderStore() {
  const page = shallowRef<Page | null>(null)
  const tree = ref<TreeShape>({ root: emptyRoot() })
  const selectedId = ref<string | null>(null)
  const loading = ref(true)
  const saving = ref(false)
  const dirty = ref(false)
  const breakpoint = ref<Breakpoint>('desktop')

  // Drag-and-drop state. draggingId = node being dragged; dropTarget computed in dragOver.
  const draggingId = ref<string | null>(null)
  const dropTarget = ref<{ parentId: string; index: number; pos: 'before' | 'after' | 'inside' } | null>(null)

  // Component masters — loaded once on builder mount, shared with palette + renderer.
  const components = useComponents()

  const selectedNode = computed<Node | null>(() => {
    if (!selectedId.value) return null
    const found = findNode(tree.value.root, selectedId.value)
    return found?.node ?? null
  })

  const canvasWidth = computed(() => BP_WIDTH[breakpoint.value])

  // --- load ---

  async function load(id: string | number) {
    loading.value = true
    try {
      const [res] = await Promise.all([
        api.get<{ data: Page }>(`/landing-pages/${id}`),
        components.load(),
      ])
      page.value = res.data.data
      tree.value = res.data.data.tree ?? { root: emptyRoot() }
      selectedId.value = null
      dirty.value = false
      resetHistory()
    } catch (e) {
      toast.error('Failed to load page')
      console.error(e)
    } finally {
      loading.value = false
    }
  }

  // --- autosave ---

  async function saveNow() {
    if (!page.value) return
    saving.value = true
    try {
      const res = await api.put<{ data: Page }>(`/landing-pages/${page.value.id}`, {
        tree: tree.value,
      })
      page.value = res.data.data
      dirty.value = false
    } catch (e) {
      toast.error('Failed to save')
      console.error(e)
    } finally {
      saving.value = false
    }
  }

  // Debounced wrapper around saveNow for autosave on tree changes.
  const persist = useDebounceFn(saveNow, 1500)

  function notifyChange() {
    dirty.value = true
    persist()
  }

  // --- undo / redo (client-side history) ---
  // ponytail: snapshot whole tree root per edit group. Trees are small JSON;
  // structuredClone can't handle Vue proxies, so cloneTree (JSON round-trip).
  const undoStack = ref<Node[]>([])
  const redoStack = ref<Node[]>([])
  let lastSnapshotAt = 0
  const COALESCE_MS = 500 // collapse rapid edits (slider drags) into one undo
  const MAX_HISTORY = 100

  // Capture PRE-mutation state. Called at the start of every tree mutation.
  // Within COALESCE_MS of the last snapshot we skip pushing so a burst of
  // changes counts as a single undo step.
  function pushHistory() {
    const now = Date.now()
    if (now - lastSnapshotAt < COALESCE_MS) {
      lastSnapshotAt = now
      return
    }
    undoStack.value.push(cloneTree(tree.value.root))
    if (undoStack.value.length > MAX_HISTORY) undoStack.value.shift()
    redoStack.value = []
    lastSnapshotAt = now
  }

  const canUndo = computed(() => undoStack.value.length > 0)
  const canRedo = computed(() => redoStack.value.length > 0)

  function undo() {
    if (!undoStack.value.length) return
    redoStack.value.push(cloneTree(tree.value.root))
    const prev = undoStack.value.pop()!
    tree.value = { root: prev }
    selectedId.value = null
    lastSnapshotAt = 0 // start a fresh edit group after navigation
    notifyChange()
  }
  function redo() {
    if (!redoStack.value.length) return
    undoStack.value.push(cloneTree(tree.value.root))
    const next = redoStack.value.pop()!
    tree.value = { root: next }
    selectedId.value = null
    lastSnapshotAt = 0
    notifyChange()
  }
  function resetHistory() {
    undoStack.value = []
    redoStack.value = []
    lastSnapshotAt = 0
  }

  // --- mutations ---

  function select(id: string | null) {
    selectedId.value = id
  }

  function patchNode(id: string, patch: Partial<Node>) {
    pushHistory()
    tree.value = { root: updateNode(tree.value.root, id, patch) }
    notifyChange()
  }

  // Rename the page (separate from tree autosave). One-shot PUT.
  async function rename(name: string) {
    if (!page.value || !name.trim()) return
    saving.value = true
    try {
      const res = await api.put<{ data: Page }>(`/landing-pages/${page.value.id}`, { name })
      page.value = res.data.data
    } catch (e) {
      toast.error('Failed to rename')
      console.error(e)
    } finally {
      saving.value = false
    }
  }

  function addNode(type: NodeType, parentId: string | null = null) {
    pushHistory()
    // Leaf nodes (text, image, etc.) cannot be parents — add as sibling instead.
    const LEAF_TYPES = new Set<NodeType>(['text', 'heading', 'image', 'divider', 'icon', 'button', 'input'])
    let resolvedParent = parentId ?? selectedId.value ?? tree.value.root.id
    // If the resolved parent is a leaf, walk up to its own parent.
    const candidate = findNode(tree.value.root, resolvedParent)
    if (candidate && LEAF_TYPES.has(candidate.node.type) && candidate.parent) {
      resolvedParent = candidate.parent.id
    }
    const node = makeNode(type)
    tree.value = { root: addChild(tree.value.root, resolvedParent, node) }
    selectedId.value = node.id
    notifyChange()
  }

  // Append a pre-built node (e.g. from the AI agent) under a parent. History-aware.
  function appendNode(node: Node, parentId: string | null = null) {
    pushHistory()
    const parent = parentId ?? selectedId.value ?? tree.value.root.id
    tree.value = { root: addChild(tree.value.root, parent, node) }
    selectedId.value = node.id
    notifyChange()
  }

  function removeNode(id: string) {
    if (id === tree.value.root.id) {
      toast.error('Root cannot be deleted')
      return
    }
    pushHistory()
    tree.value = { root: deleteNode(tree.value.root, id) }
    if (selectedId.value === id) selectedId.value = null
    notifyChange()
  }

  function duplicateNode(id: string) {
    const found = findNode(tree.value.root, id)
    if (!found || !found.parent) return
    pushHistory()
    const copy = cloneTree(found.node)
    // re-id copy + descendants to avoid collisions
    reId(copy)
    tree.value = {
      root: insertChild(tree.value.root, found.parent.id, copy, found.index + 1),
    }
    selectedId.value = copy.id
    notifyChange()
  }

  // Move a node one slot up/down within its current parent (keyboard shortcut).
  function moveSiblingNode(id: string, dir: -1 | 1) {
    pushHistory()
    tree.value = { root: moveSibling(tree.value.root, id, dir) }
    notifyChange()
  }

  // --- HTML5 drag-and-drop handlers (no dnd-kit-vue dep). Used by NodeRenderer + TreeRow. ---
  function dragStart(id: string) {
    draggingId.value = id
  }
  function dragEnd() {
    draggingId.value = null
    dropTarget.value = null
  }
  // Compute drop zone from pointer Y relative to target bounding box.
  // pos: 'before' (top third), 'after' (bottom third), 'inside' (middle, containers only).
  function dragOver(targetId: string, e: DragEvent, isContainer: boolean) {
    if (!draggingId.value || draggingId.value === targetId) return
    e.preventDefault() // allow drop
    const el = e.currentTarget as HTMLElement
    const rect = el.getBoundingClientRect()
    const y = e.clientY - rect.top
    const h = rect.height
    let pos: 'before' | 'after' | 'inside'
    let parentId: string
    let index: number
    const found = findNode(tree.value.root, targetId)
    const parent = found?.parent
    if (isContainer && y > h * 0.33 && y < h * 0.66) {
      pos = 'inside'
      parentId = targetId
      index = -1
    } else if (y < h / 2) {
      pos = 'before'
      parentId = parent?.id ?? tree.value.root.id
      index = parent ? found!.index : 0
    } else {
      pos = 'after'
      parentId = parent?.id ?? tree.value.root.id
      index = parent ? found!.index + 1 : 0
    }
    dropTarget.value = { parentId, index, pos }
    if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
  }
  function drop() {
    const t = dropTarget.value
    const d = draggingId.value
    if (!t || !d) return
    pushHistory()
    tree.value = { root: reparentTree(tree.value.root, d, t.parentId, t.index) }
    notifyChange()
    draggingId.value = null
    dropTarget.value = null
  }

  // Drop zone at the end of a container's children: append dragged node as last child.
  // Fixes "can't move top node to bottom" when the cursor lands in empty space below
  // the last sibling (no row there to fire dragover/drop).
  function dragOverEnd(parentId: string, e: DragEvent) {
    if (!draggingId.value) return
    e.preventDefault()
    dropTarget.value = { parentId, index: -1, pos: 'after' }
    if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
  }

  // --- components / instances ---

  // Save a node (and its subtree) as a new reusable component master.
  async function createComponentFromNode(nodeId: string, name: string) {
    const found = findNode(tree.value.root, nodeId)
    if (!found) return
    const master = cloneTree(found.node)
    reId(master)
    const c = await components.create(name, { root: master })
    toast.success(`Component “${c.name}” saved`)
    return c
  }

  // Insert an instance of a component master into the tree.
  // ponytail: instance = marker node { type:'component', componentId }. Renderer
  // resolves the master tree at render time. No per-instance overrides yet —
  // add instanceOverrides merge when real divergence is needed.
  function insertInstance(componentId: number, parentId: string | null = null) {
    pushHistory()
    const parent = parentId ?? selectedId.value ?? tree.value.root.id
    const node: Node = {
      id: makeNode('component').id,
      type: 'component',
      name: components.components.value.find((c) => c.id === componentId)?.name ?? 'Component',
      props: {},
      classes: [],
      children: [],
      componentId,
    }
    tree.value = { root: addChild(tree.value.root, parent, node) }
    selectedId.value = node.id
    notifyChange()
  }

  // Break link: resolve master, deep-clone into the instance node, drop componentId.
  function breakInstance(nodeId: string) {
    const found = findNode(tree.value.root, nodeId)
    if (!found || !found.node.componentId) return
    pushHistory()
    const master = components.masterRoot(found.node.componentId)
    // ponytail: master may be null if deleted — still detach as empty node.
    if (!master) {
      tree.value = {
        root: replaceNode(tree.value.root, nodeId, () => ({
          id: nodeId,
          type: 'frame',
          name: 'Detached (master deleted)',
          props: {},
          classes: ['min-h-20', 'flex', 'items-center', 'justify-center', 'border-2', 'border-dashed', 'border-red-300', 'text-xs', 'text-red-400', 'rounded-lg'],
          children: [{
            id: nodeId + '_msg',
            type: 'text',
            name: 'Message',
            props: { text: 'Component master was deleted. Replace or remove this block.' },
            classes: ['text-xs', 'text-red-400'],
            children: [],
          }],
        })),
      }
      notifyChange()
      return
    }
    const copy = cloneTree(master)
    reId(copy)
    // Preserve the instance's own id + position; absorb master's internals.
    copy.id = found.node.id
    tree.value = { root: replaceNode(tree.value.root, nodeId, () => copy) }
    notifyChange()
  }

  // --- publish ---

  async function publish() {
    if (!page.value) return
    saving.value = true
    try {
      // Force-save latest tree before publish so codegen reflects current canvas.
      await saveNow()
      const res = await api.post<{ data: Page }>(`/landing-pages/${page.value.id}/publish`)
      page.value = res.data.data
      toast.success('Page published')
    } catch (e) {
      toast.error('Failed to publish')
      console.error(e)
    } finally {
      saving.value = false
    }
  }

  return {
    // state
    page,
    tree,
    selectedId,
    selectedNode,
    loading,
    saving,
    dirty,
    breakpoint,
    canvasWidth,
    draggingId,
    dropTarget,
    components,
    // actions
    load,
    save: saveNow,
    select,
    patchNode,
    rename,
    addNode,
    appendNode,
    removeNode,
    duplicateNode,
    moveSiblingNode,
    undo,
    redo,
    canUndo,
    canRedo,
    resetHistory,
    createComponentFromNode,
    insertInstance,
    breakInstance,
    dragStart,
    dragEnd,
    dragOver,
    dragOverEnd,
    drop,
    publish,
  }
}

export type BuilderStore = ReturnType<typeof useBuilderStore>

function reId(n: Node) {
  n.id = makeNode(n.type).id
  if (Array.isArray(n.children)) n.children.forEach(reId)
}
