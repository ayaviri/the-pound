package db

import (
	"database/sql"
	"the-pound/internal"

	"github.com/google/uuid"
)

type Dog struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	IsPublic bool   `json:"is_public"`
}

func IsUsernameInUse(e DBExecutor, username string) (bool, error) {
	query := "select count(*) from dog where username = $1"
	dogCount, err := QueryCount(e, query, username)

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

func GetIdFromUsername(e DBExecutor, username string) (string, error) {
	query := "select id from dog where username = $1"
	var row *sql.Row = e.QueryRow(query, username)
	var userId string
	err = row.Scan(&userId)

	return userId, err
}

func UpdateProfileVisibility(e DBExecutor, isProtected bool, dogId string) error {
	statement := `update dog set is_public = $1 where id = $2`
	_, err = e.Exec(statement, !isProtected, dogId)

	return err
}

func IsDogPublic(e DBExecutor, dogId string) (bool, error) {
	query := `select is_public from dog where id = $1`
	var row *sql.Row = e.QueryRow(query, dogId)
	var isPublic bool
	err = row.Scan(&isPublic)

	return isPublic, err
}

func GetDogFromId(e DBExecutor, dogId string) (Dog, error) {
	query := `select id, username, is_public from dog where id = $1`
	var row *sql.Row = e.QueryRow(query, dogId)
	var d Dog
	err = row.Scan(&d.Id, &d.Username, &d.IsPublic)

	return d, err
}

func GetDogByUsername(e DBExecutor, username string) (Dog, error) {
	query := `select id, username, is_public from dog where username = $1`
	var row *sql.Row = e.QueryRow(query, username)
	var d Dog
	err = row.Scan(&d.Id, &d.Username, &d.IsPublic)

	return d, err
}

func IncrementFollowingCount(e DBExecutor, dogId string) error {
	statement := `update dog set following_count = following_count + 1 where id = $1`
	_, err = e.Exec(statement, dogId)
	return nil
}

func IncrementFollowerCount(e DBExecutor, dogId string) error {
	statement := `update dog set follower_count = follower_count + 1 where id = $1`
	_, err = e.Exec(statement, dogId)
	return nil
}

func DecrementFollowingCount(e DBExecutor, dogId string) error {
	statement := `update dog set following_count = following_count - 1 where id = $1`
	_, err = e.Exec(statement, dogId)
	return nil
}

func DecrementFollowerCount(e DBExecutor, dogId string) error {
	statement := `update dog set follower_count = follower_count - 1 where id = $1`
	_, err = e.Exec(statement, dogId)
	return nil
}
