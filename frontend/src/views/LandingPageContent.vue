<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import { IconEdit, IconEyeOff, IconEye, IconTrash, IconPlus, IconGripVertical, IconX, IconDeviceFloppy } from '@tabler/icons-vue'
import { Skeleton } from '@/components/ui/skeleton'

let dragIdx: number | null = null

interface LandingSection {
  id: number
  type: string
  content: Record<string, any>
  sort_order: number
  is_active: boolean
}

const sections = ref<LandingSection[]>([])
const loading = ref(true)
const editingId = ref<number | null>(null)
const addingSection = ref(false)
const newSectionType = ref('')
const drafts = ref<Map<number, Record<string, any>>>(new Map())

function draft(section: LandingSection) {
  if (!drafts.value.has(section.id)) {
    drafts.value.set(section.id, JSON.parse(JSON.stringify(section.content)))
  }
  return drafts.value.get(section.id)!
}

function confirmAddKey(section: LandingSection) {
  const c = draft(section)
  const key = c._newKey?.trim()
  if (!key) return
  if (!(key in c)) {
    c[key] = ''
  }
  delete c._newKey
}

function removeDraftKey(section: LandingSection, key: string) {
  const c = draft(section)
  delete c[key]
}

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/landing-sections')
    sections.value = data.data ?? []
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load')
  } finally {
    loading.value = false
  }
}

function openEdit(id: number) {
  editingId.value = id
}

function closeEdit() {
  const id = editingId.value
  if (id !== null) drafts.value.delete(id)
  editingId.value = null
}

async function saveSection(section: LandingSection) {
  const content = cleanDraft(section)
  if (!content) return
  try {
    await api.put(`/landing-sections/${section.id}`, { content })
    Object.assign(section.content, content)
    drafts.value.delete(section.id)
    toast.success('Saved')
    editingId.value = null
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Save failed')
  }
}

function cleanDraft(section: LandingSection) {
  const c = drafts.value.get(section.id)
  if (!c) return null
  const clean = { ...c }
  delete clean._newKey
  return clean
}

async function toggleSection(section: LandingSection) {
  try {
    await api.patch(`/landing-sections/${section.id}/toggle`)
    section.is_active = !section.is_active
    toast.success(section.is_active ? 'Section shown' : 'Section hidden')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Toggle failed')
  }
}

