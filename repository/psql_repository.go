package repository

import (
	"context"
	"go-rebuild/model"
	"go-rebuild/module/port"
	"strconv"

	"gorm.io/gorm"
)

type gormUser struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Username  string         `gorm:"username"`
	Email     string         `gorm:"email"`
	Password  string         `gorm:"password"`
}

// Set table name in postgresql database
func (gormUser) TableName() string {
	return "users"
}

type psqlUserRepo struct {
	db *gorm.DB
}

func NewPsqlUserRepo(db *gorm.DB) port.UserDB {
	db.AutoMigrate(&gormUser{})
	return &psqlUserRepo{db: db}
}

func (p *psqlUserRepo) Create(ctx context.Context, u *model.User) error {
	return p.db.WithContext(ctx).Create(u).Error
}

func (p *psqlUserRepo) Update(ctx context.Context, u *model.User, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	return p.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(u).Error
}

func (p *psqlUserRepo) Delete(ctx context.Context, idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	return p.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

func (p *psqlUserRepo) FindByID(ctx context.Context, idStr string) (*model.User, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = p.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return &user, err
}
