// Injection key + helper for accessing the builder store from deeply nested renderers.
import type { InjectionKey } from 'vue'
import type { BuilderStore } from '@/components/builder/useBuilderStore'

export const BUILDER_KEY: InjectionKey<BuilderStore> = Symbol('builder-store')
