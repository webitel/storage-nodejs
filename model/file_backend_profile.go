package model

import (
	"encoding/json"
	"io"
	"net/http"
)

type FileBackendProfile struct {
	Id         int64           `db:"id" json:"id"`
	Name       string          `db:"name" json:"name"`
	Domain     string          `db:"domain" json:"domain"`
	Default    bool            `db:"default" json:"default"`
	ExpireDay  int             `db:"expire_day" json:"expire_day"`
	Disabled   bool            `db:"disabled" json:"disabled"`
	MaxSizeMb  int             `db:"max_size_mb" json:"max_size_mb"`
	Properties StringInterface `db:"properties" json:"properties"`
}
type FileBackendProfilePath struct {
	Name       *string          `json:"name"`
	Default    *bool            `json:"default"`
	ExpireDay  *int             `json:"expire_day"`
	Disabled   *bool            `json:"disabled"`
	MaxSizeMb  *int             `json:"max_size_mb"`
	Properties *StringInterface `json:"properties"`
}

func (f *FileBackendProfile) IsValid() *AppError  {
	if len(f.Name) == 0 {
		return NewAppError("FileBackendProfile.IsValid", "model.file_backend_profile.name.app_error", nil, "", http.StatusBadRequest)
	}
	if len(f.Domain) == 0 {
		return NewAppError("FileBackendProfile.IsValid", "model.file_backend_profile.domain.app_error", nil, "", http.StatusBadRequest)
	}
	return nil
}

func (f *FileBackendProfile) ToJson() string {
	b, _ := json.Marshal(f)
	return string(b)
}

func (f *FileBackendProfile) Path(path *FileBackendProfilePath)  {
	if path.Name != nil {
		f.Name = *path.Name
	}

	if path.ExpireDay != nil {
		f.ExpireDay = *path.ExpireDay
	}

	if path.Disabled != nil {
		f.Disabled = *path.Disabled
	}

	if path.Default != nil {
		f.Default = *path.Default
	}

	if path.MaxSizeMb != nil {
		f.MaxSizeMb = *path.MaxSizeMb
	}

	if path.Properties != nil {
		f.Properties = StringInterface{}
		for k, v := range *path.Properties {
			f.Properties[k] = v
		}
	}
}

func FileBackendProfileFromJson(data io.Reader) *FileBackendProfile  {
	var profile FileBackendProfile
	if err := json.NewDecoder(data).Decode(&profile); err == nil {
		return &profile
	} else {
		return nil
	}
}
func FileBackendProfilePathFromJson(data io.Reader) *FileBackendProfilePath  {
	var profile FileBackendProfilePath
	if err := json.NewDecoder(data).Decode(&profile); err == nil {
		return &profile
	} else {
		return nil
	}
}

func FileBackendProfileListToJson(list []*FileBackendProfile) string  {
	b, _ := json.Marshal(list)
	return string(b)
}
