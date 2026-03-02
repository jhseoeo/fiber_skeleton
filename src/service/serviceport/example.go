package serviceport

import (
	"context"

	"github.com/jhseoeo/fiber-skeleton/src/model"
)

type ExampleServicePort interface {
	GetExample(ctx context.Context, id uint) (*model.Example, error)
	CreateExample(ctx context.Context, example *model.Example) error
	UpdateExample(ctx context.Context, example *model.Example) error
	DeleteExample(ctx context.Context, id uint) error
}
