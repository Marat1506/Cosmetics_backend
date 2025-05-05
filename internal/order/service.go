package order

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

func (s *Service) CreateOrder(ctx context.Context, dto CreateOrderDTO) (string, error) {
	order := Order{
		Username:     dto.Username,
		Phone:        dto.Phone,
		TelegramNick: dto.TelegramNick,
		Completed:    false,
	}

	return s.storage.Create(ctx, order)
}

func (s *Service) GetOrders(ctx context.Context) ([]Order, error) {

	orders, err := s.storage.GetOrders(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Service) ChangeOrder(ctx context.Context, id string) (Order, error) {
	orders, err := s.storage.ChangeOrder(ctx, id)
	if err != nil {
		return Order{}, err
	}
	return orders, nil
}

func (s *Service) DeleteOrder(ctx context.Context, id string) error {
	return s.storage.DeleteOrder(ctx, id)

}
