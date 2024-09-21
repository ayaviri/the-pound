package http

import "net/http"

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
