package shell

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	client := NewPokemonDBClient("http://localhost:0")
	client.crawlDelay = 100 * time.Millisecond // reduce for test speed

	// Simulate two consecutive rate limit checks
	client.mu.Lock()
	client.lastFetch = time.Now()
	client.mu.Unlock()

	start := time.Now()

	// This should wait for the crawl delay
	client.mu.Lock()
	elapsed := time.Since(client.lastFetch)
	if elapsed < client.crawlDelay {
		time.Sleep(client.crawlDelay - elapsed)
	}
	client.lastFetch = time.Now()
	client.mu.Unlock()

	waited := time.Since(start)
	if waited < 90*time.Millisecond {
		t.Errorf("rate limiter did not wait enough: waited %v, expected at least 90ms", waited)
	}
}

func TestNewPokemonDBClient(t *testing.T) {
	client := NewPokemonDBClient("https://pokemondb.net")

	if client.baseURL != "https://pokemondb.net" {
		t.Errorf("baseURL = %q, want %q", client.baseURL, "https://pokemondb.net")
	}
	if client.crawlDelay != 2*time.Second {
		t.Errorf("crawlDelay = %v, want %v", client.crawlDelay, 2*time.Second)
	}
	if client.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}
