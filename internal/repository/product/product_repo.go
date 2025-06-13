package product

import (
	"context"
	"go-rebuild/internal/cache"
	dbRepo "go-rebuild/internal/db"
	"go-rebuild/internal/model"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
)

type productRepo struct {
	db         dbRepo.DB
	collection string
	cacheSvc   cache.Cache
	keyGen     *cache.KeyGenerator
}

// ------------------------ Constructor ------------------------
func NewProductRepo(db dbRepo.DB, cacheSvc cache.Cache) repository.ProductRepository {
	keyGen := cache.NewKeyGenerator("products")
	return &productRepo{
		db:         db,
		collection: "products",
		cacheSvc:   cacheSvc,
		keyGen:     keyGen,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (r *productRepo) AddProduct(ctx context.Context, p *model.Product) error {
	// save to db
	if err := r.db.Create(ctx, r.collection, p); err != nil {
		return err
	}

	// clear last cache list
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("[Repo]: failed to clear cache productslist in AddProduct: ", err)
	}

	// set cache
	cacheKeyID := r.keyGen.KeyID(p.ID)
	if err := r.cacheSvc.Set(ctx, cacheKeyID, p, 15*time.Minute); err != nil {
		log.Warn("[Repo]: failed to set product cacheKeyID in AddProduct: ", err)
	}

	return nil
}

func (r *productRepo) UpdateProduct(ctx context.Context, p *model.Product, id string) error {
	// update product in db
	if err := r.db.Update(ctx, r.collection, p, id); err != nil {
		return err
	}

	// clear old cache product list
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("[Repo]: failed to clear cache products list in UpdateProduct: ", err)
	}

	// set cache
	cacheKeyID := r.keyGen.KeyID(id)
	if err := r.cacheSvc.Set(ctx, cacheKeyID, p, 15*time.Minute); err != nil {
		log.Warn("[Repo]: failed to set product cacheKeyID in UpdateProduct: ", err)
	}

	return nil
}

func (r *productRepo) DeleteProduct(ctx context.Context, id string, product *model.Product) error {
	// delete product from db
	if err := r.db.Delete(ctx, r.collection, r, id); err != nil {
		return err
	}

	// clear cache list in redis
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("[Repo]: failed to clear cache products in DeleteProduct: ", err)
	}

	// clear cache key id in redis
	cacheKeyID := r.keyGen.KeyID(id)
	if err := r.cacheSvc.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("[Repo]: failed to clear cache product in DeleteProduct: ", err)
	}

	return nil
}

// ------------------------ Method Basic Query ------------------------
func (r *productRepo) GetAllProduct(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	cacheKeyList := r.keyGen.KeyList()

	// get products from redis
	if err := r.cacheSvc.Get(ctx, cacheKeyList, &products); err == nil {
		log.Info("[Repo]: products from redis: ", products)
		return products, nil
	}

	// get products from db
	if err := r.db.GetAll(ctx, r.collection, &products); err != nil {
		log.Info("[Repo]: products from db: ", products)
		return nil, err
	}

	// set cache products in redis
	if err := r.cacheSvc.Set(ctx, cacheKeyList, products, 15*time.Minute); err != nil {
		log.Warn("[Repo]: failed to set cache products in GetAllProduct: ", err)
	}

	return products, nil
}

func (r *productRepo) GetProductByID(ctx context.Context, id string, product *model.Product) (err error) {
	// get product from redis
	cacheKeyID := r.keyGen.KeyID(id)
	if err := r.cacheSvc.Get(ctx, cacheKeyID, product); err == nil {
		log.Info("[Repo]: product from redis: ", product)
		return nil
	}

	// get product from db if fail to get that from redis
	if err := r.db.GetByID(ctx, r.collection, id, product); err != nil {
		log.Info("[Repo]: product from db: ", product)
		return err
	}

	// set product cache in redis
	if err := r.cacheSvc.Set(ctx, cacheKeyID, product, 15*time.Minute); err != nil {
		log.Warn("[Repo]: fail to set cache product in GetProductByID: ", err)
	}

	return nil
}
