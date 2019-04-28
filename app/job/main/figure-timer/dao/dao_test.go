package dao

import (
	"context"
	// "encoding/binary"
	"flag"
	"testing"
	"time"

	"go-common/app/job/main/figure-timer/conf"
	"go-common/app/job/main/figure-timer/model"
	// "go-common/library/log"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	dao *Dao
	ctx = context.TODO()
)

func init() {
	var err error
	flag.Set("conf", "../cmd/figure-timer-job-test.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	dao = New(conf.Conf)
}

func TestInit(t *testing.T) {
	Convey("TEST init", t, func() {
		var err error
		So(dao.c, ShouldNotBeNil)
		So(dao.mysql, ShouldNotBeNil)
		So(dao.hbase, ShouldNotBeNil)
		So(dao.redis, ShouldNotBeNil)
		err = dao.Ping(ctx)

		So(err, ShouldBeNil)
	})
}

func TestMysql(t *testing.T) {
	Convey("TEST figure", t, func() {
		var (
			figure = &model.Figure{
				Mid:             100,
				Score:           100,
				LawfulScore:     50,
				WideScore:       0,
				FriendlyScore:   23,
				BountyScore:     250,
				CreativityScore: 0,
				Ver:             0,
			}
			figure2    *model.Figure
			figureRank = &model.Rank{
				ScoreFrom:  2000,
				ScoreTo:    3000,
				Percentage: 50,
				Ver:        233,
			}
			err error
		)
		figure.ID, err = dao.UpsertFigure(ctx, figure)
		So(err, ShouldBeNil)
		figure2, err = dao.Figure(ctx, figure.Mid)
		So(err, ShouldBeNil)
		So(figure2, ShouldNotBeNil)
		So(figure2.Ctime, ShouldNotBeEmpty)
		So(figure2.Ctime, ShouldHappenBefore, time.Now().Add(time.Second))
		So(figure2.Mtime, ShouldNotBeEmpty)
		So(figure2.Mtime, ShouldHappenBefore, time.Now().Add(time.Second))
		_, err = dao.InsertRankHistory(ctx, figureRank)
		So(err, ShouldBeNil)
		_, err = dao.UpsertRank(ctx, figureRank)
		So(err, ShouldBeNil)
	})
	Convey("TEST figures", t, func() {
		for shard := 0; shard < 100; shard++ {
			var (
				end     bool
				fromMid = int64(shard)
				figures []*model.Figure
				err     error
			)
			for !end {
				figures, end, err = dao.Figures(ctx, fromMid, 100)
				So(err, ShouldBeNil)
				for _, f := range figures {
					So(f.Mid, ShouldBeGreaterThan, 0)
					if fromMid < f.Mid {
						fromMid = f.Mid
					}
				}
			}
		}
	})
}

func TestRedis(t *testing.T) {
	Convey("TEST figure", t, func() {
		var (
			figure = &model.Figure{
				Mid:             23,
				Score:           100,
				LawfulScore:     50,
				WideScore:       0,
				FriendlyScore:   23,
				BountyScore:     250,
				CreativityScore: 0,
				Ver:             0,
			}
			figure2 *model.Figure
			err     error
		)
		err = dao.SetFigureCache(ctx, figure)
		So(err, ShouldBeNil)
		figure2, err = dao.FigureCache(ctx, figure.Mid)
		So(err, ShouldBeNil)
		So(figure2, ShouldNotBeNil)
		So(figure, ShouldResemble, figure2)
	})
	Convey("TEST pending mids", t, func() {
		var (
			mid   int64 = 23
			ver   int64
			shard = mid % dao.c.Property.PendingMidShard
			mids  []int64
			err   error
		)
		err = dao.setPendingMidCache(ctx, mid, ver)
		So(err, ShouldBeNil)
		mids, err = dao.PendingMidsCache(ctx, ver, shard)
		So(err, ShouldBeNil)
		So(mids, ShouldNotBeEmpty)
	})
}

func (d *Dao) setPendingMidCache(c context.Context, mid int64, ver int64) (err error) {
	var (
		key  = keyPendingMids(ver, mid%d.c.Property.PendingMidShard)
		conn = d.redis.Get(c)
	)
	defer conn.Close()
	if err = conn.Send("SADD", key, mid); err != nil {
		err = errors.Wrapf(err, "conn.Send(SADD,%s,%d)", key, mid)
		return
	}
	if err = conn.Send("EXPIRE", key, d); err != nil {
		err = errors.Wrapf(err, "conn.Send(EXPIRE,%s,%v)", key, d)
		return
	}
	if err = conn.Flush(); err != nil {
		err = errors.WithStack(err)
		return
	}
	if _, err = conn.Receive(); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// func testHbase(t *testing.T) {
// 	Convey("TEST figure_activity", t, func() {
// 		var (
// 			mid          int64 = 15555180
// 			err          error
// 			weekVer      = time.Date(2017, 10, 2, 0, 0, 0, 0, time.Local).Unix()
// 			weekVerFrom  = time.Date(2016, 10, 3, 0, 0, 0, 0, time.Local).Unix()
// 			weekVerTo    = time.Date(2017, 9, 25, 0, 0, 0, 0, time.Local).Unix()
// 			userInfo     *model.UserInfo
// 			figureRecord = &model.FigureRecord{
// 				Mid:            mid,
// 				XPosLawful:     12,
// 				XNegLawful:     2,
// 				XPosWide:       12,
// 				XNegWide:       13,
// 				XPosFriendly:   233,
// 				XNegFriendly:   234,
// 				XPosCreativity: 0,
// 				XNegCreativity: 1,
// 				XPosBounty:     -250,
// 				XNegBounty:     -251,
// 			}
// 			figureRecords []*model.FigureRecord
// 		)
// 		err = dao.putSpyScore(ctx, mid, 100)
// 		So(err, ShouldBeNil)
// 		userInfo, err = dao.UserInfo(ctx, mid, weekVer)
// 		So(err, ShouldBeNil)
// 		So(userInfo, ShouldNotBeNil)

// 		err = dao.PutCalcRecord(ctx, figureRecord, time.Date(2016, 12, 3, 0, 0, 0, 0, time.Local).Unix())
// 		So(err, ShouldBeNil)

// 		figureRecords, err = dao.CalcRecords(ctx, mid, weekVerFrom, weekVerTo+1)
// 		So(err, ShouldBeNil)
// 		So(figureRecords, ShouldNotBeEmpty)
// 	})
// }

// // PutSpyScore add spy score info.
// func (d *Dao) putSpyScore(c context.Context, mid int64, score int8) (err error) {
// 	var (
// 		key         = rowKeyUserInfo(mid)
// 		scoreB      = make([]byte, 2)
// 		ctx, cancel = context.WithTimeout(c, time.Duration(conf.Conf.Hbase.WriteTimeout))
// 	)
// 	defer cancel()
// 	binary.BigEndian.PutUint16(scoreB, uint16(score))
// 	values := map[string]map[string][]byte{_hbaseUserFC: map[string][]byte{_hbaseUserQSpy: scoreB}}
// 	if _, err = d.hbase.PutStr(ctx, _hbaseUserTable, key, values); err != nil {
// 		log.Error("hbase.PutStr(%s, %s, %v) error(%v)", _hbaseUserTable, key, values, err)
// 		return
// 	}
// 	return
// }
