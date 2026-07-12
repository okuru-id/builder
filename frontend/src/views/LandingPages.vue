<script setup lang="ts">
// Landing pages list: index of builder pages. Entry point to open the canvas.
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { IconPlus, IconEdit, IconExternalLink, IconUpload } from '@tabler/icons-vue'
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
    toast.error('Gagal memuat halaman')
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function createPage() {
  try {
    const res = await api.post<{ data: Page }>('/landing-pages', { name: 'Halaman baru' })
    router.push({ name: 'builder', params: { id: res.data.data.id } })
  } catch (e) {
    toast.error('Gagal membuat halaman')
    console.error(e)
  }
}

function open(p: Page) {
  router.push({ name: 'builder', params: { id: p.id } })
}

function preview(p: Page) {
  window.open(`/?preview=${p.id}`, '_blank')
}

onMounted(load)
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Landing Pages</h1>
        <p class="text-sm text-muted-foreground">Kelola halaman landing builder canvas.</p>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" @click="showImport = true">
          <IconUpload class="size-4" /> Import
        </Button>
        <Button @click="createPage">
          <IconPlus class="size-4" /> Halaman baru
        </Button>
      </div>
    </div>

    <Card v-if="loading">
      <CardContent class="py-10 text-center text-sm text-muted-foreground">Memuat…</CardContent>
    </Card>

    <Card v-else-if="pages.length === 0">
      <CardContent class="py-10 text-center text-sm text-muted-foreground">
        Belum ada halaman. Buat halaman pertama Anda.
      </CardContent>
    </Card>

    <div v-else class="grid gap-3">
      <Card v-for="p in pages" :key="p.id">
        <CardContent class="flex items-center gap-3 py-3">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2">
              <span class="truncate font-medium">{{ p.name }}</span>
              <Badge :variant="p.status === 'published' ? 'default' : 'secondary'">
                {{ p.status }}
              </Badge>
              <span class="text-xs text-muted-foreground">v{{ p.version }}</span>
            </div>
            <div class="text-xs text-muted-foreground">/{{ p.slug }}</div>
          </div>
          <Button variant="outline" size="sm" @click="open(p)">
            <IconEdit class="size-4" /> Edit
          </Button>
          <Button variant="ghost" size="icon" @click="preview(p)">
            <IconExternalLink class="size-4" />
          </Button>
        </CardContent>
      </Card>
    </div>
  </div>

  <ImportDialog v-if="showImport" @close="showImport = false" />
</template>
