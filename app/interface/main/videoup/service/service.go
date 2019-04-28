package service

import (
	"context"
	"strings"
	"time"

	"go-common/app/interface/main/videoup/conf"
	"go-common/app/interface/main/videoup/dao/account"
	"go-common/app/interface/main/videoup/dao/archive"
	"go-common/app/interface/main/videoup/dao/bfs"
	"go-common/app/interface/main/videoup/dao/creative"
	"go-common/app/interface/main/videoup/dao/dynamic"
	"go-common/app/interface/main/videoup/dao/elec"
	"go-common/app/interface/main/videoup/dao/filter"
	"go-common/app/interface/main/videoup/dao/geetest"
	"go-common/app/interface/main/videoup/dao/mission"
	"go-common/app/interface/main/videoup/dao/order"
	"go-common/app/interface/main/videoup/dao/pay"
	"go-common/app/interface/main/videoup/dao/subtitle"
	"go-common/app/interface/main/videoup/dao/tag"
	arcmdl "go-common/app/interface/main/videoup/model/archive"
	missmdl "go-common/app/interface/main/videoup/model/mission"
	"go-common/app/interface/main/videoup/model/porder"
	"go-common/library/log"
	binfoc "go-common/library/log/infoc"
	"sync"
)

// Service is archive and videos service.
type Service struct {
	c *conf.Config
	// dao
	arc      *archive.Dao
	acc      *account.Dao
	bfs      *bfs.Dao
	miss     *mission.Dao
	elec     *elec.Dao
	order    *order.Dao
	tag      *tag.Dao
	creative *creative.Dao
	filter   *filter.Dao
	gt       *geetest.Dao
	dynamic  *dynamic.Dao
	sub      *subtitle.Dao
	pay      *pay.Dao
	wg       sync.WaitGroup
	asyncCh  chan func() error
	// type cache
	typeCache                   map[int16]*arcmdl.Type
	orderUps                    map[int64]int64
	staffUps                    map[int64]int64
	missCache                   map[int]*missmdl.Mission
	missTagsCache               map[string]int
	staffTypeCache              map[int16]*arcmdl.StaffTypeConf
	staffGary                   bool
	descFmtsCache               map[int16]map[int8][]*arcmdl.DescFormat
	exemptIDCheckUps            map[int64]int64
	exemptZeroLevelAndAnswerUps map[int64]int64
	exemptUgcPayUps             map[int64]int64
	exemptHalfMinUps            map[int64]int64
	PorderCfgs                  map[int64]*porder.Config
	PorderGames                 map[int64]*porder.Game
	infoc                       *binfoc.Infoc
}

// New new a service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		arc:      archive.New(c),
		acc:      account.New(c),
		bfs:      bfs.New(c),
		miss:     mission.New(c),
		elec:     elec.New(c),
		order:    order.New(c),
		tag:      tag.New(c),
		creative: creative.New(c),
		filter:   filter.New(c),
		gt:       geetest.New(c),
		dynamic:  dynamic.New(c),
		sub:      subtitle.New(c),
		pay:      pay.New(c),
		asyncCh:  make(chan func() error, 100),
		infoc:    binfoc.New(c.AppEditorInfoc),
	}
	s.loadType()
	s.loadMission()
	s.loadOrderUps()
	s.loadPorderCfgListAsMap()
	s.loadPorderGameListAsMap()
	s.loadExemptIDCheckUps()
	s.loadDescFormatCache()
	s.loadExemptZeroLevelAndAnswer()
	s.loadExemptUgcPayUps()
	s.loadStaffUps()
	s.loadStaffConfig()
	s.loadExemptHalfMinUps()
	go s.loadproc()
	go s.asyncproc()
	return s
}

func (s *Service) asyncproc() {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		f, ok := <-s.asyncCh
		if !ok {
			return
		}
		retry(3, 1*time.Second, f)
	}
}

func retry(attempts int, sleep time.Duration, callback func() error) (err error) {
	for i := 0; ; i++ {
		err = callback()
		if err == nil {
			return
		}
		if i >= (attempts - 1) {
			break
		}
		time.Sleep(sleep)
		log.Error("asyncproc  retry retrying after error:(%+v)", err)
	}
	log.Error("asyncproc  retry after %d attempts, last error: %s", attempts, err)
	return nil
}

func (s *Service) loadType() {
	tpm, err := s.arc.TypeMapping(context.TODO())
	if err != nil {
		log.Error("s.arc.TypeMapping error(%v)", err)
		return
	}
	s.typeCache = tpm
}

