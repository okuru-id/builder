export type Breakpoint = 'desktop' | 'tablet' | 'mobile'

// Tree JSON node types — single source of truth shared by canvas, panels, codegen.

export type NodeType =
  | 'frame'
  | 'text'
  | 'heading'
  | 'image'
  | 'icon'
  | 'button'
  | 'link'
  | 'section'
  | 'divider'
  | 'grid'
  | 'form'
  | 'input'
  | 'component'

export interface NodeProps {
  text?: string
  src?: string
  alt?: string
  href?: string
  level?: number
  placeholder?: string
  icon?: string
  iconVariant?: 'outline' | 'filled'
  label?: string
  inputType?: string
  required?: boolean
  action?: string
  method?: string
  options?: string
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
  hidden?: boolean
  /** Hide the node on the given breakpoints. Maps 1:1 to Tailwind responsive
   *  utilities in published HTML: mobile→`hidden md:block`, tablet→`md:hidden lg:block`,
   *  desktop→`lg:hidden`. */
  hiddenOn?: Breakpoint[]
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
  key?: string
  name: string
  is_system: boolean
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
  'form',
  'link',
  'component',
])
// Node types allowed as top-level palette additions.
export const PALETTE_TYPES: NodeType[] = ['frame', 'section', 'grid', 'text', 'heading', 'image', 'icon', 'button', 'link', 'divider', 'form', 'input']
