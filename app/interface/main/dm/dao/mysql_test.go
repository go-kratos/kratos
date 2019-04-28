package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm/model"
)

func TestSubject(t *testing.T) {
	var (
		c         = context.TODO()
		tp  int32 = 1
		oid int64 = 1221
	)
	s, err := testDao.Subject(c, tp, oid)
	if err != nil {
		t.Error(err)
	}
	if s == nil {
		t.Logf("oid:%d subject not exist", oid)
	} else {
		t.Logf("subject:%+v", s)
	}
}

func TestIndex(t *testing.T) {
	var (
		c          = context.TODO()
		tp   int32 = 1
		oid  int64 = 1221
		dmid int64 = 719150137
	)
	dm, err := testDao.Index(c, tp, oid, dmid)
	if err != nil {
		t.Error(err)
	}
	if dm == nil {
		t.Logf("dmid:%d not exist", dmid)
	} else {
		t.Logf("dm:%+v", dm)
	}
}

func TestIndexsByID(t *testing.T) {
	var (
		c           = context.TODO()
		tp    int32 = 1
		oid   int64 = 1221
		dmids       = []int64{719150137, 719150230, 719150141}
	)
	res, special, err := testDao.IndexsByID(c, tp, oid, dmids)
	if err != nil {
		t.Error(err)
	}
	for _, dm := range res {
		t.Logf("dm:%+v", dm)
	}
	t.Logf("special:%+v", special)
}

func TestContent(t *testing.T) {
	var (
		c          = context.TODO()
		oid  int64 = 1221
		dmid int64 = 719150137
	)
	ct, err := testDao.Content(c, oid, dmid)
	if err != nil {
		t.Error(err)
	}
	if ct == nil {
		t.Logf("content:%d not exist", dmid)
	} else {
		t.Logf("content:%+v", ct)
	}
}

func TestContents(t *testing.T) {
	var (
		c           = context.TODO()
		oid   int64 = 1221
		dmids       = []int64{719150137, 719150230, 719150141}
	)
	res, err := testDao.Contents(c, oid, dmids)
	if err != nil {
		t.Error(err)
	}
	for _, ct := range res {
		t.Logf("content:%+v", ct)
	}
}

func TestContentSpecial(t *testing.T) {
	var (
		c          = context.TODO()
		dmid int64 = 719150141
	)
	cs, err := testDao.ContentSpecial(c, dmid)
	if err != nil {
		t.Error(err)
	}
	t.Logf("content special:%+v", cs)
}

func TestContentsSpecial(t *testing.T) {
	var (
		c     = context.TODO()
		dmids = []int64{719150141, 719150141}
	)
	res, err := testDao.ContentsSpecial(c, dmids)
	if err != nil {
		t.Error(err)
	}
	for _, cs := range res {
		t.Logf("content special:%+v", cs)
	}
}

func TestCheckTransferJob(t *testing.T) {
	var (
		from, to int64 = 10108765, 10108763
	)
	job, err := testDao.CheckTransferJob(context.TODO(), from, to)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(job)
}

func TestAddTransferJob(t *testing.T) {
	var (
		from, to, mid int64 = 10108765, 10108763, 27515615
		offset              = 1.00
		state               = model.TransferJobStatInit
	)
	_, err := testDao.AddTransferJob(context.TODO(), from, to, mid, offset, state)
	if err != nil {
		t.Fatal(err)
	}
}
