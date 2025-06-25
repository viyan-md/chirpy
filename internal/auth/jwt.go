package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, secret string) (string, error) {
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		Subject:   userID.String(),
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secret))
}

func ValidateJWT(tokenString, secret string) (uuid.UUID, error) {
	var claims jwt.RegisteredClaims

	_, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (any, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return uuid.Nil, err
	}
	if claims.Issuer != "chirpy" {
		return uuid.Nil, errors.New("unexpected issuer")
	}
	return uuid.Parse(claims.Subject)
}

func GetBearerToken(headers http.Header) (string, error) {
	auth_header := headers.Get("Authorization")

	if auth_header == "" {
		return "", errors.New("token doesn't exits")
	}

	return strings.TrimPrefix(auth_header, "Bearer "), nil
}

func MakeRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("reading random: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
