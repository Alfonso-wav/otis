package core

import (
	"strings"
)

// NormalizeName normaliza el nombre de un Pokemon a minusculas y sin espacios extra.
func NormalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// FilterByClassification filtra una lista de Pokémon según su clasificación legendary/mythical.
func FilterByClassification(items []PokemonListItem, classifications map[string]SpeciesClassification, legendary, mythical bool) []PokemonListItem {
	result := make([]PokemonListItem, 0)
	for _, item := range items {
		c, ok := classifications[NormalizeName(item.Name)]
		if !ok {
			continue
		}
		if legendary && mythical {
			if c.IsLegendary || c.IsMythical {
				result = append(result, item)
			}
		} else if legendary {
			if c.IsLegendary {
				result = append(result, item)
			}
		} else if mythical {
			if c.IsMythical {
				result = append(result, item)
			}
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
