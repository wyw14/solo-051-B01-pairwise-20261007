package journal

import (
	"github.com/wyw14/cry-161/internal/model"
	"sort"
	"sync"
	"time"
)

type IncidentStore struct {
	mu     sync.RWMutex
	values map[model.ID]model.Incident
}

func NewIncidentStore() *IncidentStore {
	return &IncidentStore{values: map[model.ID]model.Incident{}}
}

func (s *IncidentStore) Raise(runID model.ID, severity string, message string) model.Incident {
	incident := model.Incident{ID: model.NewID(), RunID: runID, Severity: severity, Message: message, At: time.Now().UTC()}
	s.mu.Lock()
	s.values[incident.ID] = incident
	s.mu.Unlock()
	return incident
}

func (s *IncidentStore) List() []model.Incident {
	s.mu.RLock()
	result := make([]model.Incident, 0, len(s.values))
	for _, incident := range s.values {
		result = append(result, incident)
	}
	s.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].At.After(result[j].At) })
	return result
}

func (s *IncidentStore) ForRun(runID model.ID) []model.Incident {
	values := s.List()
	result := make([]model.Incident, 0)
	for _, incident := range values {
		if incident.RunID == runID {
			result = append(result, incident)
		}
	}
	return result
}
