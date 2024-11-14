package http

import (
	"database/sql"
	"io"
	"net/http"
	xdb "the-pound/internal/db"

	"github.com/ayaviri/goutils/timer"
	"github.com/gorilla/handlers"
)

//  _     ___   ____  ____ ___ _   _  ____
// | |   / _ \ / ___|/ ___|_ _| \ | |/ ___|
// | |  | | | | |  _| |  _ | ||  \| | |  _
// | |__| |_| | |_| | |_| || || |\  | |_| |
// |_____\___/ \____|\____|___|_| \_|\____|
//

func NewLoggingHandler(destination io.Writer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(destination, next)
	}
}

//     _   _   _ _____ _   _
//    / \ | | | |_   _| | | |
//   / _ \| | | | | | | |_| |
//  / ___ \ |_| | | | |  _  |
// /_/   \_\___/  |_| |_| |_|
//

type BearerTokenAuthMiddlewareFactory struct {
	DBExecutor xdb.DBExecutor
}

func (f BearerTokenAuthMiddlewareFactory) New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bearerToken string

		timer.WithTimer("getting bearer token from request header", func() {
			bearerToken, err = GetAuthBearerToken(r)
		})

		if err != nil {
			http.Error(
				w,
				"Could not obtain bearer token from Authorization header",
				http.StatusBadRequest,
			)
			return
		}

		var result xdb.JWTValidationResult

		timer.WithTimer("validating bearer token", func() {
			result, err = xdb.IsValidJWT(f.DBExecutor, bearerToken)
		})

		// JWT invalidation is handled weirdly by the jwt package, so
		// invalid JWTs will also through an error. I'm not going through
		// the effort of checking each error to determine whether it's a
		// server or client side issue
		if err != nil {
			http.Error(
				w,
				"Could not validate bearer token or bearer token is not valid",
				http.StatusInternalServerError,
			)
			return
		}

		if !result.IsValid {
			http.Error(w, "Bearer token is not valid", http.StatusUnauthorized)
			return
		}

		if result.NewToken != "" {
			timer.WithTimer("updating bearer token", func() {
				db, _ := f.DBExecutor.(*sql.DB)
				err = xdb.ExecuteInTransaction(db, func(e xdb.DBExecutor) error {
					return xdb.UpdateSessionToken(e, bearerToken, result.NewToken)
				})

				if err != nil {
					return
				}

				newHeader := "Bearer " + result.NewToken
				r.Header.Set("Authorization", newHeader)
				w.Header().Set("Authorization", newHeader)
			})
		}

		if err != nil {
			http.Error(
				w,
				"Could not update session token",
				http.StatusInternalServerError,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//  __  __ _____ _____ _   _  ___  ____
// |  \/  | ____|_   _| | | |/ _ \|  _ \
// | |\/| |  _|   | | | |_| | | | | | | |
// | |  | | |___  | | |  _  | |_| | |_| |
// |_|  |_|_____| |_| |_| |_|\___/|____/
//
//  _   _    _    _   _ ____  _     _____ ____  ____
// | | | |  / \  | \ | |  _ \| |   | ____|  _ \/ ___|
// | |_| | / _ \ |  \| | | | | |   |  _| | |_) \___ \
// |  _  |/ ___ \| |\  | |_| | |___| |___|  _ < ___) |
// |_| |_/_/   \_\_| \_|____/|_____|_____|_| \_\____/
//

func Get(next http.Handler) http.Handler {
	return methodHandler(http.MethodGet, next)
}

func Post(next http.Handler) http.Handler {
	return methodHandler(http.MethodPost, next)
}

func Delete(next http.Handler) http.Handler {
	return methodHandler(http.MethodDelete, next)
}

func methodHandler(httpMethod string, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != httpMethod {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			if next != nil {
				next.ServeHTTP(w, r)
			}
		},
	)
}
