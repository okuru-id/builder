<script setup lang="ts">
import type { ColumnDef } from '@tanstack/vue-table'
import { h, ref, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import DataTable from '@/components/DataTable.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { IconPlus, IconEdit, IconTrash, IconX } from '@tabler/icons-vue'

interface ProjectRow {
  id: number
  title_en: string
  title_id: string
  description_en: string
  description_id: string
  sort_order: number
  url: string
  featured: boolean
  tech_stack: Record<string, boolean>
}
interface Project {
  id?: number
  title_en: string
  title_id: string
  description_en: string
  description_id: string
  sort_order: number
  url: string
  featured: boolean
  tech_stack: string
}
const empty = (): Project => ({
  title_en: '', title_id: '', description_en: '', description_id: '',
  sort_order: 0, url: '', featured: false, tech_stack: '',
})

const projects = ref<ProjectRow[]>([])
const loading = ref(true)
const showForm = ref(false)
const form = ref<Project>(empty())

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

function startCreate() {
  form.value = empty()
  showForm.value = true
}
function startEdit(p: ProjectRow) {
  form.value = {
    id: p.id,
    title_en: p.title_en ?? '',
    title_id: p.title_id ?? '',
    description_en: p.description_en ?? '',
    description_id: p.description_id ?? '',
    sort_order: p.sort_order ?? 0,
    url: p.url ?? '',
    featured: !!p.featured,
    tech_stack: p.tech_stack ? Object.keys(p.tech_stack).join(', ') : '',
  }
  showForm.value = true
}

async function save() {
  if (!form.value.title_en) {
    toast.error('Title (EN) is required')
    return
  }
  const stack = form.value.tech_stack
    .split(',').map((s) => s.trim()).filter(Boolean)
  const payload: any = {
    title_en: form.value.title_en,
    title_id: form.value.title_id,
    description_en: form.value.description_en || null,
    description_id: form.value.description_id || null,
    sort_order: Number(form.value.sort_order) || 0,
    url: form.value.url || null,
    featured: form.value.featured,
    tech_stack: Object.fromEntries(stack.map((s) => [s, true])),
  }
  try {
    if (form.value.id) await api.put(`/projects/${form.value.id}`, payload)
    else await api.post('/projects', payload)
    toast.success('Project saved')
    showForm.value = false
    await load()
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Save failed')
  }
}

async function remove(id: number) {
  if (!confirm('Delete this project?')) return
  try {
    await api.delete(`/projects/${id}`)
    projects.value = projects.value.filter((p) => p.id !== id)
    toast.success('Deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
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
      h(Button, { variant: 'ghost', size: 'icon-sm', onClick: () => startEdit(row.original) }, () => h(IconEdit, { class: 'size-4' })),
      h(Button, { variant: 'ghost', size: 'icon-sm', onClick: () => remove(row.original.id) }, () => h(IconTrash, { class: 'size-4 text-destructive' })),
    ]),
  },
]

onMounted(load)
</script>

<template>
  <div class="p-4 lg:p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="font-heading text-2xl font-bold">Projects</h1>
      <Button @click="startCreate">
        <IconPlus class="size-4" />
        New Project
      </Button>
    </div>

    <div v-if="showForm" class="mb-6 rounded-lg border p-4">
      <div class="mb-4 flex items-center justify-between">
        <h2 class="font-semibold">{{ form.id ? 'Edit Project' : 'New Project' }}</h2>
        <Button variant="ghost" size="icon-sm" @click="showForm = false">
          <IconX class="size-4" />
        </Button>
      </div>
      <form class="grid gap-4 sm:grid-cols-2" @submit.prevent="save">
        <div class="flex flex-col gap-2"><Label>Title (EN) *</Label><Input v-model="form.title_en" /></div>
        <div class="flex flex-col gap-2"><Label>Title (ID)</Label><Input v-model="form.title_id" /></div>
        <div class="flex flex-col gap-2"><Label>Description (EN)</Label><Textarea v-model="form.description_en" /></div>
        <div class="flex flex-col gap-2"><Label>Description (ID)</Label><Textarea v-model="form.description_id" /></div>
        <div class="flex flex-col gap-2"><Label>Sort Order</Label><Input v-model.number="form.sort_order" type="number" /></div>
        <div class="flex flex-col gap-2"><Label>URL</Label><Input v-model="form.url" placeholder="https://…" /></div>
        <div class="flex flex-col gap-2 sm:col-span-2"><Label>Tech Stack (comma separated)</Label><Input v-model="form.tech_stack" placeholder="Go, Vue, Postgres" /></div>
        <div class="flex items-center gap-2 sm:col-span-2">
          <input id="featured" v-model="form.featured" type="checkbox" class="size-4" />
          <Label for="featured">Featured</Label>
        </div>
        <div class="flex gap-2 sm:col-span-2">
          <Button type="submit">Save</Button>
          <Button type="button" variant="ghost" @click="showForm = false">Cancel</Button>
        </div>
      </form>
    </div>

    <DataTable :data="projects" :columns="columns" :loading="loading" search-key="title_en" search-placeholder="Search projects…" />
  </div>
</template>
