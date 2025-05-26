package model

import "time"

type Product struct {
	ID        string     `json:"-" gorm:"primaryKey" bson:"_id,omitempty"`
	Title     string     `json:"title" gorm:"title" bson:"title"`
	Price     int        `json:"price" gorm:"price" bson:"price"`
	Detail    string     `json:"detail" gorm:"detail" bson:"detail"`
	CreatedAt time.Time  `json:"-" gorm:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"-" gorm:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index", bson:"deleted_at,omitempty"`
}
