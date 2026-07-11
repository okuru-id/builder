<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { toast } from 'vue-sonner'
import { IconLayoutGrid } from '@tabler/icons-vue'
import api from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Sheet, SheetClose, SheetContent, SheetDescription, SheetFooter, SheetHeader, SheetTitle } from '@/components/ui/sheet'
import { Skeleton } from '@/components/ui/skeleton'

interface LandingTemplate {
  id: number
  name: string
  description: string
  preview: string
  sections: { type: string }[]
  html?: string
}

const templates = ref<LandingTemplate[]>([])
const loading = ref(true)
const selectedTemplate = ref<LandingTemplate | null>(null)
const sheetOpen = ref(false)
const applyingTemplate = ref<number | null>(null)
const brokenPreviews = ref(new Set<number>())
const blobUrls = ref(new Map<number, string>())

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

function templatePreviewUrl(template: LandingTemplate) {
  if (!template.html) return ''
  if (!blobUrls.value.has(template.id)) {
    const blob = new Blob([template.html], { type: 'text/html' })
    blobUrls.value.set(template.id, URL.createObjectURL(blob))
  }
  return blobUrls.value.get(template.id)!
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

function openPreview(template: LandingTemplate) {
  selectedTemplate.value = template
  sheetOpen.value = true
}

function updateSheet(open: boolean) {
  if (!open && applyingTemplate.value !== null) return
  sheetOpen.value = open
  if (!open) selectedTemplate.value = null
}

async function applySelectedTemplate() {
  const template = selectedTemplate.value
  if (!template || !confirm(`Apply "${template.name}"? This will replace all current sections.`)) return

  applyingTemplate.value = template.id
  try {
    await api.post(`/landing-templates/${template.id}/apply`)
    toast.success('Template applied')
    sheetOpen.value = false
    selectedTemplate.value = null
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to apply template')
  } finally {
    applyingTemplate.value = null
  }
}

onMounted(loadTemplates)

onBeforeUnmount(() => {
  for (const url of blobUrls.value.values()) {
    URL.revokeObjectURL(url)
  }
  blobUrls.value.clear()
})
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
      <Card v-for="template in templates" :key="template.id" class="overflow-hidden transition-shadow hover:shadow-md">
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
          <CardHeader class="pb-2">
            <div class="flex items-center justify-between gap-2">
              <h2 class="font-semibold">{{ template.name }}</h2>
              <Badge variant="secondary">{{ template.sections?.length || 0 }} sections</Badge>
            </div>
            <p class="line-clamp-2 text-sm text-muted-foreground">{{ template.description }}</p>
          </CardHeader>
          <CardContent class="flex flex-wrap gap-1 pt-0">
            <span v-for="section in (template.sections || []).slice(0, 5)" :key="section.type" class="rounded-full bg-secondary px-2 py-0.5 text-[10px] text-secondary-foreground">
              {{ sectionLabel(section.type) }}
            </span>
            <span v-if="(template.sections?.length || 0) > 5" class="text-xs text-muted-foreground">+{{ (template.sections?.length || 0) - 5 }} more</span>
          </CardContent>
        </button>
      </Card>
    </div>

    <Sheet :open="sheetOpen" @update:open="updateSheet">
      <SheetContent v-if="selectedTemplate" class="w-screen overflow-y-auto sm:max-w-full">
        <SheetHeader>
          <SheetTitle>{{ selectedTemplate.name }}</SheetTitle>
          <SheetDescription>{{ selectedTemplate.description }}</SheetDescription>
        </SheetHeader>

        <div class="flex flex-1 flex-col gap-5 px-4">
          <div class="min-h-0">
            <iframe
              v-if="selectedTemplate.html"
              :src="templatePreviewUrl(selectedTemplate)"
              :title="`${selectedTemplate.name} preview`"
              class="h-[75vh] w-full rounded-lg border bg-white"
              sandbox="allow-scripts"
            />
            <img
              v-else-if="previewAvailable(selectedTemplate)"
              :src="selectedTemplate.preview"
              :alt="`${selectedTemplate.name} preview`"
              class="h-[75vh] w-full rounded-lg border object-cover bg-white"
              @error="markPreviewBroken(selectedTemplate.id)"
            >
            <div v-else class="flex h-[75vh] w-full flex-col items-center justify-center gap-3 rounded-lg border bg-muted text-center text-muted-foreground">
              <IconLayoutGrid class="size-10" />
              <span class="text-sm">{{ selectedTemplate.sections?.length || 0 }} sections</span>
            </div>
          </div>

          <div>
            <h3 class="mb-2 text-sm font-medium">Included sections</h3>
            <div class="flex flex-wrap gap-2">
              <Badge v-for="section in (selectedTemplate.sections || [])" :key="section.type" variant="secondary">
                {{ sectionLabel(section.type) }}
              </Badge>
              <span v-if="!selectedTemplate.sections?.length" class="text-sm text-muted-foreground">No sections listed.</span>
            </div>
          </div>
        </div>

        <SheetFooter>
          <Button :disabled="applyingTemplate !== null" @click="applySelectedTemplate">
            {{ applyingTemplate === selectedTemplate.id ? 'Applying...' : 'Apply Template' }}
          </Button>
          <SheetClose as-child>
            <Button variant="outline" :disabled="applyingTemplate !== null">Cancel</Button>
          </SheetClose>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  </div>
</template>
