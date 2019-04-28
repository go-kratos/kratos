package service

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/siddontang/go/ioutil2"

	"go-common/app/service/live/recommend/recconst"

	"go-common/app/job/live/recommend-job/internal/conf"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/sync/errgroup"

	"github.com/pkg/errors"
)

// ItemCFJob 把 hadoop 的数据写到redis
type ItemCFJob struct {
	Conf       *conf.JobConfig
	RedisConf  *redis.Config
	HadoopConf *conf.HadoopConfig
}

// Run ...
func (j *ItemCFJob) Run() {
	log.Info("Run ItemCFJob")
	processFile(j.Conf, j.HadoopConf, j.RedisConf, writeItemCFToRedis)
	log.Info("ItemCFJob Done")
}

func fetchFiles(jobConf *conf.JobConfig,
	hadoopConf *conf.HadoopConfig) (paths []string, err error) {
	path := jobConf.InputFile
	if path == "" {
		var javaHome string
		var hadoopHome string
		hadoopHome, javaHome, err = downloadHadoop(hadoopConf)
		if err != nil {
			log.Info("download hadoop err %+v ", err)
			return
		}
		log.Info("download hadoop success: " + hadoopHome)
		path, err = downloadFileFromHadoop(jobConf, hadoopHome, javaHome)
		if err != nil {
			log.Info("path %s downloaded err %+v  ", path, err)
			return
		}
		log.Info("path downloaded success: " + path)

		var file os.FileInfo
		file, err = os.Stat(path)
		if err != nil {
			log.Error("cannot open file %s err: %+v", file, err)
			return
		}
		if file.IsDir() {
			var files []os.FileInfo
			files, err = ioutil.ReadDir(path)
			if err != nil {
				log.Error("no file in dir: %d, err: %+v", path, err)
				return
			}

			for _, f := range files {
				if strings.Index(f.Name(), ".") != 0 && strings.Index(f.Name(), "_") != 0 {
					paths = append(paths, path+"/"+f.Name())
				}
			}
		} else {
			paths = []string{path}
		}
	} else {
		if strings.Index(path, "http://") == 0 {
			var httpPath = path
			path, err = downloadFileFromHttp(httpPath, "/tmp/job-downloaded.txt")
			if err != nil {
				log.Error("download from http path=%s, error=%+v", httpPath, err)
				return
			}
			log.Info("file downloaded from http %s to %s", httpPath, path)
		} else {
			var file os.FileInfo
			file, err = os.Stat(path)
			if err != nil {
				log.Error("cannot open file %s err: %+v", file, err)
				return
			}
			if file.IsDir() {
				var files []os.FileInfo
				files, err = ioutil.ReadDir(path)
				if err != nil {
					log.Error("no file in dir: %d, err: %+v", path, err)
					return
				}

				for _, f := range files {
					if strings.Index(f.Name(), ".") != 0 && strings.Index(f.Name(), "_") != 0 {
						paths = append(paths, path+"/"+f.Name())
					}
				}
			} else {
				paths = []string{path}
			}
		}
	}
	return
}

// 按行读取文件，并调用handler处理
func processFile(jobConf *conf.JobConfig,
	hadoopConf *conf.HadoopConfig,
	redisConf *redis.Config,
	handler func(line string, pool *redis.Pool) error,
) {
	paths, err := fetchFiles(jobConf, hadoopConf)
	if err != nil {
		return
	}
	var workerNum = jobConf.WorkerNum
	var r = redis.NewPool(redisConf)
	defer r.Close()

	log.Info("all of files %+v", paths)
	for _, path := range paths {

		var startLineNum int

		file, err := os.Open(path)
		if err != nil {
			log.Error("open path %s err %+v", path, errors.WithStack(err))
			panic(err)
		}

		name := filepath.Base(filepath.Dir(path)) + "-" + filepath.Base(path)
		offsetPath := "/tmp/" + name + ".offset"
		// 读取上一次的位置
		contentB, e := ioutil.ReadFile(offsetPath)
		if e == nil {
			content := string(contentB)
			offset, e := strconv.Atoi(content)
			if e == nil {
				startLineNum = offset
			} else {
				startLineNum = 1
			}
		} else {
			startLineNum = 1
		}
		log.Info("start from line: %d, file : %s offset file: %s", startLineNum, path, offsetPath)

		var ops uint64
		totalCount := lineCounter(path)
		scanner := bufio.NewScanner(file)

		g := errgroup.Group{}
		g.GOMAXPROCS(workerNum)

		var lineNum = 0
		for scanner.Scan() {
			lineNum++
			if lineNum < startLineNum {
				continue
			}

			line := scanner.Text() // 10001 [1:0.9,2:0.9]
			g.Go(func() error {
				handler(line, r)
				atomic.AddUint64(&ops, 1)
				if ops%20000 == 0 {
					fmt.Printf("progress %d / %d percent %f %s         \r", ops+uint64(startLineNum)-1, totalCount,
						float32(uint64(startLineNum)+ops-1)/float32(totalCount)*100, "%")
				} else if ops == uint64(totalCount) {
					fmt.Printf("progress %d / %d 100%%\n", ops, totalCount)
				}
				return nil
			})
			if lineNum%100000 == 0 {
				g.Wait()
				ioutil.WriteFile(offsetPath,
					[]byte(strconv.FormatInt(int64(lineNum), 10)),
					os.ModePerm)
				g = errgroup.Group{}
				g.GOMAXPROCS(workerNum)
			}
		}

		g.Wait()

		if err = scanner.Err(); err != nil {
			log.Error("err %+v", errors.WithStack(err))
			panic(err)
		}
		file.Close()
	}

}

