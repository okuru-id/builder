// Loads reusable component masters and keeps them in a Map keyed by id.
// Shared across the builder via provide/inject (BUILDER_KEY exposes this).
import { ref } from 'vue'
import api from '@/lib/api'
import type { Component, Node } from '@/types/page-builder'

export function useComponents() {
  const components = ref<Component[]>([])
  const loading = ref(false)

  async function load() {
    loading.value = true
    try {
      const res = await api.get<{ data: Component[] }>('/landing-components')
      components.value = res.data.data ?? []
    } catch (e) {
      console.error('failed to load components', e)
    } finally {
      loading.value = false
    }
  }

  // Resolve a componentId to its master root node, or null if not found.
  function masterRoot(id: number | undefined): Node | null {
    if (!id) return null
    const c = components.value.find((c) => c.id === id)
    return c?.tree.root ?? null
  }

  async function create(name: string, tree: { root: Node }) {
    const res = await api.post<{ data: Component }>('/landing-components', { name, tree })
    components.value.unshift(res.data.data)
    return res.data.data
  }

  async function remove(id: number) {
    await api.delete(`/landing-components/${id}`)
    components.value = components.value.filter((c) => c.id !== id)
  }

  async function rename(id: number, name: string) {
    const res = await api.put<{ data: Component }>(`/landing-components/${id}`, { name })
    const updated = res.data.data
    const i = components.value.findIndex((c) => c.id === id)
    if (i >= 0) components.value[i] = updated
    return updated
  }

  async function duplicate(id: number) {
    const res = await api.post<{ data: Component }>(`/landing-components/${id}/duplicate`)
    components.value.unshift(res.data.data)
    return res.data.data
  }

  return { components, loading, load, masterRoot, create, remove, rename, duplicate }
}
