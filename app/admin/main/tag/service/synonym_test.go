package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestServiceSynonym(t *testing.T) {
	var (
		tid int64 = 2233
		pn  int32 = 1
		ps  int32 = 20

		adverb = []int64{22, 33}

		keyWord = "test"
		uname   = "shunza"
		tname   = "unit test"
	)
	Convey("SynonymList", func() {
		testSvc.SynonymList(context.TODO(), keyWord, pn, ps)
	})
	Convey("SynonymAdd", func() {
		testSvc.SynonymAdd(context.TODO(), uname, tname, adverb)
	})
	Convey("SynonymInfo", func() {
		testSvc.SynonymInfo(context.TODO(), tid)
	})
	Convey("SynonymDelete", func() {
		testSvc.SynonymDelete(context.TODO(), tid)
	})
	Convey("RemoveSynonymSon", func() {
		testSvc.RemoveSynonymSon(context.TODO(), tid, adverb)
	})
	Convey("SynonymIsExist", func() {
		testSvc.SynonymIsExist(context.TODO(), tname)
	})
}
