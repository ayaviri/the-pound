package main

import (
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type RebarkRequestBody struct {
	BarkId string `json:"bark_id"`
}

func Rebark() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		}

		var b RebarkRequestBody

		timer.WithTimer("unmarshalling body of request", func() {
			err = xhttp.ReadUnmarshalRequestBody(r, &b)
		})

		if err != nil {
			http.Error(
				w,
				"Could not get bark ID from request body",
				http.StatusBadRequest,
			)
		}

		var rebarkGiven bool

		timer.WithTimer("checking if rebark has already been given", func() {
		})

		if err != nil {
			http.Error(
				w,
				"Could not check if rebark has already been given by dog",
				http.StatusInternalServerError,
			)
			return
		}

		if rebarkGiven {
			w.WriteHeader(200)
			return
		}

		switch r.Method {
		case http.MethodPost:
			PostRebark(w, r, dogId, b.BarkId, rebarkGiven)
		case http.MethodDelete:
			DeleteRebark(w, r, dogId, b.BarkId, rebarkGiven)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

//  ____   ___  ____ _____
// |  _ \ / _ \/ ___|_   _|
// | |_) | | | \___ \ | |
// |  __/| |_| |___) || |
// |_|    \___/|____/ |_|
//

func PostRebark(
	w http.ResponseWriter,
	r *http.Request,
	dogId string,
	barkId string,
	rebarkGiven bool,
) {
	if rebarkGiven {
		w.WriteHeader(200)
		return
	}

	timer.WithTimer("writing rebark to database", func() {
		err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
			err = xdb.WriteRebark(e, barkId, dogId)

			if err != nil {
				return err
			}

			return xdb.IncrementRebarkCount(e, barkId)
		})
	})

	if err != nil {
		http.Error(
			w,
			"Could not write rebark to database",
			http.StatusBadRequest,
		)
	}
}

//  ____  _____ _     _____ _____ _____
// |  _ \| ____| |   | ____|_   _| ____|
// | | | |  _| | |   |  _|   | | |  _|
// | |_| | |___| |___| |___  | | | |___
// |____/|_____|_____|_____| |_| |_____|
//

func DeleteRebark(
	w http.ResponseWriter,
	r *http.Request,
	dogId string,
	barkId string,
	rebarkGiven bool,
) {
	if !rebarkGiven {
		w.WriteHeader(200)
		return
	}

	// TODO: Check if this dog can view that dog's barks ?
	timer.WithTimer("removing rebark in database", func() {
		err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
			err = xdb.RemoveRebark(e, barkId, dogId)

			if err != nil {
				return err
			}

			return xdb.DecrementRebarkCount(e, barkId)
		})
	})

	if err != nil {
		http.Error(
			w,
			"Could not remove rebark in database",
			http.StatusBadRequest,
		)
	}
}
