<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import {
  replaceClass,
  hasClass,
  currentFromSet,
  currentArbitrary,
  BORDER_WIDTHS,
  BORDER_RADII,
} from '@/types/tokens'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Input } from '@/components/ui/input'
import InspectorSection from './InspectorSection.vue'
import { IconBorderAll } from '@tabler/icons-vue'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

const WIDTH_CLASSES = BORDER_WIDTHS.map((w) => (w === '1' ? 'border' : `border-${w}`))
const RADIUS_CLASSES = BORDER_RADII.map((r) => {
  if (r === 'none') return 'rounded-none'
  if (r === 'md') return 'rounded-md'
  if (r === 'lg') return 'rounded-lg'
  return `rounded-${r}`
})
const BORDER_COLOR_PRESETS = ['gray-300', 'gray-500', 'blue-500', 'red-500', 'green-500', 'amber-500', 'neutral-900', 'white', 'black']

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

function setBorderColor(e: Event) {
  const hex = (e.target as HTMLInputElement).value
  cls(['border-[', ...BORDER_COLOR_PRESETS.map((c) => `border-${c}`)], hex ? `border-[${hex}]` : null)
  if (hex && !hasClass(node.value!.classes, 'border')) {
    store.patchNode(node.value!.id, { classes: [...node.value!.classes, 'border'] })
  }
}

function currentWidth() {
  return currentFromSet(node.value!.classes, WIDTH_CLASSES) ?? '1'
}
function currentRadius() {
  return currentFromSet(node.value!.classes, RADIUS_CLASSES) ?? 'none'
}
function currentBorderClass() {
  const arb = currentArbitrary(node.value!.classes, 'border')
  if (arb) return `border-[${arb}]`
  return node.value!.classes.find((c) => /^border-(gray|blue|green|red|amber|neutral|slate|zinc|stone)-?\d{0,3}$/.test(c)) ?? null
}
function currentBorderHex() {
  return currentArbitrary(node.value!.classes, 'border') ?? '#000000'
}
</script>

<template>
  <InspectorSection title="Border" :icon="IconBorderAll" :show="!!node">

    <div class="flex items-center justify-between">
      <Label class="text-[11px] text-muted-foreground">Show Border</Label>
      <Switch
        :model-value="hasClass(node?.classes ?? [], 'border')"
        @update:model-value="(v) => cls(['border', ...WIDTH_CLASSES], v ? 'border' : null)"
      />
    </div>

    <template v-if="hasClass(node?.classes ?? [], 'border')">
      <div class="space-y-1.5">
        <Label class="text-[11px] text-muted-foreground">Width</Label>
        <Select
          :model-value="currentWidth()"
          @update:model-value="(v) => cls([...WIDTH_CLASSES], String(v) === '1' ? 'border' : `border-${String(v)}`)"
        >
          <SelectTrigger class="h-7 text-xs"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="w in BORDER_WIDTHS" :key="w" :value="w">{{ w === '1' ? 'default' : w }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="space-y-1.5">
        <Label class="text-[11px] text-muted-foreground">Border Color</Label>
        <div class="flex items-center gap-2">
          <input
            type="color"
            :value="currentBorderHex()"
            class="h-7 w-10 shrink-0 cursor-pointer rounded border border-input p-0.5"
            @input="setBorderColor"
          />
          <Input
            :model-value="currentBorderClass() ?? ''"
            class="h-8 flex-1 font-mono text-xs"
            placeholder="border-input"
            @update:model-value="(v) => cls(['border-[', ...BORDER_COLOR_PRESETS.map((c) => `border-${c}`)], String(v) || null)"
          />
        </div>
      </div>

      <div class="space-y-1.5">
        <Label class="text-[11px] text-muted-foreground">Radius</Label>
        <Select
          :model-value="currentRadius()"
          @update:model-value="(v) => cls([...RADIUS_CLASSES], v === 'none' ? null : (v === 'md' ? 'rounded-md' : `rounded-${String(v)}`))"
        >
          <SelectTrigger class="h-7 text-xs"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="r in BORDER_RADII" :key="r" :value="r">{{ r }}</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </template>
  </InspectorSection>
</template>
