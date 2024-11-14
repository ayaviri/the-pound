package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Bark struct {
	Id           string    `json:"id"`
	DogId        string    `json:"dog_id"`
	DogUsername  string    `json:"dog_username"`
	Bark         string    `json:"bark"`
	CreationDate time.Time `json:"created_at"`
	// TODO: This field is dependent on the dog that requests this bark
	RebarkDate time.Time `json:"rebark_date"`
	// TODO: This field doesn't make much sense outside of the context of a user's timeline
	// or their tweets
	IsRebark    bool `json:"is_rebark"`
	TreatCount  uint `json:"treat_count", omitempty`
	RebarkCount uint `json:"rebark_count", omitempty`
	PawCount    uint `json:"paw_count", omitempty`
}

func IncrementPawCount(e DBExecutor, barkId string) error {
	statement := `update bark set paw_count = paw_count + 1 where id = $1`
	_, err = e.Exec(statement, barkId)

	return err
}

func DecrementPawCount(e DBExecutor, barkId string) error {
	statement := `update bark set paw_count = paw_count - 1 where id = $1`
	_, err = e.Exec(statement, barkId)

	return err
}

func IncrementRebarkCount(e DBExecutor, barkId string) error {
	statement := `update bark set rebark_count = rebark_count + 1 where id = $1`
	_, err = e.Exec(statement, barkId)

	return err
}

func DecrementRebarkCount(e DBExecutor, barkId string) error {
	statement := `update bark set rebark_count = rebark_count - 1 where id = $1`
	_, err = e.Exec(statement, barkId)

	return err
}

func IncrementTreatCount(e DBExecutor, barkId string) error {
	statement := `update bark set treat_count = treat_count + 1 where id = $1`
	_, err = e.Exec(statement, barkId)

	return err
}

func DecrementTreatCount(e DBExecutor, barkId string) error {
	statement := `update bark set treat_count = treat_count - 1 where id = $1`
	_, err = e.Exec(statement, barkId)

	return err
}

func WriteBark(
	e DBExecutor,
	bark string,
	dogId string,
	dogUsername string,
) (string, error) {
	id := uuid.NewString()
	statement := `insert into bark (id, dog_id, bark, dog_username) values($1, $2, $3, $4)`
	_, err = e.Exec(statement, id, dogId, bark, dogUsername)

	return id, err
}

func RemoveBark(e DBExecutor, barkId string) error {
	statement := `delete from bark where id = $1`
	_, err = e.Exec(statement, barkId)

	return err
}

func GetBark(e DBExecutor, barkId string) (Bark, error) {
	// TODO: How do we know if this is a rebark ?
	query := `select id, dog_id, bark, created_at, treat_count, rebark_count, paw_count 
from bark where id = $1`
	var row *sql.Row = e.QueryRow(query, barkId)
	var b Bark
	err = row.Scan(
		&b.Id,
		&b.DogId,
		&b.Bark,
		&b.CreationDate,
		&b.TreatCount,
		&b.RebarkCount,
		&b.PawCount,
	)

	if err != nil {
		return b, err
	}

	b.IsRebark = false

	return b, nil
}

// Gets the barks from the dog with the given ID and all dogs it follows,
// constrained by the given pagination parameters
func GetDogTimeline(
	e DBExecutor,
	dogId string,
	count uint,
	offset uint,
) ([]Bark, error) {
	followingIds, err := GetFollowingIds(e, dogId)

	if err != nil {
		return []Bark{}, err
	}

	query := getDogTimelineQuery()
	var rows *sql.Rows
	rows, err = e.Query(query, append(followingIds, dogId), count, offset)

	if err != nil {
		return []Bark{}, err
	}

	return ConstructBarksFromRows(rows)
}

func GetDogBarks(e DBExecutor, dogId string, count uint, offset uint) ([]Bark, error) {
	query := getDogBarksQuery()
	var rows *sql.Rows
	rows, err = e.Query(query, dogId, count, offset)

	if err != nil {
		return []Bark{}, err
	}

	return ConstructBarksFromRows(rows)
}

func ConstructBarksFromRows(r *sql.Rows) ([]Bark, error) {
	barks := make([]Bark, 0)

	for r.Next() {
		b, err := constructBarkFromRow(r)

		if err != nil {
			return []Bark{}, err
		}

		barks = append(barks, b)
	}

	return barks, nil
}

// Scans the next row from the given set and constructs
// a bark from it
func constructBarkFromRow(r *sql.Rows) (Bark, error) {
	b := Bark{}
	var barkType string
	err = r.Scan(
		&b.Id,
		&b.DogId,
		&b.DogUsername,
		&b.Bark,
		&b.CreationDate,
		&b.RebarkDate,
		&b.TreatCount,
		&b.RebarkCount,
		&b.PawCount,
		&barkType,
	)

	if err != nil {
		return Bark{}, err
	}

	switch barkType {
	case "bark":
		b.IsRebark = false
	case "rebark":
		b.IsRebark = true
	default:
		return Bark{}, errors.New("Unknown bark type found")
	}

	return b, nil
}

//   ___  _   _ _____ ____  ___ _____ ____
//  / _ \| | | | ____|  _ \|_ _| ____/ ___|
// | | | | | | |  _| | |_) || ||  _| \___ \
// | |_| | |_| | |___|  _ < | || |___ ___) |
//  \__\_\\___/|_____|_| \_\___|_____|____/
//

func getDogBarksQuery() string {
	query := `
    with barks as (
        select 
        id, dog_id, dog_username, bark, created_at, created_at as rebarked_at,
        treat_count, rebark_count, paw_count, 'bark' as type 
        from bark where dog_id = $1 
        order by created_at desc limit $2::integer + $3::integer
    ),
    rebarks as (
        select 
        b.id, b.dog_id, b.dog_username, b.bark, b.created_at, r.created_at,
        b.treat_count, b.rebark_count, b.paw_count, 'rebark' as type
        from rebark as r join bark as b
        on r.bark_id = b.id where r.dog_id = $1
        order by r.created_at desc limit $2::integer + $3::integer
    ),
    combined as (
        select * from barks union select * from rebarks
    ),
    distinct_combined as (
        select distinct on (id) * from combined order by id, rebarked_at desc
    )
    select * from distinct_combined
    order by rebarked_at desc limit $2 offset $3
    `

	return query
}

// TODO: The union all in these two queries is causing barks to appear twice when
// they've been rebarked and recently posted. Figure out why the union all was
// here in the first place

func getDogTimelineQuery() string {
	query := `
    with barks as (
        select 
        id, dog_id, dog_username, bark, created_at, created_at as rebarked_at,
        treat_count, rebark_count, paw_count, 'bark' as type 
        from bark where dog_id = any($1) 
        order by created_at desc limit $2::integer + $3::integer
    ),
    rebarks as (
        select 
        b.id, b.dog_id, b.dog_username, b.bark, b.created_at, r.created_at,
        b.treat_count, b.rebark_count, b.paw_count, 'rebark' as type
        from rebark as r join bark as b
        on r.bark_id = b.id where r.dog_id = any($1)
        order by r.created_at desc limit $2::integer + $3::integer
    ),
    combined as (
        select * from barks union all select * from rebarks
    ),
    distinct_combined as (
        select distinct on (id) * from combined order by id, rebarked_at desc
    )
    select * from distinct_combined
    order by rebarked_at desc limit $2 offset $3
    `

	return query
}
