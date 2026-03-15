package mobile

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/alfon/pokemon-app/app"
	"github.com/alfon/pokemon-app/shell"
)

var (
	mu     sync.Mutex
	server *http.Server
)

// Start initializes the dependencies, creates the REST server and starts
// listening on the given port. dataDir is the base directory for file storage
// (teams, sprites, etc.). This function is gomobile-compatible.
func Start(port int, dataDir string) error {
	mu.Lock()
	defer mu.Unlock()

	if server != nil {
		return fmt.Errorf("server already running")
	}

	cfg := app.LoadConfig()
	fetcher := shell.NewPokeAPIClient(cfg.PokeAPIBaseURL)
	scraper := shell.NewPokemonDBClient(cfg.PokemonDBBaseURL)
	teams := shell.NewFileTeamStorage(dataDir + "/teams")
	a := app.NewApp(fetcher, scraper, teams, scraper)

	server = NewServer(a, port)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("mobile server error: %v\n", err)
		}
	}()
	return nil
}

// Stop gracefully shuts down the running server.
// This function is gomobile-compatible.
func Stop() error {
	mu.Lock()
	defer mu.Unlock()

	if server == nil {
		return fmt.Errorf("server not running")
	}

	err := server.Shutdown(context.Background())
	server = nil
	return err
}
