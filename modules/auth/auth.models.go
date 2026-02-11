package auth

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Response DTOs
type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	EmailVerified bool      `json:"emailVerified"`
	Image         string    `json:"image"`
	Region        string    `json:"region"`
	City          string    `json:"city"`
	CreatedAt     time.Time `json:"createdAt"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token,omitempty"` // Only used in non-cookie scenarios
}

type OAuthURLResponse struct {
	AuthURL string `json:"authUrl"`
}

type CSRFTokenResponse struct {
	CSRFToken string `json:"csrfToken"`
}

// JWT Claims
type JWTClaims struct {
	UserID uuid.UUID `json:"userId"`
	Email  string    `json:"email"`
	Exp    int64     `json:"exp"`
	Iat    int64     `json:"iat"`
}

type UpdateProfileRequest struct {
	Name   string  `json:"name" binding:"required"`
	Image  *string `json:"image"`
	Region *string `json:"region"`
	City   *string `json:"city"`
}

// OAuth exchange code request
type OAuthExchangeRequest struct {
	Code string `json:"code" binding:"required"`
}

// Temporary OAuth exchange code entry
type oauthExchangeEntry struct {
	User         *UserResponse
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}
