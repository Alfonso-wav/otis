package core

import "math"

// Natures contiene todas las 25 naturalezas de Pokémon
var Natures = map[string]Nature{
	"Adamant":  {Name: "Adamant", IncreasedStat: "attack", DecreasedStat: "spAttack"},
	"Bashful":  {Name: "Bashful", IncreasedStat: "", DecreasedStat: ""},
	"Bold":     {Name: "Bold", IncreasedStat: "defense", DecreasedStat: "attack"},
	"Brave":    {Name: "Brave", IncreasedStat: "attack", DecreasedStat: "speed"},
	"Calm":     {Name: "Calm", IncreasedStat: "spDefense", DecreasedStat: "attack"},
	"Careful":  {Name: "Careful", IncreasedStat: "spDefense", DecreasedStat: "spAttack"},
	"Docile":   {Name: "Docile", IncreasedStat: "", DecreasedStat: ""},
	"Gentle":   {Name: "Gentle", IncreasedStat: "spDefense", DecreasedStat: "defense"},
	"Hardy":    {Name: "Hardy", IncreasedStat: "", DecreasedStat: ""},
	"Hasty":    {Name: "Hasty", IncreasedStat: "speed", DecreasedStat: "defense"},
	"Impish":   {Name: "Impish", IncreasedStat: "defense", DecreasedStat: "spAttack"},
	"Jolly":    {Name: "Jolly", IncreasedStat: "speed", DecreasedStat: "spAttack"},
	"Lax":      {Name: "Lax", IncreasedStat: "defense", DecreasedStat: "spDefense"},
	"Lonely":   {Name: "Lonely", IncreasedStat: "attack", DecreasedStat: "defense"},
	"Mild":     {Name: "Mild", IncreasedStat: "spAttack", DecreasedStat: "defense"},
	"Modest":   {Name: "Modest", IncreasedStat: "spAttack", DecreasedStat: "attack"},
	"Naive":    {Name: "Naive", IncreasedStat: "speed", DecreasedStat: "spDefense"},
	"Naughty":  {Name: "Naughty", IncreasedStat: "attack", DecreasedStat: "spDefense"},
	"Quiet":    {Name: "Quiet", IncreasedStat: "spAttack", DecreasedStat: "speed"},
	"Quirky":   {Name: "Quirky", IncreasedStat: "", DecreasedStat: ""},
	"Rash":     {Name: "Rash", IncreasedStat: "spAttack", DecreasedStat: "spDefense"},
	"Relaxed":  {Name: "Relaxed", IncreasedStat: "defense", DecreasedStat: "speed"},
	"Sassy":    {Name: "Sassy", IncreasedStat: "spDefense", DecreasedStat: "speed"},
	"Serious":  {Name: "Serious", IncreasedStat: "", DecreasedStat: ""},
	"Timid":    {Name: "Timid", IncreasedStat: "speed", DecreasedStat: "attack"},
}

// GetNatureModifier devuelve el modificador de naturaleza para un stat
func GetNatureModifier(nature Nature, statName string) float64 {
	if nature.IncreasedStat == statName {
		return 1.1
	}
	if nature.DecreasedStat == statName {
		return 0.9
	}
	return 1.0
}

// CalculateHP calcula el HP final usando la fórmula de Gen III+
func CalculateHP(base, iv, ev, level int) int {
	return int(math.Floor(float64((2*base+iv+ev/4)*level)/100)) + level + 10
}

// CalculateStat calcula un stat (no HP) usando la fórmula de Gen III+
func CalculateStat(base, iv, ev, level int, natureModifier float64) int {
	innerCalc := math.Floor(float64((2*base+iv+ev/4)*level)/100) + 5
	return int(math.Floor(innerCalc * natureModifier))
}

// CalculateAllStats calcula todos los stats finales
func CalculateAllStats(baseStats Stats, ivs Stats, evs Stats, level int, nature Nature) Stats {
	return Stats{
		HP:        CalculateHP(baseStats.HP, ivs.HP, evs.HP, level),
		Attack:    CalculateStat(baseStats.Attack, ivs.Attack, evs.Attack, level, GetNatureModifier(nature, "attack")),
		Defense:   CalculateStat(baseStats.Defense, ivs.Defense, evs.Defense, level, GetNatureModifier(nature, "defense")),
		SpAttack:  CalculateStat(baseStats.SpAttack, ivs.SpAttack, evs.SpAttack, level, GetNatureModifier(nature, "spAttack")),
		SpDefense: CalculateStat(baseStats.SpDefense, ivs.SpDefense, evs.SpDefense, level, GetNatureModifier(nature, "spDefense")),
		Speed:     CalculateStat(baseStats.Speed, ivs.Speed, evs.Speed, level, GetNatureModifier(nature, "speed")),
	}
}

// PokemonToBaseStats convierte los stats de un Pokemon a la estructura Stats
func PokemonToBaseStats(p Pokemon) Stats {
	stats := Stats{}
	for _, s := range p.Stats {
		switch s.Name {
		case "hp":
			stats.HP = s.BaseStat
		case "attack":
			stats.Attack = s.BaseStat
		case "defense":
			stats.Defense = s.BaseStat
		case "special-attack":
			stats.SpAttack = s.BaseStat
		case "special-defense":
			stats.SpDefense = s.BaseStat
		case "speed":
			stats.Speed = s.BaseStat
		}
	}
	return stats
}

// DefaultIVs retorna IVs perfectos (31 en todo)
func DefaultIVs() Stats {
	return Stats{
		HP:        31,
		Attack:    31,
		Defense:   31,
		SpAttack:  31,
		SpDefense: 31,
		Speed:     31,
	}
}

// TotalEVs suma todos los EVs de un Stats
func TotalEVs(evs Stats) int {
	return evs.HP + evs.Attack + evs.Defense + evs.SpAttack + evs.SpDefense + evs.Speed
}
