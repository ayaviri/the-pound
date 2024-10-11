package main

import (
	"encoding/json"
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type TimelineQueryStringParameters struct {
	Count  uint
	Offset uint
}

type TimelineResponseBody struct {
	Barks []xdb.Bark `json:"barks"`
}

func Timeline() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p TimelineQueryStringParameters

		timer.WithTimer("getting query string parameters", func() {
			p, err = parseTimelineQueryStringParameters(r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not read query string parameters",
				http.StatusBadRequest,
			)
			return
		}

		var dogId string

		timer.WithTimer("getting dog's ID from Auth header JWT", func() {
			dogId, err = xhttp.GetDogIdFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not dog's ID from Auth header JWT",
				http.StatusInternalServerError,
			)
			return
		}

		var barks []xdb.Bark

		timer.WithTimer("getting dog's timeline", func() {
			barks, err = xdb.GetDogTimeline(db, dogId, p.Count, p.Offset)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get dog's timeline",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("writing barks to response body", func() {
			responseBody, err := json.Marshal(TimelineResponseBody{Barks: barks})

			if err != nil {
				return
			}

			_, err = w.Write(responseBody)
		})

		if err != nil {
			http.Error(
				w,
				"Could not write timeline to response body",
				http.StatusInternalServerError,
			)
		}
	})
}

func parseTimelineQueryStringParameters(
	r *http.Request,
) (TimelineQueryStringParameters, error) {
	return TimelineQueryStringParameters{}, nil
}