async function deleteSection(section: LandingSection) {
  if (!confirm(`Delete "${section.type}" section?`)) return
  try {
    await api.delete(`/landing-sections/${section.id}`)
    sections.value = sections.value.filter(s => s.id !== section.id)
    toast.success('Deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
  }
}

function addItem(section: LandingSection) {
  const c = draft(section)
  const items = c.items ?? []
  items.push({})
  c.items = [...items]
}

function removeItem(section: LandingSection, index: number) {
  const c = draft(section)
  const items = c.items ?? []
  items.splice(index, 1)
  c.items = [...items]
}

async function sortSections(fromIdx: number, toIdx: number) {
  const copy = [...sections.value]
  const [moved] = copy.splice(fromIdx, 1)
  copy.splice(toIdx, 0, moved)
  copy.forEach((s, i) => s.sort_order = i)
  sections.value = copy
  await Promise.all(copy.map(s =>
    api.patch(`/landing-sections/${s.id}/sort`, { sort_order: s.sort_order })
  ))
}

async function createSection() {
  const type = newSectionType.value.trim()
  if (!type) return
  try {
    const { data } = await api.post('/landing-sections', { type, content: {} })
    sections.value.push(data.data)
    newSectionType.value = ''
    addingSection.value = false
    toast.success('Section created')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Create failed')
  }
}

function onSectionDragStart(index: number, e: DragEvent) {
  dragIdx = index
  if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move'
}

function onSectionDragOver(e: DragEvent) {
  e.preventDefault()
  if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
}

function onSectionDrop(index: number) {
  if (dragIdx === null || dragIdx === index) return
  sortSections(dragIdx, index)
  dragIdx = null
}

function onItemDragStart(index: number, e: DragEvent) {
  dragIdx = index
  if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move'
}

function onItemDragOver(e: DragEvent) {
  e.preventDefault()
  if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
}

function onItemDrop(section: LandingSection, index: number) {
  if (dragIdx === null || dragIdx === index) return
  const c = draft(section)
  const copy = [...(c.items ?? [])]
  const [moved] = copy.splice(dragIdx, 1)
  copy.splice(index, 0, moved)
  c.items = copy
  dragIdx = null
}

const sectionLabel = (type: string) => {
  const labels: Record<string, string> = {
    hero: 'Hero',
    clients: 'Client Logos',
    services: 'Services',
    projects: 'Open Source Projects',
    cta: 'CTA',
  }
  return labels[type] ?? type
}

const sectionSummary = (section: LandingSection) => {
  if (section.type === 'hero') return section.content.greeting_en || section.content.greeting_id || ''
  if (section.type === 'cta') return section.content.heading || ''
  const items = section.content.items
  return items ? `${items.length} item(s)` : ''
}

const itemFields = (type: string) => {
  const fields: Record<string, { key: string; label: string; type: 'text' | 'textarea' }[]> = {
    clients: [
      { key: 'name', label: 'Name', type: 'text' },
      { key: 'logo', label: 'Logo URL', type: 'text' },
    ],
    services: [
      { key: 'title', label: 'Title', type: 'text' },
      { key: 'description_en', label: 'Description (EN)', type: 'textarea' },
      { key: 'description_id', label: 'Description (ID)', type: 'textarea' },
      { key: 'icon', label: 'Icon', type: 'text' },
    ],
    projects: [
      { key: 'title_en', label: 'Title (EN)', type: 'text' },
      { key: 'title_id', label: 'Title (ID)', type: 'text' },
      { key: 'description_en', label: 'Description (EN)', type: 'textarea' },
      { key: 'description_id', label: 'Description (ID)', type: 'textarea' },
      { key: 'github_url', label: 'GitHub URL', type: 'text' },
      { key: 'technologies', label: 'Technologies (comma-separated)', type: 'text' },
    ],
  }
  return fields[type] ?? null
}

const isListType = (type: string) => ['clients', 'services', 'projects'].includes(type)


onMounted(load)
</script>

<template>
  <div class="flex flex-col gap-6 p-8 min-w-0">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <h1 class="font-heading text-2xl font-bold">Landing Content</h1>
      <Button @click="addingSection = true" v-if="!addingSection">
        <IconPlus class="size-4" /> Add Section
      </Button>
    </div>

    <!-- Add Section Form -->
    <Card v-if="addingSection" class="border-dashed">
      <CardContent class="flex items-center gap-3 pt-6">
        <Input v-model="newSectionType" placeholder="Section type name (e.g. about, team, gallery)" class="max-w-xs" @keyup.enter="createSection" />
        <Button size="sm" @click="createSection" :disabled="!newSectionType.trim()">Create</Button>
        <Button variant="ghost" size="sm" @click="addingSection = false; newSectionType = ''">Cancel</Button>
      </CardContent>
    </Card>

    <div v-if="loading" class="flex flex-col gap-4">
      <Skeleton v-for="i in 5" :key="i" class="h-32 w-full rounded-xl" />
    </div>

    <div v-else class="flex flex-col gap-4">
      <Card
        v-for="(section, si) in sections"
        :key="section.id"
        :class="[section.is_active ? '' : 'opacity-60', si === dragIdx ? 'ring-2 ring-primary' : '']"
        draggable="true"
        @dragstart="onSectionDragStart(si, $event)"
        @dragover="onSectionDragOver"
        @drop.prevent="onSectionDrop(si)"
        @dragend="dragIdx = null"
      >
        <CardHeader class="flex flex-row items-center justify-between gap-4 pb-3">
          <div class="flex items-center gap-3">
            <IconGripVertical class="size-4 text-muted-foreground shrink-0 cursor-grab active:cursor-grabbing" />
            <Badge variant="outline" class="text-xs">{{ sectionLabel(section.type) }}</Badge>
            <span class="text-muted-foreground text-sm truncate max-w-64">{{ sectionSummary(section) }}</span>
          </div>
          <div class="flex items-center gap-1">
            <Button variant="ghost" size="icon-sm" @click="toggleSection(section)">
              <IconEye v-if="section.is_active" class="size-4" />
              <IconEyeOff v-else class="size-4" />
            </Button>
            <Button variant="ghost" size="icon-sm" @click="openEdit(section.id)">
              <IconEdit class="size-4" />
            </Button>
            <Button variant="ghost" size="icon-sm" @click="deleteSection(section)">
              <IconTrash class="size-4 text-destructive" />
            </Button>
          </div>
        </CardHeader>

        <!-- Inline Edit -->
        <CardContent v-if="editingId === section.id" class="border-t pt-4 space-y-4">
          <!-- Known non-list types (hero, cta) -->
          <template v-if="section.type === 'hero' || section.type === 'cta'">
            <div v-for="(val, key) in draft(section)" :key="key as string">
              <Label class="text-xs text-muted-foreground">{{ key }}</Label>
              <Input v-if="typeof val === 'string' && val.length < 100" v-model="draft(section)[key]" class="mt-1" />
              <Textarea v-else v-model="draft(section)[key]" class="mt-1" />
            </div>
          </template>

          <!-- Known list types (clients, services, projects) -->
          <template v-if="isListType(section.type)">
            <div class="flex items-center justify-between">
              <Label class="text-sm font-medium">Items</Label>
              <Button variant="outline" size="sm" @click="addItem(section)">
                <IconPlus class="size-4" /> Add
              </Button>
            </div>

            <div
              v-for="(item, idx) in draft(section).items"
              :key="idx"
              :class="['rounded-lg border p-4', idx === dragIdx ? 'ring-2 ring-primary opacity-60' : '']"
              draggable="true"
              @dragstart="onItemDragStart(Number(idx), $event)"
              @dragover="onItemDragOver"
              @drop.prevent="onItemDrop(section, Number(idx))"
              @dragend="dragIdx = null"
            >
              <div class="mb-3 flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <IconGripVertical class="size-4 text-muted-foreground shrink-0 cursor-grab active:cursor-grabbing" />
                  <span class="text-xs font-medium text-muted-foreground">Item {{ Number(idx) + 1 }}</span>
                </div>
                <Button variant="ghost" size="icon-sm" class="size-6" @click="removeItem(section, Number(idx))">
                  <IconTrash class="size-3 text-destructive" />
                </Button>
              </div>
              <div class="flex flex-col gap-3">
                <template v-for="field in itemFields(section.type)" :key="field.key">
                  <div>
                    <Label class="text-xs text-muted-foreground">{{ field.label }}</Label>
                    <Input
                      v-if="field.type === 'text'"
                      :model-value="item[field.key] ?? ''"
                      @update:model-value="item[field.key] = $event; draft(section).items = [...draft(section).items]"
                      class="mt-1"
                    />
                    <Textarea
                      v-else
                      :model-value="item[field.key] ?? ''"
                      @update:model-value="item[field.key] = $event; draft(section).items = [...draft(section).items]"
                      class="mt-1"
                    />
                  </div>
                </template>
                <!-- Custom fields for items in known list types -->
                <div v-for="val2 in Object.keys(item).filter(k => !(itemFields(section.type)?.some(f => f.key === k)))" :key="val2">
                  <Label class="text-xs text-muted-foreground">{{ val2 }}</Label>
                  <Input v-model="item[val2]" class="mt-1" @update:model-value="draft(section).items = [...draft(section).items]" />
                </div>
              </div>
            </div>
          </template>

          <!-- Custom / unknown types -->
          <template v-if="!isListType(section.type) && section.type !== 'hero' && section.type !== 'cta'">
            <div
              v-for="(val, key) in draft(section)"
              :key="key as string"
            >
              <div class="flex items-center gap-2">
                <Label class="text-xs text-muted-foreground flex-1">{{ key }}</Label>
                <Button variant="ghost" size="icon-sm" class="size-5" @click="removeDraftKey(section, key as string)">
                  <IconX class="size-3" />
                </Button>
              </div>
              <Input v-if="typeof val === 'string'" v-model="draft(section)[key]" class="mt-1 mb-2" />
              <Textarea v-else v-model="draft(section)[key]" class="mt-1 mb-2" />
            </div>
            <!-- Add custom field -->
            <div class="flex items-center gap-2 pt-1">
              <Input v-model="draft(section)._newKey" placeholder="New field name" class="max-w-40 text-xs" @keyup.enter="confirmAddKey(section)" />
              <Button variant="outline" size="sm" class="text-xs" @click="confirmAddKey(section)">Add Field</Button>
            </div>
          </template>

          <div class="flex items-center gap-2 pt-2">
            <Button size="sm" @click="saveSection(section)">
              <IconDeviceFloppy class="size-4" /> Save
            </Button>
            <Button variant="ghost" size="sm" @click="closeEdit">
              <IconX class="size-4" /> Cancel
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
