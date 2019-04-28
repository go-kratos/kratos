package service

import (
	"context"
	"regexp"

	"fmt"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/database/sql"
	"go-common/library/log"
)

var (
	_filenameRegexp = regexp.MustCompile(`^[A-Z0-9a-z]+$`) // only letter digital.
)

func (s *Service) filenameCheck(filename string) bool {
	return _filenameRegexp.MatchString(filename)
}

// ObtainCid  through on filename to get cid.
func (s *Service) ObtainCid(c context.Context, fn string) (cid int64, err error) {
	if v, _ := s.arc.NewVideoFn(c, fn); v != nil {
		cid = v.Cid
		return
	}
	v := &archive.Video{
		Filename: fn,
		Status:   archive.VideoStatusUploadSubmit,
	}
	var lockSuccess bool
	if lockSuccess, err = s.busCache.Lock(c, fn, 1000); !lockSuccess || err != nil {
		if err == nil {
			err = fmt.Errorf("ObtainCid SetNXLock %s locked", fn)
			log.Error("ObtainCid Lock %s locked", v.Filename)
			return
		}
		log.Error("ObtainCid Lock had run (%v,%v) and filename is (%s)", lockSuccess, err, fn)
	}
	//占据锁或者redis异常情况下允许执行
	if cid, err = s.arc.AddNewVideo(c, v); err != nil {
		log.Error("s.arc.TxAddNewVideo(%+v) error(%v)", v, err)
	}
	return
}

//FindCidByFn    get exist cid by filename
func (s *Service) FindCidByFn(c context.Context, fn string) (cid int64, err error) {
	if v, _ := s.arc.NewVideoFn(c, fn); v != nil {
		cid = v.Cid
		return
	}
	return
}

// assignCid nvs slice struct to assign cid.
func (s *Service) assignCid(c context.Context, tx *sql.Tx, nvs []*archive.Video) (err error) {
	if len(nvs) == 0 {
		return
	}
	var (
		ok, lockSuccess bool
		cfm             map[string]int64
	)
	if cfm, err = s.arc.NewCidsByFns(c, nvs); err != nil {
		log.Error("s.arc.CidsByFns(%+v) error(%v)", err)
		return
	}
	for _, v := range nvs {
		if v.Cid, ok = cfm[v.Filename]; !ok {
			if lockSuccess, err = s.busCache.Lock(c, v.Filename, 1000); !lockSuccess || err != nil {
				if err == nil {
					err = fmt.Errorf("assignCid Lock %s locked", v.Filename)
					log.Error("assignCid Lock %s locked", v.Filename)
					return
				}
				log.Error("assignCid Lock had run (%v,%v) and filename is (%s)", lockSuccess, err, v.Filename)
			}
			//占据锁或者redis异常情况下允许执行
			if v.Cid, err = s.arc.TxAddNewVideo(tx, v); err != nil {
				log.Error("s.arc.TxAddNewVideo(%+v) error(%v)", v, err)
				s.busCache.UnLock(c, v.Filename)
				return
			}
		}
	}
	return
}
