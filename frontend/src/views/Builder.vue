<script setup lang="ts">
// Builder shell: top-level fullscreen route. Wires store → provide → 3-pane layout.
import { onMounted, onUnmounted, provide, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Toolbar from '@/components/builder/Toolbar.vue'
import Canvas from '@/components/builder/Canvas.vue'
import RightPanel from '@/components/builder/RightPanel.vue'
import NodeTreePanel from '@/components/builder/NodeTreePanel.vue'
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from '@/components/ui/resizable'
import { useBuilderStore } from '@/components/builder/useBuilderStore'
import { BUILDER_KEY } from '@/components/builder/injection'

const route = useRoute()
const store = useBuilderStore()
provide(BUILDER_KEY, store)

// Side-panel visibility (UI-only, not persisted). ponytail: no localStorage.
const showLeft = ref(true)
const showRight = ref(true)

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
    // Allow Ctrl+S even from inputs.
    if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 's') {
      e.preventDefault()
      store.save()
    }
    return
  }

  // Ctrl+S — save
  if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 's') {
    e.preventDefault()
    store.save()
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

// Warn before closing/navigating away with unsaved changes.
function onBeforeUnload(e: BeforeUnloadEvent) {
  if (store.dirty.value) {
    e.preventDefault()
  }
}

onMounted(() => {
  window.addEventListener('keydown', onKey)
  window.addEventListener('beforeunload', onBeforeUnload)
})
onUnmounted(() => {
  window.removeEventListener('keydown', onKey)
  window.removeEventListener('beforeunload', onBeforeUnload)
})
</script>

<template>
  <div v-if="store.loading.value" class="flex h-screen items-center justify-center bg-background">
    <div class="text-sm text-muted-foreground">Loading builder…</div>
  </div>
  <div v-else class="flex h-screen flex-col bg-muted">
    <Toolbar :show-left="showLeft" :show-right="showRight" @toggle-left="showLeft = !showLeft" @toggle-right="showRight = !showRight" />
    <ResizablePanelGroup
      direction="horizontal"
      class="flex min-h-0 flex-1"
      auto-save-id="builder-layout"
    >
      <ResizablePanel v-if="showLeft" :default-size="18" :min-size="12" :max-size="28">
        <NodeTreePanel />
      </ResizablePanel>
      <ResizableHandle v-if="showLeft" />
      <ResizablePanel :default-size="60" :min-size="30">
        <Canvas />
      </ResizablePanel>
      <ResizableHandle v-if="showRight" />
      <ResizablePanel v-if="showRight" :default-size="22" :min-size="16" :max-size="42">
        <RightPanel />
      </ResizablePanel>
    </ResizablePanelGroup>
  </div>
</template>
