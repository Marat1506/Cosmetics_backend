package user

import "context"

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	GetByID(ctx context.Context, id string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	AddToFavorites(ctx context.Context, userID, productID string) error
	RemoveFromFavorites(ctx context.Context, userID, productID string) error
	AddToCart(ctx context.Context, userID, productID string) error
	RemoveFromCart(ctx context.Context, userID, productID string) error
}
