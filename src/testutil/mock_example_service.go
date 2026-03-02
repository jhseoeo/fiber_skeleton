package testutil

import (
	"context"

	"github.com/jhseoeo/fiber-skeleton/src/model"
)

// MockExampleService is a test double for ExampleServicePort.
// Assign the function fields to control behaviour per test case.
type MockExampleService struct {
	GetExampleFn    func(ctx context.Context, id uint) (*model.Example, error)
	CreateExampleFn func(ctx context.Context, example *model.Example) error
	UpdateExampleFn func(ctx context.Context, example *model.Example) error
	DeleteExampleFn func(ctx context.Context, id uint) error
}

func (m *MockExampleService) GetExample(ctx context.Context, id uint) (*model.Example, error) {
	return m.GetExampleFn(ctx, id)
}

func (m *MockExampleService) CreateExample(ctx context.Context, example *model.Example) error {
	return m.CreateExampleFn(ctx, example)
}

func (m *MockExampleService) UpdateExample(ctx context.Context, example *model.Example) error {
	return m.UpdateExampleFn(ctx, example)
}

func (m *MockExampleService) DeleteExample(ctx context.Context, id uint) error {
	return m.DeleteExampleFn(ctx, id)
}
