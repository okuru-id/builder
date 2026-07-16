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
  FLEX_WRAPS,
  SPACING,
} from '@/types/tokens'
import InspectorSection from './InspectorSection.vue'
import {
  IconLayoutNavbar,
  IconArrowRight,
  IconArrowDown,
  IconArrowBackUp,
  IconArrowUp,
} from '@tabler/icons-vue'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

const DIR_CLASSES = FLEX_DIRECTIONS.map((d) => `flex-${d}`)
const ITEMS_CLASSES = ALIGN_ITEMS.map((a) => `items-${a}`)
const JUSTIFY_CLASSES = JUSTIFY_CONTENTS.map((j) => `justify-${j}`)
const WRAP_CLASSES = FLEX_WRAPS.map((w) => `flex-${w}`)

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
function currentJustify() {
  return currentFromSet(node.value!.classes, JUSTIFY_CLASSES)?.slice('justify-'.length) ?? 'start'
}
function currentWrap() {
  return currentFromSet(node.value!.classes, WRAP_CLASSES)?.slice('flex-'.length) ?? 'nowrap'
}

// Direction: 4 icon buttons (row/col/row-reverse/col-reverse).
const directions = [
  { value: 'row', label: 'Row', icon: IconArrowRight },
  { value: 'col', label: 'Column', icon: IconArrowDown },
  { value: 'row-reverse', label: 'Row reversed', icon: IconArrowBackUp },
  { value: 'col-reverse', label: 'Column reversed', icon: IconArrowUp },
]
</script>

<template>
  <InspectorSection title="Layout" :icon="IconLayoutNavbar" :show="!!node">
    <!-- Flex on/off -->
    <div class="flex items-center justify-between">
      <span class="text-[11px] font-medium text-foreground/80">Flex</span>
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
          class="flex h-8 items-center justify-center rounded-md border transition-colors"
          :class="currentDir() === d.value
            ? 'border-primary bg-primary/10 text-primary'
            : 'border-border text-muted-foreground hover:bg-muted'"
          :title="d.label"
          @click="cls(DIR_CLASSES, d.value === 'row' ? null : `flex-${d.value}`)"
        >
          <component :is="d.icon" class="size-4" />
        </button>
      </div>

      <!-- Horizontal alignment (justify) — full scale -->
      <div class="flex items-center justify-between gap-2">
        <span class="text-[11px] font-medium text-foreground/80">Justify</span>
        <Select
          :model-value="currentJustify()"
          @update:model-value="(v) => cls(JUSTIFY_CLASSES, v === 'start' ? null : `justify-${String(v)}`)"
        >
          <SelectTrigger class="h-8 w-28 px-2 text-xs capitalize"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="j in JUSTIFY_CONTENTS" :key="j" :value="j">{{ j }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <!-- Vertical alignment (items) — full scale -->
      <div class="flex items-center justify-between gap-2">
        <span class="text-[11px] font-medium text-foreground/80">Align Items</span>
        <Select
          :model-value="currentItems()"
          @update:model-value="(v) => cls(ITEMS_CLASSES, v === 'stretch' ? null : `items-${String(v)}`)"
        >
          <SelectTrigger class="h-8 w-28 px-2 text-xs capitalize"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="a in ALIGN_ITEMS" :key="a" :value="a">{{ a }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <!-- Wrap — 3-way -->
      <div class="flex items-center justify-between gap-2">
        <span class="text-[11px] font-medium text-foreground/80">Wrap</span>
        <Select
          :model-value="currentWrap()"
          @update:model-value="(v) => cls(WRAP_CLASSES, v === 'nowrap' ? null : `flex-${String(v)}`)"
        >
          <SelectTrigger class="h-8 w-28 px-2 text-xs capitalize"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="w in FLEX_WRAPS" :key="w" :value="w">{{ w }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <!-- Gap -->
      <div class="flex items-center justify-between gap-2">
        <span class="text-[11px] font-medium text-foreground/80">Gap</span>
        <Select
          :model-value="currentClass(node?.classes ?? [], 'gap') ?? '0'"
          @update:model-value="(v) => cls(['gap'], v === '0' ? null : `gap-${String(v)}`)"
        >
          <SelectTrigger class="h-8 w-28 px-2 text-xs"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in SPACING" :key="s" :value="s">{{ s }}</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </template>
  </InspectorSection>
</template>
