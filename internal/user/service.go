package user

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

func (s *Service) Register(ctx context.Context, dto CreateUserDTO) (string, error) {
	user := User{
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: dto.Password, // В реальном приложении нужно хэшировать
	}

	return s.storage.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, dto LoginDTO) (User, error) {
	return s.storage.GetByEmail(ctx, dto.Email)
}

func (s *Service) GetProfile(ctx context.Context, userID string) (User, error) {
	return s.storage.GetByID(ctx, userID)
}

func (s *Service) AddToFavorites(ctx context.Context, userID, productID string) error {
	return s.storage.AddToFavorites(ctx, userID, productID)
}

func (s *Service) RemoveFromFavorites(ctx context.Context, userID, productID string) error {
	return s.storage.RemoveFromFavorites(ctx, userID, productID)
}

func (s *Service) AddToCart(ctx context.Context, userID, productID string) error {
	return s.storage.AddToCart(ctx, userID, productID)
}

func (s *Service) RemoveFromCart(ctx context.Context, userID, productID string) error {
	return s.storage.RemoveFromCart(ctx, userID, productID)
}
