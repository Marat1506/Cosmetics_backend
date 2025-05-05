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
		PasswordHash: dto.Password, // В реальном приложении здесь должен быть хэш пароля
	}
	return s.storage.Create(ctx, user)
}

func (s *Service) GetAllUsers(ctx context.Context) ([]User, error) {
	return s.storage.GetAllUsers(ctx)
}

func (s *Service) Login(ctx context.Context, dto LoginDTO) (User, error) {
	return s.storage.Login(ctx, dto.Email, dto.Password)
}