// s.missTagsCache: missionName or first tag
func (s *Service) loadMission() {
	mm, err := s.miss.Missions(context.TODO())
	if err != nil {
		log.Error("s.miss.Mission error(%v)", err)
		return
	}
	s.missTagsCache = make(map[string]int)
	for _, m := range mm {
		if len(m.Tags) > 0 {
			splitedTags := strings.Split(m.Tags, ",")
			s.missTagsCache[splitedTags[0]] = m.ID
		} else {
			s.missTagsCache[m.Name] = m.ID
		}
	}
	s.missCache = mm
}

func (s *Service) loadOrderUps() {
	orderUps, err := s.order.Ups(context.TODO())
	if err != nil {
		return
	}
	s.orderUps = orderUps
}

// loadDescFormatCache
func (s *Service) loadDescFormatCache() {
	fmts, err := s.arc.DescFormat(context.TODO())
	if err != nil {
		return
	}
	newFmts := make(map[int16]map[int8][]*arcmdl.DescFormat)
	for _, d := range fmts {
		if _, okTp := newFmts[d.TypeID]; !okTp {
			newFmts[d.TypeID] = make(map[int8][]*arcmdl.DescFormat)
		}
		if _, okCp := newFmts[d.TypeID][d.Copyright]; !okCp {
			newFmts[d.TypeID][d.Copyright] = []*arcmdl.DescFormat{}
		}
		newFmts[d.TypeID][d.Copyright] = append(newFmts[d.TypeID][d.Copyright], d)
	}
	log.Info("s.loadDescFormat: newFmts(%d), fmts(%d)", len(newFmts), len(fmts))
	s.descFmtsCache = newFmts
}

func (s *Service) loadStaffUps() {
	staffUps, err := s.arc.StaffUps(context.TODO())
	if err != nil {
		return
	}
	s.staffUps = staffUps
}

func (s *Service) loadproc() {
	for {
		time.Sleep(time.Duration(s.c.Tick))
		s.loadType()
		s.loadMission()
		s.loadExemptIDCheckUps()
		s.loadExemptZeroLevelAndAnswer()
		s.loadExemptUgcPayUps()
		s.loadExemptHalfMinUps()
		s.loadDescFormatCache()
		s.loadOrderUps()
		s.loadPorderCfgListAsMap()
		s.loadPorderGameListAsMap()
		s.loadStaffUps()
		s.loadStaffConfig()
	}
}

func (s *Service) loadPorderCfgListAsMap() {
	porderCfgList, err := s.arc.PorderCfgList(context.TODO())
	if err != nil {
		log.Error("s.arc.PorderCfgList error(%v)", err)
		return
	}
	s.PorderCfgs = porderCfgList
}

func (s *Service) loadPorderGameListAsMap() {
	list, err := s.arc.GameList(context.TODO())
	if err != nil {
		log.Error("GameList error(%v)", err)
		return
	}
	if len(list) == 0 || list == nil {
		log.Error("GameList empty error(%v)", err)
		return
	}
	s.PorderGames = list
}

// Ping ping success.
func (s *Service) Ping(c context.Context) (err error) {
	return s.arc.Ping(c)
}

// Close close resource.
func (s *Service) Close() {
	s.acc.Close()
	s.arc.Close()
	s.infoc.Close()
	close(s.asyncCh)
	s.wg.Wait()
}

// loadExemptIDCheckUps
func (s *Service) loadExemptIDCheckUps() {
	ups, err := s.arc.UpSpecial(context.TODO(), 8)
	if err != nil {
		return
	}
	s.exemptIDCheckUps = ups
}

// loadExemptIDCheckUps
func (s *Service) loadExemptZeroLevelAndAnswer() {
	ups, err := s.arc.UpSpecial(context.TODO(), 12)
	if err != nil {
		return
	}
	s.exemptZeroLevelAndAnswerUps = ups
}

// loadExemptUgcPayUps
func (s *Service) loadExemptUgcPayUps() {
	ups, err := s.arc.UpSpecial(context.TODO(), 17)
	if err != nil {
		return
	}
	s.exemptUgcPayUps = ups
}

// loadStaffConfig 加载联合投稿分区配置
func (s *Service) loadStaffConfig() {
	var err error
	if s.staffGary, s.staffTypeCache, err = s.arc.StaffTypeConfig(context.TODO()); err != nil {
		s.staffGary = true //怕权限被放出，没有获取到 creative配置的时候，默认开启了灰度策略
	}
}

// loadExemptHalfMinUps
func (s *Service) loadExemptHalfMinUps() {
	ups, err := s.arc.UpSpecial(context.TODO(), 25)
	if err != nil {
		return
	}
	s.exemptHalfMinUps = ups
}
