package email

import (
	"fmt"
	"time"
)

//FastDetector detecte speed and unique analyze
type FastDetector struct {
	lastSec            int64
	sameSecCnt         int
	sameSecThreshold   int
	overspeedCnt       int
	overspeedThreshold int
	uniqueSpeed        map[int64]int
	fastUnique         int64
}

//NewFastDetector new
func NewFastDetector(speedThreshold, overspeedThreshold int) *FastDetector {
	return &FastDetector{
		lastSec:            time.Now().Unix(),
		sameSecCnt:         0,
		sameSecThreshold:   speedThreshold,
		overspeedCnt:       0,
		overspeedThreshold: overspeedThreshold,
		uniqueSpeed:        map[int64]int{},
	}
}

//String string info
func (fd *FastDetector) String() string {
	return fmt.Sprintf("same_sec_cnt=%d,overspeed_cnt=%d,fast_unique=%d,unique_speed=%v",
		fd.sameSecCnt, fd.overspeedCnt, fd.fastUnique, fd.uniqueSpeed)
}

//Detect 快慢探查, 超限名单只能被慢速/下一个超速名单/间隔5s后空名单代替，超限名单只有在(overspeedthreshold+1) * samesecthreshold时才确定，此时返回true
func (fd *FastDetector) Detect(unique int64) (fast bool) {
	now := time.Now().Unix()
	//连续n次超限
	if now == fd.lastSec {
		fd.sameSecCnt++
		if fd.sameSecCnt == fd.sameSecThreshold {
			fd.overspeedCnt++
		}
	} else {
		if fd.sameSecCnt < fd.sameSecThreshold || (now-fd.lastSec > 5) {
			fd.overspeedCnt = 0
			fd.fastUnique = 0
			fd.uniqueSpeed = map[int64]int{}
		}
		fd.sameSecCnt = 1
		fd.lastSec = now
	}

	//连续超限后，最先超限的unique指定为超限名单
	if fd.overspeedCnt == fd.overspeedThreshold && fd.sameSecCnt == fd.sameSecThreshold {
		fd.uniqueSpeed[unique] = 0
	}
	if (fd.overspeedCnt == fd.overspeedThreshold && fd.sameSecCnt != fd.sameSecThreshold) || (fd.overspeedCnt > fd.overspeedThreshold) {
		fd.uniqueSpeed[unique]++
		if fd.uniqueSpeed[unique] >= fd.sameSecThreshold {
			fast = true
			fd.fastUnique = unique //指定超限名单
			fd.uniqueSpeed = map[int64]int{}
			fd.overspeedCnt = 0
		}
	}

	return
}

//IsFastUnique 是否为超限名单
func (fd *FastDetector) IsFastUnique(unique int64) bool {
	return fd.fastUnique == unique
}
