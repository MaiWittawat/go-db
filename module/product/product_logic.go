package product

import (
	"context"
	"go-rebuild/model"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type productService struct {
	productRepo ProductRepository 
}

// ------------------------ Constructor ------------------------
func NewProductService(productRepo ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

// ------------------------ Method Basic CUD ------------------------
func (ps *productService) Save(ctx context.Context, p *model.Product) error {
	p.ID = primitive.NewObjectID().Hex()
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
	var product model.Product
	if err := ps.productRepo.GetProductByID(ctx, id, &product); err != nil {
		return err
	}

	if err := ps.productRepo.DeleteProduct(ctx, id, &product); err != nil {
		return err
	}
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (ps *productService) GetAll(ctx context.Context) ([]model.Product, error) {
	products, err := ps.productRepo.GetAllProduct(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (ps *productService) GetByID(ctx context.Context, id string, product *model.Product) (err error) {
	if err = ps.productRepo.GetProductByID(ctx, id, product); err != nil {
		return err
	}
	return nil
}
