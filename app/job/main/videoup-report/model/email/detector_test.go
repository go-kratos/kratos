package email

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"time"
)

var fd *FastDetector
var speedThreshold = 3
var overspeedThreshold = 2
var overlimit = speedThreshold * overspeedThreshold

func TestEmailnewFastDetector(t *testing.T) {
	convey.Convey("newFastDetector", t, func(ctx convey.C) {
		fd = NewFastDetector(speedThreshold, overspeedThreshold)
		ctx.Convey("Then fd should not be nil.", func(ctx convey.C) {
			ctx.So(fd, convey.ShouldNotBeNil)
		})
	})
}

//一秒多少个
func limit(t *testing.T, ids []int64) {
	now := time.Now().UnixNano()
	limit := len(ids)
	for i := 0; i < limit; i++ {
		unique := ids[i]
		res := fd.Detect(unique)
		hit := fd.IsFastUnique(unique)
		t.Logf("index=%d, unique=%d, res=%v, hit=%v, detector(%+v)", i, unique, res, hit, fd)
	}
	if diff := now + 1e9 - time.Now().UnixNano(); diff > 0 {
		time.Sleep(time.Duration(diff))
	}
	return
}

//分片
func batch(tk []int64, length int) (path [][]int64) {
	ll := len(tk) / length
	if len(tk)%length > 0 {
		ll++
	}
	path = [][]int64{}
	item := []int64{}
	for i := 0; i < len(tk); i++ {
		if i > 0 && i%length == 0 {
			path = append(path, item)
			item = []int64{}
		}
		item = append(item, tk[i])
	}
	if len(item) > 0 {
		path = append(path, item)
	}

	return
}

func TestEmaildetect(t *testing.T) {
	TestEmailnewFastDetector(t)

	ids := []int64{222, 333}

	//convey.Convey("慢速，无超限名额", t, func() {
	//	tk := []int64{}
	//	for i := 1; i <= overlimit*2; i++ {
	//		tk = append(tk, ids[0])
	//	}
	//	path := batch(tk, speedThreshold-1)
	//	for _, tk := range path {
	//		limit(t, tk)
	//	}
	//})
	//return

	//convey.Convey("部分超限，但未超次数，无超限名额", t, func() {
	//	tk := []int64{}
	//	for i := 1; i < overlimit; i++ {
	//		tk = append(tk, ids[0])
	//	}
	//	path := [][]int64{tk}
	//	path = append(path, batch(tk, speedThreshold-1)...)
	//	for _, tk := range path {
	//		limit(t, tk)
	//	}
	//})
	//return

	convey.Convey("部分超限且超次数，没有回落,有超限名额且被替代, 几秒后再超限，但没有超限名单", t, func() {
		tk := []int64{}
		tk2 := []int64{}
		for i := 1; i < overlimit*2; i++ {
			tk = append(tk, ids[0])
			tk2 = append(tk2, ids[1])
		}
		path := [][]int64{tk}
		path = append(path, batch(tk, speedThreshold+1)...)
		path = append(path, batch(tk2, speedThreshold+1)...)
		for _, tk := range path {
			limit(t, tk)
		}
		limit(t, path[len(path)-1])

		sl := int64(5)
		time.Sleep(time.Duration(sl) * time.Second)
		t.Logf("after %ds sleep", sl)
		limit(t, path[len(path)-1])
	})

	convey.Convey("部分超限且超次数，有超限名额，但有回落，没有超限名额", t, func(ctx convey.C) {
		tk := []int64{}
		for i := 1; i < overlimit*2; i++ {
			tk = append(tk, ids[0])
		}
		path := batch(tk, speedThreshold+1)
		path = append(path, batch(tk, speedThreshold-1)...)
		for _, tk := range path {
			limit(t, tk)
		}
	})

}
