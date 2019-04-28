package academy

import (
	"context"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/academy"
	"go-common/app/interface/main/creative/dao/archive"
	"go-common/app/interface/main/creative/dao/article"
	"go-common/app/interface/main/creative/dao/resource"
	acaMdl "go-common/app/interface/main/creative/model/academy"
	"go-common/app/interface/main/creative/service"

	"go-common/library/log"
)

//Service struct
type Service struct {
	c                   *conf.Config
	aca                 *academy.Dao
	arc                 *archive.Dao
	art                 *article.Dao
	resource            *resource.Dao
	TagsCache           map[string][]*acaMdl.Tag
	TagMapCache         map[int64]*acaMdl.Tag
	parentChildMapCache map[int64]*acaMdl.Tag
	ResourceMapCache    map[int64]struct{}
	OccCache            []*acaMdl.Occupation
	OccMapCache         map[int64]*acaMdl.Occupation
	SkillCache          []*acaMdl.Skill
	SkillMapCache       map[int64]*acaMdl.Skill
	OfficialID          int64
	EditorChoiceID      int64
	NewbCourseID        int64
	ResourceID          int64
	//for recommend
	Seed          int64
	RecommendArcs []*acaMdl.RecArchive
	//task
	p *service.Public
	//keywords
	KWsCache []interface{}
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos, p *service.Public) *Service {
	s := &Service{
		c:              c,
		aca:            academy.New(c),
		arc:            rpcdaos.Arc,
		art:            rpcdaos.Art,
		resource:       resource.New(c),
		OccMapCache:    make(map[int64]*acaMdl.Occupation),
		SkillMapCache:  make(map[int64]*acaMdl.Skill),
		OfficialID:     c.Academy.OfficialID,
		EditorChoiceID: c.Academy.EditorChoiceID,
		NewbCourseID:   c.Academy.NewbCourseID,
		ResourceID:     c.Academy.ResourceID,
		p:              p,
	}
	s.loadTags()
	s.loadOccupations()
	s.loadSkills()
	s.loadResources()
	s.loadKeyWords()
	go s.loadProc()
	return s
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.aca.Ping(c); err != nil {
		log.Error("s.aca.Ping err(%v)", err)
	}
	return
}

// Close dao
func (s *Service) Close() {
	s.aca.Close()
}

// loadproc
func (s *Service) loadProc() {
	for {
		time.Sleep(3 * time.Minute)
		s.loadTags()
		s.loadOccupations()
		s.loadSkills()
		s.loadResources()
		s.loadKeyWords()
	}
}

//load tags
func (s *Service) loadTags() {
	tgList, tgMap, pcMap, err := s.aca.TagList(context.Background())
	if err != nil {
		log.Error("s.aca.TagList error(%v)", err)
		return
	}
	s.TagsCache = tgList
	s.TagMapCache = tgMap
	s.parentChildMapCache = pcMap
}

//load occupations
func (s *Service) loadOccupations() {
	ocs, err := s.aca.Occupations(context.Background())
	if err != nil {
		log.Error("s.aca.Occupations error(%v)", err)
		return
	}
	s.OccCache = ocs

	for _, v := range s.OccCache {
		s.OccMapCache[v.ID] = v
	}
}

//load skills
func (s *Service) loadSkills() {
	sks, err := s.aca.Skills(context.Background())
	if err != nil {
		log.Error("s.aca.Skills error(%v)", err)
		return
	}
	s.SkillCache = sks

	for _, v := range s.SkillCache {
		s.SkillMapCache[v.ID] = v
	}
}

//load skills
func (s *Service) loadResources() {
	res, err := s.resource.Resource(context.Background(), int(s.ResourceID))
	if err != nil {
		log.Error("loadResources ResourceID(%d) error(%v)", int(s.ResourceID), err)
		return
	}
	if res == nil {
		return
	}
	s.ResourceMapCache = res
}

//load keywords
func (s *Service) loadKeyWords() {
	res, err := s.aca.Keywords(context.Background())
	if err != nil {
		return
	}
	if res == nil {
		return
	}
	s.KWsCache = acaMdl.Trees(res, "ID", "ParentID", "Children")
}
