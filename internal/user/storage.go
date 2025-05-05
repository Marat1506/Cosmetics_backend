package user

import "context"

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, id string) (User, error)
	Login(ctx context.Context, email string, password string) (User, error)
}
