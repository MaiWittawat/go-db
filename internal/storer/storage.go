package storer

import (
	"context"
	"io"
)

type Storer interface {
	GetFileUrl(ctx context.Context, fileName string, method string) (string, error)
	Upload(ctx context.Context, fileReder io.Reader, fileSize int64, orginalFileName string, contentType string) error
}
