package main

import (
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type ProtectRequestBody struct {
	Protected bool `json:"protected"`
}

func Protect() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var dogId string

		timer.WithTimer("getting dog ID from Auth header JWT", func() {
			dogId, err = xhttp.GetDogIdFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract dog ID from JWT",
				http.StatusInternalServerError,
			)
			return
		}

		var b ProtectRequestBody

		timer.WithTimer("unmarshalling body of request", func() {
			err = xhttp.ReadUnmarshalRequestBody(r, &b)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract profile visibility from request body",
				http.StatusBadRequest,
			)
			return
		}

		timer.WithTimer("writing updated profile visibilty to database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				return xdb.UpdateProfileVisibility(e, b.Protected, dogId)
			})
		})

		if err != nil {
			http.Error(
				w,
				"Could not update profile visibility in database",
				http.StatusInternalServerError,
			)
			return
		}
	})
}
