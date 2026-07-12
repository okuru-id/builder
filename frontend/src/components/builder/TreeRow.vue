<script setup lang="ts">
// Recursive outline row with expand/collapse, type icons, and drag-and-drop.
import { computed, inject, ref } from 'vue'
import {
  IconSquare,
  IconSection,
  IconTypography,
  IconH1,
  IconPhoto,
  IconCircuitPushbutton,
  IconLink,
  IconComponents,
  IconChevronRight,
  IconChevronDown,
  IconGripVertical,
} from '@tabler/icons-vue'
import type { Node, NodeType } from '@/types/page-builder'
import { BUILDER_KEY } from '@/components/builder/injection'

defineOptions({ name: 'TreeRow' })
const props = withDefaults(
  defineProps<{
    node: Node
    depth?: number
    isLast?: boolean
    ancestorsIsLast?: boolean[]
  }>(),
  {
    depth: 0,
    isLast: false,
    ancestorsIsLast: () => [],
  },
)

const nextAncestorsIsLast = computed(() => [...props.ancestorsIsLast, props.isLast])
const ancestorLines = computed(() => props.ancestorsIsLast.slice(0, -1))

const store = inject(BUILDER_KEY, null)!
const expanded = ref(true)

// Indent geometry. guideX(d) = x of the vertical guide for children of a
// Indent geometry. guideX(d) = x of the vertical guide for children of a
// depth-d node, and the start x of the elbow for a depth-(d+1) child.
// For containers the elbow stops at the chevron center; for leaves it
// reaches the node-type icon past the spacer + gap.
const INDENT = 20
const pad = (d: number) => d * INDENT + 8
const guideX = (d: number) => d * INDENT + 16

const selected = () => store.selectedId.value === props.node.id

const CONTAINER_TYPES = new Set<NodeType>(['frame', 'section', 'link', 'component'])
// Resolve a component instance to its master root so the tree shows the
// instance's subtree (children live on the master, not the marker node).
const resolved = computed(() => {
  if (props.node.type === 'component' && props.node.componentId) {
    const master = store.components.masterRoot(props.node.componentId)
    if (master) return master
  }
  return props.node
})
const isContainer = () => CONTAINER_TYPES.has(props.node.type)
const hasChildren = () => resolved.value.children.length > 0
const hasChevron = computed(() => isContainer() && hasChildren())
const elbowWidth = computed(() => hasChevron.value ? 12 : 30)
const dragging = () => store.draggingId.value === props.node.id
const dropInside = () => {
  const t = store.dropTarget.value
  return !!t && t.pos === 'inside' && t.parentId === props.node.id
}

const NODE_ICONS: Record<NodeType, any> = {
  frame: IconSquare,
  section: IconSection,
  text: IconTypography,
  heading: IconH1,
  image: IconPhoto,
  button: IconCircuitPushbutton,
  link: IconLink,
  component: IconComponents,
}

const NODE_COLORS: Record<NodeType, string> = {
  frame: 'text-violet-500',
  section: 'text-sky-500',
  text: 'text-muted-foreground',
  heading: 'text-amber-500',
  image: 'text-emerald-500',
  button: 'text-rose-500',
  link: 'text-blue-500',
  component: 'text-purple-500',
}

function toggleExpand(e: MouseEvent) {
  e.stopPropagation()
  expanded.value = !expanded.value
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
  <div>
    <!-- Row (Only the row itself is draggable and droppable) -->
    <div
      :draggable="depth > 0"
      class="group relative flex w-full items-center gap-0.5 py-[5px] pr-2 text-left text-xs transition-colors"
      :class="[
        selected()
          ? 'bg-primary/10 text-primary font-medium'
          : 'text-foreground hover:bg-muted/50',
        dragging() ? 'opacity-30' : '',
        dropInside() ? 'ring-2 ring-blue-400 ring-inset rounded-sm' : '',
      ]"
      :style="{ paddingLeft: `${pad(depth)}px` }"
      @click.stop="store.select(node.id)"
      @dragstart="onDragStart"
      @dragover="onDragOver"
      @drop="onDrop"
      @dragend="onDragEnd"
    >
      <!-- Ancestor vertical lines -->
      <template v-for="(_parentLast, idx) in ancestorLines" :key="idx">
        <span
          v-if="!ancestorsIsLast[idx + 1]"
          class="pointer-events-none absolute border-l border-input"
          :style="{ left: `${guideX(idx)}px`, top: 0, height: '100%', zIndex: 1 }"
        />
      </template>

      <!-- Elbow connector: horizontal stub from the parent vertical guide
           pointing at this node's icon. -->
      <span
        v-if="depth > 0"
        class="pointer-events-none absolute top-1/2 border-t border-input"
        :style="{ left: `${guideX(depth - 1)}px`, width: `${elbowWidth}px`, transform: 'translateY(-50%)', zIndex: 1 }"
      />
      <!-- Vertical guide segment for this row's depth. Last sibling stops at
           row center; others span the full row height so they chain downward. -->
      <span
        v-if="depth > 0"
        class="pointer-events-none absolute border-l border-input"
        :style="{ left: `${guideX(depth - 1)}px`, top: 0, height: isLast ? '50%' : '100%', zIndex: 1 }"
      />

      <!-- Expand/collapse toggle -->
      <button
        v-if="isContainer() && hasChildren()"
        class="flex size-4 shrink-0 items-center justify-center rounded-sm text-muted-foreground hover:bg-muted hover:text-muted-foreground"
        @click="toggleExpand"
      >
        <IconChevronDown v-if="expanded" class="size-3" />
        <IconChevronRight v-else class="size-3" />
      </button>
      <span v-else class="size-4 shrink-0" />

      <!-- Node type icon -->
      <component
        :is="NODE_ICONS[node.type] ?? IconSquare"
        class="size-3.5 shrink-0"
        :class="selected() ? 'text-blue-500' : NODE_COLORS[node.type] ?? 'text-muted-foreground'"
      />

      <!-- Node label -->
      <span class="ml-1 min-w-0 flex-1 truncate select-none">{{ node.name || node.type }}</span>

      <!-- Child count badge for collapsed containers -->
      <span
        v-if="isContainer() && hasChildren() && !expanded"
        class="shrink-0 rounded-full bg-muted px-1.5 text-[10px] text-muted-foreground"
      >
        {{ resolved.children.length }}
      </span>

      <!-- Drag handle (right side, visible on hover) -->
      <IconGripVertical
        v-if="depth > 0"
        class="size-3.5 shrink-0 cursor-grab text-neutral-300 opacity-0 transition-opacity group-hover:opacity-100"
      />
    </div>

    <!-- Children (collapsible, not draggable from here directly) -->
    <div v-if="expanded && hasChildren()">
      <div class="relative">
        <TreeRow
          v-for="(child, i) in resolved.children"
          :key="child.id"
          :node="child"
          :depth="depth + 1"
          :is-last="i === resolved.children.length - 1"
          :ancestors-is-last="nextAncestorsIsLast"
        />
      </div>
    </div>
  </div>
</template>
