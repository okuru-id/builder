<script setup lang="ts">
// Left outline panel: tree view + add-node palette.
import { computed, inject } from 'vue'
import {
  IconSquare,
  IconSection,
  IconTypography,
  IconH1,
  IconPhoto,
  IconCircuitPushbutton,
  IconLink,
  IconLayersIntersect,
} from '@tabler/icons-vue'
import { PALETTE_TYPES } from '@/types/page-builder'
import type { NodeType } from '@/types/page-builder'
import { BUILDER_KEY } from '@/components/builder/injection'
import TreeRow from '@/components/builder/TreeRow.vue'
import ComponentPalette from '@/components/builder/ComponentPalette.vue'

const store = inject(BUILDER_KEY, null)!

const paletteIcons: Record<NodeType, any> = {
  frame: IconSquare,
  section: IconSection,
  text: IconTypography,
  heading: IconH1,
  image: IconPhoto,
  button: IconCircuitPushbutton,
  link: IconLink,
  component: IconSquare,
}

const paletteLabels: Record<string, string> = {
  frame: 'Frame',
  section: 'Section',
  text: 'Text',
  heading: 'Heading',
  image: 'Image',
  button: 'Button',
  link: 'Link',
}

const root = computed(() => store.tree.value.root)

// Count total nodes for header info.
function countNodes(n: { children: any[] }): number {
  let c = 1
  for (const child of n.children) c += countNodes(child)
  return c
}
const nodeCount = computed(() => countNodes(root.value))

function addType(t: NodeType) {
  store.addNode(t, store.selectedId.value ?? store.tree.value.root.id)
}
</script>

<template>
  <aside class="flex h-full w-full flex-col border-r border-neutral-200 bg-white">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-neutral-200 px-3 py-2.5">
      <div class="flex items-center gap-1.5">
        <IconLayersIntersect class="size-4 text-neutral-500" />
        <h2 class="text-sm font-semibold">Layer</h2>
      </div>
      <span class="rounded-full bg-neutral-100 px-1.5 py-0.5 text-[10px] text-neutral-400">
        {{ nodeCount }}
      </span>
    </div>

    <!-- Tree -->
    <div class="flex-1 overflow-auto py-0.5">
      <TreeRow :node="root" />
    </div>

    <!-- Component palette -->
    <div class="max-h-[35%] overflow-auto border-t border-neutral-200">
      <ComponentPalette />
    </div>

    <!-- Add node palette -->
    <div class="border-t border-neutral-200 p-2">
      <div class="mb-1.5 text-[10px] font-medium uppercase tracking-wider text-neutral-400">Add Node</div>
      <div class="grid grid-cols-4 gap-1">
        <button
          v-for="t in PALETTE_TYPES"
          :key="t"
          class="flex flex-col items-center gap-0.5 rounded-md border border-neutral-200 py-1.5 text-[10px] text-neutral-600 transition-colors hover:border-neutral-300 hover:bg-neutral-50 active:bg-neutral-100"
          :title="`Add ${paletteLabels[t] ?? t}`"
          @click="addType(t)"
        >
          <component :is="paletteIcons[t]" class="size-3.5" />
          {{ paletteLabels[t] ?? t }}
        </button>
      </div>
    </div>
  </aside>
</template>
