package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/service/main/videoup/model/archive"
	"go-common/app/service/main/videoup/model/dede"
	"go-common/app/service/main/videoup/model/prom"
	"go-common/library/ecode"
	"go-common/library/log"
)

// syncCid init cid and send sync_cid message.
func (s *Service) syncCid(aid, mid int64, nvs []*archive.Video, isUGC bool) {
	if len(nvs) == 0 {
		return
	}
	var fnm = map[string]int64{}
	var isVup = true
	for _, v := range nvs {
		fnm[v.Filename] = v.Cid
		if v.SrcType != "vupload" {
			isVup = false
		}
	}
	log.Info("aid(%d) filename() set padinfo into s.padCh", aid)
	s.padCh <- &dede.PadInfo{Aid: aid, Mid: mid, Fnm: fnm, IsUpload: isVup, IsUGC: isUGC}
}

// upDmIndex update dm_index.
//func (s *Service) upDmIndex(c context.Context, mid int64, v *archive.Video) (err error) {
//	var dmIndex string
//	log.Info("aid(%d) filename(%s) mid(%d) begin to update dm_index & dm_indexdata", v.Aid, v.Filename, mid)
//	if v.SrcType == "vupload" {
//		dmIndex = "vupload_" + strconv.FormatInt(v.Cid, 10)
//	} else if v.SrcType == "sohu" || v.SrcType == "hunan" {
//		fn := md5.Sum([]byte(v.Filename))
//		ms := hex.EncodeToString(fn[:])
//		dmIndex = v.SrcType + "_" + ms[8:16]
//	} else if v.SrcType == "qq" {
//		dmIndex = v.Filename
//	}
//
//	if _, err = s.dede.UpDmIndex(c, v.Cid, v.Aid, mid, v.Index, v.Title, v.SrcType, dmIndex); err != nil {
//		log.Error("s.dede.AddDmIndex() error(%v) | aid(%d) cid(%d) mid(%d) index(%d) title(%s) src_type(%s) dm_index(%s)", err, v.Aid, v.Cid, mid, v.Index, v.Title, v.SrcType, dmIndex)
//		return
//	}
//	if _, err = s.dede.AddDmIdxData(c, v.Cid, v.Aid, v.Filename, v.Title, v.SrcType, v.Index); err != nil {
//		log.Error("s.upDmIndex s.dede.AddDmIdxData() error(%v) | aid(%d) cid(%d) filename(%s) title(%s) src_type(%s) index(%d) ", err, v.Aid, v.Cid, v.Filename, v.Title, v.SrcType, v.Index)
//	}
//	log.Info("aid(%d) filename(%s) mid(%d) end to update dm_index & dm_indexdata", v.Aid, v.Filename, mid)
//	return
//}

// padCids update cid.
func (s *Service) padCids(c context.Context, pad *dede.PadInfo) {
	for fn, cid := range pad.Fnm {
		var (
			v   *archive.Video
			err error
		)
		if v, err = s.arc.NewVideoByFn(c, fn); err != nil {
			log.Error("cidproc s.arc.Video aid(%d) filename(%s) error(%v))", pad.Aid, fn, err)
			return
		}
		if v == nil {
			log.Error("cidproc archive(%d) filename(%s) not exists", pad.Aid, fn)
			err = ecode.NothingFound
			pad.Paded = true
			return
		}
		if err := s.upVideoCid(c, pad.Aid, fn, cid); err != nil {
			log.Error("cidproc s.upVideoCid() error(%v) | aid(%d) filename(%s) cid(%d)", err, pad.Aid, fn, cid)
			return
		}
		//if err := s.upDmIndex(c, pad.Mid, v); err != nil {
		//	log.Error("cidproc s.upDmIndex() error(%v) | aid(%d) filename(%s) cid(%d) video(%v)", err, pad.Aid, fn, cid, v)
		//	return
		//}
		if err := s.upVideoHistoryCid(c, pad.Aid, cid, fn); err != nil {
			log.Error("cidproc s.upVideoHistoryCid() error(%v) | aid(%d) filename(%s) cid(%d)", err, pad.Aid, fn, cid)
			return
		}
	}
	pad.Paded = true
}

func (s *Service) padproc() {
	// NOTE: chan
	s.wg.Add(1)
	go func() {
		var (
			c   = context.TODO()
			pad *dede.PadInfo
			ok  bool
		)
		defer s.wg.Done()
		for {
			if pad, ok = <-s.padCh; !ok {
				log.Info("padproc s.padCh proc stop")
				return
			}
			if pad.IsUpload {
				s.addMonitor("padproc-busSyncCid", strconv.FormatInt(pad.Aid, 10))
				if pad.IsUGC {
					s.busUGCSubmit(pad)
				} else {
					s.busSyncCid(pad)
				}
			}
			log.Info("aid(%d) filename() get pad from s.padCh", pad.Aid)
			s.padCids(c, pad)
			s.addMonitor("padproc-padCids", strconv.FormatInt(pad.Aid, 10))
			if !pad.Paded {
				s.dede.PushPadCache(c, pad)
				time.Sleep(100 * time.Millisecond)
				continue
			}

		}
	}()
	// NOTE: from redis list when chan error
	s.wg.Add(1)
	go func() {
		var (
			c   = context.TODO()
			pad *dede.PadInfo
			err error
		)
		defer s.wg.Done()
		for {
			if s.closed {
				return
			}
			if pad, err = s.dede.PopPadInfoCache(c); err != nil {
				log.Error("padproc s.dede.PopPadInfoCache() error(%v)", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if pad == nil {
				select {
				case <-time.After(3 * time.Minute):
					continue
				case <-s.stop:
					return
				}
			}
			s.promErr.Incr(prom.RouteDmIndexTry)
			log.Info("aid(%d) filename() get pad from redis, get cid & update dm table", pad.Aid)
			s.padCids(c, pad)
			if !pad.Paded {
				s.dede.PushPadCache(c, pad)
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}
	}()
}
