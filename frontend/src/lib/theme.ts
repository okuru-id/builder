import { ref } from 'vue'

const theme = ref<'light' | 'dark'>('dark')

function apply(t: 'light' | 'dark') {
  document.documentElement.classList.toggle('dark', t === 'dark')
  theme.value = t
  localStorage.setItem('theme', t)
}

function init() {
  const saved = localStorage.getItem('theme') as 'light' | 'dark' | null
  if (saved) {
    apply(saved)
    return
  }
  // system preference
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  apply(prefersDark ? 'dark' : 'light')
}

function toggle() {
  apply(theme.value === 'dark' ? 'light' : 'dark')
}

export function useTheme() {
  return { theme, init, toggle }
}
