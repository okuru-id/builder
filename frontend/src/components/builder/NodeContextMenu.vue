<script setup lang="ts">
// Shared right-click menu for builder nodes (tree + canvas).
// Actions: point at the chat agent, toggle visibility.
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from '@/components/ui/context-menu'
import {
  IconRobotFace,
  IconEye,
  IconEyeOff,
  IconCopy,
  IconTrash,
  IconArrowUp,
  IconArrowDown,
} from '@tabler/icons-vue'
import type { Node } from '@/types/page-builder'

const props = defineProps<{ node: Node; disabled?: boolean; isRoot?: boolean }>()
const emit = defineEmits<{
  'ask-agent': [Node]
  'toggle-hidden': [Node]
  duplicate: [Node]
  remove: [Node]
  move: [Node, dir: -1 | 1]
}>()
</script>

<template>
  <ContextMenu v-if="!disabled">
    <ContextMenuTrigger as-child>
      <slot />
    </ContextMenuTrigger>
    <ContextMenuContent class="w-48">
      <ContextMenuItem @select="emit('ask-agent', props.node)">
        <IconRobotFace class="mr-2 size-3.5" /> Ask agent
      </ContextMenuItem>
      <ContextMenuItem @select="emit('toggle-hidden', props.node)">
        <component :is="props.node.hidden ? IconEye : IconEyeOff" class="mr-2 size-3.5" />
        {{ props.node.hidden ? 'Show' : 'Hide' }}
      </ContextMenuItem>
      <ContextMenuSeparator />
      <ContextMenuItem
        :disabled="props.isRoot"
        :class="{ 'opacity-40 pointer-events-none': props.isRoot }"
        @select="!props.isRoot && emit('duplicate', props.node)"
      >
        <IconCopy class="mr-2 size-3.5" /> Duplicate
      </ContextMenuItem>
      <ContextMenuItem
        :disabled="props.isRoot"
        :class="{ 'opacity-40 pointer-events-none': props.isRoot }"
        @select="!props.isRoot && emit('move', props.node, -1)"
      >
        <IconArrowUp class="mr-2 size-3.5" /> Move up
      </ContextMenuItem>
      <ContextMenuItem
        :disabled="props.isRoot"
        :class="{ 'opacity-40 pointer-events-none': props.isRoot }"
        @select="!props.isRoot && emit('move', props.node, 1)"
      >
        <IconArrowDown class="mr-2 size-3.5" /> Move down
      </ContextMenuItem>
      <ContextMenuSeparator />
      <ContextMenuItem
        :disabled="props.isRoot"
        :class="{ 'opacity-40 pointer-events-none': props.isRoot, 'text-destructive focus:text-destructive': !props.isRoot }"
        @select="!props.isRoot && emit('remove', props.node)"
      >
        <IconTrash class="mr-2 size-3.5" /> Delete
      </ContextMenuItem>
    </ContextMenuContent>
  </ContextMenu>
  <template v-else><slot /></template>
</template>
