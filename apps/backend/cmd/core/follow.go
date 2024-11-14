package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

type FollowRequestBody struct {
	DogId string `json:"dog_id"`
}

type FollowResponseBody struct {
	xdb.FollowResult
}

func Follow() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fromDogId string

		timer.WithTimer("getting dog requesting follow from Auth header JWT", func() {
			fromDogId, err = xhttp.GetDogIdFromAuth(db, r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract ID of dog requesting follow from JWT",
				http.StatusInternalServerError,
			)
			return
		}

		var toDog xdb.Dog

		timer.WithTimer("getting dog to be followed from request body", func() {
			var b FollowRequestBody
			err = xhttp.ReadUnmarshalRequestBody(r, &b)

			if err != nil {
				return
			}

			toDog, err = xdb.GetDogFromId(db, b.DogId)
		})

		if err != nil {
			http.Error(
				w,
				"Could not extract dog to be followed from request body",
				http.StatusBadRequest,
			)
			return
		}

		var result xdb.FollowResult

		timer.WithTimer(
			"toggling follow state for dog pair in database",
			func() {
				result, err = toggleFollowStateInDB(db, fromDogId, toDog)
			},
		)

		if err != nil {
			http.Error(
				w,
				"Could not toggle follow relationship in database",
				http.StatusInternalServerError,
			)
			return
		}

		timer.WithTimer("writing result to response body", func() {
			responseBody, err := json.Marshal(FollowResponseBody{FollowResult: result})

			if err != nil {
				return
			}

			_, err = w.Write(responseBody)
		})

		if err != nil {
			http.Error(
				w,
				"Could not write follow result to response body",
				http.StatusInternalServerError,
			)
		}
	})
}

// In an atomic transaction
// 1) Checks whether a follow relationship exists for the dog pair
// 2) Writes/removes rebark accordingly
// 3) Increments/decrements bark's rebark count accordingly
// Returns a struct containing whether a follow was created and whether the
// created follow was approved
func toggleFollowStateInDB(
	db *sql.DB,
	fromDogId string,
	toDog xdb.Dog,
) (xdb.FollowResult, error) {
	resultAfterToggle := xdb.FollowResult{
		FollowRequestExists: false,
		IsApproved:          false,
	}

	return resultAfterToggle, xdb.ExecuteInTransaction(
		db,
		func(e xdb.DBExecutor) error {
			resultBeforeToggle, err := xdb.GetFollowResult(e, fromDogId, toDog.Id)

			if err != nil {
				return err
			}

			resultAfterToggle.FollowRequestExists = !resultBeforeToggle.FollowRequestExists

			if resultBeforeToggle.FollowRequestExists {
				return removeFollowFromDB(
					e,
					fromDogId,
					toDog.Id,
					resultBeforeToggle.IsApproved,
				)
			} else {
				isApproved := toDog.IsPublic
				resultAfterToggle.IsApproved = isApproved
				return WriteFollowToDB(e, fromDogId, toDog.Id, isApproved)
			}
		},
	)
}

func WriteFollowToDB(
	e xdb.DBExecutor,
	fromDogId string,
	toDogId string,
	isApproved bool,
) error {
	err = xdb.WriteFollow(e, fromDogId, toDogId, isApproved)

	if err != nil {
		return err
	}

	if !isApproved {
		return nil
	}

	err = xdb.IncrementFollowingCount(e, fromDogId)

	if err != nil {
		return err
	}

	return xdb.IncrementFollowerCount(e, toDogId)
}

func removeFollowFromDB(
	e xdb.DBExecutor,
	fromDogId string,
	toDogId string,
	// True if the following to be removed was approved and thus affected
	// the follower/following count of both dogs
	wasApproved bool,
) error {
	err = xdb.RemoveFollow(db, fromDogId, toDogId)

	if err != nil {
		return err
	}

	if !wasApproved {
		return nil
	}

	err = xdb.DecrementFollowingCount(e, fromDogId)

	if err != nil {
		return err
	}

	return xdb.DecrementFollowerCount(e, toDogId)
}
