package dao

import (
	"context"
	"errors"

	"go-common/app/service/main/antispam/conf"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// AreaReply .
	AreaReply int = iota + 1
	// AreaIMessage .
	AreaIMessage
	// AreaLiveDM .
	AreaLiveDM
	// AreaMainSiteDM .
	AreaMainSiteDM
)

const (
	// StateDefault .
	StateDefault int = iota
	// StateDeleted .
	StateDeleted
)

var (
	// ErrPingDao .
	ErrPingDao = errors.New("Ping dao error")
	// ErrResourceNotExist .
	ErrResourceNotExist = errors.New("Resource Not Exist")
	// ErrParams .
	ErrParams = errors.New("wrong params")
)

// GetTotalCounts .
func GetTotalCounts(ctx context.Context, q Querier, selectCountsSQL string) (int64, error) {
	var totalCounts int64
	if err := q.QueryRow(ctx, selectCountsSQL).Scan(&totalCounts); err != nil {
		log.Error("Error: %v, sql: %s", err, selectCountsSQL)
		return 0, err
	}
	log.Info("GetTotalCounts query sql: %s", selectCountsSQL)
	return totalCounts, nil
}

// PingMySQL .
func PingMySQL(ctx context.Context) error {
	if db != nil {
		if err := db.Ping(ctx); err != nil {
			log.Error("%v", err)
			return err
		}
	}
	return nil
}

// Close .
func Close() {
	if db != nil {
		db.Close()
	}
}

// Init .
func Init(conf *conf.Config) (ok bool) {
	if db == nil {
		db = sql.NewMySQL(conf.MySQL.AntiSpam)
	}
	return db != nil
}

var db *sql.DB
