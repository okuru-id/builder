<script setup lang="ts">
// Canvas viewport — Figma-style navigation:
//   • Ctrl+Scroll  → zoom toward cursor
//   • Space+drag   → pan (grab cursor)
//   • Middle mouse → pan
//   • +/− / Fit    → zoom HUD controls
//
// Layout model:
//   root (relative, overflow:hidden, fills panel)
//     ├─ scrollEl (overflow:auto, h-full w-full)
//     │    └─ workspace (WORKSPACE_PAD*2 + canvas*zoom — always bigger than viewport)
//     │         └─ canvas (native 1920px, scaled via CSS transform, origin: top left)
//     └─ HUD (absolute bottom-center)
//
// Key insight: workspace padding is CONSTANT (not scaled), so there is always
// a large scrollable area even when zoom is very small. This makes panning
// always possible, mirroring how Figma/Penpot work.
import { computed, inject, nextTick, onMounted, onUnmounted, ref } from 'vue'
import NodeRenderer from '@/components/builder/NodeRenderer.vue'
import { BUILDER_KEY } from '@/components/builder/injection'
import type { Breakpoint } from '@/components/builder/useBuilderStore'

const store = inject(BUILDER_KEY, null)
if (!store) throw new Error('Canvas requires a BuilderStore provider')

// ─── Dimensions ────────────────────────────────────────────────────
const CANVAS_DIMS: Record<Breakpoint, { w: number; h: number }> = {
  desktop: { w: 1920, h: 1080 },
  tablet:  { w: 768,  h: 1024 },
  mobile:  { w: 390,  h: 844  },
}

// WORKSPACE_PAD: constant extra space around the canvas in stage-px.
// NOT multiplied by zoom — so there is always a large scrollable area.
const WORKSPACE_PAD = 1200

const dims = computed(() => CANVAS_DIMS[store!.breakpoint.value])

// ─── Zoom ───────────────────────────────────────────────────────────
const zoom     = ref(0.4)
const MIN_ZOOM = 0.05
const MAX_ZOOM = 3.0
const ZOOM_STEP = 0.05

// Workspace (stage) size: canvas scaled + constant padding on each side.
const stageW = computed(() => dims.value.w * zoom.value + WORKSPACE_PAD * 2)
const stageH = computed(() => dims.value.h * zoom.value + WORKSPACE_PAD * 2)

// Canvas element: native resolution, transform-scaled, positioned at
// (WORKSPACE_PAD, WORKSPACE_PAD) inside the stage (constant, not zoomed).
const canvasStyle = computed(() => ({
  width:           `${dims.value.w}px`,
  minHeight:       `${dims.value.h}px`,
  transform:       `scale(${zoom.value})`,
  transformOrigin: 'top left',
  position:        'absolute' as const,
  top:             `${WORKSPACE_PAD}px`,
  left:            `${WORKSPACE_PAD}px`,
}))

function clampZoom(v: number) {
  return Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, parseFloat(v.toFixed(3))))
}

const scrollEl = ref<HTMLElement | null>(null)

/**
 * Zoom to `nextZoom` keeping the viewport point (pivotX, pivotY) fixed.
 *
 * World-space coordinates (what the stage pixel refers to) change with zoom.
 * Canvas top-left is always at stage coord (WORKSPACE_PAD, WORKSPACE_PAD).
 * A world point at stage coord S shows at viewport coord V when:
 *   S = scrollLeft + V
 * After zoom, the stage is rescaled:
 *   new_stage_x = WORKSPACE_PAD + (S - WORKSPACE_PAD) * (newZoom / oldZoom)
 * We want new_stage_x to appear at V:
 *   newScrollLeft + V = new_stage_x
 *   newScrollLeft = new_stage_x - V
 */
async function zoomAround(nextZoom: number, pivotX?: number, pivotY?: number) {
  const el   = scrollEl.value
  const prev = zoom.value
  const next = clampZoom(nextZoom)
  if (next === prev) return

  let targetScrollLeft: number | undefined
  let targetScrollTop: number | undefined

  if (el && pivotX !== undefined && pivotY !== undefined) {
    const ratio = next / prev
    // Stage coordinate currently under the pivot point:
    const stageX = el.scrollLeft + pivotX
    const stageY = el.scrollTop  + pivotY
    // After zoom, canvas origin moves; only the region beyond WORKSPACE_PAD scales.
    const newStageX = WORKSPACE_PAD + (stageX - WORKSPACE_PAD) * ratio
    const newStageY = WORKSPACE_PAD + (stageY - WORKSPACE_PAD) * ratio
    targetScrollLeft = newStageX - pivotX
    targetScrollTop  = newStageY - pivotY
  }

  zoom.value = next
  await nextTick() // let Vue resize the stage before repositioning scroll

  if (el && targetScrollLeft !== undefined && targetScrollTop !== undefined) {
    el.scrollLeft = targetScrollLeft
    el.scrollTop  = targetScrollTop
  }
}

