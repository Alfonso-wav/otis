package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alfon/pokemon-app/app/mobile"
)

func main() {
	port := 8080
	dataDir := "data"

	fmt.Printf("Starting mobile REST server on http://localhost:%d\n", port)
	fmt.Printf("Data directory: %s\n", dataDir)

	if err := mobile.Start(port, dataDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server running. Press Ctrl+C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("\nShutting down...")
	mobile.Stop()
}
