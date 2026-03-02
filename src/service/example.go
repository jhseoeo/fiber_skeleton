package service

import (
	"context"

	"github.com/jhseoeo/fiber-skeleton/src/model"
	repositoryport "github.com/jhseoeo/fiber-skeleton/src/repository/port"
	"github.com/jhseoeo/fiber-skeleton/src/service/serviceport"
)

var _ serviceport.ExampleServicePort = (*ExampleService)(nil)

type ExampleService struct {
	exampleRepository repositoryport.ExampleRepositoryPort
}

func NewExampleService(exampleRepository repositoryport.ExampleRepositoryPort) *ExampleService {
	return &ExampleService{
		exampleRepository: exampleRepository,
	}
}

func (s *ExampleService) GetExample(ctx context.Context, id uint) (*model.Example, error) {
	return s.exampleRepository.FindByID(ctx, id)
}

func (s *ExampleService) ListExamples(ctx context.Context, page, limit int) ([]*model.Example, int, error) {
	offset := (page - 1) * limit
	return s.exampleRepository.List(ctx, offset, limit)
}

func (s *ExampleService) CreateExample(ctx context.Context, example *model.Example) error {
	return s.exampleRepository.Create(ctx, example)
}

func (s *ExampleService) UpdateExample(ctx context.Context, example *model.Example) error {
	return s.exampleRepository.Update(ctx, example)
}

func (s *ExampleService) DeleteExample(ctx context.Context, id uint) error {
	return s.exampleRepository.Delete(ctx, id)
}
