package model

type EmailConfig struct {
	Id     int64   `db:"id" json:"id"`
	Type   string  `db:"type" json:"type"`
	Domain *string `db:"domain" json:"domain"`
	From   string  `db:"from" json:"from"`
}
