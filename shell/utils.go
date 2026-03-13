package shell

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alfon/pokemon-app/core"
)

// apiSimpleList es la estructura genérica para endpoints de lista sin paginación.
type apiSimpleList struct {
	Count   int `json:"count"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// extractIDFromURL extrae el ID numérico al final de una URL de PokéAPI.
// Ej: "https://pokeapi.co/api/v2/evolution-chain/1/" → 1
func extractIDFromURL(rawURL string) int {
	parts := strings.Split(strings.TrimRight(rawURL, "/"), "/")
	if len(parts) == 0 {
		return 0
	}
	id, _ := strconv.Atoi(parts[len(parts)-1])
	return id
}

// fetchNamedResourceList obtiene una lista genérica de recursos nombrados.
func fetchNamedResourceList(c *PokeAPIClient, endpoint string) ([]core.NamedResource, error) {
	url := fmt.Sprintf("%s/%s?limit=10000", c.baseURL, endpoint)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pokeapi returned status %d for %s", resp.StatusCode, endpoint)
	}

	var raw apiSimpleList
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decoding %s: %w", endpoint, err)
	}

	results := make([]core.NamedResource, len(raw.Results))
	for i, r := range raw.Results {
		results[i] = core.NamedResource{Name: r.Name, URL: r.URL}
	}
	return results, nil
}

// fetchOne realiza un GET y decodifica el body en dest.
// Devuelve error si el status es 404 o no es 200.
func fetchOne(c *PokeAPIClient, url string, dest interface{}) error {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found: %s", url)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pokeapi returned status %d: %s", resp.StatusCode, url)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decoding %s: %w", url, err)
	}
	return nil
}

// firstEnglish busca el primer texto en inglés de una lista de entradas
// con FlavorText/Genus + Language.
func firstEnglishFlavor(entries []struct {
	FlavorText string `json:"flavor_text"`
	Language   struct{ Name string `json:"name"` } `json:"language"`
}) string {
	for _, e := range entries {
		if e.Language.Name == "en" {
			return e.FlavorText
		}
	}
	return ""
}

// firstEnglishDesc busca la primera descripción en inglés.
func firstEnglishDesc(entries []struct {
	Description string `json:"description"`
	Language    struct{ Name string `json:"name"` } `json:"language"`
}) string {
	for _, e := range entries {
		if e.Language.Name == "en" {
			return e.Description
		}
	}
	return ""
}
