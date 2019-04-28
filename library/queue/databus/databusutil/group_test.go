package databusutil

import (
	"context"
	"encoding/json"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
)

type testMsg struct {
	Seq int64 `json:"seq"`
	Mid int64 `json:"mid"`
	Now int64 `json:"now"`
}

var (
	_sendSeqsList = make([][]int64, _groupNum)
	_recvSeqsList = make([][]int64, _groupNum)

	_sMus = make([]sync.Mutex, _groupNum)
	_rMus = make([]sync.Mutex, _groupNum)

	_groupNum = 8

	_tc = 20
	_ts = time.Now().Unix()
	_st = _ts - _ts%10 + 1000
	_ed = _bSt + int64(_groupNum*_tc) - 1

	_dsPubConf = &databus.Config{
		Key:          "0PvKGhAqDvsK7zitmS8t",
		Secret:       "0PvKGhAqDvsK7zitmS8u",
		Group:        "databus_test_group",
		Topic:        "databus_test_topic",
		Action:       "pub",
		Name:         "databus",
		Proto:        "tcp",
		Addr:         "172.16.33.158:6205",
		Active:       1,
		Idle:         1,
		DialTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		IdleTimeout:  xtime.Duration(time.Minute),
	}

	_dsSubConf = &databus.Config{
		Key:          "0PvKGhAqDvsK7zitmS8t",
		Secret:       "0PvKGhAqDvsK7zitmS8u",
		Group:        "databus_test_group",
		Topic:        "databus_test_topic",
		Action:       "sub",
		Name:         "databus",
		Proto:        "tcp",
		Addr:         "172.16.33.158:6205",
		Active:       1,
		Idle:         1,
		DialTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second * 35),
		IdleTimeout:  xtime.Duration(time.Minute),
	}
)

func TestGroup(t *testing.T) {
	for i := 0; i < _groupNum; i++ {
		_sendSeqsList[i] = make([]int64, 0)
		_recvSeqsList[i] = make([]int64, 0)
	}
	taskCounts := taskCount(_groupNum, _st, _ed)

	runtime.GOMAXPROCS(32)
	log.Init(&log.Config{
		Dir: "/data/log/queue",
	})
	c := &Config{
		Size:   200,
		Ticker: xtime.Duration(time.Second),
		Num:    _groupNum,
		Chan:   1024,
	}
	dsSub := databus.New(_dsSubConf)
	defer dsSub.Close()
	group := NewGroup(
		c,
		dsSub.Messages(),
	)
	group.New = newTestMsg
	group.Split = split
	group.Do = do
	eg, _ := errgroup.WithContext(context.Background())
	// go produce test messages
	eg.Go(func() error {
		send(_st, _ed)
		return nil
	})
	// go consume test messages
	eg.Go(func() error {
		group.Start()
		defer group.Close()
		m := make(map[int]struct{})
		for len(m) < _groupNum {
			for i := 0; i < _groupNum; i++ {
				_, ok := m[i]
				if ok {
					continue
				}
				_rMus[i].Lock()
				if len(_recvSeqsList[i]) == taskCounts[i] {
					m[i] = struct{}{}
				}
				_rMus[i].Unlock()
				log.Info("_recvSeqsList[%d] length: %d, expect: %d", i, len(_recvSeqsList[i]), taskCounts[i])
			}
			log.Info("m length: %d", len(m))
			time.Sleep(time.Millisecond * 500)
		}
		// check seqs list, sendSeqsList and recvSeqsList will not change since now, so no need to lock
		for num := 0; num < _groupNum; num++ {
			sendSeqs := _sendSeqsList[num]
			recvSeqs := _recvSeqsList[num]
			if len(sendSeqs) != taskCounts[num] {
				t.Errorf("sendSeqs length of proc %d is incorrect, expcted %d but got %d", num, taskCounts[num], len(sendSeqs))
				t.FailNow()
			}
			if len(recvSeqs) != taskCounts[num] {
				t.Errorf("recvSeqs length of proc %d is incorrect, expcted %d but got %d", num, taskCounts[num], len(recvSeqs))
				t.FailNow()
			}
			for i := range recvSeqs {
				if recvSeqs[i] != sendSeqs[i] {
					t.Errorf("res is incorrect for proc %d, expcted recvSeqs[%d] equal to sendSeqs[%d] but not, recvSeqs[%d]: %d, sendSeqs[%d]: %d", num, i, i, i, recvSeqs[i], i, sendSeqs[i])
					t.FailNow()
				}
			}
			t.Logf("proc %d processed %d messages, expected %d messages, check ok", num, taskCounts[num], len(recvSeqs))
		}
		return nil
	})
	eg.Wait()
}

