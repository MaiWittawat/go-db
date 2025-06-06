package model

import "time"

type Order struct {
	ID        string     `json:"-" gorm:"column:id;primaryKey" bson:"_id,omitempty"`
	UserID    string     `json:"user_id" gorm:"column:user_id" bson:"user_id"`
	ProductID string     `json:"product_id" gorm:"column:product_id" bson:"product_id"`
	CreatedAt time.Time  `json:"-" gorm:"column:created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"-" gorm:"column:updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"column:deleted_at;index" bson:"deleted_at,omitempty"`
}
