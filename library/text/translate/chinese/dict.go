package chinese

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go-common/library/log"

	"github.com/go-ego/cedar"
)

// dict contains the Trie and dict values
type dict struct {
	Trie   *cedar.Cedar
	Values [][]string
}

// BuildFromFile builds the da dict from fileName
func buildFromFile(fileName string) (*dict, error) {
	var err error
	trie := cedar.New()
	values := [][]string{}
	bs := raw(fileName)
	strs := strings.Split(string(bs), "\n")
	for _, line := range strs {
		items := strings.Split(strings.TrimSpace(line), "\t")
		if len(items) < 2 {
			continue
		}
		err = trie.Insert([]byte(items[0]), len(values))
		if err != nil {
			return nil, err
		}
		if len(items) > 2 {
			values = append(values, items[1:])
		} else {
			values = append(values, strings.Fields(items[1]))
		}
	}
	return &dict{Trie: trie, Values: values}, nil
}

// prefixMatch str by Dict, returns the matched string and its according values
func (d *dict) prefixMatch(str string) (map[string][]string, error) {
	if d.Trie == nil {
		return nil, fmt.Errorf("Trie is nil")
	}
	res := make(map[string][]string)
	for _, id := range d.Trie.PrefixMatch([]byte(str), 0) {
		key, err := d.Trie.Key(id)
		if err != nil {
			return nil, err
		}
		value, err := d.Trie.Value(id)
		if err != nil {
			return nil, err
		}
		res[string(key)] = d.Values[value]
	}
	return res, nil
}

var (
	defaultRead int64 = 16 * 1024 // 16kb
	defaultURL        = "http://i0.hdslb.com/bfs/static/"
)

func raw(file string) (bs []byte) {
	client := http.Client{Timeout: 10 * time.Second}
	for i := 0; i < 3; i++ {
		resp, err := client.Get(defaultURL + file)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Error("bfs client url:%s file:%s err:%+v", defaultURL, file, err)
			time.Sleep(time.Millisecond * 50)
			continue
		}
		defer resp.Body.Close()
		bs, err = readAll(resp.Body, defaultRead)
		if err == nil {
			return
		}
		log.Error("bfs client url:%s file:%s err:%+v", defaultURL, file, err)
	}
	return
}

func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
