package dao

import (
	"context"
	"strings"
	"time"

	"github.com/Dai0522/go-hash/bloomfilter"
)

// FetchMidView .
func (d *Dao) FetchMidView(c context.Context) (result []string, err error) {
	dt := time.Now().AddDate(0, 0, -1).Format("20060102")
	hdfs, err := d.scanHDFSPath(c, d.c.Berserker.API[1].URL, d.c.Berserker.Keys[0], "/"+dt+"/mid/")
	if err != nil {
		return
	}
	for _, v := range hdfs.Result {
		raw, err := d.loadHDFSFile(c, d.c.Berserker.API[1].URL, d.c.Berserker.Keys[0], "/"+dt+"/mid/"+v)
		if err != nil {
			break
		}
		lines := strings.Split(string(*raw), "\n")
		result = append(result, lines...)
	}
	return
}

// FetchBuvidView .
func (d *Dao) FetchBuvidView(c context.Context) (result []string, err error) {
	dt := time.Now().AddDate(0, 0, -1).Format("20060102")
	hdfs, err := d.scanHDFSPath(c, d.c.Berserker.API[1].URL, d.c.Berserker.Keys[0], "/"+dt+"/buvid/")
	if err != nil {
		return
	}
	for _, v := range hdfs.Result {
		raw, err := d.loadHDFSFile(c, d.c.Berserker.API[1].URL, d.c.Berserker.Keys[0], "/"+dt+"/buvid/"+v)
		if err != nil {
			break
		}
		lines := strings.Split(string(*raw), "\n")
		result = append(result, lines...)
	}
	return
}

// InsertBloomFilter 构建BF，插入redis
func (d *Dao) InsertBloomFilter(c context.Context, key string, svidList []uint64) error {
	bfK := "BBQ:BF:V1:" + key
	bf, err := bloomfilter.New(uint64(len(svidList)), 0.0001)
	if err != nil {
		return err
	}
	for _, v := range svidList {
		bf.PutUint64(v)
	}

	b := bf.Serialized()
	return d.SetBloomFilter(c, bfK, b)
}

// SetBloomFilter .
func (d *Dao) SetBloomFilter(c context.Context, key string, b *[]byte) error {
	conn := d.bfredis.Get(c)
	defer conn.Close()
	_, err := conn.Do("SETEX", []byte(key), 86400, *b)
	return err
}
