package recipe

import (
	"github.com/wyw14/cry-161/internal/model"
	"sync"
)

type DraftBook struct {
	mu     sync.RWMutex
	drafts map[model.ID]*model.Recipe
}

func NewDraftBook() *DraftBook { return &DraftBook{drafts: map[model.ID]*model.Recipe{}} }
func (b *DraftBook) Save(recipe *model.Recipe) *model.Recipe {
	b.mu.Lock()
	defer b.mu.Unlock()
	copy := cloneRecipe(recipe)
	copy.Published = false
	b.drafts[copy.ID] = copy
	return cloneRecipe(copy)
}
func (b *DraftBook) Get(id model.ID) (*model.Recipe, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	value, ok := b.drafts[id]
	if !ok {
		return nil, false
	}
	return cloneRecipe(value), true
}
func (b *DraftBook) UpdateLayer(id model.ID, layer model.Layer) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	value, ok := b.drafts[id]
	if !ok {
		return false
	}
	value.Layers[layer.Number] = layer
	value.Revision++
	return true
}
