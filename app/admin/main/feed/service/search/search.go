package search

import (
	"context"
	"fmt"
	"time"

	"go-common/app/admin/main/feed/conf"
	"go-common/app/admin/main/feed/dao/search"
	"go-common/app/admin/main/feed/dao/show"
	searchModel "go-common/app/admin/main/feed/model/search"
	Log "go-common/app/admin/main/feed/util"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"
)

var (
	ctx = context.TODO()
)

// Service is search service
type Service struct {
	dao      *search.Dao
	showDao  *show.Dao
	cronHot  *cron.Cron
	HotFre   string
	cronDark *cron.Cron
	DarkFre  string
}

const (
	_HotPubState        = "tianma_search_hot_state"
	_HotPubValue        = "tianma_search_hot_value"
	_HotPubSearchState  = "tianma_search_hot_search_state"
	_DarkPubState       = "tianma_search_dark_state"
	_DarkPubValue       = "tianma_search_dark_value"
	_DarkPubSearchState = "tianma_search_dark_search_state"
	_HotAutoPubState    = "tianma_search_auto_hot_state"
	_DarkAutoPubState   = "tianma_search_auto_dark_state"
	_HotShowUnpub       = 0
	_HotShowPub         = 1
	_HotShowUnUp        = 2
	_DarkShowUnpub      = 0
	_DarkShowPub        = 1
	_DarkShowUnUp       = 2
)

// New new a search service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao:      search.New(c),
		showDao:  show.New(c),
		cronHot:  cron.New(),
		cronDark: cron.New(),
		HotFre:   c.Cfg.HotCroFre,
		DarkFre:  c.Cfg.DarkCroFre,
	}
	go s.CrontLoad()
	s.cronHot.Start()
	s.cronDark.Start()
	return
}

//CrontLoad search box history
func (s *Service) CrontLoad() (err error) {
	if err = s.cronHot.AddFunc(s.HotFre, s.LoadHot); err != nil {
		log.Error("searchSrv.CrontLoaHot AddFunc LoadHot error(%v)", err)
		panic(err)
	}
	if err = s.cronDark.AddFunc(s.DarkFre, s.LoadDark); err != nil {
		log.Error("searchSrv.CrontLoaHot AddFunc LoadDark error(%v)", err)
		panic(err)
	}
	return
}

//LoadHot crontab auto load hot word
func (s *Service) LoadHot() {
	var (
		err    error
		status bool
	)
	timeTwelve := time.Now().Format("2006-01-02 ") + "12:00:00"
	timeTwelveStr, _ := s.parseTime(timeTwelve, "2006-01-02 15:04:05")
	timeZero := time.Now().Format("2006-01-02 ") + "00:00:00"
	timeZeroStr, _ := s.parseTime(timeZero, "2006-01-02 15:04:05")
	log.Info("searchSrv.LoadHot Auto LoadHot Start!")
	if time.Now().Unix() == timeZeroStr.Unix() {
		//0点会自动发布一次数据
		if err = s.SetHotPub(ctx, "crontabLoadHot", 0); err != nil {
			log.Error("searchSrv.LoadHot SetHotPub error(%v)", err)
			return
		}
		log.Info("searchSrv.LoadHot Auto LoadHot Success! 00:00 clock")
	} else if time.Now().Unix() >= timeTwelveStr.Unix() {
		log.Info("searchSrv.LoadHot Auto LoadHot Time > (%v)", timeTwelveStr)
		if status, err = s.isTodayAutoPubHot(ctx); err != nil {
			log.Error("searchSrv.LoadHot isTodayAutoPubHot error(%v)", err)
			return
		}
		log.Info("searchSrv.LoadHot Auto LoadHot Publish Status = (%v)", status)
		if status {
			return
		}
		if err = s.SetHotPub(ctx, "crontabLoadHot", 0); err != nil {
			log.Error("searchSrv.LoadHot SetHotPub error(%v)", err)
			return
		}
		log.Info("searchSrv.LoadHot Auto LoadHot Success! more than 12:00 clock")
	}
}

//LoadDark crontab auto load dark word
func (s *Service) LoadDark() {
	var (
		status bool
		err    error
	)
	log.Info("searchSrv.LoadDark Auto LoadDark Start!")
	timeTwelve := time.Now().Format("2006-01-02 ") + "12:00:00"
	timeTwelveStr, _ := s.parseTime(timeTwelve, "2006-01-02 15:04:05")
	timeZero := time.Now().Format("2006-01-02 ") + "00:00:00"
	timeZeroStr, _ := s.parseTime(timeZero, "2006-01-02 15:04:05")
	log.Info("searchSrv.LoadDark Auto LoadDark Start!")
	if time.Now().Unix() == timeZeroStr.Unix() {
		//0点会自动发布一次数据
		if err = s.SetDarkPub(ctx, "crontabLoadDark", 0); err != nil {
			log.Error("searchSrv.LoadDark SetDarkPub error(%v)", err)
			return
		}
		log.Info("searchSrv.LoadDark Auto LoadDark Success! 00:00 clock")
	} else if time.Now().Unix() >= timeTwelveStr.Unix() {
		log.Info("searchSrv.LoadDark Auto LoadDark Time > (%v)", timeTwelveStr)
		if status, err = s.isTodayAutoPubDark(ctx); err != nil {
			log.Error("searchSrv.LoadDark isTodayAutoPubDark error(%v)", err)
			return
		}
		log.Info("searchSrv.LoadDark Auto LoadDark Publish Status = (%v)", status)
		if status {
			return
		}
		if err = s.SetDarkPub(ctx, "crontabLoadDark", 0); err != nil {
			log.Error("searchSrv.LoadDark SetDarkPub error(%v)", err)
			return
		}
		log.Info("searchSrv.LoadDark Auto LoadDark Success!")
	}
}

