package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"net/http"
)

type SqlCdrStore struct {
	SqlStore
}

func NewSqlCdrStore(sqlStore SqlStore) store.CdrStoreData {
	us := &SqlCdrStore{sqlStore}
	return us
}

func (self *SqlCdrStore) GetLegADataByUuid(uuid string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var cdrData model.CdrData
		err := self.GetReplica().SelectOne(&cdrData, `SELECT event FROM cdr_a WHERE uuid = $1`, uuid)
		if err != nil {
			result.Err = model.NewAppError("SqlCdrStore.GetLegADataByUuid", "store.sql_cdr.get_leg_a_data.finding.app_error", nil,
				fmt.Sprintf("uuid=%s, %s", uuid, err.Error()), http.StatusInternalServerError)

			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			result.Data = &cdrData
		}
	})
}

func (self *SqlCdrStore) GetLegBDataByUuid(uuid string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var cdrData model.CdrData
		err := self.GetReplica().SelectOne(&cdrData, `SELECT event FROM cdr_b WHERE uuid = $1`, uuid)
		if err != nil {
			result.Err = model.NewAppError("SqlCdrStore.GetLegBDataByUuid", "store.sql_cdr.get_leg_b_data.finding.app_error", nil,
				fmt.Sprintf("uuid=%s, %s", uuid, err.Error()), http.StatusInternalServerError)

			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			result.Data = &cdrData
		}
	})
}

func (self *SqlCdrStore) GetByUuidCall(uuid string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		cdrCall := model.CdrCall{}

		err := self.GetReplica().SelectOne(&cdrCall, `select event as leg_a,
                    (select array_to_json(array_agg(event))
						from (
							select event
                         	FROM cdr_b WHERE cdr_b.parent_uuid = cdr_a.uuid
       					) t
                    ) as legs_b
                    from cdr_a where uuid = $1
                    limit 1`, uuid)
		if err != nil {
			result.Err = model.NewAppError("SqlCdrStore.GetByUuidCall", "store.sql_cdr.get_leg_call.finding.app_error", nil,
				fmt.Sprintf("uuid=%s, %s", uuid, err.Error()), http.StatusInternalServerError)

			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			result.Data = &cdrCall
		}
	})
}
