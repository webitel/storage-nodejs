package utils

import (
	"github.com/webitel/storage/model"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type LocalFileBackend struct {
	BaseFileBackend
	pathPattern string
	directory   string
	name        string
}

func (self *LocalFileBackend) Name() string {
	return self.name
}

func (self *LocalFileBackend) GetStoreDirectory(domain string) string {
	return path.Join(parseStorePattern(self.pathPattern, domain))
}

func (self *LocalFileBackend) TestConnection() *model.AppError {
	return nil
}

func (self *LocalFileBackend) WriteFile(src io.Reader, directory, name string) (int64, *model.AppError) {
	directory = path.Join(self.directory, directory)
	if err := os.MkdirAll(directory, 0774); err != nil {
		return 0, model.NewAppError("WriteFile", "utils.file.locally.create_dir.app_error", nil, "directory="+directory+", err="+err.Error(), http.StatusInternalServerError)
	}
	fw, err := os.OpenFile(path.Join(directory, name), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, model.NewAppError("WriteFile", "utils.file.locally.writing.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	defer fw.Close()
	written, err := io.Copy(fw, src)
	if err != nil {
		return written, model.NewAppError("WriteFile", "utils.file.locally.writing.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	self.setWriteSize(written)
	return written, nil
}

func (self *LocalFileBackend) Remove(file File) *model.AppError {
	if err := os.Remove(path.Join(self.directory, file.GetPropertyString("directory"), file.GetStoreName())); err != nil {
		return model.NewAppError("RemoveFile", "utils.file.locally.removing.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (self *LocalFileBackend) RemoveFile(directory, name string) *model.AppError {
	if err := os.Remove(path.Join(self.directory, directory, name)); err != nil {
		return model.NewAppError("RemoveFile", "utils.file.locally.removing.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (self *LocalFileBackend) Reader(file File, offset int64) (io.ReadCloser, *model.AppError) {
	if f, err := os.Open(filepath.Join(self.directory, file.GetPropertyString("directory"), file.GetStoreName())); err != nil {
		return nil, model.NewAppError("Reader", "api.file.reader.reading_local.app_error", nil, err.Error(), http.StatusInternalServerError)
	} else {
		if offset > 0 {
			f.Seek(offset, 0)
		}
		return f, nil
	}
}
