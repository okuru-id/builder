<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { IconPlus, IconEdit, IconTrash, IconX } from '@tabler/icons-vue'

const idr = new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' })
const types = ['ebook', 'template', 'source_code', 'donation']

interface Product {
  id?: number
  title: string
  slug: string
  description: string
  price: number
  type: string
  status: string
}
const empty = (): Product => ({
  title: '', slug: '', description: '', price: 0, type: 'ebook', status: 'draft',
})

const products = ref<any[]>([])
const loading = ref(true)
const showForm = ref(false)
const form = ref<Product>(empty())

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

function startCreate() {
  form.value = empty()
  showForm.value = true
}
function startEdit(p: any) {
  form.value = {
    id: p.id,
    title: p.title ?? '',
    slug: p.slug ?? '',
    description: p.description ?? '',
    price: p.price ?? 0,
    type: p.type ?? 'ebook',
    status: p.status ?? 'draft',
  }
  showForm.value = true
}

async function save() {
  if (!form.value.title || !form.value.slug) {
    toast.error('Title and Slug are required')
    return
  }
  const payload: any = {
    title: form.value.title,
    slug: form.value.slug,
    description: form.value.description || null,
    price: Number(form.value.price) || 0,
    type: form.value.type,
    status: form.value.status,
  }
  try {
    if (form.value.id) await api.put(`/products/${form.value.id}`, payload)
    else await api.post('/products', payload)
    toast.success('Saved')
    showForm.value = false
    await load()
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Save failed')
  }
}

async function remove(id: number) {
  if (!confirm('Delete this product?')) return
  try {
    await api.delete(`/products/${id}`)
    products.value = products.value.filter((p) => p.id !== id)
    toast.success('Deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
  }
}

onMounted(load)
</script>

<template>
  <div class="p-8">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="font-heading text-2xl font-bold">Products</h1>
      <Button @click="startCreate">
        <IconPlus class="size-4" />
        New Product
      </Button>
    </div>

    <div v-if="showForm" class="mb-6 rounded-lg border p-4">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="font-semibold">{{ form.id ? 'Edit Product' : 'New Product' }}</h2>
        <Button variant="ghost" size="icon-sm" @click="showForm = false">
          <IconX class="size-4" />
        </Button>
      </div>
      <form class="grid gap-4 sm:grid-cols-2" @submit.prevent="save">
        <div class="flex flex-col gap-2"><Label>Title *</Label><Input v-model="form.title" /></div>
        <div class="flex flex-col gap-2"><Label>Slug *</Label><Input v-model="form.slug" /></div>
        <div class="flex flex-col gap-2 sm:col-span-2"><Label>Description</Label><Textarea v-model="form.description" /></div>
        <div class="flex flex-col gap-2">
          <Label>Price (IDR)</Label>
          <Input v-model.number="form.price" type="number" min="0" />
        </div>
        <div class="flex flex-col gap-2">
          <Label>Type</Label>
          <Select v-model="form.type">
            <SelectTrigger class="w-full"><SelectValue /></SelectTrigger>
            <SelectContent>
              <SelectItem v-for="t in types" :key="t" :value="t">{{ t }}</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div class="flex flex-col gap-2">
          <Label>Status</Label>
          <Select v-model="form.status">
            <SelectTrigger class="w-full"><SelectValue /></SelectTrigger>
            <SelectContent>
              <SelectItem value="draft">draft</SelectItem>
              <SelectItem value="published">published</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div class="flex gap-2 sm:col-span-2">
          <Button type="submit">Save</Button>
          <Button type="button" variant="ghost" @click="showForm = false">Cancel</Button>
        </div>
      </form>
    </div>

    <div class="rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Title</TableHead>
            <TableHead>Price</TableHead>
            <TableHead>Type</TableHead>
            <TableHead>Status</TableHead>
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-if="loading">
            <TableCell :colspan="5"><Skeleton class="h-6 w-full" /></TableCell>
          </TableRow>
          <TableRow v-else-if="products.length === 0">
            <TableCell :colspan="5" class="text-center text-muted-foreground">No products.</TableCell>
          </TableRow>
          <TableRow v-for="p in products" :key="p.id">
            <TableCell class="font-medium">{{ p.title }}</TableCell>
            <TableCell>{{ idr.format(p.price) }}</TableCell>
            <TableCell>{{ p.type }}</TableCell>
            <TableCell>
              <Badge :variant="p.status === 'published' ? 'default' : 'secondary'">{{ p.status }}</Badge>
            </TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-1">
                <Button variant="ghost" size="icon-sm" @click="startEdit(p)"><IconEdit class="size-4" /></Button>
                <Button variant="ghost" size="icon-sm" @click="remove(p.id)"><IconTrash class="size-4 text-destructive" /></Button>
              </div>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
