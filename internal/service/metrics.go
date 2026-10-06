package service

import (
	"github.com/wyw14/cry-161/internal/model"
	"sort"
	"time"
)

type OperationsMetrics struct {
	CollectedAt     time.Time              `json:"collected_at"`
	RunsByState     map[model.RunState]int `json:"runs_by_state"`
	TotalRuns       int                    `json:"total_runs"`
	RecipeCount     int                    `json:"recipe_count"`
	IncidentCount   int                    `json:"incident_count"`
	LatestThickness float64                `json:"latest_thickness_nm"`
	ChamberPressure float64                `json:"chamber_pressure"`
	ActiveStates    []model.RunState       `json:"active_states"`
}

func BuildOperationsMetrics(svc *Service, runtime *Runtime) OperationsMetrics {
	values := svc.Runs.List()
	counts := map[model.RunState]int{}
	for _, value := range values {
		counts[value.State]++
	}
	states := make([]model.RunState, 0, len(counts))
	for state, count := range counts {
		if count > 0 {
			states = append(states, state)
		}
	}
	sort.Slice(states, func(i, j int) bool { return states[i] < states[j] })
	latest := runtime.Monitor.Latest()
	chamber := svc.Chamber.Snapshot()
	return OperationsMetrics{CollectedAt: time.Now().UTC(), RunsByState: counts, TotalRuns: len(values), RecipeCount: len(svc.Recipes.List()), IncidentCount: len(runtime.Incidents.List()), LatestThickness: latest.Nanometers, ChamberPressure: chamber.Pressure, ActiveStates: states}
}

func (m OperationsMetrics) Healthy() bool {
	if m.ChamberPressure < 0 {
		return false
	}
	if m.RunsByState[model.Failed] > 0 {
		return false
	}
	for _, state := range m.ActiveStates {
		switch state {
		case model.Loading, model.Pumping, model.Conditioning, model.Coating, model.Cooling, model.Venting, model.Released:
		default:
			return false
		}
	}
	return true
}
