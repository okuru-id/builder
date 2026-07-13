<script setup lang="ts">
// Top toolbar: 3-column layout. Left=identity, center=toggles+breakpoint, right=actions.
// All action buttons icon-only with native title tooltips. ponytail: title attr
// over Tooltip component — zero provider wiring, native, sufficient.
import { inject, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  IconArrowLeft,
  IconDeviceDesktop,
  IconDeviceTablet,
  IconDeviceMobile,
  IconCode,
  IconDeviceFloppy,
  IconCheck,
  IconExternalLink,
  IconLayoutSidebar,
  IconLayoutSidebarRight,
  IconRocket,
  IconArrowBackUp,
  IconArrowForwardUp,
} from '@tabler/icons-vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { BUILDER_KEY } from '@/components/builder/injection'
import type { Breakpoint } from '@/components/builder/useBuilderStore'
import ExportDialog from './ExportDialog.vue'
import HistoryDropdown from './HistoryDropdown.vue'

const store = inject(BUILDER_KEY, null)!
const router = useRouter()

defineProps<{ showLeft: boolean; showRight: boolean }>()
defineEmits<{ 'toggle-left': []; 'toggle-right': [] }>()

const editingName = ref(false)
const nameBuffer = ref('')

function startEditName() {
  nameBuffer.value = store.page.value?.name ?? ''
  editingName.value = true
}

function commitName() {
  if (nameBuffer.value.trim()) {
    store.rename(nameBuffer.value.trim())
  }
  editingName.value = false
}

const showExport = ref(false)

const bps: { key: Breakpoint; icon: any; label: string }[] = [
  { key: 'desktop', icon: IconDeviceDesktop, label: 'Desktop' },
  { key: 'tablet', icon: IconDeviceTablet, label: 'Tablet' },
  { key: 'mobile', icon: IconDeviceMobile, label: 'Mobile' },
]

// Save button label for tooltip + icon state.
const saveState = () => {
  if (store.saving.value) return { label: 'Saving…', icon: IconDeviceFloppy, spin: true }
  if (store.dirty.value) return { label: 'Save', icon: IconDeviceFloppy, spin: false }
  return { label: 'Saved', icon: IconCheck, spin: false }
}

function openPreview() {
  const id = store.page.value?.id
  if (!id) return
  // In dev, backend runs on :3000; in production same origin serves both.
  const base = window.location.port === '5174' || window.location.port === '5173'
    ? `${window.location.protocol}//${window.location.hostname}:3000`
    : ''
  window.open(`${base}/?preview=${id}`, '_blank')
}
</script>

<template>
  <header class="flex h-14 items-stretch border-b border-border bg-card">
    <!-- Left: identity -->
    <div class="flex flex-1 items-center gap-1 pl-3">
      <Button variant="ghost" size="icon" title="Back to pages" @click="router.push('/pages')">
        <IconArrowLeft class="size-4" />
      </Button>

      <div class="flex items-center gap-2">
        <Input
          v-if="editingName"
          v-model="nameBuffer"
          class="h-8 w-48"
          autofocus
          @blur="commitName"
          @keydown.enter.prevent="commitName"
          @keydown.esc.prevent="editingName = false"
        />
        <button
          v-else
          class="rounded px-2 py-1 text-sm font-medium hover:bg-muted"
          @click="startEditName"
        >
          {{ store.page.value?.name ?? '—' }}
        </button>
        <span class="text-xs text-muted-foreground">v{{ store.page.value?.version ?? 0 }}</span>
        <span v-if="store.dirty.value" class="text-xs text-amber-600">●</span>
      </div>
    </div>

    <!-- Center: breakpoint (perfectly centered) -->
    <div class="flex shrink-0 items-center gap-0.5 rounded-lg border border-border p-0.5 my-2">
      <button
        v-for="bp in bps"
        :key="bp.key"
        class="flex size-8 items-center justify-center rounded-md transition-colors"
        :class="store.breakpoint.value === bp.key ? 'bg-primary text-white' : 'text-muted-foreground hover:bg-muted'"
        :title="bp.label"
        @click="store.breakpoint.value = bp.key"
      >
        <component :is="bp.icon" class="size-4" />
      </button>
    </div>

    <!-- Right: panel toggles + actions (icon-only) -->
    <div class="flex flex-1 items-center justify-end gap-0 pr-3">
      <div class="flex items-center gap-0">
        <Button
          variant="ghost"
          size="icon"
          :class="!showLeft ? 'bg-muted' : ''"
          title="Toggle left panel"
          @click="$emit('toggle-left')"
        >
          <IconLayoutSidebar class="size-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          :class="!showRight ? 'bg-muted' : ''"
          title="Toggle right panel"
          @click="$emit('toggle-right')"
        >
          <IconLayoutSidebarRight class="size-4" />
        </Button>
      </div>

      <Button
        variant="ghost"
        size="icon"
        :disabled="!store.canUndo.value"
        title="Undo (Ctrl+Z)"
        @click="store.undo()"
      >
        <IconArrowBackUp class="size-4" />
      </Button>
      <Button
        variant="ghost"
        size="icon"
        :disabled="!store.canRedo.value"
        title="Redo (Ctrl+Shift+Z)"
        @click="store.redo()"
      >
        <IconArrowForwardUp class="size-4" />
      </Button>
      <HistoryDropdown />

      <div class="mx-0.5 h-5 w-px shrink-0 bg-border" />

      <Button
        variant="ghost"
        size="icon"
        :disabled="store.saving.value || !store.dirty.value"
        :title="saveState().label"
        @click="store.save()"
      >
        <component :is="saveState().icon" class="size-4" :class="saveState().spin ? 'animate-spin' : (!store.dirty.value && !store.saving.value ? 'text-green-600' : '')" />
      </Button>
      <Button variant="ghost" size="icon" title="Export" @click="showExport = true">
        <IconCode class="size-4" />
      </Button>
      <Button variant="ghost" size="icon" title="Preview" @click="openPreview">
        <IconExternalLink class="size-4" />
      </Button>
      <Button variant="default" size="icon" :disabled="store.saving.value" title="Publish" @click="store.publish()">
        <IconRocket class="size-4" />
      </Button>
    </div>
  </header>

  <ExportDialog v-if="showExport" @close="showExport = false" />
</template>
