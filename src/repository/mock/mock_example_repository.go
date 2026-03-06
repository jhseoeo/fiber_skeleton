package repositorymock

import (
	"context"

	"github.com/jhseoeo/fiber-skeleton/src/model"
)

// MockExampleRepository is a test double for ExampleRepositoryPort.
// Assign the function fields to control behaviour per test case.
// Calling a method whose function field is nil panics with a descriptive message.
type MockExampleRepository struct {
	FindByIDFn func(ctx context.Context, id uint) (*model.Example, error)
	ListFn     func(ctx context.Context, offset, limit int) ([]*model.Example, int, error)
	CreateFn   func(ctx context.Context, example *model.Example) error
	UpdateFn   func(ctx context.Context, example *model.Example) error
	DeleteFn   func(ctx context.Context, id uint) error
}

func (m *MockExampleRepository) FindByID(ctx context.Context, id uint) (*model.Example, error) {
	if m.FindByIDFn == nil {
		panic("MockExampleRepository.FindByIDFn not set")
	}
	return m.FindByIDFn(ctx, id)
}

func (m *MockExampleRepository) List(ctx context.Context, offset, limit int) ([]*model.Example, int, error) {
	if m.ListFn == nil {
		panic("MockExampleRepository.ListFn not set")
	}
	return m.ListFn(ctx, offset, limit)
}

func (m *MockExampleRepository) Create(ctx context.Context, example *model.Example) error {
	if m.CreateFn == nil {
		panic("MockExampleRepository.CreateFn not set")
	}
	return m.CreateFn(ctx, example)
}

func (m *MockExampleRepository) Update(ctx context.Context, example *model.Example) error {
	if m.UpdateFn == nil {
		panic("MockExampleRepository.UpdateFn not set")
	}
	return m.UpdateFn(ctx, example)
}

func (m *MockExampleRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFn == nil {
		panic("MockExampleRepository.DeleteFn not set")
	}
	return m.DeleteFn(ctx, id)
}
