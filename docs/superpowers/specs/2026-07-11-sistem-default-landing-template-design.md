# Sistem Default Landing Template

## Goal

Use supplied Sistem HTML as active default storefront landing page and first admin gallery template.

## Data

Add nullable `html` payload to landing templates. `Sistem` stores supplied HTML unchanged, including `#` links.

Existing section-based templates remain compatible: their `html` is empty and use existing section behavior.

## Seeding

The landing-template seeder must be idempotent:

- create or update `Sistem` as the first gallery item;
- store supplied HTML exactly;
- set `landing_template_html` to same HTML;
- set `landing_mode=custom` so storefront immediately serves Sistem;
- never duplicate the template on repeated seeds.

## Applying templates

Use existing apply endpoint.

- Applying HTML template atomically writes `landing_template_html` and `landing_mode=custom`.
- It does not delete or replace current landing sections.
- Applying existing section template preserves existing replace-section flow.

## Admin gallery

- Templates with `html` show supplied design in sandboxed iframe preview, on card and preview Sheet.
- Iframe uses `sandbox="allow-scripts"`; no `allow-same-origin`.
- Existing templates retain image/fallback preview behavior.
- Existing apply confirmation and error behavior remain unchanged.

## Constraints

- Preserve all supplied HTML verbatim; do not change `#` links.
- Tailwind CDN or Google Fonts can fail without network; iframe must fail safely.
- No new dependency.

## Verification

- Backend test: repeated seed does not duplicate Sistem.
- Backend test: apply Sistem sets custom mode and supplied HTML while leaving sections unchanged.
- Run `cd backend && go test ./...`.
- Run `cd frontend && bun run build`.
- Manual: `/` renders Sistem after seeding; first gallery entry previews Sistem; applying it changes storefront without deleting existing sections.
