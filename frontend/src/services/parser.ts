// HTML/Vue template → builder tree parser.
// Uses @vue/compiler-dom for robust HTML + Vue template parsing.
// ponytail: no SFC script/style extraction — strip manually before parse.
// Add proper SFC support (compiler-sfc) when .vue import demand grows.
import { parse as vueParse } from '@vue/compiler-dom'
import type { Node, NodeType } from '@/types/page-builder'
import type { TemplateChildNode, ElementNode, TextNode } from '@vue/compiler-dom'

let _counter = 0
function uid(): string {
  return `import-${++_counter}-${Date.now().toString(36)}`
}

const TAG_MAP: Record<string, NodeType> = {
  div: 'frame',
  section: 'section',
  header: 'section',
  footer: 'section',
  nav: 'section',
  aside: 'section',
  main: 'section',
  article: 'section',
  span: 'text',
  p: 'text',
  h1: 'heading',
  h2: 'heading',
  h3: 'heading',
  h4: 'heading',
  h5: 'heading',
  h6: 'heading',
  img: 'image',
  button: 'button',
  a: 'link',
}

function headingLevel(tag: string): number | undefined {
  if (tag === 'h1') return 1
  if (tag === 'h2') return 2
  if (tag === 'h3') return 3
  if (tag === 'h4') return 4
  if (tag === 'h5') return 5
  if (tag === 'h6') return 6
  return undefined
}

// ── Public API ──────────────────────────────────────────────────

/** Parse HTML string → builder tree root Node. */
export function parseHTML(html: string): Node {
  _counter = 0
  const ast = vueParse(html.trim())
  const root = ast.children?.[0]
  if (!root) return emptyRoot()
  return walkNode(root)
}

/** Parse Vue SFC template content → builder tree root Node. */
export function parseVueTemplate(template: string): Node {
  // Strip <script> and <style> blocks so the parser only sees template markup.
  const cleaned = template
    .replace(/<script[\s\S]*?<\/script>/gi, '')
    .replace(/<style[\s\S]*?<\/style>/gi, '')
    .trim()
  return parseHTML(cleaned)
}

function emptyRoot(): Node {
  return { id: uid(), type: 'frame', name: 'root', props: {}, classes: [], children: [] }
}

// ── AST walker ──────────────────────────────────────────────────

function walkNode(ast: TemplateChildNode | TextNode | ElementNode): Node {
  // Text node
  if (ast.type === 2 && ast.content != null) {
    const text = (ast.content || '').trim()
    if (!text) return null as unknown as Node // filtered in parent
    return {
      id: uid(),
      type: 'text',
      name: 'text',
      props: { text },
      classes: [],
      children: [],
    }
  }

  // Regular element
  if (ast.type === 1 && (ast.tagType === 0 || ast.tagType === undefined)) {
    const tag = ast.tag || 'div'
    const nodeType = TAG_MAP[tag] || 'frame'
    const props: Record<string, any> = {}
    const classes: string[] = []

    for (const p of ast.props || []) {
      // Static attribute
      if (p.type === 6) {
        const name = p.name
        const val = p.value?.content ?? ''
        if (name === 'class') {
          classes.push(...val.split(/\s+/).filter(Boolean))
        } else if (name === 'style') {
          // ponytail: inline style not mapped to classes. Skip.
        } else if (name === 'id') {
          // Store id in props for round-trip fidelity
          props.placeholder = `#${val}`
        } else {
          props[name] = val
        }
      }
    }

    // Heading level
    const level = headingLevel(tag)
    if (level !== undefined) props.level = level

    // Set text for leaf-like elements
    const isLeaf = tag === 'img'
    if (!isLeaf && ast.children?.length === 1 && ast.children[0].type === 2) {
      const txt = (ast.children[0].content || '').trim()
      if (txt) props.text = txt
    }

    // Build name
    const name = tag !== 'div' ? tag : 'frame'

    // Children (skip text children that were consumed as props.text)
    const children: Node[] = []
    if (!isLeaf) {
      const hasTextChild = ast.children?.length === 1 && ast.children[0].type === 2 && !!props.text
      if (!hasTextChild) {
        for (const child of ast.children || []) {
          const n = walkNode(child)
          if (n) children.push(n)
        }
      }
    }

    return {
      id: uid(),
      type: nodeType,
      name,
      props,
      classes,
      children,
    }
  }

  // ponytail: directive nodes (v-if, v-for, v-bind) → skip with a placeholder comment
  // Add proper directive parsing when full Vue SFC round-trip is needed.
  // @ts-expect-error — directive/if/for nodes have non-overlapping numeric types in the compiler DOM types
  if (ast.type === 7 || ast.type === 9 || ast.type === 10 || ast.type === 11) {
    // Return a frame with a note in name
    return {
      id: uid(),
      type: 'frame',
      name: `vue-directive-${ast.type}`,
      props: { text: `<!-- Vue directive (type ${ast.type}) → not parsed yet -->` },
      classes: [],
      children: [],
    }
  }

  return null as unknown as Node
}
