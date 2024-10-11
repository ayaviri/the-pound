package db

import "github.com/google/uuid"

func WriteRebark(e DBExecutor, barkId string, dogId string) error {
	id := uuid.NewString()
	statement := `insert into rebark (id, bark_id, dog_id) values($1, $2, $3)`
	_, err = e.Exec(statement, id, barkId, dogId)

	return err
}

func RemoveRebark(e DBExecutor, barkId string, dogId string) error {
	statement := `delete from rebark where bark_id = $1 and dog_id = $2`
	_, err = e.Exec(statement, barkId, dogId)

	return err
}

func HasDogGivenBarkRebark(e DBExecutor, barkId string, dogId string) (bool, error) {
	query := `select count(*) from rebark where bark_id = $1 and dog_id = $2`

	return QueryExists(e, query, barkId, dogId)
}

func GetBarkRebarkCount(e DBExecutor, barkId string) (uint, error) {
	query := `select count(*) from rebark where bark_id = $1`

	return QueryCount(e, query, barkId)
}
