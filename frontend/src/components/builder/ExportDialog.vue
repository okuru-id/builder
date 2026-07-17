<script setup lang="ts">
// Export/download dialog: preview generated HTML, copy, download, view safelist.
// HTML preview uses vue-prism-editor + prismjs (markup grammar) for syntax highlighting.
import { inject, ref, computed } from 'vue'
import {
  IconX,
  IconDownload,
  IconCopy,
  IconCheck,
  IconList,
  IconCode,
  IconFileCode,
  IconBrandHtml5,
} from '@tabler/icons-vue'
import { PrismEditor } from 'vue-prism-editor'
import 'vue-prism-editor/dist/prismeditor.min.css'
import 'prismjs/themes/prism-tomorrow.min.css'
import Prism from 'prismjs/components/prism-core'
import 'prismjs/components/prism-markup'
import 'prismjs/components/prism-clike'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-css'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { BUILDER_KEY } from '@/components/builder/injection'
import { generateHTML, download as dl, collectClasses, htmlEscape } from '@/composables/useCodegen'

const store = inject(BUILDER_KEY, null)!
const emit = defineEmits<{ close: [] }>()

const copied = ref<'html' | 'safelist' | null>(null)
const tab = ref<'html' | 'safelist'>('html')

const html = computed(() => {
  if (!store.tree.value) return ''
  return generateHTML(store.tree.value, store.page.value?.name ?? 'Landing')
})

const safelist = computed(() => {
  if (!store.tree.value) return ''
  return collectClasses(store.tree.value).join('\n')
})

const htmlStats = computed(() => {
  const bytes = new Blob([html.value]).size
  const lines = html.value ? html.value.split('\n').length : 0
  return { bytes, lines, kb: (bytes / 1024).toFixed(1) }
})

const safelistStats = computed(() => (safelist.value ? safelist.value.split('\n').length : 0))

// Prism highlighters. PrismEditor.highlight receives code, returns HTML string.
function highlightHTML(code: string): string {
  if (!code) return ''
  return Prism.highlight(code, Prism.languages.markup, 'markup')
}
function highlightPlain(code: string): string {
  // Safelist = bare class names; no grammar. Escape only.
  return htmlEscape(code)
}

function copy(kind: 'html' | 'safelist') {
  const text = kind === 'html' ? html.value : safelist.value
  navigator.clipboard.writeText(text).then(() => {
    copied.value = kind
    setTimeout(() => (copied.value = null), 1800)
  })
}

function doDownload() {
  const name = store.page.value?.slug ?? 'landing'
  dl(`${name}.html`, html.value)
}

function doSafelistDownload() {
  const name = store.page.value?.slug ?? 'landing'
  dl(`${name}-safelist.txt`, safelist.value, 'text/plain')
}
</script>

