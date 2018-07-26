package utils

import (
	"fmt"
	"github.com/webitel/storage/model"
	"io"
	"net/http"
	"regexp"
	"time"
)

var regCompileMask = regexp.MustCompile(`\$DOMAIN|\$Y|\$M|\$D|\$H|\$m`)

type BaseFileBackend struct {
	syncTime int64
}

func (b *BaseFileBackend) GetSyncTime() int64 {
	return b.syncTime
}

type File interface {
	GetStoreName() string
	GetPropertyString(name string) string
}

type FileBackend interface {
	TestConnection() *model.AppError
	Reader(file File, offset int64) (io.ReadCloser, *model.AppError)
	Remove(file File) *model.AppError
	WriteFile(fr io.Reader, directory, name string) (int64, *model.AppError)
	RemoveFile(directory, name string) *model.AppError
	GetStoreDirectory(domain string) string
	GetSyncTime() int64
	Name() string
}

func NewBackendStore(profile *model.FileBackendProfile) (FileBackend, *model.AppError) {
	switch profile.TypeId {
	case model.LOCAL_BACKEND:
		return &LocalFileBackend{
			BaseFileBackend: BaseFileBackend{profile.UpdatedAt},
			name:            profile.Name,
			directory:       profile.Properties.GetString("directory"),
			pathPattern:     profile.Properties.GetString("path_pattern"),
		}, nil

	}

	return nil, model.NewAppError("NewFileBackend", "api.file.no_driver.app_error", nil, "",
		http.StatusInternalServerError)
}

func parseStorePattern(pattern, domain string) string {
	now := time.Now()
	return regCompileMask.ReplaceAllStringFunc(pattern, func(s string) string {
		switch s {
		case "$DOMAIN":
			return domain
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
