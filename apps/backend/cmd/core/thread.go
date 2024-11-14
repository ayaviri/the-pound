package main

import (
	"encoding/json"
	"errors"
	"net/http"
	xdb "the-pound/internal/db"

	"github.com/ayaviri/goutils/timer"
)

type ThreadQueryStringParameters struct {
	BarkId string
}

type ThreadResponseBody struct {
	Barks []xdb.Bark `json:"barks"`
}

func Thread() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p ThreadQueryStringParameters

		timer.WithTimer("getting bark ID from query string", func() {
			p, err = parseThreadQueryStringParameters(r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not read query string parameters",
				http.StatusBadRequest,
			)
			return
		}

		var threadBarks []xdb.Bark

		timer.WithTimer("getting bark's thread", func() {
			threadBarks, err = xdb.GetBarkThread(db, p.BarkId)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get bark's thread",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("writing thread to response body", func() {
			responseBody, err := json.Marshal(ThreadResponseBody{Barks: threadBarks})

			if err != nil {
				return
			}

			_, err = w.Write(responseBody)
		})

		if err != nil {
			http.Error(
				w,
				"Could not write thread to response body",
				http.StatusInternalServerError,
			)
		}
	})
}

func parseThreadQueryStringParameters(
	r *http.Request,
) (ThreadQueryStringParameters, error) {
	var p ThreadQueryStringParameters
	var barkIds []string = r.URL.Query()["id"]

	if barkIds == nil || barkIds[0] == "" {
		return p, errors.New("Bark ID has not been set")
	}

	p.BarkId = barkIds[0]

	return p, nil
}
