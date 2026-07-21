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
  github_url: '', stars: 0, license: '', technologies: '', sort_order: 0,
})

async function load() {
  if (!isEdit.value) return
  loading.value = true
  try {
    const { data } = await api.get(`/open-source/${id.value}`)
    const o = data.data
    form.value = {
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
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load')
    router.push('/open-source')
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!form.value.title_en || !form.value.github_url) {
    toast.warning('Title (EN) and GitHub URL are required')
    return
  }
  saving.value = true
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
    if (isEdit.value) await api.put(`/open-source/${id.value}`, payload)
    else await api.post('/open-source', payload)
    toast.success('Saved')
    router.push('/open-source')
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
    <Button variant="ghost" size="sm" class="mb-4 -ml-2" @click="router.push('/open-source')">
      <IconArrowLeft class="size-4" />
      Back
    </Button>

    <h1 class="font-heading mb-6 text-2xl font-bold">{{ isEdit ? 'Edit Item' : 'New Item' }}</h1>

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
      <Label class="text-right whitespace-nowrap">GitHub URL *</Label>
      <Input v-model="form.github_url" />
      <Label class="text-right whitespace-nowrap">License</Label>
      <Input v-model="form.license" placeholder="MIT" />
      <Label class="text-right whitespace-nowrap">Stars</Label>
      <Input v-model.number="form.stars" type="number" />
      <Label class="text-right whitespace-nowrap">Sort Order</Label>
      <Input v-model.number="form.sort_order" type="number" />
      <Label class="text-right whitespace-nowrap">Technologies</Label>
      <Input v-model="form.technologies" placeholder="Go, React, Redis" />
      <div></div>
      <div class="flex gap-2">
        <Button type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</Button>
        <Button type="button" variant="ghost" @click="router.push('/open-source')">Cancel</Button>
      </div>
    </form>
  </div>
</template>
