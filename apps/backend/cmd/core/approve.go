package main

import (
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type ApproveRequestBody struct {
	DogId string `json:"dog_id"`
}

func Approve() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var toDogId string

		timer.WithTimer("getting dog ID from Auth header JWT", func() {
			toDogId, err = xhttp.GetDogIdFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get dog ID from Auth header JWT",
				http.StatusInternalServerError,
			)
			return
		}

		var b ApproveRequestBody

		timer.WithTimer(
			"getting ID of dog requesting to follow from request body",
			func() {
				err = xhttp.ReadUnmarshalRequestBody(r, &b)
			},
		)

		if err != nil {
			http.Error(
				w,
				"Could not ID of dog requesting to follow from request body",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("approving follow request in database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				fromDogId := b.DogId
				return xdb.ApproveFollowRequest(e, fromDogId, toDogId)
			})
		})

		if err != nil {
			http.Error(
				w,
				"Could not write follow request approval to database",
				http.StatusInternalServerError,
			)
			return
		}
	})
}
