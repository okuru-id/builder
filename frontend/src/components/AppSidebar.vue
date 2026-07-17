<script setup lang="ts">
import type { Component } from 'vue'
import {
  IconArticle,
  IconInbox,
  IconLayoutDashboard,
  IconLifebuoy,
  IconPackage,
  IconFileSpark,
} from '@tabler/icons-vue'

import NavMain from '@/components/NavMain.vue'
import NavSecondary from '@/components/NavSecondary.vue'
import NavUser from '@/components/NavUser.vue'
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

const data = {
  user: { name: 'Admin', email: 'admin@okuru.id', avatar: '' },
  navMain: [
    { title: 'Dashboard', to: '/', icon: IconLayoutDashboard },
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
    { title: 'Pages', to: '/pages', icon: IconFileSpark },
  ] as { title: string; to?: string; icon?: Component; items?: { title: string; to: string }[] }[],
  navSecondary: [
    { title: 'Support', to: '/', icon: IconLifebuoy },
  ],
}
</script>

<template>
  <Sidebar variant="inset">
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <a href="#">
              <img src="/images/logo.png" alt="okuru.id" class="aspect-square size-20 rounded-lg object-contain" />
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold text-lg">okuru.id</span>
              </div>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="data.navMain" />
      <NavSecondary :items="data.navSecondary" class="mt-auto" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="data.user" @logout="$emit('logout')" />
    </SidebarFooter>
  </Sidebar>
</template>
