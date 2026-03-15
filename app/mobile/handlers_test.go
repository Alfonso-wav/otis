package mobile

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alfon/pokemon-app/app"
	"github.com/alfon/pokemon-app/core"
)

func TestHandlerInvalidJSON(t *testing.T) {
	mux := newTestMux(&mockFetcher{})

	req := httptest.NewRequest("POST", "/api/battle/simulate-damage", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] == "" {
		t.Error("expected error message in response")
	}
}

func TestHandlerInvalidPathID(t *testing.T) {
	mux := newTestMux(&mockFetcher{})

	req := httptest.NewRequest("GET", "/api/evolution-chain/abc", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteTeamEndpoint(t *testing.T) {
	ts := newMockTeamStorage()
	ts.teams["TestTeam"] = core.Team{Name: "TestTeam"}
	a := app.NewApp(&mockFetcher{}, &mockScraper{}, ts, &mockSpriteDownloader{})
	mux := http.NewServeMux()
	RegisterRoutes(mux, a)

	req := httptest.NewRequest("DELETE", "/api/teams/TestTeam", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if _, err := ts.GetTeam("TestTeam"); err == nil {
		t.Error("team should have been deleted")
	}
}

func TestBattleInitEndpoint(t *testing.T) {
	mux := newTestMux(&mockFetcher{})

	body := `{"attackerMaxHP":100,"defenderMaxHP":120}`
	req := httptest.NewRequest("POST", "/api/battle/init", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var state core.BattleState
	json.NewDecoder(w.Body).Decode(&state)
	if state.AttackerHP != 100 || state.DefenderHP != 120 {
		t.Errorf("unexpected state: %+v", state)
	}
}

func TestListPokemonErrorEndpoint(t *testing.T) {
	fetcher := &mockFetcher{listErr: errTest}
	mux := newTestMux(fetcher)

	req := httptest.NewRequest("GET", "/api/pokemon?offset=0&limit=20", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

var errTest = http.ErrAbortHandler // reuse a stdlib error for tests
