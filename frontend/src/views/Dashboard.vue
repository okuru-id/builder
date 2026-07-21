<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import {
  IconArticle,
  IconMail,
  IconPackage,
  IconBriefcase,
  IconArrowRight,
  IconPlus,
  IconClock,
} from '@tabler/icons-vue'
import api from '@/lib/api'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

interface Post { id: number; title_en: string; title_id: string; status: string; published_at: string | null }
interface Message { id: number; name: string; email: string; message: string; status: string; created_at?: string }
interface Product { id: number }
interface Project { id: number; featured?: boolean }

const posts = ref<Post[]>([])
const messages = ref<Message[]>([])
const products = ref<Product[]>([])
const projects = ref<Project[]>([])
const loading = ref(true)
const me = ref<{ name?: string; email?: string }>({})

const unread = computed(() => messages.value.filter((m) => m.status === 'unread').length)
const publishedPosts = computed(() => posts.value.filter((p) => p.status === 'published').length)
const featuredProjects = computed(() => projects.value.filter((p) => p.featured).length)
const recentPosts = computed(() => posts.value.slice(0, 4))
const recentMessages = computed(() => messages.value.slice(0, 4))

function fmtDate(s: string | null) {
  if (!s) return '—'
  return new Date(s).toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}

function fmtRelative(s?: string) {
  if (!s) return ''
  const diff = Date.now() - new Date(s).getTime()
  const min = Math.floor(diff / 60000)
  if (min < 60) return `${min}m`
  const hr = Math.floor(min / 60)
  if (hr < 24) return `${hr}h`
  const day = Math.floor(hr / 24)
  if (day < 30) return `${day}d`
  return fmtDate(s)
}

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 11) return 'Good morning'
  if (h < 15) return 'Good afternoon'
  if (h < 19) return 'Good evening'
  return 'Good night'
})

const today = new Date().toLocaleDateString('id-ID', {
  weekday: 'long', day: 'numeric', month: 'long', year: 'numeric',
})

const stats = computed(() => [
  { label: 'Posts', value: posts.value.length, sub: `${publishedPosts.value} published`, icon: IconArticle, to: '/posts', tint: 'text-sky-600 bg-sky-500/10' },
  { label: 'Inbox', value: unread.value, sub: `${messages.value.length} total`, icon: IconMail, to: '/inbox', tint: 'text-amber-600 bg-amber-500/10' },
  { label: 'Products', value: products.value.length, sub: 'in catalog', icon: IconPackage, to: '/products', tint: 'text-emerald-600 bg-emerald-500/10' },
  { label: 'Projects', value: projects.value.length, sub: `${featuredProjects.value} featured`, icon: IconBriefcase, to: '/projects', tint: 'text-violet-600 bg-violet-500/10' },
])

onMounted(async () => {
  try {
    const [p, m, pr, pj, meRes] = await Promise.all([
      api.get('/posts'),
      api.get('/messages'),
      api.get('/products'),
      api.get('/projects'),
      api.get('/auth/me').catch(() => null),
    ])
    posts.value = p.data.data ?? []
    messages.value = m.data.data ?? []
    products.value = pr.data.data ?? []
    projects.value = pj.data.data ?? []
    if (meRes) me.value = meRes.data?.data ?? meRes.data ?? {}
  } finally {
    loading.value = false
  }
})

const displayName = computed(() => me.value.name || me.value.email?.split('@')[0] || 'Admin')
</script>

