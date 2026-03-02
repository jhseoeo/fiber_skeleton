package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jhseoeo/fiber-skeleton/src/model"
	repositoryerror "github.com/jhseoeo/fiber-skeleton/src/repository/error"
	"github.com/jhseoeo/fiber-skeleton/src/service"
	"github.com/stretchr/testify/assert"
)

// mockExampleRepository is a manual mock implementing ExampleRepositoryPort.
type mockExampleRepository struct {
	findByID func(ctx context.Context, id uint) (*model.Example, error)
	create   func(ctx context.Context, example *model.Example) error
	update   func(ctx context.Context, example *model.Example) error
	delete   func(ctx context.Context, id uint) error
}

func (m *mockExampleRepository) FindByID(ctx context.Context, id uint) (*model.Example, error) {
	return m.findByID(ctx, id)
}
func (m *mockExampleRepository) Create(ctx context.Context, example *model.Example) error {
	return m.create(ctx, example)
}
func (m *mockExampleRepository) Update(ctx context.Context, example *model.Example) error {
	return m.update(ctx, example)
}
func (m *mockExampleRepository) Delete(ctx context.Context, id uint) error {
	return m.delete(ctx, id)
}

// --- GetExample ---

func TestGetExample_Success(t *testing.T) {
	mock := &mockExampleRepository{
		findByID: func(_ context.Context, id uint) (*model.Example, error) {
			return &model.Example{ID: id, Content: "hello"}, nil
		},
	}
	svc := service.NewExampleService(mock)

	example, err := svc.GetExample(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), example.ID)
	assert.Equal(t, "hello", example.Content)
}

func TestGetExample_NotFound(t *testing.T) {
	mock := &mockExampleRepository{
		findByID: func(_ context.Context, id uint) (*model.Example, error) {
			return nil, repositoryerror.ErrNotFound.New("example 1")
		},
	}
	svc := service.NewExampleService(mock)

	_, err := svc.GetExample(context.Background(), 1)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, repositoryerror.ErrNotFound))
}

// --- CreateExample ---

func TestCreateExample_Success(t *testing.T) {
	mock := &mockExampleRepository{
		create: func(_ context.Context, example *model.Example) error {
			example.ID = 1
			return nil
		},
	}
	svc := service.NewExampleService(mock)

	example := &model.Example{Content: "hello"}
	err := svc.CreateExample(context.Background(), example)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), example.ID)
}

// --- UpdateExample ---

func TestUpdateExample_Success(t *testing.T) {
	mock := &mockExampleRepository{
		update: func(_ context.Context, example *model.Example) error {
			return nil
		},
	}
	svc := service.NewExampleService(mock)

	example := &model.Example{ID: 1, Content: "updated"}
	err := svc.UpdateExample(context.Background(), example)

	assert.NoError(t, err)
}

func TestUpdateExample_NotFound(t *testing.T) {
	mock := &mockExampleRepository{
		update: func(_ context.Context, example *model.Example) error {
			return repositoryerror.ErrNotFound.New("example 1")
		},
	}
	svc := service.NewExampleService(mock)

	err := svc.UpdateExample(context.Background(), &model.Example{ID: 1, Content: "updated"})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, repositoryerror.ErrNotFound))
}

// --- DeleteExample ---

func TestDeleteExample_Success(t *testing.T) {
	mock := &mockExampleRepository{
		delete: func(_ context.Context, id uint) error {
			return nil
		},
	}
	svc := service.NewExampleService(mock)

	err := svc.DeleteExample(context.Background(), 1)

	assert.NoError(t, err)
}

func TestDeleteExample_NotFound(t *testing.T) {
	mock := &mockExampleRepository{
		delete: func(_ context.Context, id uint) error {
			return repositoryerror.ErrNotFound.New("example 1")
		},
	}
	svc := service.NewExampleService(mock)

	err := svc.DeleteExample(context.Background(), 1)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, repositoryerror.ErrNotFound))
}
