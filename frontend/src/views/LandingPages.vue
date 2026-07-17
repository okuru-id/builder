<script setup lang="ts">
// Landing pages list: index of builder pages. Entry point to open the canvas.
// Each card shows the page's public routing (home / path / domain / draft-only)
// and exposes a Settings dialog for editing path, domain, and is_home.
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { toast } from 'vue-sonner'
import {
  IconPlus, IconEdit, IconExternalLink, IconUpload, IconCopy, IconTrash,
  IconSettings, IconHome, IconWorld, IconLink, IconCheck,
} from '@tabler/icons-vue'
import api from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import ImportDialog from '@/components/builder/ImportDialog.vue'

interface Page {
  id: number
  slug: string
  name: string
  status: string
  version: number
  updated_at: string
  path?: string
  domain?: string
  is_home?: boolean
}

const router = useRouter()
const pages = ref<Page[]>([])
const loading = ref(true)
const showImport = ref(false)

const origin = typeof window !== 'undefined' ? window.location.origin : ''

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

// Public URL: where this page is reachable. Falls back to draft preview.
function publicURL(p: Page): string {
  if (p.domain) return `https://${p.domain}`
  if (p.is_home) return `${origin}/`
  if (p.path) return `${origin}/${p.path}`
  return `${origin}/?preview=${p.id}`
}

function preview(p: Page) {
  window.open(publicURL(p), '_blank')
}

