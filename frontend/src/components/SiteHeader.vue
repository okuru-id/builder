<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { IconSeparator } from '@tabler/icons-vue'
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

const crumbs = computed(() => {
  const segments = route.path.replace(/^\/admin\/?/, '').split('/').filter(Boolean)
  const items: { title: string; to?: string }[] = [{ title: 'Dashboard', to: '/' }]
  let acc = ''
  for (const seg of segments) {
    acc += '/' + seg
    const title = seg
      .split('-')
      .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
      .join(' ')
    items.push({ title, to: acc })
  }
  // last = current page (no link)
  items[items.length - 1].to = undefined
  return items
})

const title = computed(() => crumbs.value[crumbs.value.length - 1]?.title ?? 'okuru.id')
</script>

<template>
  <header class="flex h-12 shrink-0 items-center gap-2 border-b">
    <div class="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
      <SidebarTrigger class="-ml-1" />
      <Separator orientation="vertical" class="mx-2 data-[orientation=vertical]:h-4" />

      <Breadcrumb>
        <BreadcrumbList>
          <template v-for="(crumb, i) in crumbs" :key="crumb.title">
            <BreadcrumbItem>
              <BreadcrumbLink v-if="crumb.to" as-child>
                <RouterLink :to="crumb.to">{{ crumb.title }}</RouterLink>
              </BreadcrumbLink>
              <BreadcrumbPage v-else>{{ crumb.title }}</BreadcrumbPage>
            </BreadcrumbItem>
            <BreadcrumbSeparator v-if="i < crumbs.length - 1">
              <IconSeparator />
            </BreadcrumbSeparator>
          </template>
        </BreadcrumbList>
      </Breadcrumb>

      <h1 class="ml-auto text-base font-medium">{{ title }}</h1>
    </div>
  </header>
</template>