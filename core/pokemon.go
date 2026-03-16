package core

import (
	"strings"
)

// NormalizeName normaliza el nombre de un Pokemon a minusculas y sin espacios extra.
func NormalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
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
