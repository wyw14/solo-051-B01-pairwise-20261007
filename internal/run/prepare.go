package run

import (
	"context"
	"github.com/wyw14/cry-161/internal/model"
	"github.com/wyw14/cry-161/internal/source"
)

func PrepareSource(ctx context.Context, controller *source.Controller, layer model.Layer) error {
	return controller.Layer(ctx, layer)
}
