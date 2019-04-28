package main

import (
	"context"
	"flag"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"go-common/app/service/main/riot-search/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"

	"github.com/ivpusic/grpool"
)

var (
	minID          uint64 = 1
	maxID          uint64 = 28731894
	keyword               = []string{"世界", "鬼畜", "自制", "搬运", "动漫", "崩坏", "搞笑", "德国", "弹幕", "乱入", "吸血鬼", "可怕", "骑士", "团长", "守护"}
	times          int
	count          int
	thread         int
	client         *bm.Client
	args           []*model.RiotSearchReq
	uri            string
	maxElapsedTime []int64
	avgElapsedTime []int64
)

//生成count个[start,end)结束的不重复的随机数
func generateRandomNumber(start uint64, end uint64, count int) []uint64 {
	if end < start || (end-start) < uint64(count) {
		return nil
	}
	nums := make([]uint64, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		num := uint64(r.Intn(int(end-start))) + start
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}

func benchmarkSearch(count int, times int, gid int) {
	var totalTime int64
	var maxTime int64
	for i := 0; i < times; i++ {
		// random chose params to use
		arg := args[rand.Intn(len(args))]
		params := url.Values{}
		var aids string
		for _, id := range arg.IDs {
			aids += strconv.FormatUint(id, 10)
		}
		text := arg.Keyword
		params.Set("aids", aids)
		params.Set("keyword", text)
		params.Set("pn", "1")
		params.Set("ps", "20")

		start := time.Now()
		err := client.Post(context.TODO(), uri, "", params, nil)
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(start)
		if int64(elapsed) > maxTime {
			maxTime = int64(elapsed)
		}
		totalTime += int64(elapsed)
	}
	avgElapsedTime[gid] = totalTime / (1000 * 1000 * int64(times))
	maxElapsedTime[gid] = maxTime / (1000 * 1000)
}

func init() {
	flag.IntVar(&times, "times", 100, "单个线程测试次数")
	flag.IntVar(&count, "count", 100000, "每次测试aid个数")
	flag.IntVar(&thread, "thread", 10, "线程数")
	flag.StringVar(&uri, "uri", "http://127.0.0.1:7871/x/internal/riot-search/arc/ids", "请求url")
	flag.Parse()
	log.Info("times: %d, count:%d, thread:%d, uri:%s", times, count, thread, uri)
	log.Info("init http client")
	app := &bm.App{
		Key:    "test",
		Secret: "test",
	}
	clientConf := &bm.ClientConfig{
		App:       app,
		Timeout:   xtime.Duration(time.Second * 1),
		Dial:      xtime.Duration(time.Second),
		KeepAlive: xtime.Duration(time.Second * 60),
	}
	client = bm.NewClient(clientConf)
	log.Info("init 10 http request params, random chose one to test")
	args = make([]*model.RiotSearchReq, 10)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		arg := &model.RiotSearchReq{
			IDs:     generateRandomNumber(minID, maxID, count),
			Keyword: keyword[rand.Intn(len(keyword))],
		}
		args[i] = arg
	}
	maxElapsedTime = make([]int64, thread)
	avgElapsedTime = make([]int64, thread)
	log.Info("init params finished")
}

func main() {
	log.Info("start test")
	if thread >= 1000 {
		panic("thread large than 1000 is not allowed")
	}
	pool := grpool.NewPool(thread, 10240)
	defer pool.Release()
	pool.WaitCount(thread)
	for i := 0; i < thread; i++ {
		threadNum := i
		pool.JobQueue <- func() {
			benchmarkSearch(count, times, threadNum)
			pool.JobDone()
		}
	}
	pool.WaitAll()
	log.Info("avg elapsed times list: %v", avgElapsedTime)
	log.Info("max elapsed times list: %v", maxElapsedTime)
	log.Info("test finished")
}
