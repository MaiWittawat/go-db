package user

import (
	"context"
	"go-rebuild/internal/cache"
	dbRepo "go-rebuild/internal/db"
	"go-rebuild/internal/model"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
)

type userRepo struct {
	db         dbRepo.DB
	collection string
	cacheSvc   cache.Cache
	keyGen     *cache.KeyGenerator
}

// ------------------------ Constructor ------------------------
func NewUserRepo(db dbRepo.DB, cacheSvc cache.Cache) repository.UserRepository {
	keyGen := cache.NewKeyGenerator("users")
	return &userRepo{
		db: db, 
		collection: "users", 
		cacheSvc: cacheSvc, 
		keyGen: keyGen,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (ur *userRepo) AddUser(ctx context.Context, u *model.User) error {
	// save to db
	if err := ur.db.Create(ctx, ur.collection, u); err != nil {
		return err
	}

	// clear last cahce list
	cacheKeyList := ur.keyGen.KeyList()
	if err := ur.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache users in AddUser: ", err)
	}

	// set cache
	cacheKeyID := ur.keyGen.KeyID(u.ID)
	if err := ur.cacheSvc.Set(ctx, cacheKeyID, u, 15*time.Minute); err != nil {
		log.Warn("failed to set cache user in AddUser: ", err)
	}

	log.Info("set cache in AddUser success")
	return nil
}

func (ur *userRepo) UpdateUser(ctx context.Context, u *model.User, id string) error {
	var oldUser model.User
	if err := ur.db.GetByID(ctx, ur.collection, id, &oldUser); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": id,
			"layer":   "repository",
			"step":    "UpdateUser",
		}).Error("failed to get user by id")
		return err
	}

	// clear old cache
	cacheKeyEmail := ur.keyGen.KeyField("email", oldUser.Email)
	if err := ur.cacheSvc.Delete(ctx, cacheKeyEmail); err != nil {
		log.Warn("failed to clear cache user in UpdateUser: ", err)
	}

	// update user data in db
	if err := ur.db.Update(ctx, ur.collection, u, id); err != nil {
		return err
	}
	log.Info("user update user: ", u)

	// clear user cache
	cacheKeyList := ur.keyGen.KeyList()
	if err := ur.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache user in UpdateUser: ", err)
	}

	// set cache
	cacheKeyEmail = ur.keyGen.KeyField("email", u.Email)
	if err := ur.cacheSvc.Set(ctx, cacheKeyEmail, u, 15*time.Minute); err != nil {
		log.Warn("failed to set cache user in UpdateUser: ", err)
	}

	cacheKeyID := ur.keyGen.KeyID(id)
	if err := ur.cacheSvc.Set(ctx, cacheKeyID, u, 15*time.Minute); err != nil {
		log.Warn("failed to clear cache user in UpdateUser: ", err)
	}

	log.Info("set cache in UpdateUser success")
	return nil

}

func (ur *userRepo) DeleteUser(ctx context.Context, id string, user *model.User) error {
	if err := ur.db.Delete(ctx, ur.collection, user, id); err != nil {
		return err
	}

	cacheKeyID := ur.keyGen.KeyID(id)
	if err := ur.cacheSvc.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache user in DeleteUser: ", err)
	}

	log.Info("clear cache in DeleteUser success")
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (ur *userRepo) GetAllUser(ctx context.Context) ([]model.User, error) {
	cacheKeyList := ur.keyGen.KeyList()
	var users []model.User

	// get data from redis
	if err := ur.cacheSvc.Get(ctx, cacheKeyList, &users); err == nil {
		return users, nil
	}

	// get data from db
	if err := ur.db.GetAll(ctx, ur.collection, &users); err != nil {
		return nil, err
	}

	// set data to redis
	if err := ur.cacheSvc.Set(ctx, cacheKeyList, users, 15*time.Minute); err != nil {
		log.Warn("failed to set cache users in GetAllUser")
	}

	log.Info("set cache in GetAllUser success")
	return users, nil
}

func (ur *userRepo) GetUserByID(ctx context.Context, id string, user *model.User) (err error) {
	log.Info("user id from userRepo: ", id)
	cacheKeyID := ur.keyGen.KeyID(id)
	if err := ur.cacheSvc.Get(ctx, cacheKeyID, &user); err == nil {
		log.Info("user from cache : ", user)
		return nil
	}

	if err := ur.db.GetByID(ctx, ur.collection, id, user); err != nil {
		log.Info("user from db : ", user)
		return err
	}

	if err := ur.cacheSvc.Set(ctx, cacheKeyID, user, 15*time.Minute); err != nil {
		log.Warn("failed to set cache user in GetUserByID")
	}

	log.Info("set cache in GetUserByID success")
	return nil
}

func (ur *userRepo) GetUserByEmail(ctx context.Context, email string, user *model.User) (err error) {
	cacheKeyEmail := ur.keyGen.KeyField("email", email)
	if err := ur.cacheSvc.Get(ctx, cacheKeyEmail, &user); err == nil {
		log.Info("user from cache : ", user)
		return nil
	}

	if err := ur.db.GetByField(ctx, ur.collection, "email", email, user); err != nil {
		log.Info("user from db : ", user)
		return err
	}

	if err := ur.cacheSvc.Set(ctx, cacheKeyEmail, user, 15*time.Minute); err != nil {
		log.Warn("failed to set cache user in GetUserByEmail")
	}

	log.Info("set cache in GetUserByField success")
	return nil
}
