package bws

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/activity/conf"
	"go-common/app/interface/main/activity/dao/bws"
	bwsmdl "go-common/app/interface/main/activity/model/bws"
	accapi "go-common/app/service/main/account/api"
	suitmdl "go-common/app/service/main/usersuit/model"
	suitrpc "go-common/app/service/main/usersuit/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/sync/pipeline/fanout"
)

const (
	_accountBlocked = 1
	_allType        = 0
	_dpType         = 1
	_gameType       = 2
	_clockinType    = 3
	_eggType        = 4
	_dp             = "dp"
	_game           = "game"
	_clockin        = "clockin"
	_egg            = "egg"
	_noAward        = 0
	_awardAlready   = 2
	_initLinkType   = 5
)

var (
	_emptPoints        = make([]*bwsmdl.Point, 0)
	_emptUserPoints    = make([]*bwsmdl.UserPointDetail, 0)
	_emptyUserAchieves = make([]*bwsmdl.UserAchieveDetail, 0)
)

// Service struct
type Service struct {
	c         *conf.Config
	dao       *bws.Dao
	accClient accapi.AccountClient
	suitRPC   *suitrpc.Service2
	// bws admin mids
	allowMids   map[int64]struct{}
	awardMids   map[int64]struct{}
	lotteryMids map[int64]struct{}
	lotteryAids map[int64]struct{}
	cache       *fanout.Fanout
}

// New Service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:       c,
		dao:     bws.New(c),
		suitRPC: suitrpc.New(c.RPCClient2.Suit),
		cache:   fanout.New("cache", fanout.Worker(1), fanout.Buffer(1024)),
	}
	var err error
	if s.accClient, err = accapi.NewClient(c.AccClient); err != nil {
		panic(err)
	}
	s.initMids()
	s.initLotteryAids()
	return
}

func (s *Service) initMids() {
	tmpMids := make(map[int64]struct{}, len(s.c.Rule.BwsMids))
	tmpAward := make(map[int64]struct{}, len(s.c.Rule.BwsMids)+len(s.c.Rule.BwsAwardMids))
	tmpLottery := make(map[int64]struct{}, len(s.c.Rule.BwsMids)+len(s.c.Rule.BwsLotteryMids))
	for _, id := range s.c.Rule.BwsMids {
		tmpMids[id] = struct{}{}
		tmpAward[id] = struct{}{}
		tmpLottery[id] = struct{}{}
	}
	for _, id := range s.c.Rule.BwsAwardMids {
		tmpAward[id] = struct{}{}
	}
	for _, id := range s.c.Rule.BwsLotteryMids {
		tmpLottery[id] = struct{}{}
	}
	s.allowMids = tmpMids
	s.awardMids = tmpAward
	s.lotteryMids = tmpLottery
}

func (s *Service) initLotteryAids() {
	tmp := make(map[int64]struct{}, len(s.c.Rule.BwsLotteryAids))
	for _, id := range s.c.Rule.BwsLotteryAids {
		tmp[id] = struct{}{}
	}
	s.lotteryAids = tmp
}

// User user info.
func (s *Service) User(c context.Context, bid, mid int64, key string) (user *bwsmdl.User, err error) {
	var (
		hp, keyID                          int64
		ac                                 *accapi.CardReply
		points, dps, games, clockins, eggs []*bwsmdl.UserPointDetail
		achErr, pointErr                   error
	)
	if key == "" {
		if key, err = s.midToKey(c, mid); err != nil {
			return
		}
	} else {
		if mid, keyID, err = s.keyToMid(c, key); err != nil {
			return
		}
	}
	user = new(bwsmdl.User)
	if mid != 0 {
		if ac, err = s.accCard(c, mid); err != nil {
			log.Error("User s.accCard(%d) error(%v)", mid, err)
			return
		}
	}
	if ac != nil && ac.Card != nil {
		user.User = &bwsmdl.UserInfo{
			Mid:  ac.Card.Mid,
			Name: ac.Card.Name,
			Face: ac.Card.Face,
			Key:  key,
		}
	} else {
		user.User = &bwsmdl.UserInfo{
			Name: strconv.FormatInt(keyID, 10),
			Key:  key,
		}
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if user.Achievements, achErr = s.userAchieves(errCtx, bid, key); achErr != nil {
			log.Error("User s.userAchieves(%d,%s) error(%v)", bid, key, achErr)
		}
		return nil
	})
	group.Go(func() error {
		if points, pointErr = s.userPoints(errCtx, bid, key); pointErr != nil {
			log.Error("User s.userPoints(%d,%s) error(%v)", bid, key, pointErr)
		}
		return nil
	})
	group.Wait()
	if len(user.Achievements) == 0 {
		user.Achievements = _emptyUserAchieves
	}
	user.Items = make(map[string][]*bwsmdl.UserPointDetail, 4)
	gidMap := make(map[int64]int64, len(points))
	for _, v := range points {
		switch v.LockType {
		case _dpType:
			dps = append(dps, v)
		case _gameType:
			if v.Points == v.Unlocked {
				if _, ok := gidMap[v.Pid]; !ok {
					games = append(games, v)
				}
				gidMap[v.Pid] = v.Pid
			}
		case _clockinType:
			clockins = append(clockins, v)
		case _eggType:
			eggs = append(eggs, v)
		}
		hp += v.Points
	}
	user.User.Hp = hp
	if len(dps) == 0 {
		user.Items[_dp] = _emptUserPoints
	} else {
		user.Items[_dp] = dps
	}
	if len(games) == 0 {
		user.Items[_game] = _emptUserPoints
	} else {
		user.Items[_game] = games
	}
	if len(clockins) == 0 {
		user.Items[_clockin] = _emptUserPoints
	} else {
		user.Items[_clockin] = clockins
	}
	if len(eggs) == 0 {
		user.Items[_egg] = _emptUserPoints
	} else {
		user.Items[_egg] = eggs
	}
	return
}

