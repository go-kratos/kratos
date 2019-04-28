package service

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/feed/conf"
	"go-common/app/service/main/feed/dao"
	feedmdl "go-common/app/service/main/feed/model"
	"go-common/library/cache/redis"
	xtime "go-common/library/time"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	_mid        = int64(27515256)
	_bangumiMid = int64(2)
	_blankMid   = int64(27515280)
	_ip         = "127.0.0.1"
	_seasonID   = int64(1)
	_dataAV     = int64(5463626)
	_arc1       = &api.Arc{Aid: 1, PubDate: xtime.Time(1500262304)}
	_arc2       = &api.Arc{Aid: 2, PubDate: xtime.Time(1500242304)}
	_arc3       = &api.Arc{Aid: 3, PubDate: xtime.Time(1500222304)}
	_arc4       = &api.Arc{Aid: 4, PubDate: xtime.Time(1500202304)}
	_arc5       = &api.Arc{Aid: 5, PubDate: xtime.Time(1500162304)}
	_arc6       = &api.Arc{Aid: 6, PubDate: xtime.Time(1500142304)}
)

var (
	s *Service
)

func WithBlankService(f func(s *Service)) func() {
	return func() {
		ss := &Service{}
		f(ss)
	}
}

func CleanCache() {
	c := context.TODO()
	pool := redis.NewPool(conf.Conf.MultiRedis.Cache)
	pool.Get(c).Do("FLUSHDB")
}

func init() {
	dir, _ := filepath.Abs("../cmd/convey-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = &Service{
		c:      conf.Conf,
		dao:    dao.New(conf.Conf),
		missch: make(chan func(), 1000),
	}
	go s.cacheproc()
}
func WithService(t *testing.T, f func(s *Service)) func() {
	return func() {
		mockCtrl := gomock.NewController(t)
		arcMock := NewMockArcRPC(mockCtrl)
		s.arcRPC = arcMock
		arcMock.EXPECT().Archive3(gomock.Any(), gomock.Any()).Return(&api.Arc{Aid: 100, Author: api.Author{}}, nil).AnyTimes()
		arcMock.EXPECT().Archives3(gomock.Any(), gomock.Any()).Return(map[int64]*api.Arc{
			1: _arc1,
			2: _arc2,
			3: _arc3,
			4: _arc4,
			5: _arc5,
			6: _arc6,
		}, nil).AnyTimes()
		arcMock.EXPECT().UpsPassed2(gomock.Any(), gomock.Any()).Return(map[int64][]*archive.AidPubTime{
			1: []*archive.AidPubTime{&archive.AidPubTime{Aid: _arc1.Aid, PubDate: _arc1.PubDate}, &archive.AidPubTime{Aid: _arc2.Aid, PubDate: _arc2.PubDate}, &archive.AidPubTime{Aid: _arc3.Aid, PubDate: _arc3.PubDate}},
			2: []*archive.AidPubTime{&archive.AidPubTime{Aid: _arc4.Aid, PubDate: _arc4.PubDate}, &archive.AidPubTime{Aid: _arc5.Aid, PubDate: _arc5.PubDate}, &archive.AidPubTime{Aid: _arc6.Aid, PubDate: _arc6.PubDate}},
		}, nil).AnyTimes()

		accMock := NewMockAccRPC(mockCtrl)
		s.accRPC = accMock
		accMock.EXPECT().Attentions3(gomock.Any(), &account.ArgMid{Mid: _mid}).Return([]int64{1, 2, 3}, nil).AnyTimes()
		accMock.EXPECT().Attentions3(gomock.Any(), &account.ArgMid{Mid: _bangumiMid}).Return([]int64{1, 2, 3}, nil).AnyTimes()
		accMock.EXPECT().Attentions3(gomock.Any(), &account.ArgMid{Mid: _blankMid}).Return([]int64{}, nil).AnyTimes()

		banMock := NewMockBangumi(mockCtrl)
		s.bangumi = banMock
		banMock.EXPECT().BangumiPull(gomock.Any(), gomock.Eq(_bangumiMid), gomock.Any()).Return([]int64{1, 2, 3, 4, 5, 6}, nil).AnyTimes()
		banMock.EXPECT().BangumiPull(gomock.Any(), gomock.Any(), gomock.Any()).Return([]int64{}, nil).AnyTimes()
		banMock.EXPECT().BangumiSeasons(gomock.Any(), gomock.Any(), gomock.Any()).Return(map[int64]*feedmdl.Bangumi{
			1: &feedmdl.Bangumi{Title: "title", SeasonID: 1, Ts: 1500142304},
			2: &feedmdl.Bangumi{Title: "title", SeasonID: 2, Ts: 1500142304},
			3: &feedmdl.Bangumi{Title: "title", SeasonID: 3, Ts: 1500142304},
			4: &feedmdl.Bangumi{Title: "title", SeasonID: 4, Ts: 1500142304},
			5: &feedmdl.Bangumi{Title: "title", SeasonID: 5, Ts: 1500142304},
			6: &feedmdl.Bangumi{Title: "title", SeasonID: 6, Ts: 1500142304},
		}, nil).AnyTimes()
		Reset(func() { CleanCache() })
		f(s)
		mockCtrl.Finish()
	}
}
