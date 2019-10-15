package cpu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	pscpu "github.com/shirou/gopsutil/cpu"
)

type cgroupCPU struct {
	frequency uint64
	quota     float64
	cores     uint64

	preSystem uint64
	preTotal  uint64
	usage     uint64
}

func newCgroupCPU() (cpu *cgroupCPU, err error) {
	var cores int
	cores, err = pscpu.Counts(true)
	if err != nil || cores == 0 {
		var cpus []uint64
		cpus, err = perCPUUsage()
		if err != nil {
			err = errors.Errorf("perCPUUsage() failed!err:=%v", err)
			return
		}
		cores = len(cpus)
	}

	sets, err := cpuSets()
	if err != nil {
		err = errors.Errorf("cpuSets() failed!err:=%v", err)
		return
	}
	quota := float64(len(sets))
	cq, err := cpuQuota()
	if err == nil && cq != -1 {
		var period uint64
		if period, err = cpuPeriod(); err != nil {
			err = errors.Errorf("cpuPeriod() failed!err:=%v", err)
			return
		}
		limit := float64(cq) / float64(period)
		if limit < quota {
			quota = limit
		}
	}
	maxFreq := cpuMaxFreq()

	preSystem, err := systemCPUUsage()
	if err != nil {
		err = errors.Errorf("systemCPUUsage() failed!err:=%v", err)
		return
	}
	preTotal, err := totalCPUUsage()
	if err != nil {
		err = errors.Errorf("totalCPUUsage() failed!err:=%v", err)
		return
	}
	cpu = &cgroupCPU{
		frequency: maxFreq,
		quota:     quota,
		cores:     uint64(cores),
		preSystem: preSystem,
		preTotal:  preTotal,
	}
	return
}

func (cpu *cgroupCPU) Usage() (u uint64, err error) {
	var (
		total  uint64
		system uint64
	)
	total, err = totalCPUUsage()
	if err != nil {
		return
	}
	system, err = systemCPUUsage()
	if err != nil {
		return
	}
	if system != cpu.preSystem {
		u = uint64(float64((total-cpu.preTotal)*cpu.cores*1e3) / (float64(system-cpu.preSystem) * cpu.quota))
	}
	cpu.preSystem = system
	cpu.preTotal = total
	return
}

func (cpu *cgroupCPU) Info() Info {
	return Info{
		Frequency: cpu.frequency,
		Quota:     cpu.quota,
	}
}

const nanoSecondsPerSecond = 1e9

// ErrNoCFSLimit is no quota limit
var ErrNoCFSLimit = errors.Errorf("no quota limit")

var clockTicksPerSecond = uint64(getClockTicks())

// systemCPUUsage returns the host system's cpu usage in
// nanoseconds. An error is returned if the format of the underlying
// file does not match.
//
// Uses /proc/stat defined by POSIX. Looks for the cpu
// statistics line and then sums up the first seven fields
// provided. See man 5 proc for details on specific field
// information.
func systemCPUUsage() (usage uint64, err error) {
	var (
		line string
		f    *os.File
	)
	if f, err = os.Open("/proc/stat"); err != nil {
		return
	}
	bufReader := bufio.NewReaderSize(nil, 128)
	defer func() {
		bufReader.Reset(nil)
		f.Close()
	}()
	bufReader.Reset(f)
	for err == nil {
		if line, err = bufReader.ReadString('\n'); err != nil {
			err = errors.WithStack(err)
			return
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "cpu":
			if len(parts) < 8 {
				err = errors.WithStack(fmt.Errorf("bad format of cpu stats"))
				return
			}
			var totalClockTicks uint64
			for _, i := range parts[1:8] {
				var v uint64
				if v, err = strconv.ParseUint(i, 10, 64); err != nil {
					err = errors.WithStack(fmt.Errorf("error parsing cpu stats"))
					return
				}
				totalClockTicks += v
			}
			usage = (totalClockTicks * nanoSecondsPerSecond) / clockTicksPerSecond
			return
		}
	}
	err = errors.Errorf("bad stats format")
	return
}

func totalCPUUsage() (usage uint64, err error) {
	var cg *cgroup
	if cg, err = currentcGroup(); err != nil {
		return
	}
	return cg.CPUAcctUsage()
}

func perCPUUsage() (usage []uint64, err error) {
	var cg *cgroup
	if cg, err = currentcGroup(); err != nil {
		return
	}
	return cg.CPUAcctUsagePerCPU()
}

func cpuSets() (sets []uint64, err error) {
	var cg *cgroup
	if cg, err = currentcGroup(); err != nil {
		return
	}
	return cg.CPUSetCPUs()
}

func cpuQuota() (quota int64, err error) {
	var cg *cgroup
	if cg, err = currentcGroup(); err != nil {
		return
	}
	return cg.CPUCFSQuotaUs()
}

func cpuPeriod() (peroid uint64, err error) {
	var cg *cgroup
	if cg, err = currentcGroup(); err != nil {
		return
	}
	return cg.CPUCFSPeriodUs()
}

func cpuFreq() uint64 {
	lines, err := readLines("/proc/cpuinfo")
	if err != nil {
		return 0
	}
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		if key == "cpu MHz" || key == "clock" {
			// treat this as the fallback value, thus we ignore error
			if t, err := strconv.ParseFloat(strings.Replace(value, "MHz", "", 1), 64); err == nil {
				return uint64(t * 1000.0 * 1000.0)
			}
		}
	}
	return 0
}

func cpuMaxFreq() uint64 {
	feq := cpuFreq()
	data, err := readFile("/sys/devices/system/cpu/cpu0/cpufreq/cpuinfo_max_freq")
	if err != nil {
		return feq
	}
	// override the max freq from /proc/cpuinfo
	cfeq, err := parseUint(data)
	if err == nil {
		feq = cfeq
	}
	return feq
}

//GetClockTicks get the OS's ticks per second
func getClockTicks() int {
	// TODO figure out a better alternative for platforms where we're missing cgo
	//
	// TODO Windows. This could be implemented using Win32 QueryPerformanceFrequency().
	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms644905(v=vs.85).aspx
	//
	// An example of its usage can be found here.
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn553408(v=vs.85).aspx

	return 100
}
