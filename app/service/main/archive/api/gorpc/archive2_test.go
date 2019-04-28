package gorpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
)

func TestDynamic(t *testing.T) {
	s := New2(nil)
	time.Sleep(5 * time.Second)
	testArchivesWithPlayer(t, s)
	testMaxAID(t, s)
	testVideo3(t, s)
	testUpsArcs3(t, s)
	testUpArcs3(t, s)
	testRecommend3(t, s)
	testArchive3(t, s)
	testArchives3(t, s)
	testAidPTime(t, s)
	testTypes2(t, s)
	testVideoshot2(t, s)
	testUpCount2(t, s)
	testView3(t, s)
	testViews3(t, s)
	testClick3(t, s)
	testRankArcs3(t, s)
	testRanksArcs3(t, s)
	testRankTopArcs3(t, s)
	testRankAllArcs3(t, s)
	testPage3(t, s)
	testRanksTopCount2(t, s)
	testStat3(t, s)
	testStats3(t, s)
	testArcCache2(t, s)
	testArcFieldCache2(t, s)
	testSetStat2(t, s)
	testUpVideo2(t, s)
	testDelVideo2(t, s)
	testDescription2(t, s)
}

func testViews3(t *testing.T, s *Service2) {
	var vs, err = s.Views3(context.TODO(), &archive.ArgAids2{Aids: []int64{10097450, 10097454}})
	if err != nil {
		t.Log(err)
		return
	}
	for _, v := range vs {
		t.Log(v.Archive3)
	}
}
func testClick3(t *testing.T, s *Service2) {
	t.Log(s.Click3(context.TODO(), &archive.ArgAid2{Aid: 10097454}))
}
func testRankArcs3(t *testing.T, s *Service2) {
	t.Log(s.RankArcs3(context.TODO(), &archive.ArgRank2{Rid: 20, Type: 0, Ps: 20, Pn: 1}))
}
func testRanksArcs3(t *testing.T, s *Service2) {
	t.Log(s.RanksArcs3(context.TODO(), &archive.ArgRanks2{Rids: []int16{20, 24}, Type: 0, Ps: 20, Pn: 1}))
}
func testRankTopArcs3(t *testing.T, s *Service2) {
	t.Log(s.RankTopArcs3(context.TODO(), &archive.ArgRankTop2{ReID: 1, Pn: 1, Ps: 10}))
}
func testRankAllArcs3(t *testing.T, s *Service2) {
	am, err := s.RankAllArcs3(context.TODO(), &archive.ArgRankAll2{Pn: 1, Ps: 40})
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("allRankTop(%d)", am.Count)
	for _, arc := range am.Archives {
		t.Logf("archive(%+v)", arc)
	}
}

func testAidPTime(t *testing.T, s *Service2) {
	res, err := s.UpsPassed2(context.TODO(), &archive.ArgUpsArcs2{Mids: []int64{15555180, 27515256}})
	fmt.Println(len(res))
	if err != nil {
		t.Logf("s.UpsPassed2 error(%v)", err)
		return
	}
	for mid, v := range res {
		t.Logf("mid:%d", mid)
		for _, aidAndPtime := range v {
			t.Logf("mid:%d, aid:%d, ptime:%d, copyright:%d", mid, aidAndPtime.Aid, aidAndPtime.PubDate, aidAndPtime.Copyright)
		}
	}
}

func testArchivesWithPlayer(t *testing.T, s *Service2) {
	arg := &archive.ArgPlayer{Aids: []int64{10097272, 1}, Qn: 16, Platform: "html5", RealIP: "222.73.196.18"}
	res, err := s.ArchivesWithPlayer(context.TODO(), arg)
	if err != nil {
		t.Logf("s.ArchivesWithPlayer(%+v) error(%v)", arg, err)
		return
	}
	for k, v := range res {
		t.Logf("aid:%d, arc:%+v, player:%+v", k, v, v.PlayerInfo)
	}
}

func testVideo3(t *testing.T, s *Service2) {
	t.Log(s.Video3(context.TODO(), &archive.ArgVideo2{Aid: 10097760, Cid: 10107336}))
}

func testRecommend3(t *testing.T, s *Service2) {
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		res, err := s.Recommend3(context.TODO(), &archive.ArgAid2{Aid: 4052562})
		if err != nil {
			t.Logf("s.Recommend3 error(%v)", err)
			continue
		}
		t.Log(res)
	}
}

