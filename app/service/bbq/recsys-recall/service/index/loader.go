package index

import (
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"go-common/app/job/bbq/recall/proto"
	"go-common/library/log"

	"github.com/golang/snappy"
)

// Loader .
type Loader interface {
	Load() (*map[uint64]*proto.ForwardIndex, error)
}

// LocalLoader .
type LocalLoader struct {
	path string
}

// Load .
func (l *LocalLoader) Load() (result *map[uint64]*proto.ForwardIndex, err error) {
	data := make(map[uint64]*proto.ForwardIndex)
	f, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	items := strings.Split(string(b), "\n")
	for _, v := range items {
		if v == "" {
			continue
		}
		compressed, err := hex.DecodeString(v)
		if err != nil {
			log.Error("hex decode: [%v] [%s]", err, string(v))
			continue
		}
		raw, err := snappy.Decode(nil, compressed)
		if err != nil {
			log.Error("snappy decode: [%v] [%s]", err, string(v))
			continue
		}
		tmp := &proto.ForwardIndex{}
		err = tmp.Unmarshal(raw)
		if err != nil {
			log.Error("proto decode: [%v] [%s]", err, string(v))
			continue
		}
		data[tmp.SVID] = tmp
	}
	result = &data

	return
}

// RemoteLoader .
type RemoteLoader struct {
	path    string
	md5Path string
	md5     string
}

// Load .
func (l *RemoteLoader) Load() (result *map[uint64]*proto.ForwardIndex, err error) {
	md5, err := l.loadMD5()
	if err != nil || md5 == l.md5 {
		return
	}

	data := make(map[uint64]*proto.ForwardIndex)
	resp, err := http.DefaultClient.Get(l.path)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	items := strings.Split(string(b), "\n")
	for _, v := range items {
		raw, err := hex.DecodeString(v)
		if err != nil {
			log.Error("hex decode: %v", err)
			continue
		}
		raw, err = snappy.Decode(nil, raw)
		if err != nil {
			log.Error("snappy decode: %v", err)
			continue
		}
		tmp := &proto.ForwardIndex{}
		err = tmp.Unmarshal(raw)
		if err != nil {
			continue
		}
		data[tmp.SVID] = tmp
	}
	result = &data
	l.md5 = md5

	return
}

func (l *RemoteLoader) loadMD5() (string, error) {
	resp, err := http.DefaultClient.Get(l.md5Path)
	if err != nil {
		return "", err
	}
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(raw), err
}
