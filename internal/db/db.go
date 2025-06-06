package db

import (
	"context"
)

type DB interface {
	// basic CRUD
	Create(ctx context.Context, collection string,  m any) error
	Update(ctx context.Context, collection string, m any, id string) error
	Delete(ctx context.Context, collection string, m any, id string) error


	// basic Query
	GetAll(ctx context.Context, collection string, results any) error
	GetByID(ctx context.Context, collection string, id string, result any) error
	GetByField(ctx context.Context, collection string, field string, value any, result any) error
}