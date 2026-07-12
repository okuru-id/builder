<script setup lang="ts">
import { inject, computed } from 'vue'
import { BUILDER_KEY } from '../injection'
import { replaceClass, currentClass, FONT_SIZES, FONT_WEIGHTS, TEXT_ALIGNS } from '@/types/tokens'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Input } from '@/components/ui/input'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

const showTypography = computed(() =>
  node.value && ['text', 'heading', 'button', 'link'].includes(node.value.type),
)

function cls(patterns: string[], add: string | null) {
  if (!node.value) return
  store.patchNode(node.value.id, { classes: replaceClass(node.value.classes, patterns, add) })
}

function setColor(e: Event) {
  const hex = (e.target as HTMLInputElement).value
  cls(['text'], hex ? `text-[${hex}]` : null)
}

function currentHex(): string {
  const v = currentClass(node.value!.classes, 'text')
  if (!v) return '#000000'
  // extract hex from arbitrary value text-[#abc123]
  const m = v.match(/\[(#[\da-fA-F]{3,8})\]/)
  return m?.[1] ?? '#000000'
}
</script>

<template>
  <div v-if="node && showTypography" class="space-y-2 border-b border-neutral-100 pb-3">
    <h3 class="text-xs font-medium uppercase tracking-wider text-neutral-500">Tipografi</h3>

    <div class="space-y-1.5">
      <Label class="text-xs">Ukuran</Label>
      <Select
        :model-value="currentClass(node.classes, 'text') ?? 'base'"
        @update:model-value="(v) => cls(['text-xs', 'text-sm', 'text-base', 'text-lg', 'text-xl', 'text-2xl', 'text-3xl', 'text-4xl', 'text-5xl', 'text-6xl', 'text-7xl', 'text-8xl', 'text-9xl'], v === 'base' ? null : `text-${String(v)}`)"
      >
        <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="s in FONT_SIZES" :key="s" :value="s">{{ s }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="space-y-1.5">
      <Label class="text-xs">Ketebalan</Label>
      <Select
        :model-value="currentClass(node.classes, 'font') ?? 'normal'"
        @update:model-value="(v) => cls(['font-thin', 'font-extralight', 'font-light', 'font-normal', 'font-medium', 'font-semibold', 'font-bold', 'font-extrabold', 'font-black'], v === 'normal' ? null : `font-${String(v)}`)"
      >
        <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="w in FONT_WEIGHTS" :key="w" :value="w">{{ w }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="space-y-1.5">
      <Label class="text-xs">Rata Teks</Label>
      <Select
        :model-value="currentClass(node.classes, 'text-align') ?? 'left'"
        @update:model-value="(v) => cls(['text-left', 'text-center', 'text-right', 'text-justify'], v === 'left' ? null : `text-${String(v)}`)"
      >
        <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem v-for="a in TEXT_ALIGNS" :key="a" :value="a">{{ a }}</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="space-y-1.5">
      <Label class="text-xs">Warna Teks</Label>
      <div class="flex items-center gap-2">
        <input
          type="color"
          :value="currentHex()"
          class="h-7 w-10 cursor-pointer rounded border border-neutral-300 p-0.5"
          @input="setColor"
        />
        <Input
          :model-value="currentClass(node.classes, 'text') ?? ''"
          class="h-8 flex-1 font-mono text-xs"
          placeholder="text-gray-900"
          @update:model-value="(v) => cls(['text'], String(v) || null)"
        />
      </div>
      <div class="flex flex-wrap gap-1">
        <button
          v-for="c in ['gray-900', 'gray-600', 'gray-400', 'blue-600', 'blue-500', 'green-600', 'red-600', 'amber-600', 'white']"
          :key="c"
          class="h-5 w-5 rounded-full border border-neutral-200"
          :class="`bg-${c}`"
          :title="c"
          @click="cls(['text'], `text-${c}`)"
        />
      </div>
    </div>

    <div class="space-y-1.5">
      <Label class="text-xs">Kapitalisasi</Label>
      <Select
        :model-value="currentClass(node.classes, 'capitalize') ?? 'none'"
        @update:model-value="(v) => cls(['uppercase', 'lowercase', 'capitalize', 'normal-case'], v === 'none' ? null : String(v))"
      >
        <SelectTrigger class="h-8"><SelectValue /></SelectTrigger>
        <SelectContent>
          <SelectItem value="none">Normal</SelectItem>
          <SelectItem value="uppercase">UPPERCASE</SelectItem>
          <SelectItem value="lowercase">lowercase</SelectItem>
          <SelectItem value="capitalize">Capitalize</SelectItem>
        </SelectContent>
      </Select>
    </div>
  </div>
</template>
