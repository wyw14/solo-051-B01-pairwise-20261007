package interlock

import (
	"context"
	"github.com/wyw14/cry-161/internal/shutter"
)

func Overcoat(ctx context.Context, d *shutter.Driver, thickness float64, limit float64) error {
	if thickness <= limit {
		return nil
	}
	return shutter.CloseSafely(ctx, d)
}
