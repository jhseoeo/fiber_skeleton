package servicemock

import (
	"context"

	"github.com/jhseoeo/fiber-skeleton/src/model"
)

// MockExampleService is a test double for ExampleServicePort.
// Assign the function fields to control behaviour per test case.
// Calling a method whose function field is nil panics with a descriptive message.
type MockExampleService struct {
	GetExampleFn    func(ctx context.Context, id uint) (*model.Example, error)
	ListExamplesFn  func(ctx context.Context, page, limit int) ([]*model.Example, int, error)
	CreateExampleFn func(ctx context.Context, example *model.Example) error
	UpdateExampleFn func(ctx context.Context, example *model.Example) error
	DeleteExampleFn func(ctx context.Context, id uint) error
}

func (m *MockExampleService) GetExample(ctx context.Context, id uint) (*model.Example, error) {
	if m.GetExampleFn == nil {
		panic("MockExampleService.GetExampleFn not set")
	}
	return m.GetExampleFn(ctx, id)
}

func (m *MockExampleService) ListExamples(ctx context.Context, page, limit int) ([]*model.Example, int, error) {
	if m.ListExamplesFn == nil {
		panic("MockExampleService.ListExamplesFn not set")
	}
	return m.ListExamplesFn(ctx, page, limit)
}

func (m *MockExampleService) CreateExample(ctx context.Context, example *model.Example) error {
	if m.CreateExampleFn == nil {
		panic("MockExampleService.CreateExampleFn not set")
	}
	return m.CreateExampleFn(ctx, example)
}

func (m *MockExampleService) UpdateExample(ctx context.Context, example *model.Example) error {
	if m.UpdateExampleFn == nil {
		panic("MockExampleService.UpdateExampleFn not set")
	}
	return m.UpdateExampleFn(ctx, example)
}

func (m *MockExampleService) DeleteExample(ctx context.Context, id uint) error {
	if m.DeleteExampleFn == nil {
		panic("MockExampleService.DeleteExampleFn not set")
	}
	return m.DeleteExampleFn(ctx, id)
}
