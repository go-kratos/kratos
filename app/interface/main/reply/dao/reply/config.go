package reply

import (
	"context"
	"encoding/json"

	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_loadConfigSQL = "SELECT id, type, oid, adminid, operator, category, config, ctime, mtime FROM reply_config WHERE type=? AND oid=? AND category = ?"
)

// ConfigDao ConfigDao
type ConfigDao struct {
	mysql *sql.DB
}

// NewConfigDao new ConfigDao and return.
func NewConfigDao(db *sql.DB) (dao *ConfigDao) {
	dao = &ConfigDao{
		mysql: db,
	}
	return
}

// LoadConfig :load a config record
func (dao *ConfigDao) LoadConfig(c context.Context, oid int64, tp, category int8) (m *reply.Config, err error) {
	m = &reply.Config{}
	row := dao.mysql.QueryRow(c, _loadConfigSQL, tp, oid, category)
	if err = row.Scan(&m.ID, &m.Type, &m.Oid, &m.AdminID, &m.Operator, &m.Category, &m.Config, &m.CTime, &m.MTime); err != nil {
		if err == sql.ErrNoRows {
			m = nil
			err = nil
			return
		}
		log.Error("row.Scan error(%v)", err)
	}
	if m.ID > 0 {
		var dat reply.Config
		if err := json.Unmarshal([]byte(m.Config), &dat); err == nil {
			m.ShowEntry = dat.ShowEntry
			m.ShowAdmin = dat.ShowAdmin
		}
	}
	return
}
