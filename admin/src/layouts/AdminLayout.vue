<script setup lang="ts">
import { RouterView, RouterLink, useRouter } from 'vue-router'
import {
  IconLayoutDashboard,
  IconArticle,
  IconFolder,
  IconBrandGithub,
  IconShoppingBag,
  IconInbox,
  IconSettings,
  IconLogout,
} from '@tabler/icons-vue'

const router = useRouter()

const nav = [
  { to: '/admin', label: 'Dashboard', icon: IconLayoutDashboard, exact: true },
  { to: '/admin/posts', label: 'Blog Posts', icon: IconArticle },
  { to: '/admin/projects', label: 'Projects', icon: IconFolder },
  { to: '/admin/open-source', label: 'Open Source', icon: IconBrandGithub },
  { to: '/admin/products', label: 'Products', icon: IconShoppingBag },
  { to: '/admin/inbox', label: 'Inbox', icon: IconInbox },
  { to: '/admin/settings', label: 'Settings', icon: IconSettings },
]

function logout() {
  localStorage.removeItem('access_token')
  router.push('/admin/login')
}
</script>

<template>
  <div class="flex h-screen bg-background">
    <aside class="flex w-60 flex-col border-r bg-sidebar text-sidebar-foreground">
      <div class="flex h-14 items-center gap-2 border-b px-5">
        <span class="font-heading text-lg font-bold tracking-tight">okuru.id</span>
        <span class="text-xs text-muted-foreground">admin</span>
      </div>
      <nav class="flex flex-1 flex-col gap-1 p-3">
        <RouterLink
          v-for="item in nav"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
          :class="{
            'bg-sidebar-accent text-sidebar-accent-foreground':
              item.exact
                ? $route.path === item.to
                : $route.path.startsWith(item.to) && item.to !== '/admin',
          }"
        >
          <component :is="item.icon" class="size-4" />
          {{ item.label }}
        </RouterLink>
      </nav>
      <div class="border-t p-3">
        <button
          class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
          @click="logout"
        >
          <IconLogout class="size-4" />
          Logout
        </button>
      </div>
    </aside>
    <main class="flex-1 overflow-auto">
      <RouterView />
    </main>
  </div>
</template>
