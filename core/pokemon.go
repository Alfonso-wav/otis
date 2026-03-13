package core

import "strings"

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
