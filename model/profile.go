package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ozzo/ozzo-validation"
)

var errorSerialize = errors.New("Error serialize")

type Profile struct {
	Id             int                `json:"id"`
	Code           string             `json:"code"`
	Name           string             `json:"name"`
	TypeName       string             `json:"-"`
	Domain         string             `json:"domain,omitempty"`
	Parameters     *ProfileParameters `json:"parameters,omitempty"`
	ArchiveProfile *int               `json:"archive_profile_id,omitempty"`
	UpdatedAt      int                `json:"updated_at"`
}

func (self *Profile) GetParameter(name string) string {
	if self.Parameters == nil {
		return ""
	}

	return self.Parameters.GetString(name)
}

func (self Profile) GetUpdateAt() int {
	return self.UpdatedAt
}

// Validate validates the Artist fields.
func (m Profile) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Id, validation.Required),
	)
}

type ProfileParameters map[string]interface{}

func (p ProfileParameters) Value() (driver.Value, error) {
	str, err := json.Marshal(p)
	return string(str), err
}

func (p *ProfileParameters) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, &p)
	}
	return errorSerialize
}

func (p ProfileParameters) GetString(name string) string {
	if val, ok := p[name]; ok {
		return fmt.Sprintf("%s", val)
	}
	return ""
}
