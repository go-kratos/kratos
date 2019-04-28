package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/urfave/cli"
)

type archs struct {
	LinuxAmd64  string `json:"linux-amd64"`
	DarwinAmd64 string `json:"darwin_amd64"`
}

type internalInfo struct {
	version    int
	maxVersion int
	up         map[string]archs
}

func upgradeAction(c *cli.Context) error {
	upgrade()
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	out.Chmod(0755)
	return out.Close()
}

func updateFile(sha, url string) {
	ex, err := os.Executable()
	if err != nil {
		fmt.Printf("fail to get download path")
		return
	}
	fpath := filepath.Dir(ex)
	tmpFilepath := fpath + "/kratos.tmp"

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("fail to download file: %v", err)
		return
	}

	out, err := os.OpenFile(tmpFilepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("fail to write download file")
		return
	}
	out.Close()
	resp.Body.Close()

	f, err := os.Open(tmpFilepath)
	if err != nil {
		fmt.Println("sha256 fail to open download file")
		return
	}

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		fmt.Println("sha256 fail to init")
		return
	}
	getsha := fmt.Sprintf("%x", h.Sum(nil))
	if getsha != sha {
		fmt.Printf("sha256 wrong. expect %s get %s", sha, getsha)
		return
	}
	f.Close()
	err = copyFile(tmpFilepath, fpath+"/kratos")
	if err != nil {
		fmt.Println("fail to install kratos")
		return
	}
	err = os.Remove(tmpFilepath)
	if err != nil {
		fmt.Println("fail to remove tmp kratos")
		return
	}
	fmt.Println("Download successfully!")
}

func upgrade() error {
	target := make(map[string]archs)
	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get("http://bazel-cabin.bilibili.co/kratos/" + Channel + "/package.json")
	if err != nil {
		return err
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&target)
	info := internalInfo{}
	info.up = target
	if info.version, err = strconv.Atoi(Version); err != nil {
		return err
	}
	for k := range target {
		ver, err := strconv.Atoi(k)
		if err != nil {
			return err
		}
		if info.maxVersion < ver {
			info.maxVersion = ver
		}
	}
	if info.maxVersion > info.version {
		fmt.Printf("kratos %d -> %d\n", info.version, info.maxVersion)
		switch runtime.GOOS + "-" + runtime.GOARCH {
		case "linux-amd64":
			updateFile(info.up[strconv.Itoa(info.maxVersion)].LinuxAmd64, "http://bazel-cabin.bilibili.co/kratos/"+Channel+"/"+strconv.Itoa(info.maxVersion)+"/linux-amd64/kratos")
		case "darwin-amd64":
			updateFile(info.up[strconv.Itoa(info.maxVersion)].DarwinAmd64, "http://bazel-cabin.bilibili.co/kratos/"+Channel+"/"+strconv.Itoa(info.maxVersion)+"/darwin-amd64/kratos")
		default:
			fmt.Println("not support this operate system")
		}
	} else {
		fmt.Println("Already up to the newest.")
	}
	return nil
}
