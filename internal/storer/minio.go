package storer

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

type storerService struct {
	minioClient *minio.Client
	bucketName  string
}

func InitMinio(url, user, password string) *minio.Client {
	minio_url := url
	minio_user := user
	minio_pass := password
	useSsl := false

	minioClient, err := minio.New(minio_url, &minio.Options{
		Creds:  credentials.NewStaticV4(minio_user, minio_pass, ""),
		Secure: useSsl,
	})

	if err != nil {
		log.Panic("Failed to connect Minio:", err)
	}

	return minioClient
}

func NewStorerService(minioClient *minio.Client, bucketName string) Storer {
	return &storerService{
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (s *storerService) GetFileUrl(ctx context.Context, objName string, method string) (string, error) {
	reqParam := make(url.Values)
	expire := 1 * time.Hour
	url, err := s.minioClient.Presign(ctx, method, s.bucketName, objName, expire, reqParam)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (s *storerService) Upload(ctx context.Context, fileReader io.Reader, fileSize int64, orginalFileName string, contentType string) error {
	log.Println("[storer]: content type", contentType)
	var uniqueObjName string
	listContentType := strings.Split(contentType, "/")
	prefix := listContentType[0]
	timestamp := time.Now().Format("20060102T150405")

	if prefix == "image" {
		uniqueObjName = fmt.Sprintf("image/%s_%s", timestamp, orginalFileName)
	} else {
		uniqueObjName = fmt.Sprintf("%s_%s", timestamp, orginalFileName)
	}

	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	info, err := s.minioClient.PutObject(ctx, s.bucketName, uniqueObjName, fileReader, fileSize, opts)
	if err != nil {
		return fmt.Errorf("failed to upload object '%s' to bucket '%s': %w", uniqueObjName, s.bucketName, err)
	}

	log.Println("[storer]: uploadfile info", info)
	return nil
}
