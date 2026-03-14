package shell

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alfon/pokemon-app/core"
)

const spriteBaseURL = "https://img.pokemondb.net"

var iconFiles = []struct {
	remotePath string
	localName  string
}{
	{"/images/icons/move-physical.png", "move-physical.png"},
	{"/images/icons/move-special.png", "move-special.png"},
	{"/images/icons/move-status.png", "move-status.png"},
}

func spriteRemotePath(category core.SpriteCategory, name string) string {
	switch category {
	case core.SpriteCategoryHomeNormal:
		return fmt.Sprintf("/sprites/home/normal/%s.png", name)
	case core.SpriteCategoryHomeShiny:
		return fmt.Sprintf("/sprites/home/shiny/%s.png", name)
	default:
		return ""
	}
}

func (c *PokemonDBClient) DownloadAllSprites(destDir string, categories []core.SpriteCategory) (core.SpriteDownloadResult, error) {
	result := core.SpriteDownloadResult{}

	entries, err := c.FetchPokedex()
	if err != nil {
		return result, fmt.Errorf("fetching pokedex for sprite names: %w", err)
	}

	for _, cat := range categories {
		if cat == core.SpriteCategoryIcons {
			catDir := filepath.Join(destDir, string(cat))
			if err := os.MkdirAll(catDir, 0755); err != nil {
				return result, fmt.Errorf("creating icon directory: %w", err)
			}
			for _, icon := range iconFiles {
				result.Total++
				dest := filepath.Join(catDir, icon.localName)
				if _, statErr := os.Stat(dest); statErr == nil {
					result.Skipped++
					continue
				}
				url := spriteBaseURL + icon.remotePath
				if dlErr := c.downloadFile(url, dest); dlErr != nil {
					result.Failed++
					result.Errors = append(result.Errors, fmt.Sprintf("icon %s: %s", icon.localName, dlErr))
				} else {
					result.Downloaded++
				}
			}
			continue
		}

		catDir := filepath.Join(destDir, string(cat))
		if err := os.MkdirAll(catDir, 0755); err != nil {
			return result, fmt.Errorf("creating directory %s: %w", catDir, err)
		}

		for i, entry := range entries {
			name := strings.ToLower(entry.Name)
			result.Total++

			dest := filepath.Join(catDir, name+".png")
			if _, statErr := os.Stat(dest); statErr == nil {
				result.Skipped++
				continue
			}

			remotePath := spriteRemotePath(cat, name)
			if remotePath == "" {
				result.Skipped++
				continue
			}

			url := spriteBaseURL + remotePath
			if dlErr := c.downloadFile(url, dest); dlErr != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("%s/%s: %s", cat, name, dlErr))
			} else {
				result.Downloaded++
			}

			if (i+1)%50 == 0 {
				log.Printf("Sprites %s: %d/%d processed", cat, i+1, len(entries))
			}
		}
	}

	return result, nil
}

func (c *PokemonDBClient) downloadFile(url, destPath string) error {
	c.mu.Lock()
	elapsed := time.Since(c.lastFetch)
	if elapsed < c.crawlDelay {
		time.Sleep(c.crawlDelay - elapsed)
	}
	c.lastFetch = time.Now()
	c.mu.Unlock()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "OtisPokedex/1.0 (pokemon-app; educational project)")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("downloading: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d for %s", resp.StatusCode, url)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
