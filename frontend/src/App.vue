<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { RouterView } from 'vue-router'
import { Toaster } from '@/components/ui/sonner'
import { useTheme } from '@/lib/theme'

const { init, theme } = useTheme()
const toasterTheme = computed(() => theme.value)

onMounted(init)
</script>

<template>
  <RouterView />
  <Toaster position="bottom-center" :theme="toasterTheme" />
</template>
