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

// Color swatches grouped by family, ordered light→dark within each group.
const BG_COLOR_GROUPS: { label: string; colors: string[] }[] = [
  { label: 'Neutral', colors: [
    'white',
    'slate-50', 'gray-50',
    'slate-100', 'gray-100',
    'slate-200', 'gray-200',
    'gray-500',
    'slate-800', 'gray-800',
    'slate-900', 'gray-900',
    'black',
  ]},
  { label: 'Red / Orange / Amber', colors: [
    'red-50', 'orange-50', 'amber-50',
    'red-100', 'orange-100', 'amber-100',
    'red-500', 'orange-500', 'amber-500',
    'red-600',
  ]},
  { label: 'Yellow / Lime / Green', colors: [
    'yellow-50', 'lime-100', 'green-50', 'emerald-50',
    'yellow-100', 'green-100', 'emerald-100',
    'lime-500', 'green-500', 'emerald-500',
    'green-600',
  ]},
  { label: 'Teal / Cyan / Sky', colors: [
    'teal-50', 'cyan-50', 'sky-50',
    'teal-100', 'cyan-100', 'sky-100',
    'teal-500', 'cyan-500', 'sky-500',
  ]},
  { label: 'Blue / Indigo / Violet', colors: [
    'blue-50', 'indigo-50', 'violet-50',
    'blue-100', 'indigo-100', 'violet-100',
    'blue-500', 'indigo-500', 'violet-500',
    'blue-600',
  ]},
  { label: 'Purple / Pink / Rose', colors: [
    'purple-50', 'fuchsia-50', 'pink-50', 'rose-50',
    'purple-100', 'fuchsia-100', 'pink-100', 'rose-100',
    'purple-500', 'fuchsia-500', 'pink-500', 'rose-500',
  ]},
]

const BG_COLOR_PRESETS = BG_COLOR_GROUPS.flatMap((g) => g.colors)

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

    <div class="space-y-1">
      <Label class="text-[11px] font-medium text-foreground/80">Background Color</Label>
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
      <div class="space-y-2">
        <div v-for="grp in BG_COLOR_GROUPS" :key="grp.label" class="space-y-0.5">
          <span class="text-[10px] font-medium uppercase tracking-wide text-muted-foreground/60">{{ grp.label }}</span>
          <div class="flex flex-wrap gap-0.5">
            <button
              v-for="c in grp.colors"
              :key="c"
              class="h-5 w-5 rounded-full border border-border"
              :class="`bg-${c}`"
              :title="c"
              @click="cls(['bg-[', ...BG_COLOR_PRESETS.map((x) => `bg-${x}`)], `bg-${c}`)"
            />
          </div>
        </div>
      </div>
    </div>

    <div class="flex items-center justify-between gap-2">
      <span class="text-[11px] font-medium text-foreground/80">Size</span>
      <Select
        :model-value="currentFromSet(node?.classes ?? [], BG_SIZE_CLASSES) ?? 'auto'"
        @update:model-value="(v) => cls([...BG_SIZE_CLASSES], v === 'auto' ? null : String(v))"
      >
        <SelectTrigger class="h-8 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem value="auto">auto</SelectItem>
          <SelectItem v-for="s in BG_SIZE_CLASSES" :key="s" :value="s">{{ s.replace('bg-', '') }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="flex items-center justify-between gap-2">
      <span class="text-[11px] font-medium text-foreground/80">Repeat</span>
      <Select
        :model-value="currentFromSet(node?.classes ?? [], BG_REPEAT_CLASSES) ?? 'repeat'"
        @update:model-value="(v) => cls([...BG_REPEAT_CLASSES], v === 'repeat' ? null : String(v))"
      >
        <SelectTrigger class="h-8 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem value="repeat">repeat</SelectItem>
          <SelectItem v-for="r in BG_REPEAT_CLASSES" :key="r" :value="r">{{ r.replace('bg-', '') }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="flex items-center justify-between gap-2">
      <span class="text-[11px] font-medium text-foreground/80">Position</span>
      <Select
        :model-value="currentFromSet(node?.classes ?? [], BG_POS_CLASSES) ?? 'bg-center'"
        @update:model-value="(v) => cls([...BG_POS_CLASSES], v === 'bg-center' ? null : String(v))"
      >
        <SelectTrigger class="h-8 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="p in BG_POS_CLASSES" :key="p" :value="p">{{ p.replace('bg-', '') }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="space-y-1">
      <Label class="text-[11px] font-medium text-foreground/80">Image URL</Label>
      <Input
        :model-value="node?.classes.find((c) => c.startsWith('bg-[url(')) ?? ''"
        class="h-8 font-mono text-xs"
        placeholder="bg-[url('/img.jpg')]"
        @update:model-value="(v) => cls(['bg-[url('], String(v) || null)"
      />
      <p class="text-[11px] font-medium text-foreground/80">Size/position/repeat classes are set above.</p>
    </div>
  </InspectorSection>
</template>
