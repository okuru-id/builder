<script setup lang="ts">
// Component palette: list reusable masters. Click to insert instance, "save" to
// create master from selected node. Lives at the bottom of the left panel.
import { inject } from 'vue'
import { IconPlus, IconBookmark, IconUnlink, IconTrash, IconPencil } from '@tabler/icons-vue'
import { BUILDER_KEY } from '@/components/builder/injection'

const store = inject(BUILDER_KEY, null)!

async function saveSelected() {
  const id = store.selectedId.value
  if (!id || id === store.tree.value.root.id) return
  const name = window.prompt('New component name?', 'Component')
  if (!name?.trim()) return
  await store.createComponentFromNode(id, name.trim())
}

function insert(id: number) {
  store.insertInstance(id)
}

function breakLink() {
  if (store.selectedId.value) store.breakInstance(store.selectedId.value)
}

async function del(id: number, name: string) {
  if (confirm(`Delete master component "${name}"? Existing instances will not resolve.`)) {
    await store.components.remove(id)
  }
}

async function rename(c: { id: number; name: string }) {
  const name = window.prompt('Rename component?', c.name)
  if (!name?.trim() || name.trim() === c.name) return
  try {
    await store.components.rename(c.id, name.trim())
  } catch (e) {
    console.error(e)
  }
}

const selectedIsInstance = () => !!store.selectedNode.value?.componentId
const canSave = () => !!store.selectedId.value && store.selectedId.value !== store.tree.value.root.id
</script>

<template>
  <div class="border-t border-neutral-200 p-2">
    <div class="mb-1.5 flex items-center justify-between">
      <span class="text-xs font-medium text-neutral-500">Components</span>
      <button
        class="rounded p-1 text-neutral-400 hover:bg-neutral-100 hover:text-neutral-700 disabled:opacity-30"
        :disabled="!canSave()"
        title="Save selected node as a component"
        @click="saveSelected"
      >
        <IconBookmark class="size-3.5" />
      </button>
    </div>

    <button
      v-if="selectedIsInstance()"
      class="mb-1.5 flex w-full items-center gap-1.5 rounded-md border border-amber-200 bg-amber-50 px-2 py-1.5 text-xs text-amber-700 hover:bg-amber-100"
      @click="breakLink"
    >
      <IconUnlink class="size-3.5" /> Detach (make a copy)
    </button>

    <div v-if="store.components.components.value.length === 0" class="text-[11px] text-neutral-400">
      No components yet.
    </div>

    <div class="space-y-1">
      <div
        v-for="c in store.components.components.value"
        :key="c.id"
        class="group flex items-center gap-1 rounded-md border border-neutral-200 px-2 py-1.5 hover:bg-neutral-50"
      >
        <button
          class="flex flex-1 items-center gap-1.5 text-left text-xs"
          :title="`Insert instance ${c.name}`"
          @click="insert(c.id)"
        >
          <IconPlus class="size-3.5 text-neutral-400" />
          <span class="truncate">{{ c.name }}</span>
        </button>
        <button
          class="rounded p-0.5 text-neutral-300 opacity-0 hover:text-neutral-700 group-hover:opacity-100"
          title="Rename"
          @click.stop="rename(c)"
        >
          <IconPencil class="size-3.5" />
        </button>
        <button
          class="rounded p-0.5 text-neutral-300 opacity-0 hover:text-red-500 group-hover:opacity-100"
          title="Delete master"
          @click.stop="del(c.id, c.name)"
        >
          <IconTrash class="size-3.5" />
        </button>
      </div>
    </div>
  </div>
</template>
