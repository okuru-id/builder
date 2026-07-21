<script setup lang="ts">
import type { ColumnDef } from '@tanstack/vue-table'
import { h, ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import DataTable from '@/components/DataTable.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { IconPlus, IconEdit, IconTrash, IconShieldX } from '@tabler/icons-vue'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'

const router = useRouter()

interface UserRow {
  id: number
  email: string
  name: string
  is_active: boolean
  is_admin: boolean
  totp_enabled: boolean
}

const users = ref<UserRow[]>([])
const loading = ref(true)
const isSuper = ref(false)
const selfId = Number(localStorage.getItem('user_id') || 0)

const showDelete = ref(false)
const deletingId = ref(0)
const showResetTotp = ref(false)
const resetTotpId = ref(0)

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/users')
    users.value = data.data ?? []
    isSuper.value = data.is_super === true
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load')
  } finally {
    loading.value = false
  }
}

function confirmRemove(id: number) {
  if (!isSuper.value) return
  deletingId.value = id
  showDelete.value = true
}
async function doDelete() {
  try {
    await api.delete(`/users/${deletingId.value}`)
    users.value = users.value.filter((u) => u.id !== deletingId.value)
    toast.success('Deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
  } finally {
    showDelete.value = false
  }
}

function confirmResetTotp(id: number) {
  if (!isSuper.value) return
  resetTotpId.value = id
  showResetTotp.value = true
}
async function doResetTotp() {
  try {
    await api.put(`/users/${resetTotpId.value}?reset_totp=1`, {})
    toast.success('2FA reset')
    await load()
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Reset failed')
  } finally {
    showResetTotp.value = false
  }
}

const columns = computed<ColumnDef<UserRow>[]>(() => [
  { accessorKey: 'name', header: 'Name', cell: ({ row }) => h('span', { class: 'font-medium' }, row.getValue('name') || '—') },
  { accessorKey: 'email', header: 'Email' },
  {
    accessorKey: 'is_admin',
    header: 'Role',
    cell: ({ row }) => row.getValue('is_admin')
      ? h(Badge, { variant: 'default' }, () => 'admin')
      : h('span', { class: 'text-muted-foreground' }, 'user'),
  },
  {
    accessorKey: 'is_active',
    header: 'Status',
    cell: ({ row }) => h(Badge, { variant: row.getValue('is_active') ? 'default' : 'secondary' }, () => (row.getValue('is_active') ? 'active' : 'disabled')),
  },
  {
    accessorKey: 'totp_enabled',
    header: '2FA',
    cell: ({ row }) => h(Badge, { variant: row.getValue('totp_enabled') ? 'default' : 'outline' }, () => (row.getValue('totp_enabled') ? 'on' : 'off')),
  },
  {
    id: 'actions',
    header: () => h('div', { class: 'text-right' }, 'Actions'),
    enableSorting: false,
    cell: ({ row }) =>
      h('div', { class: 'flex justify-end gap-1' }, [
        isSuper.value
          ? h(Button, { variant: 'ghost', size: 'icon-sm', title: 'Reset 2FA', onClick: () => confirmResetTotp(row.original.id) }, () => h(IconShieldX, { class: 'size-4' }))
          : null,
        h(Button, { variant: 'ghost', size: 'icon-sm', onClick: () => router.push(`/users/${row.original.id}/edit`) }, () => h(IconEdit, { class: 'size-4' })),
        isSuper.value && row.original.id !== selfId
          ? h(Button, { variant: 'ghost', size: 'icon-sm', onClick: () => confirmRemove(row.original.id) }, () => h(IconTrash, { class: 'size-4 text-destructive' }))
          : null,
      ]),
  },
])

onMounted(load)
</script>

<template>
  <div class="p-4 lg:p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="font-heading text-2xl font-bold">Users</h1>
      <Button v-if="isSuper" @click="router.push('/users/new')">
        <IconPlus class="size-4" />
        New User
      </Button>
    </div>

    <DataTable :data="users" :columns="columns" :loading="loading" :draggable="false" search-key="email" search-placeholder="Search users…" />
  </div>

  <ConfirmDialog v-model:open="showDelete" title="Delete user" description="Are you sure you want to delete this user?" @confirm="doDelete" />
  <ConfirmDialog v-model:open="showResetTotp" title="Reset 2FA" description="Disable two-factor auth for this user?" @confirm="doResetTotp" />
</template>
