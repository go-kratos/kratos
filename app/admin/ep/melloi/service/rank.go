package service

import (
	"net/url"
	"sort"
	"time"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/log"
)

//TreesQuery Get service tree
func (s *Service) TreesQuery() (*model.TreeList, error) {
	return s.dao.TreesQuery()
}

//TreeNumQuery Get service tree num
func (s *Service) TreeNumQuery() (*model.NumList, error) {
	return s.dao.TreeNumQuery()
}

//TopHttpQuery Get top 10 url
func (s *Service) TopHttpQuery() (res *model.TopAPIRes, err error) {
	if res, err = s.dao.TopHttpQuery(); err != nil {
		log.Error("service.rank error:(%v)", err)
		return
	}
	for _, api := range res.APIList {
		u, errURL := url.Parse(api.URL)
		if errURL != nil {
			log.Error("service.rank error:(%v)", errURL)
			return
		}
		api.URL = u.Host + u.Path
	}
	return
}

//TopGrpcQuery Top Grpc Query
func (s *Service) TopGrpcQuery() (*model.GrpcRes, error) {
	return s.dao.TopGrpcQuery()
}

//TopSceneQuery Top Scene Query
func (s *Service) TopSceneQuery() (*model.SceneRes, error) {
	return s.dao.TopSceneQuery()
}

//TopDeptQuery Get top 10 department
func (s *Service) TopDeptQuery() (*model.TopDeptRes, error) {
	return s.dao.TopDeptQuery()
}

//BuildLineQuery Get test line
func (s *Service) BuildLineQuery(rank *model.Rank, summary *model.ReportSummary) (res *model.BuildLineRes, err error) {

	var timePart time.Duration
	//根据传入时间
	timeLayout := "2006-01-02 15:04:05"
	if rank.StartTime == "" && rank.EndTime == "" {
		rank.StartTime = time.Now().Add(time.Hour * -24).Format(timeLayout)
		rank.EndTime = time.Now().Format(timeLayout)
	} else if rank.StartTime == "" {
		loc, _ := time.LoadLocation("Local")                              //重要：获取时区
		theTime, _ := time.ParseInLocation(timeLayout, rank.EndTime, loc) //使用模板在对应时区转化为time.time类型
		if timePart, err = time.ParseDuration("-24h"); err != nil {
			log.Error("service.rank error:(%v)", err)
			return
		}
		rank.StartTime = theTime.Add(timePart).Format(timeLayout)
	} else if rank.EndTime == "" {
		rank.EndTime = time.Now().Format(timeLayout)
	}

	if res, err = s.dao.BuildLineQuery(rank, summary); err != nil {
		log.Error("service.rank error:(%v)", err)
		return
	}

	var myDateMap = make(map[string]int)
	for _, bu := range res.BuildList {
		if _, ok := myDateMap[bu.Date]; ok {
			myDateMap[bu.Date]++
		} else {
			myDateMap[bu.Date] = 1
		}
	}

	sortedKeys := make([]string, 0)
	for k := range myDateMap {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var mySortedMap = make(map[string]int)
	res.BuildList = []*model.Build{}
	for _, k := range sortedKeys {
		mySortedMap[k] = myDateMap[k]
		bu := model.Build{Date: k, Count: myDateMap[k]}
		res.BuildList = append(res.BuildList, &bu)
	}
	return
}

//StateLineQuery Get test state line
func (s *Service) StateLineQuery() (*model.StateLineRes, error) {
	return s.dao.StateLineQuery()
}
