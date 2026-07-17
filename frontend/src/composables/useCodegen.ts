// TypeScript mirror of the Go codegen. Used for client-side HTML preview/download.
// Determinisitc: same tree → same bytes. 2-space indent, self-closing void elements.
import type { Node, TreeShape } from '@/types/page-builder'
import { ICONS } from '@/lib/icon-map'

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
  const props = n.props ?? {}
  for (const name of extra) {
    const v = props[name]
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
    case 'form':
      return renderForm(n, depth)
    case 'frame':
    case 'grid':
    case 'component':
    default:
      return renderContainer('div', n, depth)
    case 'divider':
      return renderSelfClosing('hr', n, depth)
    case 'input':
      return renderInput(n, depth)
    case 'icon':
      return renderIcon(n, depth)
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
  return `${ind}<${tag}${attrStr(n)}>${htmlEscape(n.props?.text ?? '')}</${tag}>\n`
}

function renderHeading(n: Node, depth: number): string {
  const lvl = Math.max(1, Math.min(6, Number(n.props?.level) || 2))
  return renderLeaf(`h${lvl}`, n, depth)
}

function renderSelfClosing(tag: string, n: Node, depth: number, ...attrs: string[]): string {
  return `${indent(depth)}<${tag}${attrStr(n, attrs)} />\n`
}

function renderIcon(n: Node, depth: number): string {
  const iconName = n.props?.icon
  let svgContent = ''
  if (iconName && ICONS[iconName]) {
    const paths = ICONS[iconName]
    svgContent = paths.map((p: [string, Record<string,string>]) => `<path d="${p[1].d}" />`).join('')
  } else {
    svgContent = `<circle cx="12" cy="12" r="4" fill="currentColor" />`
  }
  const ind = indent(depth)
  return `${ind}<span${attrStr(n)}><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">${svgContent}</svg></span>\n`
}

function renderInput(n: Node, depth: number): string {
  const ind = indent(depth)
  const p = n.props ?? {}
  const label = p.label ? `${ind}<label class="text-sm font-medium">${htmlEscape(p.label)}</label>\n` : ''
  const inputType = p.inputType ?? 'text'
  const placeholder = p.placeholder ? ` placeholder="${htmlEscape(p.placeholder)}"` : ''
  const required = p.required ? ' required' : ''
  const inputEl = `${ind}<input type="${inputType}"${placeholder}${required} class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm" />\n`
  // Wrap in a div with classes from the node
  return `${ind}<div${attrStr(n)}>\n${label}${inputEl}${ind}</div>\n`
}

function renderForm(n: Node, depth: number): string {
  const ind = indent(depth)
  const p = n.props ?? {}
  const action = p.action ? ` action="${htmlEscape(p.action)}"` : ''
  const method = p.method ? ` method="${htmlEscape(p.method)}"` : ' method="POST"'
  let out = `${ind}<form${attrStr(n)}${action}${method}>\n`
  for (const child of n.children) {
    out += renderNode(child, depth + 1)
  }
  out += `${ind}</form>\n`
  return out
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
  return `${ind}<a${attrStr(n, ['href'])}>${htmlEscape(n.props?.text ?? '')}</a>\n`
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
<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
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
  a.rel = 'noopener'
  a.style.display = 'none'
  // Firefox/Safari ignore .click() on detached elements — append first.
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  // Defer revoke: some browsers start fetch after click() returns.
  setTimeout(() => URL.revokeObjectURL(url), 1000)
}
