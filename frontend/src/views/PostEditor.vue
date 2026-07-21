<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { CalendarDate } from '@internationalized/date'
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
import { Calendar } from '@/components/ui/calendar'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { IconArrowLeft, IconCalendar } from '@tabler/icons-vue'

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

const publishedDate = computed<CalendarDate | undefined>({
  get: () => {
    const s = form.value.published_at
    if (!s) return undefined
    const [date] = s.split('T')
    const [y, m, d] = date.split('-').map(Number)
    if (!y || !m || !d) return undefined
    return new CalendarDate(y, m, d)
  },
  set: (v) => {
    if (!v) { form.value.published_at = ''; return }
    const time = form.value.published_at?.split('T')[1] ?? '09:00'
    form.value.published_at = `${v.toString()}T${time}`
  },
})

const publishedTime = computed<string>({
  get: () => form.value.published_at?.split('T')[1] ?? '',
  set: (v) => {
    const [date] = (form.value.published_at || '').split('T')
    if (!date) return
    form.value.published_at = `${date}T${v}`
  },
})

const publishedLabel = computed(() => {
  const s = form.value.published_at
  if (!s) return 'Pick date'
  const d = new Date(s)
  if (isNaN(d.getTime())) return 'Pick date'
  return d.toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })
})

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
    <form v-else class="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_320px]" @submit.prevent="save">
      <!-- main column -->
      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-2">
            <Label for="title_en">Title (EN) *</Label>
            <Input id="title_en" v-model="form.title_en" />
          </div>
          <div class="space-y-2">
            <Label for="title_id">Title (ID)</Label>
            <Input id="title_id" v-model="form.title_id" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-2">
            <Label for="excerpt_en">Excerpt (EN)</Label>
            <Textarea id="excerpt_en" v-model="form.excerpt_en" />
          </div>
          <div class="space-y-2">
            <Label for="excerpt_id">Excerpt (ID)</Label>
            <Textarea id="excerpt_id" v-model="form.excerpt_id" />
          </div>
        </div>
        <div class="space-y-2">
          <Label for="content_en">Content (EN) — Markdown</Label>
          <Textarea id="content_en" v-model="form.content_en" class="min-h-72 font-mono" />
        </div>
        <div class="space-y-2">
          <Label for="content_id">Content (ID) — Markdown</Label>
          <Textarea id="content_id" v-model="form.content_id" class="min-h-72 font-mono" />
        </div>
      </div>

      <!-- sidebar -->
      <aside class="space-y-4 lg:sticky lg:top-4 lg:self-start">
        <div class="rounded-lg border p-4 space-y-3">
          <div class="space-y-2">
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
          <div class="space-y-2">
            <Label for="published_at">Published At</Label>
            <div class="flex gap-2">
              <Popover>
                <PopoverTrigger as-child>
                  <Button variant="outline" type="button" class="flex-1 justify-start font-normal">
                    <IconCalendar class="size-4" />
                    {{ publishedLabel }}
                  </Button>
                </PopoverTrigger>
                <PopoverContent class="w-auto p-0">
                  <Calendar v-model="publishedDate" />
                </PopoverContent>
              </Popover>
              <Input v-model="publishedTime" type="time" class="w-32" />
            </div>
          </div>
          <div class="flex gap-2 pt-1">
            <Button type="submit" class="flex-1" :disabled="saving">{{ saving ? 'Saving…' : 'Save Post' }}</Button>
            <Button type="button" variant="ghost" @click="router.push('/posts')">Cancel</Button>
          </div>
        </div>

        <div class="rounded-lg border p-4 space-y-3">
          <div class="space-y-2">
            <Label for="slug">Slug *</Label>
            <Input id="slug" v-model="form.slug" placeholder="my-post" />
          </div>
          <div class="space-y-2">
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
      </aside>
    </form>
  </div>
</template>
