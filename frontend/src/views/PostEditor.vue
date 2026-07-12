<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { IconArrowLeft } from '@tabler/icons-vue'

const route = useRoute()
const router = useRouter()
const id = computed(() => (route.params.id as string) || null)
const isEdit = computed(() => !!id.value)
const loading = ref(false)
const saving = ref(false)
const categories = ref<any[]>([])

const form = ref({
  title_en: '',
  title_id: '',
  slug: '',
  excerpt_en: '',
  excerpt_id: '',
  content_en: '',
  content_id: '',
  category: '',
  status: 'draft',
  published_at: '',
})

function toLocalInput(s: string | null) {
  if (!s) return ''
  const d = new Date(s)
  if (isNaN(d.getTime())) return ''
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

async function load() {
  loading.value = true
  try {
    const [{ data: catData }] = await Promise.all([
      api.get('/categories'),
      ...(isEdit.value
        ? [api.get(`/posts/${id.value}`).then(({ data }) => {
            const p = data.data
            form.value = {
              title_en: p.title_en ?? '',
              title_id: p.title_id ?? '',
              slug: p.slug ?? '',
              excerpt_en: p.excerpt_en ?? '',
              excerpt_id: p.excerpt_id ?? '',
              content_en: p.content_en ?? '',
              content_id: p.content_id ?? '',
              category: p.category ?? '',
              status: p.status ?? 'draft',
              published_at: toLocalInput(p.published_at),
            }
          })]
        : []),
    ])
    categories.value = catData.data ?? []
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load')
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!form.value.title_en || !form.value.slug) {
    toast.warning('Title (EN) and Slug are required')
    return
  }
  saving.value = true
  const payload: any = { ...form.value }
  if (payload.published_at) {
    payload.published_at = new Date(payload.published_at).toISOString()
  } else {
    payload.published_at = null
  }
  try {
    if (isEdit.value) {
      await api.put(`/posts/${id.value}`, payload)
    } else {
      await api.post('/posts', payload)
    }
    toast.success('Post saved')
    router.push('/posts')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Save failed')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="p-8">
    <Button variant="ghost" size="sm" class="mb-4 -ml-2" @click="router.push('/posts')">
      <IconArrowLeft class="size-4" />
      Back
    </Button>
    <h1 class="mb-6 font-heading text-2xl font-bold">
      {{ isEdit ? 'Edit Post' : 'New Post' }}
    </h1>

    <div v-if="loading" class="flex flex-col gap-3">
      <Skeleton v-for="i in 4" :key="i" class="h-10 w-full" />
    </div>
    <form v-else class="grid max-w-3xl gap-4" @submit.prevent="save">
      <div class="grid gap-4 sm:grid-cols-2">
        <div class="flex flex-col gap-2">
          <Label for="title_en">Title (EN) *</Label>
          <Input id="title_en" v-model="form.title_en" />
        </div>
        <div class="flex flex-col gap-2">
          <Label for="title_id">Title (ID)</Label>
          <Input id="title_id" v-model="form.title_id" />
        </div>
      </div>

      <div class="grid gap-4 sm:grid-cols-2">
        <div class="flex flex-col gap-2">
          <Label for="slug">Slug *</Label>
          <Input id="slug" v-model="form.slug" placeholder="my-post" />
        </div>
        <div class="flex flex-col gap-2">
          <Label for="category">Category</Label>
          <Select v-model="form.category">
            <SelectTrigger class="w-full">
              <SelectValue placeholder="Select category" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="c in categories" :key="c.id" :value="c.slug">
                {{ c.name_en }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <div class="grid gap-4 sm:grid-cols-2">
        <div class="flex flex-col gap-2">
          <Label for="status">Status</Label>
          <Select v-model="form.status">
            <SelectTrigger class="w-full">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="draft">draft</SelectItem>
              <SelectItem value="published">published</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div class="flex flex-col gap-2">
          <Label for="published_at">Published At</Label>
          <Input id="published_at" v-model="form.published_at" type="datetime-local" />
        </div>
      </div>

      <div class="grid gap-4 sm:grid-cols-2">
        <div class="flex flex-col gap-2">
          <Label for="excerpt_en">Excerpt (EN)</Label>
          <Textarea id="excerpt_en" v-model="form.excerpt_en" />
        </div>
        <div class="flex flex-col gap-2">
          <Label for="excerpt_id">Excerpt (ID)</Label>
          <Textarea id="excerpt_id" v-model="form.excerpt_id" />
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <Label for="content_en">Content (EN) — Markdown</Label>
        <Textarea id="content_en" v-model="form.content_en" class="min-h-64 font-mono" />
      </div>
      <div class="flex flex-col gap-2">
        <Label for="content_id">Content (ID) — Markdown</Label>
        <Textarea id="content_id" v-model="form.content_id" class="min-h-64 font-mono" />
      </div>

      <div class="flex gap-2">
        <Button type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save Post' }}</Button>
        <Button type="button" variant="ghost" @click="router.push('/posts')">Cancel</Button>
      </div>
    </form>
  </div>
</template>
