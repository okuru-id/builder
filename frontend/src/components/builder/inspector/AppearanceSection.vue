<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import {
  replaceClass,
  currentArbitrary,
  BORDER_RADII,
} from '@/types/tokens'
import { Label } from '@/components/ui/label'
import {
  NumberField,
  NumberFieldContent,
  NumberFieldDecrement,
  NumberFieldIncrement,
  NumberFieldInput,
} from '@/components/ui/number-field'
import InspectorSection from './InspectorSection.vue'
import { IconWand } from '@tabler/icons-vue'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

// ── Class helpers ────────────────────────────────────────────────
const OPACITY_PRESETS = ['0', '5', '10', '20', '25', '30', '40', '50', '60', '70', '75', '80', '90', '95', '100'].map((s) => `opacity-${s}`)
// Global radius classes (shorthand for all 4 corners) — stripped when a
// per-side value is set so they cannot conflict with per-side overrides.
const GLOBAL_RADIUS = BORDER_RADII.map((r) => {
  if (r === 'md') return 'rounded-md'
  if (r === 'lg') return 'rounded-lg'
  return r === 'none' ? 'rounded-none' : `rounded-${r}`
})

function patch(classes: string[]) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes })
}
function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  patch(replaceClass(node.value.classes, patterns, add))
}

// ── Opacity ──────────────────────────────────────────────────────
function currentOpacityPercent() {
  const arb = currentArbitrary(node.value!.classes, 'opacity')
  if (arb) {
    const n = parseFloat(arb)
    if (!Number.isNaN(n)) return Math.round(n * 100)
  }
  const step = node.value!.classes.find((c) => /^opacity-\d+$/.test(c))
  if (step) return parseInt(step.replace('opacity-', ''), 10)
  return 100
}
function setOpacity(percent: number) {
  const clamped = Math.max(0, Math.min(100, Math.round(percent)))
  const match = OPACITY_PRESETS.includes(`opacity-${clamped}`)
  cls([...OPACITY_PRESETS, 'opacity-['], match ? `opacity-${clamped}` : `opacity-[${(clamped / 100).toFixed(2)}]`)
}

// ── Per-side corner radius ───────────────────────────────────────
type Side = 'tl' | 'tr' | 'br' | 'bl'
const SIDES: { key: Side; label: string }[] = [
  { key: 'tl', label: 'TL' },
  { key: 'tr', label: 'TR' },
  { key: 'br', label: 'BR' },
  { key: 'bl', label: 'BL' },
]
// Named per-side classes to strip when replacing (e.g. rounded-tl-md).
const sideNamed = (s: Side) => BORDER_RADII.map((r) => (r === 'none' ? `rounded-${s}-none` : `rounded-${s}-${r}`))

function sidePx(s: Side): number {
  const arb = currentArbitrary(node.value!.classes, `rounded-${s}`)
  if (arb) {
    const n = parseFloat(arb)
    if (!Number.isNaN(n)) return n
  }
  return 0
}
function setSidePx(s: Side, px: number) {
  const clamped = Math.max(0, Math.round(px))
  // Clear any global shorthand + this side's named/arbitrary, then set.
  const cleaned = replaceClass(node.value!.classes, [...GLOBAL_RADIUS, ...sideNamed(s), `rounded-${s}-[`], null)
  patch(clamped > 0 ? [...cleaned, `rounded-${s}-[${clamped}px]`] : cleaned)
}
</script>

<template>
  <InspectorSection title="Appearance" :icon="IconWand" :show="!!node">

    <!-- Opacity -->
    <div class="space-y-1">
      <Label class="text-[11px] font-medium text-foreground/80">Opacity (%)</Label>
      <NumberField
        :model-value="currentOpacityPercent()"
        :min="0"
        :max="100"
        :step="5"
        @update:model-value="(v) => setOpacity(Number(v ?? 100))"
      >
        <NumberFieldContent>
          <NumberFieldDecrement />
          <NumberFieldInput />
          <NumberFieldIncrement />
        </NumberFieldContent>
      </NumberField>
    </div>

    <!-- Corner Radius — four sides -->
    <div class="space-y-1">
      <Label class="text-[11px] font-medium text-foreground/80">Corner Radius (px)</Label>
      <div class="grid grid-cols-2 gap-1.5">
        <div v-for="s in SIDES" :key="s.key" class="space-y-1">
          <NumberField
            :model-value="sidePx(s.key)"
            :min="0"
            :max="999"
            :step="1"
            @update:model-value="(v) => setSidePx(s.key, Number(v ?? 0))"
          >
            <NumberFieldContent>
              <NumberFieldDecrement />
              <NumberFieldInput />
              <NumberFieldIncrement />
            </NumberFieldContent>
          </NumberField>
          <span class="block text-center text-[9px] font-medium uppercase tracking-wide text-muted-foreground">{{ s.label }}</span>
        </div>
      </div>
      <p class="text-[10px] text-muted-foreground">TL = top-left · TR = top-right · BR = bottom-right · BL = bottom-left. Setting a side replaces any global radius.</p>
    </div>
  </InspectorSection>
</template>
