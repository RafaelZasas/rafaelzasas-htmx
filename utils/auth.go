package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = os.Getenv("JWT_KEY")

// GenerateUID generates a unique identifier using uuid
func GenerateUID() (string, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("uuid.NewV7: %w", err)
	}
	return uid.String(), nil
}

// HashPassword generates a bcrypt hash of the password using a default cost.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plaintext password with a bcrypt hashed password.
// It returns true if the password matches the hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWTokens generates a pair of JWT tokens (access and refresh) for a given user ID.
func GenerateJWTokens(uid string) (accessTokenString, refreshTokenString string, err error) {
	// Access token expires after 5 minutes
	expirationTimeAccessToken := time.Now().Add(5 * time.Minute)
	// Refresh token expires after 7 days
	expirationTimeRefreshToken := time.Now().Add(24 * time.Hour * 7)

	// Access Token
	accessTokenClaims := &jwt.StandardClaims{
		ExpiresAt: expirationTimeAccessToken.Unix(),
		Issuer:    "htmx.proteatech.dev",
		Subject:   uid,
		IssuedAt:  time.Now().Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err = accessToken.SignedString([]byte(jwtKey))
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshTokenClaims := &jwt.StandardClaims{
		ExpiresAt: expirationTimeRefreshToken.Unix(),
		Issuer:    "htmx.proteatech.dev",
		Subject:   uid,
		// Ensure refresh token is not valid until 1 min before access token expires
		NotBefore: expirationTimeAccessToken.Add(-1 * time.Minute).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err = refreshToken.SignedString([]byte(jwtKey))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func ValidateJWT(r *http.Request) (uid string, err error) {

	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err // Token not found
	}

	tokenString := cookie.Value
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

	if err != nil {
		return "", fmt.Errorf("unable to parse JWT: %w", err) // Error parsing token
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}

func ValidateRefreshToken(r *http.Request) (uid string, err error) {

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", err // Token not found
	}

	tokenString := cookie.Value
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

	if err != nil {
		return "", fmt.Errorf("unable to parse JWT: %w", err) // Error parsing token
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return claims.Subject, nil
}

func AddAuthCookies(w *http.ResponseWriter, accessToken, refreshToken string) {
	// set access token in cookie
	http.SetCookie(*w, &http.Cookie{
		Name:     "token",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	})

	// set refresh token in cookie
	http.SetCookie(*w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	})
}
