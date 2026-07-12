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
} from '@tabler/icons-vue'
import { PALETTE_TYPES } from '@/types/page-builder'
import type { NodeType } from '@/types/page-builder'
import { BUILDER_KEY } from '@/components/builder/injection'
import TreeRow from '@/components/builder/TreeRow.vue'

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

const root = computed(() => store.tree.value.root)

function addType(t: NodeType) {
  store.addNode(t, store.tree.value.root.id)
}
</script>

<template>
  <aside class="flex w-60 flex-col border-r border-neutral-200 bg-white">
    <div class="border-b border-neutral-200 px-3 py-3">
      <h2 class="text-sm font-semibold">Layer</h2>
    </div>

    <div class="flex-1 overflow-auto py-1">
      <TreeRow :node="root" />
    </div>

    <div class="border-t border-neutral-200 p-2">
      <div class="mb-1.5 text-xs font-medium text-neutral-500">Tambah node</div>
      <div class="grid grid-cols-4 gap-1">
        <button
          v-for="t in PALETTE_TYPES"
          :key="t"
          class="flex flex-col items-center gap-1 rounded-md border border-neutral-200 py-2 text-[10px] text-neutral-600 hover:bg-neutral-50"
          :title="`Tambah ${t}`"
          @click="addType(t)"
        >
          <component :is="paletteIcons[t]" class="size-4" />
          {{ t }}
        </button>
      </div>
    </div>
  </aside>
</template>
