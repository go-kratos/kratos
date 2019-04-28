package dao

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
)

//Capture performs a HTTP Get request for the image url and upload bfs.
func (d *Dao) Capture(c context.Context, url string) (loc string, size int, err error) {
	if err = checkURL(url); err != nil {
		return
	}
	bs, ct, err := d.download(c, url)
	if err != nil {
		return
	}
	size = len(bs)
	if size == 0 {
		log.Error("capture image size(%d)|url(%s)", size, url)
		return
	}
	if ct != "image/jpeg" && ct != "image/jpg" && ct != "image/png" && ct != "image/gif" {
		log.Error("capture not allow image file type(%s)", ct)
		err = ecode.CreativeArticleImageTypeErr
		return
	}
	loc, err = d.UploadImage(c, ct, bs)
	return loc, size, err
}

func checkURL(url string) (err error) {
	// http || https
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		log.Error("capture url invalid(%s)", url)
		err = ecode.RequestErr
		return
	}
	u, err := nurl.Parse(url)
	if err != nil {
		log.Error("capture url.Parse error(%v)", err)
		err = ecode.RequestErr
		return
	}
	// make sure ip is public. avoid ssrf
	ips, err := net.LookupIP(u.Host) // take from 1st argument
	if err != nil {
		log.Error("capture url(%s) LookupIP failed", url)
		err = ecode.RequestErr
		return
	}
	if len(ips) == 0 {
		log.Error("capture url(%s) LookupIP length 0", url)
		err = ecode.RequestErr
		return
	}
	for _, v := range ips {
		if !isPublicIP(v) {
			log.Error("capture url(%s) is not public ip(%v)", url, v)
			err = ecode.RequestErr
			return
		}
	}
	return
}

func isPublicIP(IP net.IP) bool {
	if env.DeployEnv == env.DeployEnvDev || env.DeployEnv == env.DeployEnvFat1 || env.DeployEnv == env.DeployEnvUat {
		return true
	}
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func (d *Dao) download(c context.Context, url string) (bs []byte, ct string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("capture http.NewRequest error(%v)|url (%s)", err, url)
		return
	}
	// timeout
	ctx, cancel := context.WithTimeout(c, 800*time.Millisecond)
	req = req.WithContext(ctx)
	defer cancel()
	resp, err := d.bfsClient.Do(req)
	if err != nil {
		log.Error("capture d.client.Do error(%v)|url(%s)", err, url)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error("capture http.StatusCode nq http.StatusOK(%d)|url(%s)", resp.StatusCode, url)
		err = errors.New("Download out image link failed")
		return
	}
	if bs, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error("capture ioutil.ReadAll error(%v)", err)
		err = errors.New("Download out image link failed")
		return
	}
	ct = http.DetectContentType(bs)
	return
}
