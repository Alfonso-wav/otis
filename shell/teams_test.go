package shell

import (
	"os"
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
