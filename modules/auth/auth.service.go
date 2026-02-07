package auth

import (
	"context"
	"errors"
	"fmt"
	"restorapp/db"
	"restorapp/db/client"
	"restorapp/modules/email"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	queries *client.Queries
}

func NewAuthService() *AuthService {
	return &AuthService{
		queries: db.Queries,
	}
}

func userToResponse(user client.User) *UserResponse {
	return &UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		EmailVerified: user.EmailVerified.Bool,
		Image:         user.Image.String,
		CreatedAt:     user.CreatedAt.Time,
	}
}

func (s *AuthService) SignUp(ctx context.Context, req SignUpRequest) (*UserResponse, error) {
	// Check if user exists
	existingUser, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser.Email != "" {
		return nil, ErrUserExists
	}
	// If error is anything other than "not found", we can ignore it and try to create

	// Hash password
	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user, err := s.queries.CreateUser(ctx, client.CreateUserParams{
		Email:         req.Email,
		PasswordHash:  pgtype.Text{String: passwordHash, Valid: true},
		Name:          req.Name,
		EmailVerified: pgtype.Bool{Bool: false, Valid: true},
		Image:         pgtype.Text{Valid: false},
	})
	if err != nil {
		return nil, err
	}

	return userToResponse(user), nil
}

func (s *AuthService) SignIn(ctx context.Context, req SignInRequest) (*UserResponse, string, string, error) {
	// Get user
	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	// Check password
	if !user.PasswordHash.Valid || !CheckPassword(req.Password, user.PasswordHash.String) {
		return nil, "", "", ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken := GenerateRefreshToken()
	tokenHash := HashToken(refreshToken)
	expiresAt := time.Now().Add(time.Duration(AppConfig.RefreshTokenDuration) * 24 * time.Hour)

	_, err = s.queries.CreateRefreshToken(ctx, client.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, "", "", err
	}

	return userToResponse(user), accessToken, refreshToken, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
	tokenHash := HashToken(refreshToken)

	// Get refresh token from DB
	token, err := s.queries.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	// Get user
	user, err := s.queries.GetUserById(ctx, token.UserID)
	if err != nil {
		return "", "", err
	}

	// Generate new access token
	accessToken, err := GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", "", err
	}

	// Rotate refresh token
	err = s.queries.RevokeRefreshToken(ctx, tokenHash)
	if err != nil {
		return "", "", err
	}

	newRefreshToken := GenerateRefreshToken()
	newTokenHash := HashToken(newRefreshToken)
	expiresAt := time.Now().Add(time.Duration(AppConfig.RefreshTokenDuration) * 24 * time.Hour)

	_, err = s.queries.CreateRefreshToken(ctx, client.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: newTokenHash,
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *AuthService) SignOut(ctx context.Context, refreshToken string) error {
	tokenHash := HashToken(refreshToken)
	return s.queries.RevokeRefreshToken(ctx, tokenHash)
}

func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := s.queries.GetUserById(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return userToResponse(user), nil
}

func (s *AuthService) HandleGoogleOAuth(ctx context.Context, code string) (*UserResponse, string, string, error) {
	// Get user info from Google
	googleUser, oauthToken, err := GetGoogleUserInfo(ctx, code)
	if err != nil {
		return nil, "", "", err
	}

	// Check if OAuth account exists
	oauthAccount, err := s.queries.GetOAuthAccount(ctx, client.GetOAuthAccountParams{
		Provider:       "google",
		ProviderUserID: googleUser.ID,
	})

	var user client.User
	if err != nil {
		// Check if user with email exists
		existingUser, err := s.queries.GetUserByEmail(ctx, googleUser.Email)
		if err == nil {
			// User exists, link OAuth account
			user = existingUser
			// Update image if not set
			if !user.Image.Valid && googleUser.Picture != "" {
				s.queries.UpdateUserImage(ctx, client.UpdateUserImageParams{
					Image: pgtype.Text{String: googleUser.Picture, Valid: true},
					ID:    user.ID,
				})
				user.Image = pgtype.Text{String: googleUser.Picture, Valid: true}
			}
		} else {
			// Create new user
			user, err = s.queries.CreateUser(ctx, client.CreateUserParams{
				Email:         googleUser.Email,
				Name:          googleUser.Name,
				EmailVerified: pgtype.Bool{Bool: googleUser.VerifiedEmail, Valid: true},
				PasswordHash:  pgtype.Text{Valid: false}, // No password for OAuth users
				Image:         pgtype.Text{String: googleUser.Picture, Valid: googleUser.Picture != ""},
			})
			if err != nil {
				return nil, "", "", err
			}
		}

		// Create OAuth account
		expiresAt := pgtype.Timestamp{}
		if oauthToken.Expiry.IsZero() == false {
			expiresAt = pgtype.Timestamp{Time: oauthToken.Expiry, Valid: true}
		}

		_, err = s.queries.CreateOAuthAccount(ctx, client.CreateOAuthAccountParams{
			UserID:         user.ID,
			Provider:       "google",
			ProviderUserID: googleUser.ID,
			AccessToken:    pgtype.Text{String: oauthToken.AccessToken, Valid: true},
			RefreshToken:   pgtype.Text{String: oauthToken.RefreshToken, Valid: true},
			ExpiresAt:      expiresAt,
		})
		if err != nil {
			return nil, "", "", err
		}
	} else {
		// OAuth account exists, get user
		user, err = s.queries.GetUserById(ctx, oauthAccount.UserID)
		if err != nil {
			return nil, "", "", err
		}
	}

	// Generate tokens
	accessToken, err := GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken := GenerateRefreshToken()
	tokenHash := HashToken(refreshToken)
	expiresAt := time.Now().Add(time.Duration(AppConfig.RefreshTokenDuration) * 24 * time.Hour)

	_, err = s.queries.CreateRefreshToken(ctx, client.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, "", "", err
	}

	return userToResponse(user), accessToken, refreshToken, nil
}

func (s *AuthService) CreateVerificationToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err := s.queries.CreateVerificationToken(ctx, client.CreateVerificationTokenParams{
		UserID:    userID,
		Token:     token,
		Type:      "email_verification",
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	// Get verification token
	verificationToken, err := s.queries.GetVerificationToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// Update user
	err = s.queries.UpdateUserEmailVerified(ctx, verificationToken.UserID)
	if err != nil {
		return err
	}

	// Delete token
	err = s.queries.DeleteVerificationToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SendVerificationEmail(ctx context.Context, userID uuid.UUID) error {
	// Get user
	user, err := s.queries.GetUserById(ctx, userID)
	if err != nil {
		return err
	}

	if user.EmailVerified.Bool {
		return errors.New("email already verified")
	}

	// Create verification token
	token, err := s.CreateVerificationToken(ctx, userID)
	if err != nil {
		return err
	}

	// Send email
	verificationURL := fmt.Sprintf("%s/auth/verify-email?token=%s", AppConfig.BackendURL, token)

	err = email.SendVerificationEmail(user.Email, user.Name, verificationURL)
	if err != nil {
		fmt.Printf("Failed to send verification email: %v\n", err)
		fmt.Printf("Verification URL: %s\n", verificationURL)
		// Don't return error - the token was created, so the user can still verify manually
	}

	return nil
}
