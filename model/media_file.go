package model

import (
	"encoding/json"
	"net/http"
)

type MediaFile struct {
	BaseFile
	DomainRecord
	DomainName string `json:"-" db:"domain_name"`
}

type SearchMediaFile struct {
	ListRequest
	Ids []uint32
}

func (MediaFile) DefaultOrder() string {
	return "id"
}

func (a MediaFile) AllowFields() []string {
	return []string{"id", "name", "mime_type", "size", "domain_id", "created_at", "created_by", "updated_at", "updated_by"}
}

func (a MediaFile) DefaultFields() []string {
	return []string{"id", "name", "mime_type", "size"}
}

func (a MediaFile) EntityName() string {
	return "media_files_view"
}

func (self *MediaFile) PreSave() *AppError {
	self.CreatedAt = GetMillis()
	self.UpdatedAt = self.CreatedAt
	return nil
}

func (f *MediaFile) IsValid() *AppError {
	if len(f.Name) < 3 {
		return NewAppError("MediaFile.IsValid", "model.media.is_valid.name.app_error", nil, "name="+f.Name, http.StatusBadRequest)
	}

	if len(f.MimeType) < 3 {
		return NewAppError("MediaFile.IsValid", "model.media.is_valid.mime_type.app_error", nil, "name="+f.Name, http.StatusBadRequest)
	}

	if f.DomainId == 0 {
		return NewAppError("MediaFile.IsValid", "model.media.is_valid.domain_id.app_error", nil, "name="+f.Name, http.StatusBadRequest)
	}

	if f.Size == 0 {
		//FIXME
		//return NewAppError("MediaFile.IsValid", "model.media.is_valid.size.app_error", nil, "name="+f.Name, http.StatusBadRequest)
	}
	return nil
}

func (self MediaFile) GetStoreName() string {
	return self.Name
}

func (self *MediaFile) ToJson() string {
	b, _ := json.Marshal(self)
	return string(b)
}

func (self *MediaFile) Domain() int64 {
	return self.DomainId
}
