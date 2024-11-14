package main

import (
	"encoding/json"
	"errors"
	"net/http"
	xdb "the-pound/internal/db"

	"github.com/ayaviri/goutils/timer"
)

type PawsQueryStringParameters struct {
	BarkId string
}

type PawsResponseBody struct {
	Paws []xdb.Bark `json:"paws"`
}

func Paws() http.Handler {
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

		var paws []xdb.Bark

		timer.WithTimer("getting paws to bark", func() {
			paws, err = xdb.GetPawsToBark(db, p.BarkId)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get bark's paws",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("writing paws to response body", func() {
			responseBody, err := json.Marshal(PawsResponseBody{Paws: paws})

			if err != nil {
				return
			}

			_, err = w.Write(responseBody)
		})

		if err != nil {
			http.Error(
				w,
				"Could not write paws to response body",
				http.StatusInternalServerError,
			)
		}
	})
}

func parsePawsQueryStringParameters(
	r *http.Request,
) (PawsQueryStringParameters, error) {
	var p PawsQueryStringParameters
	var barkIds []string = r.URL.Query()["id"]

	if barkIds == nil || barkIds[0] == "" {
		return p, errors.New("Bark ID has not been set")
	}

	p.BarkId = barkIds[0]

	return p, nil
}
