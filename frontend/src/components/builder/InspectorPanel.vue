<script setup lang="ts">
// Right inspector: edits the selected node's props + classes. Delete/duplicate actions.
import { computed, inject } from 'vue'
import { IconCopy, IconTrash, IconArrowUp, IconArrowDown, IconUnlink } from '@tabler/icons-vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { BUILDER_KEY } from '@/components/builder/injection'
import LayoutSection from './inspector/LayoutSection.vue'
import TypographySection from './inspector/TypographySection.vue'
import SpacingSection from './inspector/SpacingSection.vue'
import BackgroundSection from './inspector/BackgroundSection.vue'
import BorderSection from './inspector/BorderSection.vue'
import SizeSection from './inspector/SizeSection.vue'

const store = inject(BUILDER_KEY, null)!

const node = computed(() => store.selectedNode.value)

function setProp(key: string, value: unknown) {
  if (!node.value) return
  store.patchNode(node.value.id, { props: { ...node.value.props, [key]: value } })
}

function del() {
  if (node.value) store.removeNode(node.value.id)
}
function dup() {
  if (node.value) store.duplicateNode(node.value.id)
}
function moveUp() {
  if (node.value) store.moveSiblingNode(node.value.id, -1)
}
function moveDown() {
  if (node.value) store.moveSiblingNode(node.value.id, 1)
}

function addClass(e: Event) {
  const input = e.target as HTMLInputElement
  const cls = input.value.trim()
  if (!cls || !node.value) return
  // Support space-separated classes
  const newClasses = cls.split(/\s+/).filter((c) => c && !node.value!.classes.includes(c))
  if (newClasses.length) {
    store.patchNode(node.value.id, { classes: [...node.value.classes, ...newClasses] })
  }
  input.value = ''
}

function removeClass(idx: number) {
  if (!node.value) return
  const classes = [...node.value.classes]
  classes.splice(idx, 1)
  store.patchNode(node.value.id, { classes })
}
</script>

<template>
  <aside class="flex w-72 flex-col border-l border-neutral-200 bg-white">
    <div class="border-b border-neutral-200 px-4 py-3">
      <h2 class="text-sm font-semibold">Inspector</h2>
    </div>

    <div v-if="!node" class="p-4 text-sm text-neutral-400">
      Pilih sebuah node di canvas untuk mengubah propertinya.
    </div>

    <div v-else class="flex-1 space-y-4 overflow-auto p-4">
      <div class="space-y-1.5">
        <Label class="text-xs">Nama</Label>
        <Input :model-value="node.name" class="h-8" @update:model-value="(v) => setProp('name', v)" />
        <!-- ponytail: name lives on the node for tree-panel display; we store it on props
             as a UI convenience but the canonical name field is node.name. -->
      </div>

      <div class="space-y-1.5">
        <Label class="text-xs">Tipe</Label>
        <div class="text-sm text-neutral-600">{{ node.type }}</div>
      </div>

      <div
        v-if="node.componentId"
        class="flex items-center justify-between rounded-md border border-blue-200 bg-blue-50 px-3 py-2"
      >
        <div class="text-xs text-blue-700">
          Instance dari komponen #{{ node.componentId }}
        </div>
        <Button variant="outline" size="sm" @click="store.breakInstance(node.id)">
          <IconUnlink class="size-3.5" /> Putus
        </Button>
      </div>

      <!-- Text-like: editable text + heading level -->
      <template v-if="node.type === 'text' || node.type === 'heading' || node.type === 'button' || node.type === 'link'">
        <div class="space-y-1.5">
          <Label class="text-xs">Teks</Label>
          <Textarea
            :model-value="String(node.props.text ?? '')"
            rows="3"
            @update:model-value="(v) => setProp('text', v)"
          />
        </div>
      </template>

      <div v-if="node.type === 'heading'" class="space-y-1.5">
        <Label class="text-xs">Level</Label>
        <Select :model-value="String(node.props.level ?? 2)" @update:model-value="(v) => setProp('level', Number(v))">
          <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="l in [1, 2, 3, 4, 5, 6]" :key="l" :value="String(l)">H{{ l }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div v-if="node.type === 'image'" class="space-y-3">
        <div class="space-y-1.5">
          <Label class="text-xs">URL gambar</Label>
          <Input :model-value="String(node.props.src ?? '')" class="h-8" @update:model-value="(v) => setProp('src', v)" />
        </div>
        <div class="space-y-1.5">
          <Label class="text-xs">Alt</Label>
          <Input :model-value="String(node.props.alt ?? '')" class="h-8" @update:model-value="(v) => setProp('alt', v)" />
        </div>
      </div>

      <div v-if="node.type === 'link'" class="space-y-1.5">
        <Label class="text-xs">Href</Label>
        <Input :model-value="String(node.props.href ?? '#')" class="h-8" @update:model-value="(v) => setProp('href', v)" />
      </div>

      <!-- Style sections: only for non-instance nodes -->
      <template v-if="!node.componentId">
        <LayoutSection />
        <TypographySection />
        <SpacingSection />
        <BackgroundSection />
        <BorderSection />
        <SizeSection />

        <!-- Custom classes: arbitrary Tailwind class input -->
        <div class="space-y-2 border-b border-neutral-100 pb-3">
          <h3 class="text-xs font-medium uppercase tracking-wider text-neutral-500">Custom Class</h3>
          <div class="flex flex-wrap gap-1">
            <span
              v-for="(cls, idx) in node.classes"
              :key="idx"
              class="group inline-flex items-center gap-0.5 rounded bg-neutral-100 px-1.5 py-0.5 text-[11px] font-mono text-neutral-700"
            >
              {{ cls }}
              <button
                class="ml-0.5 rounded-sm text-neutral-400 opacity-0 transition-opacity hover:text-red-500 group-hover:opacity-100"
                @click="removeClass(idx)"
                title="Hapus class"
              >&times;</button>
            </span>
          </div>
          <Input
            :model-value="''"
            class="h-8 font-mono text-xs"
            placeholder="Tambah class, tekan Enter"
            @keydown.enter.prevent="addClass"
          />
        </div>
      </template>
      <div v-else class="rounded-md border border-neutral-200 bg-neutral-50 px-3 py-2 text-xs text-neutral-500">
        Edit komponen master untuk mengubah instance. Klik &quot;Putus link&quot; di atas untuk jadikan copy independen.
      </div>

      <div class="flex flex-wrap gap-2 pt-2" v-if="node.id !== 'root'">
        <Button variant="outline" size="sm" @click="moveUp">
          <IconArrowUp class="size-4" /> Naik
        </Button>
        <Button variant="outline" size="sm" @click="moveDown">
          <IconArrowDown class="size-4" /> Turun
        </Button>
        <Button variant="outline" size="sm" @click="dup">
          <IconCopy class="size-4" /> Duplikat
        </Button>
        <Button variant="destructive" size="sm" @click="del">
          <IconTrash class="size-4" /> Hapus
        </Button>
      </div>
    </div>
  </aside>
</template>
