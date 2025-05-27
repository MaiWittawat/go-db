package repository

import (
	"context"
	dbRepo "go-rebuild/db"
	"go-rebuild/model"
	module "go-rebuild/module/user"
	"go-rebuild/redis"
	"time"

	log "github.com/sirupsen/logrus"
)

type UserRepo struct {
	db         dbRepo.DB
	collection string
	cache      redis.Cache
	keyGen     redis.KeyGenerator
}

func NewUserRepo(db dbRepo.DB, cache redis.Cache) module.UserRepository {
	keyGen := redis.NewKeyGenerator("users")
	return &UserRepo{db: db, collection: "users", cache: cache, keyGen: *keyGen}
}

func (ur *UserRepo) AddUser(ctx context.Context, u model.User) error {
	// save to db
	if err := ur.db.Create(ctx, ur.collection, u); err != nil {
		return err
	}

	// clear last cahce list
	cacheKeyList := ur.keyGen.KeyList()
	if err := ur.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache users in AddUser: ", err)
	}

	// set cache
	cacheKeyID := ur.keyGen.KeyID(u.ID)
	if err := ur.cache.Set(ctx, cacheKeyID, u, 15*time.Minute); err != nil {
		log.Warn("failed to set cache user in AddUser: ", err)
	}

	log.Info("set cache in AddUser success")
	return nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, u model.User, id string) error {
	// update user data in db
	if err := ur.db.Update(ctx, ur.collection, u, id); err != nil {
		return err
	}

	// clear user cache
	cacheKeyID := ur.keyGen.KeyID(id)
	if err := ur.cache.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache user in UpdateUser: ", err)
	}

	// set cache
	if err := ur.cache.Set(ctx, cacheKeyID, u, 15*time.Minute); err != nil {
		log.Warn("failed to set cache user in UpdateUser: ", err)
	}

	log.Info("set cache in UpdateUser success")
	return nil

}

func (ur *UserRepo) DeleteUser(ctx context.Context, id string, user *model.User) error {
	if err := ur.db.Delete(ctx, ur.collection, user, id); err != nil {
		return err
	}

	cacheKeyID := ur.keyGen.KeyID(id)
	if err := ur.cache.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache user in DeleteUser: ", err)
	}

	log.Info("clear cache in DeleteUser success")
	return nil
}

func (ur *UserRepo) GetAllUser(ctx context.Context) ([]model.User, error) {
	cacheKeyList := ur.keyGen.KeyList()
	var users []model.User

	// get data from redis
	if err := ur.cache.Get(ctx, cacheKeyList, &users); err == nil {
		return users, nil
	}

	// get data from db
	if err := ur.db.GetAll(ctx, ur.collection, &users); err != nil {
		return nil, err
	}

	// set data to redis
	if err := ur.cache.Set(ctx, cacheKeyList, users, 15*time.Minute); err != nil {
		log.Warn("failed to set cache users in GetAllUser")
	}

	log.Info("set cache in GetAllUser success")
	return users, nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, id string, user *model.User) (err error) {
	cacheKeyID := ur.keyGen.KeyID(id)
	if err := ur.cache.Get(ctx, cacheKeyID, &user); err == nil {
		return nil
	}

	if err := ur.db.GetByID(ctx, ur.collection, id, user); err != nil {
		return err
	}

	if err := ur.cache.Set(ctx, cacheKeyID, user, 15*time.Minute); err != nil {
		log.Warn("failed to set cache user in GetUserByID")
	}

	log.Info("set cache in GetUserByID success")
	return nil
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string, user *model.User) (err error) {
	cacheKeyEmail := ur.keyGen.KeyField("email", email)
	if err := ur.cache.Get(ctx, cacheKeyEmail, &user); err == nil {
		return nil
	}

	if err := ur.db.GetByField(ctx, ur.collection, "email", email, user); err != nil {
		return err
	}

	if err := ur.cache.Set(ctx, cacheKeyEmail, user, 15*time.Minute); err != nil {
		log.Warn("failed to set cache user in GetUserByEmail")
	}

	log.Info("set cache in GetUserByField success")
	return nil
}
