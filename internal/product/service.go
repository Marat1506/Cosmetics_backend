package product

import (
	"context"
	"server/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func NewService(storage Storage, logger *logging.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) Create(ctx context.Context, product Product) (string, error) {
	return s.storage.Create(ctx, product)
}

func (s *Service) GetAll(ctx context.Context) ([]Product, error) {
	return s.storage.GetAll(ctx)
}
