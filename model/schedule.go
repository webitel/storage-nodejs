package model

import (
	"github.com/robfig/cron"
	"github.com/webitel/wlog"
	"time"
)

type Schedule struct {
	Id             int64   `db:"id" json:"id"`
	CronExpression string  `db:"cron_expression" json:"cron_expression"`
	Type           string  `db:"type" json:"type"`
	Name           string  `db:"name" json:"name"`
	Description    *string `db:"description" json:"description"`
	TimeZone       *string `db:"time_zone" json:"time_zone"`
	CreatedAt      int64   `db:"created_at" json:"created_at"`
	Enabled        *bool   `db:"enabled" json:"enabled"`
}

func (s *Schedule) NextTime(t time.Time) int64 {
	if s.TimeZone != nil {
		if loc, _ := time.LoadLocation(*s.TimeZone); loc != nil {
			t = t.In(loc)
		}
	}

	res, err := cron.Parse(s.CronExpression)
	if err != nil {
		wlog.Critical(err.Error())
		return 0
	}

	return res.Next(t).UnixNano() / int64(time.Millisecond*1000)
}
