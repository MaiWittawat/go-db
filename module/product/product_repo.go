package product

import (
	"context"
	"go-rebuild/model"
)


type ProductRepository interface {
	AddProduct(ctx context.Context, p *model.Product) error
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
	UpdateProduct(ctx context.Context, p *model.Product, id string) error
	DeleteProduct(ctx context.Context, id string) error
}