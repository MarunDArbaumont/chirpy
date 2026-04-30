package auth

import (
	"time"
	"fmt"
	"net/http"
	"strings"
	"crypto/rand"
	"encoding/hex"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		Subject: userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
	})

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error while creating auth token: %w", err)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	validateToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("token invalid or expired: %w", err)
	}

	id, err := validateToken.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while getting ID: %w", err)
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while parsing ID: %w", err)
	}

	return parsedID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization found in header")
	}

	parts := strings.Fields(authHeader)
    if len(parts) == 2 && parts[0] == "Bearer" {
        return parts[1], nil
    }
	return "", fmt.Errorf("invalid authorization header format")
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	encodedStr := hex.EncodeToString(key)
	return encodedStr
}