package main

import (
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type BarkRequestBody struct {
	Content string `json:"content"`
}

func Bark() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b BarkRequestBody

		timer.WithTimer("unmarshalling body of request", func() {
			err = xhttp.ReadUnmarshalRequestBody(r, &b)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract bark from request body",
				http.StatusBadRequest,
			)
			return
		}

		var jwtString string

		timer.WithTimer("getting the JWT from the Auth header", func() {
			jwtString, err = xhttp.GetAuthBearerToken(r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get JWT from Auth header",
				http.StatusInternalServerError,
			)
			return
		}

		var userId string

		timer.WithTimer("getting the user ID from the JWT", func() {
			var s xdb.Session
			s, err = xdb.GetSessionByToken(db, jwtString)
			userId = s.DogId
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract user ID from JWT",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("writing bark to the database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				return xdb.Bark(e, b.Content, userId)
			})
		})

		if err != nil {
			http.Error(
				w,
				"Could not write bark to database",
				http.StatusInternalServerError,
			)
			return
		}
	})
}
