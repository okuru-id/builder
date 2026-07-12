<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import {
  replaceClass,
  hasClass,
  currentClass,
  currentFromSet,
  DISPLAY_CLASSES,
  FLEX_DIRECTIONS,
  ALIGN_ITEMS,
  JUSTIFY_CONTENTS,
  SPACING,
} from '@/types/tokens'
import InspectorSection from './InspectorSection.vue'
import {
  IconLayoutNavbar,
  IconArrowRight,
  IconArrowDown,
  IconArrowBackUp,
  IconArrowUp,
  IconLayoutAlignLeft,
  IconLayoutAlignCenter,
  IconLayoutAlignRight,
  IconLayoutAlignTop,
  IconLayoutAlignMiddle,
  IconLayoutAlignBottom,
  IconTextWrap,
  IconTextWrapDisabled,
} from '@tabler/icons-vue'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

const DIR_CLASSES = FLEX_DIRECTIONS.map((d) => `flex-${d}`)
const ITEMS_CLASSES = ALIGN_ITEMS.map((a) => `items-${a}`)
const JUSTIFY_CLASSES = JUSTIFY_CONTENTS.map((j) => `justify-${j}`)

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

function currentDir() {
  return currentFromSet(node.value!.classes, DIR_CLASSES)?.slice('flex-'.length) ?? 'row'
}
function currentItems() {
  return currentFromSet(node.value!.classes, ITEMS_CLASSES)?.slice('items-'.length) ?? 'stretch'
}
// isItemsActive: tombol aktif kalau class ada dan match, ATAU tombol=start
// tapi tidak ada items-* class (default stretch-only UI).
// ponytail: Figma treats default as stretch; start button never highlights
// stretch fallback, so we keep it explicit.
function currentJustify() {
  return currentFromSet(node.value!.classes, JUSTIFY_CLASSES)?.slice('justify-'.length) ?? 'start'
}

// Figma-style direction toggle (4 icon buttons).
const directions = [
  { value: 'row', label: 'Horizontal', icon: IconArrowRight },
  { value: 'col', label: 'Vertical', icon: IconArrowDown },
  { value: 'row-reverse', label: 'Horizontal reversed', icon: IconArrowBackUp },
  { value: 'col-reverse', label: 'Vertical reversed', icon: IconArrowUp },
]

// Figma-style: two independent alignment rows. Horizontal = justify-*,
// Vertikal = items-*. Each row is a 3-icon toggle.
const hAlign = [
  { value: 'start', icon: IconLayoutAlignLeft, label: 'Left' },
  { value: 'center', icon: IconLayoutAlignCenter, label: 'Center' },
  { value: 'end', icon: IconLayoutAlignRight, label: 'Right' },
]
const vAlign = [
  { value: 'start', icon: IconLayoutAlignTop, label: 'Top' },
  { value: 'center', icon: IconLayoutAlignMiddle, label: 'Center' },
  { value: 'end', icon: IconLayoutAlignBottom, label: 'Bottom' },
]
</script>

<template>
  <InspectorSection title="Layout" :icon="IconLayoutNavbar" :show="!!node">
    <!-- Flex on/off -->
    <div class="flex items-center justify-between">
      <span class="text-[11px] text-neutral-500">Flex</span>
      <Switch
        :model-value="hasClass(node?.classes ?? [], 'flex')"
        @update:model-value="(v) => cls(['flex', ...DISPLAY_CLASSES], v ? 'flex' : null)"
      />
    </div>

    <template v-if="hasClass(node?.classes ?? [], 'flex')">
      <!-- Direction: 4 icon buttons -->
      <div class="grid grid-cols-4 gap-1">
        <button
          v-for="d in directions"
          :key="d.value"
          class="flex h-7 items-center justify-center rounded-md border transition-colors"
          :class="currentDir() === d.value
            ? 'border-primary bg-primary/10 text-primary'
            : 'border-neutral-200 text-neutral-500 hover:bg-neutral-100'"
          :title="d.label"
          @click="cls(DIR_CLASSES, d.value === 'row' ? null : `flex-${d.value}`)"
        >
          <component :is="d.icon" class="size-4" />
        </button>
      </div>

      <!-- Alignment: 2 independent rows (horizontal + vertical) -->
      <div class="space-y-1">
        <div class="grid grid-cols-3 gap-1">
          <button
            v-for="h in hAlign"
            :key="`h-${h.value}`"
            class="flex h-7 items-center justify-center rounded-md border transition-colors"
            :class="currentJustify() === h.value
              ? 'border-primary bg-primary/10 text-primary'
              : 'border-neutral-200 text-neutral-500 hover:bg-neutral-100'"
            :title="`Horizontal ${h.label} (justify-${h.value})`"
            @click="cls(JUSTIFY_CLASSES, h.value === 'start' ? null : `justify-${h.value}`)"
          >
            <component :is="h.icon" class="size-4" />
          </button>
        </div>
        <div class="grid grid-cols-3 gap-1">
          <button
            v-for="v in vAlign"
            :key="`v-${v.value}`"
            class="flex h-7 items-center justify-center rounded-md border transition-colors"
            :class="currentItems() === v.value
              ? 'border-primary bg-primary/10 text-primary'
              : 'border-neutral-200 text-neutral-500 hover:bg-neutral-100'"
            :title="`Vertical ${v.label} (items-${v.value})`"
            @click="cls(ITEMS_CLASSES, `items-${v.value}`)"
          >
            <component :is="v.icon" class="size-4" />
          </button>
        </div>
      </div>

      <!-- Gap + Wrap row -->
      <div class="flex items-center gap-2">
        <span class="w-8 shrink-0 text-[11px] text-neutral-500">Gap</span>
        <Select
          :model-value="currentClass(node?.classes ?? [], 'gap') ?? '0'"
          class="flex-1"
          @update:model-value="(v) => cls(['gap'], v === '0' ? null : `gap-${String(v)}`)"
        >
          <SelectTrigger class="h-7 px-2 text-xs"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in SPACING" :key="s" :value="s">{{ s }}</SelectItem>
          </SelectContent>
        </Select>
        <button
          class="flex size-7 shrink-0 items-center justify-center rounded-md border transition-colors"
          :class="hasClass(node?.classes ?? [], 'flex-wrap')
            ? 'border-primary bg-primary/10 text-primary'
            : 'border-neutral-200 text-neutral-500 hover:bg-neutral-100'"
          :title="hasClass(node?.classes ?? [], 'flex-wrap') ? 'Wrap on' : 'Wrap off'"
          @click="cls(['flex-wrap', 'flex-nowrap', 'flex-wrap-reverse'], hasClass(node?.classes ?? [], 'flex-wrap') ? null : 'flex-wrap')"
        >
          <component :is="hasClass(node?.classes ?? [], 'flex-wrap') ? IconTextWrap : IconTextWrapDisabled" class="size-4" />
        </button>
      </div>
    </template>
  </InspectorSection>
</template>
