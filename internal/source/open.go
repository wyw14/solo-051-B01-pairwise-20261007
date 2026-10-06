package source

import (
	"errors"
	"fmt"
)

var ErrSourceBusy = errors.New("source busy")

type Adapter struct{ busy bool }

func NewAdapter() *Adapter            { return &Adapter{} }
func (a *Adapter) SetBusy(value bool) { a.busy = value }
func (a *Adapter) Open() error {
	if a.busy {
		return fmt.Errorf("open source: %w", ErrSourceBusy)
	}
	return nil
}
