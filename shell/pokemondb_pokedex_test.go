package shell

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

const fixturePokedexHTML = `
<html><body>
<table id="pokedex">
<thead>
<tr><th>#</th><th>Name</th><th>Type</th><th>Total</th><th>HP</th><th>Atk</th><th>Def</th><th>Sp.Atk</th><th>Sp.Def</th><th>Speed</th></tr>
</thead>
<tbody>
<tr>
  <td>0001</td>
  <td><a href="/pokedex/bulbasaur">Bulbasaur</a></td>
  <td><a href="/type/grass">Grass</a> <a href="/type/poison">Poison</a></td>
  <td>318</td><td>45</td><td>49</td><td>49</td><td>65</td><td>65</td><td>45</td>
</tr>
<tr>
  <td>0006</td>
  <td><a href="/pokedex/charizard">Charizard</a></td>
  <td><a href="/type/fire">Fire</a> <a href="/type/flying">Flying</a></td>
  <td>534</td><td>78</td><td>84</td><td>78</td><td>109</td><td>85</td><td>100</td>
</tr>
<tr>
  <td>0025</td>
  <td><a href="/pokedex/pikachu">Pikachu</a></td>
  <td><a href="/type/electric">Electric</a></td>
  <td>320</td><td>35</td><td>55</td><td>40</td><td>50</td><td>50</td><td>90</td>
</tr>
</tbody>
</table>
</body></html>
`

func TestParsePokedexTable(t *testing.T) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fixturePokedexHTML))
	if err != nil {
		t.Fatalf("parsing fixture HTML: %v", err)
	}

	entries, err := parsePokedexTable(doc)
	if err != nil {
		t.Fatalf("parsePokedexTable error: %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	// Bulbasaur
	bulba := entries[0]
	if bulba.ID != 1 {
		t.Errorf("Bulbasaur ID = %d, want 1", bulba.ID)
	}
	if bulba.Name != "Bulbasaur" {
		t.Errorf("Bulbasaur Name = %q, want %q", bulba.Name, "Bulbasaur")
	}
	if len(bulba.Types) != 2 || bulba.Types[0] != "grass" || bulba.Types[1] != "poison" {
		t.Errorf("Bulbasaur Types = %v, want [grass poison]", bulba.Types)
	}
	if bulba.Total != 318 {
		t.Errorf("Bulbasaur Total = %d, want 318", bulba.Total)
	}
	if bulba.HP != 45 || bulba.Attack != 49 || bulba.Defense != 49 {
		t.Errorf("Bulbasaur stats wrong: HP=%d Atk=%d Def=%d", bulba.HP, bulba.Attack, bulba.Defense)
	}
	if bulba.SpAtk != 65 || bulba.SpDef != 65 || bulba.Speed != 45 {
		t.Errorf("Bulbasaur stats wrong: SpA=%d SpD=%d Spe=%d", bulba.SpAtk, bulba.SpDef, bulba.Speed)
	}

	// Pikachu (single type)
	pika := entries[2]
	if pika.ID != 25 {
		t.Errorf("Pikachu ID = %d, want 25", pika.ID)
	}
	if pika.Name != "Pikachu" {
		t.Errorf("Pikachu Name = %q, want %q", pika.Name, "Pikachu")
	}
	if len(pika.Types) != 1 || pika.Types[0] != "electric" {
		t.Errorf("Pikachu Types = %v, want [electric]", pika.Types)
	}
	if pika.Speed != 90 {
		t.Errorf("Pikachu Speed = %d, want 90", pika.Speed)
	}

	// Charizard
	chari := entries[1]
	if chari.ID != 6 {
		t.Errorf("Charizard ID = %d, want 6", chari.ID)
	}
	if chari.Total != 534 {
		t.Errorf("Charizard Total = %d, want 534", chari.Total)
	}
}

func TestParsePokedexTableEmpty(t *testing.T) {
	html := `<html><body><table id="pokedex"><tbody></tbody></table></body></html>`
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	_, err := parsePokedexTable(doc)
	if err == nil {
		t.Error("expected error for empty table, got nil")
	}
}
