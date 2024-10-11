package main

import (
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type TreatRequestBody struct {
	BarkId string `json:"bark_id"`
}

func Treat() http.Handler {
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

		var b TreatRequestBody

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

		var treatGiven bool

		timer.WithTimer("checking if treat has already been given", func() {
			treatGiven, err = xdb.HasDogGivenBarkTreat(db, b.BarkId, dogId)
		})

		if err != nil {
			http.Error(
				w,
				"Could not check if treat had already been given by dog",
				http.StatusInternalServerError,
			)
			return
		}

		switch r.Method {
		case http.MethodPost:
			PostTreat(w, r, dogId, b.BarkId, treatGiven)
		case http.MethodDelete:
			DeleteTreat(w, r, dogId, b.BarkId, treatGiven)
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

func PostTreat(
	w http.ResponseWriter,
	r *http.Request,
	dogId string,
	barkId string,
	treatGiven bool,
) {
	if treatGiven {
		w.WriteHeader(200)
		return
	}

	// TODO: Check if this dog can view that dog's barks ?
	timer.WithTimer("writing treat to database", func() {
		err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
			err = xdb.WriteTreat(e, barkId, dogId)

			if err != nil {
				return err
			}

			return xdb.IncrementTreatCount(e, barkId)
		})
	})

	if err != nil {
		http.Error(
			w,
			"Could not write treat to database",
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

func DeleteTreat(
	w http.ResponseWriter,
	r *http.Request,
	dogId string,
	barkId string,
	treatGiven bool,
) {
	if !treatGiven {
		w.WriteHeader(200)
		return
	}

	// TODO: Check if this dog can view that dog's barks ?
	timer.WithTimer("removing treat in database", func() {
		err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
			err = xdb.RemoveTreat(e, barkId, dogId)

			if err != nil {
				return err
			}

			return xdb.DecrementTreatCount(e, barkId)
		})
	})

	if err != nil {
		http.Error(
			w,
			"Could not remove treat in database",
			http.StatusBadRequest,
		)
	}
}
