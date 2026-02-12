package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var authService *AuthService

// In-memory stores for OAuth state and one-time exchange codes
var (
	oauthExchangeCodes = make(map[string]*oauthExchangeEntry)
	oauthStates        = make(map[string]time.Time) // state → expiry
	oauthMu            sync.Mutex
)

func InitAuth(router *gin.Engine) {
	LoadConfig()
	InitGoogleOAuth()
	authService = NewAuthService()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", handleSignUp)
		auth.POST("/sign-in", handleSignIn)
		auth.POST("/sign-out", handleSignOut)
		auth.POST("/refresh", handleRefresh)
		auth.GET("/me", AuthMiddleware(), handleGetCurrentUser)
		auth.PUT("/me", AuthMiddleware(), handleUpdateProfile)
		auth.POST("/send-verification", AuthMiddleware(), handleSendVerification)
		auth.GET("/verify-email", handleVerifyEmail)
		auth.GET("/oauth/google", handleGoogleOAuthURL)
		auth.GET("/oauth/google/callback", handleGoogleOAuthCallback)
		auth.POST("/oauth/google/exchange", handleOAuthExchange)
	}
}

func handleSignUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := authService.SignUp(c.Request.Context(), req)
	if err != nil {
		if err == ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		println("Sign up error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	// Send verification email automatically
	userID, err := uuid.Parse(user.ID.String())
	if err == nil {
		err = authService.SendVerificationEmail(c.Request.Context(), userID)
		if err != nil {
			// Log the error but don't fail the sign-up
			println("Failed to send verification email:", err.Error())
		}
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func handleSignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := authService.SignIn(c.Request.Context(), req)
	if err != nil {
		if err == ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign in"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func handleSignOut(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&req); err == nil && req.RefreshToken != "" {
		authService.SignOut(c.Request.Context(), req.RefreshToken)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}

func handleRefresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	accessToken, newRefreshToken, err := authService.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": newRefreshToken,
	})
}

func handleGetCurrentUser(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := authService.GetUserByID(c.Request.Context(), uid)
	if err != nil {
		if err == ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func handleUpdateProfile(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := authService.UpdateProfile(c.Request.Context(), uid, req.Name, req.Image, req.Region, req.City)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func handleSendVerification(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = authService.SendVerificationEmail(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent"})
}

func handleVerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	_, err := authService.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Redirect to frontend email verified success page
	// The frontend will fetch updated user data and update the session
	c.Redirect(http.StatusFound, AppConfig.FrontendURL+"/email-verified")
}

func handleGoogleOAuthURL(c *gin.Context) {
	// Generate random state
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Store state in memory with 10-minute expiry
	oauthMu.Lock()
	oauthStates[state] = time.Now().Add(10 * time.Minute)
	oauthMu.Unlock()

	authURL := GetGoogleAuthURL(state)
	c.JSON(http.StatusOK, OAuthURLResponse{AuthURL: authURL})
}

func handleGoogleOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	// Verify state from memory
	oauthMu.Lock()
	expiry, exists := oauthStates[state]
	if exists {
		delete(oauthStates, state) // one-time use
	}
	oauthMu.Unlock()

	if !exists || time.Now().After(expiry) {
		c.Redirect(http.StatusFound, AppConfig.FrontendURL+"/sign-in?error=oauth_failed")
		return
	}

	if code == "" {
		c.Redirect(http.StatusFound, AppConfig.FrontendURL+"/sign-in?error=oauth_failed")
		return
	}

	user, accessToken, refreshToken, err := authService.HandleGoogleOAuth(c.Request.Context(), code)
	if err != nil {
		c.Redirect(http.StatusFound, AppConfig.FrontendURL+"/sign-in?error=oauth_failed")
		return
	}

	// Generate a one-time exchange code so the frontend can retrieve user data + tokens
	exchangeCode := uuid.New().String()
	oauthMu.Lock()
	oauthExchangeCodes[exchangeCode] = &oauthExchangeEntry{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(2 * time.Minute),
	}
	oauthMu.Unlock()

	// Redirect to frontend callback with the exchange code
	c.Redirect(http.StatusFound, AppConfig.FrontendURL+"/auth/google/callback?auth_code="+exchangeCode)
}

func handleOAuthExchange(c *gin.Context) {
	var req OAuthExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code is required"})
		return
	}

	oauthMu.Lock()
	entry, exists := oauthExchangeCodes[req.Code]
	if exists {
		delete(oauthExchangeCodes, req.Code) // one-time use
	}
	oauthMu.Unlock()

	if !exists || time.Now().After(entry.ExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         entry.User,
		"accessToken":  entry.AccessToken,
		"refreshToken": entry.RefreshToken,
	})
}
