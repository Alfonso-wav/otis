package core

import (
	"errors"
	"fmt"
)

const MaxTeamMembers = 6
const MaxTotalEVs = 510
const MaxSingleEV = 252
const MaxIV = 31
const MinLevel = 1
const MaxLevel = 100

var (
	ErrTeamNameEmpty = errors.New("team name cannot be empty")
	ErrTeamFull      = errors.New("team already has 6 members")
	ErrIndexOutOfRange = errors.New("member index out of range")
)

// ValidateTeam checks that a team has a non-empty name and at most 6 members.
func ValidateTeam(team Team) error {
	if team.Name == "" {
		return ErrTeamNameEmpty
	}
	if len(team.Members) > MaxTeamMembers {
		return fmt.Errorf("team has %d members, max is %d", len(team.Members), MaxTeamMembers)
	}
	return nil
}

// ValidateTeamMember checks that a member has valid EVs, IVs, and level.
func ValidateTeamMember(member TeamMember) error {
	if member.PokemonName == "" {
		return errors.New("pokemon name cannot be empty")
	}
	if member.Level < MinLevel || member.Level > MaxLevel {
		return fmt.Errorf("level must be between %d and %d", MinLevel, MaxLevel)
	}
	totalEVs := member.EVs.HP + member.EVs.Attack + member.EVs.Defense +
		member.EVs.SpAttack + member.EVs.SpDefense + member.EVs.Speed
	if totalEVs > MaxTotalEVs {
		return fmt.Errorf("total EVs %d exceeds max %d", totalEVs, MaxTotalEVs)
	}
	evFields := []int{member.EVs.HP, member.EVs.Attack, member.EVs.Defense,
		member.EVs.SpAttack, member.EVs.SpDefense, member.EVs.Speed}
	for _, ev := range evFields {
		if ev < 0 || ev > MaxSingleEV {
			return fmt.Errorf("each EV must be between 0 and %d", MaxSingleEV)
		}
	}
	ivFields := []int{member.IVs.HP, member.IVs.Attack, member.IVs.Defense,
		member.IVs.SpAttack, member.IVs.SpDefense, member.IVs.Speed}
	for _, iv := range ivFields {
		if iv < 0 || iv > MaxIV {
			return fmt.Errorf("each IV must be between 0 and %d", MaxIV)
		}
	}
	return nil
}

// AddMemberToTeam returns a new Team with the member added.
func AddMemberToTeam(team Team, member TeamMember) (Team, error) {
	if len(team.Members) >= MaxTeamMembers {
		return team, ErrTeamFull
	}
	if err := ValidateTeamMember(member); err != nil {
		return team, err
	}
	newMembers := make([]TeamMember, len(team.Members), len(team.Members)+1)
	copy(newMembers, team.Members)
	newMembers = append(newMembers, member)
	return Team{Name: team.Name, Members: newMembers}, nil
}

// RemoveMemberFromTeam returns a new Team without the member at the given index.
func RemoveMemberFromTeam(team Team, index int) (Team, error) {
	if index < 0 || index >= len(team.Members) {
		return team, ErrIndexOutOfRange
	}
	newMembers := make([]TeamMember, 0, len(team.Members)-1)
	newMembers = append(newMembers, team.Members[:index]...)
	newMembers = append(newMembers, team.Members[index+1:]...)
	return Team{Name: team.Name, Members: newMembers}, nil
}

// natureNames returns a sorted list of all nature names.
func natureNames() []string {
	names := make([]string, 0, len(Natures))
	for k := range Natures {
		names = append(names, k)
	}
	return names
}

// GenerateRandomTeamMember creates a random TeamMember from a Pokemon's data.
func GenerateRandomTeamMember(pokemon Pokemon, rng func(int) int) TeamMember {
	// Nature
	names := natureNames()
	nature := names[rng(len(names))]

	// Moves: prefer damaging moves (power > 0), pick up to 4
	var damagingMoves, otherMoves []string
	for _, m := range pokemon.Moves {
		damagingMoves = append(damagingMoves, m.Name)
		_ = otherMoves // all moves go into one pool for simplicity
	}
	movePool := damagingMoves
	moves := pickRandomUnique(movePool, 4, rng)

	// EVs: distribute 510 across 6 stats, max 252 each
	evs := randomEVSpread(rng)

	return TeamMember{
		PokemonName: pokemon.Name,
		Moves:       moves,
		Level:       50,
		Nature:      nature,
		IVs:         Stats{HP: 31, Attack: 31, Defense: 31, SpAttack: 31, SpDefense: 31, Speed: 31},
		EVs:         evs,
	}
}

// randomEVSpread distributes 510 EV points across 6 stats (max 252 each).
func randomEVSpread(rng func(int) int) Stats {
	stats := [6]int{}
	remaining := MaxTotalEVs
	for i := 0; i < 6 && remaining > 0; i++ {
		max := remaining
		if max > MaxSingleEV {
			max = MaxSingleEV
		}
		if i == 5 {
			stats[i] = remaining
			if stats[i] > MaxSingleEV {
				stats[i] = MaxSingleEV
			}
			break
		}
		stats[i] = rng(max + 1)
		remaining -= stats[i]
	}
	return Stats{
		HP: stats[0], Attack: stats[1], Defense: stats[2],
		SpAttack: stats[3], SpDefense: stats[4], Speed: stats[5],
	}
}

// pickRandomUnique selects up to n unique items from a slice.
func pickRandomUnique(items []string, n int, rng func(int) int) []string {
	if len(items) == 0 {
		return nil
	}
	if n > len(items) {
		n = len(items)
	}
	// Fisher-Yates partial shuffle
	pool := make([]string, len(items))
	copy(pool, items)
	for i := 0; i < n; i++ {
		j := i + rng(len(pool)-i)
		pool[i], pool[j] = pool[j], pool[i]
	}
	return pool[:n]
}

// FillTeamRandom fills empty slots (up to 6) with random Pokemon members.
func FillTeamRandom(team Team, availablePokemon []Pokemon, rng func(int) int) Team {
	slotsNeeded := MaxTeamMembers - len(team.Members)
	if slotsNeeded <= 0 || len(availablePokemon) == 0 {
		return team
	}

	// Collect existing pokemon names to avoid duplicates
	existing := make(map[string]bool)
	for _, m := range team.Members {
		existing[m.PokemonName] = true
	}

	// Filter out already-used pokemon
	var candidates []Pokemon
	for _, p := range availablePokemon {
		if !existing[p.Name] {
			candidates = append(candidates, p)
		}
	}

	if len(candidates) == 0 {
		return team
	}

	// Shuffle candidates and pick up to slotsNeeded
	for i := range candidates {
		j := i + rng(len(candidates)-i)
		candidates[i], candidates[j] = candidates[j], candidates[i]
	}
	if slotsNeeded > len(candidates) {
		slotsNeeded = len(candidates)
	}

	newMembers := make([]TeamMember, len(team.Members), len(team.Members)+slotsNeeded)
	copy(newMembers, team.Members)
	for i := 0; i < slotsNeeded; i++ {
		newMembers = append(newMembers, GenerateRandomTeamMember(candidates[i], rng))
	}

	return Team{Name: team.Name, Members: newMembers}
}