func writeItemCFToRedis(line string, r *redis.Pool) (err error) {
	start := strings.Index(line, "[")
	end := strings.LastIndex(line, "]")
	userIdStr := line[0 : start-1]
	items := strings.Split(line[start+1:end], ",")
	c := r.Get(context.Background())
	defer c.Close()
	userId, _ := strconv.Atoi(userIdStr)
	key := fmt.Sprintf(recconst.UserItemCFRecKey, userId)
	// 最多保留50个推荐
	var inMemoryCount int
	inMemoryCount, err = redis.Int(c.Do("ZCARD", key))
	if err != nil {
		log.Error("zcard err: %+v", err)
	} else {
		var toBeRemovedCount = inMemoryCount + len(items) - 60
		if toBeRemovedCount > 0 {
			var removed int
			removed, err = redis.Int(c.Do("ZREMRANGEBYRANK", key, 0, toBeRemovedCount-1))
			if err != nil {
				log.Error("ZREMRANGEBYRANK key:%s, err: +%v", key, err)
			} else {
				log.Info("zset removed %d count, key:%s", removed, key)
			}
		}
	}
	for _, item := range items {
		split := strings.Split(item, ":")
		itemID := split[0]
		score := split[1]

		c.Send("ZADD", key, score, itemID)
	}
	c.Send("EXPIRE", key, 86400*30)
	err = c.Flush()

	if err != nil {
		log.Error("zadd to redis error: %+v , key=%s", err, key)
		return err
	}
	for i := 0; i < len(items)+1; i++ {
		_, err = c.Receive()
		if err != nil {
			log.Error("zadd to redis error: %+v , key=%s, line=%s", err, key, line)
			return err
		}
	}
	return nil
}

func lineCounter(path string) int {
	buf := make([]byte, 32*1024)
	r, _ := os.Open(path)
	defer r.Close()
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count

		case err != nil:
			return count
		}
	}
}

func downloadHadoop(hadoopConf *conf.HadoopConfig) (hadoopHome string, javaHome string, err error) {
	if strings.LastIndex(hadoopConf.HadoopDir, "/") == len(hadoopConf.HadoopDir)-1 {
		hadoopHome = hadoopConf.HadoopDir + "hadoop-2.8.4"
	} else {
		hadoopHome = hadoopConf.HadoopDir + "/hadoop-2.8.4"
	}
	javaHome = hadoopHome + "/jdk1.8.0_60"

	if ioutil2.FileExists(hadoopHome) {
		return
	}
	var cmd = "curl -sSLf " + hadoopConf.TarUrl + " -o /tmp/hadoop.tar.gz"
	err = runCmd(cmd)
	if err != nil {
		return
	}
	cmd = "tar -C " + hadoopConf.HadoopDir + " -xf /tmp/hadoop.tar.gz"
	err = runCmd(cmd)
	if err != nil {
		return
	}

	return
}

func downloadFileFromHttp(url string, output string) (string, error) {
	var localPath = output
	var cmd = "curl -sSLf " + url + " -o " + localPath
	var err = runCmd(cmd)
	if err != nil {
		return "", err
	}
	return localPath, nil
}

func downloadFileFromHadoop(jobConf *conf.JobConfig, hadoopHome string, javaHome string) (string, error) {
	t := time.Now().AddDate(0, 0, -1)
	day := t.Format("20060102")
	localPath := fmt.Sprintf(jobConf.LocalTmpFile, day)
	if ioutil2.FileExists(localPath) {
		return localPath, nil
	}
	remotePath := fmt.Sprintf(jobConf.HadoopFile, day)
	cmd := fmt.Sprintf("export JAVA_HOME=%s; %s/bin/hdfs dfs -get %s %s",
		javaHome, hadoopHome, remotePath, localPath)
	err := runCmd(cmd)
	return localPath, err
}

// runCmd runs the cmd & print output (both stdout & stderr)
func runCmd(cmd string) (err error) {
	fmt.Printf("CMD: %s \n", cmd)
	out, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	log.Info(string(out))
	if err != nil {
		err = errors.Wrap(err, string(out))
	}
	return
}
