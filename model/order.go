package model

import "time"

type Order struct {
	UserID    string     `json:"user_id" gorm:"user_id" bson:"user_id"`
	ProductID string     `json:"product_id" gorm:"product_id" bson:"product_id"`
	CreatedAt time.Time  `json:"-" gorm:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"-" gorm:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index", bson:"deleted_at,omitempty"`
}