function zoomIn()  { zoomAround(zoom.value + ZOOM_STEP) }
function zoomOut() { zoomAround(zoom.value - ZOOM_STEP) }

function zoomFit() {
  const el = scrollEl.value
  if (!el) return
  const next = clampZoom((el.clientWidth - 80) / dims.value.w)
  zoom.value = next
  // After fit, center the canvas in the viewport.
  nextTick(() => {
    if (!el) return
    // Canvas top-left is at WORKSPACE_PAD in stage; its visual size = dim*zoom.
    // Center: scroll so canvas center aligns with viewport center.
    const canvasCenterX = WORKSPACE_PAD + (dims.value.w * zoom.value) / 2
    const canvasCenterY = WORKSPACE_PAD + (dims.value.h * zoom.value) / 2
    el.scrollLeft = canvasCenterX - el.clientWidth  / 2
    el.scrollTop  = canvasCenterY - el.clientHeight / 2
  })
}

function setZoomPct(pct: number) { zoom.value = clampZoom(pct / 100) }
const zoomPct = computed(() => Math.round(zoom.value * 100))

// ─── Ctrl+Wheel zoom ───────────────────────────────────────────────
function onWheel(e: WheelEvent) {
  if (!e.ctrlKey && !e.metaKey) return
  e.preventDefault()
  const el = scrollEl.value
  if (!el) return
  const rect   = el.getBoundingClientRect()
  const pivotX = e.clientX - rect.left
  const pivotY = e.clientY - rect.top
  // Use multiplicative step for smoother feel (same as Figma).
  const factor = e.deltaY < 0 ? (1 + ZOOM_STEP) : (1 - ZOOM_STEP)
  zoomAround(zoom.value * factor, pivotX, pivotY)
}

// ─── Pan (Space+drag / Middle-mouse) ───────────────────────────────
const spaceHeld = ref(false)
const isPanning = ref(false)
let wasPanning = false // set on mouseup, checked by click handler
let panOriginX = 0, panOriginY = 0
let panScrollL = 0, panScrollT = 0

const cursorClass = computed(() => {
  if (isPanning.value)  return 'cursor-grabbing'
  if (spaceHeld.value)  return 'cursor-grab'
  return ''
})

function startPan(clientX: number, clientY: number) {
  const el = scrollEl.value
  if (!el) return
  isPanning.value = true
  panOriginX = clientX
  panOriginY = clientY
  panScrollL = el.scrollLeft
  panScrollT = el.scrollTop
}

function onMousedown(e: MouseEvent) {
  if (e.button === 1 || (e.button === 0 && spaceHeld.value)) {
    e.preventDefault()
    startPan(e.clientX, e.clientY)
  }
}

function onMousemove(e: MouseEvent) {
  if (!isPanning.value) return
  const el = scrollEl.value
  if (!el) return
  el.scrollLeft = panScrollL - (e.clientX - panOriginX)
  el.scrollTop  = panScrollT - (e.clientY - panOriginY)
}

function onMouseup() {
  if (isPanning.value) wasPanning = true
  isPanning.value = false
  spaceHeld.value && (spaceHeld.value = spaceHeld.value) // keep space state
  setTimeout(() => { wasPanning = false }, 50)
}

function onKeydown(e: KeyboardEvent) {
  if (e.repeat) return
  const t = e.target as HTMLElement
  if (t.isContentEditable || ['INPUT', 'TEXTAREA', 'SELECT'].includes(t.tagName)) return

  if (e.code === 'Space') {
    e.preventDefault()
    spaceHeld.value = true
    return
  }
  // Keyboard shortcuts: F = fit, 0 = 100%
  if (e.key === 'f' || e.key === 'F') { zoomFit(); return }
  if (e.key === '0') { zoom.value = 1; return }
  if (e.key === '=' || e.key === '+') { zoomIn(); return }
  if (e.key === '-') { zoomOut(); return }
}

