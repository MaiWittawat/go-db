package repository

import (
	"context"
	"go-rebuild/model"
	dbRepo "go-rebuild/db"
	module "go-rebuild/module/user"
)

type UserRepo struct {
	db dbRepo.DB
	collection string
}

func NewUserRepo(db dbRepo.DB) module.UserRepository {
	return &UserRepo{db: db, collection: "users"}
}

func (ur *UserRepo) AddUser(ctx context.Context, u model.User) error {
	if err := ur.db.Create(ctx, ur.collection, u); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return nil, nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, u model.User, id string) error {
	if err := ur.db.Update(ctx, ur.collection, u, id); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, id string) error {
	if err := ur.db.Delete(ctx, ur.collection, id); err != nil {
		return err
	}
	return nil
}

