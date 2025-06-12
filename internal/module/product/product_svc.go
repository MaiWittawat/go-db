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
	ErrPermission      = errors.New("no permission")
)

type productService struct {
	productRepo repository.ProductRepository
	stockSvc    module.StockService
}

// ------------------------ Constructor ------------------------
func NewProductService(productRepo repository.ProductRepository, stockSvc module.StockService) module.ProductService {
	return &productService{
		productRepo: productRepo,
		stockSvc:    stockSvc,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (s *productService) Save(ctx context.Context, pReq *model.ProductReq, userID string) error {
	log.Println("productReq: ", pReq)

	product := pReq.ToProduct()
	product.ID = primitive.NewObjectID().Hex()
	product.CreatedBy = userID
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt


	var baseLogFields = log.Fields{
		"product_id": product.ID,
		"layer":      "product_service",
		"operation":  "product_save",
	}

	if err := s.productRepo.AddProduct(ctx, product); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to save product")
		return ErrCreateProduct
	}

	if err := s.stockSvc.Save(ctx, product.ID, pReq.Quantity); err != nil {
		return err
	}

	log.Info("product and stock created success:")
	return nil
}

func (s *productService) Update(ctx context.Context, pReq *model.ProductReq, id string, userID string) error {
	var baseLogFields = log.Fields{
		"product_id": id,
		"layer":      "product_service",
		"operation":  "product_update",
	}
	updateProduct := pReq.ToProduct()

	if err := updateProduct.Verify(); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to verify product")
		return err
	}

	var currentProduct model.Product
	if err := s.productRepo.GetProductByID(ctx, id, &currentProduct); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get product by id")
		return ErrProductNotFound
	}

	if userID != currentProduct.CreatedBy {
		return ErrPermission
	}

	if updateProduct.Title != "" {
		currentProduct.Title = updateProduct.Title
	}

	if updateProduct.Price != 0 {
		currentProduct.Price = updateProduct.Price
	}

	if updateProduct.Detail != "" {
		currentProduct.Detail = updateProduct.Detail
	}

	currentProduct.UpdatedAt = time.Now()
	if err := s.productRepo.UpdateProduct(ctx, &currentProduct, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to update product")
		return ErrUpdateProduct
	}

	if pReq.Quantity != 0 {
		if err := s.stockSvc.IncreaseQuantity(ctx, pReq.Quantity, id); err != nil {
			return err
		}
	}

	log.Info("product updated successfully:", currentProduct)
	return nil
}

func (s *productService) Delete(ctx context.Context, id string) error {
	var baseLogFields = log.Fields{
		"product_id": id,
		"layer":      "product_service",
		"operation":  "product_delete",
	}

	var product model.Product
	if err := s.productRepo.GetProductByID(ctx, id, &product); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get product by id")
		return ErrProductNotFound
	}

	if err := s.productRepo.DeleteProduct(ctx, id, &product); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to delete product")
		return ErrDeleteProduct
	}

	log.Info("product deleted success:", product)
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (s *productService) GetAll(ctx context.Context) ([]model.ProductRes, error) {
	var baseLogFields = log.Fields{
		"layer":     "product_service",
		"operation": "product_getAll",
	}

	products, err := s.productRepo.GetAllProduct(ctx)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get all product")
		return nil, ErrProductNotFound
	}

	var productsRes []model.ProductRes
	for _, product := range products {
		productRes := product.ToProductRes()
		productsRes = append(productsRes, *productRes)
	}

	log.Info("get all product success")
	return productsRes, nil
}

func (s *productService) GetByID(ctx context.Context, id string) (*model.ProductRes, error) {
	var baseLogFields = log.Fields{
		"product_id": id,
		"layer":      "product_service",
		"operation":  "product_getByID",
	}

	var product model.Product
	if err := s.productRepo.GetProductByID(ctx, id, &product); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get product by id")
		return nil, ErrProductNotFound
	}

	productRes := product.ToProductRes()
	log.Info("get product by id success:", product)
	return productRes, nil
}
