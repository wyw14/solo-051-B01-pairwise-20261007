package interlock

import (
	"errors"
	"github.com/wyw14/cry-161/internal/model"
	"github.com/wyw14/cry-161/internal/vacuum"
)

var ErrPermitDenied = errors.New("permit denied")

func Permit(c *vacuum.Controller, run model.Run) error {
	if err := vacuum.Check(c); err != nil {
		return err
	}
	if run.State != model.Conditioning && run.State != model.Coating {
		return ErrPermitDenied
	}
	return nil
}