//isTodayAutoPubHot is today publish hot word
func (s *Service) isTodayAutoPubHot(c context.Context) (status bool, err error) {
	var (
		flag bool
		date string
	)
	if flag, date, err = s.dao.GetSearchAuditStat(c, _HotAutoPubState); err != nil {
		log.Error("searchSrv.isTodayAutoPubHot GetPubState error(%v)", err)
		return
	}
	//已发布 且是今天发布的数据 则证明今天发布过
	if flag && date == time.Now().Format("2006-01-02") {
		return true, nil
	}
	return
}

//isTodayAutoPubHot is today publish hot word
func (s *Service) isTodayAutoPubDark(c context.Context) (status bool, err error) {
	var (
		flag bool
		date string
	)
	if flag, date, err = s.dao.GetSearchAuditStat(c, _DarkAutoPubState); err != nil {
		log.Error("searchSrv.isTodayAutoPubDark GetPubState error(%v)", err)
		return
	}
	//已发布 且是今天发布的数据 则证明今天发布过
	if flag && date == time.Now().Format("2006-01-02") {
		return true, nil
	}
	return
}

//parseTime parse string to unix timestamp
func (s *Service) parseTime(t string, timeLayout string) (theTime time.Time, err error) {
	//timeLayout := "2006-01-02 15:04:05"
	//timeLayout := "2006-01-02" //转化所需模板

	loc, _ := time.LoadLocation("Local") //重要：获取时区
	//使用模板在对应时区转化为time.time类型
	if theTime, err = time.ParseInLocation(timeLayout, t, loc); err != nil {
		log.Error("searchSrv.parseTime ParseInLocation(%v) error(%v)", t, err)
		return
	}
	return
}

//GetSearchValue 获取搜索的数据
func (s *Service) GetSearchValue(date string, blackSlice []string) (his []searchModel.History, err error) {
	if err = s.dao.DB.Model(&searchModel.History{}).
		Where("atime = ?", date).Where("searchword not in (?)", blackSlice).
		Where("deleted = ?", searchModel.NotDelete).
		Find(&his).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.GetSearchValue error(%v)", err)
		return
	}
	return
}

//GetSearHisValue 获取搜索热词的数据
func (s *Service) GetSearHisValue(blackSlice []string) (his []searchModel.History, err error) {
	var hisTmp searchModel.History
	if err = s.dao.DB.Model(&searchModel.History{}).
		Where("deleted = ?", searchModel.NotDelete).Order("atime desc").Limit(1).
		First(&hisTmp).Error; err != nil {
		log.Error("searchSrv.GetSearchHisValue Last Day error(%v)", err)
		return
	}
	dao := s.dao.DB.Model(&searchModel.History{}).
		Where("atime = ?", hisTmp.Atime)
	if len(blackSlice) > 0 {
		dao = dao.Where("searchword not in (?)", blackSlice)
	}
	if err = dao.Where("deleted = ?", searchModel.NotDelete).Find(&his).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.GetSearchHisValue error(%v)", err)
		return
	}
	return
}

//HotwordFromDB 从DB中取有效数据
func (s *Service) HotwordFromDB(date string) (hot []searchModel.Intervene, searchCount int, err error) {
	var (
		black      []searchModel.Black
		blackSlice []string
		his        []searchModel.History
	)
	if black, err = s.BlackList(); err != nil {
		log.Error("searchSrv.HotList Black error(%v)", err)
		return
	}
	for _, v := range black {
		blackSlice = append(blackSlice, v.Searchword)
	}
	his, err = s.GetSearchValue(date, blackSlice)
	if err != nil {
		log.Error("searchSrv.HotwordFromDB error(%v)", err)
		return
	}
	searchCount = len(his)
	//如果是取今天发布的数据 且没有取到 则以昨天的为准
	if time.Now().Format("2006-01-02") == date && len(his) == 0 {
		//如果 当天的搜索热词暂未同步过来 则取昨天的搜索热词
		if his, err = s.GetSearHisValue(blackSlice); err != nil {
			log.Error("searchSrv.HotList GetHotPubLog error(%v)", err)
			return
		}
	}
	//未结束的运营干预词
	if err = s.dao.DB.Model(&searchModel.Intervene{}).Where("etime >= ?", date).Where("searchword not in (?)", blackSlice).
		Where("deleted = ?", searchModel.NotDelete).
		Find(&hot).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.HotList Intervene error(%v)", err)
		return
	}
	//搜索历史词 默认position为-1
	for _, v := range his {
		var (
			inter = searchModel.Intervene{}
		)
		//判断搜索词是否在干预中
		if err = s.dao.DB.Model(&searchModel.Intervene{}).Where("etime >= ?", date).Where("searchword = ?", v.Searchword).
			Where("deleted = ?", searchModel.NotDelete).
			First(&inter).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Error("searchSrv.HotwordFromDB Intervene First error(%v)", err)
			return
		}
		//没有找到 则直接添加搜索的热词数据
		if err == gorm.ErrRecordNotFound {
			i := searchModel.Intervene{
				ID:         v.ID,
				Searchword: v.Searchword,
				Rank:       -1,
				Tag:        v.Tag,
				Pv:         v.Pv,
			}
			hot = append(hot, i)
		} else {
			//找到 则直接将PV值复制到运营词的数据
			for j, k := range hot {
				if k.Searchword == v.Searchword {
					hot[j].Pv = v.Pv
				}
			}
		}
	}
	return hot, searchCount, nil
}

