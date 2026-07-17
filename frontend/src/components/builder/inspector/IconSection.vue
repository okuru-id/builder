<script setup lang="ts">
// Icon picker for icon-type nodes. Lazy-loads the full Tabler icon catalog
// (outline + filled) as a separate Vite chunk only when this section opens.
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
// 'outline' shows IconX; 'filled' shows IconXFilled. Drives both the picker
// grid and the node's own iconVariant prop.
const variant = ref<'outline' | 'filled'>(node.value?.props?.iconVariant === 'filled' ? 'filled' : 'outline')

// Base list for the active variant: strip the `Icon` prefix and (for filled)
// the `Filled` suffix so search feels natural. Keys stay full ("IconStarFilled").
const baseList = computed<string[]>(() => {
  if (variant.value === 'filled') {
    return ICON_LIST.filter((n) => n.endsWith('Filled'))
  }
  return ICON_LIST.filter((n) => !n.endsWith('Filled'))
})

const allMatches = computed<string[]>(() => {
  const q = query.value.trim().toLowerCase().replace(/^icon/, '').replace(/filled$/, '')
  if (!q) return baseList.value
  return baseList.value.filter((n) =>
    n.toLowerCase().replace(/^icon/, '').replace(/filled$/, '').includes(q),
  )
})
const filtered = computed(() => allMatches.value.slice(0, limit.value))

// Lookup key into ICONS for the selected variant. Filled icons are stored
// under `${name}Filled`.
function iconKey(name: string): string {
  return variant.value === 'filled' ? name : name
}

function setIcon(name: string) {
  // Store the bare icon name (without `Filled` suffix) + variant; renderer
  // composes the lookup key from both. Keeps data portable if a filled
  // variant doesn't exist for a future icon.
  const bare = variant.value === 'filled' ? name.replace(/Filled$/, '') : name
  store.patchNode(node.value!.id, {
    props: { ...node.value!.props, icon: bare, iconVariant: variant.value },
  })
}

function setVariant(v: 'outline' | 'filled') {
  variant.value = v
  store.patchNode(node.value!.id, {
    props: { ...node.value!.props, iconVariant: v },
  })
  limit.value = 60
}

// Currently selected grid highlight: compare bare name + variant.
function isActive(name: string): boolean {
  const bare = variant.value === 'filled' ? name.replace(/Filled$/, '') : name
  return (
    node.value?.props?.icon === bare &&
    (node.value?.props?.iconVariant ?? 'outline') === variant.value
  )
}
</script>

<template>
  <InspectorSection title="Icon" :icon="IconStar" :show="!!node">
    <!-- Variant toggle -->
    <div class="inline-flex h-8 items-center rounded-lg border border-border bg-muted p-0.5 text-xs">
      <button
        class="rounded-md px-2.5 py-1 transition-colors"
        :class="variant === 'outline' ? 'bg-background shadow-sm font-medium' : 'text-muted-foreground'"
        @click="setVariant('outline')"
      >
        Outline
      </button>
      <button
        class="rounded-md px-2.5 py-1 transition-colors"
        :class="variant === 'filled' ? 'bg-background shadow-sm font-medium' : 'text-muted-foreground'"
        @click="setVariant('filled')"
      >
        Filled
      </button>
    </div>

    <!-- Search -->
    <div class="relative">
      <IconSearch class="pointer-events-none absolute left-2 top-1/2 size-3.5 -translate-y-1/2 text-muted-foreground" />
      <Input
        v-model="query"
        class="h-8 pl-7 text-xs"
        :placeholder="`Search ${variant === 'filled' ? '1000+ filled' : '6200+'} icons…`"
      />
    </div>

    <!-- Grid -->
    <div class="grid grid-cols-6 gap-1">
      <button
        v-for="name in filtered"
        :key="name"
        class="flex h-8 w-8 items-center justify-center rounded-md border text-muted-foreground transition-colors"
        :class="isActive(name) ? 'border-primary bg-primary/10 text-primary' : 'border-border hover:bg-muted'"
        :title="name"
        @click="setIcon(name)"
      >
        <svg v-if="ICONS[iconKey(name)]" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" :stroke="variant === 'filled' ? 'none' : 'currentColor'" :fill="variant === 'filled' ? 'currentColor' : 'none'" stroke-linecap="round" stroke-linejoin="round" class="size-4">
          <path v-for="(seg, i) in ICONS[iconKey(name)]" :key="i" :d="seg[1].d" />
        </svg>
      </button>
    </div>

    <!-- Footer: count + show more/less -->
    <div class="flex items-center justify-between text-[10px] text-muted-foreground">
      <span>{{ filtered.length }} / {{ allMatches.length }} icons</span>
      <div class="flex gap-2">
        <button
          v-if="limit > 60"
          class="font-medium text-primary hover:underline"
          @click="limit = 60"
        >
          Show less
        </button>
        <button
          v-if="allMatches.length > limit"
          class="font-medium text-primary hover:underline"
          @click="limit += 60"
        >
          Show more…
        </button>
      </div>
    </div>
  </InspectorSection>
</template>
