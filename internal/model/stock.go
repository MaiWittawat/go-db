package model

import (
	"errors"
	"time"
)


var (
	ErrDebtStock = errors.New("insufficient stock")
	ErrQuantity = errors.New("quantity can't be under zero")
)


type Stock struct {
	ProductID string     `gorm:"column:product_id;primaryKey"`
	Quantity  int        `gorm:"column:quantity"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

// ------------------------ Public Method ------------------------
func (s *Stock) SetQuantity(quantity int) {
	if quantity < 0 {
		s.Quantity = 0
		return
	}
	s.Quantity = quantity
}

func (s *Stock) GetQuantity() int {
	return s.Quantity
}

func (s *Stock) DecreaseQuantity(quantity int) error {
	if s.Quantity - quantity < 0 {
		return ErrDebtStock
	}
	s.Quantity -= quantity
	return nil
}

func (s *Stock) IncreaseQuantity(quantity int) {
	s.Quantity += quantity
}

// ------------------------ Private Method ------------------------