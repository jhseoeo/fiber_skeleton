package repositoryport

import (
	"context"

	"github.com/jhseoeo/fiber-skeleton/src/model"
)

type ExampleRepositoryPort interface {
	FindByID(ctx context.Context, id uint) (*model.Example, error)
	List(ctx context.Context, offset, limit int) ([]*model.Example, int, error)
	Create(ctx context.Context, example *model.Example) error
	Update(ctx context.Context, example *model.Example) error
	Delete(ctx context.Context, id uint) error
}
