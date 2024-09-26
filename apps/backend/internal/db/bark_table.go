package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Bark struct {
	DogId        string    `json:"dog_id"`
	Bark         string    `json:"bark"`
	CreationDate time.Time `json:"created_at"`
}

func WriteBark(e DBExecutor, bark string, userId string) error {
	id := uuid.NewString()
	statement := `insert into bark (id, dog_id, bark) values($1, $2, $3)`
	_, err = e.Exec(statement, id, userId, bark)

	return err
}

func GetDogBarks(e DBExecutor, dogId string, count uint, offset uint) ([]Bark, error) {
	query := `select bark, created_at from bark where dog_id = $1 order by created_at 
desc limit $2 offset $3`
	var rows *sql.Rows
	rows, err = e.Query(query, dogId, count, offset)

	if err != nil {
		return []Bark{}, err
	}

	var barks []Bark

	for rows.Next() {
		b := Bark{DogId: dogId}
		err = rows.Scan(&b.Bark, &b.CreationDate)

		if err != nil {
			return []Bark{}, err
		}

		barks = append(barks, b)
	}

	return barks, nil
}
