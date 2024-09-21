package main

import (
	"net/http"

	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type UserRegistrationRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b UserRegistrationRequestBody

		timer.WithTimer("unmarshalling body of request", func() {
			err = xhttp.ReadUnmarshalRequestBody(r, &b)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract user details from request body",
				http.StatusBadRequest,
			)
			return
		}

		var inUse bool

		timer.WithTimer("ensuring username uniqueness", func() {
			inUse, err = xdb.IsUsernameInUse(db, b.Username)
		})

		if err != nil {
			http.Error(
				w,
				"Could not determine if username in use",
				http.StatusInternalServerError,
			)
			return
		}

		if inUse {
			http.Error(w, "Username in use", http.StatusConflict)
			return
		}

		timer.WithTimer("creating dog in database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				return xdb.CreateDog(e, b.Username, b.Password)
			})
		})

		if err != nil {
			http.Error(w, "Could not create dog", http.StatusInternalServerError)
			return
		}
	})
}
