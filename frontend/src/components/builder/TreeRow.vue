<script setup lang="ts">
// Recursive outline row. Self-references by component name for nested children.
// Supports drag-and-drop reorder using the same HTML5 DnD handlers as the canvas.
import { inject } from 'vue'
import type { Node } from '@/types/page-builder'
import { BUILDER_KEY } from '@/components/builder/injection'

defineOptions({ name: 'TreeRow' })
const props = withDefaults(
  defineProps<{ node: Node; depth?: number }>(),
  { depth: 0 },
)

const store = inject(BUILDER_KEY, null)!
const selected = () => store.selectedId.value === props.node.id

const CONTAINER_TYPES = new Set(['frame', 'section', 'link', 'component'])
const isContainer = () => CONTAINER_TYPES.has(props.node.type)
const dragging = () => store.draggingId.value === props.node.id
const dropInside = () => {
  const t = store.dropTarget.value
  return !!t && t.pos === 'inside' && t.parentId === props.node.id
}

function onDragStart(e: DragEvent) {
  store.dragStart(props.node.id)
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', props.node.id)
  }
}
function onDragOver(e: DragEvent) {
  store.dragOver(props.node.id, e, isContainer())
}
function onDrop(e: DragEvent) {
  e.preventDefault()
  e.stopPropagation()
  store.drop()
}
function onDragEnd() {
  store.dragEnd()
}
</script>

<template>
  <div
    :draggable="true"
    @dragstart="onDragStart"
    @dragover="onDragOver"
    @drop="onDrop"
    @dragend="onDragEnd"
  >
    <button
      class="flex w-full items-center gap-1.5 py-1 pr-2 text-left text-xs"
      :class="[
        selected() ? 'bg-blue-50 text-blue-700' : 'hover:bg-neutral-100',
        dragging() ? 'opacity-40' : '',
        dropInside() ? 'ring-2 ring-blue-500 ring-inset' : '',
      ]"
      :style="{ paddingLeft: `${depth * 12 + 8}px` }"
      @click.stop="store.select(node.id)"
    >
      <span class="size-1.5 shrink-0 rounded-full"
        :class="selected() ? 'bg-blue-500' : 'bg-neutral-300'"
      />
      <span class="truncate">{{ node.name || node.type }}</span>
    </button>
    <TreeRow
      v-for="child in node.children"
      :key="child.id"
      :node="child"
      :depth="depth + 1"
    />
  </div>
</template>
