package optical

import (
	"fmt"
	"github.com/wyw14/cry-161/internal/journal"
)

func Analyze(record journal.SpectrumRecord) (float64, error) {
	if len(record.Wavelength) == 0 || len(record.Wavelength) != len(record.Intensity) {
		return 0, fmt.Errorf("spectrum dimensions mismatch")
	}
	total := 0.0
	for _, v := range record.Intensity {
		total += v
	}
	return total / float64(len(record.Intensity)), nil
}
