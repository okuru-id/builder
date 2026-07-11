<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarInput,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'
import { IconArchive, IconMailOpened, IconInbox, IconMailbox, IconMail, IconTrash } from '@tabler/icons-vue'

interface Message {
  id: number
  name: string
  email: string
  message: string
  status: string
  created_at: string
}

const messages = ref<Message[]>([])
const loading = ref(true)
const activeFilter = ref('all')
const q = ref('')
const selected = ref<Message | null>(null)

const filters = [
  { title: 'Inbox', value: 'all', icon: IconInbox },
  { title: 'Unread', value: 'unread', icon: IconMailbox },
  { title: 'Read', value: 'read', icon: IconMail },
  { title: 'Archived', value: 'archived', icon: IconArchive },
]

const filtered = computed(() => {
  const query = q.value.trim().toLowerCase()
  return messages.value.filter((m) => {
    if (activeFilter.value !== 'all' && m.status !== activeFilter.value) return false
    if (query && !(m.name.toLowerCase().includes(query) || m.email.toLowerCase().includes(query) || m.message.toLowerCase().includes(query))) return false
    return true
  })
})

const counts = computed(() => {
  const c: Record<string, number> = { all: messages.value.length }
  for (const f of ['unread', 'read', 'archived'])
    c[f] = messages.value.filter((m) => m.status === f).length
  return c
})

function initial(name: string) {
  return name.charAt(0).toUpperCase()
}

