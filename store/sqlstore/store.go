package sqlstore

import (
	_ "github.com/lib/pq"

	"github.com/go-gorp/gorp"
)

type SqlStore interface {
	GetMaster() *gorp.DbMap
	GetReplica() *gorp.DbMap
	GetAllConns() []*gorp.DbMap
}
