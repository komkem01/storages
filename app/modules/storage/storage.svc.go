package storage

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"storage/app/modules/entities/ent"
	entitiesinf "storage/app/modules/entities/inf"
	"storage/internal/config"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Endpoint             string `conf:"required"`
	AccessKeyId          string `conf:"required"`
	SecretAccessKey      string `conf:"required"`
	BucketName           string `conf:"required"`
	UseSSL               bool
	PublicBaseUrl        string
	PresignExpireSeconds int64
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	ent    entitiesinf.StorageEntity
}

type Service struct {
	*Options
	minio *minio.Client
}

func newService(opt *Options) *Service {
	if opt.Val.PresignExpireSeconds <= 0 {
		opt.Val.PresignExpireSeconds = 900
	}

	client, err := minio.New(opt.Val.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(opt.Val.AccessKeyId, opt.Val.SecretAccessKey, ""),
		Secure: opt.Val.UseSSL,
	})
	if err != nil {
		panic(err)
	}

	return &Service{
		Options: opt,
		minio:   client,
	}
}

func (s *Service) Upload(ctx context.Context, fileName string, contentType string, fileSize int64, reader io.Reader) (*ent.Storage, error) {
	if err := s.ensureBucket(ctx); err != nil {
		return nil, err
	}

	safeName := sanitizeFileName(fileName)
	objectKey := fmt.Sprintf("%s/%s-%s", time.Now().UTC().Format("2006/01/02"), uuid.NewString(), safeName)

	putInfo, err := s.minio.PutObject(ctx, s.Val.BucketName, objectKey, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	storage := &ent.Storage{
		ID:        uuid.New(),
		Provider:  "Railway",
		Path:      objectKey,
		URL:       s.publicURL(objectKey),
		FileSize:  putInfo.Size,
		MimeType:  contentType,
		ShortCode: randomShortCode(8),
	}

	return s.ent.CreateStorage(ctx, storage)
}

func (s *Service) Presign(ctx context.Context, storageID uuid.UUID, ttl time.Duration) (string, *ent.Storage, error) {
	if ttl <= 0 {
		ttl = time.Duration(s.Val.PresignExpireSeconds) * time.Second
	}

	storage, err := s.ent.GetStorageByID(ctx, storageID)
	if err != nil {
		return "", nil, err
	}

	url, err := s.minio.PresignedGetObject(ctx, s.Val.BucketName, storage.Path, ttl, nil)
	if err != nil {
		return "", nil, err
	}

	return url.String(), storage, nil
}

func (s *Service) ensureBucket(ctx context.Context) error {
	exists, err := s.minio.BucketExists(ctx, s.Val.BucketName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.minio.MakeBucket(ctx, s.Val.BucketName, minio.MakeBucketOptions{})
}

func (s *Service) publicURL(objectKey string) string {
	if s.Val.PublicBaseUrl != "" {
		return strings.TrimRight(s.Val.PublicBaseUrl, "/") + "/" + objectKey
	}

	scheme := "http"
	if s.Val.UseSSL {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.Val.Endpoint, s.Val.BucketName, objectKey)
}

func sanitizeFileName(name string) string {
	name = filepath.Base(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "-")
	if name == "" || name == "." || name == "/" {
		return "file"
	}
	return name
}

func randomShortCode(length int) string {
	if length <= 0 {
		length = 8
	}

	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return uuid.NewString()[:length]
	}

	code := strings.TrimRight(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf), "=")
	code = strings.ToLower(code)
	if len(code) > length {
		return code[:length]
	}
	for len(code) < length {
		code += "x"
	}
	return code
}
