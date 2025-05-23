package model

type User struct {
	Username string
	Password string
	Email string
}

func (v *User) verify() bool {
	return false
}