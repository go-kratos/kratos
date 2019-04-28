// +build darwin

package cpu

var su uint64 = 10
var tu uint64 = 10

func systemCPUUsage() (usage uint64, err error) {
	su += 1000
	return su, nil
}
func totalCPUUsage() (usage uint64, err error) {
	tu += 500
	return tu, nil
}
func perCPUUsage() (usage []uint64, err error) { return []uint64{10, 10, 10, 10}, nil }
func cpuSets() (sets []uint64, err error)      { return []uint64{0, 1, 2, 3}, nil }
func cpuQuota() (quota int64, err error)       { return 100, nil }
func cpuPeriod() (peroid uint64, err error)    { return 10, nil }
func cpuMaxFreq() (feq uint64)                 { return 10 }
