<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { IconPlus, IconEdit, IconTrash, IconX } from '@tabler/icons-vue'

interface OSS {
  id?: number
  title_en: string
  title_id: string
  description_en: string
  description_id: string
  github_url: string
  stars: number
  license: string
  technologies: string
  sort_order: number
}

const empty = (): OSS => ({
  title_en: '', title_id: '', description_en: '', description_id: '',
  github_url: '', stars: 0, license: '', technologies: '', sort_order: 0,
})

const items = ref<any[]>([])
const loading = ref(true)
const showForm = ref(false)
const form = ref<OSS>(empty())

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

function startCreate() {
  form.value = empty()
  showForm.value = true
}

function startEdit(o: any) {
  form.value = {
    id: o.id,
    title_en: o.title_en ?? '',
    title_id: o.title_id ?? '',
    description_en: o.description_en ?? '',
    description_id: o.description_id ?? '',
    github_url: o.github_url ?? '',
    stars: o.stars ?? 0,
    license: o.license ?? '',
    technologies: o.technologies ? Object.keys(o.technologies).join(', ') : '',
    sort_order: o.sort_order ?? 0,
  }
  showForm.value = true
}

async function save() {
  if (!form.value.title_en || !form.value.github_url) {
    toast.error('Title (EN) and GitHub URL are required')
    return
  }
  const techs = form.value.technologies.split(',').map((s) => s.trim()).filter(Boolean)
  const payload: any = {
    title_en: form.value.title_en,
    title_id: form.value.title_id,
    description_en: form.value.description_en || null,
    description_id: form.value.description_id || null,
    github_url: form.value.github_url,
    stars: Number(form.value.stars) || 0,
    license: form.value.license || null,
    technologies: Object.fromEntries(techs.map((t) => [t, true])),
    sort_order: Number(form.value.sort_order) || 0,
  }
  try {
    if (form.value.id) await api.put(`/open-source/${form.value.id}`, payload)
    else await api.post('/open-source', payload)
    toast.success('Saved')
    showForm.value = false
    await load()
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Save failed')
  }
}

async function remove(id: number) {
  if (!confirm('Delete this item?')) return
  try {
    await api.delete(`/open-source/${id}`)
    items.value = items.value.filter((i) => i.id !== id)
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
      <h1 class="font-heading text-2xl font-bold">Open Source</h1>
      <Button @click="startCreate">
        <IconPlus class="size-4" />
        New Item
      </Button>
    </div>

    <div v-if="showForm" class="mb-6 rounded-lg border p-4">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="font-semibold">{{ form.id ? 'Edit Item' : 'New Item' }}</h2>
        <Button variant="ghost" size="icon-sm" @click="showForm = false">
          <IconX class="size-4" />
        </Button>
      </div>
      <form class="grid gap-4 sm:grid-cols-2" @submit.prevent="save">
        <div class="flex flex-col gap-2"><Label>Title (EN) *</Label><Input v-model="form.title_en" /></div>
        <div class="flex flex-col gap-2"><Label>Title (ID)</Label><Input v-model="form.title_id" /></div>
        <div class="flex flex-col gap-2"><Label>Description (EN)</Label><Textarea v-model="form.description_en" /></div>
        <div class="flex flex-col gap-2"><Label>Description (ID)</Label><Textarea v-model="form.description_id" /></div>
        <div class="flex flex-col gap-2"><Label>GitHub URL *</Label><Input v-model="form.github_url" /></div>
        <div class="flex flex-col gap-2"><Label>License</Label><Input v-model="form.license" placeholder="MIT" /></div>
        <div class="flex flex-col gap-2"><Label>Stars</Label><Input v-model.number="form.stars" type="number" /></div>
        <div class="flex flex-col gap-2"><Label>Sort Order</Label><Input v-model.number="form.sort_order" type="number" /></div>
        <div class="flex flex-col gap-2 sm:col-span-2">
          <Label>Technologies (comma separated)</Label>
          <Input v-model="form.technologies" placeholder="Go, React, Redis" />
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
            <TableHead>Stars</TableHead>
            <TableHead>License</TableHead>
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-if="loading">
            <TableCell :colspan="4"><Skeleton class="h-6 w-full" /></TableCell>
          </TableRow>
          <TableRow v-else-if="items.length === 0">
            <TableCell :colspan="4" class="text-center text-muted-foreground">No items.</TableCell>
          </TableRow>
          <TableRow v-for="o in items" :key="o.id">
            <TableCell>
              <a :href="o.github_url" target="_blank" class="font-medium hover:underline">{{ o.title_en }}</a>
            </TableCell>
            <TableCell>{{ o.stars }}</TableCell>
            <TableCell>{{ o.license || '—' }}</TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-1">
                <Button variant="ghost" size="icon-sm" @click="startEdit(o)"><IconEdit class="size-4" /></Button>
                <Button variant="ghost" size="icon-sm" @click="remove(o.id)"><IconTrash class="size-4 text-destructive" /></Button>
              </div>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
