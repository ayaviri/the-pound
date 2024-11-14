package main

import (
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type NotificationReadRequestBody struct {
	NotificationId string `json:"notification_id"`
}

func NotificationRead() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b NotificationReadRequestBody

		timer.WithTimer(
			"getting ID of notification to be read from request body",
			func() {
				err = xhttp.ReadUnmarshalRequestBody(r, &b)
			},
		)

		if err != nil {
			http.Error(
				w,
				"Could not ID of notification to be read from request body",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("setting notification to read in database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				return xdb.SetNotificationToRead(e, b.NotificationId)
			})
		})

		if err != nil {
			http.Error(
				w,
				"Could not set notification to read in database",
				http.StatusInternalServerError,
			)
			return
		}
	})
}
