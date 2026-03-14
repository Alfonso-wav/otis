package shell

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alfon/pokemon-app/core"
)

// FetchPokedex extrae la tabla completa de Pokémon de pokemondb.net/pokedex/all.
func (c *PokemonDBClient) FetchPokedex() ([]core.PokedexDBEntry, error) {
	doc, err := c.fetchPage("/pokedex/all")
	if err != nil {
		return nil, fmt.Errorf("fetching pokedex: %w", err)
	}

	return parsePokedexTable(doc)
}

// parsePokedexTable extrae los datos de la tabla principal del pokedex.
func parsePokedexTable(doc *goquery.Document) ([]core.PokedexDBEntry, error) {
	var entries []core.PokedexDBEntry

	// La tabla principal del pokedex tiene id="pokedex"
	table := doc.Find("table#pokedex")
	if table.Length() == 0 {
		// Fallback: buscar la primera tabla grande con datos de Pokemon
		table = doc.Find("table").First()
	}

	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		entry, err := parsePokedexRow(row)
		if err != nil {
			return // skip rows that don't parse
		}
		entries = append(entries, entry)
	})

	if len(entries) == 0 {
		return nil, fmt.Errorf("no pokemon entries found in pokedex table")
	}

	return entries, nil
}

// parsePokedexRow extrae un PokedexDBEntry de una fila de la tabla.
func parsePokedexRow(row *goquery.Selection) (core.PokedexDBEntry, error) {
	cells := row.Find("td")
	if cells.Length() < 10 {
		return core.PokedexDBEntry{}, fmt.Errorf("row has %d cells, expected at least 10", cells.Length())
	}

	// Cell 0: numero nacional (con clase cell-num normalmente)
	numText := strings.TrimSpace(cells.Eq(0).Text())
	id, err := strconv.Atoi(numText)
	if err != nil {
		return core.PokedexDBEntry{}, fmt.Errorf("parsing id %q: %w", numText, err)
	}

	// Cell 1: nombre (dentro de un <a>)
	nameCell := cells.Eq(1)
	name := strings.TrimSpace(nameCell.Find("a").First().Text())
	if name == "" {
		name = strings.TrimSpace(nameCell.Text())
	}

	// Cell 2: tipos (anchors con clase type-icon)
	var types []string
	cells.Eq(2).Find("a").Each(func(_ int, a *goquery.Selection) {
		t := strings.TrimSpace(strings.ToLower(a.Text()))
		if t != "" {
			types = append(types, t)
		}
	})

	// Cells 3-8: Total, HP, Attack, Defense, Sp.Atk, Sp.Def, Speed
	// O cells 4-9 si hay una columna extra
	statStart := 3
	total := parseIntSafe(cells.Eq(statStart).Text())
	hp := parseIntSafe(cells.Eq(statStart + 1).Text())
	attack := parseIntSafe(cells.Eq(statStart + 2).Text())
	defense := parseIntSafe(cells.Eq(statStart + 3).Text())
	spAtk := parseIntSafe(cells.Eq(statStart + 4).Text())
	spDef := parseIntSafe(cells.Eq(statStart + 5).Text())
	speed := parseIntSafe(cells.Eq(statStart + 6).Text())

	return core.PokedexDBEntry{
		ID:      id,
		Name:    name,
		Types:   types,
		Total:   total,
		HP:      hp,
		Attack:  attack,
		Defense: defense,
		SpAtk:   spAtk,
		SpDef:   spDef,
		Speed:   speed,
	}, nil
}

func parseIntSafe(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
