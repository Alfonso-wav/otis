// Package mobile exposes the App methods as REST JSON endpoints.
//
// Endpoint mapping (bindings.go method → REST endpoint):
//
// --- Pokémon base ---
// ListPokemon          → GET  /api/pokemon?offset=X&limit=Y
// GetPokemon           → GET  /api/pokemon/{name}
// ListTypes            → GET  /api/types
// GetType              → GET  /api/types/{name}
// ListRegions          → GET  /api/regions
// GetRegion            → GET  /api/regions/{name}
// GetRegionPokemonByType → GET /api/regions/{region}/pokemon-by-type/{type}
// GetMove              → GET  /api/moves/{name}
// GetAbility           → GET  /api/abilities/{name}
// GetEvolutionChain    → GET  /api/evolution-chain/{id}
// GetNatures           → GET  /api/natures
//
// --- Pokémon extendido (Grupo A) ---
// GetPokemonSpecies    → GET  /api/pokemon-species/{name}
// GetPokemonForm       → GET  /api/pokemon-forms/{name}
// ListPokemonColors    → GET  /api/pokemon-colors
// ListPokemonShapes    → GET  /api/pokemon-shapes
// ListPokemonHabitats  → GET  /api/pokemon-habitats
//
// --- Naturalezas y cría (Grupo B) ---
// ListNatureNames      → GET  /api/nature-names
// GetNatureDetail      → GET  /api/nature-detail/{name}
// GetEggGroup          → GET  /api/egg-groups/{name}
// GetGender            → GET  /api/genders/{name}
// GetGrowthRate        → GET  /api/growth-rates/{name}
//
// --- Movimientos (Grupo C) ---
// GetAllMoves          → GET  /api/moves/all
// ListMoves            → GET  /api/moves?offset=X&limit=Y
// GetMoveDamageClass   → GET  /api/move-damage-classes/{name}
// GetMoveAilment       → GET  /api/move-ailments/{name}
// GetMoveTarget        → GET  /api/move-targets/{name}
// GetMachine           → GET  /api/machines/{id}
//
// --- Habilidades (Grupo D) ---
// GetAllAbilities      → GET  /api/abilities/all
// ListAbilities        → GET  /api/abilities?offset=X&limit=Y
//
// --- Ubicaciones (Grupo F) ---
// GetLocation          → GET  /api/locations/{name}
// GetLocationArea      → GET  /api/location-areas/{name}
//
// --- Stats y generaciones (Grupo G) ---
// GetStatDetail        → GET  /api/stats/{name}
// ListGenerations      → GET  /api/generations
// GetGeneration        → GET  /api/generations/{name}
// ListPokedexes        → GET  /api/pokedexes
// GetPokedex           → GET  /api/pokedexes/{name}
// ListVersionGroups    → GET  /api/version-groups
// GetVersionGroup      → GET  /api/version-groups/{name}
//
// --- Scraper ---
// ScrapePokedex        → GET  /api/scrape/pokedex
//
// --- Sprites ---
// DownloadSprites      → POST /api/sprites/download
//
// --- Batalla ---
// SimulateDamage       → POST /api/battle/simulate-damage
// InitBattle           → POST /api/battle/init
// ExecuteTurn          → POST /api/battle/execute-turn
// SimulateFullBattle   → POST /api/battle/simulate-full
// SimulateMultipleBattles    → POST /api/battle/simulate-multiple
// SimulateTeamBattle         → POST /api/battle/team-simulate
// SimulateMultipleTeamBattles → POST /api/battle/team-simulate-multiple
//
// --- Calculadoras ---
// CalculateEVs         → POST /api/calculator/evs
// CalculateStats       → POST /api/calculator/stats
//
// --- Equipos ---
// ListTeams            → GET    /api/teams
// GetTeam              → GET    /api/teams/{name}
// CreateTeam           → POST   /api/teams
// DeleteTeam           → DELETE /api/teams/{name}
// SaveToTeam           → POST   /api/teams/{name}/members
// UpdateTeamMember     → PUT    /api/teams/{name}/members/{index}
// DeleteTeamMember     → DELETE /api/teams/{name}/members/{index}
// FillTeamRandom       → POST   /api/teams/{name}/fill-random
package mobile

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alfon/pokemon-app/app"
	"github.com/alfon/pokemon-app/core"
)

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func queryInt(r *http.Request, key string, defaultVal int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}

