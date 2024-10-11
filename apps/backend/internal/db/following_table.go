package db

import (
	"database/sql"

	"github.com/google/uuid"
)

func WriteFollow(
	e DBExecutor,
	fromDogId string,
	toDogId string,
	isApproved bool,
) error {
	id := uuid.NewString()
	statement := `insert into following (id, from_dog_id, to_dog_id, is_approved)
values($1, $2, $3, $4)`
	_, err = e.Exec(statement, id, fromDogId, toDogId, isApproved)

	return err
}

// Returns true if the dog associated with fromDogId can view the barks
// by the dog associated with toDogId
func CanBarksBeViewedByDog(
	e DBExecutor,
	fromDogId string,
	toDogId string,
) (bool, error) {
	if fromDogId == toDogId {
		return true, nil
	}

	isDogPublic, err := IsDogPublic(e, toDogId)

	if err != nil {
		return false, err
	}

	if isDogPublic {
		return true, nil
	}

	doesFollowingExist, err := DoesFollowingExist(e, fromDogId, toDogId)

	if err != nil {
		return false, err
	}

	return doesFollowingExist, nil
}

// Returns true if fromDogId follows toDogId and the following is approved
func DoesFollowingExist(e DBExecutor, fromDogId string, toDogId string) (bool, error) {
	query := `select count(*) from following where from_dog_id = $1 and to_dog_id = $2 
and is_approved = true`
	followingCount, err := QueryCount(e, query, fromDogId, toDogId)

	if err != nil {
		return false, err
	}

	return followingCount == 1, nil
}

func ApproveFollowRequest(e DBExecutor, fromDogId string, toDogId string) error {
	statement := `update following set is_approved = true where from_dog_id = $1
and to_dog_id = $2`
	_, err = e.Exec(statement, fromDogId, toDogId)

	return err
}

func GetFollowingIds(e DBExecutor, fromDogId string) ([]string, error) {
	query := `select to_dog_id from following where from_dog_id = $1`
	var rows *sql.Rows
	rows, err = e.Query(query, fromDogId)

	if err != nil {
		return []string{}, err
	}

	return constructIdsFromRows(rows)
}

func constructIdsFromRows(r *sql.Rows) ([]string, error) {
	var ids []string

	for r.Next() {
		id, err := constructIdFromRow(r)

		if err != nil {
			return []string{}, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func constructIdFromRow(r *sql.Rows) (string, error) {
	var id string
	err = r.Scan(&id)

	return id, err
}
