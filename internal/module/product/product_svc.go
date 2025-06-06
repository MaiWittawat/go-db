package product

import (
	"context"
	"errors"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrCreateProduct = errors.New("fail to create product")
	ErrUpdateProduct = errors.New("fail to update product")
	ErrDeleteProduct = errors.New("fail to delete product")

	ErrProductNotFound = errors.New("product not found")
)


type productService struct {
	productRepo repository.ProductRepository
}

// ------------------------ Constructor ------------------------
func NewProductService(productRepo repository.ProductRepository) module.ProductService {
	return &productService{productRepo: productRepo}
}

// ------------------------ Method Basic CUD ------------------------
func (ps *productService) Save(ctx context.Context, p *model.Product) error {
	p.ID = primitive.NewObjectID().Hex()
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt

	var baseLogFields = log.Fields{
		"product_id": p.ID,
		"layer":      "product_service",
		"operation":  "product_save",
	}

	if err := ps.productRepo.AddProduct(ctx, p); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to save product")
		return ErrCreateProduct
	}

	log.Info("product created success:", p)
	return nil
}

func (ps *productService) Update(ctx context.Context, req *model.Product, id string) error {
	var baseLogFields = log.Fields{
		"product_id": id,
		"layer":      "product_service",
		"operation":  "product_update",
	}

	if err := req.Verify(); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to verify product")
		return err
	}

	var currentProduct model.Product
	if err := ps.productRepo.GetProductByID(ctx, id, &currentProduct); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get product by id")
		return ErrProductNotFound
	}

	if req.Title != "" {
		currentProduct.Title = req.Title
	}

	if req.Price != 0 {
		currentProduct.Price = req.Price
	}

	if req.Detail != "" {
		currentProduct.Detail = req.Detail
	}

	currentProduct.UpdatedAt = time.Now()
	if err := ps.productRepo.UpdateProduct(ctx, &currentProduct, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to update product")
		return ErrUpdateProduct
	}

	log.Info("product updated successfully:", currentProduct)
	return nil
}

func (ps *productService) Delete(ctx context.Context, id string) error {
	var baseLogFields = log.Fields{
		"product_id": id,
		"layer":      "product_service",
		"operation":  "product_delete",
	}

	var product model.Product
	if err := ps.productRepo.GetProductByID(ctx, id, &product); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get product by id")
		return ErrProductNotFound
	}

	if err := ps.productRepo.DeleteProduct(ctx, id, &product); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to delete product")
		return ErrDeleteProduct
	}

	log.Info("product deleted success:", product)
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (ps *productService) GetAll(ctx context.Context) ([]model.Product, error) {
	var baseLogFields = log.Fields{
		"layer":     "product_service",
		"operation": "product_getAll",
	}

	products, err := ps.productRepo.GetAllProduct(ctx)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get all product")
		return nil, ErrProductNotFound
	}

	log.Info("get all product success")
	return products, nil
}

func (ps *productService) GetByID(ctx context.Context, id string, product *model.Product) (err error) {
	var baseLogFields = log.Fields{
		"product_id": id,
		"layer":      "product_service",
		"operation":  "product_getByID",
	}

	if err = ps.productRepo.GetProductByID(ctx, id, product); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get product by id")
		return ErrProductNotFound
	}

	log.Info("get product by id success:", product)
	return nil
}
