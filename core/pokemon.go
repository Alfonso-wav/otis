package core

import (
	"strings"
)

// NormalizeName normaliza el nombre de un Pokemon a minusculas y sin espacios extra.
func NormalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// FilterByType retorna los Pokemon cuyo slice de tipos contiene el tipo dado.
func FilterByType(pokemons []Pokemon, typeName string) []Pokemon {
	normalized := NormalizeName(typeName)
	result := make([]Pokemon, 0)
	for _, p := range pokemons {
		for _, t := range p.Types {
			if NormalizeName(t.Name) == normalized {
				result = append(result, p)
				break
			}
		}
	}
	return result
}

// GetStat retorna el BaseStat de una estadistica por nombre, o -1 si no existe.
func GetStat(p Pokemon, statName string) int {
	normalized := NormalizeName(statName)
	for _, s := range p.Stats {
		if NormalizeName(s.Name) == normalized {
			return s.BaseStat
		}
	}
	return -1
}

// FilterMovesByType retorna los movimientos que coinciden con el tipo dado.
func FilterMovesByType(moves []Move, typeName string) []Move {
	normalized := NormalizeName(typeName)
	result := make([]Move, 0)
	for _, m := range moves {
		if NormalizeName(m.Type) == normalized {
			result = append(result, m)
		}
	}
	return result
}

// FilterMovesByCategory retorna los movimientos que coinciden con la categoría dada.
// Categorías válidas: "physical", "special", "status".
func FilterMovesByCategory(moves []Move, category string) []Move {
	normalized := NormalizeName(category)
	result := make([]Move, 0)
	for _, m := range moves {
		if NormalizeName(m.Category) == normalized {
			result = append(result, m)
		}
	}
	return result
}

// SearchMoves retorna los movimientos cuyo nombre contiene la query (case-insensitive).
func SearchMoves(moves []Move, query string) []Move {
	q := NormalizeName(query)
	if q == "" {
		return moves
	}
	result := make([]Move, 0)
	for _, m := range moves {
		if strings.Contains(NormalizeName(m.Name), q) {
			result = append(result, m)
		}
	}
	return result
}

// FilterPokedexByType retorna los nombres que aparecen tanto en la lista del pokédex regional
// como en la lista de Pokémon de un tipo dado. Ambas listas se comparan normalizadas.
func FilterPokedexByType(pokedexNames []string, typePokemonNames []string) []string {
	typeSet := make(map[string]struct{}, len(typePokemonNames))
	for _, n := range typePokemonNames {
		typeSet[NormalizeName(n)] = struct{}{}
	}
	result := make([]string, 0)
	for _, n := range pokedexNames {
		if _, ok := typeSet[NormalizeName(n)]; ok {
			result = append(result, NormalizeName(n))
		}
	}
	return result
}

// SearchAbilities retorna las habilidades cuyo nombre contiene la query (case-insensitive).
func SearchAbilities(abilities []Ability, query string) []Ability {
	q := NormalizeName(query)
	if q == "" {
		return abilities
	}
	result := make([]Ability, 0)
	for _, a := range abilities {
		if strings.Contains(NormalizeName(a.Name), q) {
			result = append(result, a)
		}
	}
	return result
}
