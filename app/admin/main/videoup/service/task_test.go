package service

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/utils"
)

func TestLockVideo(t *testing.T) {
	svr.lockVideo()
}

func TestConsumer(t *testing.T) {
	var c = context.TODO()
	e := svr.HandsUp(c, 10, "10")
	if e != nil {
		t.Fatal(e)
	}

	e = svr.HandsUp(c, 20, "20")
	if e != nil {
		t.Fatal(e)
	}

	e = svr.HandsOff(c, 10, 20)
	if e != nil {
		t.Fatal(e)
	}

	e = svr.HandsUp(c, 20, "20")
	if e != nil {
		t.Fatal(e)
	}

	e = svr.HandsOff(c, 20, 10)
	if e == nil {
		t.Fatal("只有组长能强制踢出")
	}

	cms, err := svr.Online(c)
	if err != nil {
		t.Fatal(err)
	}
	if len(cms) != 2 {
		t.Fatal("在线人数错误")
	}
	for _, v := range cms {
		if (v.UID != 10 && v.UID != 20) || v.State != 1 {
			t.Fatal("在线人信息错误")
		}
	}
}

func BenchmarkMultiGetNextTask(b *testing.B) {

	var (
		mux    sync.RWMutex
		wg     = sync.WaitGroup{}
		ConMap = make(map[int64]struct{})
	)

	Audit := func(w *sync.WaitGroup, uid int64) {
		defer w.Done()
		for {
			tl, err := svr.Next(context.TODO(), uid)
			if err != nil {
				panic(err)
			}
			if tl == nil {
				fmt.Println("任务领取完成")
				return
			}

			mux.RLock()
			_, ok := ConMap[tl.ID]
			mux.RUnlock()
			if ok {
				panic(fmt.Sprintf("%d 重复下发:%d", uid, tl.ID))
			} else {
				mux.Lock()
				ConMap[tl.ID] = struct{}{}
				mux.Unlock()
			}

			fmt.Printf("uid=%d 领取任务:%d, weight=%d\n", uid, tl.ID, tl.Weight)

			time.Sleep(time.Millisecond * 50)
			if err = svr.Delay(context.TODO(), tl.ID, uid, "test-reson"); err != nil {
				panic(err)
			}
		}
	}

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go Audit(&wg, int64(i))
	}
	wg.Wait()
}

func TestFree(t *testing.T) {
	var c = context.TODO()
	rows := svr.Free(c, 481)
	if rows == 0 {
		t.Fail()
	}
	t.Fail()
}

func Test_setWeightConf(t *testing.T) {
	err := svr.setWeightConf(context.TODO(), "7", map[int64]*archive.WCItem{
		7: {
			Radio:  4,
			Weight: 18,
			Mtime:  utils.NewFormatTime(time.Now()),
			Desc:   "指派回流任务",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
