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
