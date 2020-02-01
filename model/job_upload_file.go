package model

import (
	"fmt"
)

type JobUploadFile struct {
	Id        int64  `db:"id"`
	State     int    `db:"state"`
	Name      string `db:"name"`
	Uuid      string `db:"uuid"`
	DomainId  int64  `db:"domain_id"`
	MimeType  string `db:"mime_type"`
	Size      int64  `db:"size"`
	EmailMsg  string `db:"email_msg"`
	EmailSub  string `db:"email_sub"`
	Instance  string `db:"instance"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
	Attempts  int    `db:"attempts,default:0" json:"attempts"`
}

type JobUploadFileWithProfile struct {
	JobUploadFile
	ProfileId        int
	ProfileUpdatedAt int64
}

func (self *JobUploadFile) PreSave() {
	if self.CreatedAt == 0 {
		self.CreatedAt = GetMillis()
	}
	self.UpdatedAt = GetMillis()
}

func (f *JobUploadFile) GetSize() int64 {
	return f.Size
}

func (self *JobUploadFile) GetStoreName() string {
	return fmt.Sprintf("%s_%s", self.Uuid, self.Name)
}

//TODO
func (self *JobUploadFile) GetPropertyString(name string) string {
	return ""
}
func (self *JobUploadFile) SetPropertyString(name, value string) {

}
func (self *JobUploadFile) Domain() int64 {
	return self.DomainId
}
