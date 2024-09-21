package main

import "net/http"

//   ____ _     ___  ____    _    _     ____
//  / ___| |   / _ \| __ )  / \  | |   / ___|
// | |  _| |  | | | |  _ \ / _ \ | |   \___ \
// | |_| | |__| |_| | |_) / ___ \| |___ ___) |
//  \____|_____\___/|____/_/   \_\_____|____/
//

func Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err = w.Write([]byte("taking your thoughts to the pound."))

		if err != nil {
			http.Error(
				w,
				"Failed to write out",
				http.StatusInternalServerError,
			)
			return
		}
	})
}
