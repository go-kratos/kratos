package infoc

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"go-common/library/log"
	"go-common/library/net/ip"
)

type reporter struct {
	taskID string
	addr   string
	iip    string

	receiveCount int64
	sendCount    int64

	fails []string
}

func newReporter(taskID, addr string) (r *reporter) {
	r = &reporter{
		taskID: taskID,
		addr:   addr,
		iip:    ip.InternalIP(),
	}
	return
}

func (r *reporter) receiveIncr(delta int64) {
	atomic.AddInt64(&r.receiveCount, delta)
}

func (r *reporter) sendIncr(delta int64) {
	atomic.AddInt64(&r.sendCount, delta)
}

func (r *reporter) reportproc() {
	tick := time.NewTicker(1 * time.Minute)
	for {
		<-tick.C
		r.reporter()
	}
}

func (r *reporter) flush() {
	r.reporter()
}

func (r *reporter) reporter() {
	const _timeout = time.Second
	conn, err := net.DialTimeout("tcp", r.addr, _timeout)
	if err != nil {
		log.Error("infoc reporter flush dial error(%v)", err)
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(_timeout))
	var fails []string
	for _, fail := range r.fails {
		if _, err = conn.Write([]byte(fail)); err != nil {
			log.Error("infoc reporter write fail error(%v)", err)
			fails = append(fails, fail)
		}
	}
	for _, rc := range r.record(time.Now()) {
		if _, err = conn.Write([]byte(rc)); err != nil {
			log.Error("infoc reporter write error(%v)", err)
			fails = append(fails, rc)
		}
	}
	r.fails = fails
}

func (r *reporter) record(now time.Time) []string {
	rc := atomic.SwapInt64(&r.receiveCount, 0)
	sc := atomic.SwapInt64(&r.sendCount, 0)
	rcW := fmt.Sprintf("agent.receive.count\001%d\001%s\001%d\001%s\001\001", rc, r.iip, now.UnixNano()/int64(time.Millisecond), r.taskID)
	scW := fmt.Sprintf("agent.send.success.count\001%d\001%s\001%d\001%s\001\001", sc, r.iip, now.UnixNano()/int64(time.Millisecond), r.taskID)
	return []string{rcW, scW}
}
