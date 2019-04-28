package lancermonitor

import (
	"time"
	"net"
	"sync"
	"fmt"
	"strconv"
	"strings"
	"errors"
	"go-common/library/log"
)

const (
	_separator = "####"
)

var (
	lm      *LancerMonitor
	started bool
)

type LancerMonitor struct {
	c                *Config
	logRevStatusLock sync.Mutex
	logRevStatus     map[string]int64
	ipAddr           string
}

func InitLancerMonitor(config *Config) (l *LancerMonitor, err error) {
	if started {
		return nil, errors.New("lancer Monitor can only be init Once")
	}
	if err = config.ConfigValidate(); err != nil {
		return nil, err
	}

	l = new(LancerMonitor)
	l.c = config
	l.logRevStatus = make(map[string]int64)
	l.ipAddr = InternalIP()

	go l.reportStatus()
	started = true
	lm = l
	return l, nil
}

func (l *LancerMonitor) reportStatus() {
	reportStatusTk := time.Tick(time.Duration(60 * time.Second))
	for {
		select {
		case <-reportStatusTk:
			logCount := l.getLogCount()
			conn, error := net.DialTimeout("tcp", l.c.Addr, time.Second*5)
			if error != nil {
				log.Error("failed to connect to lancer when report status")
			} else {
				for k, v := range logCount {
					fields := strings.Split(k, _separator)
					if len(fields) == 2 {
						fmt.Fprintf(conn, fields[0]+"\u0001"+strconv.FormatInt(v, 10)+"\u0001"+l.ipAddr+"\u0001"+strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)+"\u0001"+fields[1]+"\u0001\u0001")
					}
				}
				log.Info("report status to lancer")
				conn.Close()
			}
		}
	}
}

// InternalIP get internal ip.
func InternalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

//get log count of each logid since last call
func (l *LancerMonitor) getLogCount() map[string]int64 {
	l.logRevStatusLock.Lock()
	defer l.logRevStatusLock.Unlock()
	logRevSendStatus := make(map[string]int64)
	for k, v := range l.logRevStatus {
		logRevSendStatus[k] = v
	}
	for k := range l.logRevStatus {
		delete(l.logRevStatus, k)
	}
	return logRevSendStatus
}

func IncreaseLogCount(name string, logId string) {
	if lm == nil || !started {
		return
	}
	if name == "" || logId == "" {
		return
	}
	key := name + _separator + logId
	lm.logRevStatusLock.Lock()
	defer lm.logRevStatusLock.Unlock()
	lm.logRevStatus[key] += 1
}
