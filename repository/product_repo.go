package repository

import (
	"context"
	dbRepo "go-rebuild/db"
	"go-rebuild/model"
	module "go-rebuild/module/product"
	"go-rebuild/redis"
	"time"

	log "github.com/sirupsen/logrus"
)

type ProductRepo struct {
	db         dbRepo.DB
	collection string
	cache      redis.Cache
	keyGen     *redis.KeyGenerator
}

func NewProductRepo(db dbRepo.DB, cache redis.Cache) module.ProductRepository {
	keyGen := redis.NewKeyGenerator("products")
	return &ProductRepo{db: db, collection: "products", cache: cache, keyGen: keyGen}
}

func (pr *ProductRepo) AddProduct(ctx context.Context, p *model.Product) error {
	// save to db
	if err := pr.db.Create(ctx, pr.collection, p); err != nil {
		return err
	}

	// clear last cache list
	cacheKeyList := pr.keyGen.KeyList()
	if err := pr.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache products in AddProduct: ", err)
	}

	// set cache
	cacheKeyID := pr.keyGen.KeyID(p.ID)
	if err := pr.cache.Set(ctx, cacheKeyID, p, 15*time.Minute); err != nil {
		log.Warn("failed to set cache product in AddProduct: ", err)
	}

	log.Info("set cache after AddProduct success")
	return nil
}

func (pr *ProductRepo) UpdateProduct(ctx context.Context, p *model.Product, id string) error {
	if err :=  pr.db.Update(ctx, pr.collection, p, id); err != nil {
		return err
	}
	
	cacheKeyList := pr.keyGen.KeyList()
	if err := pr.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache products in UpdateProduct: ", err)
	}

	cacheKeyID := pr.keyGen.KeyID(id)
	if err := pr.cache.Set(ctx, cacheKeyID, p, 15*time.Minute); err != nil {
		log.Warn("failed to set product cache in UpdateProduct: ", err)
	}
	
	log.Info("set cache in UpdateProduct success")
	return nil
}

func (pr *ProductRepo) DeleteProduct(ctx context.Context, id string, product *model.Product) error {
	if err := pr.db.Delete(ctx, pr.collection, product, id); err != nil {
		return err
	}

	cacheKeyID := pr.keyGen.KeyID(id)
	if err := pr.cache.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache product in DeleteProduct: ", err)
	}

	cacheKeyList := pr.keyGen.KeyList()
	if err := pr.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache products in DeleteProduct: ", err)
	}

	log.Info("set cache in DeleteProduct success")
	return nil
}



func (pr *ProductRepo) GetAllProduct(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	cacheKeyList := pr.keyGen.KeyList()

	if err := pr.cache.Get(ctx, cacheKeyList, &products); err == nil {
		return products, nil
	}

	if err := pr.db.GetAll(ctx, pr.collection, &products); err != nil {
		return nil, err
	}

	if err := pr.cache.Set(ctx, cacheKeyList, products, 15*time.Minute); err != nil {
		log.Warn("failed to set cache products in GetAllProduct: ", err)
	}

	log.Info("set cache in GetAllProduct success")
	return products, nil
}

func (pr *ProductRepo) GetProductByID(ctx context.Context, id string, product *model.Product) (err error) {
	cacheKeyID := pr.keyGen.KeyID(id)
	if err := pr.cache.Get(ctx, cacheKeyID, product); err == nil {
		return nil
	}

	if err := pr.db.GetByID(ctx, pr.collection, id, product); err != nil {
		return err
	}

	if err := pr.cache.Set(ctx, cacheKeyID, product, 15*time.Minute); err != nil {
		log.Warn("Warning: fail to set cache product in GetProductByID: ", err)
	}

	log.Info("set cache in GetProductByID success")
	return nil
}
