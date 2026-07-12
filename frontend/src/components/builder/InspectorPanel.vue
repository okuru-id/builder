<script setup lang="ts">
// Right inspector: edits the selected node's props + classes. Delete/duplicate actions.
import { computed, inject, watch, ref } from 'vue'
import { IconCopy, IconTrash, IconArrowUp, IconArrowDown } from '@tabler/icons-vue'
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

const store = inject(BUILDER_KEY, null)!

const node = computed(() => store.selectedNode.value)

// Local buffer for classes textarea; sync when selection changes.
const classesText = ref('')
watch(
  () => node.value?.id,
  () => {
    classesText.value = node.value?.classes.join('\n') ?? ''
  },
  { immediate: true },
)

function commitClasses() {
  if (!node.value) return
  const classes = classesText.value
    .split(/\s+/)
    .map((s) => s.trim())
    .filter(Boolean)
  store.patchNode(node.value.id, { classes })
}

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

      <div class="space-y-1.5">
        <Label class="text-xs">Classes (satu per baris atau dipisah spasi)</Label>
        <Textarea v-model="classesText" rows="6" @blur="commitClasses" />
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
