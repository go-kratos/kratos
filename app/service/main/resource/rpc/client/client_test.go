package client

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/resource/model"
)

func TestResource(t *testing.T) {
	s := New(nil)
	time.Sleep(2 * time.Second)
	testResourceAll(t, s)
	testAssignmentAll(t, s)
	testDefBanner(t, s)
	testResource(t, s)
	testResources(t, s)
	testBanners(t, s)
	testPasterAPP(t, s)
	testIndexIcon(t, s)
	testPlayerIcon(t, s)
	testCmtbox(t, s)
	testSideBars(t, s)
	testAbTest(t, s)
	testPasterCID(t, s)
}

// testResourceAll test rpc ResourceAll.
func testResourceAll(t *testing.T, s *Service) {
	res, err := s.ResourceAll(context.TODO())
	if err != nil {
		t.Logf("testResourceAll error(%v) \n", err)
		return
	}
	t.Logf("testResourceAll res: %+v \n", res)
}

// testAssignmentAll test rpc AssignmentAll.
func testAssignmentAll(t *testing.T, s *Service) {
	res, err := s.AssignmentAll(context.TODO())
	if err != nil {
		t.Logf("testAssignmentAll error(%v) \n", err)
		return
	}
	t.Logf("testAssignmentAll res: %+v \n", res)
}

// testDefBanner test rpc DefBanner.
func testDefBanner(t *testing.T, s *Service) {
	res, err := s.DefBanner(context.TODO())
	if err != nil {
		t.Logf("testDefBanner error(%v) \n", err)
		return
	}
	t.Logf("testDefBanner res: %+v \n", res)
}

// testResource test rpc Resource.
func testResource(t *testing.T, s *Service) {
	p := &model.ArgRes{
		ResID: 1233,
	}
	res, err := s.Resource(context.TODO(), p)
	if err != nil {
		t.Logf("testResource error(%v) \n", err)
		return
	}
	t.Logf("testResource res: %+v \n", res)
}

// testResources test rpc Resources.
func testResources(t *testing.T, s *Service) {
	p := &model.ArgRess{
		ResIDs: []int{1187, 1639},
	}
	res, err := s.Resources(context.TODO(), p)
	if err != nil {
		t.Logf("testResources error(%v) \n", err)
		return
	}
	t.Logf("testResources res: %+v \n", res)
}

// testBanners test rpc Banners.
func testBanners(t *testing.T, s *Service) {
	ab := &model.ArgBanner{
		Plat:      1,
		ResIDs:    "454,467",
		Build:     508000,
		MID:       1493031,
		Channel:   "abc",
		IP:        "211.139.80.6",
		Buvid:     "123",
		Network:   "wifi",
		MobiApp:   "iphone",
		Device:    "test",
		IsAd:      true,
		OpenEvent: "abc",
	}
	res, err := s.Banners(context.TODO(), ab)
	if err != nil {
		t.Logf("testBanners error(%v) \n", err)
		return
	}
	t.Logf("testBanners res: %+v \n", res)
}

// testPasterAPP test rpc Paster.
func testPasterAPP(t *testing.T, s *Service) {
	p := &model.ArgPaster{
		Platform: int8(4),
		AdType:   int8(1),
		Aid:      "10097274",
		TypeId:   "11",
		Buvid:    "666666",
	}
	res, err := s.PasterAPP(context.TODO(), p)
	if err != nil {
		t.Logf("testPasterAPP error(%v) \n", err)
		return
	}
	t.Logf("testPaster res: %+v \n", res)
}

// testIndexIcon test rpc IndexIcon.
func testIndexIcon(t *testing.T, s *Service) {
	res, err := s.IndexIcon(context.TODO())
	if err != nil {
		t.Logf("testIndexIcon error(%v) \n", err)
		return
	}
	t.Logf("testIndexIcon res: %+v \n", res)
}

// testPlayerIcon test rpc PlayerIcon.
func testPlayerIcon(t *testing.T, s *Service) {
	res, err := s.PlayerIcon(context.TODO())
	if err != nil {
		t.Logf("testPlayerIcon error(%v) \n", err)
		return
	}
	t.Logf("testPlayerIcon res: %+v \n", res)
}

// testCmtbox test rpc Resource.
func testCmtbox(t *testing.T, s *Service) {
	p := &model.ArgCmtbox{
		ID: 1,
	}
	res, err := s.Cmtbox(context.TODO(), p)
	if err != nil {
		t.Logf("testCmtbox error(%v) \n", err)
		return
	}
	t.Logf("testCmtbox res: %+v \n", res)
}

// testSideBars test rpc SideBars.
func testSideBars(t *testing.T, s *Service) {
	res, err := s.SideBars(context.TODO())
	if err != nil {
		t.Logf("testSideBars error(%v) \n", err)
		return
	}
	t.Logf("testSideBars res: %+v \n", res)
}

// testAbTest test rpc abtest.
func testAbTest(t *testing.T, s *Service) {
	p := &model.ArgAbTest{
		Groups: "不显示热门tab,显示热门tab",
		IP:     "127.0.0.1",
	}
	res, err := s.AbTest(context.TODO(), p)
	if err != nil {
		t.Logf("testAbTest error(%v) \n", err)
		return
	}
	t.Logf("testAbTest res: %+v \n", res)
}

// testPasterCID test rpc PasterCID
func testPasterCID(t *testing.T, s *Service) {
	res, err := s.PasterCID(context.TODO())
	if err != nil {
		t.Logf("testPasterCID error(%v) \n", err)
		return
	}
	t.Logf("testPasterCID res: %+v \n", res)
}