//GetDarkValue 获取搜索热词的数据
func (s *Service) GetDarkValue(blackSlice []string) (his []searchModel.Dark, err error) {
	var darkTmp searchModel.Dark
	if err = s.dao.DB.Model(&searchModel.Dark{}).
		Where("deleted = ?", searchModel.NotDelete).Order("atime desc").Limit(1).
		First(&darkTmp).Error; err != nil {
		log.Error("searchSrv.GetDarkValue Last Day error(%v)", err)
		return
	}
	dao := s.dao.DB.Model(&searchModel.Dark{}).
		Where("atime = ?", darkTmp.Atime)
	if len(blackSlice) > 0 {
		dao = dao.Where("searchword not in (?)", blackSlice)
	}
	if err = dao.Where("deleted = ?", searchModel.NotDelete).Find(&his).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.GetDarkValue error(%v)", err)
		return
	}
	return
}

//DarkwordFromDB 从DB中取有效数据
func (s *Service) DarkwordFromDB(date string) (darkValue []searchModel.Dark, searchCount int, err error) {
	var (
		black      []searchModel.Black
		blackSlice []string
		dark       []searchModel.Dark
	)
	if black, err = s.BlackList(); err != nil {
		log.Error("searchSrv.DarkwordFromDB BlackList error(%v)", err)
	}
	for _, v := range black {
		blackSlice = append(blackSlice, v.Searchword)
	}
	if err = s.dao.DB.Model(&searchModel.Dark{}).Where("deleted = ?", searchModel.NotDelete).
		Where("atime = ?", date).Where("deleted = ?", searchModel.NotDelete).Find(&dark).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.DarkwordFromDB Find error(%v)", err)
		return
	}
	searchCount = len(dark)
	if err = s.dao.DB.Model(&searchModel.Dark{}).Where("deleted = ?", searchModel.NotDelete).
		Where("atime = ?", date).Where("searchword not in (?)", blackSlice).
		Where("deleted = ?", searchModel.NotDelete).Find(&dark).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.DarkwordFromDB Find error(%v)", err)
		return
	}
	//若搜索没有推黑马词过来 则以昨天的数据为准
	if time.Now().Format("2006-01-02") == date && len(dark) == 0 {
		if dark, err = s.GetDarkValue(blackSlice); err != nil {
			log.Error("searchSrv.DarkwordFromDB GetDarkPubLog error(%v)", err)
			return
		}
	}
	m := make(map[string]bool)
	for _, val := range dark {
		if _, ok := m[val.Searchword]; !ok {
			m[val.Searchword] = true
			darkValue = append(darkValue, val)
		}
	}
	return
}

//OpenHotList open hotword list
func (s *Service) OpenHotList(c *bm.Context) (hotout []searchModel.Intervene, err error) {
	var (
		hot []searchModel.Intervene
	)
	if hot, err = s.GetHotPub(c); err != nil {
		log.Error("searchSrv.OpenHotList GetHotPub error(%v)", err)
		return
	}
	cTime := time.Now().Unix()
	inter := map[string]bool{}
	for _, v := range hot {
		if v.Rank != -1 && cTime >= v.Stime.Time().Unix() && cTime <= v.Etime.Time().Unix() {
			//运营词
			inter[v.Searchword] = true
		}
	}
	for _, v := range hot {
		if v.Rank == -1 {
			//如果运营词已存在 则以运营词为准
			if _, flag := inter[v.Searchword]; flag {
				continue
			}
			//-1 是ai的数据 直接添加
			hotout = append(hotout, v)
		} else if cTime >= v.Stime.Time().Unix() && cTime <= v.Etime.Time().Unix() {
			hotout = append(hotout, v)
		}
	}
	return
}

//HotList hotword list
func (s *Service) HotList(c *bm.Context, t string) (hotout searchModel.HotwordOut, err error) {
	var (
		dateStamp  time.Time
		todayStamp time.Time
		flag       bool
		hot        []searchModel.Intervene
		date       string
	)
	if dateStamp, err = s.parseTime(t, "2006-01-02"); err != nil {
		return
	}
	today := time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
	if todayStamp, err = s.parseTime(today, "2006-01-02"); err != nil {
		return
	}
	if flag, date, err = s.dao.GetSearchAuditStat(c, _HotPubState); err != nil {
		log.Error("searchSrv.HotList GetPublishCache error(%v)", err)
		return
	}
	//过去的时间 则直接从日志中取数据
	if dateStamp.Unix() < todayStamp.Unix() {
		var logFlag bool
		if hotout.Hotword, logFlag, err = s.GetHotPubLog(t); err != nil {
			log.Error("searchSrv.HotList GetHotPubLog error(%v)", err)
			return
		}
		if logFlag {
			hotout.State = _HotShowPub
		} else {
			hotout.State = _HotShowUnpub
		}
		return
	}
	//公共的逻辑
	if hot, _, err = s.HotwordFromDB(t); err != nil {
		log.Error("searchSrv.HotList HotwordFromDB error(%v)", err)
		return
	}
	hotout.Hotword = hot
	//今天的数据
	if dateStamp.Unix() == todayStamp.Unix() {
		//已发布 且是今天发布的数据
		if flag && date == time.Now().Format("2006-01-02") {
			//2.判断发布的时候 是否有搜索数据过来
			var pubStatus bool
			if pubStatus, _, err = s.dao.GetSearchAuditStat(c, _HotPubSearchState); err != nil {
				log.Error("searchSrv.SetHotPub SetSearchPubStat error(%v)", err)
				return
			}
			if pubStatus {
				//发布的时候 有搜索的数据 提示上线
				hotout.State = _HotShowPub
			} else {
				//发布的是 没有搜索的数据 提示 未更新
				hotout.State = _HotShowUnUp
			}
			return
		}
		//未发布
		hotout.State = _HotShowUnpub
		return
	}
	//未来的数据 都是未发布的
	hotout.State = _HotShowUnpub
	return
}

