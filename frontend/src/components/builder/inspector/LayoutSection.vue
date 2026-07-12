<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import { replaceClass, hasClass, currentClass, DISPLAY_CLASSES, FLEX_DIRECTIONS, ALIGN_ITEMS, JUSTIFY_CONTENTS, SPACING } from '@/types/tokens'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}
</script>

<template>
  <div v-if="node" class="space-y-2 border-b border-neutral-100 pb-3">
    <h3 class="text-xs font-medium uppercase tracking-wider text-neutral-500">Tata Letak</h3>

    <div class="flex items-center justify-between">
      <Label class="text-xs">Flex</Label>
      <Switch
        :model-value="hasClass(node.classes, 'flex')"
        @update:model-value="(v) => cls(['flex', ...DISPLAY_CLASSES], v ? 'flex' : null)"
      />
    </div>

    <template v-if="hasClass(node.classes, 'flex')">
      <div class="space-y-1.5">
        <Label class="text-xs">Arah</Label>
        <Select
          :model-value="currentClass(node.classes, 'flex') ?? 'row'"
          @update:model-value="(v) => cls(['flex-row', 'flex-col', 'flex-row-reverse', 'flex-col-reverse'], v === 'row' ? null : `flex-${String(v)}`)"
        >
          <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="d in FLEX_DIRECTIONS" :key="d" :value="d">{{ d }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="space-y-1.5">
        <Label class="text-xs">Ratakan Vertikal</Label>
        <Select
          :model-value="currentClass(node.classes, 'items') ?? 'stretch'"
          @update:model-value="(v) => cls(['items-start', 'items-center', 'items-end', 'items-stretch', 'items-baseline'], v === 'stretch' ? null : `items-${String(v)}`)"
        >
          <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="a in ALIGN_ITEMS" :key="a" :value="a">{{ a }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="space-y-1.5">
        <Label class="text-xs">Ratakan Horizontal</Label>
        <Select
          :model-value="currentClass(node.classes, 'justify') ?? 'start'"
          @update:model-value="(v) => cls(['justify-start', 'justify-center', 'justify-end', 'justify-between', 'justify-around', 'justify-evenly'], v === 'start' ? null : `justify-${String(v)}`)"
        >
          <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="j in JUSTIFY_CONTENTS" :key="j" :value="j">{{ j }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="space-y-1.5">
        <Label class="text-xs">Jarak (gap)</Label>
        <Select
          :model-value="currentClass(node.classes, 'gap') ?? ['0']"
          @update:model-value="(v) => cls(['gap'], v === '0' ? null : `gap-${String(v)}`)"
        >
          <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in SPACING" :key="s" :value="s">{{ s }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="flex items-center justify-between">
        <Label class="text-xs">Bungkus</Label>
        <Switch
          :model-value="hasClass(node.classes, 'flex-wrap')"
          @update:model-value="(v) => cls(['flex-wrap', 'flex-nowrap', 'flex-wrap-reverse'], v ? 'flex-wrap' : null)"
        />
      </div>
    </template>
  </div>
</template>
