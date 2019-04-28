package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCleanURL(t *testing.T) {
	Convey("different urls", t, func() {
		url := "/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"
		So(CleanURL("http://i0.hdslb.com/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"), ShouldEqual, url)
		So(CleanURL("https://i0.hdslb.com/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"), ShouldEqual, url)
		So(CleanURL("//i0.hdslb.com/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"), ShouldEqual, url)

		So(CleanURL("http://uat-i0.hdslb.com/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"), ShouldEqual, url)
		So(CleanURL("https://uat-i0.hdslb.com/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"), ShouldEqual, url)
		So(CleanURL("//uat-i0.hdslb.com/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"), ShouldEqual, url)
	})
}

func TestCompleteURL(t *testing.T) {
	Convey("prefix http", t, func() {
		url := "http://i0.hdslb.com/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"
		So(CompleteURL(url), ShouldEqual, url)
	})
	Convey("prefix https", t, func() {
		url := "https://i0.hdslb.com/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"
		So(CompleteURL(url), ShouldEqual, url)
	})
	Convey("prefix //", t, func() {
		url := "//i0.hdslb.com/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"
		So(CompleteURL(url), ShouldEqual, url)
	})
	Convey("prefix no prefix", t, func() {
		url := "/bfs/article/3b8b57821a94af4c171218a01d9cee44dc0dccf8.jpg"
		So(CompleteURL(url), ShouldEqual, "https://i0.hdslb.com"+url)
	})
}

func TestNoticeState(t *testing.T) {
	Convey("ToInt64", t, func() {
		a := NoticeState{"lead": true, "new": true}
		So(a.ToInt64(), ShouldEqual, 3)
	})
	Convey("ToInt64 2", t, func() {
		a := NoticeState{"lead": true, "new": false}
		So(a.ToInt64(), ShouldEqual, 1)
	})
	Convey("NewNoticeState", t, func() {
		a := NoticeState{"lead": true, "new": true}
		So(NewNoticeState(3), ShouldResemble, a)
	})
	Convey("NewNoticeState 2", t, func() {
		a := NoticeState{"lead": true, "new": false}
		So(NewNoticeState(1), ShouldResemble, a)
	})
}
