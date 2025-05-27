package product

import (
	"context"
	"go-rebuild/model"
)


type ProductRepository interface {
	AddProduct(ctx context.Context, p *model.Product) error
	UpdateProduct(ctx context.Context, p *model.Product, id string) error
	DeleteProduct(ctx context.Context, id string, p *model.Product) error


	GetAllProduct(ctx context.Context) ([]model.Product, error)
	GetProductByID(ctx context.Context, id string, p *model.Product) error
}