func testArchive3(t *testing.T, s *Service2) {
	res, err := s.Archive3(context.TODO(), &archive.ArgAid2{Aid: 5463609})
	if err != nil {
		t.Logf("s.Archive3 error(%v)", err)
		return
	}
	t.Log(res)
}

func testArchives3(t *testing.T, s *Service2) {
	res, err := s.Archives3(context.TODO(), &archive.ArgAids2{Aids: []int64{5463609, 5463608, 10097657}})
	if err != nil {
		t.Logf("s.Archives3 error(%v)", err)
		return
	}
	for _, a := range res {
		t.Log(a)
	}
}

func testMaxAID(t *testing.T, s *Service2) {
	res, err := s.MaxAID(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("maxAID is %d", res)
}

func testTypes2(t *testing.T, s *Service2) {
	res, err := s.Types2(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	for _, tp := range res {
		t.Log(tp)
	}
}

func testVideoshot2(t *testing.T, s *Service2) {
	res, err := s.Videoshot2(context.TODO(), &archive.ArgCid2{Cid: 10108203})
	t.Log(res, err)
}

func testUpArcs3(t *testing.T, s *Service2) {
	res, err := s.UpArcs3(context.TODO(), &archive.ArgUpArcs2{Mid: 27515615, Pn: 1, Ps: 20})
	if err != nil {
		t.Error(err)
		return
	}
	for _, a := range res {
		t.Log(a)
	}
}

func testUpCount2(t *testing.T, s *Service2) {
	count, err := s.UpCount2(context.TODO(), &archive.ArgUpCount2{Mid: 27515615})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(count)
}

func testUpsArcs3(t *testing.T, s *Service2) {
	res, err := s.UpsArcs3(context.TODO(), &archive.ArgUpsArcs2{Mids: []int64{27515615, 1684013}})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func testView3(t *testing.T, s *Service2) {
	t.Log(s.View3(context.TODO(), &archive.ArgAid2{Aid: 10098067}))
}

func testPage3(t *testing.T, s *Service2) {
	t.Log(s.Page3(context.TODO(), &archive.ArgAid2{Aid: 10097454}))
}

func testRanksTopCount2(t *testing.T, s *Service2) {
	t.Log(s.RanksTopCount2(context.TODO(), &archive.ArgRankTopsCount2{ReIDs: []int16{1, 2}}))
}

func testStat3(t *testing.T, s *Service2) {
	res, err := s.Stat3(context.TODO(), &archive.ArgAid2{Aid: 10989901})
	if err != nil {
		t.Errorf("testStats3 error(%v)", err)
		return
	}
	t.Log(res, 1)
}

func testStats3(t *testing.T, s *Service2) {
	res, err := s.Stats3(context.TODO(), &archive.ArgAids2{Aids: []int64{10989901, 10097453}})
	if err != nil {
		t.Logf("testStats3 error(%v)", err)
		return
	}
	for aid, r := range res {
		t.Logf("%d:%+v\n", aid, r)
	}
}

func testArcCache2(t *testing.T, s *Service2) {
	t.Log(s.ArcCache2(context.TODO(), &archive.ArgCache2{Aid: 10097454, Tp: archive.CacheUpdate}))
}

func testArcFieldCache2(t *testing.T, s *Service2) {
	t.Log(s.ArcFieldCache2(context.TODO(), &archive.ArgFieldCache2{Aid: 10097454, TypeID: 20, OldTypeID: 10}))
}

func testSetStat2(t *testing.T, s *Service2) {
	t.Log(s.SetStat2(context.TODO(), &api.Stat{Aid: 10097454, View: 0, Danmaku: 0, Reply: 10, Fav: 10, Coin: 10, Share: 10, HisRank: 10, NowRank: 10}))
}

func testUpVideo2(t *testing.T, s *Service2) {
	t.Log(s.UpVideo2(context.TODO(), &archive.ArgVideo2{Aid: 10097760, Cid: 10107336}))
}

func testDelVideo2(t *testing.T, s *Service2) {
	t.Log(s.DelVideo2(context.TODO(), &archive.ArgVideo2{Aid: 10097760, Cid: 10107336}))
}

func testDescription2(t *testing.T, s *Service2) {
	t.Log(s.Description2(context.TODO(), &archive.ArgAid{Aid: 10097454}))
}
