package model

import (
	"encoding/json"
	"fmt"
)

type File struct {
	Id         int             `db:"id" json:"id"`
	Name       string          `db:"name" json:"name"`
	Uuid       string          `db:"uuid" json:"uuid"`
	ProfileId  int             `db:"profile_id" json:"profile_id"`
	Size       int64           `db:"size" json:"size"`
	Domain     string          `db:"domain" json:"domain"`
	MimeType   string          `db:"mime_type" json:"mime_type"`
	Properties StringInterface `db:"properties" json:"properties"`
	CreatedAt  int64           `db:"created_at" json:"created_at"`
	Instance   string          `db:"instance" json:"-"`
	Removed    *bool           `db:"removed" json:"-"`
}

type RemoveFile struct {
	Id        int    `db:"id"`
	FileId    int    `db:"file_id"`
	CreatedAt int64  `db:"created_at"`
	CreatedBy string `db:"created_by"`
}

func (self *RemoveFile) PreSave() {
	self.CreatedAt = GetMillis()
}

type RemoveFileJob struct {
	File
	RemoveFile
}

type MediaFile struct {
	File
	CreatedBy string `db:"created_by"`
	UpdatedBy string `db:"updated_by"`
	UpdatedAt int64  `db:"updated_at"`
}

type FileWithProfile struct {
	File
	ProfileUpdatedAt int64 `db:"profile_updated_at"`
}

func (f *File) ToJson() string {
	b, _ := json.Marshal(f)
	return string(b)
}

func FileListToJson(list []*File) string {
	b, _ := json.Marshal(list)
	return string(b)
}

func (self File) GetPropertyString(name string) string {
	return self.Properties.GetString(name)
}

func (self File) GetStoreName() string {
	return fmt.Sprintf("%s_%s", self.Uuid, self.Name)
}
