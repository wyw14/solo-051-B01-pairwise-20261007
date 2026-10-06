package vacuum

import "fmt"

type VacuumFault struct {
	Pressure float64
	Reason   string
}

func (f *VacuumFault) Error() string {
	if f == nil {
		return "<nil>"
	}
	return fmt.Sprintf("vacuum pressure %.6g: %s", f.Pressure, f.Reason)
}

func Check(c *Controller) error {
	if c.Stable() && c.Pressure() < 0.001 {
		return nil
	}
	return &VacuumFault{Pressure: c.Pressure(), Reason: "pressure not stable"}
}
