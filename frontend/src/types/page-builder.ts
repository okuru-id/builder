// Tree JSON node types — single source of truth shared by canvas, panels, codegen.

export type NodeType =
  | 'frame'
  | 'text'
  | 'heading'
  | 'image'
  | 'button'
  | 'link'
  | 'section'
  | 'divider'
  | 'grid'
  | 'component'

export interface NodeProps {
  text?: string
  src?: string
  alt?: string
  href?: string
  level?: number
  placeholder?: string
  [key: string]: unknown
}

export interface Node {
  id: string
  type: NodeType
  name: string
  props: NodeProps
  classes: string[]
  children: Node[]
  componentId?: number
  instanceOverrides?: Partial<Node>
}

export interface TreeShape {
  root: Node
}

export interface Page {
  id: number
  slug: string
  name: string
  status: 'draft' | 'published'
  tree: TreeShape
  published_tree?: TreeShape | null
  published_html?: string
  version: number
  created_at?: string
  updated_at?: string
}

// Component master — reusable tree fragment referenced by instances via componentId.
export interface Component {
  id: number
  name: string
  tree: TreeShape
  created_at?: string
  updated_at?: string
}

// Node types that carry text content (inline-editable).
export const TEXT_TYPES: ReadonlySet<NodeType> = new Set(['text', 'heading'])
// Node types that can have children (containers).
export const CONTAINER_TYPES: ReadonlySet<NodeType> = new Set([
  'frame',
  'section',
  'grid',
  'link',
  'component',
])
// Node types allowed as top-level palette additions.
export const PALETTE_TYPES: NodeType[] = ['frame', 'section', 'grid', 'text', 'heading', 'image', 'button', 'link', 'divider']
