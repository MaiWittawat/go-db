package product

import (
	"context"
	"go-rebuild/model"
)

type ProductService interface {
	Save(ctx context.Context, p *model.Product) error
	Update(ctx context.Context, p *model.Product, id string) error
	Delete(ctx context.Context, id string) error
}