// Type declarations for prismjs component subpath imports.
// @types/prismjs only covers the main 'prismjs' entry, not the per-language components.
declare module 'prismjs/components/prism-core' {
  export const highlight: (code: string, grammar: unknown, language: string) => string
  export const languages: Record<string, unknown>
}
declare module 'prismjs/components/prism-markup'
declare module 'prismjs/components/prism-clike'
declare module 'prismjs/components/prism-javascript'
declare module 'prismjs/components/prism-css'
