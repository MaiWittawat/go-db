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

// ------------------------ Constructor ------------------------
func NewStockRepo(db dbRepo.DB, cacheSvc cache.Cache) repository.StockRepository {
	keyGen := cache.NewKeyGenerator("stocks")
	return &StockRepo{
		db:         db,
		collection: "stocks",
		cacheSvc:   cacheSvc,
		keyGen:     keyGen,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (r *StockRepo) AddStock(ctx context.Context, s *model.Stock) error {
	// save to db
	if err := r.db.Create(ctx, r.collection, s); err != nil {
		return err
	}

	// clear last cahce list
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("[Repo]: failed to clear cache stock in AddStock: ", err)
	}

	// set cache
	cacheKeyProductID := r.keyGen.KeyField("product_id", s.ProductID)
	if err := r.cacheSvc.Set(ctx, cacheKeyProductID, s, 15*time.Minute); err != nil {
		log.Warn("[Repo]: failed to set stock cacheKeyProductID in AddStock: ", err)
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
		}).Error("[Repo]: failed to get stock by id")
		return err
	}

	// clear old cache
	cacheKeyProductID := r.keyGen.KeyField("product_id", s.ProductID)
	if err := r.cacheSvc.Delete(ctx, cacheKeyProductID); err != nil {
		log.Warn("[Repo]: failed to clear cache user in UpdateStock: ", err)
	}

	// update stock data in db
	if err := r.db.Update(ctx, r.collection, s, id); err != nil {
		return err
	}

	// clear stock cache
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("[Repo]: failed to clear cache stock in UpdateStock: ", err)
	}

	// set cache
	cacheKeyProductID = r.keyGen.KeyField("product_id", s.ProductID)
	if err := r.cacheSvc.Set(ctx, cacheKeyProductID, s, 15*time.Minute); err != nil {
		log.Warn("[Repo]: failed to set cache stock in UpdateStock: ", err)
	}

	return nil
}

func (r *StockRepo) DeleteStock(ctx context.Context, id string) error {
	// delete data from db
	if err := r.db.Delete(ctx, r.collection, &model.Stock{}, id); err != nil {
		return err
	}

	// delete cacheKeyList in redis
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("[Repo]: failed to clear stock cachelist in DeletStock: ", err)
	}

	// delete cacheKeyID in redis
	cacheKeyProductID := r.keyGen.KeyField("product_id", id)
	if err := r.cacheSvc.Delete(ctx, cacheKeyProductID); err != nil {
		log.Warn("[Repo]: failed to clear stock cacheKeyProductID in DeletStock: ", err)
	} 

	return nil
}

// ------------------------ Method Basic Query ------------------------
func (r *StockRepo) GetStockByProductID(ctx context.Context, productID string, stock *model.Stock) error {
	// get stock from redis
	cacheKeyProductID := r.keyGen.KeyField("product_id", productID)
	if err := r.cacheSvc.Get(ctx, cacheKeyProductID, &stock); err == nil {
		log.Info("[Repo]: stock from cache: ", stock)
		return nil
	}

	// get stock from db
	if err := r.db.GetByField(ctx, r.collection, "product_id", productID, stock); err != nil {
		log.Info("[Repo]: stock from db: ", stock)
		return err
	}

	// set stock cache in redis
	if err := r.cacheSvc.Set(ctx, cacheKeyProductID, stock, 15*time.Minute); err != nil {
		log.Warn("[Repo]: failed to set stock cacheKeyProductID in GetStockByProductID: ", err)
	} 

	return nil
}
