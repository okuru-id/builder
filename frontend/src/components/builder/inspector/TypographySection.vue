<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import {
  replaceClass,
  currentFromSet,
  currentArbitrary,
  FONT_SIZES,
  FONT_WEIGHTS,
  FONT_FAMILIES,
  GOOGLE_FONTS,
  TEXT_ALIGNS,
} from '@/types/tokens'
import InspectorSection from './InspectorSection.vue'
import { IconTypography } from '@tabler/icons-vue'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Input } from '@/components/ui/input'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

const showTypography = computed<boolean>(() =>
  !!node.value && ['text', 'heading', 'button', 'link'].includes(node.value.type),
)

// Candidate class sets — avoids text-* prefix collision (size vs color vs align).
const SIZE_CLASSES = FONT_SIZES.map((s) => `text-${s}`)
const WEIGHT_CLASSES = FONT_WEIGHTS.map((w) => `font-${w}`)
const FAMILY_CLASSES = FONT_FAMILIES.map((f) => `font-${f}`)
const GFONT_CLASSES = GOOGLE_FONTS.map((f) => `gfont-${f.value}`)
const ALIGN_CLASSES = TEXT_ALIGNS.map((a) => `text-${a}`)
const TRANSFORM_CLASSES: string[] = ['uppercase', 'lowercase', 'capitalize', 'normal-case']
const TEXT_COLOR_PRESETS = ['gray-900', 'gray-600', 'gray-400', 'blue-600', 'blue-500', 'green-600', 'red-600', 'amber-600', 'white']

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

function currentSize() {
  return currentFromSet(node.value!.classes, SIZE_CLASSES)?.slice('text-'.length) ?? 'base'
}
function currentWeight() {
  return currentFromSet(node.value!.classes, WEIGHT_CLASSES)?.slice('font-'.length) ?? 'normal'
}
function currentFamily(): string {
  const g = currentFromSet(node.value!.classes, GFONT_CLASSES)
  if (g) return g // e.g. "gfont-inter"
  const f = currentFromSet(node.value!.classes, FAMILY_CLASSES) ?? 'font-sans'
  return f
}
function currentAlign() {
  return currentFromSet(node.value!.classes, ALIGN_CLASSES)?.slice('text-'.length) ?? 'left'
}
function currentTransform() {
  return currentFromSet(node.value!.classes, TRANSFORM_CLASSES) ?? 'normal-case'
}
function currentColorClass() {
  const arb = currentArbitrary(node.value!.classes, 'text')
  if (arb) return `text-[${arb}]`
  return node.value!.classes.find((c) => /^text-(gray|blue|green|red|amber|yellow|purple|pink|indigo|cyan|teal|emerald|orange|rose|slate|zinc|neutral|stone|white|black|sky|violet|fuchsia|lime)-?\d{0,3}$/.test(c)) ?? null
}
function currentHex(): string {
  return currentArbitrary(node.value!.classes, 'text') ?? '#000000'
}
</script>

<template>
  <InspectorSection title="Typography" :icon="IconTypography" :show="!!node && showTypography">

    <!-- Inline rows: label left | control right -->
    <div class="space-y-2">

      <div class="flex items-center justify-between gap-2">
        <span class="text-[11px] font-medium text-foreground/80">Font</span>
        <Select
          :model-value="currentFamily()"
          @update:model-value="(v) => cls([...FAMILY_CLASSES, ...GFONT_CLASSES], v === 'font-sans' ? null : String(v))"
        >
          <SelectTrigger class="h-8 w-32 px-2 text-xs"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="font-sans">System Sans</SelectItem>
            <SelectItem value="font-serif">System Serif</SelectItem>
            <SelectItem value="font-mono">System Mono</SelectItem>
            <SelectItem v-for="g in GOOGLE_FONTS" :key="g.value" :value="`gfont-${g.value}`">{{ g.label }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="flex items-center justify-between gap-2">
        <span class="text-[11px] font-medium text-foreground/80">Size</span>
        <Select
          :model-value="currentSize()"
          @update:model-value="(v) => cls(SIZE_CLASSES, v === 'base' ? null : `text-${String(v)}`)"
        >
          <SelectTrigger class="h-8 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in FONT_SIZES" :key="s" :value="s">{{ s }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="flex items-center justify-between gap-2">
        <span class="text-[11px] font-medium text-foreground/80">Weight</span>
        <Select
          :model-value="currentWeight()"
          @update:model-value="(v) => cls(WEIGHT_CLASSES, v === 'normal' ? null : `font-${String(v)}`)"
        >
          <SelectTrigger class="h-8 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="w in FONT_WEIGHTS" :key="w" :value="w">{{ w }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="flex items-center justify-between gap-2">
        <span class="text-[11px] font-medium text-foreground/80">Align</span>
        <Select
          :model-value="currentAlign()"
          @update:model-value="(v) => cls(ALIGN_CLASSES, v === 'left' ? null : `text-${String(v)}`)"
        >
          <SelectTrigger class="h-8 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="a in TEXT_ALIGNS" :key="a" :value="a">{{ a }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="flex items-center justify-between gap-2">
        <span class="text-[11px] font-medium text-foreground/80">Transform</span>
        <Select
          :model-value="currentTransform()"
          @update:model-value="(v) => cls(TRANSFORM_CLASSES, v === 'normal-case' ? null : String(v))"
        >
          <SelectTrigger class="h-8 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem value="normal-case">Normal</SelectItem>
            <SelectItem value="uppercase">UPPERCASE</SelectItem>
            <SelectItem value="lowercase">lowercase</SelectItem>
            <SelectItem value="capitalize">Capitalize</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <!-- Color: full-width row (picker + hex input) then preset swatches -->
      <div class="space-y-1">
        <div class="flex items-center justify-between gap-2">
          <span class="text-[11px] font-medium text-foreground/80">Color</span>
          <div class="flex items-center gap-2">
            <input
              type="color"
              :value="currentHex()"
              class="h-7 w-9 shrink-0 cursor-pointer rounded border border-input p-0.5"
              @input="(e) => cls(['text-['], `text-[${(e.target as HTMLInputElement).value}]`)"
            />
            <Input
              :model-value="currentColorClass() ?? ''"
              class="h-8 w-28 font-mono text-xs"
              placeholder="text-foreground"
              @update:model-value="(v) => cls(['text-[', ...TEXT_COLOR_PRESETS.map((c) => `text-${c}`)], String(v) || null)"
            />
          </div>
        </div>
        <div class="flex flex-wrap gap-1">
          <button
            v-for="c in TEXT_COLOR_PRESETS"
            :key="c"
            class="h-5 w-5 rounded-full border border-border"
            :class="`bg-${c}`"
            :title="c"
            @click="cls(['text-[', ...TEXT_COLOR_PRESETS.map((x) => `text-${x}`)], `text-${c}`)"
          />
        </div>
      </div>

    </div>
  </InspectorSection>
</template>
