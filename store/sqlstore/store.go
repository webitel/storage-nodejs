package sqlstore

import (
	_ "github.com/lib/pq"
	"github.com/webitel/storage/model"

	"github.com/go-gorp/gorp"
)

type SqlStore interface {
	GetMaster() *gorp.DbMap
	GetReplica() *gorp.DbMap
	GetAllConns() []*gorp.DbMap

	ListQuery(out interface{}, req model.ListRequest, where string, e Entity, params map[string]interface{}) error
}
