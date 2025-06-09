package model

import (
	"errors"
	"time"
)

type Product struct {
	ID        string     `json:"-" gorm:"column:id;primaryKey" bson:"_id,omitempty"`
	Title     string     `json:"title" gorm:"column:title" bson:"title"`
	Price     int        `json:"price" gorm:"column:price" bson:"price"`
	Quantity  int        `json:"quantity" gorm:"quantity" bson:"quantity"`
	Detail    string     `json:"detail" gorm:"column:detail" bson:"detail"`
	CreatedBy string     `json:"_" gorm:"column:created_by" bson:"created_by"`
	CreatedAt time.Time  `json:"-" gorm:"column:created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"-" gorm:"column:updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at;index" bson:"deleted_at,omitempty"`
}

// ------------------------ Public Method ------------------------
func (p *Product) Verify() error {
	if p.Title != "" {
		if !p.isValidTitle() {
			return errors.New("product title is too short. It must be at least 2 characters long")
		}
	}
	if p.Price != 0 {
		if !p.isValidPrice() {
			return errors.New("product price is invalid. It must be greater than zero")
		}
	}

	if p.Quantity != 0 {
		if !p.isValidQuantity() {
			return errors.New("product quantity is invalid. It must be greater than zero")
		}
	}

	if p.Detail != "" {
		if !p.isValidDetail() {
			return errors.New("product description is too short. It must be at least 4 characters long")
		}
	}
	return nil
}

// ------------------------ Private Method ------------------------
func (p *Product) isValidTitle() bool {
	return len(p.Title) >= 2
}

func (p *Product) isValidPrice() bool {
	return p.Price > 0
}

func (p *Product) isValidQuantity() bool {
	return p.Quantity > 0
}

func (p *Product) isValidDetail() bool {
	return len(p.Detail) >= 4
}
