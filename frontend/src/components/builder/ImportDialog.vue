<script setup lang="ts">
// Import HTML/Vue → new landing page.
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { IconUpload, IconX, IconAlertTriangle } from '@tabler/icons-vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { parseHTML, parseVueTemplate } from '@/services/parser'
import api from '@/lib/api'
import type { Node } from '@/types/page-builder'

const emit = defineEmits<{ close: [] }>()
const router = useRouter()

const tab = ref<'paste' | 'file'>('paste')
const rawHtml = ref('')
const pageName = ref('')
const parsing = ref(false)
const created = ref(false)

// Drag-and-drop file state
const dragOver = ref(false)
const fileName = ref('')

let parsedTree: Node | null = null
const parseError = ref('')
const nodeCount = ref(0)
const fileInput = ref<HTMLInputElement | null>(null)

function doParse() {
  const src = rawHtml.value.trim()
  if (!src) return
  parsing.value = true
  parseError.value = ''
  nodeCount.value = 0

  try {
    const isVue = fileName.value.endsWith('.vue') || src.includes('<template>')
    parsedTree = isVue ? parseVueTemplate(src) : parseHTML(src)
    if (!parsedTree || (!parsedTree.children?.length && !parsedTree.props?.text)) {
      parseError.value = 'Tidak ditemukan elemen yang bisa di-parse. Periksa HTML.'
    } else {
      countNodes(parsedTree)
      if (!pageName.value) pageName.value = fileName.value.replace(/\.\w+$/, '') || 'Imported'
    }
  } catch (e: any) {
    parseError.value = e?.message || 'Gagal parse HTML'
    parsedTree = null
  } finally {
    parsing.value = false
  }
}

function countNodes(n: Node): number {
  let c = 1
  for (const child of n.children) c += countNodes(child)
  nodeCount.value = c
  return c
}

// File upload
function onFileDrop(e: DragEvent) {
  dragOver.value = false
  const file = e.dataTransfer?.files?.[0]
  if (file) readFile(file)
}

function onFileInput(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (file) readFile(file)
}

function readFile(file: File) {
  if (!file.name.match(/\.(html?|vue|htm)$/i)) {
    toast.error('Hanya file .html, .htm, atau .vue')
    return
  }
  fileName.value = file.name
  const reader = new FileReader()
  reader.onload = () => {
    rawHtml.value = reader.result as string
    pageName.value = file.name.replace(/\.\w+$/, '')
    doParse()
  }
  reader.readAsText(file)
}

async function createPage() {
  if (!parsedTree || !pageName.value.trim()) return
  creating.value = true
  try {
    const res = await api.post<{ data: any }>('/landing-pages', {
      name: pageName.value.trim(),
      tree: { root: parsedTree },
    })
    created.value = true
    toast.success('Halaman dibuat dari import')
    router.push({ name: 'builder', params: { id: res.data.data.id } })
  } catch (e: any) {
    toast.error(e?.response?.data?.error || 'Gagal membuat halaman')
  } finally {
    creating.value = false
  }
}

const creating = ref(false)
// Handle CreatePage emit
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" @click.self="emit('close')">
    <div class="flex max-h-[85vh] w-[600px] flex-col rounded-xl border border-neutral-200 bg-white shadow-2xl">
      <div class="flex items-center justify-between border-b border-neutral-200 px-5 py-3">
        <h2 class="text-sm font-semibold">Import HTML / Vue</h2>
        <Button variant="ghost" size="icon" @click="emit('close')"><IconX class="size-4" /></Button>
      </div>

      <Tabs :default-value="tab" class="flex flex-1 flex-col overflow-hidden" @update:model-value="(v) => tab = v as 'paste' | 'file'">
        <div class="border-b border-neutral-200 px-5">
          <TabsList class="h-10">
            <TabsTrigger value="paste" class="text-xs">Paste HTML</TabsTrigger>
            <TabsTrigger value="file" class="text-xs">Upload File</TabsTrigger>
          </TabsList>
        </div>

        <div class="flex flex-1 flex-col overflow-auto p-5 pt-3 space-y-3">
          <TabsContent value="paste" class="m-0 flex flex-col gap-3">
            <Textarea
              v-model="rawHtml"
              rows="10"
              class="min-h-[200px] resize-y font-mono text-xs"
              placeholder="Tempel HTML atau Vue template di sini..."
            />
          </TabsContent>

          <TabsContent value="file" class="m-0 flex flex-col gap-3">
            <div
              class="flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-neutral-300 p-8 transition-colors"
              :class="dragOver ? 'border-blue-500 bg-blue-50' : 'hover:border-neutral-400'"
              @dragover.prevent="dragOver = true"
              @dragleave="dragOver = false"
              @drop.prevent="onFileDrop"
              @click="fileInput?.click()"
            >
              <IconUpload class="mb-2 size-8 text-neutral-400" />
              <p class="text-sm text-neutral-600">Seret file .html atau .vue ke sini</p>
              <p class="mt-1 text-xs text-neutral-400">atau klik untuk pilih file</p>
              <input ref="fileInput" type="file" accept=".html,.htm,.vue" class="hidden" @change="onFileInput" />
            </div>
          </TabsContent>

          <!-- Parse controls -->
          <div v-if="tab === 'paste'" class="flex justify-end pt-2">
            <Button size="sm" :disabled="!rawHtml.trim() || parsing" @click="doParse">
              {{ parsing ? 'Mem-parse…' : 'Parse' }}
            </Button>
          </div>

          <!-- Parse result -->
          <div v-if="parseError" class="flex items-start gap-2 rounded-md border border-red-200 bg-red-50 px-3 py-2 text-xs text-red-700">
            <IconAlertTriangle class="mt-0.5 size-3.5 shrink-0" />
            {{ parseError }}
          </div>

          <div v-if="parsedTree && !parseError" class="space-y-3">
            <div class="rounded-md border border-green-200 bg-green-50 px-3 py-2 text-xs text-green-700">
              Parse OK — {{ nodeCount }} node terdeteksi dari
              <span class="font-mono">{{ fileName || 'paste' }}</span>
            </div>

            <div class="space-y-1.5">
              <Label class="text-xs">Nama halaman</Label>
              <Input v-model="pageName" class="h-8" placeholder="Nama halaman baru" />
            </div>

            <div class="flex justify-end gap-2 pt-2">
              <Button variant="outline" size="sm" @click="emit('close')">Batal</Button>
              <Button size="sm" :disabled="!pageName.trim() || creating" @click="createPage">
                {{ creating ? 'Membuat…' : 'Buat Halaman' }}
              </Button>
            </div>
          </div>
        </div>
      </Tabs>
    </div>
  </div>
</template>
