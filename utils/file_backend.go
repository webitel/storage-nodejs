package utils

import (
	"fmt"
	"github.com/webitel/storage/model"
	"io"
	"net/http"
	"regexp"
	"sync"
	"time"
)

const (
	convert = 0.000001
)

var regCompileMask = regexp.MustCompile(`\$DOMAIN|\$Y|\$M|\$D|\$H|\$m`)

type BaseFileBackend struct {
	sync.RWMutex
	syncTime  int64
	writeSize float64
}

func (b *BaseFileBackend) GetSyncTime() int64 {
	return b.syncTime
}

func (b *BaseFileBackend) GetSize() float64 {
	b.RLock()
	defer b.RUnlock()
	return b.writeSize
}

// save to megabytes
func (b *BaseFileBackend) setWriteSize(writtenBytes int64) {
	b.Lock()
	defer b.Unlock()

	b.writeSize += float64(writtenBytes) * convert
}

type File interface {
	Domain() int64
	GetSize() int64
	GetStoreName() string
	GetPropertyString(name string) string
	SetPropertyString(name, value string)
}

type FileBackend interface {
	TestConnection() *model.AppError
	Reader(file File, offset int64) (io.ReadCloser, *model.AppError)
	Remove(file File) *model.AppError
	Write(src io.Reader, file File) (int64, *model.AppError)
	GetSyncTime() int64
	GetSize() float64
	Name() string
}

func NewBackendStore(profile *model.FileBackendProfile) (FileBackend, *model.AppError) {
	switch profile.Type {
	case model.FileDriverLocal:
		return &LocalFileBackend{
			BaseFileBackend: BaseFileBackend{
				syncTime:  profile.UpdatedAt,
				writeSize: 0,
			},
			name:        profile.Name,
			directory:   profile.Properties.GetString("directory"),
			pathPattern: profile.Properties.GetString("path_pattern"),
		}, nil
	case model.FileDriverS3:
		d := &S3FileBackend{
			BaseFileBackend: BaseFileBackend{
				syncTime:  profile.UpdatedAt,
				writeSize: 0,
			},
			name:        profile.Name,
			pathPattern: profile.Properties.GetString("path_pattern"),
			bucket:      profile.Properties.GetString("bucket_name"),
			accessKey:   profile.Properties.GetString("key_id"),
			accessToken: profile.Properties.GetString("access_key"),
			endpoint:    profile.Properties.GetString("endpoint"),
			region:      profile.Properties.GetString("region"),
		}
		if err := d.TestConnection(); err != nil {
			return nil, err
		}
		return d, nil
	}

	return nil, model.NewAppError("NewFileBackend", "api.file.no_driver.app_error", nil, "",
		http.StatusInternalServerError)
}

func parseStorePattern(pattern string, domain int64) string {
	now := time.Now()
	return regCompileMask.ReplaceAllStringFunc(pattern, func(s string) string {
		switch s {
		case "$DOMAIN":
			return fmt.Sprintf("%d", domain)
		case "$Y":
			return fmt.Sprintf("%d", now.Year())
		case "$M":
			return fmt.Sprintf("%d", now.Month())
		case "$D":
			return fmt.Sprintf("%d", now.Day())
		case "$H":
			return fmt.Sprintf("%d", now.Hour())
		case "$m":
			return fmt.Sprintf("%d", now.Minute())
		}
		return s
	})
}
