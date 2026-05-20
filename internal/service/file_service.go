package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileService struct {
	minioClient   *minio.Client
	bucket        string
	publicBaseURL string
	useSSL        bool
	endpoint      string
}

func NewFileService(endpoint, accessKey, secretKey, bucket string, useSSL bool, publicBase string) (*FileService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &FileService{
		minioClient:   client,
		bucket:        bucket,
		publicBaseURL: publicBase,
		useSSL:        useSSL,
		endpoint:      endpoint,
	}, nil
}

func (s *FileService) MapUploadTypeToDir(fileType string) string {
	switch fileType {
	case "avatar":
		return "avatar"
	case "background":
		return "background"
	case "image":
		return "image"
	default:
		return "misc"
	}
}

func (s *FileService) SaveUploadedFile(ctx context.Context, fh *multipart.FileHeader, fileType string) (string, error) {
	objectName, err := s.GenerateObjectName(fileType, path.Ext(fh.Filename))
	if err != nil {
		return "", err
	}

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		_ = src.Close()
	}()

	contentType := resolveContentType(fh.Header.Get("Content-Type"), fh.Filename)
	_, err = s.minioClient.PutObject(ctx, s.bucket, objectName, src, fh.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return s.BuildPublicURL(objectName), nil
}

func (s *FileService) EnsureBucket(ctx context.Context) error {
	exists, err := s.minioClient.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.minioClient.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
}

func (s *FileService) BuildPublicURL(objectName string) string {
	if s.publicBaseURL != "" {
		return strings.TrimRight(s.publicBaseURL, "/") + "/" + s.bucket + "/" + objectName
	}
	scheme := "http"
	if s.useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.endpoint, s.bucket, objectName)
}

func (s *FileService) GenerateObjectName(fileType, ext string) (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return path.Join(s.MapUploadTypeToDir(fileType), "uploads", hex.EncodeToString(randomBytes)+time.Now().Format("20060102150405")+ext), nil
}

func resolveContentType(headerContentType string, fileName string) string {
	contentType := strings.TrimSpace(headerContentType)
	// Some clients send generic octet-stream; infer from extension for better browser rendering.
	if contentType == "" || strings.EqualFold(contentType, "application/octet-stream") {
		if detected := mime.TypeByExtension(path.Ext(fileName)); detected != "" {
			contentType = detected
		}
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

