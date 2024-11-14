package main

import (
	"database/sql"
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type TreatRequestBody struct {
	BarkId string `json:"bark_id"`
}

// Adds a treat if the dog hasn't given one to the bark.
// Deletes the treat otherwise
func Treat() http.Handler {
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

		var b TreatRequestBody

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

		timer.WithTimer("toggling treat state for dog/bark pair in database", func() {
			err = toggleTreatStateInDB(db, b.BarkId, dogId)
		})

		if err != nil {
			http.Error(
				w,
				"Could not toggle treat's presence in database",
				http.StatusBadRequest,
			)
		}
	})
}

// In an atomic transaction
// 1) Checks whether a treat exists for the given dog/bark pair
// 2) Writes/removes treat accordingly
// 3) Increments/decrements bark's treat count accordingly
func toggleTreatStateInDB(db *sql.DB, barkId string, dogId string) error {
	return xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
		treatGiven, err := xdb.HasDogGivenBarkTreat(e, barkId, dogId)

		if err != nil {
			return err
		}

		if treatGiven {
			return removeTreatFromDB(e, barkId, dogId)
		} else {
			return writeTreatToDB(e, barkId, dogId)
		}
	})
}

func writeTreatToDB(e xdb.DBExecutor, barkId string, dogId string) error {
	err = xdb.WriteTreat(e, barkId, dogId)

	if err != nil {
		return err
	}

	return xdb.IncrementTreatCount(e, barkId)
}

func removeTreatFromDB(e xdb.DBExecutor, barkId string, dogId string) error {
	err = xdb.RemoveTreat(e, barkId, dogId)

	if err != nil {
		return err
	}

	return xdb.DecrementTreatCount(e, barkId)
}
