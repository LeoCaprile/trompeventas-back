# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Security

**NEVER commit sensitive information** such as API keys, secrets, tokens, passwords, or credentials to the repository. This includes:
- `.env` files (already in `.gitignore`)
- Hardcoded secrets in source code
- Private keys or certificates

Use environment variables for all sensitive configuration. If a secret is accidentally committed, it must be rotated immediately — removing it from future commits does NOT remove it from git history.

## Project Overview

This is a Go REST API backend for an e-commerce store, built with Gin framework. It uses PostgreSQL for data persistence, sqlc for type-safe SQL queries, and custom JWT authentication. The API runs on port 8080 and serves a React frontend at `http://localhost:5173`.

## Development Commands

- **Run the server**: `go run main.go` (starts on :8080)
- **Install dependencies**: `go mod download`
- **Format code**: `go fmt ./...`
- **Run database**: `docker-compose up -d` (starts PostgreSQL on port 5432)
- **Run migrations**: `source goose.config.sh && goose up`
- **Rollback migration**: `source goose.config.sh && goose down`
- **Generate SQL client code**: `sqlc generate` (regenerates `db/client/*` from SQL queries)

## Architecture

### Module-based Structure

The codebase follows a modular pattern where each domain has its own package:

```
modules/
├── auth/          # Custom JWT authentication
├── products/      # Product CRUD operations
├── categories/    # Category management
└── email/         # Email service (Resend integration)
```

Each module typically contains:
- `*.controller.go` - Route registration and handler mapping
- `*.service.go` - Handler functions (business logic)
- `*.models.go` - Domain-specific types (optional)

### Database Layer (sqlc Pattern)

The project uses **sqlc** for type-safe database access:

1. **Write SQL**: Define queries in `db/queries/*.sql` using sqlc annotations
2. **Generate code**: Run `sqlc generate` to create Go code in `db/client/`
3. **Use in handlers**: Access via global `db.Queries` (initialized in `db.InitDBClient()`)

**Schema management**:
- Migrations in `db/migrations/` using goose (timestamp-based naming)
- Migrations must include `-- +goose up` and `-- +goose down` directives

**Database client**:
- Connection pooling via `pgxpool`
- Global `db.Queries` variable provides access to all generated query methods
- Each query method is strongly typed based on SQL definitions

### Authentication (Custom JWT)

Authentication uses a custom JWT implementation with a single-cookie architecture:

- **Token strategy**: Short-lived access tokens (15 min) + refresh tokens (7 days)
- **Token delivery**: Tokens returned in JSON response body, stored by frontend in a single httpOnly `__session` cookie
- **Password hashing**: bcrypt with cost 12
- **Middleware**: `AuthMiddleware()` reads `Authorization: Bearer <token>` header
- **OAuth**: Google OAuth with in-memory state and one-time exchange codes
- **Email verification**: Sends verification emails via Resend

**Auth endpoints** (`modules/auth/auth.controller.go`):
- `POST /auth/sign-up` - Create account (returns `{ user }`)
- `POST /auth/sign-in` - Login (returns `{ user, accessToken, refreshToken }`)
- `POST /auth/sign-out` - Logout (accepts `{ refreshToken }` in body)
- `POST /auth/refresh` - Refresh tokens (accepts `{ refreshToken }`, returns new tokens)
- `GET /auth/me` - Get current user (protected, reads Authorization header)
- `GET /auth/oauth/google` - Get Google OAuth URL
- `GET /auth/oauth/google/callback` - Handle OAuth callback
- `POST /auth/oauth/google/exchange` - Exchange one-time code for user + tokens

### Routing Pattern

Routes are registered in module controllers, called from `main.go`:

```go
// In main.go
products.ProductsController(router)

// In modules/products/products.controller.go
func ProductsController(router *gin.Engine) {
    router.GET("/products", getProductsHandler)      // Public
    router.GET("/products/:id", getProductByIdHandler) // Public

    protected := router.Group("/products")
    protected.Use(auth.AuthMiddleware())
    protected.POST("/", createProductHandler)
    protected.GET("/me", getMyProductsHandler)
    protected.DELETE("/me/:id", deleteMyProductHandler)
    protected.POST("/me/:id", updateMyProductHandler)

    // Requires email verification
    publish := router.Group("/products")
    publish.Use(auth.AuthMiddleware(), auth.EmailVerifiedMiddleware())
    publish.POST("/publish", publishProductHandler)
}
```

**Public routes**:
- `GET /products` - List all products (with images and categories). Supports `?q=` query param for search by name/description (case-insensitive, partial match via `ILIKE`)
- `GET /products/:id` - Get single product (includes seller info)

**Protected routes** (require Authorization header):
- `POST /products` - Create product
- `GET /products/me` - List current user's products
- `POST /products/me/:id` - Update own product
- `DELETE /products/me/:id` - Delete own product
- `POST /products/publish` - Publish product (requires email verification)
- `POST /products/:id` - Update product (admin)
- `DELETE /products/:id` - Delete product (admin)

### CORS Configuration

CORS is configured in `main.go` to allow:
- **Origin**: `http://localhost:5173` (React frontend)
- **Methods**: GET, POST, PUT, DELETE, OPTIONS
- **Headers**: Origin, Content-Type, Cookie, Authorization
- **Credentials**: `true`

### Environment Variables

Required environment variables (in `.env`):

- `DB_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for signing JWT tokens
- `FRONTEND_URL` - Frontend URL (default: `http://localhost:5173`)
- `BACKEND_URL` - Backend URL (default: `http://localhost:8080`)
- `GOOGLE_CLIENT_ID` - Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret
- `RESEND_API_KEY` - Resend API key for sending emails

## Key Conventions

- **Error handling**: Log errors with `charmbracelet/log`, return appropriate HTTP status codes
- **UUID handling**: Use `google/uuid` package; all IDs are UUIDs
- **Handler naming**: `{action}{Resource}Handler` (e.g., `getProductsHandler`, `createProductHandler`)
- **JSON responses**: Use `ctx.JSON()` with appropriate status codes
- **Package naming**: Module packages named after domain (e.g., `package products`)

## Database Schema

Core tables:
- `users` - User accounts (id, email, password_hash, name, email_verified, image)
- `refresh_tokens` - JWT refresh tokens (hashed, with expiry and revocation)
- `verification_tokens` - Email verification tokens
- `oauth_accounts` - Linked OAuth providers (Google)
- `products` - Product details (id, name, description, price, user_id, condition, state, negotiable)
- `product_images` - Product images (many-to-one with products)
- `categories` - Category master data
- `products_category` - Product-category junction table (many-to-many)

All tables include `created_at` timestamps. Most include `updated_at`.

## Adding New Features

1. **Create migration**: Add SQL file in `db/migrations/` with goose directives (timestamp naming)
2. **Run migration**: `source goose.config.sh && goose up`
3. **Add queries**: Write sqlc queries in `db/queries/` (or update existing)
4. **Generate code**: Run `sqlc generate`
5. **Create module**: Add handlers in `modules/{domain}/` following existing patterns
6. **Register routes**: Call controller function from `main.go`
