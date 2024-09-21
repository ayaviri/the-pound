package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func getDatabaseServerUrl() string {
	databaseUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("POSTGRES_DB"),
	)

	return databaseUrl
}

func EstablishConnection(dbPtr **sql.DB) error {
	var url string = getDatabaseServerUrl()
	*dbPtr, err = sql.Open("pgx", url)

	if err != nil {
		return err
	}

	return (*dbPtr).Ping()
}
