package model

import (
	"encoding/json"
)

type FileBackendProfile struct {
	Id         int64           `db:"id"`
	Name       string          `db:"name"`
	Domain     string          `db:"domain"`
	Default    bool            `db:"default"`
	ExpireDay  int             `db:"expire_day"`
	Disabled   bool            `db:"disabled"`
	MaxSizeMb  int             `db:"max_size_mb"`
	Properties StringInterface `db:"properties"`
}

func (er *FileBackendProfile) ToJson() string {
	b, _ := json.Marshal(er)
	return string(b)
}
