package shutter

import (
	"errors"
	"sync"
)

var ErrShutterUnavailable = errors.New("shutter unavailable")

type Driver struct {
	mu        sync.Mutex
	open      bool
	available bool
}

func NewDriver() *Driver { return &Driver{available: true} }
func (d *Driver) Open() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.available {
		return ErrShutterUnavailable
	}
	d.open = true
	return nil
}
func (d *Driver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if !d.available {
		return ErrShutterUnavailable
	}
	d.open = false
	return nil
}
func (d *Driver) Opened() bool            { d.mu.Lock(); defer d.mu.Unlock(); return d.open }
func (d *Driver) SetAvailable(value bool) { d.mu.Lock(); d.available = value; d.mu.Unlock() }
