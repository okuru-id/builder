<script setup lang="ts">
// Component palette: list reusable masters. Click to insert instance, "save" to
// create master from selected node. Lives at the bottom of the left panel.
import { inject, ref, computed } from 'vue'
import { IconPlus, IconBookmark, IconUnlink, IconTrash, IconPencil } from '@tabler/icons-vue'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
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

const showDelete = ref(false)
const deletingComponent = ref<{ id: number; name: string } | null>(null)
const deleteDesc = computed(() => deletingComponent.value
  ? `Delete master component "${deletingComponent.value.name}"? Existing instances will not resolve.`
  : '')

async function del(id: number, name: string) {
  deletingComponent.value = { id, name }
  showDelete.value = true
}
async function doDelete() {
  const dc = deletingComponent.value
  if (!dc) return
  try {
    await store.components.remove(dc.id)
  } finally {
    showDelete.value = false
    deletingComponent.value = null
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
  <div class="border-t border-border p-2">
    <div class="mb-1.5 flex items-center justify-between">
      <span class="text-xs font-medium text-muted-foreground">Components</span>
      <button
        class="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground disabled:opacity-30"
        :disabled="!canSave()"
        title="Save selected node as a component"
        @click="saveSelected"
      >
        <IconBookmark class="size-3.5" />
      </button>
    </div>

    <button
      v-if="selectedIsInstance()"
      class="mb-1.5 flex w-full items-center gap-1.5 rounded-md border border-amber-500/20 bg-amber-500/10 px-2 py-1.5 text-xs text-amber-600 dark:text-amber-400 hover:bg-amber-100"
      @click="breakLink"
    >
      <IconUnlink class="size-3.5" /> Detach (make a copy)
    </button>

    <div v-if="store.components.components.value.length === 0" class="text-[11px] text-muted-foreground">
      No components yet.
    </div>

    <div class="space-y-1">
      <div
        v-for="c in store.components.components.value"
        :key="c.id"
        class="group flex items-center gap-1 rounded-md border border-border px-2 py-1.5 hover:bg-muted/50"
      >
        <button
          class="flex flex-1 items-center gap-1.5 text-left text-xs"
          :title="`Insert instance ${c.name}`"
          @click="insert(c.id)"
        >
          <IconPlus class="size-3.5 text-muted-foreground" />
          <span class="truncate">{{ c.name }}</span>
        </button>
        <button
          class="rounded p-0.5 text-neutral-300 opacity-0 hover:text-foreground group-hover:opacity-100"
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

  <ConfirmDialog v-model:open="showDelete" title="Delete component" :description="deleteDesc" @confirm="doDelete" />
</template>
