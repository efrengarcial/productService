package products

import (
	"context"
	"io"

	"gocloud.dev/docstore"
)

// DynamoDBRepository -
type DynamoDBRepository struct {
	coll *docstore.Collection
}

// NewDynamoDBRepository -
func NewDynamoDBRepository(c *docstore.Collection) *DynamoDBRepository {
	return &DynamoDBRepository{c}
}

// Get a user
func (r *DynamoDBRepository) Get(ctx context.Context, id string) (*Product, error) {
	product := &Product{
		ID: id,
	}

	err := r.coll.Get(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetAll products
func (r *DynamoDBRepository) GetAll(ctx context.Context) ([]*Product, error) {
	products := make([]*Product, 0)
	iter := r.coll.Query().Get(ctx)
	defer iter.Stop() // Always call Stop on an iterator.
	// Query.Get returns an iterator. Call Next on it until io.EOF.
	for {
		var p Product
		err := iter.Next(ctx, &p)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		products = append(products,  &p)
	}

	return products, nil
}

// Update a user
func (r *DynamoDBRepository) Update(ctx context.Context,product *UpdateProduct) error {
	// Set the score to a new value.
	return r.coll.Actions().Put(product).Do(ctx)
}

// Create a user
func (r *DynamoDBRepository) Create(ctx context.Context, product *Product) error {
	return r.coll.Create(ctx, product)
}

// Delete a user
func (r *DynamoDBRepository) Delete(ctx context.Context, id string) error {
	product := &Product{
		ID: id,
	}
	return r.coll.Actions().Delete(product).Do(ctx)
}
