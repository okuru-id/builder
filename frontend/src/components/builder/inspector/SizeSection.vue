<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import { replaceClass, currentClass, SIZES, SPACING } from '@/types/tokens'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Input } from '@/components/ui/input'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

function sizeOpt(prefix: string) {
  // current size value: check for size-* (e.g. w-full, h-auto) or w-<number>/h-<number>
  const val = node.value!.classes.find((c) => c.startsWith(prefix + '-'))
  return val ? val.slice(prefix.length + 1) : ''
}
</script>

<template>
  <div v-if="node" class="space-y-2 border-b border-neutral-100 pb-3">
    <h3 class="text-xs font-medium uppercase tracking-wider text-neutral-500">Ukuran</h3>

    <div class="space-y-1.5">
      <Label class="text-xs">Lebar (w-*)</Label>
      <div class="flex gap-2">
        <Select
          :model-value="currentClass(node.classes, 'w') ?? ''"
          @update:model-value="(v) => cls(['w', 'min-w', 'max-w'], v ? `w-${String(v)}` : null)"
          class="flex-1"
        >
          <SelectTrigger class="h-8"><SelectValue placeholder="auto" /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in SIZES" :key="s" :value="s">{{ s }}</SelectItem>
            <SelectItem v-for="s in SPACING" :key="'w-'+s" :value="s">{{ s }}</SelectItem>
          </SelectContent>
        </Select>
        <Input
          :model-value="sizeOpt('w')"
          class="h-8 w-20 font-mono text-xs"
          placeholder="custom"
          @update:model-value="(v) => cls(['w'], v ? `w-[${String(v)}]` : null)"
        />
      </div>
    </div>

    <div class="space-y-1.5">
      <Label class="text-xs">Tinggi (h-*)</Label>
      <div class="flex gap-2">
        <Select
          :model-value="currentClass(node.classes, 'h') ?? ''"
          @update:model-value="(v) => cls(['h', 'min-h', 'max-h'], v ? `h-${String(v)}` : null)"
          class="flex-1"
        >
          <SelectTrigger class="h-8"><SelectValue placeholder="auto" /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in SIZES" :key="s" :value="s">{{ s }}</SelectItem>
            <SelectItem v-for="s in SPACING" :key="'h-'+s" :value="s">{{ s }}</SelectItem>
          </SelectContent>
        </Select>
        <Input
          :model-value="sizeOpt('h')"
          class="h-8 w-20 font-mono text-xs"
          placeholder="custom"
          @update:model-value="(v) => cls(['h'], v ? `h-[${String(v)}]` : null)"
        />
      </div>
    </div>

    <div class="space-y-1.5">
      <Label class="text-xs">Min/Max Lebar</Label>
      <div class="flex gap-2">
        <Input
          :model-value="currentClass(node.classes, 'min-w') ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="min-w-0"
          @update:model-value="(v) => cls(['min-w'], String(v) || null)"
        />
        <Input
          :model-value="currentClass(node.classes, 'max-w') ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="max-w-4xl"
          @update:model-value="(v) => cls(['max-w'], String(v) || null)"
        />
      </div>
    </div>

    <div class="space-y-1.5">
      <Label class="text-xs">Min/Max Tinggi</Label>
      <div class="flex gap-2">
        <Input
          :model-value="currentClass(node.classes, 'min-h') ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="min-h-0"
          @update:model-value="(v) => cls(['min-h'], String(v) || null)"
        />
        <Input
          :model-value="currentClass(node.classes, 'max-h') ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="max-h-screen"
          @update:model-value="(v) => cls(['max-h'], String(v) || null)"
        />
      </div>
    </div>
  </div>
</template>