func (s *Service) accCard(c context.Context, mid int64) (ac *accapi.CardReply, err error) {
	var (
		arg = &accapi.MidReq{Mid: mid}
	)
	if ac, err = s.accClient.Card3(c, arg); err != nil || ac == nil {
		log.Error("s.accRPC.Card3(%d) error(%v)", mid, err)
		err = ecode.AnswerAccCallErr
	} else if ac.Card.Silence == _accountBlocked {
		err = ecode.UserDisabled
	}
	return
}

// Binding binding by mid
func (s *Service) Binding(c context.Context, loginMid int64, p *bwsmdl.ParamBinding) (err error) {
	var (
		achieves *bwsmdl.Achievements
		users    *bwsmdl.Users
		checkMid int64
	)
	if _, err = s.accCard(c, loginMid); err != nil {
		log.Error("s.accCard(%d) error(%v)", loginMid, err)
		return
	}
	if checkMid, _, err = s.keyToMid(c, p.Key); err != nil {
		return
	}
	if checkMid != 0 {
		err = ecode.ActivityKeyBindAlready
		return
	}
	if users, err = s.dao.UsersMid(c, loginMid); err != nil {
		err = ecode.ActivityKeyFail
		return
	}
	if users != nil && users.Key != "" {
		err = ecode.ActivityMidBindAlready
		return
	}
	if err = s.dao.Binding(c, loginMid, p); err != nil {
		log.Error("s.dao.Binding mid(%d) key(%s)  error(%v)", loginMid, p.Key, err)
		return
	}
	if s.c.Rule.NeedInitAchieve {
		if achieves, err = s.dao.Achievements(c, p.Bid); err != nil {
			log.Error("s.dao.Achievements error(%v)", err)
			err = ecode.ActivityAchieveFail
			return
		}
		if achieves == nil || len(achieves.Achievements) == 0 {
			err = ecode.ActivityNoAchieve
			return
		}
		for _, achieve := range achieves.Achievements {
			if achieve.LockType == _initLinkType {
				s.addAchieve(c, loginMid, achieve, p.Key)
				break
			}
		}
		var userAchieves []*bwsmdl.UserAchieveDetail
		if userAchieves, err = s.userAchieves(c, p.Bid, p.Key); err != nil {
			log.Error("Binding add suit key(%s) mid(%d) %+v", p.Key, loginMid, err)
			err = nil
		} else {
			for _, v := range userAchieves {
				if v.LockType == _initLinkType {
					continue
				}
				if suitID := v.SuitID; suitID != 0 {
					log.Warn("Binding suit mid(%d) suitID(%d) expire(%d)", loginMid, suitID, s.c.Rule.BwsSuitExpire)
					s.cache.Do(c, func(c context.Context) {
						arg := &suitmdl.ArgGrantByMids{Mids: []int64{loginMid}, Pid: suitID, Expire: s.c.Rule.BwsSuitExpire}
						if e := s.suitRPC.GrantByMids(c, arg); e != nil {
							log.Error("Binding s.suit.GrantByMids(%d,%d) error(%v)", loginMid, suitID, e)
						}
					})
				}
				if _, ok := s.lotteryAids[v.Aid]; ok {
					lotteryAid := v.Aid
					log.Warn("Binding lottery mid(%d) achieve id(%d) expire(%d)", loginMid, lotteryAid, s.c.Rule.BwsSuitExpire)
					s.cache.Do(c, func(c context.Context) {
						s.dao.AddLotteryMidCache(c, lotteryAid, loginMid)
					})
				}
			}
		}
	}
	s.dao.DelCacheUsersKey(c, p.Key)
	s.dao.DelCacheUsersMid(c, loginMid)
	return
}

func (s *Service) isAdmin(mid int64) bool {
	if _, ok := s.allowMids[mid]; ok {
		return true
	}
	return false
}

func (s *Service) midToKey(c context.Context, mid int64) (key string, err error) {
	var users *bwsmdl.Users
	if users, err = s.dao.UsersMid(c, mid); err != nil {
		err = ecode.ActivityKeyFail
		return
	}
	if users == nil || users.Key == "" {
		err = ecode.ActivityNotBind
		return
	}
	key = users.Key
	return
}

func (s *Service) keyToMid(c context.Context, key string) (mid, keyID int64, err error) {
	var users *bwsmdl.Users
	if users, err = s.dao.UsersKey(c, key); err != nil {
		err = ecode.ActivityKeyFail
		return
	}
	if users == nil || users.Key == "" {
		err = ecode.ActivityKeyNotExists
		return
	}
	if users.Mid > 0 {
		mid = users.Mid
	}
	keyID = users.ID
	return
}

func today() string {
	return time.Now().Format("20060102")
}
