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

// EstimateEVFromHP calcula el EV estimado para HP dado el stat actual
func EstimateEVFromHP(currentHP, base, iv, level int) int {
	// HP = floor((2*Base + IV + floor(EV/4)) * Level / 100) + Level + 10
	// Despejando: EV = ((HP - Level - 10) * 100 / Level - 2*Base - IV) * 4
	ev := ((currentHP - level - 10) * 100 / level - 2*base - iv) * 4
	return clampEV(ev)
}

// EstimateEVFromStat calcula el EV estimado para un stat (no HP)
func EstimateEVFromStat(currentStat, base, iv, level int, natureModifier float64) int {
	// Stat = floor((floor((2*Base + IV + floor(EV/4)) * Level / 100) + 5) * Nature)
	// Despejando aproximado: EV = ((Stat/Nature - 5) * 100 / Level - 2*Base - IV) * 4
	statBeforeNature := float64(currentStat) / natureModifier
	ev := int((statBeforeNature-5)*100/float64(level)-float64(2*base)-float64(iv)) * 4
	return clampEV(ev)
}

// EstimateEVRangeFromHP calcula el rango posible de EVs para HP
func EstimateEVRangeFromHP(currentHP, base, iv, level int) StatRange {
	// Buscar todos los EVs que producen el mismo stat
	minEV := -1
	maxEV := -1

	for ev := 0; ev <= 252; ev++ {
		calculated := CalculateHP(base, iv, ev, level)
		if calculated == currentHP {
			if minEV == -1 {
				minEV = ev
			}
			maxEV = ev
		}
	}

	if minEV == -1 {
		// No se encontró match exacto, usar estimación
		estimated := EstimateEVFromHP(currentHP, base, iv, level)
		return StatRange{Min: estimated, Max: estimated}
	}

	return StatRange{Min: minEV, Max: maxEV}
}

// EstimateEVRangeFromStat calcula el rango posible de EVs para un stat
func EstimateEVRangeFromStat(currentStat, base, iv, level int, natureModifier float64) StatRange {
	minEV := -1
	maxEV := -1

	for ev := 0; ev <= 252; ev++ {
		calculated := CalculateStat(base, iv, ev, level, natureModifier)
		if calculated == currentStat {
			if minEV == -1 {
				minEV = ev
			}
			maxEV = ev
		}
	}

	if minEV == -1 {
		estimated := EstimateEVFromStat(currentStat, base, iv, level, natureModifier)
		return StatRange{Min: estimated, Max: estimated}
	}

	return StatRange{Min: minEV, Max: maxEV}
}

// clampEV asegura que el EV esté en rango válido
func clampEV(ev int) int {
	if ev < 0 {
		return 0
	}
	if ev > 252 {
		return 252
	}
	return ev
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
