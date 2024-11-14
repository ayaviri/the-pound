package main

import (
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type RejectRequestBody struct {
	DogId          string `json:"dog_id"`
	NotificationId string `json:"notification_id"`
}

func Reject() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var toDogId string

		timer.WithTimer("getting dog ID from Auth header JWT", func() {
			toDogId, err = xhttp.GetDogIdFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get dog ID from Auth header JWT",
				http.StatusInternalServerError,
			)
			return
		}

		var b RejectRequestBody

		timer.WithTimer(
			"getting ID of dog requesting to follow from request body",
			func() {
				err = xhttp.ReadUnmarshalRequestBody(r, &b)
			},
		)

		if err != nil {
			http.Error(
				w,
				"Could not ID of dog requesting to follow from request body",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("rejecting follow request in database", func() {
			err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
				fromDogId := b.DogId
				err = xdb.RemoveFollow(e, fromDogId, toDogId)

				if err != nil {
					return err
				}

				// TODO: My current design of the notification schema results
				// in the possibility of a read follow notification referring
				// to a follow relationship that no longer exists

				return xdb.SetNotificationToRead(e, b.NotificationId)
			})
		})

		if err != nil {
			http.Error(
				w,
				"Could not write follow request rejection io database",
				http.StatusInternalServerError,
			)
			return
		}
	})
}