func do(msgs []interface{}) {
	for _, m := range msgs {
		if msg, ok := m.(*testMsg); ok {
			shard := int(msg.Mid) % _groupNum
			if msg.Seq < _st {
				log.Info("proc %d processed old seq: %d, mid: %d", shard, msg.Seq, msg.Mid)
				continue
			}
			_rMus[shard].Lock()
			_recvSeqsList[shard] = append(_recvSeqsList[shard], msg.Seq)
			_rMus[shard].Unlock()
			log.Info("proc %d processed seq: %d, mid: %d", shard, msg.Seq, msg.Mid)
		}
	}
}

func send(st, ed int64) error {
	dsPub := databus.New(_dsPubConf)
	defer dsPub.Close()
	ts := time.Now().Unix()
	for i := st; i <= ed; i++ {
		mid := int64(i)
		seq := i
		k := _dsPubConf.Topic + strconv.FormatInt(mid, 10)
		n := &testMsg{
			Seq: seq,
			Mid: mid,
			Now: ts,
		}
		dsPub.Send(context.TODO(), k, n)
		// NOTE: sleep here to avoid network latency caused message out of sequence
		time.Sleep(time.Millisecond * 500)
		shard := int(mid) % _groupNum
		_sMus[shard].Lock()
		_sendSeqsList[shard] = append(_sendSeqsList[shard], seq)
		_sMus[shard].Unlock()
	}
	return nil
}

func newTestMsg(msg *databus.Message) (res interface{}, err error) {
	res = new(testMsg)
	if err = json.Unmarshal(msg.Value, &res); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
	}
	return
}

func split(msg *databus.Message, data interface{}) int {
	t, ok := data.(*testMsg)
	if !ok {
		return 0
	}
	return int(t.Mid)
}

func taskCount(num int, st, ed int64) []int {
	res := make([]int, num)
	for i := st; i <= ed; i++ {
		res[int(i)%num]++
	}
	return res
}

func TestTaskCount(t *testing.T) {
	groupNum := 10
	c := 100
	ts := time.Now().Unix()
	st := ts - ts%10 + 1000
	ed := st + int64(groupNum*c) - 1
	res := taskCount(groupNum, st, ed)
	for i, v := range res {
		if v != c {
			t.Errorf("res is incorrect, expected task count 10 for proc %d but got %d", i, v)
			t.FailNow()
		}
		t.Logf("i: %d, v: %d", i, v)
	}
}

var (
	_bGroupNum = 3

	_bSendSeqsList = make([][]int64, _bGroupNum)
	_bRecvSeqsList = make([][]int64, _bGroupNum)

	_bSMus = make([]sync.Mutex, _bGroupNum)
	_bRMus = make([]sync.Mutex, _bGroupNum)

	_bTc = 20
	_bTs = time.Now().Unix()
	_bSt = _bTs - _bTs%10 + 1000
	_bEd = _bSt + int64(_bGroupNum*_bTc) - 1

	_bTaskCounts = taskCount(_bGroupNum, _bSt, _bEd)

	_blockDo   = true
	_blockDoMu sync.Mutex

	_blocked = false
)

