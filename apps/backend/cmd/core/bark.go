package main

import (
	"encoding/json"
	"errors"
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

func Bark() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var dogId string

		timer.WithTimer("getting dog ID from Auth header JWT", func() {
			dogId, err = xhttp.GetDogIdFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract dog ID from JWT",
				http.StatusInternalServerError,
			)
			return
		}

		switch r.Method {
		case http.MethodPost:
			PostBark(w, r, dogId)
		case http.MethodGet:
			GetBark(w, r, dogId)
		case http.MethodDelete:
			DeleteBark(w, r, dogId)
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

type PostBarkRequestBody struct {
	Content string `json:"content"`
}

func PostBark(w http.ResponseWriter, r *http.Request, dogId string) {
	var b PostBarkRequestBody

	timer.WithTimer("unmarshalling body of request", func() {
		err = xhttp.ReadUnmarshalRequestBody(r, &b)
	})

	if err != nil {
		http.Error(
			w,
			"Could not extract bark from request body",
			http.StatusBadRequest,
		)
		return
	}

	timer.WithTimer("writing bark to the database", func() {
		err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
			_, err = xdb.WriteBark(e, b.Content, dogId)
			return err
		})
	})

	if err != nil {
		http.Error(
			w,
			"Could not write bark to database",
			http.StatusInternalServerError,
		)
		return
	}
}

//   ____ _____ _____
//  / ___| ____|_   _|
// | |  _|  _|   | |
// | |_| | |___  | |
//  \____|_____| |_|
//

type BarkQueryStringParameters struct {
	BarkId string
}

type GetBarkResponseBody struct {
	Bark xdb.Bark
	// TODO: Consider throwing the counts of the various interactions with this bark
}

func GetBark(w http.ResponseWriter, r *http.Request, dogId string) {
	var p BarkQueryStringParameters

	timer.WithTimer("getting bark ID from query string", func() {
		p, err = parseBarkQueryStringParameters(r)
	})

	if err != nil {
		http.Error(
			w,
			"Could not read query string parameters",
			http.StatusBadRequest,
		)
		return
	}

	var bark xdb.Bark

	timer.WithTimer("getting bark from database", func() {
		bark, err = xdb.GetBark(db, p.BarkId)
	})

	if err != nil {
		http.Error(
			w,
			"Could not get bark from database",
			http.StatusInternalServerError,
		)
		return
	}

	timer.WithTimer("writing barks to response body", func() {
		responseBody, err := json.Marshal(GetBarkResponseBody{Bark: bark})

		if err != nil {
			return
		}

		_, err = w.Write(responseBody)
	})

	if err != nil {
		http.Error(
			w,
			"Could not write bark to response body",
			http.StatusInternalServerError,
		)
	}
}

//  ____  _____ _     _____ _____ _____
// |  _ \| ____| |   | ____|_   _| ____|
// | | | |  _| | |   |  _|   | | |  _|
// | |_| | |___| |___| |___  | | | |___
// |____/|_____|_____|_____| |_| |_____|
//

func DeleteBark(w http.ResponseWriter, r *http.Request, dogId string) {
	var p BarkQueryStringParameters

	timer.WithTimer("getting bark ID from query string", func() {
		p, err = parseBarkQueryStringParameters(r)
	})

	if err != nil {
		http.Error(
			w,
			"Could not read query string parameters",
			http.StatusBadRequest,
		)
		return
	}

	var isBarkPaw bool
	// ID of the bark that was pawed by this one, in case it is a paw
	var originalBarkId string

	timer.WithTimer("checking if bark is paw", func() {
		isBarkPaw, err = xdb.IsBarkPaw(db, p.BarkId)

		if err != nil {
			return
		}

		if isBarkPaw {
			originalBarkId, err = xdb.GetOriginalBarkId(db, p.BarkId)
		}
	})

	if err != nil {
		http.Error(w, "Could not check if bark is paw", http.StatusInternalServerError)
		return
	}

	timer.WithTimer("removing bark from database", func() {
		err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
			err = xdb.RemoveBark(e, p.BarkId)

			if err == nil && isBarkPaw {
				return xdb.DecrementPawCount(e, originalBarkId)
			} else {
				return err
			}
		})
	})

	if err != nil {
		http.Error(
			w,
			"Could not remove bark from database",
			http.StatusInternalServerError,
		)
	}
}

func parseBarkQueryStringParameters(
	r *http.Request,
) (BarkQueryStringParameters, error) {
	var p BarkQueryStringParameters
	var barkIds []string = r.URL.Query()["id"]

	if barkIds == nil || barkIds[0] == "" {
		return p, errors.New("Bark ID has not been set")
	}

	p.BarkId = barkIds[0]

	return p, nil
}
