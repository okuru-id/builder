<script setup lang="ts">
// Canvas viewport: renders the tree with interactive NodeRenderer. Clicking the
// blank background deselects. Breakpoint width is driven by the store.
import { computed, inject } from 'vue'
import NodeRenderer from '@/components/builder/NodeRenderer.vue'
import { BUILDER_KEY } from '@/components/builder/injection'

const store = inject(BUILDER_KEY, null)
if (!store) throw new Error('Canvas requires a BuilderStore provider')

const widthStyle = computed(() => {
  const w = store!.canvasWidth.value
  return w ? { maxWidth: `${w}px`, width: `${w}px` } : { width: '100%' }
})

function onBackgroundClick() {
  store!.select(null)
}
</script>

<template>
  <div
    class="flex h-full w-full justify-center overflow-auto bg-neutral-200 p-6 min-h-0"
    @click="onBackgroundClick"
  >
    <div
      :style="widthStyle"
      class="min-h-full bg-white shadow-lg transition-[max-width] duration-200"
    >
      <NodeRenderer :node="store.tree.value.root" />
    </div>
  </div>
</template>
