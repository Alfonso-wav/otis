package main

import (
	"embed"
	"log"

	"github.com/alfon/pokemon-app/app"
	"github.com/alfon/pokemon-app/shell"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend
var assets embed.FS

func main() {
	cfg := app.LoadConfig()
	fetcher := shell.NewPokeAPIClient(cfg.PokeAPIBaseURL)
	a := app.NewApp(fetcher)

	err := wails.Run(&options.App{
		Title:  "Pokédex",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: a.Startup,
		Bind:      []interface{}{a},
	})
	if err != nil {
		log.Fatal(err)
	}
}
