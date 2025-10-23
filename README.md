# 🐦 Chirpy — A Minimal Twitter-Like REST API in Go

## Overview

**Chirpy** is a minimalist backend application written in Go that allows users to create and manage short posts called *chirps* — similar to tweets.  
It’s built to demonstrate best practices in Go web development, focusing on clean architecture, authentication, and database interactions.

---

## ✨ Features

- **User Management:** Register, log in, update, and upgrade users.  
- **JWT Authentication:** Access and refresh tokens for secure session handling.  
- **Chirp Management:** Create, read, and delete chirps.  
- **Admin Controls:** View and reset server metrics.  
- **Payment Webhooks:** Simulated user upgrade via Polka webhook events.  
- **Middleware Metrics:** Request counter for tracking endpoint hits.

---

## 🚀 Why Chirpy Matters

Chirpy isn’t just a demo — it’s a **compact, realistic backend system** that reflects how professional Go services are built.  
It teaches you how to handle:
- API design and route organization
- Error handling and HTTP status codes
- Secure authentication using JWT and refresh tokens
- Database migrations and SQLC integration
- Environment-based configurations and API key handling

If you’re learning **Go for backend development**, Chirpy is a perfect project to understand structure, composition, and maintainability.

---

## ⚙️ Project Structure

```
.
├── handler_admin.go         # Admin-related routes (metrics, reset)
├── handler_auth.go          # Authentication logic (login, refresh, revoke)
├── handler_chirps.go        # Chirp CRUD handlers
├── handler_payments.go      # Payment/webhook-related endpoints
├── handler_users.go         # User registration and update logic
├── internal/
│   ├── auth/                # JWT, refresh tokens, API key, password hashing
│   └── database/            # SQLC-generated files and models
├── sql/                     # SQL queries and migration schemas
├── utils.go                 # Helper functions (JSON responses, error handling)
└── main.go                  # Application entry point
```

---

## 🧭 Routes

| Method | Route | Description |
|--------|--------|-------------|
| `GET` | `/api/healthz` | Health check endpoint |
| `POST` | `/api/users` | Register new user |
| `POST` | `/api/login` | Login user and receive tokens |
| `POST` | `/api/refresh` | Generate new access token via refresh token |
| `POST` | `/api/revoke` | Revoke refresh token |
| `PUT` | `/api/users` | Update user information |
| `GET` | `/api/chirps` | Get all chirps (supports sort and author_id params) |
| `GET` | `/api/chirps/{chirpID}` | Get a single chirp |
| `POST` | `/api/chirps` | Create a new chirp |
| `DELETE` | `/api/chirps/{chirpID}` | Delete a chirp |
| `POST` | `/api/polka/webhooks` | Upgrade user via Polka webhook |
| `GET` | `/admin/metrics` | View server metrics |
| `POST` | `/admin/reset` | Reset metrics |

---

## 🛠️ Installation

### 1. Clone the project
```bash
git clone https://github.com/yourusername/chirpy.git
cd chirpy
```

### 2. Install dependencies
```bash
go mod tidy
```

### 3. Setup your PostgreSQL database
Create a database:
```bash
createdb chirpy
```

Run migrations (using [Goose](https://github.com/pressly/goose)):
```bash
goose -dir ./sql/schema postgres "postgres://user:password@localhost:5432/chirpy?sslmode=disable" up
```

### 4. Generate SQLC code
```bash
sqlc generate
```

### 5. Start the server
```bash
go run main.go
```
Visit the API at: **http://localhost:8080**

---

## 🔐 Authentication Flow

1. User registers via `POST /api/users`
2. Logs in via `POST /api/login` → Receives `access_token` + `refresh_token`
3. Uses the `access_token` in `Authorization: Bearer <token>` header
4. When the access token expires, calls `POST /api/refresh`
5. Can revoke tokens via `POST /api/revoke`

---

## 🧩 Example Request

### Create a Chirp
```bash
curl -X POST http://localhost:8080/api/chirps   -H "Authorization: Bearer <access_token>"   -H "Content-Type: application/json"   -d '{"body": "Hello from Chirpy!"}'
```

---

## 🧠 Learning Takeaways

- How to use `net/http` idiomatically
- How to organize routes and handlers
- Using JWT securely with expiration and refresh logic
- Structuring Go apps with `internal/` and `sqlc`
- Writing database queries and migrations cleanly

---
---

## 🐘 Tech Stack

- **Language:** Go 1.24+  
- **Database:** PostgreSQL  
- **ORM/Query Layer:** SQLC  
- **Auth:** JWT & Refresh Tokens  
- **Migrations:** Goose


