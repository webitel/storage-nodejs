package model

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	ROOT_FILE_BACKEND_DOMAIN  = "root_domain"
	ACTIVE_BACKEND_CACHE_SIZE = 1000
	LOCAL_BACKEND             = 1
	CACHE_DIR                 = "./cache"
)

type FileBackendProfileType struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Code string `db:"code" json:"code"`
}

type FileBackendProfile struct {
	Id         int64           `db:"id" json:"id"`
	Name       string          `db:"name" json:"name"`
	Domain     string          `db:"domain" json:"domain"`
	ExpireDay  int             `db:"expire_day" json:"expire_day"`
	Priority   int             `db:"priority" json:"priority"`
	Disabled   bool            `db:"disabled" json:"disabled"`
	MaxSizeMb  int             `db:"max_size_mb" json:"max_size_mb"`
	Properties StringInterface `db:"properties" json:"properties"`
	TypeId     int             `db:"type_id" json:"type_id"`
	CreatedAt  int64           `db:"created_at" json:"created_at"`
	UpdatedAt  int64           `db:"updated_at" json:"updated_at"`
}

type FileBackendProfilePath struct {
	Name       *string          `json:"name"`
	ExpireDay  *int             `json:"expire_day"`
	Priority   *int             `json:"priority"`
	Disabled   *bool            `json:"disabled"`
	MaxSizeMb  *int             `json:"max_size_mb"`
	Properties *StringInterface `json:"properties"`
	TypeId     *int             `json:"type_id"`
}

func (f *FileBackendProfile) PreSave() {
	f.CreatedAt = GetMillis()
	f.UpdatedAt = f.CreatedAt
}

func (f *FileBackendProfile) IsValid() *AppError {
	if len(f.Name) == 0 {
		return NewAppError("FileBackendProfile.IsValid", "model.file_backend_profile.name.app_error", nil, "", http.StatusBadRequest)
	}
	if len(f.Domain) == 0 {
		return NewAppError("FileBackendProfile.IsValid", "model.file_backend_profile.domain.app_error", nil, "", http.StatusBadRequest)
	}

	if f.TypeId != 1 {
		return NewAppError("FileBackendProfile.IsValid", "model.file_backend_profile.type_id.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func (f *FileBackendProfile) ToJson() string {
	b, _ := json.Marshal(f)
	return string(b)
}

func (f *FileBackendProfile) Path(path *FileBackendProfilePath) {
	if path.Name != nil {
		f.Name = *path.Name
	}

	if path.ExpireDay != nil {
		f.ExpireDay = *path.ExpireDay
	}

	if path.Priority != nil {
		f.Priority = *path.Priority
	}

	if path.Disabled != nil {
		f.Disabled = *path.Disabled
	}

	if path.MaxSizeMb != nil {
		f.MaxSizeMb = *path.MaxSizeMb
	}

	if path.TypeId != nil {
		f.TypeId = *path.TypeId
	}

	if path.Properties != nil {
		f.Properties = StringInterface{}
		for k, v := range *path.Properties {
			f.Properties[k] = v
		}
	}
}

func FileBackendProfileFromJson(data io.Reader) *FileBackendProfile {
	var profile FileBackendProfile
	if err := json.NewDecoder(data).Decode(&profile); err == nil {
		return &profile
	} else {
		return nil
	}
}
func FileBackendProfilePathFromJson(data io.Reader) *FileBackendProfilePath {
	var profile FileBackendProfilePath
	if err := json.NewDecoder(data).Decode(&profile); err == nil {
		return &profile
	} else {
		return nil
	}
}

func FileBackendProfileListToJson(list []*FileBackendProfile) string {
	b, _ := json.Marshal(list)
	return string(b)
}
