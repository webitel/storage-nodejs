package model

type Schedule struct {
	Id             int64  `db:"id" json:"id"`
	CronExpression string `db:"cron_expression" json:"cron_expression"`
	Name           string `db:"name" json:"name"`
	Description    string `db:"description" json:"description"`
	CreatedAt      int64  `db:"created_at" json:"created_at"`
	Enabled        bool   `db:"enabled" json:"enabled"`
}
