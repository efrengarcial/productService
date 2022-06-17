package products

import (
	"context"
)

// ProductService is the top level signature of this service
type ProductService interface {
	Get(ctx context.Context, id string) (*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	Update(ctx context.Context, product *UpdateProduct) error
	Create(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
}
