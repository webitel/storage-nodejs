package model

import (
	"encoding/json"
)

type MediaFile struct {
	BaseFile
	DomainRecord
	DomainName string `json:"-" db:"domain_name"`
}

func (self *MediaFile) PreSave() *AppError {
	self.CreatedAt = GetMillis()
	self.UpdatedAt = self.CreatedAt
	return nil
}

func (self *MediaFile) IsValid() *AppError {

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
