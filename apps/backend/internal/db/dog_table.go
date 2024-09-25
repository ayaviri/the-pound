package db

import (
	"database/sql"
	"the-pound/internal"

	"github.com/google/uuid"
)

func IsUsernameInUse(e DBExecutor, username string) (bool, error) {
	statement := "select count(*) from dog where username = $1"
	dogCount, err := QueryCount(e, statement, username)

	return dogCount == 1, err
}

func CreateDog(e DBExecutor, username string, password string) error {
	id := uuid.NewString()
	passwordHash, err := internal.HashString(password)

	if err != nil {
		return err
	}

	statement := `insert into dog (id, username, password_hash) 
values($1, $2, $3)`
	_, err = e.Exec(statement, id, username, passwordHash)

	return err
}

func DoCredentialsMatch(e DBExecutor, username string, password string) (bool, error) {
	passwordHash, err := internal.HashString(password)

	if err != nil {
		return false, err
	}

	statement := "select count(*) from dog where username = $1 and password_hash = $2;"
	c, err := QueryCount(e, statement, username, passwordHash)

	return c == 1, err
}

func GetUserId(e DBExecutor, username string) (string, error) {
	query := "select id from dog where username = $1"
	var row *sql.Row
	row = e.QueryRow(query, username)
	var userId string
	err = row.Scan(&userId)

	return userId, err
}
