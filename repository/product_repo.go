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
	return pr.db.Create(ctx, pr.collection, p)
}

func (pr *ProductRepo) UpdateProduct(ctx context.Context, p *model.Product, id string) error {
	return pr.db.Update(ctx, pr.collection, p, id)
}

func (pr *ProductRepo) DeleteProduct(ctx context.Context, id string) error {
	var product model.Product
	if err := pr.db.GetByID(ctx, pr.collection, id, &product); err != nil {
		return err
	}
	return pr.db.Delete(ctx, pr.collection, product, id)
}

func (pr *ProductRepo) GetAllProduct(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	if err := pr.db.GetAll(ctx, pr.collection, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (pr *ProductRepo) GetProductByID(ctx context.Context, id string, product *model.Product) (err error) {
	return pr.db.GetByID(ctx, pr.collection, id, product)
}

