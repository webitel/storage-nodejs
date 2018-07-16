package model

import "encoding/json"

type Recoding struct {
	Id         int             `db:"id" json:"id"`
	Name       string          `db:"name" json:"name"`
	Uuid       string          `db:"uuid" json:"uuid"`
	ProfileId  int             `db:"profile_id" json:"profile_id"`
	Size       int             `db:"size" json:"size"`
	Domain     string          `db:"domain" json:"domain"`
	MimeType   string          `db:"mime_type" json:"mime_type"`
	Properties StringInterface `db:"properties" json:"properties"`
	CreatedAt  int             `db:"created_at" json:"created_at"`
}

func (f *Recoding) ToJson() string {
	b, _ := json.Marshal(f)
	return string(b)
}
