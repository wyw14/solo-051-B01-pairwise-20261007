package shutter

import "context"

func CloseSafely(ctx context.Context, d *Driver) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return d.Close()
}
