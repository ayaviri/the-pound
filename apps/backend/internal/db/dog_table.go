package db

import (
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
	statement := `insert into dog (id, username, password_hash) 
values($1, $2, $3);`
	_, err = e.Exec(statement, id, username, passwordHash)

	return err
}
