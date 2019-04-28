package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	model "go-common/app/interface/main/reply/model/reply"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// RplyNotice return a reply notice from memory.
func (s *Service) RplyNotice(c context.Context, plat int8, build int64, app, buvid string) (n *model.Notice) {
	n = new(model.Notice)
	// NOTE skip crash build, wrong plat ipad and ipad HD
	if plat == model.PlatIPad && (build >= 10400 && build <= 10420) {
		return
	}
	if (app == "iphone" || app == "iphone_i") && plat == model.PlatIPad && (build >= 4270 && build <= 4350) {
		return
	}
	if nts, ok := s.notice[plat]; ok {
		for _, notice := range nts {
			if notice.CheckBuild(plat, build, app) {
				*n = *notice
				break
			}
		}
	}
	// 如果是空 返回给客户端null
	if n.ID == 0 {
		n = nil
		return
	}
	if plat != model.PlatWeb && strings.Contains(n.Link, "https://ad-bili-data.biligame.com/api/mobile") {
		n.Link = strings.Replace(n.Link, "__MID__", strconv.FormatInt(metadata.Int64(c, metadata.Mid), 10), 1)
		n.Link = strings.Replace(n.Link, "__IP__", metadata.String(c, metadata.RemoteIP), 1)
		n.Link = strings.Replace(n.Link, "__BUVID__", buvid, 1)
		n.Link = strings.Replace(n.Link, "__TS__", strconv.FormatInt(time.Now().Unix(), 10), 1)
	}
	return
}

// loadRplNotice load reply notice
func (s *Service) loadRplNotice() (err error) {
	nts, err := s.dao.Notice.ReplyNotice(context.Background())
	if err != nil {
		log.Error("s.ReplynNotice err(%v)", err)
		return
	}
	tmp := make(map[int8][]*model.Notice, len(nts))
	for _, nt := range nts {
		tmp[nt.Plat] = append(tmp[nt.Plat], nt)
	}
	//map的是引用类型（底层实现是指针），所以不必对s.notice加锁
	s.notice = tmp
	return
}
