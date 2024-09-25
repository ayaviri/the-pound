package main

import (
	"encoding/json"
	"net/http"
	"the-pound/internal"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type UserLoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponseBody struct {
	Token string `json:"token"`
}

func Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b UserLoginRequestBody

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

		var credentialsMatch bool

		timer.WithTimer("ensuring username/password match", func() {
			credentialsMatch, err = xdb.DoCredentialsMatch(
				db,
				b.Username,
				b.Password,
			)
		})

		if err != nil {
			http.Error(
				w,
				"Could not check match of credentials",
				http.StatusInternalServerError,
			)
			return
		}

		if !credentialsMatch {
			http.Error(w, "Credentials did not match", http.StatusBadRequest)
			return
		}

		var userId string

		timer.WithTimer("getting user ID", func() {
			userId, err = xdb.GetUserId(db, b.Username)
		})

		if err != nil {
			http.Error(w, "Could not fetch user ID", http.StatusInternalServerError)
			return
		}

		var jwt internal.JWT

		timer.WithTimer("generating JWT", func() {
			jwt, err = internal.GenerateJWT(userId)
		})

		if err != nil {
			http.Error(
				w,
				"Could not generate signed JWT",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("writing JWT to database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				return xdb.WriteJWT(e, jwt)
			})
		})

		if err != nil {
			http.Error(
				w,
				"Could not write JWT to database",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("writing JWT to response body", func() {
			responseBody, err := json.Marshal(UserLoginResponseBody{Token: jwt.Token})

			if err != nil {
				return
			}

			_, err = w.Write(responseBody)
		})

		if err != nil {
			http.Error(
				w,
				"Could not write JWT to response body",
				http.StatusInternalServerError,
			)
		}
	})
}
