# Trompeventas - Backend API

RESTful API backend for Trompeventas, a Chilean marketplace platform. Built with Go, Gin framework, and PostgreSQL.

## 🚀 Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gin (HTTP web framework)
- **Database:** PostgreSQL
- **ORM:** sqlc (type-safe SQL code generation)
- **Authentication:** JWT tokens (access + refresh)
- **Email:** Resend API
- **OAuth:** Google OAuth 2.0
- **Migrations:** golang-migrate

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Resend API account (for email verification)
- Google OAuth credentials (for social login)

## 🛠️ Installation

```bash
# Clone the repository
git clone git@github.com:LeoCaprile/trompeventas-back.git
cd trompeventas-back

# Install dependencies
go mod download

# Install development tools
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## ⚙️ Environment Variables

Create a `.env` file in the root directory:

```env
# Database
DATABASE_URL=postgres://user:password@localhost:5432/trompeventas?sslmode=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
ACCESS_TOKEN_DURATION=15  # minutes
REFRESH_TOKEN_DURATION=10080  # 7 days in minutes

# Email (Resend API)
RESEND_API_KEY=re_your_resend_api_key
EMAIL_FROM=Trompeventas <no-reply@contacto.trompeventas.cl>

# Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/oauth/google/callback

# Frontend URL (for redirects)
FRONTEND_URL=http://localhost:5173

# Server
PORT=8080
```

## 🗄️ Database Setup

```bash
# Create database
createdb trompeventas

# Run migrations
migrate -path db/migrations -database "$DATABASE_URL" up

# Or run all migrations at once
make migrate-up

# Rollback migrations
make migrate-down
```

## 🚦 Development

```bash
# Run the server
go run main.go

# Run with hot reload (using air)
air

# Generate sqlc code after modifying queries
sqlc generate

# Create a new migration
migrate create -ext sql -dir db/migrations -seq migration_name
```

## 📁 Project Structure

```
main.go                 # Application entry point

modules/
├── auth/              # Authentication module
│   ├── auth.controller.go
│   ├── auth.service.go
│   ├── auth.middlewares.go
│   ├── jwt.go
│   ├── oauth.go
│   └── types.go
├── products/          # Products module
├── comments/          # Comments module
└── email/            # Email service

db/
├── migrations/        # SQL migration files
├── queries/          # SQL queries for sqlc
└── sqlc/             # Generated type-safe Go code

config/               # Configuration files
```

## 🔐 Authentication

The API uses JWT-based authentication with refresh tokens:

### Token Flow

1. **Sign In/Sign Up** → Returns `{ accessToken, refreshToken, user }`
2. **Access Token** - Short-lived (15 min), used for API requests
3. **Refresh Token** - Long-lived (7 days), stored in database
4. **Token Refresh** - Exchange refresh token for new access token
5. **Logout** - Invalidates refresh token in database

### Protected Endpoints

Use `Authorization: Bearer <access_token>` header for protected routes.

### Middleware

- `AuthMiddleware()` - Validates access token, extracts user ID
- `EmailVerifiedMiddleware()` - Checks if user's email is verified
- `OptionalAuthMiddleware()` - Doesn't require auth, but extracts user if present

## 📧 Email Verification

- Verification emails sent via Resend API
- Tokens stored in `verification_tokens` table with expiry
- Email verification required for publishing products
- Verification link: `GET /auth/verify-email?token=xxx`

## 🔌 API Endpoints

### Authentication

```
POST   /auth/sign-up                    # Create new account
POST   /auth/sign-in                    # Sign in with email/password
POST   /auth/sign-out                   # Sign out (invalidate refresh token)
POST   /auth/refresh                    # Refresh access token
GET    /auth/me                         # Get current user (protected)
PUT    /auth/me                         # Update user profile (protected)
POST   /auth/send-verification          # Send verification email (protected)
GET    /auth/verify-email?token=xxx     # Verify email
GET    /auth/oauth/google               # Get Google OAuth URL
GET    /auth/oauth/google/callback      # Google OAuth callback
POST   /auth/oauth/google/exchange      # Exchange auth code for tokens
```

### Products

```
GET    /products                        # List all products (optional auth)
GET    /products/:id                    # Get product details
POST   /products                        # Create product (protected, verified)
PUT    /products/:id                    # Update product (protected, owner only)
DELETE /products/:id                    # Delete product (protected, owner only)
GET    /products/user/:userId           # Get user's products
```

### Comments

```
GET    /products/:id/comments           # Get product comments
POST   /products/:id/comments           # Add comment (protected, verified)
DELETE /comments/:id                    # Delete comment (protected, owner only)
```

### Other

```
POST   /presign                         # Get presigned URL for S3 upload (protected)
```

## 🔍 Query Parameters

### Products List

- `?q=search` - Search by product name or description
- More filters coming soon (category, price range, location)

## 🗃️ Database Schema

### Key Tables

- `users` - User accounts
- `products` - Product listings
- `product_images` - Product images
- `product_categories` - Product category mappings
- `categories` - Available categories
- `comments` - Product comments
- `refresh_tokens` - Active refresh tokens
- `verification_tokens` - Email verification tokens

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific module tests
go test ./modules/auth/...

# Run with verbose output
go test -v ./...
```

## 🚀 Production Build

```bash
# Build binary
go build -o trompeventas-api

# Run binary
./trompeventas-api
```

## 🐳 Docker

```bash
# Build image
docker build -t trompeventas-api .

# Run container
docker run -p 8080:8080 --env-file .env trompeventas-api

# Using docker-compose
docker-compose up
```

## 📦 Key Dependencies

- `github.com/gin-gonic/gin` - HTTP framework
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/google/uuid` - UUID generation
- `github.com/joho/godotenv` - Environment variable loading
- `golang.org/x/oauth2` - OAuth 2.0 client
- `golang.org/x/crypto` - Password hashing (bcrypt)

## 🔒 Security Features

- ✅ Password hashing with bcrypt
- ✅ JWT signature verification
- ✅ Token expiration validation
- ✅ Refresh token rotation
- ✅ CORS configuration
- ✅ SQL injection prevention (sqlc + parameterized queries)
- ✅ Email verification requirement for sensitive actions
- ✅ Authorization checks (users can only modify their own resources)

## 🌐 CORS Configuration

CORS is configured to allow requests from the frontend:
- Allowed origins: `http://localhost:5173` (development)
- Allowed methods: GET, POST, PUT, DELETE, OPTIONS
- Allowed headers: Authorization, Content-Type

Update CORS settings in `main.go` for production.

## 📝 Development Workflow

### Adding a New Feature

1. **Database Changes:**
   ```bash
   # Create migration
   migrate create -ext sql -dir db/migrations -seq add_feature

   # Edit migration files
   # Run migration
   migrate -path db/migrations -database "$DATABASE_URL" up
   ```

2. **Add SQL Queries:**
   ```sql
   -- In db/queries/feature.sql
   -- name: GetFeature :one
   SELECT * FROM feature WHERE id = $1;
   ```

3. **Generate Code:**
   ```bash
   sqlc generate
   ```

4. **Implement Module:**
   - Create controller in `modules/feature/`
   - Add service logic
   - Register routes in `main.go`

5. **Test:**
   ```bash
   go test ./modules/feature/...
   ```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

[Add your license here]

## 🔗 Related

- Frontend: [trompeventas-front](https://github.com/LeoCaprile/trompeventas-front)
