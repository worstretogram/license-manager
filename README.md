# License Service

Minimal licensing service with a Go backend and Next.js frontend.  
Backend provides REST API endpoints for authentication and license management (generate / verify / list / download / CRUD).  
Frontend contains a simple admin UI for license operations.

## 📦 Contents

### **Backend** — Go API & Migrations  
Located in `backend/`

- **Auth handler:** `auth.Login`  
  → `backend/internal/auth/auth.go`
- **License handlers:**  
  `license.Generate`, `license.Verify`, `license.ListLicenses`,  
  `license.GetLicense`, `license.UpdateLicense`,  
  `license.DeleteLicense`, `license.DownloadLicense`  
  → `backend/internal/license/license.go`
- **JWT middleware:** `middleware.JWTMiddleware`  
  → `backend/internal/middleware/jwt.go`
- **Database init:** `db.InitDB`  
  → `backend/internal/db/db.go`
- **Migrations:**  
  → `backend/migrations/0001_init.up.sql`

### **Frontend** — Next.js Admin UI  
Located in `frontend/`

- **Auth flow:** `login` → stores token in `localStorage`  
  → `frontend/src/api/auth.ts`
- **License API calls:**  
  `getLicenses`, `createLicense`, `downloadLicense`, `deleteLicense`  
  → `frontend/src/api/licence.ts`
- **Global state:** MobX store  
  → `frontend/src/store/store.ts`
- **Pages:**  
  - `/` → `frontend/src/app/page.tsx`  
  - `/dashboard` → `frontend/src/app/dashboard/page.tsx`

### **Docker**
Located in `docker/`

- Compose: `docker/composes/docker-compose.local.yaml`
- Multi-stage Dockerfile: `docker/dockerfiles/Dockerfile.local`
- Nginx config: `docker/configs/nginx.local.conf`

---

## 🚀 Prerequisites

- **Docker & docker-compose** (recommended)
- **Node.js 18+**
- **Go 1.23+**
- **PostgreSQL** (if not using docker-compose)

---

## ⚡ Quick Start (Docker recommended)

From repository root:

```sh
docker compose -f docker/composes/docker-compose.local.yaml up --build
````

* Frontend → **[http://localhost:3000](http://localhost:3000)**
* Backend → **[http://localhost:8080](http://localhost:8080)**
  (When using nginx: API proxied at `/api`)

**Notes:**

* Environment variables defined in compose file.
* Dockerfile copies RSA keys + backend binary during build.

---

## 🧪 Running Locally (without Docker)

### **Backend**

Set environment variables:

```
DATABASE_URL
JWT_SECRET
RSA_PRIVATE_KEY_PATH
RSA_PUBLIC_KEY_PATH
```

Run migrations (`psql` or `migrate`) using:

```
backend/migrations/0001_init.up.sql
```

Start server:

```sh
cd backend
go run ./cmd/main.go
```

### **Frontend**

```sh
cd frontend
npm install
npm run dev
```

---

## 🔌 API Overview

### **Auth**

`POST /api/auth/login` — login
→ handler: `auth.Login`

### **License**

* `POST /api/license/generate` — generate license
* `POST /api/license/verify` — verify license

### **Admin (JWT required: Authorization: JWT <token>)**

* `GET /api/admin/licenses` — list
* `GET /api/admin/licenses/:id` — get
* `PUT /api/admin/licenses/:id` — update
* `DELETE /api/admin/licenses/:id` — delete
* `GET /api/admin/licenses/:id/download` — download license file

---

## 📝 Project Notes

* JWT auth via `middleware.JWTMiddleware`
* License format: **JSON → RSA-signed → base64**
* Frontend uses MobX store + simple UI components
* Swagger documentation available in `docs/`

---

## 🤝 Contributing

1. Create a branch
2. Run & test locally
3. Open a PR

---

