package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_Emoji(t *testing.T) {
	// insert emoji
	Convey("CreateEmoji", t, WithDao(func(d *Dao) {
		id, err := d.CreateEmoji(context.Background(), 1, "[小电视_1]", "baidu.com", 0, 0, "ssss")
		So(err, ShouldBeNil)
		So(id, ShouldNotEqual, 0)
		defer d.DelEmojiByID(context.Background(), id)
		insertID := id
		//get all emoji
		Convey("EmojiList", WithDao(func(d *Dao) {
			data, err := d.EmojiList(context.Background())
			So(err, ShouldBeNil)
			for _, v := range data {
				t.Logf("v.Id= %d, v.PackageID= %d, v.Name= %s, v.Url= %s, v.Remark= %s, v.State= %d, v.Sort= %d",
					v.ID, v.PackageID, v.Name, v.URL, v.Remark, v.State, v.Sort)
			}
		}))

		// get emoji by  package_id
		Convey("EmojiListByPid", WithDao(func(d *Dao) {
			data, err := d.EmojiListByPid(context.Background(), 1)
			So(err, ShouldBeNil)
			for _, v := range data {
				t.Logf("v.Id= %d, v.PackageID= %d, v.Name= %s, v.Url= %s, v.Remark= %s, v.State= %d, v.Sort= %d",
					v.ID, v.PackageID, v.Name, v.URL, v.Remark, v.State, v.Sort)
			}
		}))

		//update emoji sort
		Convey("UpEmojiSort", WithDao(func(d *Dao) {
			tx, _ := d.BeginTran(context.Background())
			err := d.UpEmojiSort(tx, "2,1")
			if err != nil {
				tx.Rollback()
				t.Errorf("UpEmojiSort err (%v)", err)
				return
			}
			tx.Commit()
		}))

		//update emoji state
		Convey("test UpdateEmojis", WithDao(func(d *Dao) {
			id, err := d.UpEmojiStateByID(context.Background(), 1, 70)
			So(err, ShouldBeNil)
			So(id, ShouldNotEqual, 0)
			t.Logf("id= %d", id)
		}))

		Convey("test select emoji by name", WithDao(func(d *Dao) {
			emojis, err := d.EmojiByName(context.Background(), "[小电视_1]")
			So(err, ShouldBeNil)
			for _, v := range emojis {
				t.Logf("v.ID= %d", v.ID)
			}
		}))

		// update emoji remark
		Convey("test SortEmojis", WithDao(func(d *Dao) {
			id, err := d.UpEmoji(context.Background(), "[小电视]", "cccxxx", "google.com", insertID)
			So(err, ShouldBeNil)
			So(id, ShouldNotEqual, 0)
			t.Logf("id= %d", id)
		}))

		Convey("test delEmoji", WithDao(func(d *Dao) {
			id, err := d.DelEmojiByID(context.Background(), insertID)
			So(err, ShouldBeNil)
			So(id, ShouldNotEqual, 0)
			t.Logf("id= %d", id)
		}))
	}))

}
