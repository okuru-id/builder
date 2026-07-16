// Converts arbitrary Tailwind bracket classes (colors, sizes, gradients) into
// inline styles for the builder canvas. The vite-built admin CSS only emits
// utilities present at build time, so runtime DB classes like bg-[#1a2e1f]
// would be invisible. The published page uses the Tailwind browser CDN which
// scans the DOM and handles arbitrary values natively; this resolver bridges
// the builder to match.
// ponytail: covers base-state solid colors, sizes, and linear gradients.
// Pseudo-state variants (hover:/focus:) are ignored — they can't be expressed
// inline. Upgrade path: iframe canvas with the browser CDN for full parity.
const HEX_RE = /^#([0-9a-f]{3,8})$/i

function hexToRgb(hex: string): [number, number, number] | null {
  const m = HEX_RE.exec(hex)
  if (!m) return null
  let h = m[1]
  if (h.length === 3) h = h.split('').map((c) => c + c).join('')
  else if (h.length === 4) h = h.split('').slice(0, 3).map((c) => c + c).join('')
  else if (h.length === 8) h = h.slice(0, 6)
  if (h.length !== 6) return null
  const n = parseInt(h, 16)
  return [(n >> 16) & 255, (n >> 8) & 255, n & 255]
}

function colorVal(hex: string, op?: string): string | null {
  const rgb = hexToRgb(hex)
  if (!rgb) return null
  if (op == null) return `rgb(${rgb.join(',')})`
  const a = Math.max(0, Math.min(100, parseInt(op, 10))) / 100
  return `rgba(${rgb.join(',')},${a})`
}

const GRAD_DIR: Record<string, string> = {
  t: 'to top', b: 'to bottom', l: 'to left', r: 'to right',
  tl: 'to top left', tr: 'to top right', bl: 'to bottom left', br: 'to bottom right',
}

export function resolveArbitraryStyles(classes: string[]): Record<string, string> {
  const style: Record<string, string> = {}
  let gradDir: string | null = null
  let gradFrom: string | null = null
  let gradVia: string | null = null
  let gradTo: string | null = null
  for (const raw of classes) {
    // Skip pseudo-state variants entirely (can't be expressed as inline
    // styles). Only strip responsive prefixes (sm/md/lg/xl/2xl) so the
    // builder emulates the current breakpoint.
    if (/^(hover:|focus:|active:|group-hover:|focus-within:)/.test(raw)) continue
    const c = raw.replace(/^(sm:|md:|lg:|xl:|2xl:)+/, '')
    const gm = c.match(/^bg-gradient-to-([a-z]{1,2})$/)
    if (gm) {
      gradDir = GRAD_DIR[gm[1]!] ?? gradDir
      continue
    }
    const cm = c.match(/^(bg|text|border|ring|placeholder|from|via|to)-\[#([0-9a-f]{3,8})\](?:\/(\d{1,3}))?$/i)
    if (cm) {
      const [, prop, hex, op] = cm
      const val = colorVal('#' + hex, op)
      if (!val) continue
      switch (prop) {
        case 'bg': style['background-color'] = val; break
        case 'text': style['color'] = val; break
        case 'border': style['border-color'] = val; break
        case 'ring': style['--tw-ring-color'] = val; break
        case 'placeholder': style['--tw-placeholder-color'] = val; break
        case 'from': gradFrom = val; break
        case 'via': gradVia = val; break
        case 'to': gradTo = val; break
      }
      continue
    }
    const sm = c.match(/^(w|h)-\[\s*(\d+(?:\.\d+)?)(px|rem|em|%|vw|vh|fr)?\s*\]$/)
    if (sm) {
      const [, prop, num, unit = 'px'] = sm
      style[prop === 'w' ? 'width' : 'height'] = `${num}${unit}`
      continue
    }
    const om = c.match(/^opacity-\[(\d+(?:\.\d+)?)\]$/)
    if (om) style['opacity'] = om[1]!
  }
  if (gradDir && (gradFrom || gradVia || gradTo)) {
    const stops = [gradFrom, gradVia, gradTo].filter(Boolean).join(',')
    style['background-image'] = `linear-gradient(${gradDir}, ${stops})`
  }
  return style
}
