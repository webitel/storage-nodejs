package model

import (
	"fmt"
	"path"
)

/*
type Foo struct {
    ID int `db:"id"`
    Name string `db:"name"`
}

type Bar struct {
    ID int `db:"id"`
    Foos []Foo `db:"-,fkey=foo`
}

func Bars() []Bar {
    var bars []Bar
    dbMap.Select(&bars, "select bar.id, foo.name as foo_name from bar left outer join foo on bar.foo_id = foo.id")
}
*/

type JobUploadFile struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Uuid        string `db:"uuid"`
	Domain      string `db:"domain"`
	MimeType    string `db:"mime_type"`
	Size        int    `db:"size"`
	EmailMsg    string `db:"email_msg"`
	EmailSub    string `db:"email_sub"`
	Instance    string `db:"instance"`
	FailedCount int    `db:"failed_count"`
}

func (self *JobUploadFile) GetPath(root string) string {
	return path.Join(root, fmt.Sprintf("%s_%s", self.Uuid, self.Name))
}
