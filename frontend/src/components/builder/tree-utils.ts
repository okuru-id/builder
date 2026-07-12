// Pure tree mutation helpers. No Vue reactivity inside — callers wrap results.
import type { Node, NodeType, NodeProps, TreeShape } from '@/types/page-builder'

export function emptyRoot(): Node {
  return {
    id: 'root',
    type: 'frame',
    name: 'Page',
    props: {},
    classes: ['min-h-screen', 'bg-white', 'text-neutral-900'],
    children: [],
  }
}

let counter = 0
export function uid(prefix = 'n'): string {
  counter += 1
  return `${prefix}_${Date.now().toString(36)}_${counter}`
}

export function makeNode(type: NodeType, partial?: Partial<Node>): Node {
  const base: Node = {
    id: uid(),
    type,
    name: defaultName(type),
    props: defaultProps(type),
    classes: defaultClasses(type),
    children: [],
  }
  return { ...base, ...partial }
}

export function defaultName(type: NodeType): string {
  const map: Record<NodeType, string> = {
    frame: 'Frame',
    section: 'Section',
    text: 'Text',
    heading: 'Heading',
    image: 'Image',
    button: 'Button',
    link: 'Link',
    component: 'Component',
  }
  return map[type] ?? 'Node'
}

export function defaultProps(type: NodeType): NodeProps {
  switch (type) {
    case 'text':
      return { text: 'Teks baru' }
    case 'heading':
      return { text: 'Judul', level: 2 }
    case 'button':
      return { text: 'Tombol' }
    case 'link':
      return { href: '#', text: 'Tautan' }
    case 'image':
      return { src: '', alt: '' }
    default:
      return {}
  }
}

export function defaultClasses(type: NodeType): string[] {
  switch (type) {
    case 'frame':
      return ['flex', 'flex-col', 'gap-4']
    case 'section':
      return ['py-16', 'px-6']
    case 'text':
      return ['text-base', 'text-neutral-700']
    case 'heading':
      return ['text-3xl', 'font-bold']
    case 'button':
      return ['inline-flex', 'items-center', 'rounded-lg', 'bg-blue-600', 'px-5', 'py-2.5', 'text-white']
    case 'link':
      return ['text-blue-600', 'underline']
    case 'image':
      return ['w-full', 'object-cover']
    default:
      return []
  }
}

export interface FoundNode {
  node: Node
  parent: Node | null
  index: number
}

// DFS search; root has parent = null.
export function findNode(root: Node, id: string): FoundNode | null {
  if (root.id === id) return { node: root, parent: null, index: -1 }
  return findInChildren(root.children, id, root)
}

function findInChildren(children: Node[], id: string, parent: Node): FoundNode | null {
  for (let i = 0; i < children.length; i++) {
    const c = children[i]
    if (c.id === id) return { node: c, parent, index: i }
    const deeper = findInChildren(c.children, id, c)
    if (deeper) return deeper
  }
  return null
}

// Returns a new tree with the patch applied to the matched node. Immutable.
export function updateNode(root: Node, id: string, patch: Partial<Node>): Node {
  if (root.id === id) return { ...root, ...patch }
  return {
    ...root,
    children: root.children.map((c) => updateNode(c, id, patch)),
  }
}

// Returns a new tree without the matched node (root cannot be deleted).
export function deleteNode(root: Node, id: string): Node {
  return {
    ...root,
    children: root.children
      .filter((c) => c.id !== id)
      .map((c) => deleteNode(c, id)),
  }
}

// Returns a new tree with child appended to the matched parent (root if not found).
export function addChild(root: Node, parentId: string | null, child: Node): Node {
  if (parentId === null || root.id === parentId) {
    return { ...root, children: [...root.children, child] }
  }
  return {
    ...root,
    children: root.children.map((c) => addChild(c, parentId, child)),
  }
}

// Insert child at a specific index within parentId's children. -1 = append.
export function insertChild(root: Node, parentId: string, child: Node, index: number): Node {
  if (root.id === parentId) {
    const kids = [...root.children]
    const at = index < 0 || index > kids.length ? kids.length : index
    kids.splice(at, 0, child)
    return { ...root, children: kids }
  }
  return {
    ...root,
    children: root.children.map((c) => insertChild(c, parentId, child, index)),
  }
}

// Move a node one slot up/down within its parent's children. No-op at edges.
export function moveSibling(root: Node, id: string, dir: -1 | 1): Node {
  const found = findNode(root, id)
  if (!found || !found.parent) return root
  const siblings = found.parent.children
  const i = found.index
  const j = i + dir
  if (j < 0 || j >= siblings.length) return root
  // Rebuild root with the swapped parent.
  return replaceNode(root, found.parent.id, (p) => {
    const kids = [...p.children]
    ;[kids[i], kids[j]] = [kids[j], kids[i]]
    return { ...p, children: kids }
  })
}

// Reparent: remove nodeId from its current location, insert under newParentId at index.
// Guards: cannot drop a node into itself or one of its own descendants.
export function reparent(root: Node, nodeId: string, newParentId: string, index: number): Node {
  const moving = findNode(root, nodeId)
  if (!moving || nodeId === newParentId || isDescendant(moving.node, newParentId)) {
    return root
  }
  const detached = deleteNode(root, nodeId)
  return insertChild(detached, newParentId, moving.node, index)
}

// True if `descendantId` is found anywhere under `node`.
export function isDescendant(node: Node, descendantId: string): boolean {
  return node.children.some(
    (c) => c.id === descendantId || isDescendant(c, descendantId),
  )
}

// Immutable node replacement by id via callback.
export function replaceNode(root: Node, id: string, fn: (n: Node) => Node): Node {
  if (root.id === id) return fn(root)
  return {
    ...root,
    children: root.children.map((c) => replaceNode(c, id, fn)),
  }
}

// Deep clone via structuredClone (available in modern browsers + Node 17+).
export function cloneTree<T>(t: T): T {
  return structuredClone(t)
}

export function countNodes(root: Node): number {
  return 1 + root.children.reduce((n, c) => n + countNodes(c), 0)
}

export function treeFromRoot(root: Node): TreeShape {
  return { root }
}
