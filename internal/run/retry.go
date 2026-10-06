package run

import (
	"errors"
	"github.com/wyw14/cry-161/internal/source"
)

func ShouldRetry(err error) bool { return errors.Is(err, source.ErrSourceBusy) }
func RetryBudget(err error, attempts int) int {
	if ShouldRetry(err) && attempts < 3 {
		return 3 - attempts
	}
	return 0
}
