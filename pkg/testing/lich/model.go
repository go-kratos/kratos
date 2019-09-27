package lich

import (
	"strings"
	"time"
)

// Container docker inspect resp.
type Container struct {
	ID      string    `json:"Id"`
	Created time.Time `json:"Created"`
	Path    string    `json:"Path"`
	Args    []string  `json:"Args"`
	State   struct {
		Status     string    `json:"Status"`
		Running    bool      `json:"Running"`
		Paused     bool      `json:"Paused"`
		Restarting bool      `json:"Restarting"`
		OOMKilled  bool      `json:"OOMKilled"`
		Dead       bool      `json:"Dead"`
		Pid        int       `json:"Pid"`
		ExitCode   int       `json:"ExitCode"`
		Error      string    `json:"Error"`
		StartedAt  time.Time `json:"StartedAt"`
		FinishedAt time.Time `json:"FinishedAt"`
		Health     struct {
			Status        string `json:"Status"`
			FailingStreak int    `json:"FailingStreak"`
			Log           []struct {
				Start    time.Time `json:"Start"`
				End      time.Time `json:"End"`
				ExitCode int       `json:"ExitCode"`
				Output   string    `json:"Output"`
			} `json:"Log"`
		} `json:"Health"`
	} `json:"State"`
	Config struct {
		Hostname   string   `json:"Hostname"`
		Domainname string   `json:"Domainname"`
		User       string   `json:"User"`
		Tty        bool     `json:"Tty"`
		OpenStdin  bool     `json:"OpenStdin"`
		StdinOnce  bool     `json:"StdinOnce"`
		Env        []string `json:"Env"`
		Cmd        []string `json:"Cmd"`
		Image      string   `json:"Image"`
		WorkingDir string   `json:"WorkingDir"`
		Entrypoint []string `json:"Entrypoint"`
	} `json:"Config"`
	Image           string `json:"Image"`
	ResolvConfPath  string `json:"ResolvConfPath"`
	HostnamePath    string `json:"HostnamePath"`
	HostsPath       string `json:"HostsPath"`
	LogPath         string `json:"LogPath"`
	Name            string `json:"Name"`
	RestartCount    int    `json:"RestartCount"`
	Driver          string `json:"Driver"`
	Platform        string `json:"Platform"`
	MountLabel      string `json:"MountLabel"`
	ProcessLabel    string `json:"ProcessLabel"`
	AppArmorProfile string `json:"AppArmorProfile"`
	NetworkSettings struct {
		Bridge      string `json:"Bridge"`
		SandboxID   string `json:"SandboxID"`
		HairpinMode bool   `json:"HairpinMode"`
		Ports       map[string][]struct {
			HostIP   string `json:"HostIp"`
			HostPort string `json:"HostPort"`
		} `json:"Ports"`
	} `json:"NetworkSettings"`
}

// GetImage get image name at container
func (c *Container) GetImage() (image string) {
	image = c.Config.Image
	if images := strings.Split(image, ":"); len(images) > 0 {
		image = images[0]
	}
	return
}

// GetID get id at container
func (c *Container) GetID() (id string) {
	if id = c.ID; len(id) > 9 {
		id = id[0:9]
	}
	return
}
