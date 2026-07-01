<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
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
import { IconPlus, IconEdit, IconTrash } from '@tabler/icons-vue'

const router = useRouter()
const posts = ref<any[]>([])
const loading = ref(true)

function fmtDate(s: string | null) {
  if (!s) return '—'
  return new Date(s).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' })
}

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/posts')
    posts.value = data.data ?? []
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load posts')
  } finally {
    loading.value = false
  }
}

async function remove(id: number) {
  if (!confirm('Delete this post?')) return
  try {
    await api.delete(`/posts/${id}`)
    posts.value = posts.value.filter((p) => p.id !== id)
    toast.success('Post deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
  }
}

onMounted(load)
</script>

<template>
  <div class="p-8">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="font-heading text-2xl font-bold">Blog Posts</h1>
      <Button @click="router.push('/admin/posts/new')">
        <IconPlus class="size-4" />
        New Post
      </Button>
    </div>

    <div class="rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Title (EN)</TableHead>
            <TableHead>Title (ID)</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Published</TableHead>
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-if="loading">
            <TableCell :colspan="5">
              <Skeleton class="h-6 w-full" />
            </TableCell>
          </TableRow>
          <TableRow v-else-if="posts.length === 0">
            <TableCell :colspan="5" class="text-center text-muted-foreground">
              No posts yet.
            </TableCell>
          </TableRow>
          <TableRow v-for="post in posts" :key="post.id">
            <TableCell class="font-medium">{{ post.title_en || '—' }}</TableCell>
            <TableCell class="text-muted-foreground">{{ post.title_id || '—' }}</TableCell>
            <TableCell>
              <Badge :variant="post.status === 'published' ? 'default' : 'secondary'">
                {{ post.status }}
              </Badge>
            </TableCell>
            <TableCell>{{ fmtDate(post.published_at) }}</TableCell>
            <TableCell class="text-right">
              <div class="flex justify-end gap-1">
                <Button variant="ghost" size="icon-sm" @click="router.push(`/admin/posts/${post.id}/edit`)">
                  <IconEdit class="size-4" />
                </Button>
                <Button variant="ghost" size="icon-sm" @click="remove(post.id)">
                  <IconTrash class="size-4 text-destructive" />
                </Button>
              </div>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
