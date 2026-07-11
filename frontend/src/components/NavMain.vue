<script setup lang="ts">
import type { Component } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { IconChevronRight } from '@tabler/icons-vue'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from '@/components/ui/sidebar'

interface SubItem {
  title: string
  to: string
}

interface NavItem {
  title: string
  to?: string
  icon?: Component
  items?: SubItem[]
  isActive?: boolean
}

const props = defineProps<{
  items: NavItem[]
}>()

const route = useRoute()
</script>

<template>
  <SidebarGroup>
    <SidebarGroupLabel>Platform</SidebarGroupLabel>
    <SidebarMenu>
      <template v-for="item in props.items" :key="item.title">
        <!-- Collapsible group: whole button toggles -->
        <Collapsible
          v-if="item.items?.length"
          as-child
          :default-open="true"
        >
          <SidebarMenuItem>
            <CollapsibleTrigger as-child>
              <SidebarMenuButton :tooltip="item.title">
                <component :is="item.icon" v-if="item.icon" />
                <span>{{ item.title }}</span>
                <IconChevronRight class="ml-auto transition-transform duration-200 group-data-[state=open]/menu-button:rotate-90" />
              </SidebarMenuButton>
            </CollapsibleTrigger>
            <CollapsibleContent>
              <SidebarMenuSub>
                <SidebarMenuSubItem v-for="subItem in item.items" :key="subItem.title">
                  <SidebarMenuSubButton as-child
                    :is-active="route.path === subItem.to ? true : undefined"
                  >
                    <RouterLink :to="subItem.to">
                      <span>{{ subItem.title }}</span>
                    </RouterLink>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              </SidebarMenuSub>
            </CollapsibleContent>
          </SidebarMenuItem>
        </Collapsible>

        <!-- Flat link -->
        <SidebarMenuItem v-else>
          <SidebarMenuButton as-child :tooltip="item.title"
            :data-active="item.to && route.path === item.to ? true : undefined"
          >
            <RouterLink :to="item.to!">
              <component :is="item.icon" v-if="item.icon" />
              <span>{{ item.title }}</span>
            </RouterLink>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </template>
    </SidebarMenu>
  </SidebarGroup>
</template>
