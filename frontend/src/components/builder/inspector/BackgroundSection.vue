<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import { replaceClass, currentFromSet, currentArbitrary } from '@/types/tokens'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import InspectorSection from './InspectorSection.vue'
import { IconColorSwatch } from '@tabler/icons-vue'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

const BG_SIZE_CLASSES = ['bg-cover', 'bg-contain'] as const
const BG_REPEAT_CLASSES = ['bg-repeat', 'bg-no-repeat', 'bg-repeat-x', 'bg-repeat-y'] as const
const BG_POS_CLASSES = ['bg-center', 'bg-top', 'bg-bottom', 'bg-left', 'bg-right', 'bg-left-top', 'bg-right-top', 'bg-left-bottom', 'bg-right-bottom'] as const
const ALL_BG_ASPECT_CLASSES = [...BG_SIZE_CLASSES, ...BG_REPEAT_CLASSES, ...BG_POS_CLASSES]

// Curated swatch palette spanning neutral + brand families. Rendered as a
// 2-row wrap grid; arbitrary hex still available via the color picker above.
const BG_COLOR_PRESETS = [
  'white', 'black',
  'slate-50', 'slate-100', 'slate-200', 'slate-800', 'slate-900',
  'gray-50', 'gray-100', 'gray-200', 'gray-500', 'gray-800', 'gray-900',
  'red-50', 'red-100', 'red-500', 'red-600',
  'orange-50', 'orange-100', 'orange-500',
  'amber-50', 'amber-100', 'amber-500',
  'yellow-50', 'yellow-100',
  'lime-100', 'lime-500',
  'green-50', 'green-100', 'green-500', 'green-600',
  'emerald-50', 'emerald-100', 'emerald-500',
  'teal-50', 'teal-100', 'teal-500',
  'cyan-50', 'cyan-100', 'cyan-500',
  'sky-50', 'sky-100', 'sky-500',
  'blue-50', 'blue-100', 'blue-500', 'blue-600',
  'indigo-50', 'indigo-100', 'indigo-500',
  'violet-50', 'violet-100', 'violet-500',
  'purple-50', 'purple-100', 'purple-500',
  'fuchsia-50', 'fuchsia-100', 'fuchsia-500',
  'pink-50', 'pink-100', 'pink-500',
  'rose-50', 'rose-100', 'rose-500',
]

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

// Raw bg color class shown in the text input: arbitrary hex OR preset color.
function currentBgClass() {
  const arb = currentArbitrary(node.value!.classes, 'bg')
  if (arb) return `bg-[${arb}]`
  return node.value!.classes.find((c) => /^bg-(gray|blue|green|red|amber|yellow|purple|pink|indigo|cyan|teal|emerald|orange|rose|slate|zinc|neutral|stone|white|black|sky|violet|fuchsia|lime)-?\d{0,3}$/.test(c)) ?? null
}
function currentBgHex() {
  return currentArbitrary(node.value!.classes, 'bg') ?? '#ffffff'
}
</script>

<template>
  <InspectorSection title="Background" :icon="IconColorSwatch" :show="!!node">

    <div class="space-y-1.5">
      <Label class="text-[11px] text-muted-foreground">Background Color</Label>
      <div class="flex items-center gap-2">
        <input
          type="color"
          :value="currentBgHex()"
          class="h-7 w-10 shrink-0 cursor-pointer rounded border border-input p-0.5"
          @input="(e) => cls(['bg-[', ...BG_COLOR_PRESETS.map((c) => `bg-${c}`)], `bg-[${(e.target as HTMLInputElement).value}]`)"
        />
        <Input
          :model-value="currentBgClass() ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="bg-background"
          @update:model-value="(v) => cls(['bg-[', ...BG_COLOR_PRESETS.map((c) => `bg-${c}`)], String(v) || null)"
        />
      </div>
      <div class="flex flex-wrap gap-1">
        <button
          v-for="c in BG_COLOR_PRESETS"
          :key="c"
          class="h-5 w-5 rounded-full border border-border"
          :class="`bg-${c}`"
          :title="c"
          @click="cls(['bg-[', ...BG_COLOR_PRESETS.map((x) => `bg-${x}`)], `bg-${c}`)"
        />
      </div>
    </div>

    <div class="flex items-center justify-between gap-2">
      <span class="text-[11px] text-muted-foreground">Size</span>
      <Select
        :model-value="currentFromSet(node?.classes ?? [], BG_SIZE_CLASSES) ?? 'auto'"
        @update:model-value="(v) => cls([...BG_SIZE_CLASSES], v === 'auto' ? null : String(v))"
      >
        <SelectTrigger class="h-7 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem value="auto">auto</SelectItem>
          <SelectItem v-for="s in BG_SIZE_CLASSES" :key="s" :value="s">{{ s.replace('bg-', '') }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="flex items-center justify-between gap-2">
      <span class="text-[11px] text-muted-foreground">Repeat</span>
      <Select
        :model-value="currentFromSet(node?.classes ?? [], BG_REPEAT_CLASSES) ?? 'repeat'"
        @update:model-value="(v) => cls([...BG_REPEAT_CLASSES], v === 'repeat' ? null : String(v))"
      >
        <SelectTrigger class="h-7 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem value="repeat">repeat</SelectItem>
          <SelectItem v-for="r in BG_REPEAT_CLASSES" :key="r" :value="r">{{ r.replace('bg-', '') }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="flex items-center justify-between gap-2">
      <span class="text-[11px] text-muted-foreground">Position</span>
      <Select
        :model-value="currentFromSet(node?.classes ?? [], BG_POS_CLASSES) ?? 'center'"
        @update:model-value="(v) => cls([...BG_POS_CLASSES], v === 'center' ? null : String(v))"
      >
        <SelectTrigger class="h-7 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="p in BG_POS_CLASSES" :key="p" :value="p">{{ p.replace('bg-', '') }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="space-y-1.5">
      <Label class="text-[11px] text-muted-foreground">Image URL</Label>
      <Input
        :model-value="''"
        class="h-8 font-mono text-xs"
        placeholder="url(…) or bg-image class"
        @update:model-value="(v) => cls(ALL_BG_ASPECT_CLASSES, String(v) || null)"
      />
      <p class="text-[11px] text-muted-foreground">Size/position/repeat classes are set above.</p>
    </div>
  </InspectorSection>
</template>