<template>
  <div class="@container/main flex flex-col gap-6">
    <!-- Welcome header -->
    <div class="flex flex-wrap items-end justify-between gap-4">
      <div class="space-y-1">
        <p class="text-xs uppercase tracking-wider text-muted-foreground">{{ today }}</p>
        <h1 class="font-heading text-3xl font-bold tracking-tight">
          {{ greeting }}, <span class="text-primary">{{ displayName }}</span>
        </h1>
        <p class="text-sm text-muted-foreground">
          {{ unread > 0 ? `${unread} unread message${unread > 1 ? 's' : ''} waiting` : 'Inbox zero. Nice.' }}
        </p>
      </div>
      <Button as-child class="gap-1">
        <RouterLink to="/posts/new">
          <IconPlus class="size-4" />
          New Post
        </RouterLink>
      </Button>
    </div>

    <!-- Stat cards -->
    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
      <RouterLink v-for="s in stats" :key="s.label" :to="s.to" class="group">
        <Card class="relative overflow-hidden p-4 transition-colors group-hover:bg-muted/40">
          <div class="flex items-start justify-between">
            <div class="space-y-1">
              <p class="text-xs font-medium uppercase tracking-wide text-muted-foreground">{{ s.label }}</p>
              <Skeleton v-if="loading" class="h-8 w-12" />
              <div v-else class="font-heading text-3xl font-bold tabular-nums">{{ s.value }}</div>
              <p class="text-xs text-muted-foreground">{{ s.sub }}</p>
            </div>
            <div :class="['flex size-9 items-center justify-center rounded-lg', s.tint]">
              <component :is="s.icon" class="size-5" />
            </div>
          </div>
        </Card>
      </RouterLink>
    </div>

    <!-- Main grid -->
    <div class="grid gap-4 lg:grid-cols-5">
      <!-- Recent posts -->
      <Card class="lg:col-span-3 p-0">
        <div class="flex items-center justify-between border-b p-4">
          <div>
            <h2 class="font-heading font-semibold">Recent Posts</h2>
            <p class="text-xs text-muted-foreground">Your latest writings</p>
          </div>
          <Button as-child variant="ghost" size="sm" class="gap-1">
            <RouterLink to="/posts">
              All
              <IconArrowRight class="size-4" />
            </RouterLink>
          </Button>
        </div>

        <div v-if="loading" class="space-y-3 p-4">
          <Skeleton v-for="i in 4" :key="i" class="h-14 w-full" />
        </div>

        <div v-else-if="!recentPosts.length" class="flex flex-col items-center gap-2 p-10 text-center">
          <IconArticle class="size-8 text-muted-foreground/50" />
          <p class="text-sm text-muted-foreground">No posts yet.</p>
          <Button as-child size="sm" variant="outline">
            <RouterLink to="/posts/new">Write your first</RouterLink>
          </Button>
        </div>

        <ul v-else class="divide-y">
          <li v-for="post in recentPosts" :key="post.id">
            <RouterLink
              :to="`/posts/${post.id}/edit`"
              class="flex items-center gap-3 px-4 py-3 transition-colors hover:bg-muted/50"
            >
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium">
                  {{ post.title_en || post.title_id || 'Untitled' }}
                </p>
                <p class="flex items-center gap-1 text-xs text-muted-foreground">
                  <IconClock class="size-3" />
                  {{ fmtDate(post.published_at) }}
                </p>
              </div>
              <Badge
                :variant="post.status === 'published' ? 'default' : 'secondary'"
                class="text-[10px] uppercase"
              >
                {{ post.status }}
              </Badge>
            </RouterLink>
          </li>
        </ul>
      </Card>

      <!-- Recent messages -->
      <Card class="lg:col-span-2 p-0">
        <div class="flex items-center justify-between border-b p-4">
          <div>
            <h2 class="font-heading font-semibold">Inbox</h2>
            <p class="text-xs text-muted-foreground">{{ unread }} unread</p>
          </div>
          <Button as-child variant="ghost" size="sm" class="gap-1">
            <RouterLink to="/inbox">
              Open
              <IconArrowRight class="size-4" />
            </RouterLink>
          </Button>
        </div>

        <div v-if="loading" class="space-y-3 p-4">
          <Skeleton v-for="i in 4" :key="i" class="h-12 w-full" />
        </div>

        <div v-else-if="!recentMessages.length" class="flex flex-col items-center gap-2 p-10 text-center">
          <IconMail class="size-8 text-muted-foreground/50" />
          <p class="text-sm text-muted-foreground">No messages yet.</p>
        </div>

        <ul v-else class="divide-y">
          <li v-for="msg in recentMessages" :key="msg.id">
            <RouterLink
              to="/inbox"
              class="flex items-start gap-3 px-4 py-3 transition-colors hover:bg-muted/50"
            >
              <div class="flex size-8 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-semibold text-primary">
                {{ (msg.name || msg.email || '?')[0].toUpperCase() }}
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-2">
                  <p class="truncate text-sm font-medium">{{ msg.name || msg.email }}</p>
                  <span v-if="msg.created_at" class="ml-auto shrink-0 text-[10px] text-muted-foreground">
                    {{ fmtRelative(msg.created_at) }}
                  </span>
                </div>
                <p class="truncate text-xs text-muted-foreground">{{ msg.message || msg.email }}</p>
              </div>
              <span v-if="msg.status === 'unread'" class="mt-1 size-1.5 shrink-0 rounded-full bg-primary" />
            </RouterLink>
          </li>
        </ul>
      </Card>
    </div>
  </div>
</template>
