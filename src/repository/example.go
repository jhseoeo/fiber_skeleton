// actual repository implementation is not implemented yet
// this is just a placeholder

package repository

import (
	"context"
	"fmt"
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

func (r *ExampleRepository) FindByID(ctx context.Context, id uint) (*model.Example, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if example, ok := r.examples[id]; ok {
		return example, nil
	}
	return nil, repositoryerror.ErrNotFound.New(fmt.Sprintf("example %d", id))
}

func (r *ExampleRepository) Create(ctx context.Context, example *model.Example) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	example.ID = r.nextID
	r.nextID++
	r.examples[example.ID] = example
	return nil
}

func (r *ExampleRepository) Update(ctx context.Context, example *model.Example) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.examples[example.ID]; !ok {
		return repositoryerror.ErrNotFound.New(fmt.Sprintf("example %d", example.ID))
	}
	r.examples[example.ID] = example
	return nil
}

func (r *ExampleRepository) Delete(ctx context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.examples[id]; !ok {
		return repositoryerror.ErrNotFound.New(fmt.Sprintf("example %d", id))
	}
	delete(r.examples, id)
	return nil
}
