<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { RouterView } from 'vue-router'
import { Toaster } from 'vue-sonner'
import 'vue-sonner/style.css'
import { useTheme } from '@/lib/theme'

const { init, theme } = useTheme()
const toasterTheme = computed(() => theme.value)

onMounted(init)
</script>

<template>
  <RouterView />
  <Toaster position="top-right" rich-colors close-button :theme="toasterTheme" />
</template>
