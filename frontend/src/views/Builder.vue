<script setup lang="ts">
// Builder shell: top-level fullscreen route. Wires store → provide → 3-pane layout.
import { onMounted, onUnmounted, provide, watch } from 'vue'
import { useRoute } from 'vue-router'
import Toolbar from '@/components/builder/Toolbar.vue'
import Canvas from '@/components/builder/Canvas.vue'
import InspectorPanel from '@/components/builder/InspectorPanel.vue'
import NodeTreePanel from '@/components/builder/NodeTreePanel.vue'
import { useBuilderStore } from '@/components/builder/useBuilderStore'
import { BUILDER_KEY } from '@/components/builder/injection'

const route = useRoute()
const store = useBuilderStore()
provide(BUILDER_KEY, store)

onMounted(() => {
  const id = route.params.id as string
  store.load(id)
})

// Reload when navigating between builder pages without unmount.
watch(
  () => route.params.id,
  (id) => {
    if (id) store.load(id as string)
  },
)

function onKey(e: KeyboardEvent) {
  // Skip when typing in a form field or editing text.
  const target = e.target as HTMLElement
  if (
    target &&
    (target.isContentEditable ||
      ['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName))
  ) {
    return
  }
  const id = store.selectedId.value
  if (!id || id === store.tree.value.root.id) return

  if (e.key === 'Delete' || e.key === 'Backspace') {
    e.preventDefault()
    store.removeNode(id)
  } else if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'd') {
    e.preventDefault()
    store.duplicateNode(id)
  } else if ((e.metaKey || e.ctrlKey) && (e.key === 'ArrowUp' || e.key === 'ArrowDown')) {
    e.preventDefault()
    store.moveSiblingNode(id, e.key === 'ArrowUp' ? -1 : 1)
  }
}

onMounted(() => window.addEventListener('keydown', onKey))
onUnmounted(() => window.removeEventListener('keydown', onKey))
</script>

<template>
  <div v-if="store.loading.value" class="flex h-screen items-center justify-center bg-white">
    <div class="text-sm text-neutral-400">Memuat builder…</div>
  </div>
  <div v-else class="flex h-screen flex-col bg-neutral-50">
    <Toolbar />
    <div class="flex min-h-0 flex-1">
      <NodeTreePanel />
      <Canvas />
      <InspectorPanel />
    </div>
  </div>
</template>
