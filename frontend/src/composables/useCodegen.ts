// TypeScript mirror of the Go codegen. Used for client-side HTML preview/download.
// Determinisitc: same tree → same bytes. 2-space indent, self-closing void elements.
import type { Node, TreeShape } from '@/types/page-builder'

export function htmlEscape(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
}

function indent(depth: number): string {
  return '  '.repeat(depth)
}

function classStr(n: Node): string {
  if (!n.classes || n.classes.length === 0) return ''
  const deduped = [...new Set(n.classes.filter(Boolean))]
  return deduped.join(' ')
}

function attrStr(n: Node, extra: string[] = []): string {
  const parts: string[] = []
  const cls = classStr(n)
  if (cls) parts.push(` class="${cls}"`)
  for (const name of extra) {
    const v = n.props[name]
    if (v) parts.push(` ${name}="${htmlEscape(String(v))}"`)
  }
  return parts.join('')
}

// ── Renderers ──────────────────────────────────────────────────

function renderNode(n: Node, depth: number): string {
  switch (n.type) {
    case 'text':
      return renderLeaf('span', n, depth)
    case 'heading':
      return renderHeading(n, depth)
    case 'image':
      return renderSelfClosing('img', n, depth, 'src', 'alt')
    case 'button':
      return renderLeaf('button', n, depth)
    case 'link':
      return renderLink(n, depth)
    case 'section':
      return renderContainer('section', n, depth)
    case 'frame':
    case 'grid':
    case 'component':
    default:
      return renderContainer('div', n, depth)
    case 'divider':
      return renderSelfClosing('hr', n, depth)
  }
}

function renderContainer(tag: string, n: Node, depth: number): string {
  const ind = indent(depth)
  let out = `${ind}<${tag}${attrStr(n)}>\n`
  for (const child of n.children) {
    out += renderNode(child, depth + 1)
  }
  out += `${ind}</${tag}>\n`
  return out
}

function renderLeaf(tag: string, n: Node, depth: number): string {
  const ind = indent(depth)
  return `${ind}<${tag}${attrStr(n)}>${htmlEscape(n.props.text ?? '')}</${tag}>\n`
}

function renderHeading(n: Node, depth: number): string {
  const lvl = Math.max(1, Math.min(6, Number(n.props.level) || 2))
  return renderLeaf(`h${lvl}`, n, depth)
}

function renderSelfClosing(tag: string, n: Node, depth: number, ...attrs: string[]): string {
  return `${indent(depth)}<${tag}${attrStr(n, attrs)} />\n`
}

function renderLink(n: Node, depth: number): string {
  const ind = indent(depth)
  if (n.children.length > 0) {
    let out = `${ind}<a${attrStr(n, ['href'])}>\n`
    for (const child of n.children) {
      out += renderNode(child, depth + 1)
    }
    out += `${ind}</a>\n`
    return out
  }
  return `${ind}<a${attrStr(n, ['href'])}>${htmlEscape(n.props.text ?? '')}</a>\n`
}

// ── Public API ──────────────────────────────────────────────────

/** Full HTML document with Tailwind CDN. */
export function generateHTML(tree: TreeShape, title: string): string {
  const body = renderNode(tree.root, 2)
  return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>${htmlEscape(title)}</title>
<script src="https://cdn.tailwindcss.com"></script>
</head>
<body>
${body}</body>
</html>`
}

/** Fragment only (no wrapper). */
export function generateFragment(tree: TreeShape): string {
  return renderNode(tree.root, 0)
}

// ── Safelist ────────────────────────────────────────────────────

/** Collect every distinct class from the tree. */
export function collectClasses(tree: TreeShape): string[] {
  const set = new Set<string>()
  function walk(n: Node) {
    for (const c of n.classes) {
      if (c) set.add(c)
    }
    for (const child of n.children) walk(child)
  }
  walk(tree.root)
  return [...set].sort()
}

// ── Download helpers ────────────────────────────────────────────

/** Trigger browser download of a file. */
export function download(filename: string, content: string, mime = 'text/html') {
  const blob = new Blob([content], { type: mime })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}
