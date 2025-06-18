package stock

import (
	"context"
	"errors"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// error
	ErrCreateStock   = errors.New("fail to create stock")
	ErrUpdateStock   = errors.New("fail to update stock")
	ErrDeleteStock   = errors.New("fail to delete stock")
	ErrStockNotFound = errors.New("stock not found")
	ErrStockQuantity = errors.New("fail to increase or decrease stock quantity")
)

type stockService struct {
	repo repository.StockRepository
}

// ------------------------ Constructor ------------------------
func NewStockService(repo repository.StockRepository) module.StockService {
	return &stockService{repo: repo}
}

// ------------------------ Method Basic CUD ------------------------
func (s *stockService) Save(ctx context.Context, productID string, quantity int) error {
	var stock = model.Stock{
		ProductID: productID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	var baseLogFields = log.Fields{
		"product_id": productID,
		"layer":      "stock_service",
		"method":     "stock_save",
	}

	stock.SetQuantity(quantity)

	if err := s.repo.AddStock(ctx, &stock); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("add stock")
		return ErrCreateStock
	}

	return nil
}

func (s *stockService) Update(ctx context.Context, productID string, quantity int) error {
	var currentStock model.Stock
	var baseLogFields = log.Fields{
		"product_id": productID,
		"layer":      "stock_service",
		"method":     "stock_update",
	}

	if err := s.repo.GetStockByProductID(ctx, productID, &currentStock); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get stock by product id")
		return ErrUpdateStock
	}

	currentStock.SetQuantity(quantity)
	if err := s.repo.UpdateStock(ctx, &currentStock); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("update stock")
		return ErrUpdateStock
	}

	return nil
}

func (s *stockService) IncreaseQuantity(ctx context.Context, quantity int, productID string) error {
	var currentStock model.Stock
	var baseLogFields = log.Fields{
		"product_id": productID,
		"layer":      "stock_service",
		"method":     "stock_increaseQuantity",
	}

	if err := s.repo.GetStockByProductID(ctx, productID, &currentStock); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get stock by product id")
		return ErrStockNotFound
	}

	currentStock.IncreaseQuantity(quantity)
	currentStock.UpdatedAt = time.Now()

	if err := s.repo.UpdateStock(ctx, &currentStock); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("update stock")
		return ErrUpdateStock
	}

	return nil
}

func (s *stockService) DecreaseQuantity(ctx context.Context, quantity int, productID string) error {
	var currentStock model.Stock
	var baseLogFields = log.Fields{
		"product_id": productID,
		"layer":      "stock_service",
		"method":     "stock_decreaseQuantity",
	}

	if err := s.repo.GetStockByProductID(ctx, productID, &currentStock); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get stock by product id")
		return ErrStockNotFound
	}

	if err := currentStock.DecreaseQuantity(quantity); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("decrease quantity")
		return ErrUpdateStock
	}

	currentStock.UpdatedAt = time.Now()
	if err := s.repo.UpdateStock(ctx, &currentStock); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("update stock")
		return ErrUpdateStock
	}

	return nil
}

func (s *stockService) Delete(ctx context.Context, id string) error {
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (s *stockService) GetAll(ctx context.Context) ([]model.Stock, error) {
	stocks, err := s.repo.GetAllStock(ctx)
	if err != nil {
		return nil, ErrStockNotFound
	}
	return stocks, nil
}

func (s *stockService) GetByProductID(ctx context.Context, productID string) (*model.Stock, error) {
	var stock model.Stock
	if err := s.repo.GetStockByID(ctx, productID, &stock); err != nil {
		return nil, ErrStockNotFound
	}
	return &stock, nil
}
