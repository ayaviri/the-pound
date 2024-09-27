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

		timer.WithTimer("writing bark to the database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				return xdb.WriteBark(e, b.Content, dogId)
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
