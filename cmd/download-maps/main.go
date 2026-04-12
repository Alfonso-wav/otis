// Command download-maps fetches Pokemon region map images from the
// Fandom wiki using the MediaWiki API and saves them to
// frontend/public/assets/maps/.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	wikiAPI   = "https://pokemon.fandom.com/api.php"
	userAgent = "OtisPokedex/1.0 (pokemon-app; educational project)"
	destDir   = "frontend/public/assets/maps"
	retroURL  = "https://blog.vjeux.com/wp-content/uploads/2023/12/pokemon_blue-1.png"
)

// regionMap maps local filename → Fandom wiki image filename.
var regionMap = []struct {
	local    string
	wikiFile string
}{
	{"kanto", "Kanto_Let's_Go,_Pikachu!_and_Let's_Go,_Eevee!.png"},
	{"johto", "Johto_HGSS.png"},
	{"hoenn", "Hoenn_ORAS.png"},
	{"sinnoh", "Sinnoh_BDSP.png"},
	{"unova", "Unova.jpg"},
	{"kalos", "Kalos.jpg"},
	{"alola", "Alola.png"},
	{"galar", "Galar.jpg"},
	{"paldea", "Paldea.jpg"},
	{"hisui", "Hisui_Legends_Arceus.png"},
}

type mediaWikiResponse struct {
	Query struct {
		Pages map[string]struct {
			ImageInfo []struct {
				URL string `json:"url"`
			} `json:"imageinfo"`
		} `json:"pages"`
	} `json:"query"`
}

func resolveWikiImageURL(client *http.Client, wikiFile string) (string, error) {
	params := url.Values{}
	params.Set("action", "query")
	params.Set("titles", "File:"+wikiFile)
	params.Set("prop", "imageinfo")
	params.Set("iiprop", "url")
	params.Set("format", "json")

	reqURL := wikiAPI + "?" + params.Encode()

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var mw mediaWikiResponse
	if err := json.NewDecoder(resp.Body).Decode(&mw); err != nil {
		return "", fmt.Errorf("decoding API response: %w", err)
	}

	for _, page := range mw.Query.Pages {
		if len(page.ImageInfo) > 0 {
			return page.ImageInfo[0].URL, nil
		}
	}
	return "", fmt.Errorf("no imageinfo found for %s", wikiFile)
}

func downloadFile(client *http.Client, fileURL, destPath string) error {
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d for %s", resp.StatusCode, fileURL)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		os.Remove(destPath)
		return err
	}

	log.Printf("  saved %s (%d KB)", filepath.Base(destPath), n/1024)
	return nil
}

func main() {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		log.Fatalf("creating dest dir: %s", err)
	}

	client := &http.Client{Timeout: 60 * time.Second}

	downloaded, skipped, failed := 0, 0, 0

	// Download region maps from Fandom wiki
	for _, region := range regionMap {
		ext := filepath.Ext(region.wikiFile)
		localFile := region.local + ext
		destPath := filepath.Join(destDir, localFile)

		if _, err := os.Stat(destPath); err == nil {
			log.Printf("[skip] %s already exists", localFile)
			skipped++
			continue
		}

		log.Printf("[fetch] %s → resolving wiki URL...", region.local)
		imgURL, err := resolveWikiImageURL(client, region.wikiFile)
		if err != nil {
			log.Printf("[FAIL] %s: resolve URL: %s", region.local, err)
			failed++
			continue
		}

		// Rate limit: 2s between requests
		time.Sleep(2 * time.Second)

		log.Printf("[download] %s from %s", region.local, truncURL(imgURL))
		if err := downloadFile(client, imgURL, destPath); err != nil {
			log.Printf("[FAIL] %s: download: %s", region.local, err)
			failed++
			continue
		}
		downloaded++

		// Rate limit between iterations
		time.Sleep(2 * time.Second)
	}

	// Download world map (Reddit source - direct imgur/reddit image)
	worldDest := filepath.Join(destDir, "world.png")
	if _, err := os.Stat(worldDest); err == nil {
		log.Printf("[skip] world.png already exists")
		skipped++
	} else {
		log.Printf("[fetch] world map from Fandom wiki...")
		worldURL, err := resolveWikiImageURL(client, "Regions.jpg")
		if err != nil {
			log.Printf("[FAIL] world: resolve URL: %s", err)
			failed++
		} else {
			time.Sleep(2 * time.Second)
			worldFile := filepath.Join(destDir, "world.jpg")
			if err := downloadFile(client, worldURL, worldFile); err != nil {
				log.Printf("[FAIL] world: download: %s", err)
				failed++
			} else {
				downloaded++
			}
		}
	}

	// Download retro Game Boy map
	retroDest := filepath.Join(destDir, "retro-blue.png")
	if _, err := os.Stat(retroDest); err == nil {
		log.Printf("[skip] retro-blue.png already exists")
		skipped++
	} else {
		log.Printf("[fetch] retro map...")
		time.Sleep(2 * time.Second)
		if err := downloadFile(client, retroURL, retroDest); err != nil {
			log.Printf("[FAIL] retro: download: %s", err)
			failed++
		} else {
			downloaded++
		}
	}

	log.Println("---")
	log.Printf("Done: %d downloaded, %d skipped, %d failed", downloaded, skipped, failed)
}

func truncURL(u string) string {
	if len(u) > 80 {
		return u[:77] + "..."
	}
	return u
}
