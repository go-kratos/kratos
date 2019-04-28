package reply

import (
	"context"
	"go-common/library/database/sql"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestReplyNewEmojiDao(t *testing.T) {
	convey.Convey("NewEmojiDao", t, func(ctx convey.C) {
		var (
			db = &sql.DB{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			dao := NewEmojiDao(db)
			ctx.Convey("Then dao should not be nil.", func(ctx convey.C) {
				ctx.So(dao, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyEmojiList(t *testing.T) {
	convey.Convey("EmojiList", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			emo, err := d.Emoji.EmojiList(c)
			ctx.Convey("Then err should be nil.emo should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(emo, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyEmojiListByPid(t *testing.T) {
	convey.Convey("EmojiListByPid", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			pid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			emo, err := d.Emoji.EmojiListByPid(c, pid)
			ctx.Convey("Then err should be nil.emo should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(emo, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestReplyListEmojiPack(t *testing.T) {
	convey.Convey("ListEmojiPack", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			packs, err := d.Emoji.ListEmojiPack(c)
			ctx.Convey("Then err should be nil.packs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(packs, convey.ShouldNotBeNil)
			})
		})
	})
}
