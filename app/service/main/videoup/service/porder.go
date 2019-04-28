package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

// PorderCfgList fn
func (s *Service) PorderCfgList(c context.Context) (pcfgs []*archive.Pconfig, err error) {
	if pcfgs, err = s.arc.PorderCfgList(c); err != nil {
		log.Error("s.arc.PorderCfgList() error(%v)", err)
	}
	return
}

// PorderArcList fn
func (s *Service) PorderArcList(c context.Context, begin, end string) (porders []*archive.PorderArc, err error) {
	time.Now().Format("2006-01-02 15:04:05")
	var (
		beginTs int64
		endTs   int64
	)
	if beginTs, err = strconv.ParseInt(begin, 10, 64); err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", begin, err)
		err = ecode.RequestErr
		return
	}
	btime := time.Unix(beginTs, 0)
	endTs, err = strconv.ParseInt(end, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt(%s) error(%v)", end, err)
		err = ecode.RequestErr
		return
	}
	etime := time.Unix(endTs, 0)
	if porders, err = s.arc.PorderArcList(c, btime, etime); err != nil {
		log.Error("s.arc.PorderArcList() error(%v)|begin(%+v)|end(%+v)", err, begin, end)
	}
	return
}
