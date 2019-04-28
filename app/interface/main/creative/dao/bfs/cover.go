package bfs

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"fmt"
	"go-common/library/log"
)

// UpVideoCovers upload video covers.
func (d *Dao) UpVideoCovers(c context.Context, covers []string) (cvs []string, err error) {
	var (
		nfsURL string
	)
	for _, cv := range covers {
		// get nfs file
		bs, err := d.bvcCover(cv)
		if err != nil || len(bs) == 0 {
			log.Error("d.UpVideoCovers(%s) error(%v) or bs==0", nfsURL, err)
			continue
		}
		// up to bfs
		bfsPath, err := d.UploadArc(c, http.DetectContentType(bs), bytes.NewReader(bs))
		if err != nil {
			log.Error("d.UpVideoCovers raw url(%s) error(%v)", cv, err)
			continue
		}
		// parse bfs return path
		if err != nil {
			log.Error("url.Parse(%v) error(%v)", bfsPath, err)
			continue
		}
		cvs = append(cvs, bfsPath)
		log.Info("UpVideoCovers cover(%s) bfs (%s)", cv, bfsPath)
	}
	return
}

// bvcCover http get bvc cover bytes.
func (d *Dao) bvcCover(url string) (bs []byte, err error) {
	resp, err := d.client.Get(url)
	if err != nil {
		log.Error("s.client.Get(%v) error(%v)", url, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("get NFS file faild, url(%s) http status: %d", url, resp.StatusCode)
		return
	}
	// read bytes
	if bs, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error("ioutil.ReadAll error(%v)", err)
	}
	return
}
