<script setup lang="ts">
// Landing pages list: index of builder pages. Entry point to open the canvas.
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { toast } from 'vue-sonner'
import { IconPlus, IconEdit, IconExternalLink, IconUpload, IconCopy, IconTrash } from '@tabler/icons-vue'
import api from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import ImportDialog from '@/components/builder/ImportDialog.vue'

interface Page {
  id: number
  slug: string
  name: string
  status: string
  version: number
  updated_at: string
}

const router = useRouter()
const pages = ref<Page[]>([])
const loading = ref(true)
const showImport = ref(false)

async function load() {
  loading.value = true
  try {
    const res = await api.get<{ data: Page[] }>('/landing-pages')
    pages.value = res.data.data
  } catch (e) {
    toast.error('Failed to load pages')
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function createPage() {
  try {
    const res = await api.post<{ data: Page }>('/landing-pages', { name: 'New page' })
    router.push({ name: 'builder', params: { id: res.data.data.id } })
  } catch (e) {
    toast.error('Failed to create page')
    console.error(e)
  }
}

function open(p: Page) {
  router.push({ name: 'builder', params: { id: p.id } })
}

function preview(p: Page) {
  window.open(`/?preview=${p.id}`, '_blank')
}

async function duplicate(p: Page) {
  try {
    const res = await api.post<{ data: Page }>(`/landing-pages/${p.id}/duplicate`)
    toast.success(`"${p.name}" duplicated`)
    pages.value.unshift(res.data.data)
  } catch (e) {
    toast.error('Failed to duplicate page')
    console.error(e)
  }
}

const showDelete = ref(false)
const deletingPage = ref<Page | null>(null)
const deleteDesc = computed(() => deletingPage.value
  ? `Are you sure you want to delete "${deletingPage.value.name}"? This action cannot be undone.`
  : '')

function confirmDelete(p: Page) {
  deletingPage.value = p
  showDelete.value = true
}
async function doDelete() {
  const p = deletingPage.value
  if (!p) return
  try {
    await api.delete(`/landing-pages/${p.id}`)
    pages.value = pages.value.filter((x) => x.id !== p.id)
    toast.success(`"${p.name}" deleted`)
  } catch (e) {
    toast.error('Failed to delete page')
    console.error(e)
  } finally {
    showDelete.value = false
    deletingPage.value = null
  }
}

function timeAgo(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diff = Math.floor((now.getTime() - d.getTime()) / 1000)
  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  if (diff < 2592000) return `${Math.floor(diff / 86400)}d ago`
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

onMounted(load)
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Landing Pages</h1>
        <p class="text-sm text-muted-foreground">Manage builder canvas landing pages.</p>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" @click="showImport = true">
          <IconUpload class="size-4" /> Import
        </Button>
        <Button @click="createPage">
          <IconPlus class="size-4" /> New page
        </Button>
      </div>
    </div>

    <!-- Loading skeleton grid -->
    <div v-if="loading" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      <Card v-for="i in 8" :key="i">
        <div class="aspect-[4/3] animate-pulse bg-muted" />
        <CardContent class="space-y-2 p-3">
          <div class="h-4 w-2/3 animate-pulse rounded bg-muted" />
          <div class="h-3 w-1/3 animate-pulse rounded bg-muted" />
        </CardContent>
      </Card>
    </div>

    <!-- Empty state -->
    <Card v-else-if="pages.length === 0">
      <CardContent class="flex flex-col items-center gap-3 py-16 text-center">
        <div class="flex size-12 items-center justify-center rounded-full bg-muted">
          <IconPlus class="size-6 text-muted-foreground" />
        </div>
        <div>
          <p class="font-medium">No pages yet</p>
          <p class="text-sm text-muted-foreground">Create your first page.</p>
        </div>
        <Button @click="createPage">
          <IconPlus class="size-4" /> New page
        </Button>
      </CardContent>
    </Card>

    <!-- Gallery grid -->
    <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
      <Card
        v-for="p in pages"
        :key="p.id"
        class="group overflow-hidden py-0 transition-shadow hover:shadow-md"
      >
        <!-- Thumbnail / preview area -->
        <button
          class="relative block aspect-[4/3] w-full overflow-hidden bg-gradient-to-br from-neutral-100 to-neutral-200"
          @click="open(p)"
        >
          <span class="pointer-events-none absolute inset-0 flex items-center justify-center text-2xl font-bold text-neutral-300">
            {{ (p.name || '?').charAt(0).toUpperCase() }}
          </span>
          <span class="absolute left-2 top-2">
            <Badge :variant="p.status === 'published' ? 'default' : 'secondary'" class="backdrop-blur">
              {{ p.status }}
            </Badge>
          </span>
          <!-- Hover overlay actions -->
          <span class="absolute inset-0 flex items-center justify-center gap-1.5 bg-black/40 opacity-0 backdrop-blur-sm transition-opacity group-hover:opacity-100">
            <span
              class="flex size-9 items-center justify-center rounded-md bg-white/95 text-neutral-700 shadow hover:bg-white"
              title="Edit"
              @click.stop="open(p)"
            >
              <IconEdit class="size-4" />
            </span>
            <span
              class="flex size-9 items-center justify-center rounded-md bg-white/95 text-neutral-700 shadow hover:bg-white"
              title="Preview"
              @click.stop="preview(p)"
            >
              <IconExternalLink class="size-4" />
            </span>
            <span
              class="flex size-9 items-center justify-center rounded-md bg-white/95 text-neutral-700 shadow hover:bg-white"
              title="Duplicate"
              @click.stop="duplicate(p)"
            >
              <IconCopy class="size-4" />
            </span>
            <span
              class="flex size-9 items-center justify-center rounded-md bg-white/95 text-red-600 shadow hover:bg-white"
              title="Delete"
              @click.stop="confirmDelete(p)"
            >
              <IconTrash class="size-4" />
            </span>
          </span>
        </button>

        <!-- Footer: name + meta -->
        <CardContent class="space-y-1 p-3">
          <div class="flex items-center gap-2">
            <span class="min-w-0 flex-1 truncate text-sm font-medium">{{ p.name }}</span>
            <span class="shrink-0 text-[11px] text-muted-foreground">v{{ p.version }}</span>
          </div>
          <div class="flex items-center justify-between gap-2">
            <span class="truncate text-xs text-muted-foreground">/{{ p.slug }}</span>
            <span v-if="p.updated_at" class="shrink-0 text-[11px] text-muted-foreground">{{ timeAgo(p.updated_at) }}</span>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>

  <ConfirmDialog v-model:open="showDelete" title="Delete page" :description="deleteDesc" @confirm="doDelete" />
  <ImportDialog v-if="showImport" @close="showImport = false" />
</template>
