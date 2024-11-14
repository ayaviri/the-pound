package db

import (
	"database/sql"

	"github.com/google/uuid"
)

func WritePaw(e DBExecutor, originalBarkId, barkId, dogId string) error {
	id := uuid.NewString()
	statement := `insert into paw (id, original_bark_id, bark_id, dog_id)
    values($1, $2, $3, $4)`
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

func GetPawsToBark(e DBExecutor, barkId string) ([]Bark, error) {
	query := getPawsFromBarkQuery()
	var rows *sql.Rows
	rows, err = e.Query(query, barkId)

	if err != nil {
		return []Bark{}, err
	}

	return ConstructBarksFromRows(rows)
}

// TODO: How do we want to order these ?
func getPawsFromBarkQuery() string {
	query := `
    select b.id, b.dog_id, b.dog_username, b.bark, b.created_at,
    b.created_at as rebarked_at, b.treat_count, b.rebark_count, b.paw_count,
    'bark' as type
    from (
        select bark_id from paw where original_bark_id = $1
    ) as p join bark b on b.id = p.bark_id order by b.created_at desc
    `

	return query
}
