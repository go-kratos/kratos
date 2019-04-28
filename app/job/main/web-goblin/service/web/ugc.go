package web

import (
	"context"
	"time"

	"go-common/app/job/main/web-goblin/model/web"
	"go-common/library/log"
)

const (
	_opadd   = "add"
	_opdel   = "del"
	_insert  = "insert"
	_update  = "update"
	_archive = "archive"
)

// UgcIncrement ugc increment .
func (s *Service) UgcIncrement(ctx context.Context, arg *web.ArcMsg) (err error) {
	m := make(map[string]interface{})
	if arg.New.CTime != "" {
		m["ctime"] = arg.New.CTime
	}
	m["mtime"] = time.Now().Format("2006-01-02 15:04:05")
	if arg.New.PubTime != "" {
		m["ptime"] = arg.New.PubTime
	}
	m["mid"] = arg.New.Mid
	m["aid"] = arg.New.Aid
	if arg.Action == _insert {
		m["action"] = _opadd
	}
	if arg.Action == _update {
		if arg.New.State != arg.Old.State && arg.New.State < 0 {
			m["action"] = _opdel
		} else {
			m["action"] = _update
		}
	}
	if err = s.dao.UgcSearch(ctx, m); err != nil {
		log.Error("s.dao.UgcIncre error", err)
	}
	return
}
