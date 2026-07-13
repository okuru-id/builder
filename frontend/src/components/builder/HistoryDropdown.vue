<script setup lang="ts">
// Version history dropdown. Loads server-side revisions on open; clicking an
// item restores that snapshot (bumps version server-side, swaps client tree,
// clears local undo/redo). ponytail: reload every open — limit 100, cheap,
// always fresh; no client cache/staleness to manage.
import { inject, ref } from 'vue'
import { IconHistory, IconLoader2, IconRotate, IconClock } from '@tabler/icons-vue'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { BUILDER_KEY } from '@/components/builder/injection'
import { emptyRoot } from '@/components/builder/tree-utils'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import type { Page } from '@/types/page-builder'

interface Revision {
  id: number
  message: string
  tree: unknown
  created_at?: string
  updated_at?: string
}

const store = inject(BUILDER_KEY, null)!

const revisions = ref<Revision[]>([])
const loading = ref(false)
const open = ref(false)
const restoringId = ref<number | null>(null)

async function onOpenChange(v: boolean) {
  open.value = v
  if (v) await load()
}

async function load() {
  if (!store.page.value) return
  loading.value = true
  try {
    const res = await api.get<{ data: Revision[] }>(
      `/landing-pages/${store.page.value.id}/revisions`,
    )
    revisions.value = res.data.data ?? []
  } catch {
    toast.error('Failed to load history')
  } finally {
    loading.value = false
  }
}

// Relative "x ago" — ponytail: no date lib.
function timeAgo(iso?: string): string {
  if (!iso) return ''
  const then = new Date(iso).getTime()
  if (Number.isNaN(then)) return ''
  const s = Math.floor((Date.now() - then) / 1000)
  if (s < 60) return 'just now'
  const m = Math.floor(s / 60)
  if (m < 60) return `${m}m ago`
  const h = Math.floor(m / 60)
  if (h < 24) return `${h}h ago`
  const d = Math.floor(h / 24)
  if (d < 30) return `${d}d ago`
  return new Date(iso).toLocaleDateString()
}

async function restore(rev: Revision) {
  if (!store.page.value || restoringId.value !== null) return
  restoringId.value = rev.id
  try {
    const res = await api.post<{ data: Page }>(
      `/landing-pages/${store.page.value.id}/revisions/${rev.id}/restore`,
    )
    const page = res.data.data
    store.page.value = page
    store.tree.value = page.tree ?? { root: emptyRoot() }
    store.selectedId.value = null
    store.resetHistory()
    toast.success('Restored to previous version')
  } catch {
    toast.error('Failed to restore')
  } finally {
    restoringId.value = null
  }
}
</script>

<template>
  <DropdownMenu :open="open" @update:open="onOpenChange">
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" size="icon" title="Version history">
        <IconHistory class="size-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end" class="flex max-h-80 w-72 flex-col overflow-hidden p-0">
      <!-- header pinned outside the scroll area: never scrolls, zero bleed. -->
      <DropdownMenuLabel class="flex shrink-0 items-center gap-2 px-2 py-1.5 text-xs">
        <IconHistory class="size-3.5" /> Version history
      </DropdownMenuLabel>
      <DropdownMenuSeparator class="shrink-0" />

      <!-- scroll area for the revision list -->
      <div class="min-h-0 flex-1 overflow-y-auto p-1">
        <!-- loading -->
        <div v-if="loading" class="flex items-center justify-center gap-2 py-6 text-xs text-muted-foreground">
          <IconLoader2 class="size-4 animate-spin" /> Loading…
        </div>

        <!-- empty -->
        <div v-else-if="revisions.length === 0" class="px-2 py-6 text-center text-xs text-muted-foreground">
          No saved versions yet.
        </div>

        <!-- list -->
        <template v-else>
          <DropdownMenuItem
            v-for="rev in revisions"
            :key="rev.id"
            class="flex items-start gap-2 py-2"
            :disabled="restoringId !== null"
            @select="restore(rev)"
          >
            <IconClock class="mt-0.5 size-3.5 shrink-0 text-muted-foreground" />
            <div class="min-w-0 flex-1">
              <div class="truncate text-xs font-medium">{{ rev.message || `Revision #${rev.id}` }}</div>
              <div class="text-[10px] text-muted-foreground">{{ timeAgo(rev.created_at || rev.updated_at) }}</div>
            </div>
            <IconRotate
              v-if="restoringId !== rev.id"
              class="mt-0.5 size-3.5 shrink-0 text-muted-foreground"
            />
            <IconLoader2 v-else class="mt-0.5 size-3.5 shrink-0 animate-spin" />
          </DropdownMenuItem>
        </template>
      </div>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
