package model

import "time"

type User struct {
	ID        string     `json:"-" gorm:"primaryKey" bson:"_id,omitempty"`
	Role      string     `json:"-" gorm:"role" bson:"role"`
	Username  string     `json:"username" gorm:"username" bson:"username"`
	Password  string     `json:"password" gorm:"password" bson:"password"`
	Email     string     `json:"email" gorm:"email;unique" bson:"email"`
	CreatedAt time.Time  `json:"-" gorm:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"-" gorm:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"-" gorm:"index", bson:"deleted_at,omitempty"`
}

func (u *User) verify() bool {
	return false
}
