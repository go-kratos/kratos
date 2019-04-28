package web

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	tagmdl "go-common/app/interface/main/tag/model"
	"go-common/app/interface/main/web-goblin/dao/web"
	webmdl "go-common/app/interface/main/web-goblin/model/web"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

var (
	_emptyMiArc = make([]*webmdl.Mi, 0)
)

const (
	_tagBlkSize = 50
	_tagArcType = 3
)

// FullShort  xiao mi  FullShort .
func (s *Service) FullShort(c context.Context, pn, ps int64, source string) (res []*webmdl.Mi, err error) {
	var (
		aids []int64
		ip   = metadata.String(c, metadata.RemoteIP)
		m    = make(map[int64]string)
	)
	if aids, err = s.aids(c, pn, ps); err != nil {
		return
	}
	if res, err = s.archiveWithTag(c, aids, ip, m, source); err != nil {
		log.Error("s.archiveWithTag  error(%v)", err)
	}
	return
}

func (s *Service) archiveWithTag(c context.Context, aids []int64, ip string, op map[int64]string, source string) (list []*webmdl.Mi, err error) {
	var (
		arcErr, tagErr error
		archives       map[int64]*api.Arc
		pages          []*api.Page
		pageInfo       map[int64][]*api.Page
		tags           map[int64][]*tagmdl.Tag
		mutex          = sync.Mutex{}
		tempTags       []string
	)
	group := new(errgroup.Group)
	group.Go(func() error {
		if archives, arcErr = s.arc.Archives3(context.Background(), &archive.ArgAids2{Aids: aids, RealIP: ip}); arcErr != nil {
			web.PromError("Archives3接口错误", "s.arc.Archives3(%d,%s) error %v", aids, ip, err)
			return arcErr
		}
		return nil
	})
	pageInfo = make(map[int64][]*api.Page)
	for _, aid := range aids {
		group.Go(func() error {
			pages = []*api.Page{}
			if pages, err = s.arc.Page3(context.Background(), &archive.ArgAid2{Aid: aid, RealIP: ip}); err != nil {
				log.Error("s.arc.Page3 error(%v)", err)
				return err
			}
			mutex.Lock()
			pageInfo[aid] = pages
			mutex.Unlock()
			return nil
		})
	}
	aidsLen := len(aids)
	tags = make(map[int64][]*tagmdl.Tag, aidsLen)
	for i := 0; i < aidsLen; i += _tagBlkSize {
		var partAids []int64
		if i+_tagBlkSize > aidsLen {
			partAids = aids[i:]
		} else {
			partAids = aids[i : i+_tagBlkSize]
		}
		group.Go(func() (err error) {
			var tmpRes map[int64][]*tagmdl.Tag
			arg := &tagmdl.ArgResTags{Oids: partAids, Type: _tagArcType, RealIP: ip}
			if tmpRes, tagErr = s.tag.ResTags(context.Background(), arg); tagErr != nil {
				web.PromError("ResTags接口错误", "s.tag.ResTag(%+v) error(%v)", arg, tagErr)
				return
			}
			mutex.Lock()
			for aid, tmpTags := range tmpRes {
				tags[aid] = tmpTags
			}
			mutex.Unlock()
			return nil
		})
	}
	if err = group.Wait(); err != nil {
		return
	}
	for _, aid := range aids {
		if arc, ok := archives[aid]; ok && arc.IsNormal() {
			miArc := new(webmdl.Mi)
			tempTags = []string{}
			miArc.FromArchive(arc, pageInfo[aid], op[aid], source)
			if tag, ok := tags[aid]; ok {
				for _, v := range tag {
					tempTags = append(tempTags, v.Name)
				}
			}
			if len(tempTags) == 0 {
				miArc.Tags = ""
			} else {
				miArc.Tags = strings.Join(tempTags, ",")
			}
			list = append(list, miArc)
		}
	}
	if len(list) == 0 {
		list = _emptyMiArc
	}
	return
}

func (s *Service) aids(c context.Context, pn, ps int64) (res []int64, err error) {
	var start, end int64
	if pn > 1 {
		start = pn*ps + 1
	} else {
		start = 1
	}
	end = start + ps
	if s.maxAid > 0 && end > s.maxAid {
		log.Warn("aids(%d,%d) maxAid(%d)", pn, ps, s.maxAid)
		err = ecode.RequestErr
		return
	}
	for i := start; i < end; i++ {
		res = append(res, i)
	}
	return
}

func (s *Service) justAID() {
	var (
		maxAid int64
		err    error
	)
	for {
		if maxAid, err = s.arc.MaxAID(context.Background()); err != nil {
			web.PromError("MaxAID接口错误", "s.arc.MaxAID error(%v)", err)
			time.Sleep(time.Second)
			continue
		}
		if maxAid > 0 {
			atomic.StoreInt64(&s.maxAid, maxAid+100)
		}
		time.Sleep(time.Minute)
	}
}
