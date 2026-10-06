package journal

import (
	"encoding/json"
	"fmt"
	"github.com/wyw14/cry-161/internal/model"
	"os"
	"path/filepath"
	"sync"
)

type Snapshot struct {
	Runs      []*model.Run          `json:"runs"`
	Recipes   []*model.Recipe       `json:"recipes"`
	Chamber   model.ChamberSnapshot `json:"chamber"`
	Incidents []model.Incident      `json:"incidents"`
	Sequence  int64                 `json:"sequence"`
}

type SnapshotStore struct {
	mu  sync.Mutex
	dir string
}

func NewSnapshotStore(dir string) *SnapshotStore {
	if dir == "" {
		dir = "data"
	}
	_ = os.MkdirAll(dir, 0o755)
	return &SnapshotStore{dir: dir}
}

func (s *SnapshotStore) Save(snapshot Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if snapshot.Sequence < 0 {
		return fmt.Errorf("snapshot sequence is invalid")
	}
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	temporary := filepath.Join(s.dir, "opticcoat.snapshot.tmp")
	destination := filepath.Join(s.dir, "opticcoat.snapshot.json")
	if err := os.WriteFile(temporary, append(data, byte(10)), 0o644); err != nil {
		return err
	}
	return os.Rename(temporary, destination)
}

func (s *SnapshotStore) Load() (Snapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(filepath.Join(s.dir, "opticcoat.snapshot.json"))
	if os.IsNotExist(err) {
		return Snapshot{Runs: []*model.Run{}, Recipes: []*model.Recipe{}, Incidents: []model.Incident{}}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return Snapshot{}, err
	}
	if snapshot.Sequence < 0 {
		return Snapshot{}, fmt.Errorf("snapshot sequence is invalid")
	}
	return snapshot, nil
}
