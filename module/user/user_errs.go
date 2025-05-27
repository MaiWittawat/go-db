package user

import "errors"

var (
	ErrCreateUser = errors.New("fail to create user")
	ErrCreateSeller = errors.New("fail to create seller")
	ErrUpdateUser = errors.New("fail to update user")
	ErrDeleteUser = errors.New("fail to delete user")

	ErrUserNotFound = errors.New("user not found")
)
