package lich

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	retry    int
	yamlPath string
	pathHash string
	services map[string]*Container
)

func init() {
	flag.StringVar(&yamlPath, "f", "docker-compose.yaml", "composer yaml path.")
}

// Setup setup UT related environment dependence for everything.
func Setup() (err error) {
	if _, err = os.Stat(yamlPath); os.IsNotExist(err) {
		log.Println("composer yaml is not exist!", yamlPath)
		return
	}
	if yamlPath, err = filepath.Abs(yamlPath); err != nil {
		log.Printf("filepath.Abs(%s) error(%v)", yamlPath, err)
		return
	}
	pathHash = fmt.Sprintf("%x", md5.Sum([]byte(yamlPath)))[:9]
	var args = []string{"-f", yamlPath, "-p", pathHash, "up", "-d"}
	if err = exec.Command("docker-compose", args...).Run(); err != nil {
		log.Printf("exec.Command(docker-composer) args(%v) error(%v)", args, err)
		Teardown()
		return
	}
	// 拿到yaml文件中的服务名，同时通过服务名获取到启动的容器ID
	if _, err = getServices(); err != nil {
		Teardown()
		return
	}
	// 通过容器ID检测容器的状态，包括容器服务的状态
	if _, err = checkServices(); err != nil {
		Teardown()
		return
	}
	return
}

// Teardown unsetup all environment dependence.
func Teardown() (err error) {
	if _, err = os.Stat(yamlPath); os.IsNotExist(err) {
		log.Println("composer yaml is not exist!")
		return
	}
	if yamlPath, err = filepath.Abs(yamlPath); err != nil {
		log.Printf("filepath.Abs(%s) error(%v)", yamlPath, err)
		return
	}
	pathHash = fmt.Sprintf("%x", md5.Sum([]byte(yamlPath)))[:9]
	args := []string{"-f", yamlPath, "-p", pathHash, "down"}
	if output, err := exec.Command("docker-compose", args...).CombinedOutput(); err != nil {
		log.Fatalf("exec.Command(docker-composer) args(%v) stdout(%s) error(%v)", args, string(output), err)
		return err
	}
	return
}

func getServices() (output []byte, err error) {
	var args = []string{"-f", yamlPath, "-p", pathHash, "config", "--services"}
	if output, err = exec.Command("docker-compose", args...).CombinedOutput(); err != nil {
		log.Printf("exec.Command(docker-composer) args(%v) stdout(%s) error(%v)", args, string(output), err)
		return
	}
	services = make(map[string]*Container)
	output = bytes.TrimSpace(output)
	for _, svr := range bytes.Split(output, []byte("\n")) {
		args = []string{"-f", yamlPath, "-p", pathHash, "ps", "-a", "-q", string(svr)}
		if output, err = exec.Command("docker-compose", args...).CombinedOutput(); err != nil {
			log.Printf("exec.Command(docker-composer) args(%v) stdout(%s) error(%v)", args, string(output), err)
			return
		}
		var id = string(bytes.TrimSpace(output))
		args = []string{"inspect", id, "--format", "'{{json .}}'"}
		if output, err = exec.Command("docker", args...).CombinedOutput(); err != nil {
			log.Printf("exec.Command(docker) args(%v) stdout(%s) error(%v)", args, string(output), err)
			return
		}
		if output = bytes.TrimSpace(output); bytes.Equal(output, []byte("")) {
			err = fmt.Errorf("service: %s | container: %s fails to launch", svr, id)
			log.Printf("exec.Command(docker) args(%v) error(%v)", args, err)
			return
		}
		var c = &Container{}
		if err = json.Unmarshal(bytes.Trim(output, "'"), c); err != nil {
			log.Printf("json.Unmarshal(%s) error(%v)", string(output), err)
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
			log.Printf("healthcheck(%s) error(%v) retrying %d times...", svr, err, 5-retry)
			return
		}
		// TODO About container check and more...
	}
	return
}
