package chamber

import "github.com/wyw14/cry-161/internal/model"

func CopySnapshot(value model.ChamberSnapshot) model.ChamberSnapshot { return value }
func MergeSnapshots(primary, fallback model.ChamberSnapshot) model.ChamberSnapshot {
	result := primary
	if result.Name == "" {
		result.Name = fallback.Name
	}
	if result.UpdatedAt.IsZero() {
		result.UpdatedAt = fallback.UpdatedAt
	}
	return result
}
