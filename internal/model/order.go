package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrNilUserID    = errors.New("user id is nil")
	ErrNilProductID = errors.New("product id is nil")
)

type Order struct {
	ID        string     `gorm:"column:id;primaryKey" bson:"_id,omitempty"`
	UserID    string     `gorm:"column:user_id" bson:"user_id"`
	ProductID string     `gorm:"column:product_id" bson:"product_id"`
	Quantity  int        `gorm:"quantity" bson:"quantity"`
	Price     int        `gorm:"column:price" bson:"price"`
	Status    string     `gorm:"column:status" bson:"status"`
	Amount    int        `gorm:"column:amount" bson:"amount"`
	CreatedAt time.Time  `gorm:"column:created_at" bson:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" bson:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" bson:"deleted_at,omitempty"`
}

type OrderReq struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderResp struct {
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
	Status    string    `json:"status"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ------------------------ Public Method ------------------------
func (oReq *OrderReq) ToOrder(userID string, product *ProductResp) *Order {
	return &Order{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    userID,
		ProductID: oReq.ProductID,
		Quantity:  oReq.Quantity,
		Price:     product.Price,
		Amount:    product.Price * oReq.Quantity,
		Status:    "PENDING",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (o *Order) ToOrderResp() *OrderResp {
	return &OrderResp{
		UserID:    o.UserID,
		ProductID: o.ProductID,
		Quantity:  o.Quantity,
		Status:    o.Status,
		Amount:    o.Amount,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}

func (o *Order) VerifyNil(order Order) error {
	if o.UserID == "" {
		return ErrNilUserID
	} else if o.ProductID == "" {
		return ErrNilProductID
	}
	return nil
}

// ------------------------ Private Method ------------------------
// status check or someting
