<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { Separator } from '@/components/ui/separator'
import { SidebarTrigger } from '@/components/ui/sidebar'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'

const route = useRoute()

type Crumb = { title: string; to?: string }

// Route-name based labels with optional parent (for editor pages).
const labelMap: Record<string, { label: string; parent?: { title: string; to: string } }> = {
  dashboard: { label: 'Dashboard' },
  posts: { label: 'Posts' },
  'post-new': { label: 'New Post', parent: { title: 'Posts', to: '/posts' } },
  'post-edit': { label: 'Edit Post', parent: { title: 'Posts', to: '/posts' } },
  projects: { label: 'Projects' },
  'project-new': { label: 'New Project', parent: { title: 'Projects', to: '/projects' } },
  'project-edit': { label: 'Edit Project', parent: { title: 'Projects', to: '/projects' } },
  'open-source': { label: 'Open Source' },
  'open-source-new': { label: 'New Item', parent: { title: 'Open Source', to: '/open-source' } },
  'open-source-edit': { label: 'Edit Item', parent: { title: 'Open Source', to: '/open-source' } },
  products: { label: 'Products' },
  'product-new': { label: 'New Product', parent: { title: 'Products', to: '/products' } },
  'product-edit': { label: 'Edit Product', parent: { title: 'Products', to: '/products' } },
  inbox: { label: 'Inbox' },
  pages: { label: 'Pages' },
  users: { label: 'Users' },
  'user-new': { label: 'New User', parent: { title: 'Users', to: '/users' } },
  'user-edit': { label: 'Edit User', parent: { title: 'Users', to: '/users' } },
  profile: { label: 'Profile' },
}

const crumbs = computed<Crumb[]>(() => {
  const name = String(route.name || '')
  const entry = labelMap[name]
  if (!entry) return [{ title: 'Dashboard', to: '/' }]
  const out: Crumb[] = [{ title: 'Dashboard', to: '/' }]
  if (entry.parent) out.push({ title: entry.parent.title, to: entry.parent.to })
  out.push({ title: entry.label })
  return out
})
</script>

<template>
  <header class="flex h-14 shrink-0 items-center gap-2 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
    <div class="flex w-full items-center gap-2 px-4 lg:px-6">
      <SidebarTrigger class="-ml-1 size-7" />
      <Separator orientation="vertical" class="mr-1 data-[orientation=vertical]:h-4" />

      <Breadcrumb>
        <BreadcrumbList class="text-sm">
          <template v-for="(crumb, i) in crumbs" :key="i">
            <BreadcrumbSeparator v-if="i > 0" />
            <BreadcrumbItem>
              <BreadcrumbLink v-if="crumb.to" as-child class="text-muted-foreground">
                <RouterLink :to="crumb.to">{{ crumb.title }}</RouterLink>
              </BreadcrumbLink>
              <BreadcrumbPage v-else class="font-medium">{{ crumb.title }}</BreadcrumbPage>
            </BreadcrumbItem>
          </template>
        </BreadcrumbList>
      </Breadcrumb>
    </div>
  </header>
</template>
