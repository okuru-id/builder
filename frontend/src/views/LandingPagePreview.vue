<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { IconArrowLeft } from '@tabler/icons-vue'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'

interface LandingTemplate {
  id: number
  name: string
  description: string
  preview: string
  sections: { type: string }[]
  html?: string
}

const route = useRoute()
const router = useRouter()

const template = ref<LandingTemplate | null>(null)
const loading = ref(true)
const applying = ref(false)
const blobUrl = ref<string>('')

const sectionLabels: Record<string, string> = {
  hero: 'Hero',
  clients: 'Client Logos',
  services: 'Services',
  projects: 'Open Source Projects',
  cta: 'CTA',
}

function sectionLabel(type: string) {
  return sectionLabels[type] ?? type
}

const sections = computed(() => template.value?.sections || [])

function makeBlobUrl(html: string) {
  const blob = new Blob([html], { type: 'text/html' })
  return URL.createObjectURL(blob)
}

async function load() {
  loading.value = true
  try {
    const id = route.params.id
    const { data } = await api.get(`/landing-templates/${id}`)
    template.value = data.data ?? null
    if (template.value?.html) blobUrl.value = makeBlobUrl(template.value.html)
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load template')
    router.push('/landing-page')
  } finally {
    loading.value = false
  }
}

async function applyTemplate() {
  if (!template.value || !confirm(`Apply "${template.value.name}"? This will replace all current sections.`)) return
  applying.value = true
  try {
    await api.post(`/landing-templates/${template.value.id}/apply`)
    toast.success('Template applied')
    router.push('/landing-page')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to apply template')
  } finally {
    applying.value = false
  }
}

onMounted(load)

onBeforeUnmount(() => {
  if (blobUrl.value) URL.revokeObjectURL(blobUrl.value)
})
</script>

<template>
  <div class="flex h-screen flex-col">
    <header class="flex h-14 flex-shrink-0 items-center justify-between gap-3 border-b px-4">
      <div class="flex min-w-0 items-center gap-3">
        <Button variant="ghost" size="icon" as-child>
          <RouterLink to="/landing-page" aria-label="Back">
            <IconArrowLeft class="size-4" />
          </RouterLink>
        </Button>
        <div class="min-w-0">
          <h1 v-if="template" class="truncate font-heading text-base font-semibold">{{ template.name }}</h1>
          <p v-if="template" class="truncate text-xs text-muted-foreground">{{ template.description }}</p>
        </div>
      </div>
      <Button :disabled="applying || !template" @click="applyTemplate">
        {{ applying ? 'Applying...' : 'Apply Template' }}
      </Button>
    </header>

    <div class="flex-1 min-h-0 bg-muted">
      <div v-if="loading" class="flex size-full items-center justify-center text-sm text-muted-foreground">
        Loading preview...
      </div>
      <iframe
        v-else-if="blobUrl"
        :src="blobUrl"
        :title="template?.name ?? 'preview'"
        class="size-full border-0 bg-white"
        sandbox="allow-scripts"
      />
      <img
        v-else-if="template?.preview"
        :src="template.preview"
        :alt="template?.name ?? 'preview'"
        class="size-full object-cover bg-white"
      >
      <div v-else class="flex size-full flex-col items-center justify-center gap-3 text-center text-muted-foreground">
        <span class="text-sm">{{ sections.length }} sections</span>
        <div class="flex max-w-md flex-wrap justify-center gap-2">
          <span v-for="section in sections" :key="section.type" class="rounded-full bg-background px-3 py-1 text-xs">
            {{ sectionLabel(section.type) }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
