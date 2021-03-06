package storageproviders

import (
	"bytes"
	"context"
	"time"

	"dflimg"
)

// StorageProvider is an interface all custom defined storage providers must conform to
type StorageProvider interface {
	GenerateKey(string) string
	SupportsSignedURLs() bool
	GetSize(context.Context, *dflimg.Resource) (int, error)
	Get(context.Context, *dflimg.Resource) ([]byte, *time.Time, error)
	PrepareUpload(ctx context.Context, key, contentType string, expiry time.Duration) (string, error)
	Upload(ctx context.Context, key, contentType string, file bytes.Buffer) error
}
