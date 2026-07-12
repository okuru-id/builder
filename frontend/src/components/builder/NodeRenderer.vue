<script setup lang="ts">
// Recursive node renderer. Two modes:
//   - interactive (default): click to select, dblclick text/heading to inline-edit
//   - readonly: pure render for export preview / iframe snapshot
import { computed, inject, nextTick, ref } from 'vue'
import type { Node } from '@/types/page-builder'
import { TEXT_TYPES } from '@/types/page-builder'
import { BUILDER_KEY } from '@/components/builder/injection'

const props = withDefaults(
  defineProps<{ node: Node; depth?: number; readonly?: boolean }>(),
  { depth: 0, readonly: false },
)

const store = inject(BUILDER_KEY, null)
const editing = ref(false)
const elRef = ref<HTMLElement | null>(null)

const tag = computed<string>(() => tagFor(props.node))
const isTextLike = computed(() => TEXT_TYPES.has(props.node.type))
const selected = computed(
  () => !props.readonly && store?.selectedId.value === props.node.id,
)
const classList = computed(() => {
  const base = props.node.classes
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
    case 'button':
      return 'button'
    case 'link':
      return 'a'
    case 'section':
      return 'section'
    case 'frame':
    case 'component':
    default:
      return 'div'
  }
}

function attrsFor(n: Node): Record<string, unknown> {
  const a: Record<string, unknown> = {}
  if (n.type === 'image') {
    if (n.props.src) a.src = n.props.src
    if (n.props.alt) a.alt = n.props.alt
  } else if (n.type === 'link') {
    if (n.props.href) a.href = n.props.href
  }
  return a
}

function onClick(e: MouseEvent) {
  if (props.readonly) return
  e.stopPropagation()
  store?.select(props.node.id)
}

async function onDblClick(e: MouseEvent) {
  if (props.readonly || !isTextLike.value || !store) return
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
</script>

<template>
  <component
    :is="tag"
    :ref="(el: unknown) => { if (editing) elRef = el as HTMLElement }"
    :class="classList"
    v-bind="{ ...attrsFor(node), ...interactiveAttrs }"
    @click="onClick"
    @dblclick="onDblClick"
    @blur="onBlur"
    @keydown.enter.prevent="commitText"
    @keydown.esc.prevent="commitText"
  >
    <template v-if="isTextLike">
      {{ node.props.text }}
    </template>
    <template v-else-if="node.type === 'button'">
      {{ node.props.text }}
    </template>
    <template v-else-if="node.type === 'image'">
      <!-- void element, no children -->
    </template>
    <template v-else>
      <NodeRenderer
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :depth="depth + 1"
        :readonly="readonly"
      />
    </template>
  </component>
</template>
