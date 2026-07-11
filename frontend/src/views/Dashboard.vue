<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import {
  IconArticle,
  IconMail,
  IconPackage,
  IconBriefcase,
  IconArrowUpRight,
} from '@tabler/icons-vue'
import api from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

interface Post { id: number; title_en: string; title_id: string; status: string; published_at: string | null }
interface Message { id: number; name: string; email: string; message: string; status: string }
interface Product { id: number }
interface Project { id: number; featured?: boolean }

const posts = ref<Post[]>([])
const messages = ref<Message[]>([])
const products = ref<Product[]>([])
const projects = ref<Project[]>([])
const loading = ref(true)

const unread = computed(() => messages.value.filter((m) => m.status === 'unread').length)
const recentPosts = computed(() => posts.value.slice(0, 5))
const recentMessages = computed(() => messages.value.slice(0, 5))

function fmtDate(s: string | null) {
  if (!s) return '—'
  return new Date(s).toLocaleDateString('id-ID', { year: 'numeric', month: 'short', day: 'numeric' })
}

const initials = (n: string) =>
  n.split(' ').map((w) => w.charAt(0).toUpperCase()).slice(0, 2).join('')

onMounted(async () => {
  try {
    const [p, m, pr, pj] = await Promise.all([
      api.get('/posts'),
      api.get('/messages'),
      api.get('/products'),
      api.get('/projects'),
    ])
    posts.value = p.data.data ?? []
    messages.value = m.data.data ?? []
    products.value = pr.data.data ?? []
    projects.value = pj.data.data ?? []
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="@container/main flex flex-col gap-4">
    <!-- Stat cards -->
    <div class="grid gap-4 md:grid-cols-2 md:gap-8 lg:grid-cols-4">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Blog Posts</CardTitle>
          <IconArticle class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <Skeleton v-if="loading" class="h-8 w-12" />
          <div v-else class="text-2xl font-bold">{{ posts.length }}</div>
          <p class="text-xs text-muted-foreground">
            {{ posts.filter((p) => p.status === 'published').length }} published
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Unread Messages</CardTitle>
          <IconMail class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <Skeleton v-if="loading" class="h-8 w-12" />
          <div v-else class="text-2xl font-bold">{{ unread }}</div>
          <p class="text-xs text-muted-foreground">of {{ messages.length }} total</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Products</CardTitle>
          <IconPackage class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <Skeleton v-if="loading" class="h-8 w-12" />
          <div v-else class="text-2xl font-bold">{{ products.length }}</div>
          <p class="text-xs text-muted-foreground">in catalog</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Projects</CardTitle>
          <IconBriefcase class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <Skeleton v-if="loading" class="h-8 w-12" />
          <div v-else class="text-2xl font-bold">{{ projects.length }}</div>
          <p class="text-xs text-muted-foreground">{{ projects.filter((p) => p.featured).length }} featured</p>
        </CardContent>
      </Card>
    </div>

    <!-- Table + Recent messages -->
    <div class="grid gap-4 md:gap-8 lg:grid-cols-2 xl:grid-cols-3">
      <Card class="xl:col-span-2">
        <CardHeader class="flex flex-row items-center">
          <div class="grid gap-2">
            <CardTitle>Recent Posts</CardTitle>
            <CardDescription>Latest blog posts from your site.</CardDescription>
          </div>
          <Button as-child size="sm" class="ml-auto gap-1">
            <RouterLink to="/posts">
              View All
              <IconArrowUpRight class="h-4 w-4" />
            </RouterLink>
          </Button>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Title</TableHead>
                <TableHead class="hidden xl:table-column">Status</TableHead>
                <TableHead class="hidden md:table-cell lg:hidden xl:table-column">Published</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="loading">
                <TableCell colspan="3"><Skeleton class="h-5 w-full" /></TableCell>
              </TableRow>
              <TableRow v-else-if="!recentPosts.length">
                <TableCell colspan="3" class="text-muted-foreground">No posts yet.</TableCell>
              </TableRow>
              <TableRow v-for="post in recentPosts" :key="post.id">
                <TableCell class="font-medium">
                  <RouterLink :to="`/posts/${post.id}/edit`" class="hover:underline">
                    {{ post.title_en || post.title_id || 'Untitled' }}
                  </RouterLink>
                </TableCell>
                <TableCell class="hidden xl:table-column">
                  <Badge class="text-xs" variant="outline">
                    {{ post.status }}
                  </Badge>
                </TableCell>
                <TableCell class="hidden md:table-cell lg:hidden xl:table-column text-muted-foreground">
                  {{ fmtDate(post.published_at) }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center">
          <CardTitle>Recent Messages</CardTitle>
          <Button as-child size="sm" variant="outline" class="ml-auto gap-1">
            <RouterLink to="/inbox">
              Inbox
              <IconArrowUpRight class="h-4 w-4" />
            </RouterLink>
          </Button>
        </CardHeader>
        <CardContent class="grid gap-8">
          <div v-if="loading" class="grid gap-8">
            <Skeleton v-for="i in 5" :key="i" class="h-9 w-full" />
          </div>
          <div v-else-if="!recentMessages.length" class="text-muted-foreground">
            No messages yet.
          </div>
          <div
            v-for="message in recentMessages"
            :key="message.id"
            class="flex items-center gap-4"
          >
            <div
              class="bg-muted text-muted-foreground hidden h-9 w-9 shrink-0 items-center justify-center rounded-full text-xs font-medium sm:flex"
            >
              {{ initials(message.name) }}
            </div>
            <div class="grid gap-1">
              <p class="text-sm font-medium leading-none">{{ message.name }}</p>
              <p class="text-sm text-muted-foreground truncate">{{ message.email }}</p>
            </div>
            <Badge v-if="message.status === 'unread'" class="ml-auto text-xs">new</Badge>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
