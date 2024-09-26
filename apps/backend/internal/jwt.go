package internal

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Token        string
	CreationDate time.Time
	// The ID of the user this token authenticates
	Subject string
}

func ParseJWT(tokenString string) (JWT, error) {
	var t *jwt.Token
	verificationCallback := func(parsedToken *jwt.Token) (any, error) {
		if parsedToken.Method != jwt.SigningMethodHS256 {
			return []byte(""), errors.New("Invalid signing method")
		}

		key, err := getJWTSymmetricKey()

		return []byte(key), err
	}
	t, err = jwt.Parse(tokenString, verificationCallback)

	if err != nil {
		return JWT{}, err
	}

	var creationDate *jwt.NumericDate
	creationDate, err = t.Claims.GetExpirationTime()

	if err != nil {
		return JWT{}, err
	}

	var subject string
	subject, err = t.Claims.GetSubject()

	if err != nil {
		return JWT{}, err
	}

	return JWT{
		Token:        tokenString,
		CreationDate: creationDate.Time,
		Subject:      subject,
	}, nil
}

func GenerateJWT(userId string) (JWT, error) {
	key, err := getJWTSymmetricKey()

	if err != nil {
		return JWT{}, err
	}

	var creationDate time.Time = time.Now()
	var expirationDate time.Time = creationDate.Add(time.Hour)
	claims := jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(expirationDate),
	}
	var t *jwt.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(key))

	return JWT{
		Token:        s,
		CreationDate: creationDate,
		Subject:      userId,
	}, err
}

func getJWTSymmetricKey() (string, error) {
	key, isPresent := os.LookupEnv("JWT_SECRET_KEY")

	if !isPresent {
		return "", errors.New("JWT secret key not set")
	}

	return key, nil
}
