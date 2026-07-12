<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import { replaceClass, currentClass } from '@/types/tokens'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

function setColor(e: Event) {
  const hex = (e.target as HTMLInputElement).value
  cls(['bg'], hex ? `bg-[${hex}]` : null)
}

function currentBgHex(): string {
  const v = currentClass(node.value!.classes, 'bg')
  if (!v) return '#ffffff'
  const m = v.match(/\[(#[\da-fA-F]{3,8})\]/)
  return m?.[1] ?? '#ffffff'
}
</script>

<template>
  <div v-if="node" class="space-y-2 border-b border-neutral-100 pb-3">
    <h3 class="text-xs font-medium uppercase tracking-wider text-neutral-500">Latar</h3>

    <div class="space-y-1.5">
      <Label class="text-xs">Warna Latar</Label>
      <div class="flex items-center gap-2">
        <input
          type="color"
          :value="currentBgHex()"
          class="h-7 w-10 cursor-pointer rounded border border-neutral-300 p-0.5"
          @input="setColor"
        />
        <Input
          :model-value="currentClass(node.classes, 'bg') ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="bg-white"
          @update:model-value="(v) => cls(['bg'], String(v) || null)"
        />
      </div>
      <div class="flex flex-wrap gap-1">
        <button
          v-for="c in ['white', 'gray-50', 'gray-100', 'gray-200', 'blue-50', 'blue-100', 'blue-500', 'green-50', 'green-500', 'amber-50', 'amber-100', 'red-50', 'red-500', 'neutral-900']"
          :key="c"
          class="h-5 w-5 rounded-full border border-neutral-200"
          :class="`bg-${c}`"
          :title="c"
          @click="cls(['bg'], `bg-${c}`)"
        />
      </div>
    </div>

    <div class="space-y-1.5">
      <Label class="text-xs">Gambar Latar (class)</Label>
      <Input
        :model-value="currentClass(node.classes, 'bg-cover') ?? ''"
        class="h-8 font-mono text-xs"
        placeholder="bg-cover bg-center"
        @update:model-value="(v) => cls(['bg-cover', 'bg-contain', 'bg-repeat', 'bg-no-repeat', 'bg-center', 'bg-top', 'bg-bottom'], String(v) || null)"
      />
      <p class="text-[11px] text-neutral-400">Gunakan Tailwind kelas, misal: <code class="rounded bg-neutral-100 px-1">bg-cover bg-center</code></p>
    </div>
  </div>
</template>
