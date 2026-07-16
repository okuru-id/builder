<script setup lang="ts">
// Shared right-click menu for builder nodes (tree + canvas).
// Actions: point at the chat agent, toggle visibility.
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger,
} from '@/components/ui/context-menu'
import { IconRobotFace, IconEye, IconEyeOff } from '@tabler/icons-vue'
import type { Node } from '@/types/page-builder'

const props = defineProps<{ node: Node; disabled?: boolean }>()
const emit = defineEmits<{ 'ask-agent': [Node]; 'toggle-hidden': [Node] }>()
</script>

<template>
  <ContextMenu v-if="!disabled">
    <ContextMenuTrigger as-child>
      <slot />
    </ContextMenuTrigger>
    <ContextMenuContent class="w-44">
      <ContextMenuItem @select="emit('ask-agent', props.node)">
        <IconRobotFace class="mr-2 size-3.5" /> Ask agent
      </ContextMenuItem>
      <ContextMenuItem @select="emit('toggle-hidden', props.node)">
        <component :is="props.node.hidden ? IconEye : IconEyeOff" class="mr-2 size-3.5" />
        {{ props.node.hidden ? 'Show' : 'Hide' }}
      </ContextMenuItem>
    </ContextMenuContent>
  </ContextMenu>
  <template v-else><slot /></template>
</template>
