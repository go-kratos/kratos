package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// GoGet go get path.
func GoGet(path ...string) error {
	for _, p := range path {
		fmt.Printf("go get -u %s\n", p)
		cmd := exec.Command("go", "get", "-u", p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

type ReleaseInfo struct {
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	PublishedAt string `json:"published_at"`
	Body        string `json:"body"`
	HtmlUrl     string `json:"html_url"`
}

type CommitInfo struct {
	Sha    string `json:"sha,omitempty"`
	Commit struct {
		Author struct {
			Date string `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	HtmlUrl string `json:"html_url"`
}

type ErrorInfo struct {
	Message string
}

// GetReleaseInfo for getting kratos release info
func GetReleaseInfo(version string) (*ReleaseInfo, error) {
	cli := http.Client{Timeout: 60 * time.Second}
	api := "https://api.github.com/repos/go-kratos/kratos/releases/latest"
	if version != "latest" {
		api = fmt.Sprintf("https://api.github.com/repos/go-kratos/kratos/releases/tags/%s", version)
	}
	resp, err := cli.Get(api)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		errorInfo := &ErrorInfo{}
		err = json.Unmarshal(body, errorInfo)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorInfo.Message)
	}
	releaseInfo := &ReleaseInfo{}
	err = json.Unmarshal(body, releaseInfo)
	return releaseInfo, err
}

// GetCommitsInfo for getting kratos commits info
func GetCommitsInfo() ([]CommitInfo, error) {
	// Get the latest release time
	info, err := GetReleaseInfo("latest")
	if err != nil {
		return nil, err
	}
	cli := &http.Client{Timeout: 60 * time.Second}
	resp, err := cli.Get(fmt.Sprintf("https://api.github.com/repos/go-kratos/kratos/commits?per_page=5&since=%s", info.PublishedAt))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		errorInfo := &ErrorInfo{}
		err = json.Unmarshal(body, errorInfo)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorInfo.Message)
	}
	var res []CommitInfo
	err = json.Unmarshal(body, &res)
	return res, err
}
