<script setup lang="ts">
// Left outline panel: tree view + add-node palette.
import { computed, inject, ref, watch, useTemplateRef } from 'vue'
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

const treeHeight = ref(Number(localStorage.getItem('builder:left:tree-height') ?? 300))
const componentsHeight = ref(Number(localStorage.getItem('builder:left:components-height') ?? 240))

watch(treeHeight, (v) => localStorage.setItem('builder:left:tree-height', String(v)))
watch(componentsHeight, (v) => localStorage.setItem('builder:left:components-height', String(v)))

function startResize(which: 'tree' | 'components', event: MouseEvent) {
  event.preventDefault()
  const panel = asideRef.value
  const hdr = headerRef.value
  if (!panel || !hdr) return
  const panelHeight = panel.clientHeight
  const hdrHeight = hdr.offsetHeight
  const startY = event.clientY
  const start = which === 'tree' ? treeHeight.value : componentsHeight.value
  const sectionOverhead = 60
  const addNodeOverhead = 60

  const onMove = (e: MouseEvent) => {
    const delta = e.clientY - startY
    let next = start + delta
    next = Math.max(180, next)

    // compute total occupied height if this resize goes through, then clamp
    const otherH = (which === 'tree' ? componentsHeight.value : treeHeight.value) + sectionOverhead
    const total = hdrHeight + (next + sectionOverhead) + otherH + addNodeOverhead
    const overflow = total - panelHeight
    if (overflow > 0) {
      next = Math.max(180, next - overflow)
    }

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
  <aside ref="aside" class="flex h-full w-full flex-col border-r border-border bg-background">
    <!-- Header -->
    <div ref="header" class="flex shrink-0 items-center justify-between border-b border-border px-3 py-2.5">
      <div class="flex items-center gap-1.5">
        <IconLayersIntersect class="size-4 text-muted-foreground" />
        <h2 class="text-sm font-semibold">Layer</h2>
      </div>
      <span class="rounded-full bg-muted px-1.5 py-0.5 text-[10px] text-muted-foreground">
        {{ nodeCount }}
      </span>
    </div>

    <div class="flex min-h-0 flex-1 flex-col overflow-hidden">
      <!-- Layer Tree -->
      <div class="shrink-0 border-b border-border/50">
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
      <div class="shrink-0 border-b border-border/50">
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
