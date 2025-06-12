package stock

import (
	"context"
	"go-rebuild/internal/cache"
	dbRepo "go-rebuild/internal/db"
	"go-rebuild/internal/model"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
)

type StockRepo struct {
	db         dbRepo.DB
	collection string
	cacheSvc   cache.Cache
	keyGen     *cache.KeyGenerator
}

func NewStockRepo(db dbRepo.DB, cacheSvc cache.Cache) repository.StockRepository {
	keyGen := cache.NewKeyGenerator("stocks")
	return &StockRepo{
		db:         db,
		collection: "stocks",
		cacheSvc:   cacheSvc,
		keyGen:     keyGen,
	}
}

func (r *StockRepo) AddStock(ctx context.Context, s *model.Stock) error {
	// save to db
	if err := r.db.Create(ctx, r.collection, s); err != nil {
		return err
	}

	// clear last cahce list
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache stock in AddStock: ", err)
	}

	// set cache
	cacheKeyID := r.keyGen.KeyID(s.ID)
	if err := r.cacheSvc.Set(ctx, cacheKeyID, s, 15*time.Minute); err != nil {
		log.Warn("failed to set cache stock in AddStock: ", err)
	} else {
		log.Info("set cache in AddStock success")
	}

	return nil
}

func (r *StockRepo) UpdateStock(ctx context.Context, s *model.Stock, id string) error {
	var currentStock model.Stock
	if err := r.db.GetByID(ctx, r.collection, id, &currentStock); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"stock_id": id,
			"layer":    "repository",
			"step":     "update.stock",
		}).Error("failed to get stock by id")
		return err
	}

	// clear old cache
	cacheKeyProductID := r.keyGen.KeyField("product_id", s.ProductID)
	if err := r.cacheSvc.Delete(ctx, cacheKeyProductID); err != nil {
		log.Warn("failed to clear cache user in UpdateStock: ", err)
	}

	// update stock data in db
	if err := r.db.Update(ctx, r.collection, s, id); err != nil {
		return err
	}
	log.Info("stock update stock: ", s)

	// clear stock cache
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache stock in UpdateStock: ", err)
	}

	// set cache
	cacheKeyProductID = r.keyGen.KeyField("product_id", s.ProductID)
	if err := r.cacheSvc.Set(ctx, cacheKeyProductID, s, 15*time.Minute); err != nil {
		log.Warn("failed to set cache stock in UpdateStock: ", err)
	}

	cacheKeyID := r.keyGen.KeyID(id)
	if err := r.cacheSvc.Set(ctx, cacheKeyID, s, 15*time.Minute); err != nil {
		log.Warn("failed to clear cache stock in UpdateStock: ", err)
	} else {
		log.Info("set cache in UpdateStock success")
	}

	return nil
}

func (r *StockRepo) DeleteStock(ctx context.Context, id string, s *model.Stock) error {
	if err := r.db.Delete(ctx, r.collection, s, id); err != nil {
		return err
	}

	cacheKeyID := r.keyGen.KeyID(id)
	if err := r.cacheSvc.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache stock in DeletStock: ", err)
	} else {
		log.Info("clear cache in DeleteStock success")
	}

	return nil
}

func (r *StockRepo) GetStockByProductID(ctx context.Context, productID string, stock *model.Stock) error {
	cacheKeyProductID := r.keyGen.KeyField("product_id", productID)
	if err := r.cacheSvc.Get(ctx, cacheKeyProductID, &stock); err == nil {
		log.Info("stock from cache : ", stock)
		return nil
	}
	if err := r.db.GetByField(ctx, r.collection, "product_id", productID, stock); err != nil {
		return err
	}

	log.Info("stock from db : ", stock)
	if err := r.cacheSvc.Set(ctx, cacheKeyProductID, stock, 15*time.Minute); err != nil {
		log.Warn("failed to set cache stock in GetStockByProductID: ", err)
	} else {
		log.Info("set cache in GetStockByProductID success")
	}

	return nil
}
