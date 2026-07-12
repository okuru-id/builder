<script setup lang="ts">
// Export/download dialog: preview generated HTML, copy, download, view safelist.
import { inject, ref, computed } from 'vue'
import { IconX, IconDownload, IconCopy, IconCheck, IconCode, IconList } from '@tabler/icons-vue'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { BUILDER_KEY } from '@/components/builder/injection'
import { generateHTML, download as dl, collectClasses } from '@/composables/useCodegen'

const store = inject(BUILDER_KEY, null)!
const emit = defineEmits<{ close: [] }>()

const copied = ref(false)

const html = computed(() => {
  if (!store.tree.value) return ''
  return generateHTML(store.tree.value, store.page.value?.name ?? 'Landing')
})

const safelist = computed(() => {
  if (!store.tree.value) return ''
  return collectClasses(store.tree.value).join('\n')
})

function doCopy() {
  navigator.clipboard.writeText(html.value).then(() => {
    copied.value = true
    setTimeout(() => (copied.value = false), 2000)
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
  <!-- Backdrop -->
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" @click.self="emit('close')">
    <div class="flex max-h-[85vh] w-[700px] flex-col rounded-xl border border-border bg-background shadow-2xl">
      <!-- header -->
      <div class="flex items-center justify-between border-b border-border px-5 py-3">
        <h2 class="text-sm font-semibold">Export / Download</h2>
        <Button variant="ghost" size="icon" @click="emit('close')"><IconX class="size-4" /></Button>
      </div>

      <!-- body -->
      <Tabs default-value="html" class="flex flex-1 flex-col overflow-hidden">
        <div class="border-b border-border px-5">
          <TabsList class="h-10">
            <TabsTrigger value="html" class="flex items-center gap-1.5 text-xs">
              <IconCode class="size-3.5" /> HTML
            </TabsTrigger>
            <TabsTrigger value="safelist" class="flex items-center gap-1.5 text-xs">
              <IconList class="size-3.5" /> Safelist
            </TabsTrigger>
          </TabsList>
        </div>

        <TabsContent value="html" class="flex flex-1 flex-col overflow-hidden p-5 pt-3">
          <textarea
            readonly
            class="min-h-0 flex-1 resize-none whitespace-pre rounded-md border border-border bg-muted/50 p-3 font-mono text-[11px] leading-relaxed text-foreground focus:outline-none"
            :value="html"
          />
        </TabsContent>

        <TabsContent value="safelist" class="flex flex-1 flex-col overflow-hidden p-5 pt-3">
          <p class="mb-2 text-xs text-muted-foreground">
            Unique list of every Tailwind class used on this page. Paste it into
            <code class="rounded bg-muted px-1">tailwind.config.safelist</code>
            to prevent purging.
          </p>
          <textarea
            readonly
            class="min-h-0 flex-1 resize-none whitespace-pre rounded-md border border-border bg-muted/50 p-3 font-mono text-[11px] leading-relaxed text-foreground focus:outline-none"
            :value="safelist"
          />
        </TabsContent>
      </Tabs>

      <!-- footer -->
      <div class="flex items-center justify-between border-t border-border px-5 py-3">
        <div class="flex items-center gap-2">
          <Button variant="outline" size="sm" @click="doSafelistDownload">
            <IconDownload class="size-3.5" /> Safelist
          </Button>
        </div>
        <div class="flex items-center gap-2">
          <Button variant="outline" size="sm" @click="doCopy">
            <IconCopy v-if="!copied" class="size-3.5" />
            <IconCheck v-else class="size-3.5 text-green-600" />
            {{ copied ? 'Copied' : 'Copy HTML' }}
          </Button>
          <Button variant="default" size="sm" @click="doDownload">
            <IconDownload class="size-3.5" /> Download HTML
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>
