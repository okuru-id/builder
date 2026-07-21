<script setup lang="ts">
import type { Component } from 'vue'
import { computed, ref, onMounted } from 'vue'
import {
  IconArticle,
  IconInbox,
  IconLayoutDashboard,
  IconPackage,
  IconFileSpark,
  IconUsers,
} from '@tabler/icons-vue'

import NavMain from '@/components/NavMain.vue'
import NavUser from '@/components/NavUser.vue'
import api from '@/lib/api'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

defineEmits<{ logout: [] }>()

const isSuper = ref(localStorage.getItem('is_super') === '1')

onMounted(async () => {
  isSuper.value = localStorage.getItem('is_super') === '1'
  try {
    const { data } = await api.get('/users')
    const sup = data?.is_super ? '1' : '0'
    localStorage.setItem('is_super', sup)
    isSuper.value = sup === '1'
  } catch {
    /* keep stored value */
  }
})

const data = {
  user: { name: 'Admin', email: 'admin@okuru.id', avatar: '' },
}

const navMain = computed(() => {
  const items: { title: string; to?: string; icon?: Component; items?: { title: string; to: string }[] }[] = [
    { title: 'Dashboard', to: '/', icon: IconLayoutDashboard },
    { title: 'Pages', to: '/pages', icon: IconFileSpark },
    {
      title: 'Blog',
      icon: IconArticle,
      items: [
        { title: 'All Posts', to: '/posts' },
        { title: 'New Post', to: '/posts/new' },
      ],
    },
    {
      title: 'Catalog',
      icon: IconPackage,
      items: [
        { title: 'Products', to: '/products' },
        { title: 'Projects', to: '/projects' },
        { title: 'Open Source', to: '/open-source' },
      ],
    },
    { title: 'Inbox', to: '/inbox', icon: IconInbox },
  ]
  if (isSuper.value) items.push({ title: 'Users', to: '/users', icon: IconUsers })
  return items
})
</script>

<template>
  <Sidebar variant="inset">
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <a href="#" style="gap:0.5rem">
              <img src="/images/logo.png" alt="okuru.id" class="aspect-square size-10 rounded-lg object-contain" />
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold text-lg">Builder</span>
                <span class="truncate text-xs">okuru.id</span>
              </div>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="navMain" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="data.user" @logout="$emit('logout')" />
    </SidebarFooter>
  </Sidebar>
</template>
