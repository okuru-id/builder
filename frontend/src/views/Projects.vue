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

const router = useRouter()

interface ProjectRow {
  id: number
  title_en: string
  sort_order: number
  featured: boolean
}

const projects = ref<ProjectRow[]>([])
const loading = ref(true)
const showDelete = ref(false)
const deletingId = ref(0)

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/projects')
    projects.value = data.data ?? []
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load')
  } finally {
    loading.value = false
  }
}

function confirmRemove(id: number) {
  deletingId.value = id
  showDelete.value = true
}
async function doDelete() {
  try {
    await api.delete(`/projects/${deletingId.value}`)
    projects.value = projects.value.filter((p) => p.id !== deletingId.value)
    toast.success('Deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
  } finally {
    showDelete.value = false
  }
}

const columns: ColumnDef<ProjectRow>[] = [
  { accessorKey: 'title_en', header: 'Title', cell: ({ row }) => h('span', { class: 'font-medium' }, row.getValue('title_en')) },
  { accessorKey: 'sort_order', header: 'Sort' },
  {
    accessorKey: 'featured',
    header: 'Featured',
    cell: ({ row }) => row.getValue('featured')
      ? h(Badge, { variant: 'default' }, () => 'featured')
      : h('span', { class: 'text-muted-foreground' }, '—'),
  },
  {
    id: 'actions',
    header: () => h('div', { class: 'text-right' }, 'Actions'),
    enableSorting: false,
    cell: ({ row }) => h('div', { class: 'flex justify-end gap-1' }, [
      h(Button, { variant: 'ghost', size: 'icon-sm', onClick: () => router.push(`/projects/${row.original.id}/edit`) }, () => h(IconEdit, { class: 'size-4' })),
      h(Button, { variant: 'ghost', size: 'icon-sm', onClick: () => confirmRemove(row.original.id) }, () => h(IconTrash, { class: 'size-4 text-destructive' })),
    ]),
  },
]

onMounted(load)
</script>

<template>
  <div class="p-4 lg:p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="font-heading text-2xl font-bold">Projects</h1>
      <Button @click="router.push('/projects/new')">
        <IconPlus class="size-4" />
        New Project
      </Button>
    </div>

    <DataTable :data="projects" :columns="columns" :loading="loading" search-key="title_en" search-placeholder="Search projects…" />
  </div>

  <ConfirmDialog v-model:open="showDelete" title="Delete project" description="Are you sure you want to delete this project?" @confirm="doDelete" />
</template>
