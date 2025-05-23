package repository

import (
	"context"
	"go-rebuild/model"
	"go-rebuild/module/port"
)

type UserRepo struct {
	db port.UserDB
}

func NewUserRepo(db port.UserDB) port.UserRepository {
	return &UserRepo{db: db}
}

func (ur *UserRepo) AddUser(ctx context.Context, u model.User) error{
	if err := ur.db.Create(ctx, &u); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := ur.db.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepo) UpdateUser(ctx context.Context, u model.User, id string) error {
	if err := ur.db.Update(ctx, &u, id); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) DeleteUser(ctx context.Context, id string) error {
	if err := ur.db.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
