package core

// AggregateEncounters merges encounters from multiple location areas,
// deduplicating by PokemonName and keeping the highest MaxChance.
func AggregateEncounters(areas []LocationArea) []PokemonEncounter {
	best := make(map[string]int)
	for _, area := range areas {
		for _, enc := range area.PokemonEncounters {
			if enc.MaxChance > best[enc.PokemonName] {
				best[enc.PokemonName] = enc.MaxChance
			}
		}
	}

	result := make([]PokemonEncounter, 0, len(best))
	for name, chance := range best {
		if chance > 100 {
			chance = 100
		}
		result = append(result, PokemonEncounter{
			PokemonName: name,
			MaxChance:   chance,
		})
	}
	return result
}
