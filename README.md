# Chirpy — Minimal Twitter-Like REST API in Go 

Chirpy is a compact backend application written in Go that demonstrates clean web API design, authentication with JWT + refresh tokens, database interaction via SQL (sqlc), and small but realistic features such as payment webhooks and admin metrics.

## Features
- User registration, login, update
- JWT access tokens + refresh tokens stored in DB
- Chirp CRUD: create, read, delete; simple profanity filter
- Admin endpoints (metrics, reset)
- Simulated payment webhook to upgrade users
- Request counting middleware

## Project layout & architecture
- `main.go` — app entrypoint and route registration
- `handler_*.go` — grouped HTTP handlers (auth, users, chirps, admin, payments)
- `utils.go` — JSON responses & error helpers
- `internal/auth` — authentication helpers (JWT, refresh token, passwords, api key)
- `internal/database` — sqlc generated queries & models (DB access layer)
- `sql/`— SQL schema and queries used by sqlc and migrations
- `index.html` — simple static file served under /app/
- `.env` — local environment config

## Prerequisites
- Go 1.24+
- PostgreSQL with pgcrypto extension (gen_random_uuid)
- sqlc (for generating internal/database types if you change queries)
- goose (or your chosen migration tool) for running migrations
- git

## Quickstart (local)
1. Clone:
    ```
   git clone https://github.com/deltron-fr/chirpy.git
   cd chirpy
    ```

2. Copy and edit .env (example values included in repo). Required vars:
   `DB_URL, PLATFORM, SECRET_KEY, POLKA_KEY`

3. Create DB:
    ```
   createdb chirpy
    ```

4. Run migrations (example using goose):
   ```
   goose -dir ./sql/schema postgres "postgres://user:pass@localhost:5432/chirpy" up
   ```

5. Generate sqlc code (if you modified sql files):
   ```
   sqlc generate
   ```

6. Fetch deps and run:
    ```
   go mod tidy
   go run main.go
    ```

By default server runs on :8080. Static files are available at /app/.

## Environment variables
- `DB_URL` — Postgres connection string, e.g. postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
- `PLATFORM` — "dev" or "prod" (admin reset is restricted to dev)
- `SECRET_KEY` — secret used to sign JWT access tokens
- `POLKA_KEY` — API key expected by payment webhook endpoint

## Database: schema, queries & migrations
- `sql/schema/*.sql` — Goose migrations; run them in order. Key migrations:
  - users table changes (hashed_password, is_chirpy_red)
  - chirps table with FK to users
  - refresh_tokens table for revocation/tracking
- `sql/queries/*.sql` — sqlc queries for DB operations (CreateUser, GetUser, CreateChirp, GetChirps, CreateRefreshToken, GetRefreshToken, UpdateRevokedAt, UpgradeUser, etc.)
- The repository uses `gen_random_uuid()` so ensure pgcrypto is enabled in your DB (CREATE EXTENSION IF NOT EXISTS pgcrypto;)

## API Reference (core)
Common headers:
- Content-Type: application/json
- Authorization: Bearer <access_token> (for protected endpoints)

#### 1) Health
- GET `/api/healthz`
  - `200: "OK"`

#### 2) Users
- POST `/api/users — Register`
    - Request:
        `{ "email": "me@example.com", "password": "pa$$" }`
    - Response: `201 user object (id, created_at, updated_at, email, is_chirpy_red)`

- POST `/api/login` — Login & receive tokens
    - Request:
        `{ "email": "me@example.com", "password": "pa$$" }`
    - Response: `200 includes token (JWT), refresh_token, and user metadata`

- PUT `/api/users` — Update user
    - Headers: `Authorization: Bearer <token>`
        - Request:
            `{ "email": "new@example.com", "password": "newpass" }`
    - Response: 200 updated user object (token is echoed)

#### 3) Tokens
- POST `/api/refresh` — Exchange refresh token for new access token
  - Headers: `Authorization: Bearer <refresh_token>`
  - Response: `200 { "token": "<new_jwt>" }`

- POST `/api/revoke` — Revoke refresh token
  - Headers: `Authorization: Bearer <refresh_token>`
  - Response: `204`

#### 4) Chirps
- POST `/api/chirps` — Create chirp
  - Headers: Authorization: Bearer <access_token>
  - Body: `{ "body": "Hello, world!" }`
  - Response: `201 created chirp (id, created_at, updated_at, body, user_id)`
  - Notes: body is trimmed; max length = 140; basic bad-word replacement applied.

- GET `/api/chirps` — List chirps
  - Query params:
    - `author_id (optional)`: UUID to filter by author
    - `sort (optional)`: asc | desc (default desc)
  - Response: `200` array of chirps

- GET `/api/chirps/{chirpID}` — Get single chirp
  - Response: `200` chirp object

- DELETE `/api/chirps/{chirpID}`
    - Headers: `Authorization: Bearer <access_token>`
    - Response: `204` on success

#### 5) Payments / Webhook
- POST `/api/polka/webhooks`
  - Headers: `Authorization: Bearer <POLKA_KEY>`
  - Body:
     `{ "event": "user.upgraded", "data": { "user_id": "<uuid>" } }`
  Response:
    - `204` if user upgraded
    - `401` if API key missing/invalid
  - Notes: This endpoint calls UpgradeUser (sets is_chirpy_red=true)

#### 6) Admin
- GET `/admin/metrics` — View request counter (HTML)
- POST `/admin/reset` — Reset server metrics and, in dev only, delete users (platform must be "dev")

## Authentication & tokens
- Access tokens: JWT (HS256) signed with SECRET_KEY. Short lived (handler uses 60m in places).
- Refresh tokens: random hex token (32 bytes) stored in refresh_tokens table with expires_at and revoked_at. Exchange via /api/refresh.
- Use Authorization: `Bearer <token>` for both access (JWT) and refresh token endpoints.

## Testing
- Unit tests exist for internal/auth package (jwt, password hashing, bearer parsing).
  Run:
    ```
    go test ./...
    ```

- If you alter SQL queries, regenerate sqlc code and run tests again.