func pathInt(r *http.Request, key string) (int, bool) {
	v := r.PathValue(key)
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

func decodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// RegisterRoutes registers all REST endpoints on the given mux.
func RegisterRoutes(mux *http.ServeMux, a *app.App) {
	// --- Pokémon base ---
	mux.HandleFunc("GET /api/pokemon", func(w http.ResponseWriter, r *http.Request) {
		// If a "name" query param is present, treat as GetPokemon for compatibility
		offset := queryInt(r, "offset", 0)
		limit := queryInt(r, "limit", 20)
		result, err := a.ListPokemon(offset, limit)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/pokemon/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		result, err := a.GetPokemon(name)
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/types", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ListTypes()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/types/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetType(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/regions", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ListRegions()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/regions/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetRegion(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/regions/{region}/pokemon-by-type/{type}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetRegionPokemonByType(r.PathValue("region"), r.PathValue("type"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/moves/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetMove(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/abilities/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetAbility(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/evolution-chain/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := pathInt(r, "id")
		if !ok {
			jsonError(w, "invalid id", http.StatusBadRequest)
			return
		}
		result, err := a.GetEvolutionChain(id)
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/natures", func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, a.GetNatures())
	})

	// --- Grupo A: Pokémon extendido ---

	mux.HandleFunc("GET /api/pokemon-species/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetPokemonSpecies(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/pokemon-forms/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetPokemonForm(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/pokemon-colors", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ListPokemonColors()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/pokemon-shapes", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ListPokemonShapes()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/pokemon-habitats", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ListPokemonHabitats()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	// --- Grupo B: Naturalezas y cría ---

	mux.HandleFunc("GET /api/nature-names", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ListNatureNames()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/nature-detail/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetNatureDetail(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/egg-groups/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetEggGroup(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/genders/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetGender(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/growth-rates/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetGrowthRate(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	// --- Grupo C: Movimientos ---

	mux.HandleFunc("GET /api/moves/all", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetAllMoves()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/moves", func(w http.ResponseWriter, r *http.Request) {
		offset := queryInt(r, "offset", 0)
		limit := queryInt(r, "limit", 20)
		result, err := a.ListMoves(offset, limit)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/move-damage-classes/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetMoveDamageClass(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/move-ailments/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetMoveAilment(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/move-targets/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetMoveTarget(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/machines/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, ok := pathInt(r, "id")
		if !ok {
			jsonError(w, "invalid id", http.StatusBadRequest)
			return
		}
		result, err := a.GetMachine(id)
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	// --- Grupo D: Habilidades ---

	mux.HandleFunc("GET /api/abilities/all", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetAllAbilities()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/abilities", func(w http.ResponseWriter, r *http.Request) {
		offset := queryInt(r, "offset", 0)
		limit := queryInt(r, "limit", 20)
		result, err := a.ListAbilities(offset, limit)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	// --- Grupo F: Ubicaciones ---

	mux.HandleFunc("GET /api/locations/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetLocation(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/location-areas/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetLocationArea(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	// --- Grupo G: Stats y generaciones ---

	mux.HandleFunc("GET /api/stats/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetStatDetail(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/generations", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ListGenerations()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/generations/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetGeneration(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/pokedexes", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ListPokedexes()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/pokedexes/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetPokedex(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/version-groups", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ListVersionGroups()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/version-groups/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetVersionGroup(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	// --- Scraper ---

	mux.HandleFunc("GET /api/scrape/pokedex", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ScrapePokedex()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	// --- Sprites ---

	mux.HandleFunc("POST /api/sprites/download", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Categories []core.SpriteCategory `json:"categories"`
		}
		if err := decodeBody(r, &body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		result, err := a.DownloadSprites(body.Categories)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	// --- Batalla ---

	mux.HandleFunc("POST /api/battle/simulate-damage", func(w http.ResponseWriter, r *http.Request) {
		var input core.DamageInput
		if err := decodeBody(r, &input); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		result, err := a.SimulateDamage(input)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("POST /api/battle/init", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			AttackerMaxHP int `json:"attackerMaxHP"`
			DefenderMaxHP int `json:"defenderMaxHP"`
		}
		if err := decodeBody(r, &body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		jsonResponse(w, a.InitBattle(body.AttackerMaxHP, body.DefenderMaxHP))
	})

	mux.HandleFunc("POST /api/battle/execute-turn", func(w http.ResponseWriter, r *http.Request) {
		var input core.TurnInput
		if err := decodeBody(r, &input); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		jsonResponse(w, a.ExecuteTurn(input))
	})

	mux.HandleFunc("POST /api/battle/simulate-full", func(w http.ResponseWriter, r *http.Request) {
		var input core.FullBattleInput
		if err := decodeBody(r, &input); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		jsonResponse(w, a.SimulateFullBattle(input))
	})

	mux.HandleFunc("POST /api/battle/simulate-multiple", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			core.FullBattleInput
			N int `json:"n"`
		}
		if err := decodeBody(r, &body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		result, err := a.SimulateMultipleBattles(body.FullBattleInput, body.N)
		if err != nil {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("POST /api/battle/team-simulate", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Team1Name string `json:"team1Name"`
			Team2Name string `json:"team2Name"`
		}
		if err := decodeBody(r, &body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		result, err := a.SimulateTeamBattle(body.Team1Name, body.Team2Name)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("POST /api/battle/team-simulate-multiple", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Team1Name string `json:"team1Name"`
			Team2Name string `json:"team2Name"`
			N         int    `json:"n"`
		}
		if err := decodeBody(r, &body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		result, err := a.SimulateMultipleTeamBattles(body.Team1Name, body.Team2Name, body.N)
		if err != nil {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResponse(w, result)
	})

	// --- Calculadoras ---

	mux.HandleFunc("POST /api/calculator/evs", func(w http.ResponseWriter, r *http.Request) {
		var input core.EVCalculatorInput
		if err := decodeBody(r, &input); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		result, err := a.CalculateEVs(input)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("POST /api/calculator/stats", func(w http.ResponseWriter, r *http.Request) {
		var input core.StatCalculatorInput
		if err := decodeBody(r, &input); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		result, err := a.CalculateStats(input)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	// --- Equipos ---

	mux.HandleFunc("GET /api/teams", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.ListTeams()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("GET /api/teams/{name}", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.GetTeam(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonResponse(w, result)
	})

	mux.HandleFunc("POST /api/teams", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name string `json:"name"`
		}
		if err := decodeBody(r, &body); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := a.CreateTeam(body.Name); err != nil {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, map[string]string{"status": "created"})
	})

	mux.HandleFunc("DELETE /api/teams/{name}", func(w http.ResponseWriter, r *http.Request) {
		if err := a.DeleteTeam(r.PathValue("name")); err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, map[string]string{"status": "deleted"})
	})

	mux.HandleFunc("POST /api/teams/{name}/members", func(w http.ResponseWriter, r *http.Request) {
		var member core.TeamMember
		if err := decodeBody(r, &member); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := a.SaveToTeam(r.PathValue("name"), member); err != nil {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		jsonResponse(w, map[string]string{"status": "added"})
	})

	mux.HandleFunc("PUT /api/teams/{name}/members/{index}", func(w http.ResponseWriter, r *http.Request) {
		idx, ok := pathInt(r, "index")
		if !ok {
			jsonError(w, "invalid index", http.StatusBadRequest)
			return
		}
		var member core.TeamMember
		if err := decodeBody(r, &member); err != nil {
			jsonError(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := a.UpdateTeamMember(r.PathValue("name"), idx, member); err != nil {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResponse(w, map[string]string{"status": "updated"})
	})

	mux.HandleFunc("DELETE /api/teams/{name}/members/{index}", func(w http.ResponseWriter, r *http.Request) {
		idx, ok := pathInt(r, "index")
		if !ok {
			jsonError(w, "invalid index", http.StatusBadRequest)
			return
		}
		if err := a.DeleteTeamMember(r.PathValue("name"), idx); err != nil {
			jsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonResponse(w, map[string]string{"status": "deleted"})
	})

	mux.HandleFunc("POST /api/teams/{name}/fill-random", func(w http.ResponseWriter, r *http.Request) {
		result, err := a.FillTeamRandom(r.PathValue("name"))
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, result)
	})
}
