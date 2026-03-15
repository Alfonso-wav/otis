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

	"github.com/PuerkitoBio/goquery"
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

// BattleSpriteURLs holds the scraped Gen 1 sprite URLs for a Pokémon.
type BattleSpriteURLs struct {
	Front string
	Back  string
}

// ScrapeBattleSpriteURLs visits /sprites/{name} and extracts the oldest generation
// Normal (front) and Back sprite image URLs.
func ScrapeBattleSpriteURLs(doc *goquery.Document) BattleSpriteURLs {
	var result BattleSpriteURLs

	// Find all generation headings (h2 elements containing "Generation")
	var genSections []*goquery.Selection
	doc.Find("h2").Each(func(_ int, h2 *goquery.Selection) {
		text := strings.TrimSpace(h2.Text())
		if strings.HasPrefix(text, "Generation") {
			genSections = append(genSections, h2)
		}
	})

	if len(genSections) == 0 {
		return result
	}

	// Try each generation section starting from the oldest (first on page)
	for _, h2 := range genSections {
		// The table follows the h2 heading — find next sibling elements until we hit a table
		var table *goquery.Selection
		for sib := h2.Next(); sib.Length() > 0; sib = sib.Next() {
			if goquery.NodeName(sib) == "table" {
				table = sib
				break
			}
			// Stop if we hit the next h2 (next generation)
			if goquery.NodeName(sib) == "h2" {
				break
			}
		}
		if table == nil {
			continue
		}

		// Parse column headers to find "Normal" and "Back" indices
		var colHeaders []string
		table.Find("thead th").Each(func(_ int, th *goquery.Selection) {
			colHeaders = append(colHeaders, strings.TrimSpace(th.Text()))
		})

		normalIdx := -1
		backIdx := -1
		for i, h := range colHeaders {
			lower := strings.ToLower(h)
			if lower == "normal" && normalIdx == -1 {
				normalIdx = i
			}
			if lower == "back" && backIdx == -1 {
				backIdx = i
			}
		}

		// Look through table rows for the first row with sprites
		table.Find("tbody tr").Each(func(_ int, tr *goquery.Selection) {
			if result.Front != "" && result.Back != "" {
				return
			}
			cells := tr.Find("td")

			if normalIdx >= 0 && result.Front == "" {
				img := cells.Eq(normalIdx).Find("span:nth-child(2) a img").First()
				if src, exists := img.Attr("src"); exists && src != "" {
					result.Front = src
				}
			}

			if backIdx >= 0 && result.Back == "" {
				img := cells.Eq(backIdx).Find("span:nth-child(2) a img").First()
				if src, exists := img.Attr("src"); exists && src != "" {
					result.Back = src
				}
			}
		})

		// If we found at least one sprite in this generation, use it
		if result.Front != "" || result.Back != "" {
			return result
		}
	}

	return result
}

func (c *PokemonDBClient) downloadBattleSprites(destDir string, entries []core.PokedexDBEntry, result *core.SpriteDownloadResult) {
	backDir := filepath.Join(destDir, string(core.SpriteCategoryBattleBack))
	frontDir := filepath.Join(destDir, string(core.SpriteCategoryBattleFront))

	if err := os.MkdirAll(backDir, 0755); err != nil {
		result.Failed++
		result.Errors = append(result.Errors, fmt.Sprintf("creating battle-back dir: %s", err))
		return
	}
	if err := os.MkdirAll(frontDir, 0755); err != nil {
		result.Failed++
		result.Errors = append(result.Errors, fmt.Sprintf("creating battle-front dir: %s", err))
		return
	}

	for i, entry := range entries {
		name := strings.ToLower(entry.Name)

		backDest := filepath.Join(backDir, name+".png")
		frontDest := filepath.Join(frontDir, name+".png")

		_, backExists := os.Stat(backDest)
		_, frontExists := os.Stat(frontDest)

		// Count both sprites
		result.Total += 2

		if backExists == nil && frontExists == nil {
			result.Skipped += 2
			continue
		}

		// Need to scrape the page to get URLs
		doc, err := c.fetchPage(fmt.Sprintf("/sprites/%s", name))
		if err != nil {
			if backExists != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("battle-back/%s: %s", name, err))
			} else {
				result.Skipped++
			}
			if frontExists != nil {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("battle-front/%s: %s", name, err))
			} else {
				result.Skipped++
			}
			continue
		}

		urls := ScrapeBattleSpriteURLs(doc)

		// Download back sprite
		if backExists != nil {
			if urls.Back != "" {
				if dlErr := c.downloadFile(urls.Back, backDest); dlErr != nil {
					result.Failed++
					result.Errors = append(result.Errors, fmt.Sprintf("battle-back/%s: %s", name, dlErr))
				} else {
					result.Downloaded++
				}
			} else {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("battle-back/%s: no back sprite found", name))
			}
		} else {
			result.Skipped++
		}

		// Download front sprite
		if frontExists != nil {
			if urls.Front != "" {
				if dlErr := c.downloadFile(urls.Front, frontDest); dlErr != nil {
					result.Failed++
					result.Errors = append(result.Errors, fmt.Sprintf("battle-front/%s: %s", name, dlErr))
				} else {
					result.Downloaded++
				}
			} else {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("battle-front/%s: no front sprite found", name))
			}
		} else {
			result.Skipped++
		}

		if (i+1)%50 == 0 {
			log.Printf("Battle sprites: %d/%d processed", i+1, len(entries))
		}
	}
}

func (c *PokemonDBClient) DownloadAllSprites(destDir string, categories []core.SpriteCategory) (core.SpriteDownloadResult, error) {
	result := core.SpriteDownloadResult{}

	entries, err := c.FetchPokedex()
	if err != nil {
		return result, fmt.Errorf("fetching pokedex for sprite names: %w", err)
	}

	// Check if battle sprite categories are requested
	hasBattle := false
	for _, cat := range categories {
		if cat == core.SpriteCategoryBattleBack || cat == core.SpriteCategoryBattleFront {
			hasBattle = true
			break
		}
	}

	if hasBattle {
		c.downloadBattleSprites(destDir, entries, &result)
	}

	for _, cat := range categories {
		// Battle sprites are handled together above
		if cat == core.SpriteCategoryBattleBack || cat == core.SpriteCategoryBattleFront {
			continue
		}

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
