package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinorringer/pal-pad-go/db"
	"github.com/justinorringer/pal-pad-go/endpoints"
	"github.com/justinorringer/pal-pad-go/sockets"
)

func main() {
	rc, err := db.NewRedisClient("localhost:6379", 0, 0)
	if err != nil { // if the database connection fails, panic
		panic(err)
	}

	hub := sockets.NewHub()
	go hub.Run(rc) // hub handles all websocket connections

	ws := func(w http.ResponseWriter, r *http.Request) {
		sockets.ServeWs(hub, w, r)
	}

	r := chi.NewRouter()
	sub := chi.NewRouter()

	r.Mount("/api/v1", sub)
	r.Get("/", sockets.ServeHome)
	r.Get("/ws", ws)
	r.Get("/ws-test", sockets.ServeHome)

	sub.Post("/user", endpoints.Lubdub) // create a new user, return the user id

	srv := http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
