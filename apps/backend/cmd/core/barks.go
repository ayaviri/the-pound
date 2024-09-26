package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	xdb "the-pound/internal/db"

	"github.com/ayaviri/goutils/timer"
)

type BarksQueryStringParameters struct {
	DogId  string
	Count  uint
	Offset uint
}

type BarksResponseBody struct {
	Barks []xdb.Bark `json:"barks"`
}

func Barks() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p BarksQueryStringParameters

		timer.WithTimer("getting query string parameters", func() {
			p, err = parseQueryStringParameters(r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not read query string parameters",
				http.StatusBadRequest,
			)
			return
		}

		// 2) Error if requested user's tweets cannot be viewed
		// by the user that made the request

		var barks []xdb.Bark

		timer.WithTimer("getting dog bark's", func() {
			barks, err = xdb.GetDogBarks(db, p.DogId, p.Count, p.Offset)
		})

		if err != nil {
			http.Error(w, "Could not get dog's barks", http.StatusInternalServerError)
			return
		}

		timer.WithTimer("writing barks to response body", func() {
			responseBody, err := json.Marshal(BarksResponseBody{Barks: barks})

			if err != nil {
				return
			}

			_, err = w.Write(responseBody)
		})

		if err != nil {
			http.Error(
				w,
				"Could not write barks to response body",
				http.StatusInternalServerError,
			)
		}
	})
}

//  _   _ _____ _     ____  _____ ____  ____
// | | | | ____| |   |  _ \| ____|  _ \/ ___|
// | |_| |  _| | |   | |_) |  _| | |_) \___ \
// |  _  | |___| |___|  __/| |___|  _ < ___) |
// |_| |_|_____|_____|_|   |_____|_| \_\____/
//

func parseQueryStringParameters(r *http.Request) (BarksQueryStringParameters, error) {
	var p BarksQueryStringParameters
	var dogIds []string = r.URL.Query()["dog_id"]

	if dogIds == nil || dogIds[0] == "" {
		return p, errors.New("Dog ID has not been set")
	}

	p.DogId = dogIds[0]

	var counts []string = r.URL.Query()["count"]

	if counts == nil || counts[0] == "" {
		return p, errors.New("Count has not been set")
	}

	count64, err := strconv.ParseUint(counts[0], 10, 64)

	if err != nil {
		return p, err
	}

	p.Count = uint(count64)
	var offsets []string = r.URL.Query()["offset"]

	if offsets == nil || offsets[0] == "" {
		return p, errors.New("Offset has not been set")
	}

	offset64, err := strconv.ParseUint(offsets[0], 10, 64)

	if err != nil {
		return p, err
	}

	p.Offset = uint(offset64)

	return p, nil
}
