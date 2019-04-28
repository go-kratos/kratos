package limit

import (
	"errors"
	"os"
	"fmt"
	"io/ioutil"
	"strconv"
)

var (
	LimitConfNil          = errors.New("Config of resource limit is nil")
	LimitConfError        = errors.New("LimitConfError")
	CgroupPathNotExist    = errors.New("Cgroup Path Not Exist")
	MemCgroupPathNotExist = errors.New("Mem subgroup Not Exist")
	CpuCgroupPathNotExist = errors.New("Cpu subgroup Not Exist")
)

type LimitConf struct {
	AppName         string
	CgroupPath      string
	LimitMemMB      int
	LimitMemEnabled bool
	LimitCpuCore    int
	LimitCpuEnabled bool
}

type Limit struct {
	c *LimitConf
}

// LimitRes init Limit
func LimitRes(c *LimitConf) (l *Limit, err error) {
	if c == nil {
		return nil, LimitConfError
	}
	l = new(Limit)
	l.c = c

	if c.AppName == "" {
		return nil, fmt.Errorf("AppName can't be nil")
	}

	if err = pathExists(l.c.CgroupPath, false); err != nil {
		return nil, CgroupPathNotExist
	}
	if l.c.LimitMemEnabled {
		if l.c.LimitCpuCore <= 0 {
			return nil, fmt.Errorf("LimitCpuCore must be greater than 0")
		}
		if err = l.limitMem(); err != nil {
			return nil, err
		}
	}
	if l.c.LimitCpuEnabled {
		if l.c.LimitMemMB <= 0 {
			return nil, fmt.Errorf("LimitMemMB must be greater than 0")
		}
		if err = l.limitCpu(); err != nil {
			return nil, err
		}
	}
	return
}

// limitMem limit memory by memory.limit_in_bytes
func (l *Limit) limitMem() (err error) {
	if err = pathExists(fmt.Sprintf("%s/memory", l.c.CgroupPath), false); err != nil {
		return MemCgroupPathNotExist
	}

	memPath := fmt.Sprintf("%s/memory/%s/", l.c.CgroupPath, l.c.AppName)
	if err = pathExists(memPath, true); err != nil {
		return err
	}
	// cgroup.procs
	pidPath := memPath + "cgroup.procs"
	if err = pathExists(pidPath, false); err != nil {
		return err
	}
	if err = ioutil.WriteFile(pidPath, []byte(strconv.Itoa(os.Getegid())), 0644); err != nil {
		return err
	}
	// memory.limit_in_bytes
	limitPath := memPath + "memory.limit_in_bytes"
	if err = pathExists(limitPath, false); err != nil {
		return err
	}
	if err = ioutil.WriteFile(limitPath, []byte(fmt.Sprintf("%sM", strconv.Itoa(l.c.LimitMemMB))), 0644); err != nil {
		return err
	}
	return
}

// limitCpu limit cpu by cpu.cfs_quota_us and cpu.cfs_period_us
func (l *Limit) limitCpu() (err error) {
	if err = pathExists(fmt.Sprintf("%s/cpu,cpuacct", l.c.CgroupPath), false); err != nil {
		return CpuCgroupPathNotExist
	}

	cpuPath := fmt.Sprintf("%s/cpu,cpuacct/%s/", l.c.CgroupPath, l.c.AppName)
	if err = pathExists(cpuPath, true); err != nil {
		return err
	}
	// cpu.cfs_quota_us
	quotaPath := cpuPath + "cpu.cfs_quota_us"
	if err = pathExists(quotaPath, false); err != nil {
		return err
	}
	if err = ioutil.WriteFile(quotaPath, []byte(strconv.Itoa(10000)), 0644); err != nil {
		return err
	}
	// cpu.cfs_period_us
	periodPath := cpuPath + "cpu.cfs_period_us"
	if err = pathExists(periodPath, false); err != nil {
		return err
	}
	if err = ioutil.WriteFile(periodPath, []byte(fmt.Sprintf("%s", strconv.Itoa(10000*l.c.LimitCpuCore))), 0644); err != nil {
		return err
	}
	return
}

// pathExists check if path exist
func pathExists(path string, create bool) (err error) {
	if _, err = os.Stat(path); err == nil {
		return
	}

	if os.IsNotExist(err) && create == true {
		return os.MkdirAll(path, os.ModePerm)
	}
	return
}
