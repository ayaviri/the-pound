package main

import (
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type PawRequestBody struct {
	BarkId  string `json:"bark_id"`
	Content string `json:"content"`
}

func Paw() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var dog xdb.Dog

		timer.WithTimer("getting dog info from Auth header JWT", func() {
			dog, err = xhttp.GetDogFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get dog information from Auth header JWT",
				http.StatusInternalServerError,
			)
			return
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
			return
		}

		timer.WithTimer("writing paw to database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				originalBarkId := b.BarkId
				barkId, err := xdb.WriteBark(e, b.Content, dog.Id, dog.Username)

				if err != nil {
					return err
				}

				err = xdb.WritePaw(e, originalBarkId, barkId, dog.Id)

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
