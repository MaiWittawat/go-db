package product

import (
	"context"
	"encoding/json"
	"errors"
	messagebroker "go-rebuild/internal/message_broker"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// error
	ErrCreateProduct   = errors.New("fail to create product")
	ErrUpdateProduct   = errors.New("fail to update product")
	ErrDeleteProduct   = errors.New("fail to delete product")
	ErrProductNotFound = errors.New("product not found")
	ErrPermission      = errors.New("no permission")
	ErrMarShal         = errors.New("failed to marshal object")
)

type productService struct {
	productRepo repository.ProductRepository
	producerSvc messagebroker.ProducerService
}

// ------------------------ Constructor ------------------------
func NewProductService(productRepo repository.ProductRepository, producerSvc messagebroker.ProducerService) module.ProductService {
	return &productService{
		productRepo: productRepo,
		producerSvc: producerSvc,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (s *productService) Save(ctx context.Context, pReq *model.ProductReq, userID string) error {
	product := pReq.ToProduct()
	product.ID = primitive.NewObjectID().Hex()
	product.CreatedBy = userID
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt

	var baseLogFields = log.Fields{
		"product_id": product.ID,
		"layer":      "product_service",
		"method":     "product_save",
	}

	if err := s.productRepo.AddProduct(ctx, product); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("add product")
		return ErrCreateProduct
	}
	log.Printf("[Service]: product {%s} created success", product.ID)

	bodyByte, err := json.Marshal(&model.Stock{ProductID: product.ID, Quantity: pReq.Quantity})
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("json marshal")
		return ErrMarShal
	}

	mqConf := &model.MQConfig{
		ExchangeName: messagebroker.StockExchangeName,
		ExchangeType: messagebroker.StockExchangeType,
		QueueName:    messagebroker.StockQueueName,
		RoutingKey:   "stock.create",
	}

	if err := s.producerSvc.Publishing(ctx, mqConf, bodyByte); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("publishing")
		return err
	}

	return nil
}

func (s *productService) Update(ctx context.Context, pReq *model.ProductReq, id string, userID string) error {
	var baseLogFields = log.Fields{
		"product_id": id,
		"layer":      "product_service",
		"method":     "product_update",
	}
	updateProduct := pReq.ToProduct()

	if pReq.Quantity < 0 {
		return ErrUpdateProduct
	}

	if err := updateProduct.Verify(); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("verify")
		return err
	}

	var currentProduct model.Product
	if err := s.productRepo.GetProductByID(ctx, id, &currentProduct); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get product by id")
		return ErrProductNotFound
	}

	if userID != currentProduct.CreatedBy {
		return ErrPermission
	}

	currentProduct.UpdateNotNilField(pReq)
	if err := s.productRepo.UpdateProduct(ctx, &currentProduct, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("update product")
		return ErrUpdateProduct
	}

	bodyByte, err := json.Marshal(&model.Stock{ProductID: currentProduct.ID, Quantity: pReq.Quantity})
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("json marshal")
		return ErrMarShal
	}

	mqConf := &model.MQConfig{ExchangeName: messagebroker.StockExchangeName, ExchangeType: messagebroker.StockExchangeType, QueueName: messagebroker.StockQueueName, RoutingKey: "stock.update"}
	if err := s.producerSvc.Publishing(ctx, mqConf, bodyByte); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("publishing")
		return ErrUpdateProduct
	}

	log.Printf("[Service]: product {%s} updated success\n", currentProduct.ID)
	return nil
}

func (s *productService) Delete(ctx context.Context, id string) error {
	var baseLogFields = log.Fields{
		"product_id": id,
		"layer":      "product_service",
		"method":     "product_delete",
	}

	var product model.Product
	if err := s.productRepo.GetProductByID(ctx, id, &product); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get product by id")
		return ErrProductNotFound
	}

	if err := s.productRepo.DeleteProduct(ctx, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("delete product")
		return ErrDeleteProduct
	}

	log.Printf("[Service]: product {%s} deleted success\n", product.ID)
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (s *productService) GetAll(ctx context.Context) ([]model.ProductResp, error) {
	var baseLogFields = log.Fields{
		"layer":  "product_service",
		"method": "product_getAll",
	}

	products, err := s.productRepo.GetAllProduct(ctx)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get all product")
		return nil, ErrProductNotFound
	}

	var productsRes []model.ProductResp
	for _, product := range products {
		productRes := product.ToProductRes()
		productsRes = append(productsRes, *productRes)
	}

	log.Info("[Service]: get all product success")
	return productsRes, nil
}

func (s *productService) GetByID(ctx context.Context, id string) (*model.ProductResp, error) {
	var baseLogFields = log.Fields{
		"product_id": id,
		"layer":      "product_service",
		"method":     "product_getByID",
	}

	var product model.Product
	if err := s.productRepo.GetProductByID(ctx, id, &product); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get product by id")
		return nil, ErrProductNotFound
	}

	productRes := product.ToProductRes()
	log.Printf("[Service]: get product {%s} success\n", product.ID)
	return productRes, nil
}
