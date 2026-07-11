<script setup lang="ts">
import type { Component } from 'vue'
import {
  IconArticle,
  IconInbox,
  IconInnerShadowTop,
  IconLayoutDashboard,
  IconLifebuoy,
  IconPackage,
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
              <div class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                <IconInnerShadowTop class="!size-5" />
              </div>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">okuru.id</span>
                <span class="truncate text-xs">Admin Panel</span>
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
