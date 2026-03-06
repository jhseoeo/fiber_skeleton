// actual repository implementation is not implemented yet
// this is just a placeholder

package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/jhseoeo/fiber-skeleton/src/model"
	repositoryerror "github.com/jhseoeo/fiber-skeleton/src/repository/error"
	repositoryport "github.com/jhseoeo/fiber-skeleton/src/repository/port"
)

var _ repositoryport.ExampleRepositoryPort = (*ExampleRepository)(nil)

type ExampleRepository struct {
	mu       sync.RWMutex
	examples map[uint]*model.Example
	nextID   uint
}

func NewExampleRepository() *ExampleRepository {
	return &ExampleRepository{
		examples: make(map[uint]*model.Example),
		nextID:   1,
	}
}

func (r *ExampleRepository) List(_ context.Context, offset, limit int) ([]*model.Example, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Collect IDs and sort for stable pagination order.
	ids := make([]uint, 0, len(r.examples))
	for id := range r.examples {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	total := len(ids)
	if offset >= total {
		return []*model.Example{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}

	// Return value copies to avoid shared-pointer races after the lock is released.
	page := ids[offset:end]
	result := make([]*model.Example, len(page))
	for i, id := range page {
		cp := *r.examples[id]
		result[i] = &cp
	}
	return result, total, nil
}

func (r *ExampleRepository) FindByID(_ context.Context, id uint) (*model.Example, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if example, ok := r.examples[id]; ok {
		cp := *example
		return &cp, nil
	}
	return nil, repositoryerror.ErrNotFound.New(fmt.Sprintf("example %d", id))
}

func (r *ExampleRepository) Create(_ context.Context, example *model.Example) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	example.ID = r.nextID
	r.nextID++
	cp := *example
	r.examples[example.ID] = &cp
	return nil
}

func (r *ExampleRepository) Update(_ context.Context, example *model.Example) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.examples[example.ID]; !ok {
		return repositoryerror.ErrNotFound.New(fmt.Sprintf("example %d", example.ID))
	}
	cp := *example
	r.examples[example.ID] = &cp
	return nil
}

func (r *ExampleRepository) Delete(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.examples[id]; !ok {
		return repositoryerror.ErrNotFound.New(fmt.Sprintf("example %d", id))
	}
	delete(r.examples, id)
	return nil
}
