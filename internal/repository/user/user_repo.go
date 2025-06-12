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
		db:         db,
		collection: "users",
		cacheSvc:   cacheSvc,
		keyGen:     keyGen,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (r *userRepo) AddUser(ctx context.Context, u *model.User) error {
	// save user to db
	if err := r.db.Create(ctx, r.collection, u); err != nil {
		return err
	}

	// clear last cahce list
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear users cachelist in AddUser: ", err)
	}

	// set user cache
	cacheKeyID := r.keyGen.KeyID(u.ID)
	if err := r.cacheSvc.Set(ctx, cacheKeyID, u, 15*time.Minute); err != nil {
		log.Warn("failed to set user cacheKeyID in AddUser: ", err)
	}

	return nil
}

func (r *userRepo) UpdateUser(ctx context.Context, u *model.User, id string) error {
	var oldUser model.User
	if err := r.db.GetByID(ctx, r.collection, id, &oldUser); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": id,
			"layer":   "repository",
			"step":    "UpdateUser",
		}).Error("failed to get user by id")
		return err
	}

	// update user data in db
	if err := r.db.Update(ctx, r.collection, u, id); err != nil {
		return err
	}

	// clear old cache
	cacheKeyEmail := r.keyGen.KeyField("email", oldUser.Email)
	if err := r.cacheSvc.Delete(ctx, cacheKeyEmail); err != nil {
		log.Warn("failed to clear user cacheKeyEmail in UpdateUser: ", err)
	}

	// clear user cache
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear user cachelist in UpdateUser: ", err)
	}

	// set cache
	cacheKeyEmail = r.keyGen.KeyField("email", u.Email)
	if err := r.cacheSvc.Set(ctx, cacheKeyEmail, u, 15*time.Minute); err != nil {
		log.Warn("failed to set user cacheKeyEmail in UpdateUser: ", err)
	}

	cacheKeyID := r.keyGen.KeyID(id)
	if err := r.cacheSvc.Set(ctx, cacheKeyID, u, 15*time.Minute); err != nil {
		log.Warn("failed to clear user cacheKeyID in UpdateUser: ", err)
	} 

	return nil

}

func (r *userRepo) DeleteUser(ctx context.Context, id string, user *model.User) error {
	// delete user from db
	if err := r.db.Delete(ctx, r.collection, user, id); err != nil {
		return err
	}

	// delete cachelist in redis
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clearlist cache user in DeleteUser: ", err)
	}

	// delete cacheKeyID in redis
	cacheKeyID := r.keyGen.KeyID(id)
	if err := r.cacheSvc.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear user cacheKeyID in DeleteUser: ", err)
	} 

	// delete cacheKeyEmail in redis
	cacheKeyEmail := r.keyGen.KeyField("email", user.Email)
	if err := r.cacheSvc.Delete(ctx, cacheKeyEmail); err != nil {
		log.Warn("failed to clear user cacheKeyEmail in DeleteUser: ", err)
	} 

	return nil
}

// ------------------------ Method Basic Query ------------------------
func (r *userRepo) GetAllUser(ctx context.Context) ([]model.User, error) {
	cacheKeyList := r.keyGen.KeyList()
	var users []model.User

	// get users from redis
	if err := r.cacheSvc.Get(ctx, cacheKeyList, &users); err == nil {
		log.Info("users from cache: ", users)
		return users, nil
	}

	// get users from db
	if err := r.db.GetAll(ctx, r.collection, &users); err != nil {
		log.Info("users from db: ", users)
		return nil, err
	}

	// set cache users in redis
	if err := r.cacheSvc.Set(ctx, cacheKeyList, users, 15*time.Minute); err != nil {
		log.Warn("failed to set cachelist users in GetAllUser")
	} 

	return users, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id string, user *model.User) (err error) {
	// get data from redis
	cacheKeyID := r.keyGen.KeyID(id)
	if err := r.cacheSvc.Get(ctx, cacheKeyID, &user); err == nil {
		log.Info("user from cache: ", user)
		return nil
	}

	// get data from db if redis has no cache
	if err := r.db.GetByID(ctx, r.collection, id, user); err != nil {
		log.Info("user from db: ", user)
		return err
	}

	// set user cache in redis
	if err := r.cacheSvc.Set(ctx, cacheKeyID, user, 15*time.Minute); err != nil {
		log.Warn("failed to set user cacheKeyID in GetUserByID")
	}

	return nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string, user *model.User) (err error) {
	// get user from redis
	cacheKeyEmail := r.keyGen.KeyField("email", email)
	if err := r.cacheSvc.Get(ctx, cacheKeyEmail, &user); err == nil {
		log.Info("user from cache: ", user)
		return nil
	}

	// get user from db if redis has no cache
	if err := r.db.GetByField(ctx, r.collection, "email", email, user); err != nil {
		log.Info("user from db: ", user)
		return err
	}

	// set user cache in redis
	if err := r.cacheSvc.Set(ctx, cacheKeyEmail, user, 15*time.Minute); err != nil {
		log.Warn("failed to set user cacheKeyEmail in GetUserByEmail")
	} 

	return nil
}
