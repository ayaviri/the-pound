package main

import (
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type FollowRequestBody struct {
	DogId string `json:"dog_id"`
}

func Follow() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fromDogId string

		timer.WithTimer("getting dog ID from Auth header JWT", func() {
			fromDogId, err = xhttp.GetDogIdFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract dog ID from JWT",
				http.StatusInternalServerError,
			)
			return
		}

		var b FollowRequestBody

		timer.WithTimer("unmarshalling body of request", func() {
			err = xhttp.ReadUnmarshalRequestBody(r, &b)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract dog ID from request body",
				http.StatusBadRequest,
			)
			return
		}

		timer.WithTimer("writing follow relationship to database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				return xdb.WriteFollow(e, fromDogId, b.DogId)
			})
		})

		if err != nil {
			http.Error(
				w,
				"Could not write follow relationship to database",
				http.StatusInternalServerError,
			)
			return
		}
	})
}
