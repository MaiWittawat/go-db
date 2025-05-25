package repository

import (
	"context"
	dbRepo "go-rebuild/db"
	"go-rebuild/model"
	module "go-rebuild/module/product"
)

type ProductRepo struct {
	db         dbRepo.DB
	collection string
}

func NewProductRepo(db dbRepo.DB) module.ProductRepository {
	return &ProductRepo{db: db, collection: "products"}
}

func (pr *ProductRepo) AddProduct(ctx context.Context, p *model.Product) error {
	if err := pr.db.Create(ctx, pr.collection, p); err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepo) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	return nil, nil
}

func (pr *ProductRepo) UpdateProduct(ctx context.Context, p *model.Product, id string) error {
	if err := pr.db.Update(ctx, pr.collection, p, id); err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepo) DeleteProduct(ctx context.Context, id string) error {
	if err := pr.db.Delete(ctx, pr.collection, id); err != nil {
		return err
	}
	return nil
}
