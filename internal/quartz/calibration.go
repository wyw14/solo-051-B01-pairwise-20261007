package quartz

import "math"

type Calibration struct {
	Slope  float64
	Offset float64
}

func DefaultCalibration() Calibration { return Calibration{Slope: 0.42, Offset: 1.5} }
func (c Calibration) Thickness(freq float64) float64 {
	if freq <= 0 {
		return 0
	}
	return math.Max(0, c.Slope*freq+c.Offset)
}
