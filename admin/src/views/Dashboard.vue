<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import api from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

const stats = ref({ posts: 0, unread: 0, products: 0 })
const loading = ref(true)

onMounted(async () => {
  try {
    const [posts, messages, products] = await Promise.all([
      api.get('/posts'),
      api.get('/messages'),
      api.get('/products'),
    ])
    stats.value.posts = posts.data.data?.length ?? 0
    stats.value.unread = (messages.data.data ?? []).filter((m: any) => m.status === 'unread').length
    stats.value.products = products.data.data?.length ?? 0
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="p-8">
    <h1 class="mb-6 font-heading text-2xl font-bold">Dashboard</h1>
    <div class="grid gap-4 sm:grid-cols-3">
      <Card v-for="card in [
        { label: 'Blog Posts', value: stats.posts, to: '/admin/posts' },
        { label: 'Unread Messages', value: stats.unread, to: '/admin/inbox' },
        { label: 'Products', value: stats.products, to: '/admin/products' },
      ]" :key="card.to">
        <CardHeader>
          <CardTitle class="text-sm font-medium text-muted-foreground">{{ card.label }}</CardTitle>
        </CardHeader>
        <CardContent>
          <Skeleton v-if="loading" class="h-8 w-16" />
          <RouterLink v-else :to="card.to" class="text-3xl font-bold hover:underline">
            {{ card.value }}
          </RouterLink>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
