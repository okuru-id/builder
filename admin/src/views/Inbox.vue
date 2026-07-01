<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { IconMailOpened, IconArchive, IconTrash, IconX } from '@tabler/icons-vue'

const messages = ref<any[]>([])
const loading = ref(true)
const selected = ref<any | null>(null)

function statusVariant(s: string) {
  return s === 'unread' ? 'default' : s === 'read' ? 'secondary' : 'outline'
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

async function markRead(m: any) {
  try {
    await api.post(`/messages/${m.id}/read`)
    m.status = 'read'
    if (selected.value?.id === m.id) selected.value = { ...selected.value, status: 'read' }
    toast.success('Marked as read')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed')
  }
}

async function archive(m: any) {
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

function select(m: any) {
  selected.value = m
  if (m.status === 'unread') markRead(m)
}

onMounted(load)
</script>

<template>
  <div class="flex h-full">
    <!-- List -->
    <div class="flex-1 overflow-auto p-8">
      <h1 class="mb-6 font-heading text-2xl font-bold">Inbox</h1>
      <div class="rounded-lg border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Message</TableHead>
              <TableHead>Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-if="loading">
              <TableCell :colspan="4"><Skeleton class="h-6 w-full" /></TableCell>
            </TableRow>
            <TableRow v-else-if="messages.length === 0">
              <TableCell :colspan="4" class="text-center text-muted-foreground">No messages.</TableCell>
            </TableRow>
            <TableRow
              v-for="m in messages"
              :key="m.id"
              class="cursor-pointer"
              :class="{ 'bg-muted/50': selected?.id === m.id, 'font-semibold': m.status === 'unread' }"
              @click="select(m)"
            >
              <TableCell>{{ m.name }}</TableCell>
              <TableCell class="text-muted-foreground">{{ m.email }}</TableCell>
              <TableCell class="max-w-xs truncate text-muted-foreground">{{ m.message }}</TableCell>
              <TableCell><Badge :variant="statusVariant(m.status)">{{ m.status }}</Badge></TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>
    </div>

    <!-- Detail panel -->
    <aside v-if="selected" class="flex w-96 flex-col border-l bg-card">
      <div class="flex items-center justify-between border-b p-4">
        <h2 class="font-semibold">Message</h2>
        <Button variant="ghost" size="icon-sm" @click="selected = null">
          <IconX class="size-4" />
        </Button>
      </div>
      <div class="flex-1 overflow-auto p-4">
        <div class="mb-1 text-lg font-semibold">{{ selected.name }}</div>
        <a :href="`mailto:${selected.email}`" class="mb-3 block text-sm text-primary hover:underline">{{ selected.email }}</a>
        <p class="whitespace-pre-wrap text-sm">{{ selected.message }}</p>
      </div>
      <div class="flex gap-2 border-t p-4">
        <Button v-if="selected.status !== 'read'" variant="outline" size="sm" @click="markRead(selected)">
          <IconMailOpened class="size-4" /> Read
        </Button>
        <Button variant="outline" size="sm" @click="archive(selected)">
          <IconArchive class="size-4" /> Archive
        </Button>
        <Button variant="destructive" size="sm" class="ml-auto" @click="remove(selected.id)">
          <IconTrash class="size-4" /> Delete
        </Button>
      </div>
    </aside>
  </div>
</template>
