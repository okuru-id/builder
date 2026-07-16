<script setup lang="ts">
// Left outline panel: tree view + add-node palette.
import { computed, inject, provide, ref, watch, useTemplateRef } from 'vue'
import {
  IconSquare,
  IconSection,
  IconTypography,
  IconH1,
  IconPhoto,
  IconCircuitPushbutton,
  IconLink,
  IconLayersIntersect,
  IconArrowsMinimize,
  IconGridDots,
  IconMinus,
  IconStar,
  IconForms,
  IconInputSearch,
  IconSettings,
  IconComponents,
} from '@tabler/icons-vue'
import { PALETTE_TYPES } from '@/types/page-builder'
import type { NodeType } from '@/types/page-builder'
import { BUILDER_KEY } from '@/components/builder/injection'
import TreeRow from '@/components/builder/TreeRow.vue'
import ComponentPalette from '@/components/builder/ComponentPalette.vue'
import InspectorSection from './inspector/InspectorSection.vue'

const store = inject(BUILDER_KEY, null)!
const asideRef = useTemplateRef<HTMLElement>('aside')
const headerRef = useTemplateRef<HTMLElement>('header')
const treeWrapRef = useTemplateRef<HTMLElement>('treeWrap')
const compWrapRef = useTemplateRef<HTMLElement>('compWrap')

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

// ponytail: shared collapsed-set keyed by node id. Provided to TreeRow so
// expand/collapse state survives child unmount on parent collapse. A local
// ref in TreeRow would reset every time the row leaves the DOM.
const collapsed = ref(new Set<string>())
provide('builder-tree-collapsed', collapsed)
function walkIds(node: { id: string; children: any[] }, acc: string[]) {
  acc.push(node.id)
  for (const c of node.children ?? []) walkIds(c, acc)
}
function collapseAll() {
  const ids: string[] = []
  walkIds(root.value, ids)
  collapsed.value = new Set(ids)
}

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

const treeHeight = ref(Number(localStorage.getItem('builder:left:tree-height') ?? 300))
const componentsHeight = ref(Number(localStorage.getItem('builder:left:components-height') ?? 240))

watch(treeHeight, (v) => localStorage.setItem('builder:left:tree-height', String(v)))
watch(componentsHeight, (v) => localStorage.setItem('builder:left:components-height', String(v)))

function startResize(which: 'tree' | 'components', event: MouseEvent) {
  event.preventDefault()
  const panel = asideRef.value
  const hdr = headerRef.value
  const treeWrap = treeWrapRef.value
  const compWrap = compWrapRef.value
  if (!panel || !hdr || !treeWrap || !compWrap) return
  const startY = event.clientY
  const start = which === 'tree' ? treeHeight.value : componentsHeight.value

  // The OTHER resizable section, measured live so collapsed InspectorSection
  // (chevron-collapsed to header-only) frees space the resizer can reclaim.
  // Add Node is flex-1 and absorbs whatever is left — never measure it (circular);
  // only reserve its minimum header so it stays visible.
  const otherSectionActual = which === 'tree' ? compWrap.offsetHeight : treeWrap.offsetHeight
  const OWN_OVERHEAD = 44
  const ADD_NODE_MIN = 32
  const MIN_CONTENT = 120
  const maxNext = panel.clientHeight - hdr.offsetHeight - OWN_OVERHEAD - otherSectionActual - ADD_NODE_MIN

  const onMove = (e: MouseEvent) => {
    const delta = e.clientY - startY
    const next = Math.max(MIN_CONTENT, Math.min(start + delta, maxNext))
    if (which === 'tree') treeHeight.value = next
    else componentsHeight.value = next
  }
  const onUp = () => {
    window.removeEventListener('mousemove', onMove)
    window.removeEventListener('mouseup', onUp)
  }
  window.addEventListener('mousemove', onMove)
  window.addEventListener('mouseup', onUp)
}
</script>

<template>
  <aside ref="aside" class="flex h-full w-full flex-col border-r border-border bg-card">
    <!-- Header -->
    <div ref="header" class="flex shrink-0 items-center justify-between border-b border-border px-3 py-2.5">
      <div class="flex items-center gap-1.5">
        <IconLayersIntersect class="size-4 text-muted-foreground" />
        <h2 class="text-sm font-semibold">Layer</h2>
      </div>
      <div class="flex items-center gap-1">
        <button
          class="flex size-6 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
          title="Collapse all layers"
          @click="collapseAll"
        >
          <IconArrowsMinimize class="size-3.5" />
        </button>
        <span class="rounded-full bg-muted px-1.5 py-0.5 text-[10px] text-muted-foreground">
          {{ nodeCount }}
        </span>
      </div>
    </div>

    <div class="flex min-h-0 flex-1 flex-col overflow-hidden">
      <!-- Layer Tree -->
      <div ref="treeWrap" class="shrink-0 border-b border-border/50">
        <InspectorSection title="Layer Tree" :icon="IconSettings">
          <div class="relative overflow-auto py-0.5" :style="{ height: `${treeHeight}px` }">
            <TreeRow :node="root" />
          </div>
          <div class="mt-1 flex justify-center">
            <button
              class="h-2 w-12 cursor-row-resize rounded-full bg-muted hover:bg-muted-foreground/20"
              title="Resize layer tree"
              @mousedown="startResize('tree', $event)"
            />
          </div>
        </InspectorSection>
      </div>

      <!-- Components -->
      <div ref="compWrap" class="shrink-0 border-b border-border/50">
        <InspectorSection title="Components" :icon="IconComponents">
          <div class="relative overflow-auto border-t border-border -mx-3" :style="{ height: `${componentsHeight}px` }">
            <ComponentPalette />
          </div>
          <div class="mt-1 flex justify-center">
            <button
              class="h-2 w-12 cursor-row-resize rounded-full bg-muted hover:bg-muted-foreground/20"
              title="Resize components section"
              @mousedown="startResize('components', $event)"
            />
          </div>
        </InspectorSection>
      </div>

      <!-- Add Node: fills remaining space -->
      <div class="min-h-0 flex-1 overflow-auto">
        <InspectorSection title="Add Node" :icon="IconSquare">
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
        </InspectorSection>
      </div>
    </div>
  </aside>
</template>
