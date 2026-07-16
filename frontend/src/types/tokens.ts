// Token config for brand design system.
// Tokens map to Tailwind classes at codegen time. Editor uses standard
// Tailwind utilities + arbitrary values.

export interface TokenConfig {
  colors: Record<string, string>  // token name → hex, e.g. { primary: '#3b82f6' }
  fonts: Record<string, string>   // token name → font-family, e.g. { heading: 'Inter' }
  radius: Record<string, string>  // token name → value, e.g. { default: '0.5rem' }
}

export const DEFAULT_TOKENS: TokenConfig = {
  colors: { primary: '#3b82f6', secondary: '#6b7280', accent: '#8b5cf6', neutral: '#f5f5f4', 'neutral-content': '#292524' },
  fonts: { heading: 'Inter', body: 'Inter' },
  radius: { default: '0.5rem', sm: '0.25rem', lg: '1rem' },
}

// ── Class mutation helpers ──────────────────────────────────────

// Replace classes matching any pattern with an optional new class.
// A pattern like 'gap' matches 'gap-2', 'gap-4', etc.
// A pattern like 'flex' matches 'flex' exactly.
export function replaceClass(classes: string[], patterns: string[], add: string | null): string[] {
  const out = classes.filter((c) => !patterns.some((p) => {
    if (c === p) return true
    // Arbitrary-value prefixes (e.g. 'text-[', 'bg-[', 'bg-[url(', 'shadow-[')
    // end in a bracket and must match by raw startsWith: 'text-[#fff]'
    // starts with 'text-[' but NOT 'text-[-', so the prefix-family branch
    // below would miss it and leave stale color classes behind.
    if (/[[(*]$/.test(p)) return c.startsWith(p)
    return c.startsWith(p + '-')
  }))
  if (add) out.push(add)
  return out
}

// Get first class value by prefix, e.g. currentClass(classes, 'gap') → '4'
export function currentClass(classes: string[], prefix: string): string | null {
  const found = classes.find((c) => c.startsWith(prefix + '-'))
  if (!found) return null
  // handle arbitrary values: text-[#abc] →  [#abc]
  const val = found.slice(prefix.length + 1)
  return val || null
}

// Exact-match lookup against a candidate set. Use when multiple token families
// share a prefix (text-3xl vs text-blue-600 vs text-center) and prefix-match
// would collide. Returns the candidate class present, or null.
export function currentFromSet(classes: string[], candidates: readonly string[]): string | null {
  for (const c of classes) {
    if (candidates.includes(c)) return c
  }
  return null
}

// Extract the inner value of a `prefix-[value]` arbitrary class (e.g.
// `text-[#abc]` → `#abc`, `w-[320px]` → `320px`). Returns null if absent.
export function currentArbitrary(classes: string[], prefix: string): string | null {
  const p = prefix + '-['
  for (const c of classes) {
    if (c.startsWith(p) && c.endsWith(']')) {
      return c.slice(p.length, -1)
    }
  }
  return null
}

// Check if a boolean class is present.
export function hasClass(classes: string[], cls: string): boolean {
  return classes.includes(cls)
}

// Mutually-exclusive groups (display modes).
export const DISPLAY_CLASSES = ['flex', 'block', 'grid', 'inline-flex', 'hidden', 'inline', 'inline-block'] as const

// ── Tailwind value constants ────────────────────────────────────

export const SPACING = ['0', 'px', '0.5', '1', '1.5', '2', '2.5', '3', '3.5', '4', '5', '6', '7', '8', '9', '10', '11', '12', '14', '16', '20', '24', '28', '32', '36', '40', '44', '48', '52', '56', '60', '64', '72', '80', '96'] as const

export const FONT_SIZES = ['xs', 'sm', 'base', 'lg', 'xl', '2xl', '3xl', '4xl', '5xl', '6xl', '7xl', '8xl', '9xl'] as const

export const FONT_WEIGHTS = ['thin', 'extralight', 'light', 'normal', 'medium', 'semibold', 'bold', 'extrabold', 'black'] as const

export const FONT_FAMILIES = ['sans', 'serif', 'mono'] as const

// Google Fonts loaded via <link> in index.html (admin preview) and codegen
// head (published HTML). Mapped to .gfont-* CSS classes (see style.css +
// codegen.go) — Tailwind browser CDN in publish only knows font-sans/serif/mono.
export const GOOGLE_FONTS = [
  { value: 'inter', label: 'Inter', family: "'Inter', sans-serif" },
  { value: 'poppins', label: 'Poppins', family: "'Poppins', sans-serif" },
  { value: 'roboto', label: 'Roboto', family: "'Roboto', sans-serif" },
  { value: 'montserrat', label: 'Montserrat', family: "'Montserrat', sans-serif" },
  { value: 'opensans', label: 'Open Sans', family: "'Open Sans', sans-serif" },
  { value: 'playfair', label: 'Playfair Display', family: "'Playfair Display', serif" },
  { value: 'lora', label: 'Lora', family: "'Lora', serif" },
  { value: 'merriweather', label: 'Merriweather', family: "'Merriweather', serif" },
  { value: 'jetbrains', label: 'JetBrains Mono', family: "'JetBrains Mono', monospace" },
  { value: 'sourcesans', label: 'Source Sans 3', family: "'Source Sans 3', sans-serif" },
  { value: 'nunito', label: 'Nunito', family: "'Nunito', sans-serif" },
  { value: 'raleway', label: 'Raleway', family: "'Raleway', sans-serif" },
  { value: 'jakarta', label: 'Plus Jakarta Sans', family: "'Plus Jakarta Sans', sans-serif" },
] as const

export const BORDER_RADII = ['none', 'sm', 'md', 'lg', 'xl', '2xl', '3xl', 'full'] as const

export const BORDER_WIDTHS = ['0', '1', '2', '4', '8'] as const

export const TEXT_ALIGNS = ['left', 'center', 'right', 'justify'] as const

export const TEXT_TRANSFORMS = ['none', 'uppercase', 'lowercase', 'capitalize'] as const

export const FLEX_DIRECTIONS = ['row', 'col', 'row-reverse', 'col-reverse'] as const

export const ALIGN_ITEMS = ['start', 'center', 'end', 'stretch', 'baseline'] as const

export const JUSTIFY_CONTENTS = ['start', 'center', 'end', 'between', 'around', 'evenly'] as const

export const FLEX_WRAPS = ['wrap', 'nowrap', 'wrap-reverse'] as const

export const SIZES = ['auto', 'full', 'screen', 'min', 'max', 'fit'] as const
