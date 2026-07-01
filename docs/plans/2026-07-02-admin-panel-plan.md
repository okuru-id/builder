# Okuru.id Admin Panel Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.
> **Note:** User testing policy — JANGAN buat test file. Verifikasi via lint, typecheck, manual check.

**Goal:** Backend API + admin panel untuk okuru.id personal website, dengan Goravel + SQLite + Vue admin SPA + Docker deployment via GitLab CI.

**Architecture:** Goravel single binary serve API + embed Vue admin SPA. SQLite untuk storage. Docker Compose (app + Caddy) deploy ke VPS 512MB. GitLab CI build image, push ke registry, deploy via SSH.

**Tech Stack:** Go 1.23 + Goravel v1.17, SQLite (CGO), Vue 3 + Vite + shadcn-vue (Reka UI), Go html/template + Alpine.js untuk storefront, Caddy 2, Docker, GitLab CI.

**Design Doc:** `docs/plans/2026-07-02-admin-panel-design.md`

---

## Phase 1: Project Scaffolding

### Task 1: Init monorepo structure

**Files:**
- Create: `okuru.id/` root directory (sudah ada)
- Create: `.gitignore`
- Create: `Makefile`
- Create: `README.md`

**Step 1: Init git repo**

```bash
cd /home/kurob/code/okuru.id
git init
git branch -M main
```

**Step 2: Create .gitignore**

```gitignore
# Backend
backend/.env
backend/storage/logs/*
backend/storage/app/public/*
backend/public/admin/*
backend/bin/

# Admin
admin/node_modules/
admin/dist/

# Docker
deploy/.env

# OS
.DS_Store
*.swp
*.swo

# IDE
.idea/
.vscode/
*.code-workspace

# Env
.env
.env.local
.env.*.local
```

**Step 3: Create root Makefile**

```makefile
.PHONY: dev-admin dev-backend dev build-local docker-up docker-down deploy migrate seed backup

# === LOCAL DEV ===
dev-admin:
	cd admin && pnpm run dev

dev-backend:
	cd backend && air

dev:
	@make dev-backend & make dev-admin & wait

# === DOCKER ===
build-local:
	docker compose -f deploy/docker-compose.yml build

docker-up:
	docker compose -f deploy/docker-compose.yml up -d

docker-down:
	docker compose -f deploy/docker-compose.yml down

docker-logs:
	docker compose -f deploy/docker-compose.yml logs -f

# === DATABASE ===
migrate:
	docker compose -f deploy/docker-compose.yml exec app ./artisan migrate

seed:
	docker compose -f deploy/docker-compose.yml exec app ./artisan db:seed

# === DEPLOY ===
deploy:
	docker compose -f deploy/docker-compose.yml pull
	docker compose -f deploy/docker-compose.yml up -d --remove-orphans
	docker compose -f deploy/docker-compose.yml exec -T app ./artisan migrate --force

# === BACKUP ===
backup:
	docker compose -f deploy/docker-compose.yml exec app /opt/okuru/scripts/backup.sh
```

**Step 4: Commit**

```bash
git add .gitignore Makefile
git commit -m "chore: init monorepo structure"
```

---

### Task 2: Scaffold Goravel backend

**Files:**
- Create: `backend/` (Goravel project)

**Step 1: Install Goravel installer**

```bash
go install github.com/goravel/installer/goravel@latest
```

**Step 2: Create Goravel project**

```bash
cd /home/kurob/code/okuru.id
goravel new backend
```

**Step 3: Verify setup**

```bash
cd backend
cp .env.example .env
./artisan key:generate
go run .
```

Expected: Server berjalan di `localhost:3000`.

**Step 4: Install SQLite driver**

```bash
cd backend
go get github.com/goravel/sqlite
```

**Step 5: Configure database**

Edit `backend/config/database.go`:

```go
"connections": map[string]any{
    "sqlite": map[string]any{
        "driver":   "sqlite",
        "database": config.Env("DB_DATABASE", "okuru.db"),
        "prefix":   "",
    },
},
```

Set default connection ke sqlite:

```go
"default": config.Env("DB_CONNECTION", "sqlite"),
```

Edit `backend/.env`:

```env
DB_CONNECTION=sqlite
DB_DATABASE=okuru.db
```

**Step 6: Verify SQLite works**

```bash
cd backend
go run .
# Buka browser: http://localhost:3000
# Harus tidak error
```

**Step 7: Commit**

```bash
git add backend/
git commit -m "chore: scaffold Goravel backend with SQLite"
```

---

### Task 3: Scaffold Vue admin SPA

**Files:**
- Create: `admin/` (Vue + Vite + shadcn-vue)

**Step 1: Create Vite Vue project**

```bash
cd /home/kurob/code/okuru.id
pnpm create vite admin --template vue-ts
cd admin
pnpm install
```

**Step 2: Configure TypeScript path aliases**

Edit `admin/tsconfig.json`:

```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    },
    "types": ["vite/client"]
  },
  "include": ["src/**/*.ts", "src/**/*.d.ts", "src/**/*.tsx", "src/**/*.vue"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

**Step 3: Configure Vite**

Edit `admin/vite.config.ts`:

```typescript
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      "/admin/api": {
        target: "http://localhost:3000",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "dist",
  },
});
```

**Step 4: Init shadcn-vue**

```bash
cd admin
pnpm dlx shadcn-vue@latest init --preset a4Xy5c93 --template vite --pointer
```

Ikuti prompt:
- Style: New York (atau pilihan dari preset)
- Base color: Slate
- CSS variables: Yes

**Step 5: Install dependencies dasar**

```bash
cd admin
pnpm add vue-router@4 @vee-validate/core @vee-validate/zod zod
pnpm add axios
pnpm add @vueuse/core
pnpm add lucide-vue-next
```

**Step 6: Add shadcn components dasar**

```bash
cd admin
pnpm dlx shadcn-vue@latest add button card input label form dialog dropdown-menu sonner sidebar table badge avatar separator tabs sheet skeleton
```

**Step 7: Setup Vue Router**

Create `admin/src/router/index.ts`:

```typescript
import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory("/admin/"),
  routes: [
    {
      path: "/admin/login",
      name: "login",
      component: () => import("@/views/Login.vue"),
      meta: { public: true },
    },
    {
      path: "/admin",
      component: () => import("@/layouts/AdminLayout.vue"),
      children: [
        { path: "", name: "dashboard", component: () => import("@/views/Dashboard.vue") },
        { path: "posts", name: "posts", component: () => import("@/views/Posts.vue") },
        { path: "posts/new", name: "post-new", component: () => import("@/views/PostEditor.vue") },
        { path: "posts/:id/edit", name: "post-edit", component: () => import("@/views/PostEditor.vue") },
        { path: "projects", name: "projects", component: () => import("@/views/Projects.vue") },
        { path: "open-source", name: "open-source", component: () => import("@/views/OpenSource.vue") },
        { path: "products", name: "products", component: () => import("@/views/Products.vue") },
        { path: "inbox", name: "inbox", component: () => import("@/views/Inbox.vue") },
        { path: "settings", name: "settings", component: () => import("@/views/Settings.vue") },
      ],
    },
  ],
});

export default router;
```

**Step 8: Mount router**

Edit `admin/src/main.ts`:

```typescript
import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";
import router from "./router";

