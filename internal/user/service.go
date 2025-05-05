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

func (s *Service) CreateUser(ctx context.Context, dto CreateUserDTO) (string, error) {
	user := User{
		Email:        dto.Email,
		Username:     dto.Username,
		PasswordHash: dto.Password,
		Favorites:    []string{},
		Cart:         []string{},
	}
	return s.storage.Create(ctx, user)
}

func (s *Service) GetAllUsers(ctx context.Context) ([]User, error) {
	return s.storage.GetAllUsers(ctx)
}

func (s *Service) Login(ctx context.Context, dto LoginDTO) (User, error) {
	return s.storage.Login(ctx, dto.Email, dto.Password)
}

func (s *Service) AddToFavorites(ctx context.Context, userID string, productID string) error {
	return s.storage.AddToFavorites(ctx, userID, productID)
}

func (s *Service) RemoveFromFavorites(ctx context.Context, userID string, productID string) error {
	return s.storage.RemoveFromFavorites(ctx, userID, productID)
}

func (s *Service) AddToCart(ctx context.Context, userID string, productID string) error {
	return s.storage.AddToCart(ctx, userID, productID)
}
func (s *Service) GetFavorites(ctx context.Context, userID string) ([]string, error) {
	return s.storage.GetFavorites(ctx, userID)
}

func (s *Service) RemoveFromCart(ctx context.Context, userID string, productID string) error {
	return s.storage.RemoveFromCart(ctx, userID, productID)
}

func (s *Service) UpdateCart(ctx context.Context, userID string, cart []string) error {
	return s.storage.UpdateCart(ctx, userID, cart)
}
