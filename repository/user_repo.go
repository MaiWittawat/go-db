package repository

import (
	"context"
	dbRepo "go-rebuild/db"
	"go-rebuild/model"
	module "go-rebuild/module/user"
)

type UserRepo struct {
	db dbRepo.DB
	collection string
}

func NewUserRepo(db dbRepo.DB) module.UserRepository {
	return &UserRepo{db: db, collection: "users"}
}

func (ur *UserRepo) Add(ctx context.Context, u model.User) error {
	return ur.db.Create(ctx, ur.collection, u)
}

func (ur *UserRepo) UpdateUser(ctx context.Context, u model.User, id string) error {
	return ur.db.Update(ctx, ur.collection, u, id)
}

func (ur *UserRepo) DeleteUser(ctx context.Context, id string) error {
	var user model.User
	if err := ur.db.GetByID(ctx, ur.collection, id, &user); err != nil {
		return err
	}
	return ur.db.Delete(ctx, ur.collection, user, id)
}

func (ur *UserRepo) GetAllUser(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := ur.db.GetAll(ctx, ur.collection, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, id string, user *model.User) (err error) {
	return ur.db.GetByID(ctx, ur.collection, id, user)
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, email string, user *model.User) (err error) {
	return ur.db.GetByField(ctx, ur.collection, "email", email, user)
}

