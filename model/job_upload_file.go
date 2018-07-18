package model

import (
	"fmt"
)

type JobUploadFile struct {
	Id        int64  `db:"id"`
	State     int    `db:"state"`
	Name      string `db:"name"`
	Uuid      string `db:"uuid"`
	Domain    string `db:"domain"`
	MimeType  string `db:"mime_type"`
	Size      int    `db:"size"`
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

func (self *JobUploadFile) GetStoreName() string {
	return fmt.Sprintf("%s_%s", self.Uuid, self.Name)
}

func (self *JobUploadFile) GetPropertyString(name string) string {
	return ""
}
