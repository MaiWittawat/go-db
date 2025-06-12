package db

import (
	"context"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	"go-rebuild/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type psqlRepo struct {
	db *gorm.DB
}

func InitPsqlDB() (*gorm.DB, error) {
	dns := appcore_config.Config.PostgresConnString
	dialector := postgres.Open(dns)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// ------------------------ Constructor ------------------------
func NewPsqlRepo(db *gorm.DB) DB {
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Product{})
	db.AutoMigrate(&model.Stock{})
	db.AutoMigrate(&model.Order{})
	return &psqlRepo{db: db}
}

// ------------------------ Method Basic CUD ------------------------
func (p *psqlRepo) Create(ctx context.Context, _ string, model any) error {
	return p.db.WithContext(ctx).Create(model).Error
}

func (p *psqlRepo) Update(ctx context.Context, _ string, model any, id string) error {
	result := p.db.WithContext(ctx).Model(model).Where("id = ?", id).Updates(model)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (p *psqlRepo) Delete(ctx context.Context, _ string, model any, id string) error {
	return p.db.WithContext(ctx).Delete(model, "id = ?", id).Error
}

// ------------------------ Method Basic Query ------------------------
func (p *psqlRepo) GetAll(ctx context.Context, _ string, results any) error {
	res := p.db.WithContext(ctx).Find(results)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (p *psqlRepo) GetByID(ctx context.Context, _ string, id string, result any) error {
	condition := map[string]any{"id": id}
	res := p.db.WithContext(ctx).Where(condition).First(result)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (p *psqlRepo) GetByField(ctx context.Context, _ string, field string, value any, result any) error {
	condition := map[string]any{field: value}
	res := p.db.WithContext(ctx).Where(condition).First(result)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
