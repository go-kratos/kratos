package cpu

import (
	"fmt"
	"go-common/library/log"
	"sync/atomic"
	"time"
)

var (
	cores   uint64
	maxFreq uint64
	quota   float64
	usage   uint64

	preSystem uint64
	preTotal  uint64
)

func init() {
	cpus, err := perCPUUsage()
	if err != nil {
		panic(fmt.Errorf("stat/sys/cpu: perCPUUsage() failed!err:=%v", err))
	}
	cores = uint64(len(cpus))

	sets, err := cpuSets()
	if err != nil {
		panic(fmt.Errorf("stat/sys/cpu: cpuSets() failed!err:=%v", err))
	}
	quota = float64(len(sets))
	cq, err := cpuQuota()
	if err == nil {
		if cq != -1 {
			var period uint64
			if period, err = cpuPeriod(); err != nil {
				panic(fmt.Errorf("stat/sys/cpu: cpuPeriod() failed!err:=%v", err))
			}
			limit := float64(cq) / float64(period)
			if limit < quota {
				quota = limit
			}
		}
	}
	maxFreq = cpuMaxFreq()

	preSystem, err = systemCPUUsage()
	if err != nil {
		panic(fmt.Errorf("sys/cpu: systemCPUUsage() failed!err:=%v", err))
	}
	preTotal, err = totalCPUUsage()
	if err != nil {
		panic(fmt.Errorf("sys/cpu: totalCPUUsage() failed!err:=%v", err))
	}

	go func() {
		ticker := time.NewTicker(time.Millisecond * 250)
		defer ticker.Stop()
		for {
			<-ticker.C
			cpu := refreshCPU()
			if cpu != 0 {
				atomic.StoreUint64(&usage, cpu)
			}
		}
	}()
}

func refreshCPU() (u uint64) {
	total, err := totalCPUUsage()
	if err != nil {
		log.Warn("os/stat: get totalCPUUsage failed,error(%v)", err)
		return
	}
	system, err := systemCPUUsage()
	if err != nil {
		log.Warn("os/stat: get systemCPUUsage failed,error(%v)", err)
		return
	}
	if system != preSystem {
		u = uint64(float64((total-preTotal)*cores*1e3) / (float64(system-preSystem) * quota))
	}
	preSystem = system
	preTotal = total
	return u
}

// Stat cpu stat.
type Stat struct {
	Usage uint64 // cpu use ratio.
}

// Info cpu info.
type Info struct {
	Frequency uint64
	Quota     float64
}

// ReadStat read cpu stat.
func ReadStat(stat *Stat) {
	stat.Usage = atomic.LoadUint64(&usage)
}

// GetInfo get cpu info.
func GetInfo() Info {
	return Info{
		Frequency: maxFreq,
		Quota:     quota,
	}
}
