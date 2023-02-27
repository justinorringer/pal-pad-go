package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/justinorringer/pal-pad-go/endpoints"
)

func main() {
	r := chi.NewRouter()
	sub := chi.NewRouter()

	r.Mount("/api/v1", sub)
	r.Get("/", endpoints.Lubdub)

	// GET /api/v1/pads/{pad_id} // get the canvas
	// POST /api/v1/user/ // json or something for the data, return user id for frontend
	// POST /api/v1/pads/{pad_id}/ // and maybe pass the user here too so we know who made the change

	sub.Post("/pads", endpoints.Lubdub)

	srv := http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
