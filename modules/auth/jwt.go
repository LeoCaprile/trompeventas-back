package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

func GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	exp := now.Add(time.Duration(AppConfig.AccessTokenDuration) * time.Minute)

	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Iat:    now.Unix(),
		Exp:    exp.Unix(),
	}

	return generateToken(claims)
}

func generateToken(claims JWTClaims) (string, error) {
	// Create header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Create payload
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Create signature
	message := headerEncoded + "." + payloadEncoded
	signature := createSignature(message, AppConfig.JWTSecret)

	// Combine all parts
	token := message + "." + signature

	return token, nil
}

func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	headerEncoded := parts[0]
	payloadEncoded := parts[1]
	signature := parts[2]

	// Verify signature
	message := headerEncoded + "." + payloadEncoded
	expectedSignature := createSignature(message, AppConfig.JWTSecret)

	if signature != expectedSignature {
		return nil, ErrInvalidToken
	}

	// Decode payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	// Check expiration
	if time.Now().Unix() > claims.Exp {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

func createSignature(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func GenerateRefreshToken() string {
	return uuid.New().String()
}

func HashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return fmt.Sprintf("%x", h.Sum(nil))
}
