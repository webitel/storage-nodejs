package model

const (
	SESSION_CACHE_SIZE = 35000
)

type Session struct {
	Key    string  `db:"key" json:"key"`
	Token  string  `db:"token" json:"token"`
	UserId string  `db:"user_id" json:"user_id"`
	Domain *string `db:"domain" json:"domain"`
	//RoleId int     `db:"role_id" json:"role_id"`
}

func (self *Session) IsExpired() bool {
	//TODO
	return false
}
