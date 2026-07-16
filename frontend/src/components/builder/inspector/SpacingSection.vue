<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import { replaceClass, currentFromSet, SPACING } from '@/types/tokens'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import InspectorSection from './InspectorSection.vue'
import { IconArrowsHorizontal } from '@tabler/icons-vue'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

function spacingOpts(prefix: string) {
  const cands = SPACING.map((s) => `${prefix}-${s}`)
  const v = currentFromSet(node.value!.classes, cands)?.slice(prefix.length + 1)
  return {
    model: v ?? '0',
    set: (val: unknown) => cls([prefix], val === '0' ? null : `${prefix}-${String(val)}`),
  }
}

// Two groups: padding (p/px/py) and margin (m/mx/my). Rendered as compact
// label-left / control-right rows.
const groups = [
  { title: 'Padding', items: [
    { prefix: 'p', label: 'All' },
    { prefix: 'px', label: 'Horizontal' },
    { prefix: 'py', label: 'Vertical' },
  ]},
  { title: 'Margin', items: [
    { prefix: 'm', label: 'All' },
    { prefix: 'mx', label: 'Horizontal' },
    { prefix: 'my', label: 'Vertical' },
  ]},
]
</script>

<template>
  <InspectorSection title="Spacing" :icon="IconArrowsHorizontal" :show="!!node">
    <div v-for="grp in groups" :key="grp.title" class="space-y-1">
      <div class="text-[10px] font-medium uppercase tracking-wide text-muted-foreground">{{ grp.title }}</div>
      <div
        v-for="item in grp.items"
        :key="item.prefix"
        class="flex items-center justify-between gap-2"
      >
        <span class="text-[11px] font-medium text-foreground/80">{{ item.label }}</span>
        <Select
          :model-value="spacingOpts(item.prefix).model"
          @update:model-value="spacingOpts(item.prefix).set"
        >
          <SelectTrigger class="h-8 w-20 px-2 text-xs"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in SPACING" :key="s" :value="s">{{ s }}</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  </InspectorSection>
</template>
