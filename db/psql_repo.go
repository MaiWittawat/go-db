package db

import (
	"context"
	"go-rebuild/model"
	"strconv"

	"gorm.io/gorm"
)

type psqlRepo struct {
	db *gorm.DB
}

func NewPsqlRepo(db *gorm.DB) DB {
	db.AutoMigrate(&model.User{})
	return &psqlRepo{db: db}
}

func (p *psqlRepo) Create(ctx context.Context, _ string,  model any) error {
	return p.db.WithContext(ctx).Create(model).Error
}

func (p *psqlRepo) Update(ctx context.Context, _ string, model any, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	return p.db.WithContext(ctx).Model(model).Where("id = ?", id).Updates(model).Error
}

func (p *psqlRepo) Delete(ctx context.Context, _ string, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	return p.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}
