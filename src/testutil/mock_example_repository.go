package testutil

import (
	"context"

	"github.com/jhseoeo/fiber-skeleton/src/model"
)

// MockExampleRepository is a test double for ExampleRepositoryPort.
// Assign the function fields to control behaviour per test case.
type MockExampleRepository struct {
	FindByIDFn func(ctx context.Context, id uint) (*model.Example, error)
	ListFn     func(ctx context.Context, offset, limit int) ([]*model.Example, int, error)
	CreateFn   func(ctx context.Context, example *model.Example) error
	UpdateFn   func(ctx context.Context, example *model.Example) error
	DeleteFn   func(ctx context.Context, id uint) error
}

func (m *MockExampleRepository) FindByID(ctx context.Context, id uint) (*model.Example, error) {
	return m.FindByIDFn(ctx, id)
}

func (m *MockExampleRepository) List(ctx context.Context, offset, limit int) ([]*model.Example, int, error) {
	return m.ListFn(ctx, offset, limit)
}

func (m *MockExampleRepository) Create(ctx context.Context, example *model.Example) error {
	return m.CreateFn(ctx, example)
}

func (m *MockExampleRepository) Update(ctx context.Context, example *model.Example) error {
	return m.UpdateFn(ctx, example)
}

func (m *MockExampleRepository) Delete(ctx context.Context, id uint) error {
	return m.DeleteFn(ctx, id)
}
