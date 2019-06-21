package cpu

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type psutilCPU struct {
	interval time.Duration
}

func newPsutilCPU(interval time.Duration) (cpu *psutilCPU, err error) {
	cpu = &psutilCPU{interval: interval}
	_, err = cpu.Usage()
	if err != nil {
		return
	}
	return
}

func (ps *psutilCPU) Usage() (u uint64, err error) {
	var percents []float64
	percents, err = cpu.Percent(ps.interval, false)
	if err == nil {
		u = uint64(percents[0] * 10)
	}
	return
}

func (ps *psutilCPU) Info() (info Info) {
	stats, err := cpu.Info()
	if err != nil {
		return
	}
	cores, err := cpu.Counts(true)
	if err != nil {
		return
	}
	info = Info{
		Frequency: uint64(stats[0].Mhz),
		Quota:     float64(cores),
	}
	return
}
