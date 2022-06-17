package products

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var (
	validate *validator.Validate
)

type repository interface {
	Get(ctx context.Context, id string) (*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	Update(ctx context.Context, product *UpdateProduct) error
	Create(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
}

// Usecase for interacting with users
type Usecase struct {
	Repository repository
}


// Get a single user
func (u *Usecase) Get(ctx context.Context, id string) (*Product, error) {
	user, err := u.Repository.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching a single product")
	}
	return user, nil
}

// GetAll gets all users
func (u *Usecase) GetAll(ctx context.Context) ([]*Product, error) {
	users, err := u.Repository.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error fetching all products")
	}
	return users, nil
}

// Update a single user
func (u *Usecase) Update(ctx context.Context, product *UpdateProduct) error {
	validate = validator.New()
	if err := validate.Struct(product); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}

	if err := u.Repository.Update(ctx, product); err != nil {
		return errors.Wrap(err, "error updating product")
	}
	return nil
}

// Create a single user
func (u *Usecase) Create(ctx context.Context, product *Product) error {
	validate = validator.New()
	if err := validate.Struct(*product); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}

	if err := u.Repository.Create(ctx, product); err != nil {
		return errors.Wrap(err, "error creating new product")
	}

	return nil
}

// Delete a single user
func (u *Usecase) Delete(ctx context.Context, id string) error {
	if err := u.Repository.Delete(ctx, id); err != nil {
		return errors.Wrap(err, "error deleting product")
	}
	return nil
}
