package product

import (
	"context"
	"go-rebuild/model"
	"time"
)


type productService struct {
	productRepo ProductRepository 
}

func NewProductService(productRepo ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (ps *productService) Save(ctx context.Context, p *model.Product) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt
	if err := ps.productRepo.AddProduct(ctx, p); err != nil {
		return err
	}
	return nil
}

func (ps *productService) Update(ctx context.Context, p *model.Product, id string) error {
	p.UpdatedAt = time.Now()
	if err := ps.productRepo.UpdateProduct(ctx, p, id); err != nil {
		return err
	}
	return nil
}

func (ps *productService) Delete(ctx context.Context, id string) error {
	if err := ps.productRepo.DeleteProduct(ctx, id); err != nil {
		return err
	}
	return nil
}