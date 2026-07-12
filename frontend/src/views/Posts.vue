<script setup lang="ts">
import type { ColumnDef } from '@tanstack/vue-table'
import { h, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import DataTable from '@/components/DataTable.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { IconPlus, IconEdit, IconTrash } from '@tabler/icons-vue'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'

interface Post {
  id: number
  title_en: string
  title_id: string
  status: string
  published_at: string | null
}

const router = useRouter()
const posts = ref<Post[]>([])
const loading = ref(true)

function fmtDate(s: string | null) {
  if (!s) return '—'
  return new Date(s).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' })
}

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/posts')
    posts.value = data.data ?? []
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load posts')
  } finally {
    loading.value = false
  }
}

const showDelete = ref(false)
const deletingId = ref(0)

function confirmRemove(id: number) {
  deletingId.value = id
  showDelete.value = true
}
async function doDelete() {
  try {
    await api.delete(`/posts/${deletingId.value}`)
    posts.value = posts.value.filter((p) => p.id !== deletingId.value)
    toast.success('Post deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
  } finally {
    showDelete.value = false
  }
}

const columns: ColumnDef<Post>[] = [
  {
    accessorKey: 'title_en',
    header: 'Title (EN)',
    cell: ({ row }) => h('span', { class: 'font-medium' }, row.getValue('title_en') || '—'),
  },
  {
    accessorKey: 'title_id',
    header: 'Title (ID)',
    cell: ({ row }) => h('span', { class: 'text-muted-foreground' }, row.getValue('title_id') || '—'),
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => h(Badge, {
      variant: row.getValue('status') === 'published' ? 'default' : 'secondary',
    }, () => String(row.getValue('status'))),
  },
  {
    accessorKey: 'published_at',
    header: 'Published',
    cell: ({ row }) => fmtDate(row.getValue('published_at')),
  },
  {
    id: 'actions',
    header: () => h('div', { class: 'text-right' }, 'Actions'),
    enableSorting: false,
    cell: ({ row }) => h('div', { class: 'flex justify-end gap-1' }, [
      h(Button, {
        variant: 'ghost', size: 'icon-sm',
        onClick: () => router.push(`/posts/${row.original.id}/edit`),
      }, () => h(IconEdit, { class: 'size-4' })),
      h(Button, {
        variant: 'ghost', size: 'icon-sm',
        onClick: () => confirmRemove(row.original.id),
      }, () => h(IconTrash, { class: 'size-4 text-destructive' })),
    ]),
  },
]

onMounted(load)
</script>

<template>
  <div class="p-4 lg:p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="font-heading text-2xl font-bold">Blog Posts</h1>
      <Button @click="router.push('/posts/new')">
        <IconPlus class="size-4" />
        New Post
      </Button>
    </div>

    <DataTable :data="posts" :columns="columns" :loading="loading" search-key="title_en" search-placeholder="Search posts…" />
  </div>

  <ConfirmDialog v-model:open="showDelete" title="Delete post" description="Are you sure you want to delete this post?" @confirm="doDelete" />
</template>
