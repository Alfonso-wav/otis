package shell

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// PokemonDBClient es el cliente HTTP para scraping de pokemondb.net.
// Respeta el Crawl-delay de 2 segundos entre requests.
type PokemonDBClient struct {
	baseURL    string
	httpClient *http.Client
	mu         sync.Mutex
	lastFetch  time.Time
	crawlDelay time.Duration
}

// NewPokemonDBClient crea un nuevo cliente para pokemondb.net.
func NewPokemonDBClient(baseURL string) *PokemonDBClient {
	return &PokemonDBClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		crawlDelay: 2 * time.Second,
	}
}

// fetchPage hace GET a la ruta indicada, respeta el rate limit y retorna el documento HTML parseado.
func (c *PokemonDBClient) fetchPage(path string) (*goquery.Document, error) {
	c.mu.Lock()
	elapsed := time.Since(c.lastFetch)
	if elapsed < c.crawlDelay {
		time.Sleep(c.crawlDelay - elapsed)
	}
	c.lastFetch = time.Now()
	c.mu.Unlock()

	url := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request for %s: %w", url, err)
	}
	req.Header.Set("User-Agent", "OtisPokedex/1.0 (pokemon-app; educational project)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pokemondb returned status %d for %s", resp.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML from %s: %w", url, err)
	}

	return doc, nil
}
