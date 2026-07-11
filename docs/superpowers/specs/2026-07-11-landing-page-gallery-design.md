# Landing Page Template Gallery

## Goal

Make `/admin/landing-page` a template gallery. Keep section-content editing on a separate route.

## Routes

- `/admin/landing-page`: template gallery, replacing current builder/custom UI.
- `/admin/landing-page/content`: existing section editor, moved from current page.

Gallery header includes an `Edit Content` link to `/admin/landing-page/content`.

## Gallery

- Load templates with existing `GET /landing-templates` API.
- Show responsive template cards.
- Card shows `preview` image when supplied. If absent or unusable, show a safe fallback containing section names and count.
- Card click opens a `Sheet` preview panel. Panel shows name, description, section list, preview/fallback, Cancel, and Apply Template actions.
- Applying requires native `confirm()` because it replaces all landing sections.
- Confirmed apply calls existing `POST /landing-templates/:id/apply`.
- On success: show toast and close dialog. On failure: show toast and keep dialog open.
- Loading uses skeleton cards. Empty results show an empty state.

## Content editor

Move builder section controls into a dedicated `LandingPageContent.vue` view:

- list sections;
- add, edit, delete, reorder, and toggle sections;
- retain generic and known section field editing.

Remove custom HTML mode and default/minimal/bold template switches from admin UI. Existing settings and APIs remain unchanged.

## Constraints

- Frontend-only change. No migration or backend API change.
- Reuse installed shadcn-vue sheet primitives for preview and existing API client; do not add dependencies.
- Preserve authenticated route behavior and lazy-loaded views.

## Error handling

- Load and mutation API errors use existing toast pattern.
- Preview sheet remains open on apply error.
- Broken/missing preview must not break card or modal; fallback stays visible.

## Verification

Run `bun run build` in `frontend/`.

Manual checks:

1. Gallery loads templates, skeletons, and empty state correctly.
2. Card opens preview sheet; Cancel does not mutate sections.
3. Apply confirmation, success, and error states behave correctly.
4. `Edit Content` opens section editor; all old section CRUD/reorder/toggle actions work.
