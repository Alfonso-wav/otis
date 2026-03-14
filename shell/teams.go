package shell

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alfon/pokemon-app/core"
)

// FileTeamStorage persists teams as individual JSON files on disk.
type FileTeamStorage struct {
	dir string
}

// NewFileTeamStorage creates a new FileTeamStorage with the given base directory.
func NewFileTeamStorage(dir string) *FileTeamStorage {
	return &FileTeamStorage{dir: dir}
}

var slugRegexp = regexp.MustCompile(`[^a-z0-9-]+`)

// slugify normalizes a team name into a safe filename slug.
func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	s = slugRegexp.ReplaceAllString(s, "")
	if s == "" {
		s = "team"
	}
	return s
}

func (fs *FileTeamStorage) filePath(name string) string {
	return filepath.Join(fs.dir, slugify(name)+".json")
}

func (fs *FileTeamStorage) ensureDir() error {
	return os.MkdirAll(fs.dir, 0755)
}

// SaveTeam writes a team to disk as a JSON file.
func (fs *FileTeamStorage) SaveTeam(team core.Team) error {
	if err := core.ValidateTeam(team); err != nil {
		return err
	}
	if err := fs.ensureDir(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(team, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fs.filePath(team.Name), data, 0644)
}

// ListTeams reads all team JSON files from the directory.
func (fs *FileTeamStorage) ListTeams() ([]core.Team, error) {
	if err := fs.ensureDir(); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(fs.dir)
	if err != nil {
		return nil, err
	}
	var teams []core.Team
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(fs.dir, entry.Name()))
		if err != nil {
			continue
		}
		var team core.Team
		if err := json.Unmarshal(data, &team); err != nil {
			continue
		}
		teams = append(teams, team)
	}
	return teams, nil
}

// GetTeam reads a single team by name from disk.
func (fs *FileTeamStorage) GetTeam(name string) (core.Team, error) {
	data, err := os.ReadFile(fs.filePath(name))
	if err != nil {
		return core.Team{}, err
	}
	var team core.Team
	if err := json.Unmarshal(data, &team); err != nil {
		return core.Team{}, err
	}
	return team, nil
}

// DeleteTeam removes a team file from disk.
func (fs *FileTeamStorage) DeleteTeam(name string) error {
	return os.Remove(fs.filePath(name))
}
