<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import api from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Skeleton } from '@/components/ui/skeleton'
import { IconPlus, IconTrash } from '@tabler/icons-vue'

const settings = ref<Record<string, string>>({})
const original = ref<Record<string, string>>({})
const loading = ref(true)
const saving = ref(false)
const newKey = ref('')
const newValue = ref('')

async function load() {
  loading.value = true
  try {
    const { data } = await api.get('/settings')
    settings.value = { ...(data.data ?? {}) }
    original.value = { ...settings.value }
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Failed to load settings')
  } finally {
    loading.value = false
  }
}

function addSetting() {
  if (!newKey.value.trim()) {
    toast.error('Key is required')
    return
  }
  settings.value[newKey.value.trim()] = newValue.value
  newKey.value = ''
  newValue.value = ''
}

async function removeKey(key: string) {
  if (!confirm(`Delete setting "${key}"?`)) return
  try {
    await api.delete(`/settings/${key}`)
    delete settings.value[key]
    delete original.value[key]
    toast.success('Deleted')
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Delete failed')
  }
}

async function save() {
  saving.value = true
  const changed = Object.entries(settings.value).filter(
    ([k, v]) => original.value[k] !== v,
  )
  if (changed.length === 0) {
    toast.info('No changes')
    saving.value = false
    return
  }
  try {
    for (const [key, value] of changed) {
      await api.put('/settings', { key, value })
    }
    original.value = { ...settings.value }
    toast.success(`Saved ${changed.length} setting(s)`)
  } catch (e: any) {
    toast.error(e.response?.data?.error || 'Save failed')
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="p-8">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="font-heading text-2xl font-bold">Settings</h1>
      <Button :disabled="saving" @click="save">{{ saving ? 'Saving…' : 'Save Changes' }}</Button>
    </div>

    <div v-if="loading" class="flex flex-col gap-3">
      <Skeleton v-for="i in 4" :key="i" class="h-10 w-full" />
    </div>
    <div v-else class="flex max-w-2xl flex-col gap-3">
      <div
        v-for="(value, key) in settings"
        :key="key"
        class="flex items-end gap-2 rounded-lg border p-3"
        :class="{ 'border-primary': original[key] !== value }"
      >
        <div class="flex flex-1 flex-col gap-1">
          <Label class="text-xs text-muted-foreground">{{ key }}</Label>
          <Input v-model="settings[key]" />
        </div>
        <Button variant="ghost" size="icon-sm" @click="removeKey(key as string)">
          <IconTrash class="size-4 text-destructive" />
        </Button>
      </div>

      <div class="mt-4 flex items-end gap-2 rounded-lg border border-dashed p-3">
        <div class="flex flex-1 flex-col gap-1">
          <Label class="text-xs text-muted-foreground">New key</Label>
          <Input v-model="newKey" placeholder="site.title" />
        </div>
        <div class="flex flex-1 flex-col gap-1">
          <Label class="text-xs text-muted-foreground">Value</Label>
          <Input v-model="newValue" />
        </div>
        <Button variant="outline" size="sm" @click="addSetting">
          <IconPlus class="size-4" /> Add
        </Button>
      </div>
    </div>
  </div>
</template>
