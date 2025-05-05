package order

import "context"

type Storage interface {
	Create(ctx context.Context, order Order) (string, error)
	GetOrders(ctx context.Context) ([]Order, error)
	ChangeOrder(ctx context.Context, id string) (Order, error)
	DeleteOrder(ctx context.Context, id string) error
}
