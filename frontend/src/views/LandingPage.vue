<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { IconLayoutGrid } from '@tabler/icons-vue'
import api from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

interface LandingTemplate {
  id: number
  name: string
  description: string
  preview: string
  sections: { type: string }[]
  html?: string
}

const router = useRouter()

const templates = ref<LandingTemplate[]>([])
const loading = ref(true)
const brokenPreviews = ref(new Set<number>())

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

function previewAvailable(template: LandingTemplate) {
  return Boolean(template.preview) && !brokenPreviews.value.has(template.id)
}

function markPreviewBroken(id: number) {
  brokenPreviews.value = new Set([...brokenPreviews.value, id])
}

function openPreview(template: LandingTemplate) {
  router.push({ name: 'landing-page-preview', params: { id: template.id } })
}

async function loadTemplates() {
  loading.value = true
  try {
    const { data } = await api.get('/landing-templates')
    templates.value = data.data ?? []
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load templates')
  } finally {
    loading.value = false
  }
}

onMounted(loadTemplates)
</script>

<template>
  <div class="flex flex-col gap-6 p-8 min-w-0">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="font-heading text-2xl font-bold">Landing Page Templates</h1>
        <p class="text-sm text-muted-foreground">Choose a template for your landing page.</p>
      </div>
      <Button as-child variant="outline">
        <RouterLink to="/landing-page/content">Edit Content</RouterLink>
      </Button>
    </div>

    <div v-if="loading" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <Skeleton v-for="i in 3" :key="i" class="h-72 rounded-xl" />
    </div>

    <div v-else-if="templates.length === 0" class="py-12 text-center text-muted-foreground">
      <IconLayoutGrid class="mx-auto mb-3 size-10" />
      <p class="text-sm">No templates yet.</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <Card
        v-for="template in templates"
        :key="template.id"
        class="gap-0 overflow-hidden py-0 transition-shadow hover:shadow-md"
      >
        <button
          type="button"
          class="block w-full text-left focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
          :aria-label="`Preview ${template.name} template`"
          @click="openPreview(template)"
        >
          <div class="aspect-video overflow-hidden bg-muted">
            <img
              v-if="previewAvailable(template)"
              :src="template.preview"
              :alt="`${template.name} preview`"
              class="size-full object-cover"
              @error="markPreviewBroken(template.id)"
            >
            <div v-else class="flex size-full flex-col items-center justify-center gap-2 p-4 text-center text-muted-foreground">
              <IconLayoutGrid class="size-8" />
              <span class="text-xs">{{ template.sections?.length || 0 }} sections</span>
              <div class="flex flex-wrap justify-center gap-1">
                <span v-for="section in (template.sections || []).slice(0, 3)" :key="section.type" class="rounded-full bg-background px-2 py-0.5 text-[10px] text-foreground">
                  {{ sectionLabel(section.type) }}
                </span>
              </div>
            </div>
          </div>
          <div class="space-y-2 p-4">
            <div class="flex items-center justify-between gap-2">
              <h2 class="font-semibold">{{ template.name }}</h2>
              <Badge variant="secondary">{{ template.sections?.length || 0 }} sections</Badge>
            </div>
            <p class="line-clamp-2 text-sm text-muted-foreground">{{ template.description }}</p>
            <div class="flex flex-wrap gap-1 pt-1">
              <span v-for="section in (template.sections || []).slice(0, 5)" :key="section.type" class="rounded-full bg-secondary px-2 py-0.5 text-[10px] text-secondary-foreground">
                {{ sectionLabel(section.type) }}
              </span>
              <span v-if="(template.sections?.length || 0) > 5" class="text-xs text-muted-foreground">+{{ (template.sections?.length || 0) - 5 }} more</span>
            </div>
          </div>
        </button>
      </Card>
    </div>
  </div>
</template>
