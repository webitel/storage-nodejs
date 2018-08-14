package model

import (
	"encoding/json"
)

type MediaFile struct {
	BaseFile
	CreatedBy string `db:"created_by" json:"created_by"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	UpdatedBy string `db:"updated_by" json:"updated_by"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
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
