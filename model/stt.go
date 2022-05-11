package model

var (
	SttCacheSize = 100
)

type SttProfile struct {
	Id        int         `json:"id" db:"id"`
	DomainId  int64       `json:"domain_id" db:"domain_id"`
	Type      string      `json:"type" db:"type"`
	UpdatedAt int64       `json:"updated_at" db:"updated_at"`
	Config    []byte      `json:"config" db:"config"`
	Name      string      `json:"name" db:"name"`
	Instance  interface{} `json:"instance" db:"-"`
}

func (s *SttProfile) GetSyncTime() int64 {
	return s.UpdatedAt
}
