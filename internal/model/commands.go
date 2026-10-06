package model

import (
	"fmt"
	"time"
)

type CommandKind string

const (
	CommandPump         CommandKind = "pump"
	CommandCondition    CommandKind = "condition"
	CommandOpenShutter  CommandKind = "open_shutter"
	CommandCloseShutter CommandKind = "close_shutter"
	CommandRotate       CommandKind = "rotate"
	CommandVent         CommandKind = "vent"
)

type Command struct {
	ID            ID          `json:"id"`
	RunID         ID          `json:"run_id"`
	LayerRevision int         `json:"layer_revision"`
	Sequence      int64       `json:"sequence"`
	Kind          CommandKind `json:"kind"`
	RequestedAt   time.Time   `json:"requested_at"`
	CompletedAt   time.Time   `json:"completed_at,omitempty"`
	Status        string      `json:"status"`
	Detail        string      `json:"detail,omitempty"`
}

func NewCommand(runID ID, revision int, sequence int64, kind CommandKind) Command {
	return Command{ID: NewID(), RunID: runID, LayerRevision: revision, Sequence: sequence, Kind: kind, RequestedAt: time.Now().UTC(), Status: "pending"}
}

func (c Command) Validate() error {
	if c.ID == "" || c.RunID == "" {
		return fmt.Errorf("command identity is required")
	}
	if c.LayerRevision < 0 || c.Sequence < 1 {
		return fmt.Errorf("command revision and sequence must be positive")
	}
	switch c.Kind {
	case CommandPump, CommandCondition, CommandOpenShutter, CommandCloseShutter, CommandRotate, CommandVent:
		return nil
	default:
		return fmt.Errorf("unsupported command kind %q", c.Kind)
	}
}

func (c Command) Complete(detail string) Command {
	c.Status = "completed"
	c.Detail = detail
	c.CompletedAt = time.Now().UTC()
	return c
}

func (c Command) Fail(err error) Command {
	c.Status = "failed"
	if err != nil {
		c.Detail = err.Error()
	}
	c.CompletedAt = time.Now().UTC()
	return c
}
