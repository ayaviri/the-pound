package main

//
// import (
// 	"encoding/json"
// 	"net/http"
// 	xdb "the-pound/internal/db"
// 	xhttp "the-pound/internal/http"
//
// 	"github.com/ayaviri/goutils/timer"
// )
//
// type NotificationsResponseBody struct {
// 	Notifications []xdb.Notification `json:"notifications"`
// }
//
// func Notifications() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var dogId string
//
// 		timer.WithTimer("getting dog ID from Auth header JWT", func() {
// 			dogId, err = xhttp.GetDogIdFromAuth(db, r)
// 		})
//
// 		if err != nil {
// 			http.Error(
// 				w,
// 				"Could not get dog ID from Auth header JWT",
// 				http.StatusInternalServerError,
// 			)
// 			return
// 		}
//
// 		var notifications []xdb.Notification
//
// 		timer.WithTimer("getting notifications", func() {
// 			notifications, err = xdb.GetNotifications(db, dogId, count, offset)
// 		})
//
// 		if err != nil {
// 			http.Error(
// 				w,
// 				"Could not get notifications",
// 				http.StatusInternalServerError,
// 			)
// 			return
// 		}
//
// 		// Write notifications to response body
// 		timer.WithTimer("writing notifications to response body", func() {
// 			responseBody, err := json.Marshal(
// 				NotificationsResponseBody{Notifications: notifications},
// 			)
//
// 			if err != nil {
// 				return
// 			}
//
// 			_, err = w.Write(responseBody)
// 		})
//
// 		if err != nil {
// 			http.Error(
// 				w,
// 				"Could not write notifications to response body",
// 				http.StatusInternalServerError,
// 			)
// 		}
// 	})
// }
