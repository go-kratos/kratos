// Copyright 2016 ego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package riot

import (
	"sync"

	"github.com/go-vgo/gt/info"
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
	return info.MemPercent()
}

// MemUsed returns the amount of used memory in bytes.
func MemUsed() (uint64, error) {
	return info.MemUsed()
}

// UsedMem returns the amount of riot used memory in bytes
// after init() func.
func (engine *Engine) UsedMem() (uint64, error) {
	memUsed, err := MemUsed()
	if err != nil {
		return 0, err
	}

	return memUsed - InitMemUsed, err
}

// MemTotal returns the amount of total memory in bytes.
func MemTotal() (uint64, error) {
	return info.MemTotal()
}

// MemFree returns the amount of free memory in bytes.
func MemFree() (uint64, error) {
	return info.MemFree()
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
// func Disk(pt ...bool) ([]*disk.UsageStat, error) {
// 	return info.Disk(pt...)
// }

// DiskPercent returns the amount of use disk in percent.
func DiskPercent() (string, error) {
	return info.DiskPercent()
}

// DiskUsed returns the amount of use disk in bytes.
func DiskUsed() (uint64, error) {
	return info.DiskUsed()
}

// UsedDisk returns the amount of use disk in bytes
// after init() func.
func (engine *Engine) UsedDisk() (uint64, error) {
	diskUsed, err := DiskUsed()
	if err != nil {
		return 0, err
	}

	return diskUsed - InitDiskUsed, err
}

// DiskTotal returns the amount of total disk in bytes.
func DiskTotal() (uint64, error) {
	return info.DiskTotal()
}

// DiskFree returns the amount of free disk in bytes.
func DiskFree() (uint64, error) {
	return info.DiskFree()
}

// CPUInfo returns the cpu info
func CPUInfo(args ...int) (string, error) {
	return info.CPUInfo(args...)
}

// CPUPercent returns the amount of use cpu in percent.
func CPUPercent() ([]float64, error) {
	return info.CPUPercent()
}

// Uptime returns the system uptime in seconds.
func Uptime() (uptime uint64, err error) {
	return info.Uptime()
}

// PlatformInfo fetches system platform information.
func PlatformInfo() (platform, family, osVersion string, err error) {
	return info.PlatformInfo()
}

// Platform returns the platform name and OS Version.
func Platform() (string, error) {
	return info.Platform()
}

// KernelVer returns the kernel version as a string.
func KernelVer() (string, error) {
	return info.KernelVer()
}
