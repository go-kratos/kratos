package service

import (
	"context"
	"encoding/json"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/esports/conf"
	"go-common/app/interface/main/esports/dao"
	"go-common/app/interface/main/esports/model"
	arcclient "go-common/app/service/main/archive/api"
	favrpc "go-common/app/service/main/favorite/api/gorpc"
	"go-common/library/log"
	"go-common/library/sync/pipeline/fanout"

	"github.com/robfig/cron"
)

const (
	_perPage       = 100
	_lolType       = 1
	_dotaType      = 2
	_firstPage     = "1"
	_lolGame       = "lol/games"
	_dotaGame      = "dota2/games"
	_lolItems      = "lol/items"
	_dotaItems     = "dota2/items"
	_lolChampions  = "lol/champions"
	_lolHeroes     = "dota2/heroes"
	_lolSpells     = "lol/spells"
	_dotaAbilities = "dota2/abilities"
	_lolPlayers    = "lol/players"
	_dotaPlayers   = "dota2/players"
)

// Service service struct.
type Service struct {
	c   *conf.Config
	dao *dao.Dao
	// rpc
	fav *favrpc.Service
	// cache proc
	cache                     *fanout.Fanout
	arcClient                 arcclient.ArchiveClient
	lolGameMap, dotaGameMap   *model.SyncGame
	lolItemsMap, dotaItemsMap *model.SyncItem
	lolChampions              *model.SyncChampion
	dotaHeroes                *model.SyncHero
	lolSpells, dotaAbilities  *model.SyncInfo
	lolPlayers, dotaPlayers   *model.SyncInfo
	// cron
	cron *cron.Cron
}

// New new service.
func New(c *conf.Config) *Service {
	s := &Service{
		c:     c,
		dao:   dao.New(c),
		fav:   favrpc.New2(c.FavoriteRPC),
		cache: fanout.New("cache"),
		lolGameMap: &model.SyncGame{
			Data: make(map[int64][]*model.Game),
		},
		dotaGameMap: &model.SyncGame{
			Data: make(map[int64][]*model.Game),
		},
		lolItemsMap: &model.SyncItem{
			Data: make(map[int64]*model.Item),
		},
		dotaItemsMap: &model.SyncItem{
			Data: make(map[int64]*model.Item),
		},
		lolChampions: &model.SyncChampion{
			Data: make(map[int64]*model.Champion),
		},
		dotaHeroes: &model.SyncHero{
			Data: make(map[int64]*model.Hero),
		},
		lolSpells: &model.SyncInfo{
			Data: make(map[int64]*model.LdInfo),
		},
		dotaAbilities: &model.SyncInfo{
			Data: make(map[int64]*model.LdInfo),
		},
		lolPlayers: &model.SyncInfo{
			Data: make(map[int64]*model.LdInfo),
		},
		dotaPlayers: &model.SyncInfo{
			Data: make(map[int64]*model.LdInfo),
		},
		cron: cron.New(),
	}
	var err error
	if s.arcClient, err = arcclient.NewClient(c.ArcClient); err != nil {
		panic(err)
	}
	go s.loadKnockTreeCache()
	go s.loadLdGame()
	go s.createCron()
	return s
}

// Ping ping service.
func (s *Service) Ping(c context.Context) (err error) {
	if err = s.dao.Ping(c); err != nil {
		log.Error("s.dao.Ping error(%v)", err)
	}
	return
}

// loadCache load cache
func (s *Service) loadKnockTreeCache() {
	for {
		s.BuildKnockTree(context.Background())
		time.Sleep(time.Duration(conf.Conf.Rule.KnockTree))
	}
}