// Short label for the card footer.
function routeLabel(p: Page): string {
  if (p.domain) return p.domain
  if (p.is_home) return '/ (home)'
  if (p.path) return `/${p.path}`
  return `draft · /?preview=${p.id}`
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

// ── Settings dialog (path / domain / is_home) ──────────────────────────
const showSettings = ref(false)
const settingsPage = ref<Page | null>(null)
const settingsPath = ref('')
const settingsDomain = ref('')
const settingsIsHome = ref(false)
const settingsSaving = ref(false)

function openSettings(p: Page) {
  settingsPage.value = p
  settingsPath.value = p.path ?? ''
  settingsDomain.value = p.domain ?? ''
  settingsIsHome.value = !!p.is_home
  showSettings.value = true
}

async function saveSettings() {
  const p = settingsPage.value
  if (!p) return
  settingsSaving.value = true
  try {
    const res = await api.put<{ data: Page }>(`/landing-pages/${p.id}/settings`, {
      path: settingsPath.value,
      domain: settingsDomain.value,
      is_home: settingsIsHome.value,
    })
    // Patch in place + if this became home, clear others' is_home locally.
    Object.assign(p, res.data.data)
    if (settingsIsHome.value) {
      for (const x of pages.value) if (x.id !== p.id) x.is_home = false
    }
    toast.success('Settings saved')
    showSettings.value = false
  } catch (e: any) {
    toast.error(e?.response?.data?.error || 'Failed to save settings')
    console.error(e)
  } finally {
    settingsSaving.value = false
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
        <p class="text-sm text-muted-foreground">Publish to a path, your own domain, or as the site home.</p>
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
          <span class="absolute left-2 top-2 flex items-center gap-1">
            <Badge :variant="p.status === 'published' ? 'default' : 'secondary'" class="backdrop-blur">
              {{ p.status }}
            </Badge>
            <Badge v-if="p.is_home" variant="default" class="backdrop-blur bg-indigo-500 hover:bg-indigo-500">
              home
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
              title="Open public URL"
              @click.stop="preview(p)"
            >
              <IconExternalLink class="size-4" />
            </span>
            <span
              class="flex size-9 items-center justify-center rounded-md bg-white/95 text-neutral-700 shadow hover:bg-white"
              title="Routing & domain"
              @click.stop="openSettings(p)"
            >
              <IconSettings class="size-4" />
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

        <!-- Footer: name + public route -->
        <CardContent class="space-y-1 p-3">
          <div class="flex items-center gap-2">
            <span class="min-w-0 flex-1 truncate text-sm font-medium">{{ p.name }}</span>
            <span class="shrink-0 text-[11px] text-muted-foreground">v{{ p.version }}</span>
          </div>
          <div class="flex items-center justify-between gap-2">
            <span class="flex min-w-0 items-center gap-1 truncate text-xs text-muted-foreground">
              <component
                :is="p.domain ? IconWorld : (p.is_home ? IconHome : (p.path ? IconLink : IconExternalLink))"
                class="size-3 shrink-0"
              />
              <span class="truncate">{{ routeLabel(p) }}</span>
            </span>
            <span v-if="p.updated_at" class="shrink-0 text-[11px] text-muted-foreground">{{ timeAgo(p.updated_at) }}</span>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>

  <ConfirmDialog v-model:open="showDelete" title="Delete page" :description="deleteDesc" @confirm="doDelete" />
  <ImportDialog v-if="showImport" @close="showImport = false" />

  <!-- Settings dialog: routing + domain -->
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
    appear
  >
    <div
      v-if="showSettings"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
      @click.self="showSettings = false"
    >
      <div class="w-full max-w-md rounded-xl border border-border bg-background shadow-2xl">
        <div class="flex items-center justify-between border-b border-border px-5 py-3">
          <h2 class="text-sm font-semibold">Routing &amp; Domain</h2>
          <Button variant="ghost" size="icon-sm" @click="showSettings = false">
            <IconPlus class="size-4 rotate-45" />
          </Button>
        </div>

        <div class="space-y-4 px-5 py-4">
          <p class="text-xs text-muted-foreground">
            Publish this page to a path (<code class="rounded bg-muted px-1">/promo</code>), a custom domain, or as the site home. Only published pages are served publicly.
          </p>

          <div class="space-y-1.5">
            <Label class="text-xs">Path</Label>
            <div class="flex items-center gap-2">
              <span class="text-xs text-muted-foreground">{{ origin }}/</span>
              <Input
                v-model="settingsPath"
                class="h-8 flex-1 font-mono text-xs"
                placeholder="promo"
              />
            </div>
            <p class="text-[11px] text-muted-foreground">Leave empty if using domain or home.</p>
          </div>

          <div class="space-y-1.5">
            <Label class="text-xs">Custom domain</Label>
            <Input
              v-model="settingsDomain"
              class="h-8 font-mono text-xs"
              placeholder="promo.example.com"
            />
            <p class="text-[11px] text-muted-foreground">
              Point this domain's DNS to the server. Requests with this Host header serve this page.
            </p>
          </div>

          <label class="flex cursor-pointer items-center justify-between rounded-lg border border-border px-3 py-2">
            <span class="flex items-center gap-2 text-xs">
              <IconHome class="size-3.5 text-muted-foreground" />
              Serve as site home (<code class="rounded bg-muted px-1">/</code>)
            </span>
            <input v-model="settingsIsHome" type="checkbox" class="size-3.5" />
          </label>
          <p v-if="settingsIsHome" class="text-[11px] text-amber-600">
            Only one page can be home. Setting this will unset home on all other pages.
          </p>
        </div>

        <div class="flex items-center justify-between gap-2 border-t border-border bg-muted/30 px-4 py-3">
          <Button
            v-if="settingsPage && (settingsPage.domain || settingsPage.is_home || settingsPage.path)"
            variant="ghost"
            size="sm"
            class="h-8 text-muted-foreground"
            @click="settingsPath = ''; settingsDomain = ''; settingsIsHome = false"
          >
            Clear routing
          </Button>
          <div v-else></div>
          <div class="flex items-center gap-2">
            <Button variant="outline" size="sm" class="h-8" @click="showSettings = false">Cancel</Button>
            <Button size="sm" class="h-8" :disabled="settingsSaving" @click="saveSettings">
              <component :is="settingsSaving ? IconCheck : IconCheck" class="size-3.5" />
              {{ settingsSaving ? 'Saving…' : 'Save' }}
            </Button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>
