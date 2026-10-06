package run

import (
	"fmt"
	"github.com/wyw14/cry-161/internal/model"
)

func Advance(registry *Registry, id model.ID, next model.RunState) error {
	if !registry.Transition(id, next) {
		return fmt.Errorf("invalid transition to %s", next)
	}
	return nil
}
func Complete(registry *Registry, id model.ID) error {
	if err := Advance(registry, id, model.Cooling); err != nil {
		return err
	}
	return Advance(registry, id, model.Venting)
}
