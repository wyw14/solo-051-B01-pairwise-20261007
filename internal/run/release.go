package run

import (
	"fmt"
	"github.com/wyw14/cry-161/internal/optical"
)

func ReleaseQuality(record *optical.Monitor) error {
	sample := record.Latest()
	if !sample.Valid {
		return fmt.Errorf("invalid thickness evidence")
	}
	return nil
}
