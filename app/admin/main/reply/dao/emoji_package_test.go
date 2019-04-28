package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_EmojiPack(t *testing.T) {
	Convey("CreateEmojiPackage", t, WithDao(func(d *Dao) {
		id, err := d.CreateEmojiPackage(context.Background(), "[2233娘]", "www.baidu.com", 0, "", 1)
		So(err, ShouldBeNil)
		So(id, ShouldNotEqual, 0)
		d.DelEmojiPackage(context.Background(), id)
	}))

	Convey("EmojiPackageList", t, WithDao(func(d *Dao) {
		packs, err := d.EmojiPackageList(context.Background())
		So(err, ShouldBeNil)
		for _, v := range packs {
			t.Logf("v.Id= %d, v.Name= %s, v.Url= %s, v.Remark= %s, v.State= %d, v.Sort= %d",
				v.ID, v.Name, v.URL, v.Remark, v.State, v.Sort)
		}
	}))

	Convey("UpEmojiPackage", t, WithDao(func(d *Dao) {
		id, err := d.UpEmojiPackage(context.Background(), "[小电视x]", "xxxx", "xx", 1, 1)
		So(err, ShouldBeNil)
		t.Logf("id= %d", id)
	}))

	Convey("UpEmojiPackageSort", t, WithDao(func(d *Dao) {
		tx, _ := d.BeginTran(context.Background())
		err := d.UpEmojiPackageSort(tx, "1")
		if err != nil {
			tx.Rollback()
			t.Errorf("UpEmojiPackageSort err (%v)", err)
			return
		}
		tx.Commit()
	}))
}
