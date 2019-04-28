package dao

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testAddMidRedis int64 = 15555180
	testLowScore    int8  = 7
	testBlockNo     int64 = 100000
	testKey               = "test_lock"
	testMidList           = []int64{15555180}
)

func Test_PingRedis(t *testing.T) {
	Convey("ping redis", t, func() {
		So(d.PingRedis(context.TODO()), ShouldBeNil)
	})
}

func Test_AddBlockCache(t *testing.T) {
	Convey("AddBlockCache redis ", t, func() {
		c := context.TODO()
		ret, err := d.BlockMidCache(c, testBlockNo, 10)
		So(ret, ShouldBeEmpty)
		So(err, ShouldBeNil)

		err = d.AddBlockCache(c, testAddMidRedis, testLowScore, testBlockNo)
		So(err, ShouldBeNil)

		ret, err = d.BlockMidCache(c, testBlockNo, 10)
		So(ret, ShouldContain, testAddMidRedis)
		So(err, ShouldBeNil)

		err = d.DelBlockCache(c, testBlockNo, testAddMidRedis)
		So(err, ShouldBeNil)

		ret, err = d.BlockMidCache(c, testBlockNo, 10)
		So(ret, ShouldNotContain, testAddMidRedis)
		So(err, ShouldBeNil)

	})
}

func Test_SetNXLockCache(t *testing.T) {
	Convey("SetNXLockCache", t, func() {
		ret, err := d.SetNXLockCache(c, testKey, 2)
		So(err, ShouldBeNil)
		So(ret, ShouldBeTrue)

		err = d.DelLockCache(c, testKey)
		So(err, ShouldBeNil)

		ret, err = d.SetNXLockCache(c, testKey, 2)
		So(err, ShouldBeNil)
		So(ret, ShouldBeTrue)

		ret, err = d.SetNXLockCache(c, testKey, 2)
		So(err, ShouldBeNil)
		So(ret, ShouldBeFalse)
	})

}

func Test_SetBlockCache(t *testing.T) {
	Convey("ping SetBlockCache", t, func() {
		err := d.SetBlockCache(context.TODO(), testMidList)
		So(err, ShouldBeNil)
		Convey("ping SetBlockCache", func() {
			mid, err := d.SPOPBlockCache(context.TODO())
			So(err, ShouldBeNil)
			So(testMidList, ShouldContain, mid)
			Convey("ping SetBlockCache 2", func() {
				mid, err := d.SPOPBlockCache(context.TODO())
				So(err, ShouldBeNil)
				So(mid == 0, ShouldBeTrue)
			})
		})
	})
}

// go test  -test.v -test.run TestPfaddCache
func TestPfaddCache(t *testing.T) {
	Convey("PfaddCache", t, func() {
		idx := strings.Replace(randHex(), "-", "", -1)
		ok, err := d.PfaddCache(context.TODO(), idx)
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)
		ok, err = d.PfaddCache(context.TODO(), idx)
		So(err, ShouldBeNil)
		So(ok, ShouldBeFalse)
	})
}

// go test  -test.v -test.run TestBatchPfaddCache
func TestBatchPfaddCache(t *testing.T) {
	Convey("PfaddCache", t, func() {
		for i := 0; i < 10; i++ {
			idx := strings.Replace(randHex(), "-", "", -1)
			_, err := d.PfaddCache(context.TODO(), idx)
			So(err, ShouldBeNil)
		}
	})
}

func randHex() string {
	bs := make([]byte, 16)
	rand.Read(bs)
	return hex.EncodeToString(bs)
}
