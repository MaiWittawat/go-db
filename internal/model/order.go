package model

import (
	"errors"
	"time"
)

var (
	ErrNilUserID    = errors.New("user id is nil")
	ErrNilProductID = errors.New("product id is nil")
)

type Order struct {
	ID        string     `json:"-" gorm:"column:id;primaryKey" bson:"_id,omitempty"`
	UserID    string     `json:"user_id" gorm:"column:user_id" bson:"user_id"`
	ProductID string     `json:"product_id" gorm:"column:product_id" bson:"product_id"`
	Quantity  int        `json:"quantity" gorm:"quantity" bson:"quantity"`
	Status    string     `json:"-" gorm:"column:status" bson:"status"`
	Amount    int        `json:"-" gorm:"column:amount" bson:"amount"`
	CreatedAt time.Time  `json:"-" gorm:"column:created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"-" gorm:"column:updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at;index" bson:"deleted_at,omitempty"`
}

// ------------------------ Public Method ------------------------
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