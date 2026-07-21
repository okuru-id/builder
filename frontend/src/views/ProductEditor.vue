<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import {
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue,
} from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { IconArrowLeft } from '@tabler/icons-vue'

const route = useRoute()
const router = useRouter()
const id = computed(() => (route.params.id as string) || null)
const isEdit = computed(() => !!id.value)
const loading = ref(false)
const saving = ref(false)

const types = ['ebook', 'template', 'source_code', 'donation']

const form = ref({
  title: '', slug: '', description: '', price: 0, type: 'ebook', status: 'draft',
})

async function load() {
  if (!isEdit.value) return
  loading.value = true
  try {
    const { data } = await api.get(`/products/${id.value}`)
    const p = data.data
    form.value = {
      title: p.title ?? '',
      slug: p.slug ?? '',
      description: p.description ?? '',
      price: p.price ?? 0,
      type: p.type ?? 'ebook',
      status: p.status ?? 'draft',
    }
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load')
    router.push('/products')
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!form.value.title || !form.value.slug) {
    toast.warning('Title and Slug are required')
    return
  }
  saving.value = true
  const payload: any = {
    title: form.value.title,
    slug: form.value.slug,
    description: form.value.description || null,
    price: Number(form.value.price) || 0,
    type: form.value.type,
    status: form.value.status,
  }
  try {
    if (isEdit.value) await api.put(`/products/${id.value}`, payload)
    else await api.post('/products', payload)
    toast.success('Saved')
    router.push('/products')
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
    <Button variant="ghost" size="sm" class="mb-4 -ml-2" @click="router.push('/products')">
      <IconArrowLeft class="size-4" />
      Back
    </Button>

    <h1 class="font-heading mb-6 text-2xl font-bold">{{ isEdit ? 'Edit Product' : 'New Product' }}</h1>

    <div v-if="loading" class="max-w-2xl space-y-4">
      <Skeleton v-for="i in 4" :key="i" class="h-10 w-full" />
    </div>

    <form v-else class="grid max-w-xl grid-cols-[max-content_1fr] items-center gap-x-3 gap-y-2" @submit.prevent="save">
      <Label class="text-right whitespace-nowrap">Title *</Label>
      <Input v-model="form.title" />
      <Label class="text-right whitespace-nowrap">Slug *</Label>
      <Input v-model="form.slug" />
      <Label class="self-start pt-2 text-right whitespace-nowrap">Description</Label>
      <Textarea v-model="form.description" />
      <Label class="text-right whitespace-nowrap">Price (IDR)</Label>
      <Input v-model.number="form.price" type="number" min="0" />
      <Label class="text-right whitespace-nowrap">Type</Label>
      <Select v-model="form.type">
        <SelectTrigger class="w-full"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="t in types" :key="t" :value="t">{{ t }}</SelectItem>
        </SelectContent>
      </Select>
      <Label class="text-right whitespace-nowrap">Status</Label>
      <Select v-model="form.status">
        <SelectTrigger class="w-full"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem value="draft">draft</SelectItem>
          <SelectItem value="published">published</SelectItem>
        </SelectContent>
      </Select>
      <div></div>
      <div class="flex gap-2">
        <Button type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</Button>
        <Button type="button" variant="ghost" @click="router.push('/products')">Cancel</Button>
      </div>
    </form>
  </div>
</template>
