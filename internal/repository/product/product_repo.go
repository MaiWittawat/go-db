package product

import (
	"context"
	rclient "go-rebuild/internal/cache"
	dbRepo "go-rebuild/internal/db"
	"go-rebuild/internal/model"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
)

type productRepo struct {
	db         dbRepo.DB
	collection string
	cache      rclient.Cache
	keyGen     *rclient.KeyGenerator
}

// ------------------------ Constructor ------------------------
func NewProductRepo(db dbRepo.DB, cache rclient.Cache) repository.ProductRepository {
	keyGen := rclient.NewKeyGenerator("products")
	return &productRepo{
		db:         db,
		collection: "products",
		cache:      cache,
		keyGen:     keyGen,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (pr *productRepo) AddProduct(ctx context.Context, p *model.Product) error {
	// save to db
	if err := pr.db.Create(ctx, pr.collection, p); err != nil {
		return err
	}

	// clear last cache list
	cacheKeyList := pr.keyGen.KeyList()
	if err := pr.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache productslist in AddProduct: ", err)
	}

	// set cache
	cacheKeyID := pr.keyGen.KeyID(p.ID)
	if err := pr.cache.Set(ctx, cacheKeyID, p, 15*time.Minute); err != nil {
		log.Warn("failed to set product cacheKeyID in AddProduct: ", err)
	}

	return nil
}

func (pr *productRepo) UpdateProduct(ctx context.Context, p *model.Product, id string) error {
	// update product in db
	if err := pr.db.Update(ctx, pr.collection, p, id); err != nil {
		return err
	}

	// clear old cache product list
	cacheKeyList := pr.keyGen.KeyList()
	if err := pr.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache products list in UpdateProduct: ", err)
	}

	// set cache
	cacheKeyID := pr.keyGen.KeyID(id)
	if err := pr.cache.Set(ctx, cacheKeyID, p, 15*time.Minute); err != nil {
		log.Warn("failed to set product cacheKeyID in UpdateProduct: ", err)
	}

	return nil
}

func (pr *productRepo) DeleteProduct(ctx context.Context, id string, product *model.Product) error {
	// delete product from db
	if err := pr.db.Delete(ctx, pr.collection, product, id); err != nil {
		return err
	}

	// clear cache list in redis
	cacheKeyList := pr.keyGen.KeyList()
	if err := pr.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache products in DeleteProduct: ", err)
	}

	// clear cache key id in redis
	cacheKeyID := pr.keyGen.KeyID(id)
	if err := pr.cache.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache product in DeleteProduct: ", err)
	}

	return nil
}

// ------------------------ Method Basic Query ------------------------
func (pr *productRepo) GetAllProduct(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	cacheKeyList := pr.keyGen.KeyList()

	// get products from redis
	if err := pr.cache.Get(ctx, cacheKeyList, &products); err == nil {
		log.Info("products from redis: ", products)
		return products, nil
	}

	// get products from db
	if err := pr.db.GetAll(ctx, pr.collection, &products); err != nil {
		log.Info("products from db: ", products)
		return nil, err
	}

	// set cache products in redis
	if err := pr.cache.Set(ctx, cacheKeyList, products, 15*time.Minute); err != nil {
		log.Warn("failed to set cache products in GetAllProduct: ", err)
	}

	return products, nil
}

func (pr *productRepo) GetProductByID(ctx context.Context, id string, product *model.Product) (err error) {
	// get product from redis
	cacheKeyID := pr.keyGen.KeyID(id)
	if err := pr.cache.Get(ctx, cacheKeyID, product); err == nil {
		log.Info("product from redis: ", product)
		return nil
	}

	// get product from db if fail to get that from redis
	if err := pr.db.GetByID(ctx, pr.collection, id, product); err != nil {
		log.Info("product from db: ", product)
		return err
	}

	// set product cache in redis
	if err := pr.cache.Set(ctx, cacheKeyID, product, 15*time.Minute); err != nil {
		log.Warn("Warning: fail to set cache product in GetProductByID: ", err)
	}

	return nil
}