function onKeyup(e: KeyboardEvent) {
  if (e.code !== 'Space') return
  spaceHeld.value = false
  isPanning.value = false
}

// ─── Lifecycle ─────────────────────────────────────────────────────

onMounted(() => {
  zoomFit()
  const el = scrollEl.value!
  el.addEventListener('wheel',     onWheel,     { passive: false })
  el.addEventListener('mousedown', onMousedown)
  window.addEventListener('mousemove', onMousemove)
  window.addEventListener('mouseup',   onMouseup)
  window.addEventListener('keydown',   onKeydown)
  window.addEventListener('keyup',     onKeyup)
})

onUnmounted(() => {
  const el = scrollEl.value
  el?.removeEventListener('wheel',     onWheel)
  el?.removeEventListener('mousedown', onMousedown)
  window.removeEventListener('mousemove', onMousemove)
  window.removeEventListener('mouseup',   onMouseup)
  window.removeEventListener('keydown',   onKeydown)
  window.removeEventListener('keyup',     onKeyup)
})

function onBackgroundClick() {
  if (wasPanning) return
  store!.select(null)
}
</script>

<template>
  <!-- Root: relative + overflow:hidden so stage overflow is clipped -->
  <div class="relative h-full w-full overflow-hidden">

    <!-- Scroll container: always has overflow because stage >= viewport -->
    <div
      ref="scrollEl"
      class="h-full w-full overflow-auto select-none"
      :class="cursorClass"
      style="background-color: var(--canvas-bg); background-image: radial-gradient(circle, var(--canvas-dot) 1px, transparent 1px); background-size: 20px 20px;"
      @click.self="onBackgroundClick"
    >
      <!--
        Stage/workspace: always wider + taller than viewport thanks to WORKSPACE_PAD.
        Canvas sits at (WORKSPACE_PAD, WORKSPACE_PAD) inside this space.
      -->
      <div
        :style="{
          width:      `${stageW}px`,
          height:     `${stageH}px`,
          position:   'relative',
          flexShrink: 0,
        }"
        @click.self="onBackgroundClick"
      >
        <!-- Canvas at native resolution, CSS-scaled via transform -->
        <div
          :style="canvasStyle"
          class="bg-white shadow-2xl ring-1 ring-black/10"
          @click.self="onBackgroundClick"
        >
          <NodeRenderer :node="store.tree.value.root" />
        </div>
      </div>
    </div>

    <!-- Floating Zoom HUD — always pinned bottom-center, outside scroll -->
    <div class="pointer-events-none absolute bottom-5 left-0 right-0 z-30 flex justify-center">
      <div class="pointer-events-auto flex items-center gap-0.5 rounded-full border border-white/10 bg-black/75 px-2 py-1.5 shadow-xl backdrop-blur">

        <button
          class="flex size-7 items-center justify-center rounded-full text-white/60 hover:bg-white/10 hover:text-white transition-colors text-base leading-none font-light select-none"
          title="Zoom out (Ctrl+Scroll)"
          @click.stop="zoomOut"
        >−</button>

        <div class="flex items-center" title="Ketik persentase zoom">
          <input
            type="number"
            :value="zoomPct"
            min="5" max="300" step="5"
            class="w-12 bg-transparent text-center text-xs font-mono text-white/90 focus:outline-none [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
            @change="setZoomPct(Number(($event.target as HTMLInputElement).value))"
          />
          <span class="text-xs text-white/40 -ml-0.5">%</span>
        </div>

        <button
          class="flex size-7 items-center justify-center rounded-full text-white/60 hover:bg-white/10 hover:text-white transition-colors text-base leading-none font-light select-none"
          title="Zoom in (Ctrl+Scroll)"
          @click.stop="zoomIn"
        >+</button>

        <div class="mx-1.5 h-4 w-px bg-white/15" />

        <button
          class="rounded-full px-2.5 py-0.5 text-xs text-white/60 hover:bg-white/10 hover:text-white transition-colors select-none"
          title="Fit canvas ke panel"
          @click.stop="zoomFit"
        >Fit</button>

        <div class="mx-1.5 h-4 w-px bg-white/15" />

        <span class="pr-1 text-[10px] text-white/25 font-mono select-none">Space+drag · Mid</span>
      </div>
    </div>
  </div>
</template>
