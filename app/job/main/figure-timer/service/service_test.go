package service

import (
	"context"
	"flag"
	"math/rand"
	"testing"
	"time"

	"go-common/app/job/main/figure-timer/conf"
	"go-common/app/job/main/figure-timer/dao/mock_dao"
	"go-common/app/job/main/figure-timer/model"
	"go-common/library/log"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
	c = context.TODO()
)

func init() {
	var err error
	flag.Set("conf", "../cmd/figure-timer-job-test.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	s = New(conf.Conf)

}

func TestSplitMids(t *testing.T) {
	Convey("TEST split mids", t, func() {
		var (
			mids        []int64
			concurrency int64 = 100
			midsSize    int64 = 233
		)
		for i := int64(0); i < midsSize; i++ {
			mids = append(mids, i)
		}
		smids := splitMids(mids, concurrency)
		actualC := concurrency
		if actualC == 0 {
			actualC = 1
		}
		if actualC > midsSize {
			actualC = midsSize
		}
		So(len(smids), ShouldEqual, 233/3+1)
		total := 0
		for s := range smids {
			total += len(smids[s])
		}
		So(total, ShouldEqual, len(mids))
	})
}

func TestVersion(t *testing.T) {
	Convey("TEST week version", t, func() {
		var (
			ts = []time.Time{
				time.Date(2017, 9, 25, 0, 0, 0, 0, time.Local),
				time.Date(2017, 10, 1, 23, 59, 59, 99, time.Local),
				time.Date(2017, 9, 29, 12, 12, 12, 12, time.Local),
			}
			actual = time.Date(2017, 9, 25, 0, 0, 0, 0, time.Local)
		)
		for _, t := range ts {
			ver := weekVersion(t)
			So(ver, ShouldResemble, actual.Unix())
		}
	})
}

func TestCalcFigure(t *testing.T) {
	Convey("TEST calc offset", t, func() {
		var (
			max = 100.0
			k   = 0.1
			x1  = 10.0
			x2  = 0.0
		)
		r := calcOffset(max, k, x1)
		So(r, ShouldEqual, 63)
		r = calcOffset(max, k, x2)
		So(r, ShouldEqual, 0)
	})
	Convey("TEST calc figure", t, func() {
		var (
			userInfo = &model.UserInfo{
				Mid:          233,
				Exp:          673,
				SpyScore:     100,
				ArchiveViews: 0,
				VIPStatus:    0,
			}
			actionCounters = []*model.ActionCounter{
				{
					Mid:           233,
					CoinCount:     0,
					ReplyCount:    0,
					CoinLowRisk:   0,
					CoinHighRisk:  0,
					ReplyLowRisk:  0,
					ReplyHighRisk: 0,
					ReplyLiked:    0,
					ReplyUnliked:  0,
					Version:       time.Date(2017, 10, 2, 0, 0, 0, 0, time.Local),
				},
			}
			records = []*model.FigureRecord{
				{
					XPosCreativity: 0,
					XNegFriendly:   0,
					XPosFriendly:   0,
					XNegLawful:     80,
					Version:        time.Date(2017, 9, 25, 0, 0, 0, 0, time.Local),
				},
			}
			weekVer = time.Date(2018, 1, 2, 0, 0, 0, 0, time.Local).Unix()
		)
		figure, newRecord := s.CalcFigure(c, userInfo, actionCounters, records, weekVer)
		So(figure, ShouldNotBeNil)
		So(newRecord, ShouldNotBeNil)
		So(figure.LawfulScore, ShouldEqual, 2859)
		So(newRecord.XNegLawful, ShouldEqual, 0)
		So(newRecord.XPosLawful, ShouldEqual, 0)
	})
}

func TestRank(t *testing.T) {
	rand.Seed(time.Now().Unix())
	rank.Init()
	for i := 0; i < 1; i++ {
		rank.AddScore(rand.Int31n(5000))
	}
	s.calcRank(c, time.Now().Unix())
}

func TestFix(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDao := mock_dao.NewMockDaoInt(ctrl)
	mockDao.EXPECT().Figures(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		[]*model.Figure{
			{Mid: 1000},
			{Mid: 2000},
		},
		true,
		nil,
	).AnyTimes()
	mockDao.EXPECT().CalcRecords(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(
		[]*model.FigureRecord{
			{Version: time.Date(2018, 1, 2, 0, 0, 0, 0, time.Local)},
		},
		nil,
	).AnyTimes()
	mockDao.EXPECT().ActionCounter(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&model.ActionCounter{
			PayLiveMoney: 2000,
			Version:      time.Date(2018, 1, 2, 0, 0, 0, 0, time.Local),
		},
		nil,
	).AnyTimes()
	mockDao.EXPECT().PutCalcRecord(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	s.dao = mockDao

	Convey("TEST fix record", t, func() {
		s.fixproc()
	})
}
