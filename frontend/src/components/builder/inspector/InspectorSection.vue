<script setup lang="ts">
// Collapsible inspector section shell: chevron + icon + title trigger,
// border between sections. Slot = controls. default-open controls initial
// state; open/closed persists across node selection (Figma-style).
import { ref, watch } from 'vue'
import { IconChevronDown, IconChevronRight } from '@tabler/icons-vue'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'

const props = withDefaults(
  defineProps<{ title: string; icon?: any; defaultOpen?: boolean; show?: boolean }>(),
  { defaultOpen: true, show: true },
)

const STORAGE_PREFIX = 'builder:section:'
const stored = localStorage.getItem(STORAGE_PREFIX + props.title)
const open = ref(stored !== null ? stored === 'true' : props.defaultOpen)

watch(open, (v) => localStorage.setItem(STORAGE_PREFIX + props.title, String(v)))
</script>

<template>
  <Collapsible
    v-if="show"
    v-model:open="open"
    class="border-b border-neutral-100"
  >
    <CollapsibleTrigger
      class="flex w-full items-center gap-1.5 px-3 py-1.5 text-left text-[11px] font-medium uppercase tracking-wider text-neutral-500 hover:text-neutral-700"
    >
      <component :is="open ? IconChevronDown : IconChevronRight" class="size-3 shrink-0" />
      <component :is="icon" v-if="icon" class="size-3.5 shrink-0" />
      <span>{{ title }}</span>
    </CollapsibleTrigger>
    <CollapsibleContent class="space-y-1.5 px-3 pb-2.5 pt-0.5">
      <slot />
    </CollapsibleContent>
  </Collapsible>
</template>
