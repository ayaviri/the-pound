package db

import (
	"database/sql"
	"errors"
	"the-pound/internal"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Session struct {
	Id             string
	DogId          string
	Token          string
	ExpirationDate time.Time
	CreationDate   time.Time
}

func WriteJWT(e DBExecutor, jwt internal.JWT) error {
	id := uuid.NewString()
	// Two weeks
	var expirationDate time.Time = jwt.CreationDate.Add(time.Hour * 24 * 7 * 2)
	statement := `insert into session (id, dog_id, token, expires_at, created_at)
values($1, $2, $3, $4, $5)`
	_, err = e.Exec(
		statement,
		id,
		jwt.Subject,
		jwt.Token,
		expirationDate,
		jwt.CreationDate,
	)

	return err
}

type JWTValidationResult struct {
	IsValid  bool
	NewToken string
}

// NOTE: IsValidJWT lives here instead of in the internal package to avoid
// circular imports since a DB read is required, I don't care to find a cleaner
// solution for this right now

// Checks if the given JWT:
// 1) Has the correct format
// 2) Has the correct signature
// 3) Has the correct signing method
// 4) Is not expired (according to the expiration date made by the token)
// 5) Is not expired (according to the database)
// JWT is valid if 1-3 and 5 are met. If 4 is met, no new token is generated. If 4
// is not met, a new token is generated
func IsValidJWT(e DBExecutor, tokenString string) (JWTValidationResult, error) {
	result := JWTValidationResult{IsValid: false, NewToken: ""}
	var t internal.JWT
	t, err = internal.ParseJWT(tokenString)
	isExpired := errors.Is(err, jwt.ErrTokenExpired)

	if err != nil && !isExpired {
		return result, err
	}

	if isExpired {
		t, err = internal.GenerateJWT(t.Subject)

		if err != nil {
			return result, errors.New("Token is expired, could not generate new one")
		}

		result.NewToken = t.Token
	}

	isExpired, err := IsJWTExpired(e, tokenString)

	if err != nil {
		return result, err
	}

	result.IsValid = !isExpired

	return result, nil
}

func IsJWTExpired(e DBExecutor, tokenString string) (bool, error) {
	query := `select expires_at from session where token = $1`
	var row *sql.Row
	row = e.QueryRow(query, tokenString)
	var expirationDate time.Time
	err = row.Scan(&expirationDate)

	if err != nil {
		return false, err
	}

	return time.Now().After(expirationDate), nil
}

func GetSessionByToken(e DBExecutor, tokenString string) (Session, error) {
	query := `select id, dog_id, token, expires_at, created_at from session 
where token = $1`
	var row *sql.Row
	row = e.QueryRow(query, tokenString)
	var s Session
	err = row.Scan(&s.Id, &s.DogId, &s.Token, &s.ExpirationDate, &s.CreationDate)

	if err != nil {
		return Session{}, err
	}

	return s, nil
}

func UpdateSessionToken(e DBExecutor, oldTokenString, newTokenString string) error {
	statement := `update session set token = $1 where token = $2`
	_, err = e.Exec(statement, newTokenString, oldTokenString)

	return err
}
