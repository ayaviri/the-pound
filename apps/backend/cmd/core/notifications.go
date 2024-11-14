package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type NotificationsQueryStringParameters struct {
	Count  uint
	Offset uint
}

type NotificationsResponseBody struct {
	Notifications []xdb.Notification `json:"notifications"`
}

func Notifications() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p NotificationsQueryStringParameters

		timer.WithTimer("parsing query string parameters", func() {
			p, err = parseNotificationsQueryStringParameters(r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not parse query string parameters",
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
				"Could not get dog ID from Auth header JWT",
				http.StatusInternalServerError,
			)
			return
		}

		var notifications []xdb.Notification

		timer.WithTimer("getting notifications", func() {
			notifications, err = xdb.GetUnreadNotifications(
				db,
				dogId,
				p.Count,
				p.Offset,
			)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get notifications",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("writing notifications to response body", func() {
			responseBody, err := json.Marshal(
				NotificationsResponseBody{Notifications: notifications},
			)

			if err != nil {
				return
			}

			_, err = w.Write(responseBody)
		})

		if err != nil {
			http.Error(
				w,
				"Could not write notifications to response body",
				http.StatusInternalServerError,
			)
			return
		}
	})
}

func parseNotificationsQueryStringParameters(
	r *http.Request,
) (NotificationsQueryStringParameters, error) {
	var p NotificationsQueryStringParameters
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
