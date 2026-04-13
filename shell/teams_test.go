package shell

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alfon/pokemon-app/core"
)

func TestFileTeamStorage_CRUD(t *testing.T) {
	dir := t.TempDir()
	storage := NewFileTeamStorage(dir)

	team := core.Team{
		Name: "My Team",
		Members: []core.TeamMember{
			{
				PokemonName: "pikachu",
				Moves:       []string{"thunderbolt", "quick-attack"},
				Level:       50,
				Nature:      "Jolly",
				IVs:         core.Stats{HP: 31, Attack: 31, Defense: 31, SpAttack: 31, SpDefense: 31, Speed: 31},
				EVs:         core.Stats{HP: 0, Attack: 252, Defense: 0, SpAttack: 0, SpDefense: 4, Speed: 252},
			},
		},
	}

	// Save
	if err := storage.SaveTeam(team); err != nil {
		t.Fatalf("SaveTeam() error = %v", err)
	}

	// Get
	got, err := storage.GetTeam("My Team")
	if err != nil {
		t.Fatalf("GetTeam() error = %v", err)
	}
	if got.Name != team.Name {
		t.Errorf("name = %q, want %q", got.Name, team.Name)
	}
	if len(got.Members) != 1 {
		t.Fatalf("members count = %d, want 1", len(got.Members))
	}
	if got.Members[0].PokemonName != "pikachu" {
		t.Errorf("member name = %q, want pikachu", got.Members[0].PokemonName)
	}

	// List
	teams, err := storage.ListTeams()
	if err != nil {
		t.Fatalf("ListTeams() error = %v", err)
	}
	if len(teams) != 1 {
		t.Errorf("teams count = %d, want 1", len(teams))
	}

	// Delete
	if err := storage.DeleteTeam("My Team"); err != nil {
		t.Fatalf("DeleteTeam() error = %v", err)
	}
	teams, _ = storage.ListTeams()
	if len(teams) != 0 {
		t.Errorf("teams count after delete = %d, want 0", len(teams))
	}
}

func TestFileTeamStorage_SaveValidation(t *testing.T) {
	dir := t.TempDir()
	storage := NewFileTeamStorage(dir)

	err := storage.SaveTeam(core.Team{Name: ""})
	if err == nil {
		t.Error("expected validation error for empty name")
	}
}

func TestFileTeamStorage_GetNonExistent(t *testing.T) {
	dir := t.TempDir()
	storage := NewFileTeamStorage(dir)

	_, err := storage.GetTeam("nonexistent")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}

func TestFileTeamStorage_RoundTripWithAbility(t *testing.T) {
	dir := t.TempDir()
	storage := NewFileTeamStorage(dir)
	team := core.Team{
		Name: "weather-squad",
		Members: []core.TeamMember{
			{PokemonName: "politoed", Level: 50, Nature: "Hardy", Ability: "drizzle"},
			{PokemonName: "kingdra", Level: 50, Nature: "Hardy", Ability: "swift-swim"},
			{PokemonName: "ludicolo", Level: 50, Nature: "Hardy"}, // no ability
		},
	}
	if err := storage.SaveTeam(team); err != nil {
		t.Fatalf("SaveTeam: %v", err)
	}
	got, err := storage.GetTeam("weather-squad")
	if err != nil {
		t.Fatalf("GetTeam: %v", err)
	}
	if got.Members[0].Ability != "drizzle" {
		t.Errorf("member[0] ability: want drizzle, got %q", got.Members[0].Ability)
	}
	if got.Members[1].Ability != "swift-swim" {
		t.Errorf("member[1] ability: want swift-swim, got %q", got.Members[1].Ability)
	}
	if got.Members[2].Ability != "" {
		t.Errorf("member[2] ability: want empty, got %q", got.Members[2].Ability)
	}
}

func TestFileTeamStorage_LoadLegacyFileWithoutAbility(t *testing.T) {
	dir := t.TempDir()
	legacy := `{
  "name": "legacy",
  "members": [
    {"pokemonName": "pikachu", "level": 50, "nature": "Hardy", "moves": [], "ivs": {}, "evs": {}}
  ]
}`
	if err := os.WriteFile(filepath.Join(dir, "legacy.json"), []byte(legacy), 0644); err != nil {
		t.Fatal(err)
	}
	storage := NewFileTeamStorage(dir)
	got, err := storage.GetTeam("legacy")
	if err != nil {
		t.Fatalf("GetTeam legacy: %v", err)
	}
	if got.Members[0].Ability != "" {
		t.Errorf("legacy member ability: want empty, got %q", got.Members[0].Ability)
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"My Team", "my-team"},
		{"  Hello World  ", "hello-world"},
		{"UPPER CASE", "upper-case"},
		{"special!@#chars", "specialchars"},
		{"", "team"},
	}
	for _, tt := range tests {
		got := slugify(tt.input)
		if got != tt.want {
			t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
