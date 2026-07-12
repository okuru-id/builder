<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import { replaceClass, currentClass, SPACING } from '@/types/tokens'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

function spacingOpts(prefix: string) {
  const v = currentClass(node.value!.classes, prefix)
  return { model: v ?? '0', set: (val: unknown) => cls([prefix], val === '0' ? null : `${prefix}-${String(val)}`) }
}
</script>

<template>
  <div v-if="node" class="space-y-2 border-b border-neutral-100 pb-3">
    <h3 class="text-xs font-medium uppercase tracking-wider text-neutral-500">Spasi</h3>

    <!-- Padding: grouped all, x, y -->
    <div v-for="grp in [{prefix:'p',label:'Padding semua'},{prefix:'px',label:'Padding kiri/kanan'},{prefix:'py',label:'Padding atas/bawah'}]" :key="grp.prefix" class="space-y-1">
      <Label class="text-xs">{{ grp.label }}</Label>
      <Select
        :model-value="spacingOpts(grp.prefix).model"
        @update:model-value="spacingOpts(grp.prefix).set"
      >
        <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="s in SPACING" :key="s" :value="s">{{ s }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <!-- Margin: grouped all, x, y -->
    <div v-for="grp in [{prefix:'m',label:'Margin semua'},{prefix:'mx',label:'Margin kiri/kanan'},{prefix:'my',label:'Margin atas/bawah'}]" :key="grp.prefix" class="space-y-1">
      <Label class="text-xs">{{ grp.label }}</Label>
      <Select
        :model-value="spacingOpts(grp.prefix).model"
        @update:model-value="spacingOpts(grp.prefix).set"
      >
        <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="s in SPACING" :key="s" :value="s">{{ s }}</SelectItem>
        </SelectContent>
      </Select>
    </div>
  </div>
</template>
