package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"

	xdb "the-pound/internal/db"
	xhttp "the-pound/internal/http"

	"github.com/ayaviri/goutils/timer"
	"github.com/rs/cors"
)

//   ____ _     ___  ____    _    _     ____
//  / ___| |   / _ \| __ )  / \  | |   / ___|
// | |  _| |  | | | |  _ \ / _ \ | |   \___ \
// | |_| | |__| |_| | |_) / ___ \| |___ ___) |
//  \____|_____\___/|____/_/   \_\_____|____/
//

var err error
var db *sql.DB
var FE_SERVER_URL string

func initialiseAuxiliaryConnections() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		timer.WithTimer("connecting to database", func() {
			defer wg.Done()
			err = xdb.EstablishConnection(&db)

			if err != nil {
				log.Fatal("Could not connect to database")
			}
		})
	}()

	go func() {
		timer.WithTimer("reading environment variables", func() {
			defer wg.Done()
			var isPresent bool
			FE_SERVER_URL, isPresent = os.LookupEnv("FE_SERVER_URL")

			if !isPresent {
				log.Fatal("Could not read frontend server URL")
			}
		})
	}()

	wg.Wait()
}

func defineAppRoutes() *http.ServeMux {
	auth := xhttp.BearerTokenAuthMiddlewareFactory{DBExecutor: db}
	logging := xhttp.NewLoggingHandler(os.Stdout)
	c := cors.New(
		cors.Options{
			AllowedOrigins: []string{FE_SERVER_URL},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			ExposedHeaders:   []string{"Authorization"},
		},
	).Handler

	var s *http.ServeMux = http.NewServeMux()
	s.Handle("/health", logging(xhttp.Get(Health())))
	s.Handle("/register", c(logging(xhttp.Post(Register()))))
	s.Handle("/login", c(logging(xhttp.Post(Login()))))
	s.Handle("/bark", c(logging(auth.New(Bark()))))
	s.Handle("/barks", c(logging(auth.New(xhttp.Get(Barks())))))
	s.Handle("/protect", c(logging(auth.New(xhttp.Post(Protect())))))
	s.Handle("/approve", c(logging(auth.New(xhttp.Post(Approve())))))
	s.Handle("/reject", c(logging(auth.New(xhttp.Post(Reject())))))
	s.Handle("/notifications", c(logging(auth.New(xhttp.Get(Notifications())))))
	s.Handle("/notification_read", c(logging(auth.New(xhttp.Post(NotificationRead())))))
	s.Handle("/timeline", c(logging(auth.New(xhttp.Get(Timeline())))))
	s.Handle("/paw", c(logging(auth.New(xhttp.Post(Paw())))))
	s.Handle("/paws", c(logging(auth.New(xhttp.Get(Paws())))))
	s.Handle("/validate", c(logging(auth.New(xhttp.Nil()))))
	s.Handle("/thread", c(logging(auth.New(xhttp.Get(Thread())))))
	s.Handle("/dog", c(logging(auth.New(xhttp.Get(Dog())))))
	s.Handle("/does_follow", c(logging(auth.New(xhttp.Get(DoesFollow())))))

	//  _____ ___   ____  ____ _     _____ ____
	// |_   _/ _ \ / ___|/ ___| |   | ____/ ___|
	//   | || | | | |  _| |  _| |   |  _| \___ \
	//   | || |_| | |_| | |_| | |___| |___ ___) |
	//   |_| \___/ \____|\____|_____|_____|____/
	//

	s.Handle("/treat", c(logging(auth.New(xhttp.Post(Treat())))))
	s.Handle("/rebark", c(logging(auth.New(xhttp.Post(Rebark())))))
	s.Handle("/follow", c(logging(auth.New(xhttp.Post(Follow())))))

	return s
}

func startServer() {
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
