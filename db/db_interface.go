package db

import (
	"context"
)

type DB interface {
	Create(ctx context.Context, collection string,  m any) error
	Update(ctx context.Context, collection string, m any, id string) error
	Delete(ctx context.Context, collection string, id string) error
}