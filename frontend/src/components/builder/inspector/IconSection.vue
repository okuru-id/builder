<script setup lang="ts">
// Icon picker for icon-type nodes. Lazy-loads the full Tabler icon catalog
// (6203 icons, ~2.3MB) as a separate Vite chunk only when this section opens.
// Live search by name; grid caps visible rows for perf (no virtualization).
import { inject, computed, ref } from 'vue'
import { BUILDER_KEY } from '../injection'
import InspectorSection from './InspectorSection.vue'
import { Input } from '@/components/ui/input'
import { IconStar, IconSearch } from '@tabler/icons-vue'
import { ICONS, ICON_LIST } from '@/lib/icon-map'

const store = inject(BUILDER_KEY, null)!
const node = computed(() => store.selectedNode.value)

const query = ref('')
const limit = ref(60) // visible cap; bump on "show more"

const allMatches = computed<string[]>(() => {
  const q = query.value.trim().toLowerCase().replace(/^icon/, '')
  if (!q) return ICON_LIST
  return ICON_LIST.filter((n) => n.toLowerCase().replace(/^icon/, '').includes(q))
})
const filtered = computed(() => allMatches.value.slice(0, limit.value))

function setIcon(name: string) {
  store.patchNode(node.value!.id, { props: { ...node.value!.props, icon: name } })
}
</script>

<template>
  <InspectorSection title="Icon" :icon="IconStar" :show="!!node">
    <!-- Search -->
    <div class="relative">
      <IconSearch class="pointer-events-none absolute left-2 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
      <Input
        v-model="query"
        class="h-8 pl-7 text-xs"
        placeholder="Search 6200+ icons…"
      />
    </div>

    <!-- Grid -->
    <div class="grid grid-cols-6 gap-1">
      <button
        v-for="name in filtered"
        :key="name"
        class="flex h-8 w-8 items-center justify-center rounded-md border text-muted-foreground transition-colors"
        :class="node?.props?.icon === name ? 'border-primary bg-primary/10 text-primary' : 'border-border hover:bg-muted'"
        :title="name"
        @click="setIcon(name)"
      >
        <svg v-if="ICONS[name]" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round" class="size-4">
          <path v-for="(seg, i) in ICONS[name]" :key="i" :d="seg[1].d" />
        </svg>
      </button>
    </div>

    <!-- Footer: count + show more/less -->
    <div class="flex items-center justify-between text-[10px] text-muted-foreground">
      <span>{{ filtered.length }} / {{ allMatches.length }} icons</span>
      <button
        v-if="allMatches.length > limit"
        class="font-medium text-primary hover:underline"
        @click="limit += 60"
      >
        Show more…
      </button>
      <button
        v-else-if="limit > 60"
        class="font-medium text-primary hover:underline"
        @click="limit = 60"
      >
        Show less
      </button>
    </div>
  </InspectorSection>
</template>