//DarkList darkword list
func (s *Service) DarkList(c *bm.Context, t string) (darkout searchModel.DarkwordOut, err error) {
	var (
		dateStamp  time.Time
		todayStamp time.Time
		flag       bool
		//flagAuto   bool
		dark []searchModel.Dark
		date string
	)
	if dateStamp, err = s.parseTime(t, "2006-01-02"); err != nil {
		return
	}
	today := time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
	if todayStamp, err = s.parseTime(today, "2006-01-02"); err != nil {
		return
	}
	if flag, date, err = s.dao.GetSearchAuditStat(c, _DarkPubState); err != nil {
		log.Error("searchSrv.DarkList GetPublishCache error(%v)", err)
		return
	}
	//过去的时间 则直接从日志中取数据
	if dateStamp.Unix() < todayStamp.Unix() {
		var logFlag bool
		if darkout.Darkword, logFlag, err = s.GetDarkPubLog(t); err != nil {
			log.Error("searchSrv.HotList GetHotPubLog error(%v)", err)
			return
		}
		if logFlag {
			darkout.State = _DarkShowPub
		} else {
			darkout.State = _DarkShowUnpub
		}
		return
	}
	//公共的逻辑
	if dark, _, err = s.DarkwordFromDB(t); err != nil {
		log.Error("searchSrv.DarkList HotwordFromDB error(%v)", err)
		return
	}
	darkout.Darkword = dark
	//今天的数据
	if dateStamp.Unix() == todayStamp.Unix() {
		//已发布 且是今天发布的数据 则直接取缓存数据
		if flag && date == time.Now().Format("2006-01-02") {
			//判断发布的时候 是否有搜索数据过来
			var pubStatus bool
			if pubStatus, _, err = s.dao.GetSearchAuditStat(c, _DarkPubSearchState); err != nil {
				log.Error("searchSrv.DarkList GetPubState error(%v)", err)
				return
			}
			if pubStatus {
				//发布的时候 有搜索的数据 提示上线
				darkout.State = _DarkShowPub
			} else {
				//发布的是 没有搜索的数据 提示 未更新
				darkout.State = _DarkShowUnUp
			}
			return
		}
		//未更新
		darkout.State = _DarkShowUnpub
		return
	}
	//未来的数据 都是未发布的
	darkout.State = _HotShowUnpub
	return
}

//BlackList black list
func (s *Service) BlackList() (black []searchModel.Black, err error) {
	if err = s.dao.DB.Model(&searchModel.Black{}).
		Where("deleted = ?", searchModel.NotDelete).Find(&black).Error; err != nil {
		log.Error("searchSrv.History Index error(%v)", err)
		return
	}
	return
}

//DelBlack add black
func (s *Service) DelBlack(c *bm.Context, id int, person string, uid int64) (err error) {
	var (
		black searchModel.Black
	)
	//根据id查找热词
	if err = s.dao.DB.Model(&searchModel.Black{}).
		Where("id = ?", id).First(&black).Error; err != nil {
		log.Error("searchSrv.DelBlack Black First error(%v)", err)
		return
	}
	if err = s.dao.DB.Model(&searchModel.Black{}).
		Where("id = ?", id).Update("deleted", searchModel.Delete).Error; err != nil {
		log.Error("searchSrv.DelBlack Update error(%v)", err)
		return
	}
	//更新AI热词为删除状态
	if err = s.dao.DB.Model(&searchModel.History{}).
		Where("searchword = ?", black.Searchword).Update("deleted", searchModel.Delete).Error; err != nil {
		log.Error("searchSrv.DelBlack Update History error(%v)", err)
		return
	}
	//更新运营热词为删除状态
	if err = s.dao.DB.Model(&searchModel.Intervene{}).
		Where("searchword = ?", black.Searchword).Update("deleted", searchModel.Delete).Error; err != nil {
		log.Error("searchSrv.DelBlack Update error(%v)", err)
		return
	}
	//更新黑马词为删除状态
	if err = s.dao.DB.Model(&searchModel.Dark{}).
		Where("searchword = ?", black.Searchword).Update("deleted", searchModel.Delete).Error; err != nil {
		log.Error("searchSrv.DelBlack Update error(%v)", err)
		return
	}
	//设置黑名单之后 立即发布新数据
	if err = s.SetHotPub(c, person, uid); err != nil {
		return
	}
	if err = s.SetDarkPub(c, person, uid); err != nil {
		return
	}
	obj := map[string]interface{}{
		"id": id,
	}
	if err = Log.AddLog(searchModel.Business, person, uid, int64(id), searchModel.ActionDelBlack, obj); err != nil {
		log.Error("searchSrv.DelBlack AddLog error(%v)", err)
		return
	}
	return
}

