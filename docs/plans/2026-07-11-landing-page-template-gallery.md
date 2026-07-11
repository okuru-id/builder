# Landing Page Template Gallery Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Make `/admin/landing-page` a template gallery with preview-before-apply and move section editing to `/admin/landing-page/content`.

**Architecture:** Replace `LandingPage.vue` with a gallery-only view using existing template APIs and installed Sheet primitives. Extract existing builder section CRUD UI and state into `LandingPageContent.vue`; add route and sidebar navigation without backend changes.

**Tech Stack:** Vue 3 `<script setup>`, TypeScript, Vue Router, axios API client, shadcn-vue/Reka Sheet, Tailwind v4, vue-sonner.

---

### Task 1: Add content-editor route and navigation

**Files:**
- Modify: `frontend/src/router/index.ts`
- Modify: `frontend/src/components/AppSidebar.vue`

**Step 1: Write failing route-level verification**

No frontend test runner exists. Use TypeScript build as regression check after route is added.

**Step 2: Add lazy-loaded content route**

In `frontend/src/router/index.ts`, retain existing gallery route and add child route directly after it:

```ts
{ path: 'landing-page/content', name: 'landing-page-content', component: () => import('@/views/LandingPageContent.vue') },
```

Do not change auth guard or base route.

**Step 3: Add sidebar content link**

In `frontend/src/components/AppSidebar.vue`, make Landing Page a grouped `navSecondary` entry, matching Blog/Catalog shape:

```ts
{
  title: 'Landing Page',
  icon: IconLayout,
  items: [
    { title: 'Templates', to: '/landing-page' },
    { title: 'Content', to: '/landing-page/content' },
  ],
},
```

Update local `navSecondary` type so `items?: { title: string; to: string }[]` is accepted. Do not add a second icon or dependency.

**Step 4: Run typecheck to verify expected temporary failure**

Run: `cd frontend && bun run build`

Expected: FAIL because `@/views/LandingPageContent.vue` does not exist.

**Step 5: Commit route scaffold**

```bash
git add frontend/src/router/index.ts frontend/src/components/AppSidebar.vue
git commit -m "feat: add landing content route"
```

### Task 2: Extract builder section editor into content view

**Files:**
- Create: `frontend/src/views/LandingPageContent.vue`
- Modify: `frontend/src/views/LandingPage.vue`

**Step 1: Create content view from existing builder code**

Copy only builder-related imports, interfaces, reactive state, helpers, API calls, lifecycle hook, template blocks, and scoped styles from `LandingPage.vue` into `LandingPageContent.vue`.

Keep these capabilities unchanged:

- `GET /landing-sections` loading;
- section create/update/delete/toggle/sort APIs;
- draft, generic field, hero/CTA, and known-list item editing;
- drag reorder behavior;
- loading and error toasts.

Exclude all custom HTML code and settings logic: `PrismEditor`, Prism imports/styles, `mode`, `customHtml`, `customTab`, `previewUrl`, `updatePreview`, `saveCustomHtml`, `setMode`, `saveTemplate`, template style toggle, and `GET /settings`.

Header must be:

```vue
<div class="flex flex-wrap items-center justify-between gap-3">
  <h1 class="font-heading text-2xl font-bold">Landing Content</h1>
  <Button @click="addingSection = true" v-if="!addingSection">
    <IconPlus class="size-4" /> Add Section
  </Button>
</div>
```

Keep Add Section form and all section cards from existing builder UI. Keep a local `dragIdx` for section/item drag behavior; do not share global state with gallery.

**Step 2: Reduce `LandingPage.vue` to gallery-only later**

Do not delete current builder code until gallery replacement exists in Task 3, avoiding a temporary blank route.

**Step 3: Run build**

Run: `cd frontend && bun run build`

Expected: PASS. Route target exists; no custom editor imports are required by new content view.

**Step 4: Commit extracted editor**

```bash
git add frontend/src/views/LandingPageContent.vue
git commit -m "feat: move landing section editor to content page"
```

### Task 3: Replace landing page with template gallery and preview sheet

**Files:**
- Modify: `frontend/src/views/LandingPage.vue`

**Step 1: Replace old view with gallery-only state and API logic**

Delete builder/custom editor state and imports. Keep only `ref`, `onMounted`, `toast`, `api`, `Button`, `Badge`, `Card`, `CardHeader`, `CardContent`, `Skeleton`, `IconLayoutGrid`, plus `IconArrowRight` or another already-installed Tabler icon for content navigation.

Define:

```ts
interface LandingTemplate {
  id: number
  name: string
  description: string
  preview: string
  sections: { type: string }[]
}

const templates = ref<LandingTemplate[]>([])
const loading = ref(true)
const selectedTemplate = ref<LandingTemplate | null>(null)
const sheetOpen = ref(false)
const applyingTemplate = ref<number | null>(null)
```

