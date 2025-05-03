package product

import "context"

type Storage interface {
	Create(ctx context.Context, product Product) (string, error)
	GetAll(ctx context.Context, category string) ([]Product, error)
}
