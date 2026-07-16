<script setup lang="ts">
// Recursive node renderer. Two modes:
//   - interactive (default): click to select, dblclick text/heading to inline-edit
//   - readonly: pure render for export preview / iframe snapshot
import { computed, inject, nextTick, ref } from 'vue'
import type { Node } from '@/types/page-builder'
import { CONTAINER_TYPES, TEXT_TYPES } from '@/types/page-builder'
import type { Breakpoint } from '@/types/page-builder'
import { BUILDER_KEY } from '@/components/builder/injection'
import { ICONS } from '@/lib/icon-map'

const props = withDefaults(
  defineProps<{ node: Node; depth?: number; readonly?: boolean; ancestorId?: string }>(),
  { depth: 0, readonly: false, ancestorId: '' },
)

const store = inject(BUILDER_KEY, null)
const editing = ref(false)
const elRef = ref<HTMLElement | null>(null)

// If this node is inside a component instance, detach on first edit.
const instanceOwner = computed(() => props.ancestorId || (isInstance.value ? props.node.id : ''))

const tag = computed<string>(() => tagFor(displayNode.value))
const isTextLike = computed(() => TEXT_TYPES.has(displayNode.value.type))
const selected = computed(
  () => !props.readonly && store?.selectedId.value === props.node.id,
)
// Resolve a component instance to its master root for display. Selection/drag
// still bind to the instance node (props.node.id). Children rendered readonly.
const displayNode = computed<Node>(() => {
  if (props.node.type === 'component' && props.node.componentId && store) {
    const master = store.components.masterRoot(props.node.componentId)
    if (master) return master
  }
  return props.node
})
const isInstance = computed(() => props.node.type === 'component' && !!props.node.componentId)
const inInstance = computed(() => !!instanceOwner.value)
// Per-breakpoint visibility: node hidden on the current canvas breakpoint is
// rendered faint + dashed so it stays selectable/editable. Checks both the
// node itself (covers component-instance overrides) and its display node
// (master). ponytail: JS-driven (not a Tailwind class) so it reacts live and
// never leaks into published output (codegen emits md:hidden / lg:hidden).
const hiddenHere = computed(() => {
  if (props.readonly) return false
  const bp = store?.breakpoint.value as Breakpoint | undefined
  if (!bp) return false
  const inst = props.node.hiddenOn
  const disp = displayNode.value.hiddenOn
  return !!(inst?.includes(bp) || disp?.includes(bp))
})
const classList = computed(() => {
  const base = displayNode.value.classes
  if (props.readonly) return base
  // Selection ring + hover affordance, layered so we never overwrite real classes.
  return selected.value
    ? [...base, 'outline-2', 'outline-blue-500', 'outline-dashed', '-outline-offset-2']
    : [...base, 'hover:outline-1', 'hover:outline-blue-300', 'hover:-outline-offset-1']
})

function tagFor(n: Node): string {
  switch (n.type) {
    case 'text':
      return 'span'
    case 'heading': {
      const lvl = Number(n.props.level)
      if (lvl >= 1 && lvl <= 6) return `h${lvl}`
      return 'h2'
    }
    case 'image':
      return 'img'
    case 'divider':
      return 'hr'
    case 'button':
      return 'button'
    case 'link':
      return 'a'
    case 'section':
      return 'section'
    case 'form':
      return 'form'
    case 'frame':
    case 'grid':
    case 'input':
    case 'component':
    default:
      return 'div'
  }
}

function attrsFor(n: Node): Record<string, unknown> {
  const a: Record<string, unknown> = {}
  if (n.type === 'icon') {
    return {}
  }
  if (n.type === 'image') {
    if (n.props.src) a.src = n.props.src
    if (n.props.alt) a.alt = n.props.alt
  } else if (n.type === 'link') {
    if (n.props.href) a.href = n.props.href
  }
  return a
}
function onClick(e: MouseEvent) {
  if (props.readonly && !inInstance.value) return
  e.stopPropagation()
  // Click inside an instance selects the instance itself.
  const targetId = inInstance.value ? instanceOwner.value : props.node.id
  store?.select(targetId || props.node.id)
}

