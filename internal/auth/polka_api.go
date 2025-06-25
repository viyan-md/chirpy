package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("token doesn't exits")
	}

	return strings.TrimPrefix(authHeader, "ApiKey "), nil
}
