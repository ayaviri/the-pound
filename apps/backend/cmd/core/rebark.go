package main

import (
	"database/sql"
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type RebarkRequestBody struct {
	BarkId string `json:"bark_id"`
}

// Adds a rebark if the dog hasn't given one to the bark.
// Deletes the rebark otherwise
func Rebark() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var dogId string

		timer.WithTimer("getting dog ID from Auth header JWT", func() {
			dogId, err = xhttp.GetDogIdFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get dog ID from Auth header JWT",
				http.StatusInternalServerError,
			)
			return
		}

		var b RebarkRequestBody

		timer.WithTimer("unmarshalling body of request", func() {
			err = xhttp.ReadUnmarshalRequestBody(r, &b)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get bark ID from request body",
				http.StatusBadRequest,
			)
			return
		}

		timer.WithTimer("toggling rebark state for dog/bark pair in database", func() {
			err = toggleRebarkStateInDB(db, b.BarkId, dogId)
		})
	})
}

// In an atomic transaction
// 1) Checks whether a rebark exists for the given dog/bark pair
// 2) Writes/removes rebark accordingly
// 3) Increments/decrements bark's rebark count accordingly
func toggleRebarkStateInDB(db *sql.DB, barkId string, dogId string) error {
	return xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
		rebarkGiven, err := xdb.HasDogGivenBarkRebark(e, barkId, dogId)

		if err != nil {
			return err
		}

		if rebarkGiven {
			return removeRebarkFromDB(e, barkId, dogId)
		} else {
			return writeRebarkToDB(e, barkId, dogId)
		}
	})
}

func writeRebarkToDB(e xdb.DBExecutor, barkId string, dogId string) error {
	err = xdb.WriteRebark(e, barkId, dogId)

	if err != nil {
		return err
	}

	return xdb.IncrementRebarkCount(e, barkId)
}

func removeRebarkFromDB(e xdb.DBExecutor, barkId string, dogId string) error {
	err = xdb.RemoveRebark(e, barkId, dogId)

	if err != nil {
		return err
	}

	return xdb.DecrementRebarkCount(e, barkId)
}
