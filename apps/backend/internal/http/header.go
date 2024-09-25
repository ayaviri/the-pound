package http

import (
	"errors"
	"net/http"
	"strings"
)

func GetAuthBearerToken(request *http.Request) (string, error) {
	var authHeader string = request.Header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("Authorization header required")
	}

	var bearerToken string = strings.TrimPrefix(authHeader, "Bearer ")

	if bearerToken == authHeader {
		return "", errors.New("Invalid token format")
	}

	return bearerToken, nil
}