func (s *Service) loadLdGame() {
	var (
		contestDatas []*model.Contest
		err          error
	)
	for {
		if contestDatas, err = s.dao.ContestDatas(context.Background()); err != nil {
			log.Error("loadLeida  s.dao.ContestDatas error(%v)", err)
			time.Sleep(time.Second)
			continue
		}
		for _, data := range contestDatas {
			tmp := data
			go s.setGamesMap(tmp)
		}
		time.Sleep(time.Duration(conf.Conf.Leidata.AfterSleep))
	}
}
func (s *Service) createCron() {
	go s.lolPlayersCron()
	go s.dotaPlayersCron()
	go s.infoCron()
	s.cron.AddFunc(s.c.Leidata.LolPlayersCron, s.lolPlayersCron)
	s.cron.AddFunc(s.c.Leidata.DotaPlayersCron, s.dotaPlayersCron)
	s.cron.AddFunc(s.c.Leidata.InfoCron, s.infoCron)
	s.cron.Start()
}
func (s *Service) lolPlayersCron() {
	go s.loadLdPages(_lolPlayers)
	log.Info("createCron lolPlayersCron start")
}
func (s *Service) dotaPlayersCron() {
	go s.loadLdPages(_dotaPlayers)
	log.Info("createCron dotaPlayersCron start")
}
func (s *Service) infoCron() {
	go s.loadLdPages(_lolItems)
	go s.loadLdPages(_dotaItems)
	go s.loadLdPages(_lolSpells)
	go s.loadLdPages(_dotaAbilities)
	go s.loadLdPages(_lolChampions)
	go s.loadLdPages(_lolHeroes)
	log.Info("createCron infoCron start")
}

func (s *Service) setGamesMap(data *model.Contest) {
	var (
		err     error
		params  url.Values
		rs      json.RawMessage
		games   []*model.Game
		endTime time.Time
		isTime  bool
	)
	params = url.Values{}
	params.Set("match_id", strconv.FormatInt(data.MatchID, 10))
	if data.Etime > 0 {
		endTime = time.Unix(data.Etime, 0).Add(time.Duration(s.c.Leidata.EndSleep))
		if time.Now().Unix() > endTime.Unix() {
			isTime = true
		}
	}
	if !isTime && data.Stime > 0 && time.Now().Unix() < data.Stime {
		isTime = true
	}
	if data.DataType == _lolType {
		if _, ok := s.lolGameMap.Data[data.MatchID]; ok && isTime {
			return
		}
		if rs, _, err = s.leida(params, _lolGame); err == nil {
			if err = json.Unmarshal(rs, &games); err == nil {
				s.lolGameMap.Lock()
				s.lolGameMap.Data[data.MatchID] = games
				s.lolGameMap.Unlock()
			}
		}
	} else if data.DataType == _dotaType {
		if _, ok := s.dotaGameMap.Data[data.MatchID]; ok && isTime {
			return
		}
		if rs, _, err = s.leida(params, _dotaGame); err == nil {
			if err = json.Unmarshal(rs, &games); err == nil {
				s.dotaGameMap.Lock()
				s.dotaGameMap.Data[data.MatchID] = games
				s.dotaGameMap.Unlock()
			}
		}
	}
}

func (s *Service) loadLdPages(tp string) {
	var (
		err    error
		params url.Values
		count  int
	)
	params = url.Values{}
	params.Set("page", _firstPage)
	params.Set("per_page", strconv.Itoa(_perPage))
	if count, err = s.setPages(tp, params); err != nil {
		return
	}
	for i := 2; i <= count; i++ {
		time.Sleep(time.Second)
		params.Set("page", strconv.Itoa(i))
		params.Set("per_page", strconv.Itoa(_perPage))
		s.setPages(tp, params)
	}
}