function relativeTime(dateStr: string) {
  const d = new Date(dateStr.replace(' ', 'T') + '+08:00')
  const now = Date.now()
  const diff = now - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'now'
  if (mins < 60) return `${mins}m`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h`
  const days = Math.floor(hrs / 24)
  if (days < 7) return `${days}d`
  return d.toLocaleDateString('en-ID', { day: 'numeric', month: 'short' })
}

function messagePreview(text: string) {
  return text.length > 100 ? text.substring(0, 100) + '…' : text
}

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/messages')
    messages.value = data.data ?? []
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load')
  } finally {
    loading.value = false
  }
}

function select(m: Message) {
  selected.value = m
  if (m.status === 'unread') markRead(m)
}

async function markRead(m: Message) {
  try {
    await api.post(`/messages/${m.id}/read`)
    m.status = 'read'
    if (selected.value?.id === m.id) selected.value = { ...selected.value, status: 'read' }
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed')
  }
}

async function archive(m: Message) {
  try {
    await api.post(`/messages/${m.id}/archive`)
    m.status = 'archived'
    if (selected.value?.id === m.id) selected.value = null
    toast.success('Archived')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed')
  }
}

async function remove(id: number) {
  if (!confirm('Delete this message?')) return
  try {
    await api.delete(`/messages/${id}`)
    messages.value = messages.value.filter((m) => m.id !== id)
    if (selected.value?.id === id) selected.value = null
    toast.success('Deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
  }
}

onMounted(load)
</script>

<template>
  <div class="flex h-full overflow-hidden rounded-lg border bg-card">
    <!-- Left: sidebar-like panel -->
    <div class="bg-sidebar flex w-72 shrink-0 flex-col border-r">
      <!-- Search -->
      <div class="border-b p-3">
        <SidebarInput v-model="q" placeholder="Search messages…" class="h-8" />
      </div>

      <!-- Filter tabs with counts -->
      <SidebarGroup>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="f in filters" :key="f.value">
              <SidebarMenuButton
                :is-active="activeFilter === f.value ? true : undefined"
                class="flex w-full items-center"
                @click="activeFilter = f.value"
              >
                <component :is="f.icon" class="size-4 shrink-0" />
                <span class="flex-1 text-left">{{ f.title }}</span>
                <span
                  class="text-sidebar-foreground/50 inline-flex h-5 min-w-5 items-center justify-center rounded-full bg-sidebar-accent px-1.5 text-[11px] font-medium tabular-nums"
                >{{ counts[f.value] }}</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
      <div class="border-t" />

      <!-- Message list -->
      <div class="flex-1 overflow-auto">
        <template v-if="loading">
          <div v-for="i in 4" :key="i" class="flex items-start gap-3 border-b p-3">
            <div class="bg-muted size-8 shrink-0 animate-pulse rounded-full" />
            <div class="flex-1 space-y-1.5">
              <div class="bg-muted h-3 w-1/2 animate-pulse rounded" />
              <div class="bg-muted h-2 w-full animate-pulse rounded" />
            </div>
          </div>
        </template>
        <a
          v-for="m in filtered"
          v-else
          :key="m.id"
          href="#"
          :class="[
            'group relative flex items-start gap-3 border-b p-3 text-sm leading-tight transition-colors last:border-b-0 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground',
            selected?.id === m.id ? 'bg-sidebar-accent text-sidebar-accent-foreground' : '',
          ]"
          @click.prevent="select(m)"
        >
          <!-- Unread indicator dot -->
          <span
            v-if="m.status === 'unread'"
            class="bg-primary absolute left-1.5 top-1/2 size-1.5 -translate-y-1/2 rounded-full"
          />

          <!-- Avatar -->
          <span
            :class="[
              'flex size-8 shrink-0 items-center justify-center rounded-full text-xs font-semibold',
              m.status === 'unread'
                ? 'bg-primary/15 text-primary ring-primary/25 ring-2'
                : 'bg-muted text-muted-foreground',
            ]"
          >{{ initial(m.name) }}</span>

          <!-- Content -->
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <span
                :class="[
                  'truncate',
                  m.status === 'unread' ? 'font-semibold text-foreground' : 'font-medium text-foreground',
                ]"
              >{{ m.name }}</span>
              <span class="text-muted-foreground ml-auto shrink-0 text-[11px] tabular-nums">{{ relativeTime(m.created_at) }}</span>
            </div>
            <div class="text-muted-foreground truncate text-xs">{{ m.email }}</div>
            <div
              :class="[
                'mt-0.5 line-clamp-2 text-xs',
                m.status === 'unread' ? 'text-foreground/80' : 'text-muted-foreground',
              ]"
            >{{ messagePreview(m.message) }}</div>
          </div>
        </a>
        <div v-if="!loading && filtered.length === 0" class="text-muted-foreground flex flex-col items-center gap-2 p-8 text-center text-sm">
          <IconInbox class="text-muted-foreground/40 size-10" />
          <span>No messages</span>
        </div>
      </div>
    </div>

    <!-- Right: detail panel -->
    <div class="flex flex-1 flex-col">
      <template v-if="!selected">
        <div class="text-muted-foreground flex flex-1 flex-col items-center justify-center gap-3 text-sm">
          <IconMailOpened class="text-muted-foreground/30 size-16" />
          <span>Select a message to read</span>
        </div>
      </template>
      <template v-else>
        <div class="flex items-center justify-between border-b p-5">
          <div class="flex items-center gap-3">
            <span class="bg-primary/15 text-primary flex size-10 items-center justify-center rounded-full text-sm font-semibold">
              {{ initial(selected.name) }}
            </span>
            <div>
              <h2 class="text-lg font-semibold">{{ selected.name }}</h2>
              <a :href="`mailto:${selected.email}`" class="text-muted-foreground text-sm hover:underline">{{ selected.email }}</a>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-muted-foreground text-xs tabular-nums">{{ relativeTime(selected.created_at) }}</span>
            <Badge :variant="selected.status === 'unread' ? 'default' : 'secondary'">{{ selected.status }}</Badge>
          </div>
        </div>
        <div class="flex-1 overflow-auto p-5">
          <p class="whitespace-pre-wrap leading-relaxed">{{ selected.message }}</p>
        </div>
        <div class="flex items-center gap-2 border-t bg-muted/30 px-5 py-3">
          <Button v-if="selected.status !== 'read'" size="sm" variant="outline" @click="markRead(selected)">
            <IconMailOpened class="size-4" /> Mark read
          </Button>
          <Button size="sm" variant="outline" @click="archive(selected)">
            <IconArchive class="size-4" /> Archive
          </Button>
          <Button size="sm" variant="destructive" class="ml-auto" @click="remove(selected.id)">
            <IconTrash class="size-4" /> Delete
          </Button>
        </div>
      </template>
    </div>
  </div>
</template>