async function onDblClick(e: MouseEvent) {
  if (!isTextLike.value || !store) return
  // Allow breaking instances even in readonly mode.
  if (props.readonly && !inInstance.value) return
  // Instances: break (detach) on first edit, then let user edit the copy directly.
  if (inInstance.value || isInstance.value) {
    e.stopPropagation()
    const id = inInstance.value ? instanceOwner.value : props.node.id
    if (id) store.breakInstance(id)
    return
  }
  e.stopPropagation()
  editing.value = true
  await nextTick()
  elRef.value?.focus()
  // place caret at end
  const range = document.createRange()
  range.selectNodeContents(elRef.value!)
  range.collapse(false)
  const sel = window.getSelection()
  sel?.removeAllRanges()
  sel?.addRange(range)
}

function commitText() {
  if (!editing.value || !store || !elRef.value) return
  const text = elRef.value.innerText
  store.patchNode(props.node.id, { props: { ...props.node.props, text } })
  editing.value = false
}

function onBlur() {
  commitText()
}

const interactiveAttrs = computed(() =>
  props.readonly
    ? {}
    : {
        contenteditable: editing.value ? 'true' : 'false',
        spellcheck: 'false',
        suppressContentEditableWarning: 'true',
      },
)

// --- drag-and-drop ---
const isContainer = computed(() => CONTAINER_TYPES.has(props.node.type))
const dropInside = computed(() => {
  const t = store?.dropTarget.value
  return !!t && t.pos === 'inside' && t.parentId === props.node.id
})
const dragging = computed(
  () => !props.readonly && store?.draggingId.value === props.node.id,
)

function onDragStart(e: DragEvent) {
  if (props.readonly || !store) return
  store.dragStart(props.node.id)
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    // Firefox requires setData to start DnD.
    e.dataTransfer.setData('text/plain', props.node.id)
  }
}
function onDragOver(e: DragEvent) {
  if (props.readonly || !store) return
  store.dragOver(props.node.id, e, isContainer.value)
}
function onDragLeave() {
  // ponytail: we do not clear dropTarget here — dragover fires continuously on
  // adjacent nodes and would flicker. Cleared on drop/dragend.
}
function onDrop(e: DragEvent) {
  if (props.readonly || !store) return
  e.preventDefault()
  e.stopPropagation()
  store.drop()
}
function onDragEnd() {
  store?.dragEnd()
}
</script>

<template>
  <component
    :is="tag"
    ref="elRef"
    :style="props.node.hidden ? { display: 'none' } : undefined"
    :class="[classList, {
      'opacity-40 outline-dashed outline-2 outline-amber-400 -outline-offset-2': hiddenHere,
      'opacity-40': dragging,
      'ring-2 ring-blue-500 ring-inset': dropInside,
    }]"
    :draggable="!readonly && !editing"
    v-bind="{ ...attrsFor(displayNode), ...interactiveAttrs }"
    @click="onClick"
    @dblclick="onDblClick"
    @blur="onBlur"
    @keydown.enter.prevent="commitText"
    @keydown.esc.prevent="commitText"
    @dragstart="onDragStart"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
    @dragend="onDragEnd"
  >
    <template v-if="isTextLike">
      {{ displayNode.props.text }}
    </template>
    <template v-else-if="displayNode.type === 'button'">
      {{ displayNode.props.text }}
    </template>
    <template v-else-if="displayNode.type === 'link' && (!displayNode.children || displayNode.children.length === 0)">
      {{ displayNode.props.text }}
    </template>
    <template v-else-if="displayNode.type === 'icon'">
      <span v-if="displayNode.props.icon && ICONS[displayNode.props.icon]" class="inline-flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round">
          <template v-for="seg in ICONS[displayNode.props.icon]" :key="seg[1].key">
            <path :d="seg[1].d" />
          </template>
        </svg>
      </span>
      <span v-else>{{ displayNode.props.icon ?? '✦' }}</span>
    </template>
    <template v-else-if="displayNode.type === 'input'">
      <label v-if="displayNode.props.label" class="text-sm font-medium">{{ displayNode.props.label }}</label>
      <input
        :type="displayNode.props.inputType ?? 'text'"
        :placeholder="displayNode.props.placeholder ?? ''"
        :required="displayNode.props.required ?? false"
        class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
        @click.stop
      />
    </template>
    <template v-else-if="displayNode.type === 'image'">
      <!-- void element, no children -->
    </template>
    <template v-else>
      <NodeRenderer
        v-for="child in displayNode.children"
        :key="child.id"
        :node="child"
        :depth="depth + 1"
        :readonly="readonly || isInstance"
        :ancestor-id="instanceOwner"
      />
    </template>
  </component>
</template>
