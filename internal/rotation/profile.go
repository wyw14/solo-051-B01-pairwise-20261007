package rotation

import "github.com/wyw14/cry-161/internal/model"

type Profile struct {
	RPM             float64
	Acceleration    float64
	DurationSeconds int
}

func ProfileFor(layer model.Layer) Profile {
	rpm := 12.0
	if layer.RateNanometers > 1 {
		rpm = 18
	}
	return Profile{RPM: rpm, Acceleration: 4, DurationSeconds: int(layer.TargetNanometers/10) + 1}
}
func (p Profile) Valid() bool { return p.RPM > 0 && p.DurationSeconds > 0 }
