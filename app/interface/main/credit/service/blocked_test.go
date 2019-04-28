package service

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestBatchBLKInfos .
func TestBatchBLKInfos(t *testing.T) {
	Convey("TestBatchBLKInfos", t, func() {
		res, err := s.BatchBLKInfos(context.TODO(), []int64{1111, 222})
		So(err, ShouldBeNil)
		out, err := json.Marshal(res)
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
	})
}

// TestBlockedUserCard .
func TestBlockedUserCard(t *testing.T) {
	Convey("TestBlockedUserCard", t, func() {
		res, err := s.BlockedUserCard(context.TODO(), 21432418)
		So(err, ShouldBeNil)
		out, err := json.Marshal(res)
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
	})
}

// TestBlockedUserList .
func TestBlockedUserList(t *testing.T) {
	Convey("TestBlockedUserList", t, func() {
		res, err := s.BlockedUserList(context.TODO(), 1)
		So(err, ShouldBeNil)
		out, err := json.Marshal(res)
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
	})
}

// TestBlockedInfo .
func TestBlockedInfo(t *testing.T) {
	Convey("TestBlockedInfo", t, func() {
		res, err := s.BlockedInfo(context.TODO(), 1475)
		So(err, ShouldBeNil)
		out, err := json.Marshal(res)
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
	})
}

// TestBlockedList .
func TestBlockedList(t *testing.T) {
	Convey("TestBlockedList", t, func() {
		res, err := s.BlockedList(context.TODO(), 0, -1, 1, 5)
		So(err, ShouldBeNil)
		out, err := json.Marshal(res)
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
	})
}

// TestAnnouncementInfo .
func TestAnnouncementInfo(t *testing.T) {
	Convey("TestAnnouncementInfo", t, func() {
		s.LoadAnnouncement(context.TODO())
		res, err := s.AnnouncementInfo(context.TODO(), 48)
		So(err, ShouldBeNil)
		out, err := json.Marshal(res)
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
	})
}

// TestAnnouncementList .
func TestAnnouncementList(t *testing.T) {
	Convey("TestAnnouncementList", t, func() {
		s.LoadAnnouncement(context.TODO())
		res, err := s.AnnouncementList(context.TODO(), 1, 1, 3)
		So(err, ShouldBeNil)
		out, err := json.Marshal(res)
		So(err, ShouldBeNil)
		So(out, ShouldNotBeNil)
	})
}
