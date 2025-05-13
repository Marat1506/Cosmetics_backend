package user

import "context"

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserById(ctx context.Context, id string) (User, error)
	Login(ctx context.Context, email string, password string) (User, error)
	AddToFavorites(ctx context.Context, userID string, productID string) error
	RemoveFromFavorites(ctx context.Context, userID string, productID string) error
	AddToCart(ctx context.Context, userID string, productID string) error
	RemoveFromCart(ctx context.Context, userID string, productID string) error
	UpdateCart(ctx context.Context, userID string, cart []string) error
	GetFavorites(ctx context.Context, userID string) ([]string, error)
	GetCart(ctx context.Context, userID string) ([]string, error)
}
