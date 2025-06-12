package stock

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
	ErrCreateStock   = errors.New("fail to create stock")
	ErrUpdateStock   = errors.New("fail to update stock")
	ErrDeleteStock   = errors.New("fail to delete stock")
	ErrStockNotFound = errors.New("stock not found")
	ErrStockQuantity = errors.New("fail to increase or decrease stock quantity")
)

type stockService struct {
	repo repository.StockRepository
}

func NewStockService(repo repository.StockRepository) module.StockService {
	return &stockService{repo: repo}
}

func (s *stockService) Save(ctx context.Context, productID string, quantity int) error {
	var stock = model.Stock{
		ID:        primitive.NewObjectID().Hex(),
		ProductID: productID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var baseLogFields = log.Fields{
		"stock_id":  stock.ID,
		"layer":     "stock_service",
		"operation": "stock.create",
	}

	stock.SetQuantity(quantity)

	if err := s.repo.AddStock(ctx, &stock); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to create stock")
		return ErrCreateStock
	}

	return nil
}

func (s *stockService) IncreaseQuantity(ctx context.Context, quantity int, productID string) error {
	var currentStock model.Stock
	var baseLogFields = log.Fields{
		"product_id":  productID,
		"layer":     "stock_service",
		"operation": "stock.update",
	}

	if err := s.repo.GetStockByProductID(ctx, productID, &currentStock); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("stock not found using id")
		return ErrStockNotFound
	}

	currentStock.IncreaseQuantity(quantity)
	currentStock.UpdatedAt = time.Now()

	if err := s.repo.UpdateStock(ctx, &currentStock, currentStock.ID); err != nil {
		return ErrUpdateStock
	}

	return nil
}

func (s *stockService) DecreaseQuantity(ctx context.Context, quantity int, productID string) error {
	var currentStock model.Stock
	var baseLogFields = log.Fields{
		"product_id":  productID,
		"layer":     "stock_service",
		"operation": "stock.update",
	}

	if err := s.repo.GetStockByProductID(ctx, productID, &currentStock); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("stock not found using id")
		return ErrStockNotFound
	}

	if err := currentStock.DecreaseQuantity(quantity); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("fail to decrease quantity in stock")
		return ErrUpdateStock
	}

	currentStock.UpdatedAt = time.Now()
	if err := s.repo.UpdateStock(ctx, &currentStock, currentStock.ID); err != nil {
		return ErrUpdateStock
	}

	return nil
}

func (s *stockService) Delete(ctx context.Context, id string) error {
	return nil
}
