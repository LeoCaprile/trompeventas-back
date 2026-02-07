# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go REST API backend for an e-commerce store, built with Gin framework. It uses PostgreSQL for data persistence, sqlc for type-safe SQL queries, and GoBetterAuth for authentication. The API runs on port 8080 and serves a React frontend at `http://localhost:5173`.

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
├── auth/          # Authentication (GoBetterAuth integration)
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
- Migrations in `db/migrations/` using goose
- Migration naming: `NNN_description.sql` (e.g., `001_created_initial_tables.sql`)
- Migrations must include `-- +goose up` and `-- +goose down` directives

**Database client**:
- Connection pooling via `pgxpool`
- Global `db.Queries` variable provides access to all generated query methods
- Each query method is strongly typed based on SQL definitions

### Authentication (GoBetterAuth)

Authentication is handled by the `go-better-auth` library:

- **Configuration**: `modules/auth/auth.config.go` initializes auth with email/password strategy
- **Email verification**: Required on signup, sends verification email via Resend
- **Routes**: All auth routes exposed at `/auth/*` (handled by BetterAuth)
- **Middleware**: `auth.AuthMiddleware()` protects routes (checks session/token)
- **Templates**: Email templates in `modules/email/templates/` using gonja

**Auth endpoints** (managed by BetterAuth):
- `POST /auth/sign-up` - Create account
- `POST /auth/sign-in` - Login
- `POST /auth/sign-out` - Logout
- `GET /auth/verify-email` - Verify email token

### Routing Pattern

Routes are registered in module controllers, called from `main.go`:

```go
// In main.go
products.ProductsController(router)

// In modules/products/products.controller.go
func ProductsController(router *gin.Engine) {
    router.GET("/products", getProductsHandler)  // Public

    protected := router.Group("/products")
    protected.Use(auth.AuthMiddleware())         // Protected
    protected.POST("/", createProductHandler)
}
```

**Public routes**:
- `GET /products` - List all products (with images and categories)
- `GET /products/:id` - Get single product

**Protected routes** (require authentication):
- `POST /products` - Create product
- `POST /products/:id` - Update product
- `DELETE /products/:id` - Delete product

### CORS Configuration

CORS is configured in `main.go` to allow:
- **Origin**: `http://localhost:5173` (React frontend)
- **Methods**: GET, POST, PUT, DELETE, OPTIONS
- **Headers**: Origin, Content-Type, Cookie, Authorization
- **Credentials**: `true` (allows cookies)

### Environment Variables

Required environment variables (in `.env`):

- `DB_URL` - PostgreSQL connection string
- `GO_BETTER_AUTH_BASE_URL` - Base URL for auth callbacks (e.g., `http://localhost:8080`)
- `GO_BETTER_AUTH_SECRET` - Secret key for session encryption
- `RESEND_APIKEY` - Resend API key for sending emails

**Note**: The `.env` file is tracked in git with example values. In production, use actual secrets.

## Key Conventions

- **Error handling**: Log errors with `charmbracelet/log`, return appropriate HTTP status codes
- **UUID handling**: Use `google/uuid` package; all IDs are UUIDs
- **Handler naming**: `{action}{Resource}Handler` (e.g., `getProductsHandler`, `createProductHandler`)
- **JSON responses**: Use `ctx.JSON()` with appropriate status codes
- **Package naming**: Module packages named after domain (e.g., `package products`)

## Database Schema

Core tables:
- `products` - Product details (id, name, description, price)
- `product_images` - Product images (many-to-one with products)
- `categories` - Category master data
- `products_category` - Product-category junction table (many-to-many)

All tables include `created_at` and `updated_at` timestamps.

## Adding New Features

1. **Create migration**: Add SQL file in `db/migrations/` with goose directives
2. **Run migration**: `source goose.config.sh && goose up`
3. **Add queries**: Write sqlc queries in `db/queries/` (or update existing)
4. **Generate code**: Run `sqlc generate`
5. **Create module**: Add handlers in `modules/{domain}/` following existing patterns
6. **Register routes**: Call controller function from `main.go`