<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
    appear
  >
    <div
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
      @click.self="emit('close')"
    >
      <Transition
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="opacity-0 scale-95 translate-y-2"
        enter-to-class="opacity-100 scale-100 translate-y-0"
        leave-active-class="transition duration-150 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
        appear
      >
        <div
          class="flex max-h-[92vh] w-full max-w-[1400px] flex-col overflow-hidden rounded-xl border border-border bg-background shadow-2xl"
        >
          <!-- header -->
          <div class="flex items-start justify-between gap-3 border-b border-border bg-muted/30 px-5 py-4">
            <div class="flex items-start gap-3">
              <div class="mt-0.5 flex size-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
                <IconFileCode class="size-5" />
              </div>
              <div>
                <h2 class="text-sm font-semibold leading-tight">Export Page</h2>
                <p class="mt-0.5 text-xs text-muted-foreground">
                  Self-contained HTML with Tailwind CDN, or the safelist for your config.
                </p>
              </div>
            </div>
            <Button variant="ghost" size="icon" class="-mr-1.5 -mt-1 size-7" @click="emit('close')">
              <IconX class="size-4" />
            </Button>
          </div>

          <!-- tabs + body -->
          <Tabs v-model="tab" default-value="html" class="flex min-h-0 flex-1 flex-col">
            <div class="flex items-center justify-between border-b border-border px-3 py-2">
              <TabsList class="h-9">
                <TabsTrigger value="html" class="gap-1.5 text-xs">
                  <IconBrandHtml5 class="size-3.5" /> HTML
                </TabsTrigger>
                <TabsTrigger value="safelist" class="gap-1.5 text-xs">
                  <IconList class="size-3.5" /> Safelist
                </TabsTrigger>
              </TabsList>

              <!-- inline stats -->
              <div v-if="tab === 'html'" class="flex items-center gap-1.5 text-[11px] text-muted-foreground">
                <span class="rounded-md bg-muted px-1.5 py-0.5 font-mono">{{ htmlStats.lines }} lines</span>
                <span class="rounded-md bg-muted px-1.5 py-0.5 font-mono">{{ htmlStats.kb }} KB</span>
              </div>
              <div v-else class="flex items-center gap-1.5 text-[11px] text-muted-foreground">
                <span class="rounded-md bg-muted px-1.5 py-0.5 font-mono">{{ safelistStats }} classes</span>
              </div>
            </div>

            <!-- HTML content -->
            <TabsContent value="html" class="mt-0 flex min-h-0 flex-1 flex-col overflow-hidden">
              <PrismEditor
                class="export-editor"
                :model-value="html"
                :highlight="highlightHTML"
                :readonly="true"
                :line-numbers="true"
                :tab-size="2"
              />
            </TabsContent>

            <!-- Safelist content -->
            <TabsContent value="safelist" class="mt-0 flex min-h-0 flex-1 flex-col overflow-hidden">
              <div class="border-b border-border bg-muted/30 px-4 py-2 text-[11px] text-muted-foreground">
                Every unique Tailwind class on this page. Paste into
                <code class="rounded bg-background px-1 py-0.5 font-mono text-[10px] text-foreground">tailwind.config.safelist</code>
                to prevent purging.
              </div>
              <PrismEditor
                class="export-editor"
                :model-value="safelist"
                :highlight="highlightPlain"
                :readonly="true"
                :line-numbers="true"
                :tab-size="2"
              />
            </TabsContent>
          </Tabs>

          <!-- footer -->
          <div class="flex items-center justify-between gap-2 border-t border-border bg-muted/30 px-4 py-3">
            <div class="flex items-center gap-2">
              <Button
                v-if="tab === 'html'"
                variant="outline"
                size="sm"
                class="h-8"
                @click="copy('html')"
              >
                <component :is="copied === 'html' ? IconCheck : IconCopy" class="size-3.5" :class="copied === 'html' ? 'text-green-600' : ''" />
                {{ copied === 'html' ? 'Copied' : 'Copy' }}
              </Button>
              <Button v-else variant="outline" size="sm" class="h-8" @click="copy('safelist')">
                <component :is="copied === 'safelist' ? IconCheck : IconCopy" class="size-3.5" :class="copied === 'safelist' ? 'text-green-600' : ''" />
                {{ copied === 'safelist' ? 'Copied' : 'Copy' }}
              </Button>
            </div>

            <div class="flex items-center gap-2">
              <Button
                v-if="tab === 'safelist'"
                variant="ghost"
                size="sm"
                class="h-8 text-muted-foreground"
                @click="tab = 'html'"
              >
                <IconCode class="size-3.5" /> View HTML
              </Button>
              <Button
                v-if="tab === 'html'"
                variant="ghost"
                size="sm"
                class="h-8 text-muted-foreground"
                @click="tab = 'safelist'"
              >
                <IconList class="size-3.5" /> View Safelist
              </Button>

              <Button
                v-if="tab === 'safelist'"
                variant="outline"
                size="sm"
                class="h-8"
                @click="doSafelistDownload"
              >
                <IconDownload class="size-3.5" /> .txt
              </Button>
              <Button v-if="tab === 'html'" size="sm" class="h-8" @click="doDownload">
                <IconDownload class="size-3.5" /> Download .html
              </Button>
            </div>
          </div>
        </div>
      </Transition>
    </div>
  </Transition>
</template>

<style scoped>
/* vue-prism-editor root = .prism-editor-wrapper (class lands on it). Give it a
   concrete viewport height so the inner overflow:auto scrolls reliably;
   flex chains through reka TabsContent are unreliable for % heights. */
.export-editor {
  height: 62vh;
  background: #1e1e2e;
  color: #e4e4e7;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 12px;
  line-height: 1.6;
}

/* Editor inner padding (overrides component default of 0). */
.export-editor :deep(.prism-editor__container) {
  padding: 12px 14px;
}

/* No line wrapping — allow horizontal scroll for long HTML lines. */
.export-editor :deep(.prism-editor__editor),
.export-editor :deep(.prism-editor__textarea) {
  white-space: pre !important;
}

/* Line numbers gutter tint + alignment. */
.export-editor :deep(.prism-editor__line-numbers) {
  padding-top: 12px;
  margin-right: 14px;
  color: #6b6b7b;
  min-width: 2.5rem;
}

/* Selection color inside the editor. */
.export-editor :deep(.prism-editor__textarea::selection) {
  background: rgba(99, 102, 241, 0.35);
}

/* Scrollbar: subtle, matches dark theme. */
.export-editor :deep(.prism-editor-wrapper::-webkit-scrollbar) {
  width: 10px;
  height: 10px;
}
.export-editor :deep(.prism-editor-wrapper::-webkit-scrollbar-thumb) {
  background: #3f3f52;
  border-radius: 6px;
}
.export-editor :deep(.prism-editor-wrapper::-webkit-scrollbar-thumb:hover) {
  background: #525266;
}
.export-editor :deep(.prism-editor-wrapper::-webkit-scrollbar-track) {
  background: transparent;
}
</style>
