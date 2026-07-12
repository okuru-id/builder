// Reverse mapping from Tailwind utility class → inspector property.
// Used when loading a parsed tree back into the builder so the inspector
// populates controls correctly.
//
// ponytail: partial coverage — add prefixes as you encounter them.
// A complete mapping would require parsing Tailwind's entire config.

export interface ClassInfo {
  section: 'layout' | 'typography' | 'spacing' | 'background' | 'border' | 'size' | 'other'
  property: string
  value: string
}

// Prefix → section + property name (within that section).
// The value is whatever follows the prefix + '-'.
const CLASS_MAP: Record<string, { section: ClassInfo['section']; property: string }> = {
  // Layout
  flex: { section: 'layout', property: 'display' },
  'flex-row': { section: 'layout', property: 'direction' },
  'flex-col': { section: 'layout', property: 'direction' },
  'flex-wrap': { section: 'layout', property: 'wrap' },
  'flex-nowrap': { section: 'layout', property: 'wrap' },
  'items-start': { section: 'layout', property: 'align' },
  'items-center': { section: 'layout', property: 'align' },
  'items-end': { section: 'layout', property: 'align' },
  'items-stretch': { section: 'layout', property: 'align' },
  'items-baseline': { section: 'layout', property: 'align' },
  'justify-start': { section: 'layout', property: 'justify' },
  'justify-center': { section: 'layout', property: 'justify' },
  'justify-end': { section: 'layout', property: 'justify' },
  'justify-between': { section: 'layout', property: 'justify' },
  'justify-around': { section: 'layout', property: 'justify' },
  'justify-evenly': { section: 'layout', property: 'justify' },
  gap: { section: 'layout', property: 'gap' },
  block: { section: 'layout', property: 'display' },
  grid: { section: 'layout', property: 'display' },
  hidden: { section: 'layout', property: 'display' },

  // Typography
  'font-thin': { section: 'typography', property: 'fontWeight' },
  'font-extralight': { section: 'typography', property: 'fontWeight' },
  'font-light': { section: 'typography', property: 'fontWeight' },
  'font-normal': { section: 'typography', property: 'fontWeight' },
  'font-medium': { section: 'typography', property: 'fontWeight' },
  'font-semibold': { section: 'typography', property: 'fontWeight' },
  'font-bold': { section: 'typography', property: 'fontWeight' },
  'font-extrabold': { section: 'typography', property: 'fontWeight' },
  'font-black': { section: 'typography', property: 'fontWeight' },
  'text-left': { section: 'typography', property: 'textAlign' },
  'text-center': { section: 'typography', property: 'textAlign' },
  'text-right': { section: 'typography', property: 'textAlign' },
  'text-justify': { section: 'typography', property: 'textAlign' },
  uppercase: { section: 'typography', property: 'textTransform' },
  lowercase: { section: 'typography', property: 'textTransform' },
  capitalize: { section: 'typography', property: 'textTransform' },
  'normal-case': { section: 'typography', property: 'textTransform' },
  'text-xs': { section: 'typography', property: 'fontSize' },
  'text-sm': { section: 'typography', property: 'fontSize' },
  'text-base': { section: 'typography', property: 'fontSize' },
  'text-lg': { section: 'typography', property: 'fontSize' },
  'text-xl': { section: 'typography', property: 'fontSize' },
  'text-2xl': { section: 'typography', property: 'fontSize' },
  'text-3xl': { section: 'typography', property: 'fontSize' },
  'text-4xl': { section: 'typography', property: 'fontSize' },
  'text-5xl': { section: 'typography', property: 'fontSize' },
  'text-6xl': { section: 'typography', property: 'fontSize' },
  'text-7xl': { section: 'typography', property: 'fontSize' },
  'text-8xl': { section: 'typography', property: 'fontSize' },
  'text-9xl': { section: 'typography', property: 'fontSize' },

  // Spacing
  p: { section: 'spacing', property: 'padding' },
  px: { section: 'spacing', property: 'paddingX' },
  py: { section: 'spacing', property: 'paddingY' },
  pt: { section: 'spacing', property: 'paddingTop' },
  pr: { section: 'spacing', property: 'paddingRight' },
  pb: { section: 'spacing', property: 'paddingBottom' },
  pl: { section: 'spacing', property: 'paddingLeft' },
  m: { section: 'spacing', property: 'margin' },
  mx: { section: 'spacing', property: 'marginX' },
  my: { section: 'spacing', property: 'marginY' },
  mt: { section: 'spacing', property: 'marginTop' },
  mr: { section: 'spacing', property: 'marginRight' },
  mb: { section: 'spacing', property: 'marginBottom' },
  ml: { section: 'spacing', property: 'marginLeft' },

  // Border
  border: { section: 'border', property: 'borderWidth' },
  'border-0': { section: 'border', property: 'borderWidth' },
  'border-2': { section: 'border', property: 'borderWidth' },
  'border-4': { section: 'border', property: 'borderWidth' },
  'border-8': { section: 'border', property: 'borderWidth' },
  'rounded-none': { section: 'border', property: 'borderRadius' },
  'rounded-sm': { section: 'border', property: 'borderRadius' },
  'rounded': { section: 'border', property: 'borderRadius' },
  'rounded-md': { section: 'border', property: 'borderRadius' },
  'rounded-lg': { section: 'border', property: 'borderRadius' },
  'rounded-xl': { section: 'border', property: 'borderRadius' },
  'rounded-2xl': { section: 'border', property: 'borderRadius' },
  'rounded-3xl': { section: 'border', property: 'borderRadius' },
  'rounded-full': { section: 'border', property: 'borderRadius' },

  // Size
  'w-auto': { section: 'size', property: 'width' },
  'w-full': { section: 'size', property: 'width' },
  'w-screen': { section: 'size', property: 'width' },
  'w-min': { section: 'size', property: 'width' },
  'w-max': { section: 'size', property: 'width' },
  'w-fit': { section: 'size', property: 'width' },
  'h-auto': { section: 'size', property: 'height' },
  'h-full': { section: 'size', property: 'height' },
  'h-screen': { section: 'size', property: 'height' },
  'h-min': { section: 'size', property: 'height' },
  'h-max': { section: 'size', property: 'height' },
  'h-fit': { section: 'size', property: 'height' },
  'min-w-': { section: 'size', property: 'minWidth' },
  'max-w-': { section: 'size', property: 'maxWidth' },
  'min-h-': { section: 'size', property: 'minHeight' },
  'max-h-': { section: 'size', property: 'maxHeight' },
}

