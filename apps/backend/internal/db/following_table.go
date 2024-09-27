package db

import (
	"github.com/google/uuid"
)

func WriteFollow(e DBExecutor, fromDogId string, toDogId string) error {
	id := uuid.NewString()
	statement := "insert into following (id, from_dog_id, to_dog_id) values($1, $2, $3)"
	_, err = e.Exec(statement, id, fromDogId, toDogId)

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

// Returns true if fromDogId follows toDogId
func DoesFollowingExist(e DBExecutor, fromDogId string, toDogId string) (bool, error) {
	query := `select count(*) from following where from_dog_id = $1 and to_dog_id = $2`
	followingCount, err := QueryCount(e, query, fromDogId, toDogId)

	if err != nil {
		return false, err
	}

	return followingCount == 1, nil
}
