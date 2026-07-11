<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'
import AppSidebar from '@/components/AppSidebar.vue'
import SiteHeader from '@/components/SiteHeader.vue'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'

const router = useRouter()

function logout() {
  localStorage.removeItem('access_token')
  router.push('/login')
}
</script>

<template>
  <SidebarProvider>
    <AppSidebar @logout="logout" />
    <SidebarInset class="min-w-0">
      <SiteHeader />
      <div class="flex flex-1 flex-col gap-4 p-4 min-w-0">
        <RouterView />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
