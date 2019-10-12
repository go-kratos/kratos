package lich

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/bilibili/kratos/pkg/log"
)

var (
	retry    int
	noDown   bool
	yamlPath string
	pathHash string
	services map[string]*Container
)

func init() {
	flag.StringVar(&yamlPath, "f", "docker-compose.yaml", "composer yaml path.")
	flag.BoolVar(&noDown, "nodown", false, "containers are not recycled.")
}

func runCompose(args ...string) (output []byte, err error) {
	if _, err = os.Stat(yamlPath); os.IsNotExist(err) {
		log.Error("os.Stat(%s) composer yaml is not exist!", yamlPath)
		return
	}
	if yamlPath, err = filepath.Abs(yamlPath); err != nil {
		log.Error("filepath.Abs(%s) error(%v)", yamlPath, err)
		return
	}
	pathHash = fmt.Sprintf("%x", md5.Sum([]byte(yamlPath)))[:9]
	args = append([]string{"-f", yamlPath, "-p", pathHash}, args...)
	if output, err = exec.Command("docker-compose", args...).CombinedOutput(); err != nil {
		log.Error("exec.Command(docker-compose) args(%v) stdout(%s) error(%v)", args, string(output), err)
		return
	}
	return
}

// Setup setup UT related environment dependence for everything.
func Setup() (err error) {
	if _, err = runCompose("up", "-d"); err != nil {
		return
	}
	defer func() {
		if err != nil {
			Teardown()
		}
	}()
	if _, err = getServices(); err != nil {
		return
	}
	_, err = checkServices()
	return
}

// Teardown unsetup all environment dependence.
func Teardown() (err error) {
	if !noDown {
		_, err = runCompose("down")
	}
	return
}

func getServices() (output []byte, err error) {
	if output, err = runCompose("config", "--services"); err != nil {
		return
	}
	services = make(map[string]*Container)
	output = bytes.TrimSpace(output)
	for _, svr := range bytes.Split(output, []byte("\n")) {
		if output, err = runCompose("ps", "-a", "-q", string(svr)); err != nil {
			return
		}
		var (
			id   = string(bytes.TrimSpace(output))
			args = []string{"inspect", id, "--format", "'{{json .}}'"}
		)
		if output, err = exec.Command("docker", args...).CombinedOutput(); err != nil {
			log.Error("exec.Command(docker) args(%v) stdout(%s) error(%v)", args, string(output), err)
			return
		}
		if output = bytes.TrimSpace(output); bytes.Equal(output, []byte("")) {
			err = fmt.Errorf("service: %s | container: %s fails to launch", svr, id)
			log.Error("exec.Command(docker) args(%v) error(%v)", args, err)
			return
		}
		var c = &Container{}
		if err = json.Unmarshal(bytes.Trim(output, "'"), c); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(output), err)
			return
		}
		services[string(svr)] = c
	}
	return
}

func checkServices() (output []byte, err error) {
	defer func() {
		if err != nil && retry < 4 {
			retry++
			getServices()
			time.Sleep(time.Second * 5)
			output, err = checkServices()
			return
		}
		retry = 0
	}()
	for svr, c := range services {
		if err = c.Healthcheck(); err != nil {
			log.Error("healthcheck(%s) error(%v) retrying %d times...", svr, err, 5-retry)
			return
		}
		// TODO About container check and more...
	}
	return
}
