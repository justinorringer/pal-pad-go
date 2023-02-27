package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinorringer/pal-pad-go/endpoints"
	"github.com/justinorringer/pal-pad-go/sockets"
)

func main() {
	hub := sockets.NewHub()
	go hub.Run()

	ws := func(w http.ResponseWriter, r *http.Request) {
		sockets.ServeWs(hub, w, r)
	}

	r := chi.NewRouter()
	sub := chi.NewRouter()

	r.Mount("/api/v1", sub)
	r.Get("/", endpoints.Lubdub)
	r.Get("/ws", ws)

	sub.Post("/user", endpoints.Lubdub) // create a new user, return the user id

	srv := http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
