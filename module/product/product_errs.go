package product

import "errors"

var (
	ErrCreateProduct = errors.New("fail to create product")
	ErrUpdateProduct = errors.New("fail to update product")
	ErrDeleteProduct = errors.New("fail to delete product")

	ErrProductNotFound = errors.New("product not found")
)
