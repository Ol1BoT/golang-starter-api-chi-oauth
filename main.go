package main

import (
	"log"
	"net/http"

	"github.com/Ol1BoT/api-server/routes"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

func main() {

	r := chi.NewMux()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("server is running"))

	})

	r.Route("/oauth", routes.OAuthIndex)

	handler := cors.Default().Handler(r)

	if err := http.ListenAndServe("127.0.0.1:3000", handler); err != nil {

		log.Fatalln(err)

	}

}
