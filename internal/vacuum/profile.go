package vacuum

import (
	"fmt"
	"time"
)

type PumpStage struct {
	Name           string        `json:"name"`
	TargetPressure float64       `json:"target_pressure"`
	Timeout        time.Duration `json:"timeout"`
}

type PumpProfile struct {
	Name   string      `json:"name"`
	Stages []PumpStage `json:"stages"`
}

func ProductionProfile() PumpProfile {
	return PumpProfile{Name: "optic-production", Stages: []PumpStage{{Name: "roughing", TargetPressure: 1e-1, Timeout: 90 * time.Second}, {Name: "high-vacuum", TargetPressure: 1e-4, Timeout: 180 * time.Second}, {Name: "stabilize", TargetPressure: 1e-5, Timeout: 60 * time.Second}}}
}

func (p PumpProfile) Validate() error {
	if p.Name == "" || len(p.Stages) == 0 {
		return fmt.Errorf("pump profile is empty")
	}
	previous := 2.0
	for _, stage := range p.Stages {
		if stage.Name == "" || stage.TargetPressure <= 0 || stage.Timeout <= 0 {
			return fmt.Errorf("pump stage is invalid")
		}
		if stage.TargetPressure >= previous {
			return fmt.Errorf("pump stages must lower pressure")
		}
		previous = stage.TargetPressure
	}
	return nil
}

func (p PumpProfile) TotalTimeout() time.Duration {
	total := time.Duration(0)
	for _, stage := range p.Stages {
		total += stage.Timeout
	}
	return total
}