//AddBlack add black
func (s *Service) AddBlack(c *bm.Context, black string, person string, uid int64) (err error) {
	var (
		word searchModel.Black
	)
	if err = s.dao.DB.Model(&searchModel.Black{}).
		Where("deleted = ?", searchModel.NotDelete).Where("searchword = ?", black).
		First(&word).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.AddBlack get First error(%v)", err)
		return
	}
	if err != gorm.ErrRecordNotFound {
		err = fmt.Errorf("黑名单已存在")
		return
	}
	w := searchModel.AddBlack{
		Searchword: black,
	}
	if err = s.dao.DB.Model(&searchModel.Black{}).
		Create(w).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.AddBlack Create error(%v)", err)
		return
	}
	//设置黑名单之后 立即发布新数据
	if err = s.SetHotPub(c, person, uid); err != nil {
		return
	}
	if err = s.SetDarkPub(c, person, uid); err != nil {
		return
	}
	obj := map[string]interface{}{
		"blackword": word,
	}
	if err = Log.AddLog(searchModel.Business, person, uid, 0, searchModel.ActionAddBlack, obj); err != nil {
		log.Error("searchSrv.AddBlack AddLog error(%v)", err)
		return
	}
	return
}

//checkBlack checkout blacklist
func (s *Service) checkBlack(word string) (state bool, err error) {
	var (
		black searchModel.Black
	)
	if err = s.dao.DB.Model(&searchModel.Black{}).
		Where("deleted = ?", searchModel.NotDelete).Where("searchword = ?", word).
		First(&black).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.checkBlack get First error(%v)", err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return true, nil
}

//checkInter checkout intervene
func (s *Service) checkInter(word string, id int) (state bool, err error) {
	var (
		intervene searchModel.Intervene
	)
	dataStr := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	query := s.dao.DB.Model(&searchModel.Intervene{}).
		Where("searchword = ?", word)
	if id != 0 {
		query = query.Where("id != ?", id)
	}
	//取未删除且结束时间大于当前时间的词
	query = query.Where("deleted = ?", searchModel.NotDelete).Where("etime > ?", dataStr)
	if err = query.First(&intervene).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.checkInter get First error(%v)", err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return true, nil
}

//checkTimeConflict checkout intervene time conflict
func (s *Service) checkTimeConflict(i searchModel.InterveneAdd, id int) (state bool, err error) {
	var (
		c          int
		black      []searchModel.Black
		blackSlice []string
	)
	if black, err = s.BlackList(); err != nil {
		log.Error("searchSrv.HotList Black error(%v)", err)
	}
	for _, v := range black {
		blackSlice = append(blackSlice, v.Searchword)
	}
	query := s.dao.DB.Model(&searchModel.Intervene{}).
		Where("rank = ?", i.Rank).
		Where("stime < ?", i.Etime).
		Where("etime > ?", i.Stime).
		Where("searchword not in (?)", blackSlice).
		Where("deleted = ?", searchModel.NotDelete)
	if id != 0 {
		query = query.Where("id != ?", id)
	}
	if err = query.Count(&c).Error; err != nil {
		log.Error("searchSrv.checkTimeConflict Count error(%v)", err)
		return
	}
	if c > 0 {
		return true, nil
	}
	return false, nil
}

//AddInter add intervene word
func (s *Service) AddInter(c *bm.Context, v searchModel.InterveneAdd, person string, uid int64) (err error) {
	var (
		state bool
	)
	if state, err = s.checkBlack(v.Searchword); err != nil {
		log.Error("searchSrv.addInter checkBlack error(%v)", err)
		return
	}
	if state {
		err = fmt.Errorf("所添加的词在黑名单中已存在")
		return
	}
	if state, err = s.checkInter(v.Searchword, 0); err != nil {
		log.Error("searchSrv.addInter checkBlack error(%v)", err)
		return
	}
	if state {
		err = fmt.Errorf("干预词已存在")
		return
	}
	if state, err = s.checkTimeConflict(v, 0); err != nil {
		log.Error("searchSrv.addInter checkTimeConflict error(%v)", err)
		return
	}
	if state {
		err = fmt.Errorf("相同时间内，该位置已存在搜索词")
		return
	}
	if err = s.dao.DB.Model(&searchModel.InterveneAdd{}).Create(&v).Error; err != nil {
		log.Error("searchSrv.AddIntervene Create error(%v)", err)
		return
	}
	s.dao.SetSearchAuditStat(c, _HotPubState, false)
	if err = s.dao.SetSearchAuditStat(c, _HotPubState, false); err != nil {
		log.Error("searchSrv.DelBlack SetPubStat error(%v)", err)
		return
	}
	obj := map[string]interface{}{
		"value": v,
	}
	if err = Log.AddLog(searchModel.Business, person, uid, 0, searchModel.ActionAddInter, obj); err != nil {
		log.Error("searchSrv.addInter AddLog error(%v)", err)
		return
	}
	return
}

