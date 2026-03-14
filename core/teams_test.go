package core

import (
	"testing"
)

func TestValidateTeam(t *testing.T) {
	tests := []struct {
		name    string
		team    Team
		wantErr bool
	}{
		{"valid team", Team{Name: "My Team", Members: []TeamMember{{PokemonName: "pikachu", Level: 50}}}, false},
		{"empty name", Team{Name: "", Members: nil}, true},
		{"too many members", Team{Name: "Big", Members: make([]TeamMember, 7)}, true},
		{"max members ok", Team{Name: "Full", Members: make([]TeamMember, 6)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTeam(tt.team)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTeam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTeamMember(t *testing.T) {
	validMember := TeamMember{
		PokemonName: "pikachu",
		Moves:       []string{"thunderbolt"},
		Level:       50,
		Nature:      "Adamant",
		IVs:         Stats{HP: 31, Attack: 31, Defense: 31, SpAttack: 31, SpDefense: 31, Speed: 31},
		EVs:         Stats{HP: 0, Attack: 252, Defense: 0, SpAttack: 0, SpDefense: 4, Speed: 252},
	}

	tests := []struct {
		name    string
		member  TeamMember
		wantErr bool
	}{
		{"valid", validMember, false},
		{"empty pokemon name", TeamMember{PokemonName: "", Level: 50}, true},
		{"level too low", TeamMember{PokemonName: "pikachu", Level: 0}, true},
		{"level too high", TeamMember{PokemonName: "pikachu", Level: 101}, true},
		{"ev over 252", TeamMember{PokemonName: "pikachu", Level: 50, EVs: Stats{HP: 253}}, true},
		{"total evs over 510", TeamMember{PokemonName: "pikachu", Level: 50, EVs: Stats{HP: 252, Attack: 252, Defense: 252}}, true},
		{"iv over 31", TeamMember{PokemonName: "pikachu", Level: 50, IVs: Stats{HP: 32}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTeamMember(tt.member)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTeamMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddMemberToTeam(t *testing.T) {
	member := TeamMember{PokemonName: "pikachu", Level: 50, Nature: "Hardy"}

	t.Run("add to empty team", func(t *testing.T) {
		team := Team{Name: "Test"}
		got, err := AddMemberToTeam(team, member)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Members) != 1 {
			t.Errorf("expected 1 member, got %d", len(got.Members))
		}
		if got.Members[0].PokemonName != "pikachu" {
			t.Errorf("expected pikachu, got %s", got.Members[0].PokemonName)
		}
	})

	t.Run("team full", func(t *testing.T) {
		team := Team{Name: "Full", Members: make([]TeamMember, 6)}
		_, err := AddMemberToTeam(team, member)
		if err != ErrTeamFull {
			t.Errorf("expected ErrTeamFull, got %v", err)
		}
	})

	t.Run("does not mutate original", func(t *testing.T) {
		team := Team{Name: "Test", Members: []TeamMember{{PokemonName: "charmander", Level: 50}}}
		got, _ := AddMemberToTeam(team, member)
		if len(team.Members) != 1 {
			t.Error("original team was mutated")
		}
		if len(got.Members) != 2 {
			t.Errorf("expected 2 members, got %d", len(got.Members))
		}
	})
}

func TestRemoveMemberFromTeam(t *testing.T) {
	team := Team{
		Name: "Test",
		Members: []TeamMember{
			{PokemonName: "pikachu", Level: 50},
			{PokemonName: "charmander", Level: 50},
			{PokemonName: "bulbasaur", Level: 50},
		},
	}

	t.Run("remove middle", func(t *testing.T) {
		got, err := RemoveMemberFromTeam(team, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got.Members) != 2 {
			t.Errorf("expected 2 members, got %d", len(got.Members))
		}
		if got.Members[0].PokemonName != "pikachu" || got.Members[1].PokemonName != "bulbasaur" {
			t.Error("wrong members remaining")
		}
	})

	t.Run("index out of range", func(t *testing.T) {
		_, err := RemoveMemberFromTeam(team, 5)
		if err != ErrIndexOutOfRange {
			t.Errorf("expected ErrIndexOutOfRange, got %v", err)
		}
	})

	t.Run("negative index", func(t *testing.T) {
		_, err := RemoveMemberFromTeam(team, -1)
		if err != ErrIndexOutOfRange {
			t.Errorf("expected ErrIndexOutOfRange, got %v", err)
		}
	})
}
