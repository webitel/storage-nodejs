package model

import (
	"encoding/json"
	"fmt"
)

type BaseFile struct {
	Id         int64           `db:"id" json:"id"`
	Domain     string          `db:"domain" json:"domain"`
	Name       string          `db:"name" json:"name"`
	Size       int64           `db:"size" json:"size"`
	MimeType   string          `db:"mime_type" json:"mime_type"`
	Properties StringInterface `db:"properties" json:"properties"`
	Instance   string          `db:"instance" json:"-"`
}

type File struct {
	BaseFile
	Uuid      string `db:"uuid" json:"uuid"`
	ProfileId int    `db:"profile_id" json:"profile_id"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	Removed   *bool  `db:"removed" json:"-"`
}

type RemoveFile struct {
	Id        int    `db:"id"`
	FileId    int64  `db:"file_id"`
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

func (self BaseFile) GetPropertyString(name string) string {
	return self.Properties.GetString(name)
}

func (self BaseFile) SetPropertyString(name, value string) {
	self.Properties[name] = value
}

func (self BaseFile) DomainName() string {
	return self.Domain
}

func (self File) GetStoreName() string {
	return fmt.Sprintf("%s_%s", self.Uuid, self.Name)
}
