package main

import (
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type PawRequestBody struct {
	BarkId string `json:"bark_id"`
	Paw    string `json:"paw"`
}

func Paw() http.Handler {
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
		}

		var b PawRequestBody

		timer.WithTimer("unmarshalling body of request", func() {
			err = xhttp.ReadUnmarshalRequestBody(r, &b)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get bark ID from request body",
				http.StatusBadRequest,
			)
		}

		timer.WithTimer("writing paw to database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				originalBarkId := b.BarkId
				barkId, err := xdb.WriteBark(e, b.Paw, dogId)

				if err != nil {
					return err
				}

				err = xdb.WritePaw(e, originalBarkId, barkId, dogId)

				if err != nil {
					return err
				}

				return xdb.IncrementPawCount(e, originalBarkId)
			})
		})

		if err != nil {
			http.Error(
				w,
				"Could not write paw to database",
				http.StatusInternalServerError,
			)
		}
	})
}
