<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import { replaceClass, hasClass, currentClass, BORDER_WIDTHS, BORDER_RADII } from '@/types/tokens'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Input } from '@/components/ui/input'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

function setBorderColor(e: Event) {
  const hex = (e.target as HTMLInputElement).value
  cls(['border-gray', 'border-blue', 'border-green', 'border-red', 'border-amber', 'border-neutral', 'border', 'border-white', 'border-black'], hex ? `border-[${hex}]` : null)
  // Ensure border is visible when color is set
  if (hex && !hasClass(node.value!.classes, 'border')) {
    store.patchNode(node.value!.id, { classes: [...node.value!.classes, 'border'] })
  }
}

function currentBorderHex(): string {
  const v = currentClass(node.value!.classes, 'border')
  if (!v) return '#000000'
  const m = v.match(/\[(#[\da-fA-F]{3,8})\]/)
  return m?.[1] ?? '#000000'
}
</script>

<template>
  <div v-if="node" class="space-y-2 border-b border-neutral-100 pb-3">
    <h3 class="text-xs font-medium uppercase tracking-wider text-neutral-500">Garis</h3>

    <div class="flex items-center justify-between">
      <Label class="text-xs">Tampilkan Garis</Label>
      <Switch
        :model-value="hasClass(node.classes, 'border')"
        @update:model-value="(v) => cls(['border', 'border-0', 'border-2', 'border-4', 'border-8'], v ? 'border' : null)"
      />
    </div>

    <template v-if="hasClass(node.classes, 'border')">
      <div class="space-y-1.5">
        <Label class="text-xs">Tebal</Label>
        <Select
          :model-value="currentClass(node.classes, 'border') ?? '1'"
          @update:model-value="(v) => cls(['border-0', 'border-2', 'border-4', 'border-8'], v === '1' ? null : `border-${String(v)}`)"
        >
          <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="w in BORDER_WIDTHS" :key="w" :value="w">{{ w === '1' ? 'default' : w }}</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div class="space-y-1.5">
        <Label class="text-xs">Warna Garis</Label>
        <div class="flex items-center gap-2">
          <input
            type="color"
            :value="currentBorderHex()"
            class="h-7 w-10 cursor-pointer rounded border border-neutral-300 p-0.5"
            @input="setBorderColor"
          />
          <Input
            :model-value="currentClass(node.classes, 'border') ?? ''"
            class="h-8 flex-1 font-mono text-xs"
            placeholder="border-gray-300"
            @update:model-value="(v) => cls(['border'], String(v) || null)"
          />
        </div>
      </div>

      <div class="space-y-1.5">
        <Label class="text-xs">Sudut Lengkung</Label>
        <Select
          :model-value="currentClass(node.classes, 'rounded') ?? 'none'"
          @update:model-value="(v) => cls(['rounded-none', 'rounded-sm', 'rounded', 'rounded-md', 'rounded-lg', 'rounded-xl', 'rounded-2xl', 'rounded-3xl', 'rounded-full'], v === 'none' ? null : (v === 'md' ? 'rounded-md' : `rounded-${String(v)}`))"
        >
          <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="r in BORDER_RADII" :key="r" :value="r">{{ r }}</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </template>
  </div>
</template>