//UpdateInter update intervene word
func (s *Service) UpdateInter(c *bm.Context, v searchModel.InterveneAdd, id int, person string, uid int64) (err error) {
	var (
		state bool
	)
	if state, err = s.checkBlack(v.Searchword); err != nil {
		log.Error("searchSrv.UpdateInter checkBlack error(%v)", err)
		return
	}
	if state {
		err = fmt.Errorf("所添加的词在黑名单中已存在")
		return
	}
	if state, err = s.checkInter(v.Searchword, id); err != nil {
		log.Error("searchSrv.UpdateInter checkInter error(%v)", err)
		return
	}
	if state {
		err = fmt.Errorf("干预词已存在")
		return
	}
	if state, err = s.checkTimeConflict(v, id); err != nil {
		log.Error("searchSrv.UpdateInter checkTimeConflict error(%v)", err)
		return
	}
	if state {
		err = fmt.Errorf("相同时间内，该位置已存在搜索词")
		return
	}
	if v.Tag == "" {
		v2 := map[string]interface{}{
			"searchword": v.Searchword,
			"rank":       v.Rank,
			"tag":        v.Tag,
			"stime":      v.Stime,
			"etime":      v.Etime,
		}
		if err = s.dao.DB.Model(&searchModel.InterveneAdd{}).
			Where("id = ?", id).Updates(v2).Error; err != nil {
			log.Error("searchSrv.UpdateInter Update error(%v)", err)
			return
		}
	} else {
		if err = s.dao.DB.Model(&searchModel.InterveneAdd{}).
			Where("id = ?", id).Updates(&v).Error; err != nil {
			log.Error("searchSrv.UpdateInter Update error(%v)", err)
			return
		}
	}

	if err = s.dao.SetSearchAuditStat(c, _HotPubState, false); err != nil {
		log.Error("searchSrv.DelBlack SetPubStat error(%v)", err)
		return
	}
	obj := map[string]interface{}{
		"value": v,
		"id":    id,
	}
	if err = Log.AddLog(searchModel.Business, person, uid, int64(id), searchModel.ActionUpdateInter, obj); err != nil {
		log.Error("searchSrv.addInter AddLog error(%v)", err)
		return
	}
	return
}

//UpdateSearch update search hot tag
func (s *Service) UpdateSearch(c *bm.Context, tag string, id int, person string, uid int64) (err error) {
	if err = s.dao.DB.Model(&searchModel.History{}).
		Where("id = ?", id).Update("tag", tag).Error; err != nil {
		log.Error("searchSrv.UpdateSearch Update error(%v)", err)
		return
	}
	if err = s.dao.SetSearchAuditStat(c, _HotPubState, false); err != nil {
		log.Error("searchSrv.DelBlack SetPubStat error(%v)", err)
		return
	}
	obj := map[string]interface{}{
		"value": tag,
		"id":    id,
	}
	if err = Log.AddLog(searchModel.Business, person, uid, int64(id), searchModel.ActionUpdateSearch, obj); err != nil {
		log.Error("searchSrv.UpdateSearch AddLog error(%v)", err)
		return
	}
	return
}

//DeleteHot delete hot word
func (s *Service) DeleteHot(c context.Context, id int, t uint8, person string, uid int64) (err error) {
	if t == searchModel.HotAI {
		//删除AI热词
		if err = s.dao.DB.Model(&searchModel.History{}).
			Where("id = ?", id).Update("deleted", searchModel.Delete).Error; err != nil {
			log.Error("searchSrv.DeleteHot Update AI error(%v)", err)
			return
		}
	} else if t == searchModel.HotOpe {
		//删除运营热词
		if err = s.dao.DB.Model(&searchModel.Intervene{}).
			Where("id = ?", id).Update("deleted", searchModel.Delete).Error; err != nil {
			log.Error("searchSrv.DeleteHot Update Operate error(%v)", err)
			return
		}
	}
	//删除热词之后 立即发布新数据
	if err = s.SetHotPub(c, person, uid); err != nil {
		return
	}
	obj := map[string]interface{}{
		"type": t,
		"id":   id,
	}
	if err = Log.AddLog(searchModel.Business, person, uid, int64(id), searchModel.ActionDeleteHot, obj); err != nil {
		log.Error("searchSrv.DeleteHot AddLog error(%v)", err)
		return
	}
	return
}

//DeleteDark delete dark word
func (s *Service) DeleteDark(c context.Context, id int, person string, uid int64) (err error) {
	if err = s.dao.DB.Model(&searchModel.Dark{}).
		Where("id = ?", id).Update("deleted", searchModel.Delete).Error; err != nil {
		log.Error("searchSrv.DeleteDark Update error(%v)", err)
		return
	}
	if err = s.SetDarkPub(c, person, uid); err != nil {
		return
	}
	obj := map[string]interface{}{
		"id": id,
	}
	if err = Log.AddLog(searchModel.Business, person, uid, int64(id), searchModel.ActionDeleteDark, obj); err != nil {
		log.Error("searchSrv.DeleteDark AddLog error(%v)", err)
		return
	}
	return
}

//OpenAddDarkword open api for search add dark word
func (s *Service) OpenAddDarkword(c context.Context, values searchModel.OpenDark) (err error) {
	if err = s.dao.DB.Model(&searchModel.Dark{}).
		Where("atime = ?", values.Date).Update("deleted", searchModel.Delete).Error; err != nil {
		log.Error("searchSrv.OpenAddDarkword Update error(%v)", err)
		return
	}
	for _, v := range values.Values {
		dark := searchModel.Dark{
			Searchword: v.Searchword,
			PV:         v.PV,
			Atime:      values.Date,
		}
		if err = s.dao.DB.Model(&searchModel.Dark{}).Create(&dark).Error; err != nil {
			log.Error("searchSrv.OpenAddDarkword Create error(%v)", err)
			return
		}
	}
	//如果有黑马词同步过来 则更新发布状态为false
	if err = s.dao.SetSearchAuditStat(c, _DarkAutoPubState, false); err != nil {
		log.Error("searchSrv.DelBlack SetPubStat error(%v)", err)
		return
	}
	obj := map[string]interface{}{
		"value": values,
	}
	if err = Log.AddLog(searchModel.Business, "SEARCH", 0, 0, searchModel.ActionOpenAddDark, obj); err != nil {
		log.Error("searchSrv.OpenAddDarkword AddLog error(%v)", err)
		return
	}
	return
}

