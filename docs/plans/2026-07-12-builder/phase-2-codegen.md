# Phase 2: Backend Publish + Codegen Foundation

**Status:** ⬜ todo
**Goal:** Publish endpoint run codegen tree → HTML, simpan `published_html`. Storefront serve HTML tersebut.

## Checklist

- [ ] Service `backend/app/services/codegen.go` — walk tree JSON → HTML string
- [ ] Tambah field `PublishedHTML string` di model `landing_page.go`
- [ ] Migration alter table tambah kolom `published_html`
- [ ] Controller Publish: jalankan codegen → simpan ke `published_html`
- [ ] `routes/web.go`: storefront `/` cari published page → serve `published_html`
- [ ] Fallback: tidak ada published page → placeholder HTML
- [ ] Codegen deterministic (urutan atribut stabil)
- [ ] Wrap output dengan `<!DOCTYPE html>` + Tailwind CDN script
- [ ] `go build ./...` OK

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
