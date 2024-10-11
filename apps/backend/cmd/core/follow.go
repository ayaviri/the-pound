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

		timer.WithTimer("getting dog requesting follow from Auth header JWT", func() {
			fromDogId, err = xhttp.GetDogIdFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract ID of dog requesting follow from JWT",
				http.StatusInternalServerError,
			)
			return
		}

		var toDog xdb.Dog

		timer.WithTimer("getting dog to be followed from request body", func() {
			var b FollowRequestBody
			err = xhttp.ReadUnmarshalRequestBody(r, &b)

			if err != nil {
				return
			}

			toDog, err = xdb.GetDogFromId(db, b.DogId)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract dog to be followed from request body",
				http.StatusBadRequest,
			)
			return
		}

		timer.WithTimer("writing follow relationship to database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				isApproved := toDog.IsPublic
				return xdb.WriteFollow(e, fromDogId, toDog.Id, isApproved)
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
