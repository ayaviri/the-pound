package db

import "github.com/google/uuid"

func WriteTreat(e DBExecutor, barkId string, dogId string) error {
	id := uuid.NewString()
	statement := `insert into treat (id, bark_id, dog_id) values($1, $2, $3)`
	_, err = e.Exec(statement, id, barkId, dogId)

	return err
}

func RemoveTreat(e DBExecutor, barkId string, dogId string) error {
	statement := `delete from treat where bark_id = $1 and dog_id = $2`
	_, err = e.Exec(statement, barkId, dogId)

	return err
}

func HasDogGivenBarkTreat(e DBExecutor, barkId string, dogId string) (bool, error) {
	query := `select count(*) from treat where bark_id = $1 and dog_id = $2`

	return QueryExists(e, query, barkId, dogId)
}

func GetBarkTreatCount(e DBExecutor, barkId string) (uint, error) {
	query := `select count(*) from treat where bark_id = $1`

	return QueryCount(e, query, barkId)
}
