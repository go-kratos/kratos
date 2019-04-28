// Copyright 2017 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/gt/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package info

import (
	"errors"
	"fmt"
	"sync"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

var (
	lck sync.RWMutex

	// InitMemUsed init mem used
	InitMemUsed uint64
	// InitDiskUsed init disk used
	InitDiskUsed uint64
)

func init() {
	lck.Lock()
	InitMemUsed, _ = MemUsed()
	InitDiskUsed, _ = DiskUsed()
	lck.Unlock()
}

// MemPercent returns the amount of use memory in percent.
func MemPercent() (string, error) {
	memInfo, err := mem.VirtualMemory()
	useMem := fmt.Sprintf("%.2f", memInfo.UsedPercent)

	return useMem, err
}

// MemUsed returns the amount of used memory in bytes.
func MemUsed() (uint64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return memInfo.Used, err
}

// UsedMem returns the amount of riot used memory in bytes
// after init() func.
func UsedMem() (uint64, error) {
	memInfo, err := mem.VirtualMemory()
	// memInfo, err := MemUsed()
	if err != nil {
		return 0, err
	}

	return memInfo.Used - InitMemUsed, err
}

// MemTotal returns the amount of total memory in bytes.
func MemTotal() (uint64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return memInfo.Total, err
}

// MemFree returns the amount of free memory in bytes.
func MemFree() (uint64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return memInfo.Free, err
}

// ToKB bytes to kb
func ToKB(data uint64) uint64 {
	return data / 1024
}

// ToMB bytes to mb
func ToMB(data uint64) uint64 {
	return data / 1024 / 1024
}

// ToGB bytes to gb
func ToGB(data uint64) uint64 {
	return data / 1024 / 1024 / 1024
}

// Disk init the disk
func Disk(pt ...bool) ([]*disk.UsageStat, error) {
	var ptBool bool
	if len(pt) > 0 {
		ptBool = pt[0]
	}

	var usage []*disk.UsageStat
	parts, err := disk.Partitions(ptBool)

	for _, part := range parts {
		use, err := disk.Usage(part.Mountpoint)
		if err != nil {
			return usage, err
		}
		usage = append(usage, use)
		// printUsage(use)
	}

	return usage, err
}

// DiskPercent returns the amount of use disk in percent.
func DiskPercent() (string, error) {
	usage, err := Disk()
	if len(usage) > 0 {
		useDisk := fmt.Sprintf("%.2f", usage[0].UsedPercent)
		return useDisk, err
	}

	return "0.00", err
}

// DiskUsed returns the amount of use disk in bytes.
func DiskUsed() (uint64, error) {
	usage, err := Disk()
	// for i := 0; i < len(usage); i++ {
	if len(usage) > 0 {
		useDisk := usage[0].Used
		return useDisk, err
	}

	return 0, err
}

// UsedDisk returns the amount of use disk in bytes
// after init() func.
func UsedDisk() (uint64, error) {
	diskUsed, err := DiskUsed()
	if err != nil {
		return 0, err
	}

	return diskUsed - InitDiskUsed, err
}

// DiskTotal returns the amount of total disk in bytes.
func DiskTotal() (uint64, error) {
	usage, err := Disk()
	// for i := 0; i < len(usage); i++ {
	if len(usage) > 0 {
		totalDisk := usage[0].Total
		return totalDisk, err
	}

	return 0, err
}

// DiskFree returns the amount of free disk in bytes.
func DiskFree() (uint64, error) {
	usage, err := Disk()
	// for i := 0; i < len(usage); i++ {
	if len(usage) > 0 {
		freeDisk := usage[0].Free
		return freeDisk, err
	}

	return 0, err
}

// CPUInfo returns the cpu info
func CPUInfo(args ...int) (string, error) {
	info, err := cpu.Info()
	if err != nil {
		return "", err
	}

	if len(info) == 0 {
		return "", errors.New("no CPU detected")
	}

	if len(args) > 0 {
		return info[args[0]].ModelName, nil
	}

	return info[0].ModelName, nil
}

// CPUPercent returns the amount of use cpu in percent.
func CPUPercent() ([]float64, error) {
	used, err := cpu.Percent(0, true)
	return used, err
}

// Uptime returns the system uptime in seconds.
func Uptime() (uptime uint64, err error) {
	hostInfo, err := host.Info()
	if err != nil {
		return 0, err
	}
	return hostInfo.Uptime, nil
}

// PlatformInfo fetches system platform information.
func PlatformInfo() (platform, family, osVersion string, err error) {
	platform, family, osVersion, err = host.PlatformInformation()

	return
}

// Platform returns the platform name and OS Version.
func Platform() (string, error) {
	platform, _, osVersion, err := host.PlatformInformation()

	return platform + " " + osVersion, err
}

// KernelVer returns the kernel version as a string.
func KernelVer() (string, error) {
	return host.KernelVersion()
}