Implement `loadTemplates()` with `api.get('/landing-templates')`, mapping absent `data.data` to `[]`, toast on failure, and `finally` setting `loading` false. Call it from `onMounted`.

Implement `openPreview(tmpl)` to set selection then `sheetOpen.value = true`.

Implement `applySelectedTemplate()`:

```ts
const tmpl = selectedTemplate.value
if (!tmpl || !confirm(`Apply "${tmpl.name}"? This will replace all current sections.`)) return
applyingTemplate.value = tmpl.id
try {
  await api.post(`/landing-templates/${tmpl.id}/apply`)
  toast.success('Template applied')
  sheetOpen.value = false
  selectedTemplate.value = null
} catch (e: any) {
  toast.error(e.response?.data?.error || 'Failed to apply template')
} finally {
  applyingTemplate.value = null
}
```

Do not reload templates after apply; templates are immutable for this flow. Error must leave sheet and selection intact.

**Step 2: Add safe preview fallback**

Use a reusable template block or small local component-free markup for both card and sheet. Render `<img>` only while URL is usable; track failed image IDs in `ref(new Set<number>())` and handle image failure:

```ts
const brokenPreviews = ref(new Set<number>())
function previewAvailable(tmpl: LandingTemplate) {
  return Boolean(tmpl.preview) && !brokenPreviews.value.has(tmpl.id)
}
function markPreviewBroken(id: number) {
  brokenPreviews.value = new Set([...brokenPreviews.value, id])
}
```

Fallback must always include section count and section labels. Use known labels (`Hero`, `Client Logos`, `Services`, `Open Source Projects`, `CTA`), defaulting to raw type.

**Step 3: Build gallery UI**

Header:

```vue
<div class="flex flex-wrap items-center justify-between gap-3">
  <h1 class="font-heading text-2xl font-bold">Landing Page Templates</h1>
  <Button as-child variant="outline">
    <RouterLink to="/landing-page/content">Edit Content</RouterLink>
  </Button>
</div>
```

Import `RouterLink` from `vue-router`.

- While loading, show three skeleton cards in existing responsive grid.
- On empty result, show centered `No templates yet.` state.
- Otherwise use `grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4` cards.
- Cards expose keyboard interaction: use a native `<button type="button">` card wrapper or a button inside Card, not a click-only `div`.
- Card shows preview/fallback, name, section count, description, and up to five section labels.
- Clicking card opens Sheet; no direct apply action exists on card.

**Step 4: Add accessible Sheet preview**

Import all installed primitives from `@/components/ui/sheet`:

```ts
import { Sheet, SheetClose, SheetContent, SheetDescription, SheetFooter, SheetHeader, SheetTitle } from '@/components/ui/sheet'
```

Bind Sheet to `sheetOpen`. Render selected content only when `selectedTemplate` exists. Sheet includes title, description, large preview/fallback, complete section labels, Cancel (`SheetClose as-child` with a Button), and Apply Template (`Button` calling `applySelectedTemplate`). Disable Apply while `applyingTemplate` is non-null and change label to `Applying...`.

Close handler must clear selection after close animation only when no apply is in progress; simplest valid approach uses `@update:open="(open) => { sheetOpen = open; if (!open && !applyingTemplate) selectedTemplate = null }"`.

**Step 5: Remove dead legacy UI**

After gallery works, delete all legacy builder/custom template imports, API logic, template markup, Prism styles, and related fields from `LandingPage.vue`. Those builder capabilities now exist only in `LandingPageContent.vue`. Do not touch backend settings or APIs.

**Step 6: Run build**

Run: `cd frontend && bun run build`

Expected: PASS with `vue-tsc` followed by Vite build.

**Step 7: Commit gallery**

```bash
git add frontend/src/views/LandingPage.vue
git commit -m "feat: make landing page a template gallery"
```

### Task 4: Verify complete user flow

**Files:**
- Modify: none unless verification finds defect

**Step 1: Run production build**

Run: `cd frontend && bun run build`

Expected: successful Vue typecheck and Vite bundle.

**Step 2: Manual authenticated browser checks**

1. Open `/admin/landing-page`; verify skeleton then cards or empty state.
2. Check template with preview; then trigger broken preview URL in available test data or DevTools and confirm fallback stays visible.
3. Open card through mouse and keyboard; Sheet shows full template data.
4. Cancel Sheet; verify no API mutation occurred.
5. Click Apply, reject native confirmation; verify no request.
6. Click Apply, confirm; verify one `POST /landing-templates/:id/apply`, toast, and closed Sheet.
7. Force apply API error; verify toast and Sheet remains open.
8. Open `/admin/landing-page/content`; verify add/edit/delete/toggle/reorder section behavior matches previous page.

**Step 3: Commit only repair changes if any**

```bash
git add frontend
git commit -m "fix: complete landing gallery flow"
```

Skip this commit when no verification repairs were needed.
