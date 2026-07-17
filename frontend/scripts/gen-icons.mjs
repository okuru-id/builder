// Generate frontend/src/lib/icon-map.ts + backend/app/services/icons.json from
// ALL @tabler/icons-vue icons (outline + filled) via SSR render + path-d extraction.
// Filled variants are keyed as `${name}Filled` (e.g. IconStarFilled).
// Usage: node scripts/gen-icons.mjs
import { createSSRApp, h } from 'vue'
import { renderToString } from 'vue/server-renderer'
import * as T from '@tabler/icons-vue'
import { writeFileSync } from 'node:fs'

const names = Object.keys(T).filter((k) => k.startsWith('Icon'))
const out = {}
for (const name of names) {
  const Comp = T[name]
  const app = createSSRApp({ render: () => h(Comp, { size: 24 }) })
  const html = await renderToString(app)
  const segs = []
  for (const m of html.matchAll(/<(path|circle|rect|line|polyline|polygon)\b([^>]*)>/g)) {
    const tag = m[1]
    const d = (m[2].match(/d="([^"]+)"/) || [])[1]
    if (d) segs.push([tag, { d, key: `${name}-${segs.length}` }])
  }
  out[name] = segs
}

const list = Object.keys(out)
const filledCount = list.filter((n) => n.endsWith('Filled')).length
const ts =
  '// Auto-generated from @tabler/icons-vue — ALL icons (' + list.length +
  ', ' + filledCount + ' filled).\n' +
  '// Regenerate: node scripts/gen-icons.mjs\n' +
  'export const ICONS: Record<string, [string, Record<string,string>][]> = ' +
  JSON.stringify(out) + '\n\nexport const ICON_LIST = Object.keys(ICONS)\n'
writeFileSync('src/lib/icon-map.ts', ts)
writeFileSync('../backend/app/services/icons.json', JSON.stringify(out))
console.log('icons:', list.length, '(filled:', filledCount + ')')
console.log('icon-map.ts:', (ts.length / 1024).toFixed(0), 'KB')
console.log('icons.json:', (JSON.stringify(out).length / 1024).toFixed(0), 'KB')