//OpenAddHotword open api for search add hotword
func (s *Service) OpenAddHotword(c *bm.Context, values searchModel.OpenHot) (err error) {
	if err = s.dao.DB.Model(&searchModel.Hot{}).Where("atime = ?", values.Date).Delete(&searchModel.Hot{}).Error; err != nil {
		log.Error("searchSrv.OpenAddHotword Delete error(%v)", err)
		return
	}
	for _, v := range values.Values {
		hot := searchModel.Hot{
			Searchword: v.Searchword,
			PV:         v.PV,
			Atime:      values.Date,
		}
		if err = s.dao.DB.Model(&searchModel.Hot{}).Create(&hot).Error; err != nil {
			log.Error("searchSrv.OpenAddHotword Create error(%v)", err)
			return
		}
	}
	//如果有搜索热词同步过来 则更新自动发布状态为false
	if err = s.dao.SetSearchAuditStat(c, _HotAutoPubState, false); err != nil {
		log.Error("searchSrv.DelBlack SetPubStat error(%v)", err)
		return
	}
	obj := map[string]interface{}{
		"value": values,
	}
	if err = Log.AddLog(searchModel.Business, "SEARCH", 0, 0, searchModel.ActionOpenAddHot, obj); err != nil {
		log.Error("searchSrv.OpenAddHotword AddLog error(%v)", err)
		return
	}
	return
}

//GetHotPub get hotword publish from mc
func (s *Service) GetHotPub(c *bm.Context) (hot []searchModel.Intervene, err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
	)
	conn = s.dao.MC.Get(c)
	defer conn.Close()
	if item, err = conn.Get(_HotPubValue); err != nil {
		if err == memcache.ErrNotFound {
			return nil, nil
		}
		return
	}
	if err = conn.Scan(item, &hot); err != nil {
		return
	}
	return
}

//GetDarkPub get darkword publish from mc
func (s *Service) GetDarkPub(c *bm.Context) (dark []searchModel.Dark, err error) {
	var (
		conn memcache.Conn
		item *memcache.Item
	)
	conn = s.dao.MC.Get(c)
	defer conn.Close()
	if item, err = conn.Get(_DarkPubValue); err != nil {
		if err == memcache.ErrNotFound {
			return nil, nil
		}
		return
	}
	if err = conn.Scan(item, &dark); err != nil {
		return
	}
	return
}

//SetHotPub set hotword publish to mc
func (s *Service) SetHotPub(c context.Context, person string, uid int64) (err error) {
	var (
		conn        memcache.Conn
		hot         []searchModel.Intervene
		searchCount int
	)
	conn = s.dao.MC.Get(c)
	defer conn.Close()
	//只能发布当天的数据
	//从DB中取今天的数据
	if hot, searchCount, err = s.HotwordFromDB(time.Now().Format("2006-01-02")); err != nil {
		log.Error("searchSrv.SetHoGetSearHisValuetPub HotwordFromDB error(%v)", err)
		return
	}
	itemJSON := &memcache.Item{
		Key:        _HotPubValue,
		Flags:      memcache.FlagJSON,
		Object:     hot,
		Expiration: 0,
	}
	if err = conn.Set(itemJSON); err != nil {
		log.Error("searchSrv.SetHotPub conn.Set error(%v)", err)
		return
	}
	if searchCount == 0 {
		//证明搜索没有推数据过来 设置搜索的数据为假
		if err = s.dao.SetSearchAuditStat(c, _HotPubSearchState, false); err != nil {
			log.Error("searchSrv.SetHotPub SetSearchPubStat error(%v)", err)
			return
		}
	} else {
		//证明搜索有推数据过来 设置搜索的数据为真
		if err = s.dao.SetSearchAuditStat(c, _HotPubSearchState, true); err != nil {
			log.Error("searchSrv.SetHotPub SetSearchPubStat error(%v)", err)
			return
		}
	}
	//设置自动发布状态为true
	if err = s.dao.SetSearchAuditStat(c, _HotAutoPubState, true); err != nil {
		log.Error("searchSrv.SetHotPub SetPubStat error(%v)", err)
		return
	}
	//设置运营发布状态为true
	if err = s.dao.SetSearchAuditStat(c, _HotPubState, true); err != nil {
		log.Error("searchSrv.SetHotPub SetPubStat error(%v)", err)
		return
	}
	obj := map[string]interface{}{
		"value": hot,
	}
	if err = Log.AddLog(searchModel.Business, person, uid, 0, searchModel.ActionPublishHot, obj); err != nil {
		log.Error("searchSrv.SetHotPub AddLog error(%v)", err)
		return
	}
	if err = s.HotPubLog(hot); err != nil {
		log.Error("searchSrv.SetHotPub HotPubLog error(%v)", err)
		return
	}
	return
}

