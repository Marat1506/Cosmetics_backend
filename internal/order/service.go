package order

import (
	"context"
	"server/pkg/logging"
	"time"
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

func (s *Service) CreateOrder(ctx context.Context, dto CreateOrderDTO) (string, error) {
	order := Order{
		UserID:     dto.UserID,
		Products:   dto.Products,
		TotalPrice: dto.TotalPrice,
		Status:     "created",
		CreatedAt:  time.Now(),
	}

	return s.storage.Create(ctx, order)
}

func (s *Service) GetOrder(ctx context.Context, id string) (Order, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *Service) GetUserOrders(ctx context.Context, userID string) ([]Order, error) {
	return s.storage.GetByUserID(ctx, userID)
}

func (s *Service) UpdateStatus(ctx context.Context, id, status string) error {
	return s.storage.UpdateStatus(ctx, id, status)
}

func (s *Service) CancelOrder(ctx context.Context, id string) error {
	return s.storage.Cancel(ctx, id)
}
