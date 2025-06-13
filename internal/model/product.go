package model

import (
	"errors"
	"time"
)

type Product struct { // parse to db
	ID        string     `gorm:"column:id;primaryKey" bson:"_id,omitempty"`
	Title     string     `gorm:"column:title" bson:"title"`
	Price     int        `gorm:"column:price" bson:"price"`
	Detail    string     `gorm:"column:detail" bson:"detail"`
	CreatedBy string     `gorm:"column:created_by" bson:"created_by"`
	CreatedAt time.Time  `gorm:"column:created_at" bson:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" bson:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" bson:"deleted_at,omitempty"`
}

type ProductReq struct { // input form user
	Title     string `json:"title"`
	Price     int    `json:"price"`
	Detail    string `json:"detail"`
	Quantity  int    `json:"quantity"`
	CreatedBy string `json:"created_by"`
}

type ProductRes struct { // show output to user
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Price     int       `json:"price"`
	Detail    string    `json:"detail"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_by"`
}

func (pReq *ProductReq) ToProduct() *Product {
	product := Product{
		Title : pReq.Title,
		Price : pReq.Price,
		Detail: pReq.Detail,
		CreatedBy : pReq.CreatedBy,
	}
	return &product
}

func (p *Product) ToProductRes() *ProductRes {
	productRes := ProductRes{
		ID : p.ID,
		Title: p.Title,
		Price: p.Price,
		Detail: p.Detail,
		CreatedBy: p.CreatedBy,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	return &productRes
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

	if p.Detail != "" {
		if !p.isValidDetail() {
			return errors.New("product description is too short. It must be at least 4 characters long")
		}
	}
	return nil
}

func (p *Product) UpdateNotNilField(pReq *ProductReq) {
	if pReq.Title != "" {
		p.Title = pReq.Title
	}

	if pReq.Price != 0 {
		p.Price = pReq.Price
	}

	if pReq.Detail != "" {
		p.Detail = pReq.Detail
	}

	p.UpdatedAt = time.Now()
}

// ------------------------ Private Method ------------------------
func (p *Product) isValidTitle() bool {
	return len(p.Title) >= 2
}

func (p *Product) isValidPrice() bool {
	return p.Price > 0
}

func (p *Product) isValidDetail() bool {
	return len(p.Detail) >= 4
}
