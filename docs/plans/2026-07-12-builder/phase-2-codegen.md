# Phase 2: Backend Publish + Codegen Foundation

**Status:** ✅ done
**Commit:** _(see git log)_
**Goal:** Publish endpoint run codegen tree → HTML, simpan `published_html`. Storefront serve HTML tersebut.

## Checklist

- [x] Service `backend/app/services/codegen.go` — walk tree JSON → HTML string
- [x] Field `PublishedHTML string` di model `landing_page.go` (sudah ada dari Phase 1)
- [x] Kolom `published_html` di migration (sudah ada dari Phase 1)
- [x] Controller Publish: jalankan codegen → simpan ke `published_html`
- [x] `routes/web.go`: storefront `/` cari published page → serve `published_html`
- [x] Fallback: tidak ada published page → placeholder HTML (`defaultStorefrontHTML`)
- [x] Codegen deterministic (urutan atribut stabil: class → src → alt → href)
- [x] Wrap output dengan `<!DOCTYPE html>` + Tailwind CDN script
- [x] `go build ./...` OK
- [x] Tests: 8 unit tests pass (deterministic, nested, attrs, XSS escape, document wrap)

## Codegen Spec

```
node.type    → element
─────────────────────────────
text         <span class="...">{props.text}</span>
heading      <h{level} class="...">{props.text}</h{level}>
frame        <div class="...">{children}</div>
section      <section class="...">{children}</section>
image        <img class="..." src={props.src} alt={props.alt} />
button       <button class="...">{props.text}</button>
link         <a class="..." href={props.href}>{children|text}</a>
component    resolve instanceOverrides → render master tree
```

`class=""` = `join(' ', node.classes)`. Atribut urut: class, src, alt, href, dst (stabil).

## Hasil Verifikasi

```
login 200, update 200, publish 200 (html_len=610)
storefront /  200 (610 bytes)
```

Codegen output deterministik:
```html
<section class="min-h-screen flex flex-col items-center justify-center">
  <h1 class="text-5xl font-bold">okuru.id</h1>
  <span class="mt-4 text-lg">Build faster.</span>
  <div class="mt-8 flex gap-4">
    <button class="bg-blue-600 text-white px-6 py-3 rounded-lg">Get Started</button>
    <a class="text-blue-600" href="https://docs.okuru.id"><span>Learn more</span></a>
  </div>
</section>
```

## Catatan

- Field `PublishedHTML` + kolom `published_html` sudah ditambah di Phase 1 (migration awal). Tidak perlu migration alter.
- `ponytail:` component instance resolution (`componentId` → master tree) deferred ke Phase 5.
- `ponytail:` placeholder storefront = static string, no template engine. Upgrade saat butuh server-side data injection.
- `ponytail:` class dedup belum — tambah saat style panel emit duplikat (Phase 6).
- Tests di `app/services/codegen_test.go`, bukan `tests/feature/` — pure function, no Goravel bootstrap perlu.

## Files

- Create: `backend/app/services/codegen.go`
- Modify: `backend/app/models/landing_page.go`
- Modify: `backend/database/migrations/*_add_published_html_to_landing_pages.go`
- Modify: `backend/app/http/controllers/admin/landing_page_controller.go`
- Modify: `backend/routes/web.go`

## Commit

```bash
git add backend/app/services/codegen.go
git add -A
git commit -m "Phase 2: publish pipeline, tree→HTML codegen, storefront render"
```
