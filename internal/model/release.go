package model

import (
	"fmt"
	"sort"
	"time"
)

type ReleaseDecision string

const (
	ReleasePending  ReleaseDecision = "pending"
	ReleaseApproved ReleaseDecision = "approved"
	ReleaseRejected ReleaseDecision = "rejected"
)

type ReleaseCheck struct {
	Name   string    `json:"name"`
	Passed bool      `json:"passed"`
	Detail string    `json:"detail"`
	At     time.Time `json:"at"`
}

type ReleaseRecord struct {
	ID          ID              `json:"id"`
	RunID       ID              `json:"run_id"`
	Decision    ReleaseDecision `json:"decision"`
	Checks      []ReleaseCheck  `json:"checks"`
	RequestedAt time.Time       `json:"requested_at"`
	DecidedAt   time.Time       `json:"decided_at,omitempty"`
	Operator    string          `json:"operator"`
}

func NewReleaseRecord(runID ID, operator string) ReleaseRecord {
	return ReleaseRecord{ID: NewID(), RunID: runID, Decision: ReleasePending, Checks: make([]ReleaseCheck, 0), RequestedAt: time.Now().UTC(), Operator: operator}
}

func (r ReleaseRecord) AddCheck(name string, passed bool, detail string) ReleaseRecord {
	check := ReleaseCheck{Name: name, Passed: passed, Detail: detail, At: time.Now().UTC()}
	replaced := false
	for index, existing := range r.Checks {
		if existing.Name == name {
			r.Checks[index] = check
			replaced = true
			break
		}
	}
	if !replaced {
		r.Checks = append(r.Checks, check)
	}
	sort.Slice(r.Checks, func(i, j int) bool { return r.Checks[i].Name < r.Checks[j].Name })
	return r
}

func (r ReleaseRecord) Decide() ReleaseRecord {
	r.Decision = ReleaseApproved
	if len(r.Checks) == 0 {
		r.Decision = ReleaseRejected
	}
	for _, check := range r.Checks {
		if !check.Passed {
			r.Decision = ReleaseRejected
			break
		}
	}
	r.DecidedAt = time.Now().UTC()
	return r
}

func (r ReleaseRecord) Validate() error {
	if r.ID == "" || r.RunID == "" {
		return fmt.Errorf("release identity is required")
	}
	if r.Operator == "" {
		return fmt.Errorf("release operator is required")
	}
	if r.Decision == ReleasePending && !r.DecidedAt.IsZero() {
		return fmt.Errorf("pending release has decision time")
	}
	if r.Decision != ReleasePending && r.DecidedAt.IsZero() {
		return fmt.Errorf("decided release has no decision time")
	}
	names := map[string]bool{}
	for _, check := range r.Checks {
		if check.Name == "" {
			return fmt.Errorf("release check name is required")
		}
		if names[check.Name] {
			return fmt.Errorf("release check is duplicated")
		}
		names[check.Name] = true
	}
	return nil
}

func EvaluateRelease(run Run, evidence []LayerEvidence, incidents []Incident, operator string) ReleaseRecord {
	record := NewReleaseRecord(run.ID, operator)
	statePassed := run.State == Cooling || run.State == Venting || run.State == Released
	record = record.AddCheck("run_state", statePassed, string(run.State))
	evidencePassed := len(evidence) > 0
	for _, value := range evidence {
		if value.Validate() != nil {
			evidencePassed = false
			break
		}
	}
	record = record.AddCheck("layer_evidence", evidencePassed, fmt.Sprintf("layers=%d", len(evidence)))
	incidentsPassed := true
	for _, incident := range incidents {
		if incident.Severity == "critical" {
			incidentsPassed = false
			break
		}
	}
	record = record.AddCheck("incidents", incidentsPassed, fmt.Sprintf("incidents=%d", len(incidents)))
	return record.Decide()
}
