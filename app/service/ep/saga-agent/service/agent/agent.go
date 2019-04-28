package agent

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime/debug"
	"strings"

	"go-common/app/service/ep/saga-agent/conf"
	"go-common/library/log"

	"github.com/pkg/errors"
)

var (
	hostName string
)

const _runnerWorkDir = "/data/gitlab-runner"

func shellExec(dir string, name string, args ...string) (stdout string, stderr string, err error) {
	defer func() {
		log.Info("Shell Exec dir[%s] name[%s] args[%+v] stdout[%s] stderr[%s] err[%+v]", dir, name, args, stdout, stderr, err)
	}()
	if _, err = exec.LookPath(name); err != nil {
		err = errors.Wrapf(err, "LookPath(%s) failed", name)
		return
	}
	var (
		stdo = new(bytes.Buffer)
		stde = new(bytes.Buffer)
		cmd  *exec.Cmd
	)
	cmd = exec.Command(name, args...)
	cmd.Stdout = stdo
	cmd.Stderr = stde
	cmd.Dir = dir
	if err = cmd.Run(); err != nil {
		err = errors.Wrapf(err, "cmd.Run(%s,%s,%v) failed", dir, name, args)
	}
	stdout = stdo.String()
	stderr = stde.String()
	return
}

// RunnerRegister ...
func RunnerRegister(url, token, name string) (stdout string, stderr string, err error) {
	var args string

	if (url == "") || (token == "") {
		err = fmt.Errorf("RunnerRegister invalid parameters(url:%s token:%s)", url, token)
		log.Error(err.Error())
		return
	}

	log.Info("gitlab-runner register url:%s token:%s", url, token)
	args = fmt.Sprintf("gitlab-runner register -u %s -r %s --name %s --executor shell -n true -c /data/gitlab-runner/config.toml",
		url, token, name)
	if stdout, stderr, err = shellExec(_runnerWorkDir, "/bin/sh", "-c", args); err != nil {
		log.Error(strings.TrimSpace(stdout + "\n" + stderr))
	}
	return
}

// RunnerUnRegister ...
func RunnerUnRegister(name string) (stdout string, stderr string, err error) {
	var args string
	log.Info("gitlab-runner unregister name:%s", name)
	args = fmt.Sprintf("gitlab-runner unregister -n %s", name)
	if stdout, stderr, err = shellExec(_runnerWorkDir, "/bin/sh", "-c", args); err != nil {
		log.Error(strings.TrimSpace(stdout + "\n" + stderr))
	}
	return
}

// RunnerUnRegisterAll ...
func RunnerUnRegisterAll() (stdout string, stderr string, err error) {
	var args string
	log.Info("gitlab-runner RunnerUnregisterAll")
	args = "gitlab-runner unregister --all-runners"
	if stdout, stderr, err = shellExec(_runnerWorkDir, "/bin/sh", "-c", args); err != nil {
		log.Error(strings.TrimSpace(stdout + "\n" + stderr))
	}
	return
}

// RunnerStart ...
func RunnerStart() (stdout string, stderr string, err error) {
	var args string
	log.Info("gitlab-runner start...")
	args = fmt.Sprintf("gitlab-runner run -d %s", _runnerWorkDir)
	if stdout, stderr, err = shellExec(_runnerWorkDir, "/bin/sh", "-c", args); err != nil {
		log.Error(strings.TrimSpace(stdout + "\n" + stderr))
	}
	return

}

// ExecRegister ...
func ExecRegister() {
	var tmpMap map[string]conf.Runner

	defer func() {
		if x := recover(); x != nil {
			log.Error("execRegister: %+v %s", x, debug.Stack())
		}
	}()

	log.Info("execRegister... runner.offline %t,runner.len:%d", conf.Conf.Offline, len(conf.Conf.Runner))
	if conf.Conf.Offline {
		conf.RunnerMap = make(map[string]conf.Runner)
		RunnerUnRegisterAll()
		return
	}
	tmpMap = make(map[string]conf.Runner)
	for _, v := range conf.Conf.Runner {
		tmpMap[v.Token] = v
		_, ok := conf.RunnerMap[v.Token]
		if !ok {
			RunnerRegister(v.URL, v.Token, hostName+"-"+v.Name)
			conf.RunnerMap[v.Token] = v
		}
	}
	for _, v := range conf.RunnerMap {
		_, ok := tmpMap[v.Token]
		if !ok {
			RunnerUnRegister(hostName + "-" + v.Name)
			delete(conf.RunnerMap, v.Token)
		}
	}
}

// UpdateRegister ...
func UpdateRegister() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("updateRegister: %+v %s", x, debug.Stack())
		}
	}()

	for range conf.ReloadEvent() {
		log.Info("updateRegister")
		ExecRegister()
	}
}
