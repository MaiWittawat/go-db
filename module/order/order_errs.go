package order

import "errors"

var (
	ErrCreateOrder = errors.New("fail to create order")
	ErrUpdateOrder = errors.New("fail to update order")
	ErrDeleteOrder = errors.New("fail to delete order")

	ErrOrderNotFound = errors.New("order not found")
)
