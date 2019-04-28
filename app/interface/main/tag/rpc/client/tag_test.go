package tag

import (
	"context"
	"testing"

	"go-common/app/interface/main/tag/model"
	"go-common/library/net/rpc"
)

func TestTag(t *testing.T) {
	s := &Service{}
	s.client = rpc.NewDiscoveryCli("", nil)
	testInfoByID(t, s)
	testInfoByName(t, s)
	testInfosByIDs(t, s)
	testInfosByNames(t, s)
	testArcTags(t, s)
	testSubTags(t, s)
}

func testInfoByID(t *testing.T, svr *Service) {
	arg := &model.ArgID{ID: 10031, Mid: 27515274}
	tag, err := svr.InfoByID(context.TODO(), arg)
	if err != nil {
		t.Error("err:", err)
		t.Fail()
		return
	}
	t.Log("tag:", tag)
}

func testInfoByName(t *testing.T, svr *Service) {
	arg := &model.ArgName{Name: "朱杰测试8", Mid: 27515274}
	tag, err := svr.InfoByName(context.TODO(), arg)
	if err != nil {
		t.Error("err:", err)
		t.Fail()
		return

	}
	t.Log("tag:", tag)
}

func testInfosByIDs(t *testing.T, svr *Service) {
	arg := &model.ArgIDs{IDs: []int64{10031, 2}, Mid: 27515274}
	tags, err := svr.InfoByIDs(context.TODO(), arg)
	if err != nil {
		t.Error("err:", err)
		t.Fail()
		return
	}
	t.Log("tags:", tags)
}

func testInfosByNames(t *testing.T, svr *Service) {
	arg2 := &model.ArgNames{Names: []string{"朱杰测试8", "2012"}, Mid: 27515274}
	tags, err := svr.InfoByNames(context.TODO(), arg2)
	if err != nil {
		t.Error("err:", err)
		t.Fail()
		return
	}
	t.Log("tags:", tags)

}

func testArcTags(t *testing.T, svr *Service) {
	arg := &model.ArgAid{Aid: 4053003, Mid: 27515274}
	tags, err := svr.ArcTags(context.TODO(), arg)
	if err != nil {
		t.Error("err:", err)
		t.Fail()
		return
	}
	t.Log("tags:", tags)
}

func testSubTags(t *testing.T, svr *Service) {
	arg := &model.ArgSub{Mid: 15555180, Pn: 1, Ps: 20, Order: -1}
	sub, err := svr.SubTags(context.TODO(), arg)
	if err != nil {
		t.Error("err:", err)
		t.Fail()
		return
	}
	t.Log("sub:", sub)
}
