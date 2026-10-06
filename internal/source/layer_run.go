package source

import (
	"fmt"
	"github.com/wyw14/cry-161/internal/shutter"
)

func ExecuteLayer(run func() error, closeShutter func() error) (err error) {
	defer func() {
		if closeErr := closeShutter(); closeErr != nil {
			err = fmt.Errorf("shutter close: %w", closeErr)
		}
	}()
	err = run()
	return
}
func ExecuteWithDriver(run func() error, driver *shutter.Driver) error {
	return ExecuteLayer(run, driver.Close)
}
