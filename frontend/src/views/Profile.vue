<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { Badge } from '@/components/ui/badge'
import { IconArrowLeft, IconEdit } from '@tabler/icons-vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(true)
const saving = ref(false)
const editing = ref(false)
const isSuper = ref(false)

const profile = ref({ name: '', email: '' })
const form = ref({ name: '', email: '', password: '' })

const userId = computed(() => Number(localStorage.getItem('user_id') || 0))

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/users')
    isSuper.value = !!data.is_super
    const me = Array.isArray(data.data) ? data.data[0] : data.data
    profile.value = {
      name: me?.name || me?.email?.split('@')[0] || '',
      email: me?.email ?? '',
    }
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load')
    router.push('/')
  } finally {
    loading.value = false
  }
}

function startEdit() {
  form.value = { name: profile.value.name, email: profile.value.email, password: '' }
  editing.value = true
}

async function save() {
  if (!form.value.email) {
    toast.warning('Email is required')
    return
  }
  saving.value = true
  const payload: any = { name: form.value.name, email: form.value.email }
  if (form.value.password) payload.password = form.value.password
  try {
    await api.put(`/users/${userId.value}`, payload)
    profile.value = { name: form.value.name, email: form.value.email }
    toast.success('Profile saved')
    editing.value = false
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
    <Button variant="ghost" size="sm" class="mb-4 -ml-2" @click="router.push('/')">
      <IconArrowLeft class="size-4" />
      Back
    </Button>

    <div class="mb-6 flex items-center gap-3">
      <h1 class="font-heading text-2xl font-bold">My Profile</h1>
      <Badge v-if="isSuper" variant="default">Super Admin</Badge>
    </div>

    <div v-if="loading" class="max-w-xl space-y-4">
      <Skeleton v-for="i in 3" :key="i" class="h-10 w-full" />
    </div>

    <!-- read-only view -->
    <div v-else-if="!editing" class="max-w-xl">
      <dl class="divide-y rounded-lg border">
        <div class="grid grid-cols-[140px_1fr] gap-2 p-4">
          <dt class="text-sm text-muted-foreground">Name</dt>
          <dd class="text-sm font-medium">{{ profile.name || '—' }}</dd>
        </div>
        <div class="grid grid-cols-[140px_1fr] gap-2 p-4">
          <dt class="text-sm text-muted-foreground">Email</dt>
          <dd class="text-sm font-medium">{{ profile.email }}</dd>
        </div>
      </dl>
      <Button class="mt-4" @click="startEdit">
        <IconEdit class="size-4" />
        Edit Profile
      </Button>
    </div>

    <!-- edit form -->
    <form v-else class="grid max-w-xl grid-cols-[max-content_1fr] items-center gap-x-3 gap-y-2" @submit.prevent="save">
      <Label class="text-right whitespace-nowrap">Name</Label>
      <Input v-model="form.name" />

      <Label class="text-right whitespace-nowrap">Email *</Label>
      <Input v-model="form.email" type="email" />

      <Label class="text-right whitespace-nowrap">New Password</Label>
      <Input v-model="form.password" type="password" autocomplete="new-password" placeholder="leave blank to keep" />

      <div></div>
      <div class="flex gap-2 pt-2">
        <Button type="submit" :disabled="saving">{{ saving ? 'Saving…' : 'Save' }}</Button>
        <Button type="button" variant="ghost" @click="editing = false">Cancel</Button>
      </div>
    </form>
  </div>
</template>
