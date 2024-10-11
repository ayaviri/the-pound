package http

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	xdb "the-pound/internal/db"
)

// Requires a JOIN between the session and dog tables, so prefer
// GetDogFromAuth unless more information on the dog is needed
func GetDogFromAuth(db *sql.DB, r *http.Request) (xdb.Dog, error) {
	var d xdb.Dog
	var jwt string
	jwt, err := GetAuthBearerToken(r)

	if err != nil {
		return d, err
	}

	d, err = xdb.GetDogByToken(db, jwt)

	if err != nil {
		return d, err
	}

	return d, nil
}

func GetDogIdFromAuth(db *sql.DB, r *http.Request) (string, error) {
	var jwt string
	jwt, err := GetAuthBearerToken(r)

	if err != nil {
		return "", err
	}

	var s xdb.Session
	s, err = xdb.GetSessionByToken(db, jwt)

	if err != nil {
		return "", err
	}

	return s.DogId, nil
}

func GetAuthBearerToken(r *http.Request) (string, error) {
	var authHeader string = r.Header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("Authorization header required")
	}

	var bearerToken string = strings.TrimPrefix(authHeader, "Bearer ")

	if bearerToken == authHeader {
		return "", errors.New("Invalid token format")
	}

	return bearerToken, nil
}
