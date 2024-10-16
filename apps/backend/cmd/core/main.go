package main

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
)

//   ____ _     ___  ____    _    _     ____
//  / ___| |   / _ \| __ )  / \  | |   / ___|
// | |  _| |  | | | |  _ \ / _ \ | |   \___ \
// | |_| | |__| |_| | |_) / ___ \| |___ ___) |
//  \____|_____\___/|____/_/   \_\_____|____/
//

var err error
var db *sql.DB

func initialiseAuxiliaryConnections() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		timer.WithTimer("connecting to database", func() {
			defer wg.Done()
			err = xdb.EstablishConnection(&db)

			if err != nil {
				log.Fatal("Could not connect to database")
			}
		})
	}()

	wg.Wait()
}

func defineAppRoutes() *http.ServeMux {
	auth := xhttp.BearerTokenAuthMiddlewareFactory{DBExecutor: db}
	var s *http.ServeMux = http.NewServeMux()
	s.Handle("/health", xhttp.Get(Health()))
	s.Handle("/register", xhttp.Post(Register()))
	s.Handle("/login", xhttp.Post(Login()))
	s.Handle("/bark", auth.New(Bark()))
	s.Handle("/barks", auth.New(xhttp.Get(Barks())))
	s.Handle("/follow", auth.New(xhttp.Post(Follow())))
	s.Handle("/protect", auth.New(xhttp.Post(Protect())))
	s.Handle("/approve", auth.New(xhttp.Post(Approve())))
	// s.Handle("/notifications", auth.New(xhttp.Get(Notifications())))
	s.Handle("/treat", auth.New(Treat()))
	s.Handle("/rebark", auth.New(Rebark()))
	// s.Handle("/unrebark", auth.New(xhttp.Post(Unrebark())))
	s.Handle("/timeline", auth.New(xhttp.Get(Timeline())))
	s.Handle("/paw", auth.New(xhttp.Post(Paw())))

	return s
}

func startServer() {
	// TODO: Add middleware here and give them to `defineAppRoutes`
	var s *http.ServeMux = defineAppRoutes()
	log.Fatal(http.ListenAndServe(":8000", s))
}

func main() {
	timer.WithTimer(
		"initialising auxiliary connections",
		initialiseAuxiliaryConnections,
	)
	timer.WithTimer("starting server", startServer)
}
