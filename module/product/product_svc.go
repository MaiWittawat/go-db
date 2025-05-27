package product

import (
	"context"
	"go-rebuild/model"
	"time"

	log "github.com/sirupsen/logrus"
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
		log.WithError(err).WithFields(log.Fields{
			"product_id": p.ID,
			"layer":      "service",
			"step":       "Save",
		}).Error("failed to save product")
		return ErrCreateProduct
	}
	return nil
}

func (ps *productService) Update(ctx context.Context, p *model.Product, id string) error {
	p.UpdatedAt = time.Now()
	log.Info("updating product with ID:", id)
	if err := ps.productRepo.GetProductByID(ctx, id, p); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"product_id": id,
			"layer":      "service",
			"step":       "Update",
		}).Error("failed to get product by id")
		return ErrProductNotFound
	}
	log.Info("product found, proceeding to update:", p)
	if err := ps.productRepo.UpdateProduct(ctx, p, id); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"product_id": p.ID,
			"layer":      "service",
			"step":       "Update",
		}).Error("failed to update product")
		return ErrUpdateProduct
	}
	log.Info("product updated successfully:", p)
	return nil
}

func (ps *productService) Delete(ctx context.Context, id string) error {
	var product model.Product
	if err := ps.productRepo.GetProductByID(ctx, id, &product); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"product_id": id,
			"layer":      "service",
			"step":       "Delete",
		}).Error("failed to get product by id")
		return ErrProductNotFound
	}

	if err := ps.productRepo.DeleteProduct(ctx, id, &product); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"product_id": id,
			"layer":      "service",
			"step":       "Delete",
		}).Error("failed to delete product")
		return ErrDeleteProduct
	}
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (ps *productService) GetAll(ctx context.Context) ([]model.Product, error) {
	products, err := ps.productRepo.GetAllProduct(ctx)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"layer": "service",
			"step":  "GetAll",
		}).Error("failed to get all product")
		return nil, ErrProductNotFound
	}
	return products, nil
}

func (ps *productService) GetByID(ctx context.Context, id string, product *model.Product) (err error) {
	if err = ps.productRepo.GetProductByID(ctx, id, product); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"product_id": id,
			"layer":      "service",
			"step":       "GetByID",
		}).Error("failed to get product by id")
		return ErrProductNotFound
	}
	return nil
}
