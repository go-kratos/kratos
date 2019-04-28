package operation

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/creative"
	operMdl "go-common/app/interface/main/creative/model/operation"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
	"time"
)

//Service struct
type Service struct {
	c           *conf.Config
	creative    *creative.Dao
	NoticeStr   string
	CreativeStr string
	// cache
	operCache           []*operMdl.Operation
	allRelOperCache     []*operMdl.Operation
	toolCache           map[string][]*operMdl.Operation
	WebRelOperCache     map[string][]*operMdl.Operation
	CreatorRelOperCache map[string][]*operMdl.Operation
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:        c,
		creative: creative.New(c),
	}
	s.loadOper()
	s.loadTool()
	s.loadRelOper()
	go s.loadproc()
	return s
}

func (s *Service) loadOper() {
	var (
		op       []*operMdl.Operation
		creative *operMdl.Operation
		notice   *operMdl.Operation
		err      error
	)
	op, err = s.creative.Operations(context.TODO(), operMdl.FullTypes())
	if err != nil {
		log.Error("s.oper.Operations error(%v)", err)
		return
	}
	log.Warn("loadOper fulltypes (%+v)", op)
	s.operCache = op
	// generate rank 最小的 creative and notice str
	for _, v := range op {
		if v.Ty == "creative" && (v.Platform == 0 || v.Platform == 1) {
			if creative == nil || (v.Rank < creative.Rank) {
				creative = v
			}
		}
		if v.Ty == "notice" && (v.Platform == 0 || v.Platform == 1) {
			if notice == nil || (v.Rank < notice.Rank) {
				notice = v
			}
		}
	}
	log.Warn("loadOper CreativeStr(%+v) | NoticeStr(%+v)", creative, notice)
	if creative != nil {
		s.CreativeStr = creative.Content
	} else {
		s.CreativeStr = ""
	}
	if notice != nil {
		s.NoticeStr = notice.Content
	} else {
		s.NoticeStr = ""
	}
}

func (s *Service) loadRelOper() {
	var (
		op  []*operMdl.Operation
		err error
	)
	op, err = s.creative.AllOperByTypeSQL(context.TODO(), []string{"'play'", "'notice'", "'collect_arc'"})
	if err != nil {
		log.Error("s.oper.AllOperByTypeSQL error(%v)", err)
		return
	}
	s.allRelOperCache = op
	tmpWebRelOperCache := make(map[string][]*operMdl.Operation)
	tmpCreatorRelOperCache := make(map[string][]*operMdl.Operation)
	for _, v := range op {
		vs := &operMdl.Operation{
			ID:       v.ID,
			Ty:       v.Ty,
			Rank:     v.Rank,
			Pic:      v.Pic,
			Link:     v.Link,
			Content:  v.Content,
			Remark:   v.Remark,
			Note:     v.Note,
			Stime:    v.Stime,
			Etime:    v.Etime,
			AppPic:   v.AppPic,
			Platform: v.Platform,
		}
		// 合并 collect_arc + play => play, for Web + Creator
		if vs.Ty == "collect_arc" {
			vs.Ty = "play"
		}
		if v.Platform == 0 { //all platform
			tmpWebRelOperCache[vs.Ty] = append(tmpWebRelOperCache[vs.Ty], vs)
			tmpCreatorRelOperCache[vs.Ty] = append(tmpCreatorRelOperCache[vs.Ty], vs)
		} else if v.Platform == 1 { //app
			tmpCreatorRelOperCache[vs.Ty] = append(tmpCreatorRelOperCache[vs.Ty], vs)
		} else if v.Platform == 2 { //web
			tmpWebRelOperCache[vs.Ty] = append(tmpWebRelOperCache[vs.Ty], vs)
		}
	}
	s.WebRelOperCache = tmpWebRelOperCache
	s.CreatorRelOperCache = tmpCreatorRelOperCache
}

func (s *Service) loadTool() {
	var (
		icon  []*operMdl.Operation
		sicon []*operMdl.Operation
		err   error
	)
	icon, err = s.creative.Tool(context.TODO(), "icon")
	if err != nil {
		log.Error("s.oper.Tool(%s) error(%v)", "icon", err)
		return
	}
	sicon, err = s.creative.Tool(context.TODO(), "submit_icon")
	if err != nil {
		log.Error("s.oper.Tool(%s) error(%v)", "submit_icon", err)
		return
	}
	var tmp = map[string][]*operMdl.Operation{}
	tmp["icon"] = icon
	tmp["submit_icon"] = sicon
	s.toolCache = tmp
}

// loadproc
func (s *Service) loadproc() {
	for {
		time.Sleep(1 * time.Minute)
		s.loadOper()
		s.loadTool()
		s.loadRelOper()
	}
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.creative.Ping(c); err != nil {
		log.Error("s.operationDao.PingDb err(%v)", err)
	}
	return
}

// Close dao
func (s *Service) Close() {
	s.creative.Close()
}
