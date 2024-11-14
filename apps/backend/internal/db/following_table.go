package db

import (
	"database/sql"

	"github.com/google/uuid"
)

type Following struct {
	FromDogId string
	ToDogId   string
	// NOTE: Only to be considered valid within the lifespan of the transaction
	// from which this was queried
	IsApproved bool
}

// A more client-friendly representation of the following relationship between two
// dogs
type FollowResult struct {
	FollowRequestExists bool `json:"follow_request_exists"`
	// True if the follow request exists and has been approved. False otherwise
	IsApproved bool `json:"is_approved"`
}

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

func RemoveFollow(
	e DBExecutor,
	fromDogId string,
	toDogId string,
) error {
	statement := `delete from following where from_dog_id = $1 and to_dog_id = $2`
	_, err = e.Exec(statement, fromDogId, toDogId)

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

	return DoesFollowingExist(e, fromDogId, toDogId)
}

func DoesFollowingExist(e DBExecutor, fromDogId string, toDogId string) (bool, error) {
	query := `select count(*) from following where from_dog_id = $1 and to_dog_id = $2 
and is_approved = true`

	return QueryExists(e, query, fromDogId, toDogId)
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

func GetFollowing(e DBExecutor, fromDogId string, toDogId string) (Following, error) {
	query := `select from_dog_id, to_dog_id, is_approved from following where 
from_dog_id = $1 and to_dog_id = $2`
	var row *sql.Row = e.QueryRow(query, fromDogId, toDogId)
	var f Following
	err = row.Scan(&f.FromDogId, &f.ToDogId, &f.IsApproved)

	return f, err
}

// Wraps GetFollowing in a more client readable form. Determines whether a follow
// request exists between the two dogs and if that request has been approved
func GetFollowResult(
	e DBExecutor,
	fromDogId string,
	toDogId string,
) (FollowResult, error) {
	following, err := GetFollowing(e, fromDogId, toDogId)

	if err != nil && err != sql.ErrNoRows {
		return FollowResult{}, err
	}

	followRequestExists := err != sql.ErrNoRows
	r := FollowResult{
		FollowRequestExists: followRequestExists,
		IsApproved:          false,
	}

	if followRequestExists {
		r.IsApproved = following.IsApproved
	}

	return r, nil
}

// Scans the given result set into a list of string IDs
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
