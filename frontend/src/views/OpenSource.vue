<script setup lang="ts">
import type { ColumnDef } from '@tanstack/vue-table'
import { h, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import DataTable from '@/components/DataTable.vue'
import { Button } from '@/components/ui/button'
import { IconPlus, IconEdit, IconTrash } from '@tabler/icons-vue'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'

const router = useRouter()

interface OSSRow {
  id: number
  title_en: string
  github_url: string
  stars: number
  license: string
}

const items = ref<OSSRow[]>([])
const loading = ref(true)
const showDelete = ref(false)
const deletingId = ref(0)

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/open-source')
    items.value = data.data ?? []
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
    await api.delete(`/open-source/${deletingId.value}`)
    items.value = items.value.filter((i) => i.id !== deletingId.value)
    toast.success('Deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
  } finally {
    showDelete.value = false
  }
}

const columns: ColumnDef<OSSRow>[] = [
  {
    accessorKey: 'title_en',
    header: 'Title',
    cell: ({ row }) => h('a', {
      href: row.original.github_url, target: '_blank', class: 'font-medium hover:underline',
    }, row.getValue('title_en')),
  },
  { accessorKey: 'stars', header: 'Stars' },
  {
    accessorKey: 'license',
    header: 'License',
    cell: ({ row }) => String(row.getValue('license') || '—'),
  },
  {
    id: 'actions',
    header: () => h('div', { class: 'text-right' }, 'Actions'),
    enableSorting: false,
    cell: ({ row }) => h('div', { class: 'flex justify-end gap-1' }, [
      h(Button, { variant: 'ghost', size: 'icon-sm', onClick: () => router.push(`/open-source/${row.original.id}/edit`) }, () => h(IconEdit, { class: 'size-4' })),
      h(Button, { variant: 'ghost', size: 'icon-sm', onClick: () => confirmRemove(row.original.id) }, () => h(IconTrash, { class: 'size-4 text-destructive' })),
    ]),
  },
]

onMounted(load)
</script>

<template>
  <div class="p-4 lg:p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="font-heading text-2xl font-bold">Open Source</h1>
      <Button @click="router.push('/open-source/new')">
        <IconPlus class="size-4" />
        New Item
      </Button>
    </div>

    <DataTable :data="items" :columns="columns" :loading="loading" search-key="title_en" search-placeholder="Search items…" />
  </div>

  <ConfirmDialog v-model:open="showDelete" title="Delete item" description="Are you sure you want to delete this item?" @confirm="doDelete" />
</template>
