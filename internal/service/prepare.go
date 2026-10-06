package service

import (
	"context"
	"fmt"
	"github.com/wyw14/cry-161/internal/model"
	"github.com/wyw14/cry-161/internal/run"
)

func (s *Service) Prepare(ctx context.Context, layer model.Layer) error {
	if err := run.PrepareSource(ctx, s.Source, layer); err != nil {
		if run.ShouldRetry(err) {
			return fmt.Errorf("source temporarily busy: %w", err)
		}
		return err
	}
	return nil
}
