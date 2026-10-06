package chamber

import (
	"context"
	"github.com/wyw14/cry-161/internal/vacuum"
)

func Vent(c *Chamber, pump *vacuum.Controller, ctx context.Context) error {
	if err := pump.Vent(ctx); err != nil {
		return err
	}
	c.SetVacuum(1, false)
	c.SetInterlock("vacuum", false)
	return nil
}
func Cool(c *Chamber) { c.SetInterlock("cooling", true); c.SetSourcePower(0) }