func (s *Service) setPages(tp string, params url.Values) (count int, err error) {
	var (
		rs        json.RawMessage
		items     []*model.Item
		infos     []*model.LdInfo
		champions []*model.Champion
		heroes    []*model.Hero
	)
	switch tp {
	case _lolItems:
		if rs, count, err = s.leida(params, _lolItems); err == nil {
			if err = json.Unmarshal(rs, &items); err == nil {
				for _, item := range items {
					s.lolItemsMap.Lock()
					s.lolItemsMap.Data[item.ID] = item
					s.lolItemsMap.Unlock()
				}
			}
		}
	case _dotaItems:
		if rs, count, err = s.leida(params, _dotaItems); err == nil {
			if err = json.Unmarshal(rs, &items); err == nil {
				for _, item := range items {
					s.dotaItemsMap.Lock()
					s.dotaItemsMap.Data[item.ID] = item
					s.dotaItemsMap.Unlock()
				}
			}
		}
	case _lolSpells:
		if rs, count, err = s.leida(params, _lolSpells); err == nil {
			if err = json.Unmarshal(rs, &infos); err == nil {
				for _, info := range infos {
					s.lolSpells.Lock()
					s.lolSpells.Data[info.ID] = info
					s.lolSpells.Unlock()
				}
			}
		}
	case _dotaAbilities:
		if rs, count, err = s.leida(params, _dotaAbilities); err == nil {
			if err = json.Unmarshal(rs, &infos); err == nil {
				for _, info := range infos {
					s.dotaAbilities.Lock()
					s.dotaAbilities.Data[info.ID] = info
					s.dotaAbilities.Unlock()
				}
			}
		}
	case _lolPlayers:
		if rs, count, err = s.leida(params, _lolPlayers); err == nil {
			if err = json.Unmarshal(rs, &infos); err == nil {
				for _, info := range infos {
					s.lolPlayers.Lock()
					s.lolPlayers.Data[info.ID] = info
					s.lolPlayers.Unlock()
				}
			}
		}
	case _dotaPlayers:
		if rs, count, err = s.leida(params, _dotaPlayers); err == nil {
			if err = json.Unmarshal(rs, &infos); err == nil {
				for _, info := range infos {
					s.dotaPlayers.Lock()
					s.dotaPlayers.Data[info.ID] = info
					s.dotaPlayers.Unlock()
				}
			}
		}
	case _lolChampions:
		if rs, count, err = s.leida(params, _lolChampions); err == nil {
			if err = json.Unmarshal(rs, &champions); err == nil {
				for _, champion := range champions {
					s.lolChampions.Lock()
					s.lolChampions.Data[champion.ID] = champion
					s.lolChampions.Unlock()
				}
			}
		}
	case _lolHeroes:
		if rs, count, err = s.leida(params, _lolHeroes); err == nil {
			if err = json.Unmarshal(rs, &heroes); err == nil {
				for _, hero := range heroes {
					s.dotaHeroes.Lock()
					s.dotaHeroes.Data[hero.ID] = hero
					s.dotaHeroes.Unlock()
				}
			}
		}
	}
	return
}

func (s *Service) leida(params url.Values, route string) (rs []byte, count int, err error) {
	var body, orginBody []byte
	params.Del("route")
	params.Set("key", s.c.Leidata.Key)
	url := s.c.Leidata.URL + "/" + route + "?" + params.Encode()
	for i := 0; i < s.c.Leidata.Retry; i++ {
		if body, err = s.dao.Leida(context.Background(), url); err != nil {
			time.Sleep(time.Second)
			continue
		}
		bodyStr := string(body[:])
		if bodyStr == "" {
			time.Sleep(time.Second)
			continue
		}
		rsPos := strings.Index(bodyStr, "[")
		if rsPos > -1 {
			orginBody = body
			body = []byte(bodyStr[rsPos:])
		} else {
			time.Sleep(time.Second)
			continue
		}
		rs = body
		totalPos := strings.Index(bodyStr, "X-Total:")
		if totalPos > 0 {
			s := string(orginBody[totalPos+9 : rsPos-4])
			if t, e := strconv.ParseFloat(s, 64); e == nil {
				count = int(math.Ceil(t / float64(_perPage)))
			}
		}
		break
	}
	if err != nil {
		log.Error("json.Unmarshal url(%s) body(%s) error(%v)", url, string(body), err)
	}
	return
}
