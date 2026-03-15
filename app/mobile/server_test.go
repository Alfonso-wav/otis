package mobile

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alfon/pokemon-app/app"
	"github.com/alfon/pokemon-app/core"
)

// mockFetcher implements core.PokemonFetcher for tests.
type mockFetcher struct {
	pokemon    core.Pokemon
	pokemonErr error
	list       core.PokemonListResponse
	listErr    error
}

func (m *mockFetcher) FetchPokemon(name string) (core.Pokemon, error) {
	return m.pokemon, m.pokemonErr
}
func (m *mockFetcher) FetchPokemonList(offset int, limit int) (core.PokemonListResponse, error) {
	return m.list, m.listErr
}
func (m *mockFetcher) FetchTypeList() (core.TypeListResponse, error)                { return core.TypeListResponse{}, nil }
func (m *mockFetcher) FetchType(name string) (core.PokemonTypeDetail, error)         { return core.PokemonTypeDetail{}, nil }
func (m *mockFetcher) FetchRegions() ([]core.Region, error)                          { return nil, nil }
func (m *mockFetcher) FetchRegion(name string) (core.Region, error)                  { return core.Region{}, nil }
func (m *mockFetcher) FetchMove(name string) (core.Move, error)                      { return core.Move{}, nil }
func (m *mockFetcher) FetchAbility(name string) (core.Ability, error)                { return core.Ability{}, nil }
func (m *mockFetcher) FetchEvolutionChain(id int) (core.EvolutionChain, error)       { return core.EvolutionChain{}, nil }
func (m *mockFetcher) FetchPokemonSpecies(name string) (core.PokemonSpecies, error)  { return core.PokemonSpecies{}, nil }
func (m *mockFetcher) FetchPokemonForm(name string) (core.PokemonForm, error)        { return core.PokemonForm{}, nil }
func (m *mockFetcher) FetchPokemonColors() ([]core.NamedResource, error)             { return nil, nil }
func (m *mockFetcher) FetchPokemonShapes() ([]core.NamedResource, error)             { return nil, nil }
func (m *mockFetcher) FetchPokemonHabitats() ([]core.NamedResource, error)           { return nil, nil }
func (m *mockFetcher) FetchNatureList() ([]core.NamedResource, error)                { return nil, nil }
func (m *mockFetcher) FetchNatureDetail(name string) (core.NatureDetail, error)      { return core.NatureDetail{}, nil }
func (m *mockFetcher) FetchEggGroup(name string) (core.EggGroup, error)              { return core.EggGroup{}, nil }
func (m *mockFetcher) FetchGender(name string) (core.Gender, error)                  { return core.Gender{}, nil }
func (m *mockFetcher) FetchGrowthRate(name string) (core.GrowthRate, error)          { return core.GrowthRate{}, nil }
func (m *mockFetcher) FetchAllMoves() ([]core.Move, error)                           { return nil, nil }
func (m *mockFetcher) FetchMoveList(offset int, limit int) (core.MoveListResponse, error) { return core.MoveListResponse{}, nil }
func (m *mockFetcher) FetchMoveDamageClass(name string) (core.MoveDamageClass, error) { return core.MoveDamageClass{}, nil }
func (m *mockFetcher) FetchMoveAilment(name string) (core.MoveAilment, error)        { return core.MoveAilment{}, nil }
func (m *mockFetcher) FetchMoveTarget(name string) (core.MoveTarget, error)          { return core.MoveTarget{}, nil }
func (m *mockFetcher) FetchMachine(id int) (core.Machine, error)                     { return core.Machine{}, nil }
func (m *mockFetcher) FetchAllAbilities() ([]core.Ability, error)                    { return nil, nil }
func (m *mockFetcher) FetchAbilityList(offset int, limit int) (core.AbilityListResponse, error) { return core.AbilityListResponse{}, nil }
func (m *mockFetcher) FetchLocation(name string) (core.LocationDetail, error)        { return core.LocationDetail{}, nil }
func (m *mockFetcher) FetchLocationArea(name string) (core.LocationArea, error)      { return core.LocationArea{}, nil }
func (m *mockFetcher) FetchStat(name string) (core.StatDetail, error)                { return core.StatDetail{}, nil }
func (m *mockFetcher) FetchGenerations() ([]core.NamedResource, error)               { return nil, nil }
func (m *mockFetcher) FetchGeneration(name string) (core.Generation, error)          { return core.Generation{}, nil }
func (m *mockFetcher) FetchPokedexList() ([]core.NamedResource, error)               { return nil, nil }
func (m *mockFetcher) FetchPokedex(name string) (core.Pokedex, error)                { return core.Pokedex{}, nil }
func (m *mockFetcher) FetchVersionGroups() ([]core.NamedResource, error)             { return nil, nil }
func (m *mockFetcher) FetchVersionGroup(name string) (core.VersionGroup, error)      { return core.VersionGroup{}, nil }

