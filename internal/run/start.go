package run

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-161/internal/interlock"
	"github.com/wyw14/cry-161/internal/model"
	"github.com/wyw14/cry-161/internal/vacuum"
)

type Starter struct {
	Registry *Registry
	Vacuum   *vacuum.Controller
}

func (s *Starter) Permit(ctx context.Context, runID model.ID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	value, ok := s.Registry.Get(runID)
	if !ok {
		return fmt.Errorf("run not found")
	}
	if err := interlock.Permit(s.Vacuum, *value); err != nil {
		return fmt.Errorf("start permit: %w", err)
	}
	return nil
}

func (s *Starter) StartCoating(ctx context.Context, runID model.ID) error {
	if err := s.Permit(ctx, runID); err != nil {
		return err
	}
	if !s.Registry.Transition(runID, model.Coating) {
		return fmt.Errorf("cannot enter coating")
	}
	return nil
}
