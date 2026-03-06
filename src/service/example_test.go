package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jhseoeo/fiber-skeleton/src/model"
	repositoryerror "github.com/jhseoeo/fiber-skeleton/src/repository/error"
	repositorymock "github.com/jhseoeo/fiber-skeleton/src/repository/mock"
	"github.com/jhseoeo/fiber-skeleton/src/service"
)

func TestGetExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &repositorymock.MockExampleRepository{
			FindByIDFn: func(_ context.Context, id uint) (*model.Example, error) {
				return &model.Example{ID: id, Content: "hello"}, nil
			},
		}
		svc := service.NewExampleService(mock)

		example, err := svc.GetExample(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), example.ID)
		assert.Equal(t, "hello", example.Content)
	})

	t.Run("not found", func(t *testing.T) {
		mock := &repositorymock.MockExampleRepository{
			FindByIDFn: func(_ context.Context, id uint) (*model.Example, error) {
				return nil, repositoryerror.ErrNotFound.New("example 1")
			},
		}
		svc := service.NewExampleService(mock)

		_, err := svc.GetExample(context.Background(), 1)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repositoryerror.ErrNotFound))
	})
}

func TestCreateExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &repositorymock.MockExampleRepository{
			CreateFn: func(_ context.Context, example *model.Example) error {
				example.ID = 1
				return nil
			},
		}
		svc := service.NewExampleService(mock)

		example := &model.Example{Content: "hello"}
		err := svc.CreateExample(context.Background(), example)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), example.ID)
	})
}

func TestUpdateExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &repositorymock.MockExampleRepository{
			UpdateFn: func(_ context.Context, example *model.Example) error {
				return nil
			},
		}
		svc := service.NewExampleService(mock)

		err := svc.UpdateExample(context.Background(), &model.Example{ID: 1, Content: "updated"})

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock := &repositorymock.MockExampleRepository{
			UpdateFn: func(_ context.Context, example *model.Example) error {
				return repositoryerror.ErrNotFound.New("example 1")
			},
		}
		svc := service.NewExampleService(mock)

		err := svc.UpdateExample(context.Background(), &model.Example{ID: 1, Content: "updated"})

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repositoryerror.ErrNotFound))
	})
}

func TestDeleteExample(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &repositorymock.MockExampleRepository{
			DeleteFn: func(_ context.Context, id uint) error {
				return nil
			},
		}
		svc := service.NewExampleService(mock)

		err := svc.DeleteExample(context.Background(), 1)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock := &repositorymock.MockExampleRepository{
			DeleteFn: func(_ context.Context, id uint) error {
				return repositoryerror.ErrNotFound.New("example 1")
			},
		}
		svc := service.NewExampleService(mock)

		err := svc.DeleteExample(context.Background(), 1)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repositoryerror.ErrNotFound))
	})
}
