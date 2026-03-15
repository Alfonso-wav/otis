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

func deterministicRng(seq []int) func(int) int {
	i := 0
	return func(n int) int {
		v := seq[i%len(seq)] % n
		i++
		return v
	}
}

func TestGenerateRandomTeamMember(t *testing.T) {
	pokemon := Pokemon{
		Name: "pikachu",
		Moves: []PokemonMoveEntry{
			{Name: "thunderbolt"}, {Name: "quick-attack"}, {Name: "iron-tail"},
			{Name: "electro-ball"}, {Name: "thunder"}, {Name: "volt-tackle"},
		},
	}
	rng := deterministicRng([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	member := GenerateRandomTeamMember(pokemon, rng)

	if member.PokemonName != "pikachu" {
		t.Errorf("expected pikachu, got %s", member.PokemonName)
	}
	if member.Level != 50 {
		t.Errorf("expected level 50, got %d", member.Level)
	}
	if len(member.Moves) != 4 {
		t.Errorf("expected 4 moves, got %d", len(member.Moves))
	}
	// Check no duplicate moves
	seen := make(map[string]bool)
	for _, m := range member.Moves {
		if seen[m] {
			t.Errorf("duplicate move: %s", m)
		}
		seen[m] = true
	}
	// Check EVs valid
	totalEVs := member.EVs.HP + member.EVs.Attack + member.EVs.Defense +
		member.EVs.SpAttack + member.EVs.SpDefense + member.EVs.Speed
	if totalEVs > MaxTotalEVs {
		t.Errorf("total EVs %d exceeds max %d", totalEVs, MaxTotalEVs)
	}
	evFields := []int{member.EVs.HP, member.EVs.Attack, member.EVs.Defense,
		member.EVs.SpAttack, member.EVs.SpDefense, member.EVs.Speed}
	for _, ev := range evFields {
		if ev < 0 || ev > MaxSingleEV {
			t.Errorf("EV %d out of range [0, %d]", ev, MaxSingleEV)
		}
	}
	// Check IVs all 31
	if member.IVs.HP != 31 || member.IVs.Attack != 31 || member.IVs.Speed != 31 {
		t.Errorf("expected all IVs 31, got %+v", member.IVs)
	}
	// Check nature is valid
	if _, ok := Natures[member.Nature]; !ok {
		t.Errorf("invalid nature: %s", member.Nature)
	}
}

func TestFillTeamRandom(t *testing.T) {
	available := []Pokemon{
		{Name: "pikachu", Moves: []PokemonMoveEntry{{Name: "thunderbolt"}, {Name: "quick-attack"}, {Name: "iron-tail"}, {Name: "electro-ball"}}},
		{Name: "charmander", Moves: []PokemonMoveEntry{{Name: "flamethrower"}, {Name: "scratch"}, {Name: "ember"}, {Name: "fire-blast"}}},
		{Name: "bulbasaur", Moves: []PokemonMoveEntry{{Name: "vine-whip"}, {Name: "razor-leaf"}, {Name: "solar-beam"}, {Name: "tackle"}}},
		{Name: "squirtle", Moves: []PokemonMoveEntry{{Name: "water-gun"}, {Name: "hydro-pump"}, {Name: "bubble"}, {Name: "surf"}}},
		{Name: "eevee", Moves: []PokemonMoveEntry{{Name: "tackle"}, {Name: "swift"}, {Name: "bite"}, {Name: "shadow-ball"}}},
		{Name: "jigglypuff", Moves: []PokemonMoveEntry{{Name: "sing"}, {Name: "pound"}, {Name: "double-slap"}, {Name: "hyper-voice"}}},
		{Name: "meowth", Moves: []PokemonMoveEntry{{Name: "scratch"}, {Name: "pay-day"}, {Name: "bite"}, {Name: "slash"}}},
	}
	rng := deterministicRng([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20})

	t.Run("fill empty team to 6", func(t *testing.T) {
		team := Team{Name: "Test"}
		got := FillTeamRandom(team, available, rng)
		if len(got.Members) != 6 {
			t.Errorf("expected 6 members, got %d", len(got.Members))
		}
		// Check no duplicate pokemon
		seen := make(map[string]bool)
		for _, m := range got.Members {
			if seen[m.PokemonName] {
				t.Errorf("duplicate pokemon: %s", m.PokemonName)
			}
			seen[m.PokemonName] = true
		}
	})

	t.Run("fill partial team", func(t *testing.T) {
		team := Team{
			Name:    "Partial",
			Members: []TeamMember{{PokemonName: "pikachu", Level: 50}, {PokemonName: "charmander", Level: 50}},
		}
		got := FillTeamRandom(team, available, rng)
		if len(got.Members) != 6 {
			t.Errorf("expected 6 members, got %d", len(got.Members))
		}
		// Check pikachu and charmander are not duplicated
		count := make(map[string]int)
		for _, m := range got.Members {
			count[m.PokemonName]++
		}
		if count["pikachu"] != 1 || count["charmander"] != 1 {
			t.Error("existing members should not be duplicated")
		}
	})

	t.Run("full team unchanged", func(t *testing.T) {
		team := Team{Name: "Full", Members: make([]TeamMember, 6)}
		got := FillTeamRandom(team, available, rng)
		if len(got.Members) != 6 {
			t.Errorf("expected 6 members unchanged, got %d", len(got.Members))
		}
	})
}
