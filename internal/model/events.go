package model

import "time"

type Event struct {
	ID    ID             `json:"id"`
	Kind  string         `json:"kind"`
	RunID ID             `json:"run_id"`
	Data  map[string]any `json:"data"`
	At    time.Time      `json:"at"`
}

func NewEvent(kind string, runID ID, data map[string]any) Event {
	if data == nil {
		data = map[string]any{}
	}
	return Event{ID: NewID(), Kind: kind, RunID: runID, Data: data, At: time.Now().UTC()}
}
