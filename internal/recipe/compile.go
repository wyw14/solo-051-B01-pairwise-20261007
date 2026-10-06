package recipe

import (
	"fmt"
	"github.com/wyw14/cry-161/internal/model"
	"sort"
)

type Plan struct {
	RecipeID model.ID
	Revision int
	Layers   []model.Layer
}

func Compile(recipe *model.Recipe) (Plan, error) {
	if recipe == nil || len(recipe.Layers) == 0 {
		return Plan{}, fmt.Errorf("recipe has no layers")
	}
	keys := make([]int, 0, len(recipe.Layers))
	for k := range recipe.Layers {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	layers := make([]model.Layer, 0, len(keys))
	for _, k := range keys {
		layers = append(layers, recipe.Layers[k])
	}
	return Plan{RecipeID: recipe.ID, Revision: recipe.Revision, Layers: layers}, nil
}
