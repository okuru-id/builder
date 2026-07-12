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
  IconGridDots,
  IconMinus,
  IconStar,
  IconForms,
  IconInputSearch,
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
  grid: IconGridDots,
  text: IconTypography,
  heading: IconH1,
  image: IconPhoto,
  button: IconCircuitPushbutton,
  link: IconLink,
  divider: IconMinus,
  icon: IconStar,
  form: IconForms,
  input: IconInputSearch,
  component: IconSquare,
}

const paletteLabels: Record<string, string> = {
  frame: 'Frame',
  section: 'Section',
  grid: 'Grid',
  text: 'Text',
  heading: 'Heading',
  image: 'Image',
  button: 'Button',
  link: 'Link',
  divider: 'Divider',
  icon: 'Icon',
  form: 'Form',
  input: 'Input',
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
  <aside class="flex h-full w-full flex-col border-r border-border bg-background">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-border px-3 py-2.5">
      <div class="flex items-center gap-1.5">
        <IconLayersIntersect class="size-4 text-muted-foreground" />
        <h2 class="text-sm font-semibold">Layer</h2>
      </div>
      <span class="rounded-full bg-muted px-1.5 py-0.5 text-[10px] text-muted-foreground">
        {{ nodeCount }}
      </span>
    </div>

    <!-- Tree -->
    <div class="flex-1 overflow-auto py-0.5">
      <TreeRow :node="root" />
    </div>

    <!-- Component palette -->
    <div class="max-h-[35%] overflow-auto border-t border-border">
      <ComponentPalette />
    </div>

    <!-- Add node palette -->
    <div class="border-t border-border p-2">
      <div class="mb-1.5 text-[10px] font-medium uppercase tracking-wider text-muted-foreground">Add Node</div>
      <div class="grid grid-cols-4 gap-1">
        <button
          v-for="t in PALETTE_TYPES"
          :key="t"
          class="flex flex-col items-center gap-0.5 rounded-md border border-border py-1.5 text-[10px] text-muted-foreground transition-colors hover:border-input hover:bg-muted/50 active:bg-muted"
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
