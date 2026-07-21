<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Skeleton } from '@/components/ui/skeleton'
import { IconArrowLeft } from '@tabler/icons-vue'

const route = useRoute()
const router = useRouter()
const id = computed(() => (route.params.id as string) || null)
const isEdit = computed(() => !!id.value)
const loading = ref(false)
const saving = ref(false)

const form = ref({
  title_en: '', title_id: '', description_en: '', description_id: '',
  sort_order: 0, url: '', featured: false, tech_stack: '',
})

async function load() {
  if (!isEdit.value) return
  loading.value = true
  try {
    const { data } = await api.get(`/projects/${id.value}`)
    const p = data.data
    form.value = {
      title_en: p.title_en ?? '',
      title_id: p.title_id ?? '',
      description_en: p.description_en ?? '',
      description_id: p.description_id ?? '',
      sort_order: p.sort_order ?? 0,
      url: p.url ?? '',
      featured: !!p.featured,
      tech_stack: p.tech_stack ? Object.keys(p.tech_stack).join(', ') : '',
    }
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load')
    router.push('/projects')
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!form.value.title_en) {
    toast.warning('Title (EN) is required')
    return
  }
  saving.value = true
  const stack = form.value.tech_stack.split(',').map((s) => s.trim()).filter(Boolean)
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
    if (isEdit.value) await api.put(`/projects/${id.value}`, payload)
    else await api.post('/projects', payload)
    toast.success('Saved')
    router.push('/projects')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Save failed')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="p-4 lg:p-6">
    <Button variant="ghost" size="sm" class="mb-4 -ml-2" @click="router.push('/projects')">
      <IconArrowLeft class="size-4" />
      Back
    </Button>

    <h1 class="font-heading mb-6 text-2xl font-bold">{{ isEdit ? 'Edit Project' : 'New Project' }}</h1>

    <div v-if="loading" class="max-w-2xl space-y-4">
      <Skeleton v-for="i in 4" :key="i" class="h-10 w-full" />
    </div>

    <form v-else class="grid max-w-xl grid-cols-[max-content_1fr] items-center gap-x-3 gap-y-2" @submit.prevent="save">
      <Label class="text-right whitespace-nowrap">Title (EN) *</Label>
      <Input v-model="form.title_en" />
      <Label class="text-right whitespace-nowrap">Title (ID)</Label>
      <Input v-model="form.title_id" />
      <Label class="self-start pt-2 text-right whitespace-nowrap">Description (EN)</Label>
      <Textarea v-model="form.description_en" />
      <Label class="self-start pt-2 text-right whitespace-nowrap">Description (ID)</Label>
      <Textarea v-model="form.description_id" />
      <Label class="text-right whitespace-nowrap">Sort Order</Label>
      <Input v-model.number="form.sort_order" type="number" />
      <Label class="text-right whitespace-nowrap">URL</Label>
      <Input v-model="form.url" placeholder="https://…" />
      <Label class="text-right whitespace-nowrap">Tech Stack</Label>
      <Input v-model="form.tech_stack" placeholder="Go, Vue, Postgres" />
      <Label class="text-right whitespace-nowrap">Featured</Label>
      <div class="flex items-center gap-2">
        <input id="featured" v-model="form.featured" type="checkbox" class="size-4" />
        <Label for="featured" class="text-sm font-normal text-muted-foreground">show on homepage</Label>
      </div>
      <div></div>
      <div class="flex gap-2">
        <Button type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</Button>
        <Button type="button" variant="ghost" @click="router.push('/projects')">Cancel</Button>
      </div>
    </form>
  </div>
</template>
