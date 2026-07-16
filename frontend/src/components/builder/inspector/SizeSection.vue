<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import { replaceClass, currentClass, currentArbitrary, SIZES, SPACING } from '@/types/tokens'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Input } from '@/components/ui/input'
import InspectorSection from './InspectorSection.vue'
import { IconResize } from '@tabler/icons-vue'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

// Named size for the Select (w-full, h-auto, w-32). Excludes arbitrary w-[…]
// so the dropdown never shows raw bracket syntax.
function namedOpt(prefix: string): string {
  const found = node.value!.classes.find((c) => c.startsWith(prefix + '-') && !c.startsWith(prefix + '-['))
  return found ? found.slice(prefix.length + 1) : ''
}
// Arbitrary custom value for the text input (w-[320px] → "320px", no brackets).
// Fixes bracket leak that produced w-[[320px]] on retype.
function customOpt(prefix: string): string {
  return currentArbitrary(node.value!.classes, prefix) ?? ''
}
</script>

<template>
  <InspectorSection title="Size" :icon="IconResize" :show="!!node" :default-open="false">

    <div class="space-y-1">
      <Label class="text-[11px] font-medium text-foreground/80">Width (w-*)</Label>
      <div class="flex gap-2">
        <Select
          :model-value="namedOpt('w')"
          @update:model-value="(v) => cls(['w', 'min-w', 'max-w'], v ? `w-${String(v)}` : null)"
          class="flex-1"
        >
          <SelectTrigger class="h-8 text-xs"><SelectValue placeholder="auto" /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in SIZES" :key="s" :value="s">{{ s }}</SelectItem>
            <SelectItem v-for="s in SPACING" :key="'w-'+s" :value="s">{{ s }}</SelectItem>
          </SelectContent>
        </Select>
        <Input
          :model-value="customOpt('w')"
          class="h-8 w-20 font-mono text-xs"
          placeholder="custom"
          @update:model-value="(v) => cls(['w'], v ? `w-[${String(v)}]` : null)"
        />
      </div>
    </div>

    <div class="space-y-1">
      <Label class="text-[11px] font-medium text-foreground/80">Height (h-*)</Label>
      <div class="flex gap-2">
        <Select
          :model-value="namedOpt('h')"
          @update:model-value="(v) => cls(['h', 'min-h', 'max-h'], v ? `h-${String(v)}` : null)"
          class="flex-1"
        >
          <SelectTrigger class="h-8 text-xs"><SelectValue placeholder="auto" /></SelectTrigger>
          <SelectContent>
            <SelectItem v-for="s in SIZES" :key="s" :value="s">{{ s }}</SelectItem>
            <SelectItem v-for="s in SPACING" :key="'h-'+s" :value="s">{{ s }}</SelectItem>
          </SelectContent>
        </Select>
        <Input
          :model-value="customOpt('h')"
          class="h-8 w-20 font-mono text-xs"
          placeholder="custom"
          @update:model-value="(v) => cls(['h'], v ? `h-[${String(v)}]` : null)"
        />
      </div>
    </div>

    <div class="space-y-1">
      <Label class="text-[11px] font-medium text-foreground/80">Min/Max Width</Label>
      <div class="flex gap-2">
        <Input
          :model-value="currentClass(node?.classes ?? [], 'min-w') ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="min-w-0"
          @update:model-value="(v) => cls(['min-w'], String(v) || null)"
        />
        <Input
          :model-value="currentClass(node?.classes ?? [], 'max-w') ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="max-w-4xl"
          @update:model-value="(v) => cls(['max-w'], String(v) || null)"
        />
      </div>
    </div>

    <div class="space-y-1">
      <Label class="text-[11px] font-medium text-foreground/80">Min/Max Height</Label>
      <div class="flex gap-2">
        <Input
          :model-value="currentClass(node?.classes ?? [], 'min-h') ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="min-h-0"
          @update:model-value="(v) => cls(['min-h'], String(v) || null)"
        />
        <Input
          :model-value="currentClass(node?.classes ?? [], 'max-h') ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="max-h-screen"
          @update:model-value="(v) => cls(['max-h'], String(v) || null)"
        />
      </div>
    </div>
  </InspectorSection>
</template>
