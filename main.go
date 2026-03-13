package main

import (
	"log"
	"net/http"

	"github.com/alfon/pokemon-app/app"
	"github.com/alfon/pokemon-app/shell"
)

func main() {
	cfg := app.LoadConfig()
	fetcher := shell.NewPokeAPIClient(cfg.PokeAPIBaseURL)
	mux := app.NewServer(cfg, fetcher)

	log.Printf("Server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}
