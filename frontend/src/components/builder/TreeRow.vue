<script setup lang="ts">
// Recursive outline row. Self-references by component name for nested children.
import { inject } from 'vue'
import type { Node } from '@/types/page-builder'
import { BUILDER_KEY } from '@/components/builder/injection'

defineOptions({ name: 'TreeRow' })
const props = withDefaults(
  defineProps<{ node: Node; depth?: number }>(),
  { depth: 0 },
)

const store = inject(BUILDER_KEY, null)!
const selected = () => store.selectedId.value === props.node.id
</script>

<template>
  <div>
    <button
      class="flex w-full items-center gap-1.5 py-1 pr-2 text-left text-xs"
      :class="selected() ? 'bg-blue-50 text-blue-700' : 'hover:bg-neutral-100'"
      :style="{ paddingLeft: `${depth * 12 + 8}px` }"
      @click.stop="store.select(node.id)"
    >
      <span class="size-1.5 shrink-0 rounded-full"
        :class="selected() ? 'bg-blue-500' : 'bg-neutral-300'"
      />
      <span class="truncate">{{ node.name || node.type }}</span>
    </button>
    <TreeRow
      v-for="child in node.children"
      :key="child.id"
      :node="child"
      :depth="depth + 1"
    />
  </div>
</template>
