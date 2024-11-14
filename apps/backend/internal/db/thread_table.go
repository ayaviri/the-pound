package db

import "database/sql"

func GetBarkThread(e DBExecutor, barkId string) ([]Bark, error) {
	query := getBarkThreadQuery()
	var rows *sql.Rows
	rows, err = e.Query(query, barkId)

	if err != nil {
		return []Bark{}, err
	}

	return ConstructBarksFromRows(rows)
}

func getBarkThreadQuery() string {
	query := `
    select id, dog_id, dog_username, bark, created_at, created_at as rebark_date,
    treat_count, rebark_count, paw_count, 'bark' as type
    from bark where id = $1 and not exists (select 1 from thread where bark_id = $1)

    union 

    select b.id, b.dog_id, b.dog_username, b.bark, b.created_at, b.created_at,
    b.treat_count, b.rebark_count, b.paw_count, 'bark' as type
    from (
        select unnest(string_to_array(thread_path::text, '.')) as bark_id
        from thread where bark_id = $1
    ) as t join bark b on t.bark_id = b.id

    order by created_at asc
    `

	return query
}
