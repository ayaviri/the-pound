package http

import "net/http"

// "Endware" that serves as the end of a chain of handlers for an HTTP request.
// Consumes the previous handler and returns a handler that does nothing
func Nil() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
