// Round-trip test: generate → parse → generate = identical.
// Run with: bun test frontend/tests/parser.test.ts
import { describe, it, expect } from 'bun:test'
import { parseHTML } from '@/services/parser'
import { generateFragment, collectClasses } from '@/composables/useCodegen'
import type { Node, TreeShape } from '@/types/page-builder'

// ── Helpers ─────────────────────────────────────────────────────

function makeNode(overrides: Partial<Node>): Node {
  return {
    id: 'test-1',
    type: 'frame',
    name: 'test',
    props: {},
    classes: [],
    children: [],
    ...overrides,
  }
}

function wrap(n: Node): TreeShape {
  return { root: n }
}

function classesOnly(n: Node): string[] {
  return [...new Set(n.classes.filter(Boolean))].sort()
}

// ── Tests ───────────────────────────────────────────────────────

describe('HTML → tree parsing', () => {
  it('parses a plain div', () => {
    const n = parseHTML('<div class="flex gap-4">hello</div>')
    expect(n.type).toBe('frame')
    expect(classesOnly(n)).toEqual(['flex', 'gap-4'])
    expect(n.props.text).toBe('hello')
  })

  it('parses heading with level', () => {
    const n = parseHTML('<h2 class="text-2xl">Title</h2>')
    expect(n.type).toBe('heading')
    expect(n.props.level).toBe(2)
    expect(n.props.text).toBe('Title')
  })

  it('parses image with src/alt', () => {
    const n = parseHTML('<img src="/pic.jpg" alt="Pic" class="rounded-lg" />')
    expect(n.type).toBe('image')
    expect(n.props.src).toBe('/pic.jpg')
    expect(n.props.alt).toBe('Pic')
  })

  it('parses nested structure', () => {
    const n = parseHTML(`<section class="p-4">
      <h1>Judul</h1>
      <p>Deskripsi</p>
    </section>`)
    expect(n.type).toBe('section')
    expect(n.children.length).toBe(2)
    expect(n.children[0].type).toBe('heading')
    expect(n.children[0].props.level).toBe(1)
  })

  it('parses link with href', () => {
    const n = parseHTML('<a href="/contact" class="text-blue-600">Kontak</a>')
    expect(n.type).toBe('link')
    expect(n.props.href).toBe('/contact')
    expect(n.props.text).toBe('Kontak')
  })
})

describe('Round-trip: generate → parse → generate', () => {
  function rt(n: Partial<Node>): { gen1: string; gen2: string; tree: TreeShape } {
    const node = makeNode(n)
    const tree = wrap(node)
    const gen1 = generateFragment(tree)
    const reparsed = parseHTML(gen1)
    const gen2 = generateFragment(wrap(reparsed))
    return { gen1, gen2, tree: wrap(reparsed) }
  }

  it('empty frame round-trips', () => {
    const { gen1, gen2 } = rt({
      type: 'frame',
      name: 'root',
      classes: [],
    })
    expect(gen1).toBe(gen2)
  })

  it('text with class round-trips', () => {
    const { gen1, gen2 } = rt({
      type: 'text',
      props: { text: 'Hello World' },
      classes: ['text-lg', 'font-bold'],
    })
    expect(gen1).toBe(gen2)
  })

  it('heading round-trips', () => {
    const { gen1, gen2 } = rt({
      type: 'heading',
      props: { text: 'Section Title', level: 3 },
      classes: ['text-3xl'],
    })
    expect(gen1).toBe(gen2)
  })

  it('button round-trips', () => {
    const { gen1, gen2 } = rt({
      type: 'button',
      props: { text: 'Click Me' },
      classes: ['bg-blue-500', 'text-white'],
    })
    expect(gen1).toBe(gen2)
  })

  it('link round-trips', () => {
    const { gen1, gen2 } = rt({
      type: 'link',
      props: { text: 'Visit', href: 'https://example.com' },
      classes: ['text-blue-600'],
    })
    expect(gen1).toBe(gen2)
  })

  it('image round-trips', () => {
    const { gen1, gen2 } = rt({
      type: 'image',
      props: { src: '/hero.jpg', alt: 'Hero image' },
      classes: ['rounded-lg', 'w-full'],
    })
    expect(gen1).toBe(gen2)
  })

  it('section with children round-trips', () => {
    const node = makeNode({
      type: 'section',
      classes: ['p-4', 'flex', 'gap-4'],
      children: [
        makeNode({ type: 'text', props: { text: 'A' }, classes: ['font-bold'] }),
        makeNode({ type: 'text', props: { text: 'B' }, classes: [] }),
      ],
    })
    const tree = wrap(node)
    const gen1 = generateFragment(tree)
    const reparsed = parseHTML(gen1)
    const gen2 = generateFragment(wrap(reparsed))
    expect(gen1).toBe(gen2)
  })
})

describe('Safelist', () => {
  it('collectClasses extracts all unique classes', () => {
    const tree = wrap(makeNode({
      classes: ['flex', 'p-4', 'gap-4', 'p-4'],
      children: [
        makeNode({ type: 'text', classes: ['text-lg', 'font-bold', 'text-lg'] }),
      ],
    }))
    const classes = collectClasses(tree)
    expect(classes).toEqual(['flex', 'font-bold', 'gap-4', 'p-4', 'text-lg'])
  })
})
