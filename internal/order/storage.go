package order

import "context"

type Storage interface {
	Create(ctx context.Context, order Order) (string, error)
	GetByID(ctx context.Context, id string) (Order, error)
	GetByUserID(ctx context.Context, userID string) ([]Order, error)
	UpdateStatus(ctx context.Context, id, status string) error
	Cancel(ctx context.Context, id string) error
}
