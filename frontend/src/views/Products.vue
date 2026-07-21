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

interface ProductRow {
  id: number
  title: string
  slug: string
  description: string
  price: number
  type: string
  status: string
}

const idr = new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' })

const products = ref<ProductRow[]>([])
const loading = ref(true)
const showDelete = ref(false)
const deletingId = ref(0)

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/products')
    products.value = data.data ?? []
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
    await api.delete(`/products/${deletingId.value}`)
    products.value = products.value.filter((p) => p.id !== deletingId.value)
    toast.success('Deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
  } finally {
    showDelete.value = false
  }
}

const columns: ColumnDef<ProductRow>[] = [
  { accessorKey: 'title', header: 'Title', cell: ({ row }) => h('span', { class: 'font-medium' }, row.getValue('title')) },
  { accessorKey: 'price', header: 'Price', cell: ({ row }) => idr.format(Number(row.getValue('price'))) },
  { accessorKey: 'type', header: 'Type' },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => h(Badge, { variant: row.getValue('status') === 'published' ? 'default' : 'secondary' }, () => String(row.getValue('status'))),
  },
  {
    id: 'actions',
    header: () => h('div', { class: 'text-right' }, 'Actions'),
    enableSorting: false,
    cell: ({ row }) => h('div', { class: 'flex justify-end gap-1' }, [
      h(Button, { variant: 'ghost', size: 'icon-sm', onClick: () => router.push(`/products/${row.original.id}/edit`) }, () => h(IconEdit, { class: 'size-4' })),
      h(Button, { variant: 'ghost', size: 'icon-sm', onClick: () => confirmRemove(row.original.id) }, () => h(IconTrash, { class: 'size-4 text-destructive' })),
    ]),
  },
]

onMounted(load)
</script>

<template>
  <div class="p-4 lg:p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="font-heading text-2xl font-bold">Products</h1>
      <Button @click="router.push('/products/new')">
        <IconPlus class="size-4" />
        New Product
      </Button>
    </div>

    <DataTable :data="products" :columns="columns" :loading="loading" :draggable="false" search-key="title" search-placeholder="Search products…" />
  </div>

  <ConfirmDialog v-model:open="showDelete" title="Delete product" description="Are you sure you want to delete this product?" @confirm="doDelete" />
</template>
