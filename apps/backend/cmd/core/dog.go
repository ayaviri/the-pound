package main

import (
	"encoding/json"
	"errors"
	"net/http"
	xdb "the-pound/internal/db"

	"github.com/ayaviri/goutils/timer"
)

type DogQueryStringParameters struct {
	DogUsername *string
	DogId       *string
}

type DogResponseBody struct {
	Dog xdb.Dog `json:"dog"`
}

func Dog() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p DogQueryStringParameters

		timer.WithTimer("getting username from query string", func() {
			p, err = parseDogQueryStringParameters(r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get username from query string",
				http.StatusBadRequest,
			)
			return
		}

		var dog xdb.Dog

		timer.WithTimer("searching for dog in database", func() {
			if p.DogId != nil {
				dog, err = xdb.GetDogFromId(db, *p.DogId)
			} else {
				dog, err = xdb.GetDogByUsername(db, *p.DogUsername)
			}
		})

		if err != nil {
			http.Error(
				w,
				"Could not search for dog with username",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("writing dog to response body", func() {
			responseBody, err := json.Marshal(
				DogResponseBody{Dog: dog},
			)

			if err != nil {
				return
			}

			_, err = w.Write(responseBody)
		})

		if err != nil {
			http.Error(
				w,
				"Could not write dog to response body",
				http.StatusInternalServerError,
			)
		}
	})
}

func parseDogQueryStringParameters(
	r *http.Request,
) (DogQueryStringParameters, error) {
	var p DogQueryStringParameters
	idSet := false
	usernameSet := false
	var ids []string = r.URL.Query()["id"]

	if ids != nil && ids[0] != "" {
		p.DogId = &ids[0]
		idSet = true
	}

	var usernames []string = r.URL.Query()["username"]

	if usernames != nil && usernames[0] != "" {
		p.DogUsername = &usernames[0]
		usernameSet = true
	}

	if idSet || usernameSet {
		return p, nil
	} else {
		return p, errors.New("Neither ID nor username set")
	}
}
