<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Checkbox } from '@/components/ui/checkbox'
import { Skeleton } from '@/components/ui/skeleton'
import { IconArrowLeft } from '@tabler/icons-vue'

const route = useRoute()
const router = useRouter()
const id = computed(() => (route.params.id as string) || null)
const isEdit = computed(() => !!id.value)
const isSelf = computed(() => isEdit.value && Number(id.value) === Number(localStorage.getItem('user_id')))
const loading = ref(false)
const saving = ref(false)
const isSuper = ref(false)

const form = ref({
  email: '',
  name: '',
  password: '',
  is_active: true,
  is_admin: false,
  reset_totp: false,
})

async function load() {
  // Non-super users reach this page only for their own profile (guarded by the list).
  const { data } = await api.get('/users')
  isSuper.value = data.is_super === true

  if (!isEdit.value) return
  loading.value = true
  try {
    const { data } = await api.get(`/users/${id.value}`)
    form.value = {
      email: data.data.email,
      name: data.data.name,
      password: '',
      is_active: data.data.is_active,
      is_admin: data.data.is_admin === true,
      reset_totp: false,
    }
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load')
    router.push('/users')
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!form.value.email) {
    toast.warning('Email is required')
    return
  }
  if (!isEdit.value && !form.value.password) {
    toast.warning('Password is required for new users')
    return
  }
  saving.value = true
  const payload: any = { email: form.value.email, name: form.value.name, is_active: form.value.is_active }
  if (!isSelf.value) payload.is_admin = form.value.is_admin
  if (form.value.password) payload.password = form.value.password
  try {
    if (isEdit.value) {
      await api.put(`/users/${id.value}${form.value.reset_totp ? '?reset_totp=1' : ''}`, payload)
    } else {
      await api.post('/users', payload)
    }
    toast.success('Saved')
    router.push('/users')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Save failed')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="p-4 lg:p-6">
    <Button variant="ghost" size="sm" class="mb-4 -ml-2" @click="router.push('/users')">
      <IconArrowLeft class="size-4" />
      Back
    </Button>

    <h1 class="font-heading mb-6 text-2xl font-bold">{{ isEdit ? 'Edit User' : 'New User' }}</h1>

    <div v-if="loading" class="max-w-xl space-y-4">
      <Skeleton class="h-10 w-full" />
      <Skeleton class="h-10 w-full" />
      <Skeleton class="h-10 w-full" />
    </div>

    <form v-else class="grid max-w-xl grid-cols-[max-content_1fr] items-center gap-x-3 gap-y-2" @submit.prevent="save">
      <Label class="text-right whitespace-nowrap">Email *</Label>
      <Input v-model="form.email" type="email" />
      <Label class="text-right whitespace-nowrap">Name</Label>
      <Input v-model="form.name" />
      <Label class="text-right whitespace-nowrap">{{ isEdit ? 'New Password (leave blank to keep)' : 'Password *' }}</Label>
      <Input v-model="form.password" type="password" autocomplete="new-password" />
      <Label class="text-right whitespace-nowrap">Active</Label>
      <div class="flex items-center gap-3">
        <Switch :model-value="form.is_active" :disabled="!isSuper" @update:model-value="(v: boolean) => (form.is_active = v)" />
      </div>
      <template v-if="isSuper && !isSelf">
        <Label class="text-right whitespace-nowrap">Admin</Label>
        <div class="flex items-center gap-3">
          <Switch :model-value="form.is_admin" @update:model-value="(v: boolean) => (form.is_admin = v)" />
          <span class="text-sm text-muted-foreground">grant full access</span>
        </div>
      </template>
      <template v-if="isSuper && isEdit">
        <Label class="text-right whitespace-nowrap">Reset 2FA</Label>
        <div class="flex items-center gap-2">
          <Checkbox :model-value="form.reset_totp" @update:model-value="(v: boolean | 'indeterminate') => (form.reset_totp = v === true)" id="reset-totp" />
          <Label for="reset-totp" class="text-sm font-normal text-muted-foreground">on save</Label>
        </div>
      </template>
      <div></div>
      <div class="flex gap-2">
        <Button type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</Button>
        <Button type="button" variant="ghost" @click="router.push('/users')">Cancel</Button>
      </div>
    </form>
  </div>
</template>
