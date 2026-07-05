# Landing Page Builder — Design

**Date:** 2026-07-05
**Status:** Approved

## Overview

Replace the existing key-value settings editor at `/admin/landing-page` with a
dynamic section-based landing page builder. Store each landing page section as a
row in a single `landing_sections` table with typed JSON content.

## Model

One table — `landing_sections`:

| Field       | Type         | Notes                        |
|-------------|--------------|------------------------------|
| id          | auto-increment |                              |
| type        | string       | unique — `hero`, `clients`, `services`, `projects`, `cta` |
| content     | jsonb        | structured per type (see below) |
| sort_order  | int          | section ordering             |
| is_active   | bool         | show/hide section            |
| timestamps  |              |                              |

### JSON Content Schemas

**hero** (singleton):
```json
{
  "greeting_en": "string",
  "greeting_id": "string",
  "description_en": "string",
  "description_id": "string",
  "profile_image": "string",
  "profile_image_mobile": "string",
  "cta_text": "string",
  "cta_link": "string"
}
```

**clients** (array):
```json
[
  { "name": "string", "logo": "string" }
]
```

**services** (array):
```json
[
  {
    "title": "string",
    "description_en": "string",
    "description_id": "string",
    "icon": "string"
  }
]
```

**projects** (array):
```json
[
  {
    "title_en": "string",
    "title_id": "string",
    "description_en": "string",
    "description_id": "string",
    "github_url": "string",
    "technologies": ["string"]
  }
]
```

**cta** (singleton):
```json
{
  "heading": "string",
  "email": "string",
  "whatsapp": "string"
}
```

## API

### Public (no auth)

| Method | Path             | Description                     |
|--------|------------------|---------------------------------|
| GET    | `/api/v1/landing` | All active sections, ordered   |

Response:
```json
{
  "hero": { ... },
  "clients": [ ... ],
  "services": [ ... ],
  "projects": [ ... ],
  "cta": { ... }
}
```

### Admin (auth + totp)

| Method | Path                                     | Description              |
|--------|------------------------------------------|--------------------------|
| GET    | `/admin/api/landing-sections`            | List all sections        |
| POST   | `/admin/api/landing-sections`            | Create section           |
| PUT    | `/admin/api/landing-sections/{id}`       | Update section content   |
| DELETE | `/admin/api/landing-sections/{id}`       | Delete section           |
| PATCH  | `/admin/api/landing-sections/{id}/toggle` | Toggle is_active        |
| PATCH  | `/admin/api/landing-sections/{id}/sort`  | Update sort_order        |

## Admin UI

Page at `/admin/landing-page` (replacing current view):

- Each section rendered as a card with type name, summary preview, and action buttons (edit, toggle, delete)
- List-type sections (clients, services, projects) show an inline table of items
- "Edit" opens a modal/drawer with form fields tailored to the section type
- "Add Item" for list sections inserts a new item into the content array
- Save per section via `PUT`
- Drag-to-reorder sections (or up/down buttons)

## Landing Page

`frontend/src/landing/App.vue`:

- On mount, fetch `GET /api/v1/landing`
- Map section `type` to render component: `<LandingHero>`, `<LandingClients>`, etc.
- Fall back to hardcoded data if API fails or section is missing

## Backend Changes

1. Migration: create `landing_sections` table
2. Model: `LandingSection`
3. Controller: `app/http/controllers/landing_controller.go` (public)
4. Controller: `app/http/controllers/admin/landing_section_controller.go` (admin CRUD)
5. Routes: register in `routes/api.go`
6. Seeder: populate with current hardcoded landing page data
