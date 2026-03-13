package core

// StatComparison representa la comparación de un stat individual entre dos Pokémon.
type StatComparison struct {
	Name   string
	StatA  int
	StatB  int
	Diff   int    // StatA - StatB (positivo = A gana, negativo = B gana)
	Winner string // "a", "b", "tie"
}

// PokemonComparison es el resultado de comparar dos Pokémon.
type PokemonComparison struct {
	PokemonA Pokemon
	PokemonB Pokemon
	Stats    []StatComparison
	TotalA   int
	TotalB   int
	Winner   string // "a", "b", "tie"
}

// ComparePokemons compara los stats base de dos Pokémon y devuelve el resultado detallado.
// Es una función pura: no produce efectos secundarios.
func ComparePokemons(a, b Pokemon) PokemonComparison {
	statMap := func(p Pokemon) map[string]int {
		m := make(map[string]int, len(p.Stats))
		for _, s := range p.Stats {
			m[s.Name] = s.BaseStat
		}
		return m
	}

	statsA := statMap(a)
	statsB := statMap(b)

	// Orden canónico de stats
	order := []string{"hp", "attack", "defense", "special-attack", "special-defense", "speed"}

	comparisons := make([]StatComparison, 0, len(order))
	totalA, totalB := 0, 0

	for _, name := range order {
		valA := statsA[name]
		valB := statsB[name]
		totalA += valA
		totalB += valB

		diff := valA - valB
		winner := "tie"
		if diff > 0 {
			winner = "a"
		} else if diff < 0 {
			winner = "b"
		}

		comparisons = append(comparisons, StatComparison{
			Name:   name,
			StatA:  valA,
			StatB:  valB,
			Diff:   diff,
			Winner: winner,
		})
	}

	overallWinner := "tie"
	if totalA > totalB {
		overallWinner = "a"
	} else if totalB > totalA {
		overallWinner = "b"
	}

	return PokemonComparison{
		PokemonA: a,
		PokemonB: b,
		Stats:    comparisons,
		TotalA:   totalA,
		TotalB:   totalB,
		Winner:   overallWinner,
	}
}
