import { reactive } from 'vue'
import { DEFAULT_TOKENS } from '@/types/tokens'
import type { TokenConfig } from '@/types/tokens'

const STORAGE_KEY = 'okuru-builder-tokens'

const tokens = reactive<TokenConfig>({ ...DEFAULT_TOKENS })

function load() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) Object.assign(tokens, JSON.parse(raw))
  } catch { /* ignore corrupt data */ }
}

function save() {
  localStorage.setItem(STORAGE_KEY, JSON.stringify({ ...tokens }))
}

function reset() {
  Object.assign(tokens, DEFAULT_TOKENS)
  save()
}

export function useTokens() {
  return { tokens, load, save, reset }
}