/** Classify a single Tailwind class → inspector info, or null if unknown. */
export function classifyClass(cls: string): ClassInfo | null {
  // Check exact match first
  if (CLASS_MAP[cls]) return { ...CLASS_MAP[cls], value: cls }

  // Try prefix match (e.g. 'gap-4' matches 'gap')
  const dash = cls.indexOf('-')
  if (dash > 0) {
    const prefix = cls.slice(0, dash)
    const base = CLASS_MAP[prefix]
    if (base) return { ...base, value: cls.slice(dash + 1) }
  }

  // Arbitrary value: text-[#abc], bg-[#abc], w-[200px]
  const arbMatch = cls.match(/^(\w+)-\[(.+)\]$/)
  if (arbMatch) {
    const base = CLASS_MAP[arbMatch[1]]
    if (base) return { ...base, value: `[${arbMatch[2]}]` }
  }

  // Color class: text-blue-500, bg-gray-100, border-red-600
  const colorMatch = cls.match(/^(text|bg|border)-(\w[-\w]*)$/)
  if (colorMatch) {
    const prefix = colorMatch[1]
    const val = colorMatch[2]
    const section: ClassInfo['section'] = prefix === 'text' ? 'typography' : prefix === 'bg' ? 'background' : 'border'
    return { section, property: prefix === 'text' ? 'color' : prefix === 'bg' ? 'bgColor' : 'borderColor', value: val }
  }

  return null
}