func TestGroup_Blocking(t *testing.T) {
	for i := 0; i < _bGroupNum; i++ {
		_bSendSeqsList[i] = make([]int64, 0)
		_bRecvSeqsList[i] = make([]int64, 0)
	}

	runtime.GOMAXPROCS(32)
	log.Init(&log.Config{
		Dir: "/data/log/queue",
	})
	c := &Config{
		Size:   20,
		Ticker: xtime.Duration(time.Second),
		Num:    _bGroupNum,
		Chan:   5,
	}

	dsSub := databus.New(_dsSubConf)
	defer dsSub.Close()
	g := NewGroup(
		c,
		dsSub.Messages(),
	)
	g.New = newTestMsg
	g.Split = split
	g.Do = func(msgs []interface{}) {
		blockingDo(t, g, msgs)
	}
	eg, _ := errgroup.WithContext(context.Background())
	// go produce test messages
	eg.Go(func() error {
		dsPub := databus.New(_dsPubConf)
		defer dsPub.Close()
		ts := time.Now().Unix()
		for i := _bSt; i <= _bEd; i++ {
			mid := int64(i)
			seq := i
			k := _dsPubConf.Topic + strconv.FormatInt(mid, 10)
			n := &testMsg{
				Seq: seq,
				Mid: mid,
				Now: ts,
			}
			dsPub.Send(context.TODO(), k, n)
			// NOTE: sleep here to avoid network latency caused message out of sequence
			time.Sleep(time.Millisecond * 500)
			shard := int(mid) % _bGroupNum
			_bSMus[shard].Lock()
			_bSendSeqsList[shard] = append(_bSendSeqsList[shard], seq)
			_bSMus[shard].Unlock()
		}
		return nil
	})
	// go consume test messages
	eg.Go(func() error {
		g.Start()
		defer g.Close()
		m := make(map[int]struct{})
		// wait until all proc process theirs messages done
		for len(m) < _bGroupNum {
			for i := 0; i < _bGroupNum; i++ {
				_, ok := m[i]
				if ok {
					continue
				}
				_bRMus[i].Lock()
				if len(_bRecvSeqsList[i]) == _bTaskCounts[i] {
					m[i] = struct{}{}
				}
				_bRMus[i].Unlock()
				log.Info("_bRecvSeqsList[%d] length: %d, expect: %d, blockDo: %t", i, len(_bRecvSeqsList[i]), _bTaskCounts[i], _blockDo)
			}
			log.Info("m length: %d", len(m))
			time.Sleep(time.Millisecond * 500)
		}
		return nil
	})
	eg.Wait()
}

func blockingDo(t *testing.T, g *Group, msgs []interface{}) {
	_blockDoMu.Lock()
	if !_blockDo {
		_blockDoMu.Unlock()
		processMsg(msgs)
		return
	}
	// blocking to see if consume proc blocks finally
	lastGLen := 0
	cnt := 0
	for i := 0; i < 60; i++ {
		// print seqs status, not lock because final stable
		for i, v := range _bRecvSeqsList {
			log.Info("_bRecvSeqsList[%d] length: %d, expect: %d", i, len(v), _bTaskCounts[i])
		}
		gLen := 0
		for h := g.head; h != nil; h = h.next {
			gLen++
		}
		if gLen == lastGLen {
			cnt++
		} else {
			cnt = 0
		}
		lastGLen = gLen
		log.Info("blocking test: gLen: %d, cnt: %d, _bSt: %d, _bEd: %d", gLen, cnt, _bSt, _bEd)
		if cnt == 5 {
			_blocked = true
			log.Info("blocking test: consumeproc now is blocked, now trying to unblocking do callback")
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	// assert blocked
	if !_blocked {
		t.Errorf("res is incorrect, _blocked should be true but got false")
		t.FailNow()
	}
	// unblocking and check if consume proc unblocking too
	_blockDo = false
	_blockDoMu.Unlock()
	processMsg(msgs)
}

func processMsg(msgs []interface{}) {
	for _, m := range msgs {
		if msg, ok := m.(*testMsg); ok {
			shard := int(msg.Mid) % _bGroupNum
			if msg.Seq < _bSt {
				log.Info("proc %d processed old seq: %d, mid: %d", shard, msg.Seq, msg.Mid)
				continue
			}
			_bRMus[shard].Lock()
			_bRecvSeqsList[shard] = append(_bRecvSeqsList[shard], msg.Seq)
			log.Info("appended: %d", msg.Seq)
			_bRMus[shard].Unlock()
			log.Info("proc %d processed seq: %d, mid: %d", shard, msg.Seq, msg.Mid)
		}
	}
}
