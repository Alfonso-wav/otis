package shell

import (
	"fmt"

	"github.com/alfon/pokemon-app/core"
)

type apiBerryList struct {
	Count   int `json:"count"`
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type apiBerry struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	GrowthTime       int    `json:"growth_time"`
	MaxHarvest       int    `json:"max_harvest"`
	NaturalGiftPower int    `json:"natural_gift_power"`
	Size             int    `json:"size"`
	Smoothness       int    `json:"smoothness"`
	SoilDryness      int    `json:"soil_dryness"`
	Firmness         struct {
		Name string `json:"name"`
	} `json:"firmness"`
	Flavors []struct {
		Flavor  struct{ Name string `json:"name"` } `json:"flavor"`
		Potency int                                  `json:"potency"`
	} `json:"flavors"`
	NaturalGiftType struct {
		Name string `json:"name"`
	} `json:"natural_gift_type"`
	Item struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"item"`
}

type apiItem struct {
	Sprites struct {
		Default string `json:"default"`
	} `json:"sprites"`
	EffectEntries []struct {
		ShortEffect string `json:"short_effect"`
		Language    struct {
			Name string `json:"name"`
		} `json:"language"`
	} `json:"effect_entries"`
}

func (c *PokeAPIClient) FetchBerryList() (core.BerryListResponse, error) {
	url := fmt.Sprintf("%s/berry?limit=100", c.baseURL)
	var raw apiBerryList
	if err := fetchOne(c, url, &raw); err != nil {
		return core.BerryListResponse{}, fmt.Errorf("fetching berry list: %w", err)
	}

	results := make([]core.BerryListItem, len(raw.Results))
	for i, r := range raw.Results {
		results[i] = core.BerryListItem{Name: r.Name, URL: r.URL}
	}
	return core.BerryListResponse{Count: raw.Count, Results: results}, nil
}

func (c *PokeAPIClient) FetchBerry(name string) (core.Berry, error) {
	url := fmt.Sprintf("%s/berry/%s", c.baseURL, name)
	var raw apiBerry
	if err := fetchOne(c, url, &raw); err != nil {
		return core.Berry{}, fmt.Errorf("fetching berry %s: %w", name, err)
	}

	flavors := make([]core.BerryFlavor, 0, len(raw.Flavors))
	for _, f := range raw.Flavors {
		if f.Potency > 0 {
			flavors = append(flavors, core.BerryFlavor{
				Flavor:  f.Flavor.Name,
				Potency: f.Potency,
			})
		}
	}

	// Fetch item sprite and effect
	itemSprite := ""
	effect := ""
	if raw.Item.Name != "" {
		itemURL := fmt.Sprintf("%s/item/%s", c.baseURL, raw.Item.Name)
		var item apiItem
		if err := fetchOne(c, itemURL, &item); err == nil {
			itemSprite = item.Sprites.Default
			for _, e := range item.EffectEntries {
				if e.Language.Name == "en" {
					effect = e.ShortEffect
					break
				}
			}
		}
	}

	return core.Berry{
		ID:               raw.ID,
		Name:             raw.Name,
		GrowthTime:       raw.GrowthTime,
		MaxHarvest:       raw.MaxHarvest,
		NaturalGiftPower: raw.NaturalGiftPower,
		Size:             raw.Size,
		Smoothness:       raw.Smoothness,
		SoilDryness:      raw.SoilDryness,
		Firmness:         raw.Firmness.Name,
		Flavors:          flavors,
		NaturalGiftType:  raw.NaturalGiftType.Name,
		ItemName:         raw.Item.Name,
		ItemSprite:       itemSprite,
		Effect:           effect,
	}, nil
}
