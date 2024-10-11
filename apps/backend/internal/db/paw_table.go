package db

import (
	"database/sql"

	"github.com/google/uuid"
)

func WritePaw(e DBExecutor, originalBarkId, barkId, dogId string) error {
	id := uuid.NewString()
	statement := `insert into paw (id, original_bark_id, bark_id, dog_id)`
	_, err = e.Exec(statement, id, originalBarkId, barkId, dogId)

	return err
}

func GetBarkPawCount(e DBExecutor, barkId string) (uint, error) {
	query := `select count(*) from paw where original_bark_id = $1`

	return QueryCount(e, query, barkId)
}

func IsBarkPaw(e DBExecutor, barkId string) (bool, error) {
	query := `select count(*) from paw where bark_id = $1`

	return QueryExists(e, query, barkId)
}

func GetOriginalBarkId(e DBExecutor, barkId string) (string, error) {
	var originalBarkId string
	query := `select original_bark_id from paw where bark_id = $1`
	var row *sql.Row = e.QueryRow(query, barkId)
	err = row.Scan(&originalBarkId)

	return originalBarkId, err
}