//HotPubLog hotword publish log
func (s *Service) HotPubLog(hot []searchModel.Intervene) (err error) {
	t := time.Now().Unix()
	for _, v := range hot {
		w := searchModel.HotPubLog{
			Searchword: v.Searchword,
			Position:   v.Rank,
			Pv:         v.Pv,
			Tag:        v.Tag,
			Stime:      v.Stime,
			Etime:      v.Etime,
			Atime:      time.Now().Format("2006-01-02"),
			Groupid:    t,
		}
		if err = s.dao.DB.Model(&searchModel.HotPubLog{}).Create(&w).Error; err != nil {
			log.Error("searchSrv.DarkPubLog Create error(%v)", err)
			return
		}
	}
	return
}

//GetHotPubLog get hotword publish log
func (s *Service) GetHotPubLog(date string) (hotout []searchModel.Intervene, pub bool, err error) {
	var (
		//hotout searchModel.HotwordOut
		logs []searchModel.HotPubLog
	)
	l := searchModel.HotPubLog{}
	if err = s.dao.DB.Model(&searchModel.HotPubLog{}).Where("atime = ?", date).Order("groupid desc").
		First(&l).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.GetHotPubLog First error(%v)", err)
		return
	}
	//证明没有发布过
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	//取最大的groupid的值
	if err = s.dao.DB.Model(&searchModel.HotPubLog{}).Where("groupid = ?", l.Groupid).
		Find(&logs).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.GetHotPubLog Find error(%v)", err)
		return
	}
	for _, v := range logs {
		a := searchModel.Intervene{
			Searchword: v.Searchword,
			Rank:       v.Position,
			Pv:         v.Pv,
			Tag:        v.Tag,
			Stime:      v.Stime,
			Etime:      v.Etime,
		}
		hotout = append(hotout, a)
	}
	return hotout, true, nil
}

//GetDarkPubLog get darkword publish log
func (s *Service) GetDarkPubLog(date string) (darkout []searchModel.Dark, pub bool, err error) {
	var (
		//hotout searchModel.HotwordOut
		logs []searchModel.DarkPubLog
	)
	l := searchModel.DarkPubLog{}
	if err = s.dao.DB.Model(&searchModel.DarkPubLog{}).Where("atime = ?", date).Order("groupid desc").
		First(&l).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.GetDarkPubLog First error(%v)", err)
		return
	}
	//证明没有发布过
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	//取最大的groupid的值
	if err = s.dao.DB.Model(&searchModel.DarkPubLog{}).Where("groupid = ?", l.Groupid).
		Find(&logs).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("searchSrv.GetDarkPubLog Find error(%v)", err)
		return
	}
	for _, v := range logs {
		a := searchModel.Dark{
			Searchword: v.Searchword,
			PV:         v.Pv,
		}
		darkout = append(darkout, a)
	}
	return darkout, true, nil
}

//SetDarkPub set darkword to mc
func (s *Service) SetDarkPub(c context.Context, person string, uid int64) (err error) {
	var (
		conn        memcache.Conn
		dark        []searchModel.Dark
		searchCount int
	)
	conn = s.dao.MC.Get(c)
	defer conn.Close()
	//只能发布当天的数据
	//从DB中取今天的数据
	if dark, searchCount, err = s.DarkwordFromDB(time.Now().Format("2006-01-02")); err != nil {
		log.Error("searchSrv.SetHotPub HotwordFromDB error(%v)", err)
		return
	}
	itemJSON := &memcache.Item{
		Key:        _DarkPubValue,
		Flags:      memcache.FlagJSON,
		Object:     dark,
		Expiration: 0,
	}
	if err = conn.Set(itemJSON); err != nil {
		log.Error("searchSrv.SetHotPub conn.Set error(%v)", err)
		return
	}
	if searchCount == 0 {
		//证明搜索没有推数据过来 设置搜索的数据为假
		if err = s.dao.SetSearchAuditStat(c, _DarkPubSearchState, false); err != nil {
			log.Error("searchSrv.SetDarkPub SetSearchPubStat error(%v)", err)
			return
		}
	} else {
		//证明搜索有推数据过来 设置搜索的数据为真
		if err = s.dao.SetSearchAuditStat(c, _DarkPubSearchState, true); err != nil {
			log.Error("searchSrv.SetDarkPub SetSearchPubStat error(%v)", err)
			return
		}
	}
	if err = s.dao.SetSearchAuditStat(c, _DarkPubState, true); err != nil {
		log.Error("searchSrv.SetHotPub SetPubStat error(%v)", err)
		return
	}
	if err = s.dao.SetSearchAuditStat(c, _DarkAutoPubState, true); err != nil {
		log.Error("searchSrv.SetHotPub SetPubStat error(%v)", err)
		return
	}
	obj := map[string]interface{}{
		"value": dark,
	}
	if err = Log.AddLog(searchModel.Business, person, uid, 0, searchModel.ActionPublishDark, obj); err != nil {
		log.Error("searchSrv.SetDarkPub AddLog error(%v)", err)
		return
	}
	if err = s.DarkPubLog(dark); err != nil {
		log.Error("searchSrv.SetDarkPub DarkPubLog error(%v)", err)
		return
	}
	return
}

//DarkPubLog get darkword publish log
func (s *Service) DarkPubLog(dark []searchModel.Dark) (err error) {
	t := time.Now().Unix()
	for _, v := range dark {
		w := searchModel.DarkPubLog{
			Searchword: v.Searchword,
			Pv:         v.PV,
			Atime:      time.Now().Format("2006-01-02"),
			Groupid:    t,
		}
		if err = s.dao.DB.Model(&searchModel.DarkPubLog{}).Create(&w).Error; err != nil {
			log.Error("searchSrv.DarkPubLog Create error(%v)", err)
			return
		}
	}
	return
}
