# Phase 6: Style Panel & Token System

**Status:** ⬜ todo
**Goal:** Inspector kanan dengan UI controls untuk Tailwind classes. Bukan raw textarea. Token brand system independent dari admin shadcn CSS vars.

## Checklist

- [ ] 6.1 Layout — direction (flex-row/col), align, justify, gap → dropdown + slider
- [ ] 6.2 Typography — font family, size, weight, color → select + color picker
- [ ] 6.3 Spacing — margin/padding per direction
- [ ] 6.4 Background — color, image
- [ ] 6.5 Border — width, color, radius
- [ ] 6.6 Width/Height — min/max/fixed
- [ ] 6.7 Token system — define brand tokens (primary, secondary, font-heading, spacing-unit) → apply sebagai Tailwind classes
- [ ] 6.8 Setiap control mutasi `node.classes[]` → rebuild class string
- [ ] 6.9 Separate Tailwind config editor vs output (token independence)

## Class Mutation Pattern

```typescript
// control → add/remove class dari node.classes[]
function setClass(node: Node, prefix: string, value: string | null) {
  node.classes = node.classes.filter(c => !c.startsWith(prefix + '-'))
  if (value) node.classes.push(`${prefix}-${value}`)
}
// e.g. setClass(node, 'gap', '4') → remove 'gap-*', add 'gap-4'
```

## Token Config

```json
{
  "colors": { "primary": "#3b82f6", "secondary": "#8b5cf6" },
  "fonts":  { "heading": "Inter", "body": "Inter" },
  "radius": { "default": "0.5rem" }
}
```

Map ke `tailwind.config` builder output. Editor pakai config sendiri.

## Files

- Create: `frontend/src/components/builder/inspector/{LayoutSection,TypographySection,SpacingSection,BackgroundSection,BorderSection,SizeSection}.vue`
- Create: `frontend/src/composables/useTokens.ts`
- Create: `frontend/src/types/tokens.ts`

## Commit

```bash
git commit -m "Phase 6: style panel + token system"
```
