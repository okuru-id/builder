<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import {
  replaceClass,
  currentFromSet,
  currentArbitrary,
  FONT_SIZES,
  FONT_WEIGHTS,
  TEXT_ALIGNS,
} from '@/types/tokens'
import InspectorSection from './InspectorSection.vue'
import { IconTypography } from '@tabler/icons-vue'
import { Label } from '@/components/ui/label'
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
const ALIGN_CLASSES = TEXT_ALIGNS.map((a) => `text-${a}`)
const TRANSFORM_CLASSES: string[] = ['uppercase', 'lowercase', 'capitalize', 'normal-case']
const TEXT_COLOR_PRESETS = ['gray-900', 'gray-600', 'gray-400', 'blue-600', 'blue-500', 'green-600', 'red-600', 'amber-600', 'white']

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

function currentSize() {
  // Strip arbitrary color first so text-[#hex] doesn't shadow a size match.
  return currentFromSet(node.value!.classes, SIZE_CLASSES)?.slice('text-'.length) ?? 'base'
}
function currentWeight() {
  return currentFromSet(node.value!.classes, WEIGHT_CLASSES)?.slice('font-'.length) ?? 'normal'
}
function currentAlign() {
  return currentFromSet(node.value!.classes, ALIGN_CLASSES)?.slice('text-'.length) ?? 'left'
}
function currentTransform() {
  return currentFromSet(node.value!.classes, TRANSFORM_CLASSES) ?? 'normal-case'
}
// Raw color class for the text field (preset or arbitrary).
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

    <div class="space-y-1.5">
      <Label class="text-[11px] text-neutral-400">Size</Label>
      <Select
        :model-value="currentSize()"
        @update:model-value="(v) => cls(SIZE_CLASSES, v === 'base' ? null : `text-${String(v)}`)"
      >
        <SelectTrigger class="h-7 text-xs"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="s in FONT_SIZES" :key="s" :value="s">{{ s }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="space-y-1.5">
      <Label class="text-[11px] text-neutral-400">Weight</Label>
      <Select
        :model-value="currentWeight()"
        @update:model-value="(v) => cls(WEIGHT_CLASSES, v === 'normal' ? null : `font-${String(v)}`)"
      >
        <SelectTrigger class="h-7 text-xs"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="w in FONT_WEIGHTS" :key="w" :value="w">{{ w }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="space-y-1.5">
      <Label class="text-[11px] text-neutral-400">Align</Label>
      <Select
        :model-value="currentAlign()"
        @update:model-value="(v) => cls(ALIGN_CLASSES, v === 'left' ? null : `text-${String(v)}`)"
      >
        <SelectTrigger class="h-7 text-xs"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="a in TEXT_ALIGNS" :key="a" :value="a">{{ a }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="space-y-1.5">
      <Label class="text-[11px] text-neutral-400">Text Color</Label>
      <div class="flex items-center gap-2">
        <input
          type="color"
          :value="currentHex()"
          class="h-7 w-10 shrink-0 cursor-pointer rounded border border-neutral-300 p-0.5"
          @input="(e) => cls(['text-['], `text-[${(e.target as HTMLInputElement).value}]`)"
        />
        <Input
          :model-value="currentColorClass() ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="text-gray-900"
          @update:model-value="(v) => cls(['text-[', ...TEXT_COLOR_PRESETS.map((c) => `text-${c}`)], String(v) || null)"
        />
      </div>
      <div class="flex flex-wrap gap-1">
        <button
          v-for="c in TEXT_COLOR_PRESETS"
          :key="c"
          class="h-5 w-5 rounded-full border border-neutral-200"
          :class="`bg-${c}`"
          :title="c"
          @click="cls(['text-[', ...TEXT_COLOR_PRESETS.map((x) => `text-${x}`)], `text-${c}`)"
        />
      </div>
    </div>

    <div class="space-y-1.5">
      <Label class="text-[11px] text-neutral-400">Transform</Label>
      <Select
        :model-value="currentTransform()"
        @update:model-value="(v) => cls(TRANSFORM_CLASSES, v === 'normal-case' ? null : String(v))"
      >
        <SelectTrigger class="h-7 text-xs"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem value="normal-case">Normal</SelectItem>
          <SelectItem value="uppercase">UPPERCASE</SelectItem>
          <SelectItem value="lowercase">lowercase</SelectItem>
          <SelectItem value="capitalize">Capitalize</SelectItem>
        </SelectContent>
      </Select>
    </div>
  </InspectorSection>
</template>