createApp(App).use(router).mount("#app");
```

**Step 9: Verify build**

```bash
cd admin
pnpm run build
```

Expected: `dist/index.html` + `dist/assets/*.js` generated.

**Step 10: Commit**

```bash
git add admin/
git commit -m "chore: scaffold Vue admin SPA with shadcn-vue"
```

---

## Phase 2: Database Migrations & Models

### Task 4: Create database migrations

**Files:**
- Create: `backend/database/migrations/*.go`

**Step 1: Generate migration files**

```bash
cd backend
./artisan make:migration create_users_table
./artisan make:migration create_categories_table
./artisan make:migration create_posts_table
./artisan make:migration create_projects_table
./artisan make:migration create_open_source_projects_table
./artisan make:migration create_products_table
./artisan make:migration create_orders_table
./artisan make:migration create_messages_table
./artisan make:migration create_settings_table
```

**Step 2: Write users migration**

Edit `backend/database/migrations/<timestamp>_create_users_table.go`:

```go
package migrations

import (
    "github.com/goravel/framework/contracts/database/schema"
    "github.com/goravel/framework/facades"
)

func init() {
    facades.Schema().Create("users", func(table schema.Blueprint) {
        table.ID()
        table.String("email").Unique()
        table.String("password")
        table.String("totp_secret").Nullable()
        table.Boolean("totp_verified").Default(false)
        table.Timestamps()
    })
}
```

**Step 3: Write categories migration**

```go
func init() {
    facades.Schema().Create("categories", func(table schema.Blueprint) {
        table.ID()
        table.String("slug").Unique()
        table.String("name_en")
        table.String("name_id")
        table.Timestamps()
    })
}
```

**Step 4: Write posts migration**

```go
func init() {
    facades.Schema().Create("posts", func(table schema.Blueprint) {
        table.ID()
        table.String("slug").Unique()
        table.String("title_en")
        table.String("title_id")
        table.Text("excerpt_en").Nullable()
        table.Text("excerpt_id").Nullable()
        table.LongText("content_en").Nullable()
        table.LongText("content_id").Nullable()
        table.String("category").Nullable()
        table.JSON("tags").Nullable()
        table.String("thumbnail").Nullable()
        table.String("status").Default("draft")
        table.Timestamp("published_at").Nullable()
        table.Integer("read_time").Default(0)
        table.Timestamps()
    })
}
```

**Step 5: Write projects migration**

```go
func init() {
    facades.Schema().Create("projects", func(table schema.Blueprint) {
        table.ID()
        table.Integer("sort_order").Default(0)
        table.String("title_en")
        table.String("title_id")
        table.Text("description_en").Nullable()
        table.Text("description_id").Nullable()
        table.JSON("tech_stack").Nullable()
        table.String("url").Nullable()
        table.Boolean("featured").Default(false)
        table.Timestamps()
    })
}
```

**Step 6: Write open_source_projects migration**

```go
func init() {
    facades.Schema().Create("open_source_projects", func(table schema.Blueprint) {
        table.ID()
        table.Integer("sort_order").Default(0)
        table.String("title_en")
        table.String("title_id")
        table.Text("description_en").Nullable()
        table.Text("description_id").Nullable()
        table.String("github_url")
        table.JSON("technologies").Nullable()
        table.Integer("stars").Default(0)
        table.String("license").Nullable()
        table.Timestamps()
    })
}
```

**Step 7: Write products migration**

```go
func init() {
    facades.Schema().Create("products", func(table schema.Blueprint) {
        table.ID()
        table.String("slug").Unique()
        table.String("title")
        table.Text("description").Nullable()
        table.Integer("price").Default(0)
        table.String("type")
        table.String("file_path").Nullable()
        table.String("thumbnail").Nullable()
        table.String("status").Default("active")
        table.Timestamps()
    })
}
```

**Step 8: Write orders migration**

```go
func init() {
    facades.Schema().Create("orders", func(table schema.Blueprint) {
        table.ID()
        table.String("reference").Unique()
        table.UnsignedBigInteger("product_id")
        table.String("buyer_email")
        table.String("buyer_name")
        table.Integer("amount").Default(0)
        table.String("status").Default("pending")
        table.String("payment_ref").Nullable()
        table.String("download_token").Nullable()
        table.Integer("download_count").Default(0)
        table.Integer("max_downloads").Default(3)
        table.Timestamp("expires_at").Nullable()
        table.Timestamp("paid_at").Nullable()
        table.Timestamps()
    })
}
```

**Step 9: Write messages migration**

```go
func init() {
    facades.Schema().Create("messages", func(table schema.Blueprint) {
        table.ID()
        table.String("name")
        table.String("email")
        table.Text("message")
        table.String("status").Default("unread")
        table.Timestamps()
    })
}
```

**Step 10: Write settings migration**

```go
func init() {
    facades.Schema().Create("settings", func(table schema.Blueprint) {
        table.ID()
        table.String("key").Unique()
        table.Text("value").Nullable()
        table.Timestamps()
    })
}
```

**Step 11: Run migrations**

```bash
cd backend
./artisan migrate
```

Expected: semua tabel tercreate. Cek `okuru.db` dengan `sqlite3 okuru.db ".tables"`.

**Step 12: Commit**

```bash
git add backend/database/migrations/
git commit -m "feat: add database migrations for all tables"
```

---

### Task 5: Create ORM models

**Files:**
- Create: `backend/app/models/*.go`

**Step 1: Generate model files**

```bash
cd backend
./artisan make:model User
./artisan make:model Category
./artisan make:model Post
./artisan make:model Project
./artisan make:model OpenSourceProject
./artisan make:model Product
./artisan make:model Order
./artisan make:model Message
./artisan make:model Setting
```

**Step 2: Write User model**

Edit `backend/app/models/user.go`:

```go
package models

import "github.com/goravel/framework/database/orm"

type User struct {
    orm.Model
    Email        string `gorm:"uniqueIndex" json:"email"`
    Password     string `gorm:"-" json:"-"`
    TotpSecret   *string `json:"-"`
    TotpVerified bool   `json:"-"`
    orm.Timestamps
}
```

**Step 3: Write Post model**

Edit `backend/app/models/post.go`:

```go
package models

import "github.com/goravel/framework/database/orm"

type Post struct {
    orm.Model
    Slug        string `gorm:"uniqueIndex" json:"slug"`
    TitleEn     string `json:"title_en"`
    TitleId     string `json:"title_id"`
    ExcerptEn   *string `json:"excerpt_en"`
    ExcerptId   *string `json:"excerpt_id"`
    ContentEn   *string `json:"content_en"`
    ContentId   *string `json:"content_id"`
    Category    *string `json:"category"`
    Tags        orm.JSON `json:"tags"`
    Thumbnail   *string `json:"thumbnail"`
    Status      string `json:"status"`
    PublishedAt *orm.Time `json:"published_at"`
    ReadTime    int `json:"read_time"`
    orm.Timestamps
}
```

**Step 4: Tulis model lainnya mengikuti pattern** (Category, Project, OpenSourceProject, Product, Order, Message, Setting) — field sesuai migration.

**Step 5: Verify compile**

```bash
cd backend
go build ./...
```

Expected: no errors.

**Step 6: Commit**

```bash
git add backend/app/models/
git commit -m "feat: add ORM models for all tables"
```

---

### Task 6: Create seeder

**Files:**
- Create: `backend/database/seeders/seeder.go`

**Step 1: Generate seeder**

```bash
cd backend
./artisan make:seeder Seeder
```

**Step 2: Write seeder**

Edit `backend/database/seeders/seeder.go`:

```go
package seeders

import (
    "github.com/goravel/framework/facades"

    "okuru/app/models"
)

type Seeder struct{}

func (s *Seeder) Run() error {
    if err := s.seedAdmin(); err != nil {
        return err
    }
    if err := s.seedCategories(); err != nil {
        return err
    }
    if err := s.seedSettings(); err != nil {
        return err
    }
    return nil
}

func (s *Seeder) seedAdmin() error {
    password, _ := facades.Hash().Make("changeme123")
    user := models.User{
        Email:    "admin@okuru.id",
        Password: password,
    }
    return facades.DB().Create(&user).Error
}

func (s *Seeder) seedCategories() error {
    categories := []models.Category{
        {Slug: "web-development", NameEn: "Web Development", NameId: "Pengembangan Web"},
        {Slug: "devops", NameEn: "DevOps", NameId: "DevOps"},
        {Slug: "javascript", NameEn: "JavaScript", NameId: "JavaScript"},
        {Slug: "laravel", NameEn: "Laravel", NameId: "Laravel"},
        {Slug: "react", NameEn: "React", NameId: "React"},
        {Slug: "docker", NameEn: "Docker", NameId: "Docker"},
    }
    for _, c := range categories {
        if err := facades.DB().Create(&c).Error; err != nil {
            return err
        }
    }
    return nil
}

func (s *Seeder) seedSettings() error {
    settings := []models.Setting{
        {Key: "hero_title_en", Value: "Hi, I'm Kurob"},
        {Key: "hero_title_id", Value: "Halo, Saya Kurob"},
        {Key: "hero_desc_en", Value: "I'm from Indonesia with 8+ years of experience as a full-stack web developer."},
        {Key: "hero_desc_id", Value: "Saya dari Indonesia dengan pengalaman 8+ tahun sebagai full-stack web developer."},
        {Key: "start_year", Value: "2016"},
    }
    for _, st := range settings {
        if err := facades.DB().Create(&st).Error; err != nil {
            return err
        }
    }
    return nil
}
```

**Step 3: Run seeder**

```bash
cd backend
./artisan db:seed
```

**Step 4: Commit**

```bash
git add backend/database/seeders/
git commit -m "feat: add database seeder for admin user, categories, settings"
```

---

## Phase 3: Authentication

### Task 7: Setup auth guard & JWT

**Files:**
- Modify: `backend/config/auth.go`
- Modify: `backend/config/jwt.go`

**Step 1: Generate JWT secret**

```bash
cd backend
./artisan jwt:secret
```

**Step 2: Configure auth guard**

Edit `backend/config/auth.go`:

```go
"defaults": config.Env("AUTH_GUARD", "admin"),

"guards": map[string]any{
    "admin": map[string]any{
        "driver":   "jwt",
        "provider": "users",
        "ttl":      15,
        "refresh_ttl": 20160, // 7 days in minutes
    },
},

"providers": map[string]any{
    "users": map[string]any{
        "driver": "orm",
        "model":  "models.User",
    },
},
```

**Step 3: Configure JWT**

Edit `backend/config/jwt.go` — pastikan secret terisi dari env.

**Step 4: Verify**

```bash
cd backend
go build ./...
```

**Step 5: Commit**

```bash
git add backend/config/
git commit -m "feat: configure JWT auth guard for admin"
```

---

### Task 8: TOTP library & encryption

**Files:**
- Modify: `backend/go.mod`
- Create: `backend/app/services/totp_service.go`

**Step 1: Install TOTP library**

```bash
cd backend
go get github.com/pquerna/otp/totp
```

**Step 2: Write TOTP service**

Create `backend/app/services/totp_service.go`:

```go
package services

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base32"
    "errors"
    "io"

    "github.com/pquerna/otp"
    "github.com/pquerna/otp/totp"
)

type TotpService struct {
    appKey []byte
}

func NewTotpService(appKey string) *TotpService {
    key := []byte(appKey)
    if len(key) < 32 {
        padded := make([]byte, 32)
        copy(padded, key)
        key = padded
    }
    return &TotpService{appKey: key[:32]}
}

// GenerateSecret membuat secret baru untuk user
func (s *TotpService) GenerateSecret(email string) (string, string, error) {
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "Okuru.id",
        AccountName: email,
    })
    if err != nil {
        return "", "", err
    }
    return key.Secret(), key.URL(), nil
}

// ValidateCode memverifikasi kode TOTP
func (s *TotpService) ValidateCode(secret, code string) bool {
    return totp.Validate(code, secret)
}

// EncryptSecret encrypt secret dengan AES-256-GCM
func (s *TotpService) EncryptSecret(plaintext string) (string, error) {
    block, err := aes.NewCipher(s.appKey)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
    return base32.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSecret decrypt secret
func (s *TotpService) DecryptSecret(encoded string) (string, error) {
    ciphertext, err := base32.StdEncoding.DecodeString(encoded)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(s.appKey)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return "", errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
```

**Step 3: Verify compile**

```bash
cd backend
go build ./...
```

**Step 4: Commit**

```bash
git add backend/app/services/totp_service.go backend/go.mod backend/go.sum
git commit -m "feat: add TOTP service with AES-256-GCM encryption"
```

---

### Task 9: Auth controller (login + TOTP)

**Files:**
- Create: `backend/app/http/controllers/admin/auth_controller.go`
- Create: `backend/app/http/requests/admin/login_request.go`
- Create: `backend/app/http/requests/admin/totp_request.go`

**Step 1: Write login request validation**

Create `backend/app/http/requests/admin/login_request.go`:

```go
package admin

import "github.com/goravel/framework/contracts/http"

type LoginRequest struct {
    Email    string `form:"email" json:"email" validate:"required,email"`
    Password string `form:"password" json:"password" validate:"required,min:8"`
}

func (r *LoginRequest) Authorize(ctx http.Context) bool {
    return true
}

func (r *LoginRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "email":    "required|email",
        "password": "required|min:8",
    }
}
```

**Step 2: Write TOTP request validation**

Create `backend/app/http/requests/admin/totp_request.go`:

```go
package admin

import "github.com/goravel/framework/contracts/http"

type TotpRequest struct {
    TempToken string `json:"temp_token" validate:"required"`
    Code      string `json:"code" validate:"required,len:6"`
}

func (r *TotpRequest) Authorize(ctx http.Context) bool {
    return true
}

func (r *TotpRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "temp_token": "required",
        "code":       "required|len:6",
    }
}
```

**Step 3: Write auth controller**

Create `backend/app/http/controllers/admin/auth_controller.go`:

```go
package admin

import (
    "errors"
    "time"

    "github.com/goravel/framework/auth"
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"

    "okuru/app/models"
    "okuru/app/services"
)

type AuthController struct {
    totpService *services.TotpService
}

func NewAuthController() *AuthController {
    appKey := facades.Config().GetString("app.key")
    return &AuthController{
        totpService: services.NewTotpService(appKey),
    }
}

// POST /admin/api/auth/login
func (c *AuthController) Login(ctx http.Context) {
    var req LoginRequest
    if err := ctx.Request().Bind(&req); err != nil {
        ctx.Response().Json(http.StatusBadRequest, http.Json{
            "error": "Invalid request",
        })
        return
    }

    // Find user by email
    var user models.User
    if err := facades.DB().Where("email = ?", req.Email).First(&user); err != nil {
        ctx.Response().Json(http.StatusUnauthorized, http.Json{
            "error": "Invalid credentials",
        })
        return
    }

    // Verify password
    if !facades.Hash().Check(req.Password, user.Password) {
        ctx.Response().Json(http.StatusUnauthorized, http.Json{
            "error": "Invalid credentials",
        })
        return
    }

    // Check if TOTP is setup
    if user.TotpSecret != nil && user.TotpVerified {
        // Generate temp token (short-lived, scope: totp_verify)
        tempToken, err := c.generateTempToken(user.ID)
        if err != nil {
            ctx.Response().Json(http.StatusInternalServerError, http.Json{
                "error": "Failed to generate token",
            })
            return
        }
        ctx.Response().Json(http.StatusOK, http.Json{
            "requires_totp": true,
            "temp_token":    tempToken,
        })
    } else {
        // TOTP not setup, generate full token
        token, err := c.generateFullToken(user.ID)
        if err != nil {
            ctx.Response().Json(http.StatusInternalServerError, http.Json{
                "error": "Failed to generate token",
            })
            return
        }
        ctx.Response().Json(http.StatusOK, http.Json{
            "requires_totp": false,
            "access_token":  token,
            "totp_setup_required": true,
        })
    }
}

// POST /admin/api/auth/totp
func (c *AuthController) VerifyTotp(ctx http.Context) {
    var req TotpRequest
    if err := ctx.Request().Bind(&req); err != nil {
        ctx.Response().Json(http.StatusBadRequest, http.Json{
            "error": "Invalid request",
        })
        return
    }

    // Parse temp token
    payload, err := facades.Auth(ctx).Parse(req.TempToken)
    if err != nil {
        ctx.Response().Json(http.StatusUnauthorized, http.Json{
            "error": "Invalid or expired token",
        })
        return
    }

    userID := payload.Key

    // Get user
    var user models.User
    if err := facades.DB().Where("id = ?", userID).First(&user); err != nil {
        ctx.Response().Json(http.StatusUnauthorized, http.Json{
            "error": "User not found",
        })
        return
    }

    // Decrypt TOTP secret
    if user.TotpSecret == nil {
        ctx.Response().Json(http.StatusBadRequest, http.Json{
            "error": "TOTP not setup",
        })
        return
    }

    secret, err := c.totpService.DecryptSecret(*user.TotpSecret)
    if err != nil {
        ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "error": "Failed to decrypt secret",
        })
        return
    }

    // Validate code
    if !c.totpService.ValidateCode(secret, req.Code) {
        ctx.Response().Json(http.StatusUnauthorized, http.Json{
            "error": "Invalid TOTP code",
        })
        return
    }

    // Generate full token with TOTP verified claim
    token, err := c.generateFullToken(user.ID)
    if err != nil {
        ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "error": "Failed to generate token",
        })
        return
    }

    ctx.Response().Json(http.StatusOK, http.Json{
        "access_token": token,
    })
}

// POST /admin/api/auth/totp/setup
func (c *AuthController) SetupTotp(ctx http.Context) {
    userID, err := facades.Auth(ctx).ID()
    if err != nil {
        ctx.Response().Json(http.StatusUnauthorized, http.Json{
            "error": "Unauthorized",
        })
        return
    }

    var user models.User
    if err := facades.DB().Where("id = ?", userID).First(&user); err != nil {
        ctx.Response().Json(http.StatusNotFound, http.Json{
            "error": "User not found",
        })
        return
    }

    secret, qrUrl, err := c.totpService.GenerateSecret(user.Email)
    if err != nil {
        ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "error": "Failed to generate TOTP secret",
        })
        return
    }

    // Encrypt and store
    encrypted, err := c.totpService.EncryptSecret(secret)
    if err != nil {
        ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "error": "Failed to encrypt secret",
        })
        return
    }

    user.TotpSecret = &encrypted
    user.TotpVerified = false
    facades.DB().Save(&user)

    ctx.Response().Json(http.StatusOK, http.Json{
        "secret": secret,
        "qr_url": qrUrl,
    })
}

// POST /admin/api/auth/totp/verify-setup
func (c *AuthController) VerifyTotpSetup(ctx http.Context) {
    var req struct {
        Code string `json:"code"`
    }
    if err := ctx.Request().Bind(&req); err != nil {
        ctx.Response().Json(http.StatusBadRequest, http.Json{
            "error": "Invalid request",
        })
        return
    }

    userID, _ := facades.Auth(ctx).ID()

    var user models.User
    if err := facades.DB().Where("id = ?", userID).First(&user); err != nil {
        ctx.Response().Json(http.StatusNotFound, http.Json{
            "error": "User not found",
        })
        return
    }

    if user.TotpSecret == nil {
        ctx.Response().Json(http.StatusBadRequest, http.Json{
            "error": "TOTP not setup",
        })
        return
    }

    secret, err := c.totpService.DecryptSecret(*user.TotpSecret)
    if err != nil {
        ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "error": "Failed to decrypt secret",
        })
        return
    }

    if !c.totpService.ValidateCode(secret, req.Code) {
        ctx.Response().Json(http.StatusUnauthorized, http.Json{
            "error": "Invalid code",
        })
        return
    }

    user.TotpVerified = true
    facades.DB().Save(&user)

    ctx.Response().Json(http.StatusOK, http.Json{
        "verified": true,
    })
}

// GET /admin/api/auth/me
func (c *AuthController) Me(ctx http.Context) {
    userID, _ := facades.Auth(ctx).ID()

    var user models.User
    if err := facades.DB().Where("id = ?", userID).First(&user); err != nil {
        ctx.Response().Json(http.StatusNotFound, http.Json{
            "error": "User not found",
        })
        return
    }

    ctx.Response().Json(http.StatusOK, http.Json{
        "id":            user.ID,
        "email":         user.Email,
        "totp_verified": user.TotpVerified,
    })
}

// Helper: generate temp token (2 min TTL)
func (c *AuthController) generateTempToken(userID uint) (string, error) {
    return facades.Auth(ctx).LoginUsingID(userID)
}

// Helper: generate full token
func (c *AuthController) generateFullToken(userID uint) (string, error) {
    return facades.Auth(ctx).LoginUsingID(userID)
}
```

Note: `ctx` scope issue di helper — refactor untuk pass ctx. Lihat implementasi aktual.

**Step 4: Verify compile**

```bash
cd backend
go build ./...
```

**Step 5: Commit**

```bash
git add backend/app/http/controllers/admin/auth_controller.go backend/app/http/requests/admin/
git commit -m "feat: add auth controller with TOTP two-phase login"
```

---

### Task 10: Auth middleware (JWT + TOTP)

**Files:**
- Create: `backend/app/http/middleware/totp_verified.go`

**Step 1: Write TOTP verified middleware**

Create `backend/app/http/middleware/totp_verified.go`:

```go
package middleware

import (
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"

    "okuru/app/models"
)

func TotpVerified() http.Middleware {
    return func(ctx http.Context) {
        // Parse token dari Authorization header
        token := ctx.Request().Header("Authorization")
        if token == "" {
            ctx.Response().Json(http.StatusUnauthorized, http.Json{
                "error": "Token required",
            })
            ctx.Abort()
            return
        }

        // Parse JWT
        payload, err := facades.Auth(ctx).Parse(token)
        if err != nil {
            ctx.Response().Json(http.StatusUnauthorized, http.Json{
                "error": "Invalid token",
            })
            ctx.Abort()
            return
        }

        // Cek user TOTP status
        var user models.User
        if err := facades.DB().Where("id = ?", payload.Key).First(&user); err != nil {
            ctx.Response().Json(http.StatusUnauthorized, http.Json{
                "error": "User not found",
            })
            ctx.Abort()
            return
        }

        // Jika TOTP secret set tapi belum verified, tolak akses
        if user.TotpSecret != nil && !user.TotpVerified {
            ctx.Response().Json(http.StatusForbidden, http.Json{
                "error":          "TOTP verification required",
                "totp_required":  true,
            })
            ctx.Abort()
            return
        }

        ctx.With("user_id", payload.Key)
        ctx.Request().Next()
    }
}
```

**Step 2: Register middleware**

Edit `backend/app/providers/route_service_provider.go` atau `backend/kernel.go` — tambahkan middleware alias:

```go
"TotpVerified": middleware.TotpVerified(),
```

**Step 3: Verify compile**

```bash
cd backend
go build ./...
```

**Step 4: Commit**

```bash
git add backend/app/http/middleware/totp_verified.go
git commit -m "feat: add TOTP verified middleware"
```

---

### Task 11: Auth routes

**Files:**
- Modify: `backend/routes/api.go`

**Step 1: Write auth routes**

Edit `backend/routes/api.go`:

```go
package routes

import (
    "github.com/goravel/framework/facades"

    "okuru/app/http/controllers/admin"
)

func Api() {
    authController := admin.NewAuthController()

    facades.Route().Prefix("admin/api").Group(func(router route.Router) {
        // Public auth routes
        router.Post("/auth/login", authController.Login)
        router.Post("/auth/totp", authController.VerifyTotp)

        // Protected routes (JWT only, untuk setup TOTP pertama kali)
        authGroup := router.Group(func(r route.Router) {
            r.Use(middleware.Auth())
        })
        authGroup.Post("/auth/totp/setup", authController.SetupTotp)
        authGroup.Post("/auth/totp/verify-setup", authController.VerifyTotpSetup)
        authGroup.Get("/auth/me", authController.Me)

        // Fully protected routes (JWT + TOTP verified)
        totpGroup := router.Group(func(r route.Router) {
            r.Use(middleware.Auth())
            r.Use(middleware.TotpVerified())
        })
        // CRUD routes akan ditambahkan di sini (Task selanjutnya)
    })
}
```

**Step 2: Verify compile & run**

```bash
cd backend
go build ./...
go run .
```

Test login manual:
```bash
curl -X POST http://localhost:3000/admin/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@okuru.id","password":"changeme123"}'
```

Expected: `{"requires_totp":false,"access_token":"...","totp_setup_required":true}`

**Step 3: Commit**

```bash
git add backend/routes/api.go
git commit -m "feat: add auth routes for login and TOTP"
```

---

## Phase 4: Public API

### Task 12: Content controller (public endpoints)

**Files:**
- Create: `backend/app/http/controllers/content_controller.go`

**Step 1: Write content controller**

Create `backend/app/http/controllers/content_controller.go`:

```go
package controllers

import (
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"

    "okuru/app/models"
)

type ContentController struct{}

func NewContentController() *ContentController {
    return &ContentController{}
}

// GET /api/v1/content — batch endpoint untuk React frontend
func (c *ContentController) Index(ctx http.Context) {
    content := http.Json{}

    // Hero / profile settings
    settings := c.getSettingsMap()
    content["meta"] = settings["meta"]
    content["hero"] = http.Json{
        "title":         settings["hero_title"],
        "description":   settings["hero_desc"],
        "startYear":     settings["start_year"],
        "highlightedRole": settings["highlighted_role"],
        "cta":           settings["cta"],
        "socialLinks":   settings["social_links"],
    }

    // Projects
    var projects []models.Project
    facades.DB().OrderBy("sort_order", "asc").Get(&projects)
    content["projects"] = projects

    // Open source
    var oss []models.OpenSourceProject
    facades.DB().OrderBy("sort_order", "asc").Get(&oss)
    content["openSource"] = oss

    // Latest blog posts (5)
    var posts []models.Post
    facades.DB().Where("status = ?", "published").
        OrderBy("published_at", "desc").
        Limit(5).
        Get(&posts)
    content["blog"] = posts

    // Skills
    content["skills"] = settings["skills"]

    // Profile
    content["profile"] = settings["profile"]

    ctx.Response().Json(http.StatusOK, content)
}

// GET /api/v1/posts
func (c *ContentController) Posts(ctx http.Context) {
    page := ctx.Request().QueryInt("page", 1)
    category := ctx.Request().Query("category", "")

    query := facades.DB().Where("status = ?", "published").
        OrderBy("published_at", "desc")

    if category != "" {
        query = query.Where("category = ?", category)
    }

    var posts []models.Post
    query.Offset((page - 1) * 10).Limit(10).Get(&posts)

    ctx.Response().Json(http.StatusOK, http.Json{
        "data": posts,
        "page": page,
    })
}

// GET /api/v1/posts/:slug
func (c *ContentController) Post(ctx http.Context) {
    slug := ctx.Request().Route("slug")

    var post models.Post
    if err := facades.DB().Where("slug = ?", slug).
        Where("status = ?", "published").
        First(&post); err != nil {
        ctx.Response().Json(http.StatusNotFound, http.Json{
            "error": "Post not found",
        })
        return
    }

    ctx.Response().Json(http.StatusOK, post)
}

// GET /api/v1/projects
func (c *ContentController) Projects(ctx http.Context) {
    var projects []models.Project
    facades.DB().OrderBy("sort_order", "asc").Get(&projects)
    ctx.Response().Json(http.StatusOK, projects)
}

// GET /api/v1/open-source
func (c *ContentController) OpenSource(ctx http.Context) {
    var oss []models.OpenSourceProject
    facades.DB().OrderBy("sort_order", "asc").Get(&oss)
    ctx.Response().Json(http.StatusOK, oss)
}

// GET /api/v1/categories
func (c *ContentController) Categories(ctx http.Context) {
    var categories []models.Category
    facades.DB().Get(&categories)
    ctx.Response().Json(http.StatusOK, categories)
}

// GET /api/v1/settings/public
func (c *ContentController) PublicSettings(ctx http.Context) {
    ctx.Response().Json(http.StatusOK, c.getSettingsMap())
}

// POST /api/v1/contact
func (c *ContentController) Contact(ctx http.Context) {
    var req struct {
        Name    string `json:"name" validate:"required,min:2"`
        Email   string `json:"email" validate:"required,email"`
        Message string `json:"message" validate:"required,min:10"`
    }
    if err := ctx.Request().Bind(&req); err != nil {
        ctx.Response().Json(http.StatusBadRequest, http.Json{
            "error": "Invalid request",
        })
        return
    }

    message := models.Message{
        Name:    req.Name,
        Email:   req.Email,
        Message: req.Message,
        Status:  "unread",
    }
    facades.DB().Create(&message)

    ctx.Response().Json(http.StatusOK, http.Json{
        "success": true,
    })
}

// Helper: load settings ke map
func (c *ContentController) getSettingsMap() map[string]any {
    var settings []models.Setting
    facades.DB().Get(&settings)

    m := make(map[string]any)
    for _, s := range settings {
        m[s.Key] = s.Value
    }
    return m
}
```

**Step 2: Add public routes**

Edit `backend/routes/api.go` — tambahkan:

```go
func Api() {
    contentController := controllers.NewContentController()

    // Public API
    facades.Route().Prefix("api/v1").Group(func(router route.Router) {
        router.Get("/content", contentController.Index)
        router.Get("/posts", contentController.Posts)
        router.Get("/posts/{slug}", contentController.Post)
        router.Get("/projects", contentController.Projects)
        router.Get("/open-source", contentController.OpenSource)
        router.Get("/categories", contentController.Categories)
        router.Get("/settings/public", contentController.PublicSettings)
        router.Post("/contact", contentController.Contact)
    })

    // ... admin routes dari Task 11
}
```

**Step 3: Verify**

```bash
cd backend
go build ./...
go run .
```

Test:
```bash
curl http://localhost:3000/api/v1/content
curl -X POST http://localhost:3000/api/v1/contact \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@test.com","message":"Hello world test"}'
```

**Step 4: Commit**

```bash
git add backend/app/http/controllers/content_controller.go backend/routes/api.go
git commit -m "feat: add public content API endpoints"
```

---

### Task 13: CORS configuration

**Files:**
- Modify: `backend/config/cors.go`

**Step 1: Configure CORS**

Create or edit `backend/config/cors.go`:

```go
package config

import "github.com/goravel/framework/contracts/config"

func init() {
    facades.Config().Add("cors", map[string]any{
        "paths": []string{"api/*"},
        "allowed_methods": []string{"GET", "POST"},
        "allowed_origins": []string{
            "https://okuru.id",
            "https://www.okuru.id",
            "https://shop.okuru.id",
            "http://localhost:*",
        },
        "allowed_headers": []string{"Content-Type", "Authorization"},
        "exposed_headers": []string{},
        "max_age":         300,
    })
}
```

**Step 2: Register CORS middleware untuk public routes**

Edit `backend/routes/api.go`:

```go
facades.Route().Prefix("api/v1").Group(func(router route.Router) {
    router.Middleware(middleware.Cors())
    // ... routes
})
```

**Step 3: Verify**

```bash
cd backend
go build ./...
```

**Step 4: Commit**

```bash
git add backend/config/cors.go backend/routes/api.go
git commit -m "feat: configure CORS for public API"
```

---

## Phase 5: Admin CRUD Controllers

### Task 14: CRUD controllers (posts, projects, etc)

**Files:**
- Create: `backend/app/http/controllers/admin/post_controller.go`
- Create: `backend/app/http/controllers/admin/project_controller.go`
- Create: `backend/app/http/controllers/admin/open_source_controller.go`
- Create: `backend/app/http/controllers/admin/product_controller.go`
- Create: `backend/app/http/controllers/admin/message_controller.go`
- Create: `backend/app/http/controllers/admin/setting_controller.go`
- Modify: `backend/routes/api.go`

**Step 1: Write generic CRUD pattern**

Create `backend/app/http/controllers/admin/post_controller.go`:

```go
package admin

import (
    "strconv"

    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"

    "okuru/app/models"
)

type PostController struct{}

func NewPostController() *PostController {
    return &PostController{}
}

func (c *PostController) Index(ctx http.Context) {
    var posts []models.Post
    facades.DB().OrderBy("created_at", "desc").Get(&posts)
    ctx.Response().Json(http.StatusOK, posts)
}

func (c *PostController) Show(ctx http.Context) {
    id := ctx.Request().Route("id")
    var post models.Post
    if err := facades.DB().Where("id = ?", id).First(&post); err != nil {
        ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
        return
    }
    ctx.Response().Json(http.StatusOK, post)
}

func (c *PostController) Store(ctx http.Context) {
    var post models.Post
    if err := ctx.Request().Bind(&post); err != nil {
        ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
        return
    }
    // Auto-calc read time (200 wpm)
    if post.ContentEn != nil {
        words := len(strings.Fields(*post.ContentEn))
        post.ReadTime = int(math.Ceil(float64(words) / 200.0))
    }
    facades.DB().Create(&post)
    ctx.Response().Json(http.StatusCreated, post)
}

func (c *PostController) Update(ctx http.Context) {
    id := ctx.Request().Route("id")
    var post models.Post
    if err := facades.DB().Where("id = ?", id).First(&post); err != nil {
        ctx.Response().Json(http.StatusNotFound, http.Json{"error": "Not found"})
        return
    }
    if err := ctx.Request().Bind(&post); err != nil {
        ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
        return
    }
    facades.DB().Save(&post)
    ctx.Response().Json(http.StatusOK, post)
}

func (c *PostController) Destroy(ctx http.Context) {
    id := ctx.Request().Route("id")
    facades.DB().Where("id = ?", id).Delete(&models.Post{})
    ctx.Response().Json(http.StatusOK, http.Json{"deleted": true})
}
```

**Step 2: Tulis controller lainnya** (ProjectController, OpenSourceController, ProductController, MessageController, SettingController) mengikuti pattern yang sama, sesuaikan model.

**Step 3: Message controller** — special: mark as read/archive:

```go
func (c *MessageController) MarkRead(ctx http.Context) {
    id := ctx.Request().Route("id")
    facades.DB().Model(&models.Message{}).Where("id = ?", id).Update("status", "read")
    ctx.Response().Json(http.StatusOK, http.Json{"updated": true})
}

func (c *MessageController) Archive(ctx http.Context) {
    id := ctx.Request().Route("id")
    facades.DB().Model(&models.Message{}).Where("id = ?", id).Update("status", "archived")
    ctx.Response().Json(http.StatusOK, http.Json{"updated": true})
}
```

**Step 4: Setting controller** — key-value:

```go
func (c *SettingController) Index(ctx http.Context) {
    var settings []models.Setting
    facades.DB().Get(&settings)
    ctx.Response().Json(http.StatusOK, settings)
}

func (c *SettingController) Update(ctx http.Context) {
    var req struct {
        Key   string `json:"key"`
        Value string `json:"value"`
    }
    if err := ctx.Request().Bind(&req); err != nil {
        ctx.Response().Json(http.StatusBadRequest, http.Json{"error": err.Error()})
        return
    }
    facades.DB().Model(&models.Setting{}).
        Where("key = ?", req.Key).
        Update("value", req.Value)
    ctx.Response().Json(http.StatusOK, http.Json{"updated": true})
}
```

**Step 5: Register admin CRUD routes**

Edit `backend/routes/api.go` — tambahkan ke totpGroup:

```go
postController := admin.NewPostController()
projectController := admin.NewProjectController()
osController := admin.NewOpenSourceController()
productController := admin.NewProductController()
messageController := admin.NewMessageController()
settingController := admin.NewSettingController()

totpGroup.Get("/posts", postController.Index)
totpGroup.Post("/posts", postController.Store)
totpGroup.Get("/posts/{id}", postController.Show)
totpGroup.Put("/posts/{id}", postController.Update)
totpGroup.Delete("/posts/{id}", postController.Destroy)

totpGroup.Get("/projects", projectController.Index)
totpGroup.Post("/projects", projectController.Store)
totpGroup.Get("/projects/{id}", projectController.Show)
totpGroup.Put("/projects/{id}", projectController.Update)
totpGroup.Delete("/projects/{id}", projectController.Destroy)

totpGroup.Get("/open-source", osController.Index)
totpGroup.Post("/open-source", osController.Store)
totpGroup.Put("/open-source/{id}", osController.Update)
totpGroup.Delete("/open-source/{id}", osController.Destroy)

totpGroup.Get("/products", productController.Index)
totpGroup.Post("/products", productController.Store)
totpGroup.Put("/products/{id}", productController.Update)
totpGroup.Delete("/products/{id}", productController.Destroy)

totpGroup.Get("/messages", messageController.Index)
totpGroup.Get("/messages/{id}", messageController.Show)
totpGroup.Delete("/messages/{id}", messageController.Destroy)
totpGroup.Put("/messages/{id}/read", messageController.MarkRead)
totpGroup.Put("/messages/{id}/archive", messageController.Archive)

totpGroup.Get("/settings", settingController.Index)
totpGroup.Put("/settings", settingController.Update)

totpGroup.Get("/categories", categoryController.Index)
totpGroup.Post("/categories", categoryController.Store)
totpGroup.Put("/categories/{id}", categoryController.Update)
totpGroup.Delete("/categories/{id}", categoryController.Destroy)
```

**Step 6: Verify**

```bash
cd backend
go build ./...
go run .
```

Test dengan token dari login:
```bash
TOKEN="..." # dari login response
curl http://localhost:3000/admin/api/posts -H "Authorization: Bearer $TOKEN"
```

**Step 7: Commit**

```bash
git add backend/app/http/controllers/admin/ backend/routes/api.go
git commit -m "feat: add CRUD controllers for admin (posts, projects, products, etc)"
```

---

### Task 15: File upload controller

**Files:**
- Create: `backend/app/http/controllers/admin/upload_controller.go`

**Step 1: Write upload controller**

Create `backend/app/http/controllers/admin/upload_controller.go`:

```go
package admin

import (
    "crypto/rand"
    "encoding/hex"
    "path/filepath"
    "strings"

    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
)

type UploadController struct{}

func NewUploadController() *UploadController {
    return &UploadController{}
}

// POST /admin/api/upload
func (c *UploadController) Store(ctx http.Context) {
    file, err := ctx.Request().File("file")
    if err != nil {
        ctx.Response().Json(http.StatusBadRequest, http.Json{
            "error": "No file uploaded",
        })
        return
    }

    // Validate size (max 5MB)
    if file.GetSize() > 5*1024*1024 {
        ctx.Response().Json(http.StatusBadRequest, http.Json{
            "error": "File too large (max 5MB)",
        })
        return
    }

    // Validate extension
    ext := strings.ToLower(filepath.Ext(file.GetExtension()))
    allowed := map[string]bool{
        ".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
        ".webp": true, ".pdf": true,
    }
    if !allowed[ext] {
        ctx.Response().Json(http.StatusBadRequest, http.Json{
            "error": "File type not allowed",
        })
        return
    }

    // Generate random filename
    randomBytes := make([]byte, 16)
    rand.Read(randomBytes)
    filename := hex.EncodeToString(randomBytes) + ext

    // Store
    path, err := file.StoreAs("public/uploads", filename)
    if err != nil {
        ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "error": "Failed to store file",
        })
        return
    }

    ctx.Response().Json(http.StatusOK, http.Json{
        "path": path,
        "url":  "/storage/" + path,
    })
}
```

**Step 2: Register route**

Edit `backend/routes/api.go`:

```go
uploadController := admin.NewUploadController()
totpGroup.Post("/upload", uploadController.Store)
```

**Step 3: Verify**

```bash
cd backend
go build ./...
```

**Step 4: Commit**

```bash
git add backend/app/http/controllers/admin/upload_controller.go backend/routes/api.go
git commit -m "feat: add file upload controller with validation"
```

---

## Phase 6: Admin Vue SPA

### Task 16: Admin layout & navigation

**Files:**
- Create: `admin/src/layouts/AdminLayout.vue`
- Create: `admin/src/components/Sidebar.vue`
- Create: `admin/src/lib/api.ts`

**Step 1: Write API client**

Create `admin/src/lib/api.ts`:

```typescript
import axios from "axios";

const api = axios.create({
  baseURL: "/admin/api",
  withCredentials: true,
});

// Attach token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("access_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle 401
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("access_token");
      window.location.href = "/admin/login";
    }
    return Promise.reject(error);
  }
);

export default api;
```

**Step 2: Write admin layout**

Create `admin/src/layouts/AdminLayout.vue`:

```vue
<script setup lang="ts">
import { RouterView, RouterLink, useRouter } from "vue-router";
import { useColorMode } from "@vueuse/core";

const router = useRouter();
const mode = useColorMode();

const logout = () => {
  localStorage.removeItem("access_token");
  router.push("/admin/login");
};

const navItems = [
  { name: "Dashboard", path: "/admin", icon: "LayoutDashboard" },
  { name: "Blog Posts", path: "/admin/posts", icon: "FileText" },
  { name: "Projects", path: "/admin/projects", icon: "Briefcase" },
  { name: "Open Source", path: "/admin/open-source", icon: "Github" },
  { name: "Products", path: "/admin/products", icon: "Package" },
  { name: "Inbox", path: "/admin/inbox", icon: "Inbox" },
  { name: "Settings", path: "/admin/settings", icon: "Settings" },
];
</script>

<template>
  <div class="flex h-screen bg-background">
    <!-- Sidebar -->
    <aside class="w-64 border-r bg-muted/30 p-4">
      <div class="mb-8">
        <h1 class="text-xl font-bold">Okuru Admin</h1>
      </div>
      <nav class="space-y-1">
        <RouterLink
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors hover:bg-accent"
          active-class="bg-accent"
        >
          {{ item.name }}
        </RouterLink>
      </nav>
      <div class="absolute bottom-4 left-4 right-4">
        <button
          @click="logout"
          class="w-full rounded-lg px-3 py-2 text-sm text-muted-foreground hover:bg-accent"
        >
          Logout
        </button>
      </div>
    </aside>

    <!-- Main content -->
    <main class="flex-1 overflow-auto p-8">
      <RouterView />
    </main>
  </div>
</template>
```

**Step 3: Create placeholder views**

Create minimal placeholder untuk setiap view (`admin/src/views/Dashboard.vue`, `Posts.vue`, dll):

```vue
<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">Dashboard</h1>
    <p>Welcome back!</p>
  </div>
</template>
```

**Step 4: Add route guard**

Edit `admin/src/router/index.ts`:

```typescript
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem("access_token");
  if (to.meta.public) {
    next();
  } else if (!token) {
    next("/admin/login");
  } else {
    next();
  }
});
```

**Step 5: Verify build**

```bash
cd admin
pnpm run build
```

**Step 6: Commit**

```bash
git add admin/src/
git commit -m "feat: add admin layout, sidebar, API client, route guards"
```

---

### Task 17: Login page with TOTP

**Files:**
- Create: `admin/src/views/Login.vue`

**Step 1: Write login view**

Create `admin/src/views/Login.vue`:

```vue
<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import api from "@/lib/api";
import { toast } from "vue-sonner";

const router = useRouter();

const step = ref<"login" | "totp" | "totp-setup">("login");
const email = ref("");
const password = ref("");
const tempToken = ref("");
const totpCode = ref("");
const loading = ref(false);

const handleLogin = async () => {
  loading.value = true;
  try {
    const res = await api.post("/auth/login", {
      email: email.value,
      password: password.value,
    });

    if (res.data.requires_totp) {
      tempToken.value = res.data.temp_token;
      step.value = "totp";
    } else if (res.data.totp_setup_required) {
      localStorage.setItem("access_token", res.data.access_token);
      await setupTotp();
    } else {
      localStorage.setItem("access_token", res.data.access_token);
      router.push("/admin");
    }
  } catch (err: any) {
    toast.error(err.response?.data?.error || "Login failed");
  } finally {
    loading.value = false;
  }
};

const handleTotp = async () => {
  loading.value = true;
  try {
    const res = await api.post("/auth/totp", {
      temp_token: tempToken.value,
      code: totpCode.value,
    });
    localStorage.setItem("access_token", res.data.access_token);
    router.push("/admin");
  } catch (err: any) {
    toast.error(err.response?.data?.error || "Invalid code");
  } finally {
    loading.value = false;
  }
};

const setupTotp = async () => {
  const res = await api.post("/auth/totp/setup");
  step.value = "totp-setup";
  qrUrl.value = res.data.qr_url;
  secret.value = res.data.secret;
};

const verifySetup = async () => {
  loading.value = true;
  try {
    await api.post("/auth/totp/verify-setup", { code: totpCode.value });
    toast.success("TOTP verified!");
    router.push("/admin");
  } catch (err: any) {
    toast.error(err.response?.data?.error || "Invalid code");
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-muted/30">
    <Card class="w-full max-w-md">
      <CardHeader>
        <CardTitle>Okuru Admin</CardTitle>
        <CardDescription v-if="step === 'login'">Sign in to your account</CardDescription>
        <CardDescription v-else-if="step === 'totp'">Enter your authenticator code</CardDescription>
        <CardDescription v-else>Scan QR & verify setup</CardDescription>
      </CardHeader>
      <CardContent>
        <!-- Login step -->
        <form v-if="step === 'login'" @submit.prevent="handleLogin" class="space-y-4">
          <div class="space-y-2">
            <Label for="email">Email</Label>
            <Input id="email" v-model="email" type="email" required />
          </div>
          <div class="space-y-2">
            <Label for="password">Password</Label>
            <Input id="password" v-model="password" type="password" required />
          </div>
          <Button type="submit" class="w-full" :disabled="loading">
            {{ loading ? "Signing in..." : "Sign In" }}
          </Button>
        </form>

        <!-- TOTP step -->
        <form v-else-if="step === 'totp'" @submit.prevent="handleTotp" class="space-y-4">
          <div class="space-y-2">
            <Label for="totp">Authentication Code</Label>
            <Input id="totp" v-model="totpCode" maxlength="6" placeholder="000000" required />
          </div>
          <Button type="submit" class="w-full" :disabled="loading">Verify</Button>
        </form>

        <!-- TOTP setup step -->
        <div v-else class="space-y-4">
          <img :src="qrUrl" alt="TOTP QR Code" class="mx-auto" v-if="qrUrl" />
          <p class="text-center text-sm text-muted-foreground">Secret: {{ secret }}</p>
          <form @submit.prevent="verifySetup" class="space-y-2">
            <Input v-model="totpCode" maxlength="6" placeholder="Enter code from app" required />
            <Button type="submit" class="w-full" :disabled="loading">Verify Setup</Button>
          </form>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
```

**Step 2: Verify build**

```bash
cd admin
pnpm run build
```

**Step 3: Commit**

```bash
git add admin/src/views/Login.vue
git commit -m "feat: add login page with TOTP two-phase flow"
```

---

### Task 18: CRUD views (Posts, Projects, etc)

**Files:**
- Modify: `admin/src/views/Posts.vue`
- Create: `admin/src/views/PostEditor.vue`
- Modify views lainnya

**Step 1: Write Posts list view**

Edit `admin/src/views/Posts.vue`:

```vue
<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import api from "@/lib/api";
import { toast } from "vue-sonner";

interface Post {
  ID: number;
  Slug: string;
  TitleEn: string;
  TitleId: string;
  Status: string;
  PublishedAt: string;
}

const router = useRouter();
const posts = ref<Post[]>([]);
const loading = ref(true);

const load = async () => {
  const res = await api.get("/posts");
  posts.value = res.data;
  loading.value = false;
};

const destroy = async (id: number) => {
  if (!confirm("Delete this post?")) return;
  await api.delete(`/posts/${id}`);
  toast.success("Post deleted");
  load();
};

onMounted(load);
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">Blog Posts</h1>
      <Button @click="router.push('/admin/posts/new')">New Post</Button>
    </div>

    <div v-if="loading" class="text-muted-foreground">Loading...</div>

    <Table v-else>
      <TableHeader>
        <TableRow>
          <TableHead>Title (EN)</TableHead>
          <TableHead>Title (ID)</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Published</TableHead>
          <TableHead>Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="post in posts" :key="post.ID">
          <TableCell>{{ post.TitleEn }}</TableCell>
          <TableCell>{{ post.TitleId }}</TableCell>
          <TableCell>
            <Badge :variant="post.Status === 'published' ? 'default' : 'secondary'">
              {{ post.Status }}
            </Badge>
          </TableCell>
          <TableCell>{{ post.PublishedAt }}</TableCell>
          <TableCell class="space-x-2">
            <Button size="sm" variant="ghost" @click="router.push(`/admin/posts/${post.ID}/edit`)">
              Edit
            </Button>
            <Button size="sm" variant="ghost" @click="destroy(post.ID)">Delete</Button>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  </div>
</template>
```

**Step 2: Write Post Editor**

Create `admin/src/views/PostEditor.vue`:

```vue
<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import api from "@/lib/api";
import { toast } from "vue-sonner";

const route = useRoute();
const router = useRouter();
const isEdit = computed(() => !!route.params.id);

const form = ref({
  slug: "",
  title_en: "",
  title_id: "",
  excerpt_en: "",
  excerpt_id: "",
  content_en: "",
  content_id: "",
  category: "",
  status: "draft",
  published_at: "",
});

const categories = ref([]);

const load = async () => {
  if (isEdit.value) {
    const res = await api.get(`/posts/${route.params.id}`);
    form.value = res.data;
  }
  const catRes = await api.get("/categories");
  categories.value = catRes.data;
};

const save = async () => {
  try {
    if (isEdit.value) {
      await api.put(`/posts/${route.params.id}`, form.value);
    } else {
      await api.post("/posts", form.value);
    }
    toast.success("Post saved");
    router.push("/admin/posts");
  } catch (err: any) {
    toast.error(err.response?.data?.error || "Save failed");
  }
};

onMounted(load);
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-6">{{ isEdit ? "Edit Post" : "New Post" }}</h1>

    <div class="space-y-6 max-w-3xl">
      <div class="grid grid-cols-2 gap-4">
        <div class="space-y-2">
          <Label>Title (EN)</Label>
          <Input v-model="form.title_en" />
        </div>
        <div class="space-y-2">
          <Label>Title (ID)</Label>
          <Input v-model="form.title_id" />
        </div>
      </div>

      <div class="space-y-2">
        <Label>Slug</Label>
        <Input v-model="form.slug" placeholder="my-blog-post" />
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div class="space-y-2">
          <Label>Excerpt (EN)</Label>
          <Textarea v-model="form.excerpt_en" rows="3" />
        </div>
        <div class="space-y-2">
          <Label>Excerpt (ID)</Label>
          <Textarea v-model="form.excerpt_id" rows="3" />
        </div>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div class="space-y-2">
          <Label>Content (EN) - Markdown</Label>
          <Textarea v-model="form.content_en" rows="15" class="font-mono" />
        </div>
        <div class="space-y-2">
          <Label>Content (ID) - Markdown</Label>
          <Textarea v-model="form.content_id" rows="15" class="font-mono" />
        </div>
      </div>

      <div class="grid grid-cols-3 gap-4">
        <div class="space-y-2">
          <Label>Category</Label>
          <Select v-model="form.category">
            <SelectOption v-for="cat in categories" :key="cat.ID" :value="cat.Slug">
              {{ cat.NameEn }}
            </SelectOption>
          </Select>
        </div>
        <div class="space-y-2">
          <Label>Status</Label>
          <Select v-model="form.status">
            <SelectOption value="draft">Draft</SelectOption>
            <SelectOption value="published">Published</SelectOption>
          </Select>
        </div>
        <div class="space-y-2">
          <Label>Published At</Label>
          <Input v-model="form.published_at" type="datetime-local" />
        </div>
      </div>

      <div class="flex gap-3">
        <Button @click="save">Save</Button>
        <Button variant="ghost" @click="router.push('/admin/posts')">Cancel</Button>
      </div>
    </div>
  </div>
</template>
```

**Step 3: Tulis views lainnya** (Projects, OpenSource, Products, Inbox, Settings, Dashboard) mengikuti pattern yang sama. Inbox = list messages dengan mark read/archive. Settings = form untuk edit hero, social links, dll.

**Step 4: Verify build**

```bash
cd admin
pnpm run build
```

**Step 5: Commit**

```bash
git add admin/src/views/
git commit -m "feat: add CRUD views for posts, projects, products, inbox, settings"
```

---

## Phase 7: Storefront (Go template + Alpine.js)

### Task 19: Storefront templates & routes

**Files:**
- Create: `backend/resources/views/layouts/shop.tmpl`
- Create: `backend/resources/views/shop/index.tmpl`
- Create: `backend/resources/views/shop/product.tmpl`
- Modify: `backend/routes/web.go`

**Step 1: Write shop layout**

Create `backend/resources/views/layouts/shop.tmpl`:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Okuru Shop</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
</head>
<body class="bg-gray-50 min-h-screen">
    <nav class="bg-white border-b">
        <div class="max-w-5xl mx-auto px-4 py-4 flex items-center justify-between">
            <a href="/" class="text-xl font-bold">Okuru Shop</a>
            <div class="flex gap-4">
                <a href="/" class="text-gray-600 hover:text-gray-900">Products</a>
                <a href="https://okuru.id" class="text-gray-600 hover:text-gray-900">Blog</a>
            </div>
        </div>
    </nav>
    <main class="max-w-5xl mx-auto px-4 py-8">
        {{template "content" .}}
    </main>
</body>
</html>
```

**Step 2: Write index template**

Create `backend/resources/views/shop/index.tmpl`:

```html
{{template "layouts/shop" .}}

{{define "content"}}
<h1 class="text-3xl font-bold mb-8">Digital Products</h1>

<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
    {{range .Products}}
    <div class="bg-white rounded-lg shadow p-6">
        {{if .Thumbnail}}
        <img src="{{.Thumbnail}}" alt="{{.Title}}" class="w-full h-40 object-cover rounded mb-4">
        {{end}}
        <h2 class="text-lg font-semibold mb-2">{{.Title}}</h2>
        <p class="text-gray-600 text-sm mb-4">{{.Description}}</p>
        <div class="flex items-center justify-between">
            <span class="text-xl font-bold">Rp {{formatRupiah .Price}}</span>
            <a href="/product/{{.Slug}}" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
                View
            </a>
        </div>
    </div>
    {{end}}
</div>
{{end}}
```

**Step 3: Write product detail template**

Create `backend/resources/views/shop/product.tmpl`:

```html
{{template "layouts/shop" .}}

{{define "content"}}
<div class="max-w-3xl mx-auto" x-data="{ loading: false }">
    <h1 class="text-3xl font-bold mb-4">{{.Product.Title}}</h1>
    <p class="text-2xl text-gray-900 mb-6">Rp {{formatRupiah .Product.Price}}</p>
    <div class="prose mb-8">{{.Product.Description}}</div>

    <form @submit.prevent="loading=true; $refs.form.submit()" action="/checkout/{{.Product.Slug}}" method="POST" x-ref="form">
        <div class="mb-4">
            <label class="block text-sm font-medium mb-1">Email</label>
            <input type="email" name="email" required class="w-full border rounded px-3 py-2">
        </div>
        <div class="mb-4">
            <label class="block text-sm font-medium mb-1">Name</label>
            <input type="text" name="name" required class="w-full border rounded px-3 py-2">
        </div>
        <button type="submit" :disabled="loading"
            class="bg-blue-600 text-white px-6 py-3 rounded hover:bg-blue-700 disabled:opacity-50">
            {{if loading}}Processing...{{else}}Buy Now{{end}}
        </button>
    </form>
</div>
{{end}}
```

**Step 4: Write shop controller**

Create `backend/app/http/controllers/shop_controller.go`:

```go
package controllers

import (
    "fmt"

    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"

    "okuru/app/models"
)

type ShopController struct{}

func NewShopController() *ShopController {
    return &ShopController{}
}

func (c *ShopController) Index(ctx http.Context) {
    var products []models.Product
    facades.DB().Where("status = ?", "active").Get(&products)

    ctx.Response().View("shop.index", http.Json{
        "Title":    "Products",
        "Products": products,
    })
}

func (c *ShopController) Product(ctx http.Context) {
    slug := ctx.Request().Route("slug")

    var product models.Product
    if err := facades.DB().Where("slug = ?", slug).
        Where("status = ?", "active").
        First(&product); err != nil {
        ctx.Response().String(http.StatusNotFound, "Product not found")
        return
    }

    ctx.Response().View("shop.product", http.Json{
        "Title":   product.Title,
        "Product": product,
    })
}

// Checkout & payment — Phase 2 (Sumopod integration)
func (c *ShopController) Checkout(ctx http.Context) {
    // TODO: implement when Sumopod API ready
    ctx.Response().String(http.StatusNotImplemented, "Payment integration coming soon")
}
```

**Step 5: Register shop routes**

Edit `backend/routes/web.go`:

```go
package routes

import (
    "github.com/goravel/framework/facades"

    "okuru/app/http/controllers"
)

func Web() {
    shop := controllers.NewShopController()

    facades.Route().Get("/", shop.Index)
    facades.Route().Get("/product/{slug}", shop.Product)
    facades.Route().Post("/checkout/{slug}", shop.Checkout)
}
```

**Step 6: Verify**

```bash
cd backend
go build ./...
go run .
```

Buka browser: `http://localhost:3000/` (storefront).

**Step 7: Commit**

```bash
git add backend/resources/views/ backend/app/http/controllers/shop_controller.go backend/routes/web.go
git commit -m "feat: add storefront with Go templates + Alpine.js"
```

---

## Phase 8: Docker Deployment

### Task 20: Dockerfile & docker-compose

**Files:**
- Create: `deploy/Dockerfile`
- Create: `deploy/docker-compose.yml`
- Create: `deploy/docker-compose.prod.yml`
- Create: `deploy/Caddyfile`
- Create: `deploy/.env.example`

**Step 1: Write Dockerfile**

Create `deploy/Dockerfile`:

```dockerfile
# ============== Stage 1: Build Admin SPA ==============
FROM node:20-alpine AS admin-builder

WORKDIR /build

RUN corepack enable && corepack prepare pnpm@latest --activate

COPY admin/package.json admin/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY admin/ ./
RUN pnpm run build

# ============== Stage 2: Build Go Backend ==============
FROM golang:1.23-alpine AS backend-builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

# Copy admin build output ke embed dir
COPY --from=admin-builder /build/dist ./public/admin/

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /okuru .

# ============== Stage 3: Production Runtime ==============
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata sqlite-libs wget && \
    addgroup -S okuru && adduser -S okuru -G okuru

WORKDIR /opt/okuru

COPY --from=backend-builder /okuru ./okuru
COPY backend/.env.example ./.env.example

RUN mkdir -p /var/lib/okuru /opt/okuru/storage/app/public /opt/okuru/storage/logs && \
    chown -R okuru:okuru /opt/okuru /var/lib/okuru

COPY scripts/backup.sh /opt/okuru/scripts/backup.sh
RUN chmod +x /opt/okuru/scripts/backup.sh

USER okuru

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./okuru"]
```

**Step 2: Write docker-compose untuk local dev**

Create `deploy/docker-compose.yml`:

```yaml
services:
  app:
    build:
      context: ..
      dockerfile: deploy/Dockerfile
    container_name: okuru-app
    restart: unless-stopped
    env_file: ../backend/.env
    volumes:
      - db-data:/var/lib/okuru
      - storage-data:/opt/okuru/storage/app/public
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
    networks:
      - okuru-net

volumes:
  db-data:
  storage-data:

networks:
  okuru-net:
    driver: bridge
```

**Step 3: Write docker-compose untuk production**

Create `deploy/docker-compose.prod.yml`:

```yaml
services:
  app:
    image: registry.gitlab.com/${CI_PROJECT_PATH:-kurob1993/okuru}:latest
    container_name: okuru-app
    restart: unless-stopped
    env_file: /opt/okuru/backend/.env
    volumes:
      - db-data:/var/lib/okuru
      - storage-data:/opt/okuru/storage/app/public
      - /opt/okuru/logs:/opt/okuru/storage/logs
      - /opt/okuru/backups:/var/backups/okuru
    ports:
      - "127.0.0.1:8080:8080"
    deploy:
      resources:
        limits:
          memory: 256M
        reservations:
          memory: 128M
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
    networks:
      - okuru-net

  caddy:
    image: caddy:2-alpine
    container_name: okuru-caddy
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy-data:/data
      - caddy-config:/config
    depends_on:
      app:
        condition: service_healthy
    deploy:
      resources:
        limits:
          memory: 64M
    networks:
      - okuru-net

volumes:
  db-data:
  storage-data:
  caddy-data:
  caddy-config:

networks:
  okuru-net:
    driver: bridge
```

**Step 4: Write Caddyfile**

Create `deploy/Caddyfile`:

```caddy
api.okuru.id {
    reverse_proxy app:8080

    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains"
        X-Content-Type-Options nosniff
        X-Frame-Options DENY
        Referrer-Policy strict-origin-when-cross-origin
    }
}

shop.okuru.id {
    reverse_proxy app:8080
    header Strict-Transport-Security "max-age=31536000"
}
```

**Step 5: Write backup script**

Create `scripts/backup.sh`:

```bash
#!/bin/sh
set -e

BACKUP_DIR=/var/backups/okuru
DB_PATH=/var/lib/okuru/okuru.db
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"

sqlite3 "$DB_PATH" ".backup '$BACKUP_DIR/okuru_$DATE.db'"
gzip "$BACKUP_DIR/okuru_$DATE.db"
find "$BACKUP_DIR" -name "okuru_*.db.gz" -mtime +7 -delete

echo "Backup complete: okuru_$DATE.db.gz"
```

**Step 6: Write .env.example**

Create `deploy/.env.example`:

```env
# Docker Compose
COMPOSE_PROJECT_NAME=okuru

# CI Project Path (untuk image tag)
CI_PROJECT_PATH=kurob1993/okuru
```

**Step 7: Test local build**

```bash
cd /home/kurob/code/okuru.id
docker compose -f deploy/docker-compose.yml build
docker compose -f deploy/docker-compose.yml up -d
curl http://localhost:8080/api/v1/content
docker compose -f deploy/docker-compose.yml down
```

Expected: API respond JSON.

**Step 8: Commit**

```bash
git add deploy/ scripts/
git commit -m "feat: add Dockerfile, docker-compose, Caddyfile for deployment"
```

---

## Phase 9: GitLab CI/CD

### Task 21: GitLab CI pipeline

**Files:**
- Create: `.gitlab-ci.yml`
- Create: `scripts/setup-vps.sh`

**Step 1: Write .gitlab-ci.yml**

Create `.gitlab-ci.yml`:

```yaml
stages:
  - lint
  - test
  - build
  - deploy

variables:
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: "/certs"
  IMAGE_NAME: $CI_REGISTRY_IMAGE
  IMAGE_TAG: $CI_COMMIT_SHORT_SHA
  IMAGE_LATEST: $CI_REGISTRY_IMAGE:latest

# === LINT ===

lint:go:
  stage: lint
  image: golangci/golangci-lint:v1.61-alpine
  script:
    - cd backend
    - golangci-lint run --timeout 5m ./...
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH

lint:vue:
  stage: lint
  image: node:20-alpine
  before_script:
    - corepack enable && corepack prepare pnpm@latest --activate
  script:
    - cd admin
    - pnpm install --frozen-lockfile
    - pnpm run lint
  cache:
    key: admin-deps-${CI_COMMIT_REF_SLUG}
    paths:
      - admin/node_modules/
    policy: pull-push
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH

# === BUILD ===

build:image:
  stage: build
  image: docker:24
  services:
    - docker:24-dind
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - docker build -t $IMAGE_NAME:$IMAGE_TAG -t $IMAGE_LATEST -f deploy/Dockerfile .
    - docker push $IMAGE_NAME:$IMAGE_TAG
    - docker push $IMAGE_LATEST
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
    - if: $CI_COMMIT_BRANCH == "develop"
    - if: $CI_COMMIT_TAG

# === DEPLOY ===

.deploy_template: &deploy_template
  image: alpine:3.20
  before_script:
    - apk add --no-cache openssh-client
    - eval $(ssh-agent -s)
    - echo "$VPS_SSH_KEY" | tr -d '\r' | ssh-add -
    - mkdir -p ~/.ssh && chmod 700 ~/.ssh
    - echo "$VPS_SSH_FINGERPRINT" >> ~/.ssh/known_hosts
  script:
    - |
      ssh $VPS_USER@$VPS_HOST << 'EOF'
        set -e
        cd /opt/okuru
        docker compose -f deploy/docker-compose.prod.yml pull app
        docker compose -f deploy/docker-compose.prod.yml up -d --remove-orphans app
        docker compose -f deploy/docker-compose.prod.yml exec -T app ./artisan migrate --force
        docker image prune -f
      EOF

deploy:staging:
  <<: *deploy_template
  stage: deploy
  environment:
    name: staging
    url: https://staging.okuru.id
  rules:
    - if: $CI_COMMIT_BRANCH == "develop"
      when: on_success

deploy:production:
  <<: *deploy_template
  stage: deploy
  environment:
    name: production
    url: https://api.okuru.id
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
      when: manual

# === SECURITY ===

include:
  - template: Security/SAST.gitlab-ci.yml
  - template: Security/Dependency-Scanning.gitlab-ci.yml
  - template: Security/Container-Scanning.gitlab-ci.yml

container_scanning:
  variables:
    CS_IMAGE: $IMAGE_NAME:$IMAGE_TAG
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
```

**Step 2: Write VPS setup script**

Create `scripts/setup-vps.sh`:

```bash
#!/bin/bash
set -e

VPS_DIR=/opt/okuru
REGISTRY=registry.gitlab.com/kurob1993/okuru

echo "=== Installing Docker ==="
curl -fsSL https://get.docker.com | sh
systemctl enable --now docker

echo "=== Creating directories ==="
mkdir -p $VPS_DIR/deploy $VPS_DIR/backend $VPS_DIR/logs $VPS_DIR/backups

echo "=== Login to GitLab Registry ==="
echo "$GITLAB_DEPLOY_TOKEN" | docker login $REGISTRY -u deploy-token --password-stdin

echo ""
echo "=== Setup complete ==="
echo "Next steps:"
echo "1. Copy deploy/docker-compose.prod.yml and deploy/Caddyfile to $VPS_DIR/deploy/"
echo "2. Create backend/.env with production secrets"
echo "3. Run: docker compose -f deploy/docker-compose.prod.yml up -d"
echo "4. Run: docker compose exec app ./artisan key:generate"
echo "5. Run: docker compose exec app ./artisan jwt:secret"
echo "6. Run: docker compose exec app ./artisan migrate --seed"
echo "7. Setup cron backup:"
echo "   (crontab -l 2>/dev/null; echo '0 3 * * * cd $VPS_DIR && docker compose -f deploy/docker-compose.prod.yml exec app /opt/okuru/scripts/backup.sh') | crontab -"
```

**Step 3: Commit**

```bash
git add .gitlab-ci.yml scripts/setup-vps.sh
git commit -m "feat: add GitLab CI/CD pipeline and VPS setup script"
```

---

## Phase 10: Final Integration

### Task 22: Health check endpoint

**Files:**
- Modify: `backend/routes/web.go`

**Step 1: Add health route**

Edit `backend/routes/web.go`:

```go
facades.Route().Get("/health", func(ctx http.Context) {
    ctx.Response().Json(http.StatusOK, http.Json{
        "status": "ok",
    })
})
```

**Step 2: Verify**

```bash
cd backend
go build ./...
go run .
curl http://localhost:3000/health
```

Expected: `{"status":"ok"}`

**Step 3: Commit**

```bash
git add backend/routes/web.go
git commit -m "feat: add health check endpoint"
```

---

### Task 23: README & AGENTS.md

**Files:**
- Create: `README.md`
- Create: `AGENTS.md`

**Step 1: Write README**

Create `README.md`:

```markdown
# Okuru.id Backend & Admin Panel

Backend API + admin panel untuk okuru.id personal website.

## Tech Stack

- **Backend:** Goravel (Go) + SQLite
- **Admin:** Vue 3 + shadcn-vue (Reka UI)
- **Storefront:** Go template + Alpine.js
- **Deployment:** Docker + GitLab CI

## Local Development

### Prerequisites

- Go 1.23+
- Node.js 20+ & pnpm
- Docker (untuk containerized dev)
- Air (Go live reload): `go install github.com/air-verse/air@latest`

### Run backend

```bash
cd backend
cp .env.example .env
./artisan key:generate
./artisan jwt:secret
./artisan migrate --seed
air
```

Backend: http://localhost:3000

### Run admin SPA

```bash
cd admin
pnpm install
pnpm run dev
```

Admin: http://localhost:5173/admin/

### Run via Docker

```bash
docker compose -f deploy/docker-compose.yml up -d --build
```

## Deployment

See `docs/plans/2026-07-02-admin-panel-design.md` for full deployment guide.

## Project Structure

```
okuru.id/
├── backend/      # Goravel API
├── admin/        # Vue admin SPA
├── deploy/       # Docker, Caddy
├── scripts/      # Setup & backup
└── docs/         # Design docs
```
```

**Step 2: Write AGENTS.md**

Create `AGENTS.md`:

```markdown
# AGENTS.md - Okuru.id Backend

## Deskripsi Proyek

Backend API + admin panel untuk okuru.id. Goravel (Go) + SQLite + Vue admin SPA + Docker.

## Build/Lint/Test Commands

### Backend (Goravel)
```bash
cd backend
go build ./...                    # Build check
golangci-lint run --timeout 5m    # Lint
go run .                          # Run dev server
./artisan migrate                 # Run migrations
./artisan db:seed                 # Seed database
./artisan make:model Name         # Generate model
./artisan make:migration name     # Generate migration
./artisan make:controller Name    # Generate controller
```

### Admin (Vue + shadcn-vue)
```bash
cd admin
pnpm install                      # Install deps
pnpm run dev                      # Dev server (HMR)
pnpm run build                    # Production build
pnpm run lint                     # ESLint
pnpm run type-check               # TypeScript check
pnpm dlx shadcn-vue@latest add [component]  # Add component
```

### Docker
```bash
docker compose -f deploy/docker-compose.yml build      # Build image
docker compose -f deploy/docker-compose.yml up -d      # Start
docker compose -f deploy/docker-compose.yml logs -f    # Logs
```

## Tech Stack

- Go 1.23 + Goravel v1.17
- SQLite (CGO)
- Vue 3 + Vite + shadcn-vue (Reka UI)
- Tailwind CSS v4
- Docker + Caddy
- GitLab CI/CD

## Code Style

- Go: standard gofmt + goimports
- Vue: Composition API + `<script setup lang="ts">`
- Bilingual content: field `*_en` dan `*_id`

## Language Preference

Selalu jawab dalam Bahasa Indonesia.

## Git Workflow

Jangan commit/push kecuali diminta eksplisit.
```

**Step 3: Commit**

```bash
git add README.md AGENTS.md
git commit -m "docs: add README and AGENTS.md"
```

---

## Task Summary

| # | Task | Phase |
|---|---|---|
| 1 | Init monorepo structure | Scaffolding |
| 2 | Scaffold Goravel backend | Scaffolding |
| 3 | Scaffold Vue admin SPA | Scaffolding |
| 4 | Database migrations | Database |
| 5 | ORM models | Database |
| 6 | Database seeder | Database |
| 7 | Auth guard & JWT | Auth |
| 8 | TOTP service | Auth |
| 9 | Auth controller | Auth |
| 10 | Auth middleware | Auth |
| 11 | Auth routes | Auth |
| 12 | Content controller (public API) | Public API |
| 13 | CORS configuration | Public API |
| 14 | CRUD controllers | Admin API |
| 15 | File upload controller | Admin API |
| 16 | Admin layout & navigation | Admin SPA |
| 17 | Login page with TOTP | Admin SPA |
| 18 | CRUD views | Admin SPA |
| 19 | Storefront templates & routes | Storefront |
| 20 | Dockerfile & docker-compose | Docker |
| 21 | GitLab CI pipeline | CI/CD |
| 22 | Health check endpoint | Integration |
| 23 | README & AGENTS.md | Integration |

## Out of Scope (Phase 2)

- Sumopod payment integration (webhook, checkout flow)
- Email delivery untuk digital products
- Download token verification
- Analytics dashboard
- Multi-admin / roles
```