type mockScraper struct{}

func (m *mockScraper) FetchPokedex() ([]core.PokedexDBEntry, error) { return nil, nil }

type mockSpriteDownloader struct{}

func (m *mockSpriteDownloader) DownloadAllSprites(destDir string, categories []core.SpriteCategory) (core.SpriteDownloadResult, error) {
	return core.SpriteDownloadResult{}, nil
}

type mockTeamStorage struct {
	teams map[string]core.Team
}

func newMockTeamStorage() *mockTeamStorage {
	return &mockTeamStorage{teams: make(map[string]core.Team)}
}

func (m *mockTeamStorage) SaveTeam(team core.Team) error   { m.teams[team.Name] = team; return nil }
func (m *mockTeamStorage) ListTeams() ([]core.Team, error) {
	var result []core.Team
	for _, t := range m.teams {
		result = append(result, t)
	}
	return result, nil
}
func (m *mockTeamStorage) GetTeam(name string) (core.Team, error) {
	t, ok := m.teams[name]
	if !ok {
		return core.Team{}, errors.New("not found")
	}
	return t, nil
}
func (m *mockTeamStorage) DeleteTeam(name string) error { delete(m.teams, name); return nil }

func newTestMux(fetcher *mockFetcher) *http.ServeMux {
	a := app.NewApp(fetcher, &mockScraper{}, newMockTeamStorage(), &mockSpriteDownloader{})
	mux := http.NewServeMux()
	RegisterRoutes(mux, a)
	return mux
}

func TestGetPokemonEndpoint(t *testing.T) {
	mux := newTestMux(&mockFetcher{
		pokemon: core.Pokemon{ID: 25, Name: "pikachu"},
	})

	req := httptest.NewRequest("GET", "/api/pokemon/pikachu", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var p core.Pokemon
	json.NewDecoder(w.Body).Decode(&p)
	if p.Name != "pikachu" || p.ID != 25 {
		t.Errorf("got %+v, want pikachu/25", p)
	}
}

func TestGetPokemonNotFound(t *testing.T) {
	mux := newTestMux(&mockFetcher{
		pokemonErr: errors.New("not found"),
	})

	req := httptest.NewRequest("GET", "/api/pokemon/unknown", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestListPokemonEndpoint(t *testing.T) {
	mux := newTestMux(&mockFetcher{
		list: core.PokemonListResponse{
			Count:   2,
			Results: []core.PokemonListItem{{Name: "bulbasaur"}, {Name: "ivysaur"}},
		},
	})

	req := httptest.NewRequest("GET", "/api/pokemon?offset=0&limit=2", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp core.PokemonListResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.Count != 2 || len(resp.Results) != 2 {
		t.Errorf("got %+v", resp)
	}
}

func TestCORSHeaders(t *testing.T) {
	a := app.NewApp(&mockFetcher{}, &mockScraper{}, newMockTeamStorage(), &mockSpriteDownloader{})
	srv := NewServer(a, 0)

	req := httptest.NewRequest("OPTIONS", "/api/pokemon", nil)
	w := httptest.NewRecorder()
	srv.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for OPTIONS, got %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected CORS origin *, got %q", got)
	}
}

func TestCreateTeamEndpoint(t *testing.T) {
	ts := newMockTeamStorage()
	a := app.NewApp(&mockFetcher{}, &mockScraper{}, ts, &mockSpriteDownloader{})
	mux := http.NewServeMux()
	RegisterRoutes(mux, a)

	body := `{"name":"MyTeam"}`
	req := httptest.NewRequest("POST", "/api/teams", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	if _, err := ts.GetTeam("MyTeam"); err != nil {
		t.Errorf("team not saved: %v", err)
	}
}

func TestNaturesEndpoint(t *testing.T) {
	mux := newTestMux(&mockFetcher{})

	req := httptest.NewRequest("GET", "/api/natures", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var natures []core.Nature
	json.NewDecoder(w.Body).Decode(&natures)
	if len(natures) == 0 {
		t.Error("expected natures list, got empty")
	}
}
