package model

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `gorm:"column:id;primaryKey" bson:"_id,omitempty"`
	Role      string     `gorm:"column:role" bson:"role"`
	Username  string     `gorm:"column:username" bson:"username"`
	Password  string     `gorm:"column:password" bson:"password"`
	Email     string     `gorm:"column:email;unique" bson:"email"`
	CreatedAt time.Time  `gorm:"column:created_at" bson:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" bson:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" bson:"deleted_at,omitempty"`
}

// ------------------------ Setter ------------------------
func (u *User) SetPassword(password string) error {
	cost := 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

// ------------------------ Public Method ------------------------
func (u *User) Verify() error {
	if u.Username != "" {
		if !u.isValidUsername() {
			return errors.New("username is too short. It must be at least 4 characters long")
		}
	}
	if u.Password != "" {
		if !u.isValidPassword() {
			return errors.New("password is too short. It must be at least 4 characters long")
		}
	}
	if u.Email != "" {
		if !u.isValidEmail() {
			return errors.New("eamil is invalid. please try again")
		}
	}
	return nil
}

// ------------------------ Private Method ------------------------
func (u *User) isValidUsername() bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{4,20}$`)
	return re.MatchString(u.Username)
}

func (u *User) isValidPassword() bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9.!@#%&*]{4,20}$`)
	return re.MatchString(u.Password)
}

func (u *User) isValidEmail() bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(u.Email)
}

func (u *User) IsAdmin() bool {
	return u.Role == "ADMIN"
}

func (u *User) IsSeller() bool {
	return u.Role == "SELLER"
}
