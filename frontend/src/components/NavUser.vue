<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  IconCircleCheck,
  IconCreditCard,
  IconBell,
  IconLogout,
  IconSelector,
  IconUserCircle,
  IconSun,
  IconMoon,
} from '@tabler/icons-vue'
import api from '@/lib/api'
import { useTheme } from '@/lib/theme'
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'

defineProps<{
  user?: {
    name: string
    email: string
    avatar: string
  }
}>()

const emit = defineEmits<{ logout: [] }>()
const { isMobile } = useSidebar()
const { theme, toggle } = useTheme()

const me = ref<{ name?: string; email?: string; avatar?: string }>({})

onMounted(async () => {
  try {
    const { data } = await api.get('/auth/me')
    me.value = data?.data ?? data ?? {}
  } catch {
    // fall back to props
  }
})

const name = computed(() => me.value.name || me.value.email?.split('@')[0] || 'Admin')
const email = computed(() => me.value.email || 'admin@okuru.id')
const avatar = computed(() => me.value.avatar || '')

function initials(s: string) {
  return s
    .split(/[\s@._-]+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((w) => w[0]?.toUpperCase())
    .join('')
}
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <Avatar class="size-8 rounded-lg">
              <AvatarImage :src="avatar" :alt="name" />
              <AvatarFallback class="rounded-lg">{{ initials(name) }}</AvatarFallback>
            </Avatar>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-medium">{{ name }}</span>
              <span class="truncate text-xs">{{ email }}</span>
            </div>
            <IconSelector class="ml-auto size-4" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          class="w-(--reka-dropdown-menu-trigger-width) min-w-56 rounded-lg"
          :side="isMobile ? 'bottom' : 'right'"
          align="end"
          :side-offset="4"
        >
          <DropdownMenuLabel class="p-0 font-normal">
            <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
              <Avatar class="size-8 rounded-lg">
                <AvatarImage :src="avatar" :alt="name" />
                <AvatarFallback class="rounded-lg">{{ initials(name) }}</AvatarFallback>
              </Avatar>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-medium">{{ name }}</span>
                <span class="truncate text-xs">{{ email }}</span>
              </div>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem>
              <IconUserCircle />
              Account
            </DropdownMenuItem>
            <DropdownMenuItem>
              <IconCircleCheck />
              Profile
            </DropdownMenuItem>
            <DropdownMenuItem>
              <IconCreditCard />
              Billing
            </DropdownMenuItem>
            <DropdownMenuItem>
              <IconBell />
              Notifications
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuItem @select="toggle">
            <component :is="theme === 'dark' ? IconSun : IconMoon" />
            {{ theme === 'dark' ? 'Light mode' : 'Dark mode' }}
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem @select="emit('logout')">
            <IconLogout />
            Log out
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  </SidebarMenu>
</template>
