package substrate

import "github.com/wyw14/cry-161/internal/model"

type Inspection struct {
	Surface string
	Count   int
	Passed  bool
}

func Inspect(run model.Run) Inspection {
	return Inspection{Surface: run.Batch, Count: 1, Passed: run.State != model.Failed}
}
func (i Inspection) Releaseable() bool { return i.Passed && i.Count > 0 }
