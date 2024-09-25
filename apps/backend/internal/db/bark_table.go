package db

import "github.com/google/uuid"

func Bark(e DBExecutor, bark string, userId string) error {
	id := uuid.NewString()
	statement := `insert into bark (id, dog_id, bark) values($1, $2, $3)`
	_, err = e.Exec(statement, id, userId, bark)

	return err
}
