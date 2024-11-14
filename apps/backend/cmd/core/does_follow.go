package main

import (
	"encoding/json"
	"errors"
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type DoesFollowQueryStringParameters struct {
	DogId string
}

type DoesFollowResponseBody struct {
	xdb.FollowResult
}

func DoesFollow() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p DoesFollowQueryStringParameters

		timer.WithTimer("getting dog ID from query string", func() {
			p, err = parseDoesFollowQueryStringParameters(r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get dog ID from query string",
				http.StatusBadRequest,
			)
			return
		}

		var fromDogId string

		timer.WithTimer("getting dog ID from Auth header JWT", func() {
			fromDogId, err = xhttp.GetDogIdFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get dog ID from Auth header JWT",
				http.StatusInternalServerError,
			)
			return
		}

		var result xdb.FollowResult

		timer.WithTimer("checking to see if follow exists", func() {
			toDogId := p.DogId
			result, err = xdb.GetFollowResult(db, fromDogId, toDogId)
		})

		if err != nil {
			http.Error(
				w,
				"Could not check if dog is followed by requestor",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("writing to response body", func() {
			responseBody, err := json.Marshal(
				DoesFollowResponseBody{FollowResult: result},
			)

			if err != nil {
				return
			}

			_, err = w.Write(responseBody)
		})

		if err != nil {
			http.Error(
				w,
				"Could not write to response body",
				http.StatusInternalServerError,
			)
		}
	})
}

func parseDoesFollowQueryStringParameters(
	r *http.Request,
) (DoesFollowQueryStringParameters, error) {
	var p DoesFollowQueryStringParameters
	var ids []string = r.URL.Query()["id"]

	if ids == nil || ids[0] == "" {
		return p, errors.New("Dog ID has not been set")
	}

	p.DogId = ids[0]

	return p, nil
}
