package data

import (
	"context"
	"strings"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/account"
	"go-common/app/interface/main/creative/dao/activity"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/article"
	"go-common/app/interface/main/creative/dao/bfs"
	"go-common/app/interface/main/creative/dao/data"
	"go-common/app/interface/main/creative/dao/medal"
	"go-common/app/interface/main/creative/dao/tag"
	arcMdl "go-common/app/interface/main/creative/model/archive"
	dataMdl "go-common/app/interface/main/creative/model/data"
	"go-common/app/interface/main/creative/service"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

//Service struct
type Service struct {
	c          *conf.Config
	data       *data.Dao
	acc        *account.Dao
	arc        *archive.Dao
	art        *article.Dao
	dtag       *tag.Dao
	dbfs       *bfs.Dao
	medal      *medal.Dao
	act        *activity.Dao
	p          *service.Public
	pCacheHit  *prom.Prom
	pCacheMiss *prom.Prom
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos, p *service.Public) *Service {
	s := &Service{
		c:          c,
		data:       data.New(c),
		acc:        rpcdaos.Acc,
		arc:        rpcdaos.Arc,
		art:        rpcdaos.Art,
		dtag:       tag.New(c),
		dbfs:       bfs.New(c),
		medal:      medal.New(c),
		act:        activity.New(c),
		p:          p,
		pCacheHit:  prom.CacheHit,
		pCacheMiss: prom.CacheMiss,
	}
	return s
}

// Stat get stat.
func (s *Service) Stat(c context.Context, mid int64, ip string) (r *dataMdl.Stat, err error) {
	if r, err = s.data.Stat(c, ip, mid); err != nil {
		return
	}
	pfl, err := s.acc.ProfileWithStat(c, mid)
	if err != nil {
		return
	}
	r.Fan = int64(pfl.Follower)
	return
}

// Tags get predict tag.
func (s *Service) Tags(c context.Context, mid int64, tid uint16, title, filename, desc, cover string, tagFrom int8) (res *dataMdl.Tags, err error) {
	res = &dataMdl.Tags{
		Tags: make([]string, 0),
	}
	var ctags []*dataMdl.CheckedTag
	ctags, err = s.TagsWithChecked(c, mid, tid, title, filename, desc, cover, tagFrom)
	if len(ctags) > 0 {
		for _, t := range ctags {
			res.Tags = append(res.Tags, t.Tag)
		}
	}
	return
}

// TagsWithChecked get predict tag with checked mark.
func (s *Service) TagsWithChecked(c context.Context, mid int64, tid uint16, title, filename, desc, cover string, tagFrom int8) (retTags []*dataMdl.CheckedTag, err error) {
	retTags = make([]*dataMdl.CheckedTag, 0)
	var checkedTags []*dataMdl.CheckedTag
	if checkedTags, err = s.data.TagsWithChecked(c, mid, tid, title, filename, desc, cover, tagFrom); err != nil {
		log.Error("s.data.TagsWithChecked error(%v)", err)
		err = nil
	}
	// no recommend tag, return
	if len(checkedTags) == 0 {
		return
	}
	actTagsMap := make(map[string]int64)
	acts, _ := s.act.Activities(c)
	// no acts, no filter and return
	if len(acts) == 0 {
		retTags = checkedTags
		return
	}
	// filter act tags
	for _, act := range acts {
		tagsSlice := strings.Split(act.Tags, ",")
		if len(tagsSlice) > 0 {
			for _, tag := range tagsSlice {
				actTagsMap[tag] = act.ID
			}
		}
		actTagsMap[act.Name] = act.ID
	}
	for _, checkedTag := range checkedTags {
		if _, ok := actTagsMap[checkedTag.Tag]; !ok {
			retTags = append(retTags, checkedTag)
		}
	}
	return
}

// VideoQuitPoints get video quit data.
func (s *Service) VideoQuitPoints(c context.Context, cid, mid int64, ip string) (res []int64, err error) {
	// check cid owner
	isWhite := false
	for _, m := range s.c.Whitelist.DataMids {
		if m == mid {
			isWhite = true
			break
		}
	}
	if !isWhite {
		var video *arcMdl.Video
		if video, err = s.arc.VideoByCid(c, int64(cid), ip); err != nil {
			err = ecode.AccessDenied
			return
		}
		if arc, errr := s.arc.Archive(c, video.Aid, ip); errr == nil {
			if arc.Author.Mid != mid {
				err = ecode.AccessDenied
				return
			}
		}
	}
	res, err = s.data.VideoQuitPoints(c, cid)
	return
}

// ArchiveStat get video quit data.
func (s *Service) ArchiveStat(c context.Context, aid, mid int64, ip string) (stat *dataMdl.ArchiveData, err error) {
	// check aid valid and owner
	isWhite := false
	for _, m := range s.c.Whitelist.DataMids {
		if m == mid {
			isWhite = true
			break
		}
	}
	if !isWhite {
		a, re := s.arc.Archive(c, aid, ip)
		if re != nil {
			err = re
			return
		}
		if a == nil {
			log.Warn("arc not exist")
			err = ecode.AccessDenied
			return
		}
		if a.Author.Mid != mid {
			log.Warn("owner mid(%d) and viewer mid(%d) are different ", a.Author.Mid, mid)
			err = ecode.AccessDenied
			return
		}
	}
	if stat, err = s.data.ArchiveStat(c, aid); err != nil {
		return
	}
	if stat == nil {
		return
	}
	if stat.ArchiveStat != nil && stat.ArchiveStat.Play >= 100 {
		if stat.ArchiveAreas, err = s.data.ArchiveArea(c, aid); err != nil {
			log.Error(" s.data.ArchiveArea aid(%d)|mid(%d), err(%+v)", aid, mid, err)
			return
		}
	}
	//获取播放数据
	if stat.ArchivePlay, err = s.data.UpArcPlayAnalysis(c, aid); err != nil {
		log.Error(" s.data.UpArcPlayAnalysis aid(%d)|mid(%d), err(%+v)", aid, mid, err)
		return
	}
	return
}

// Covers get covers from ai
func (s *Service) Covers(c context.Context, mid int64, fns []string, ip string) (cvs []string, err error) {
	if cvs, err = s.data.RecommendCovers(c, mid, fns); err != nil {
		log.Error("s.data.RecommendCovers err(%v)", mid, err)
		return
	}
	return
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.data.Ping(c); err != nil {
		log.Error("s.data.Dao.PingDb err(%v)", err)
	}
	return
}

// Close dao
func (s *Service) Close() {
	s.data.Close()
}

// IsWhite check white.
func (s *Service) IsWhite(mid int64) (isWhite bool) {
	for _, m := range s.c.Whitelist.DataMids {
		if m == mid {
			isWhite = true
			break
		}
	}
	return
}

// IsForbidVideoup fn.
func (s *Service) IsForbidVideoup(mid int64) (isForbid bool) {
	for _, m := range s.c.Whitelist.ForbidVideoupMids {
		if m == mid {
			isForbid = true
			break
		}
	}
	return
}
