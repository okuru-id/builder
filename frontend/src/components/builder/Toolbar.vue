<script setup lang="ts">
// Top toolbar: back, page name, breakpoint switch, save state, publish.
import { inject, ref } from 'vue'
import { useRouter } from 'vue-router'
import { IconArrowLeft, IconDeviceDesktop, IconDeviceTablet, IconDeviceMobile } from '@tabler/icons-vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { BUILDER_KEY } from '@/components/builder/injection'
import type { Breakpoint } from '@/components/builder/useBuilderStore'

const store = inject(BUILDER_KEY, null)!
const router = useRouter()

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

const bps: { key: Breakpoint; icon: any; label: string }[] = [
  { key: 'desktop', icon: IconDeviceDesktop, label: 'Desktop' },
  { key: 'tablet', icon: IconDeviceTablet, label: 'Tablet' },
  { key: 'mobile', icon: IconDeviceMobile, label: 'Mobile' },
]
</script>

<template>
  <header class="flex h-14 items-center gap-3 border-b border-neutral-200 bg-white px-3">
    <Button variant="ghost" size="icon" @click="router.push('/')">
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
        class="rounded px-2 py-1 text-sm font-medium hover:bg-neutral-100"
        @click="startEditName"
      >
        {{ store.page.value?.name ?? '—' }}
      </button>
      <span class="text-xs text-neutral-400">v{{ store.page.value?.version ?? 0 }}</span>
      <span v-if="store.dirty.value" class="text-xs text-amber-600">●</span>
      <span v-if="store.saving.value" class="text-xs text-neutral-400">menyimpan…</span>
    </div>

    <div class="ml-4 flex items-center gap-0.5 rounded-lg border border-neutral-200 p-0.5">
      <button
        v-for="bp in bps"
        :key="bp.key"
        class="flex size-8 items-center justify-center rounded-md transition-colors"
        :class="store.breakpoint.value === bp.key ? 'bg-neutral-900 text-white' : 'text-neutral-500 hover:bg-neutral-100'"
        :title="bp.label"
        @click="store.breakpoint.value = bp.key"
      >
        <component :is="bp.icon" class="size-4" />
      </button>
    </div>

    <div class="ml-auto flex items-center gap-2">
      <Button variant="default" size="sm" :disabled="store.saving.value" @click="store.publish()">
        Publish
      </Button>
    </div>
  </header>
</template>
