package realname

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// Upload upload ID card.
func (s *Service) Upload(c context.Context, mid int64, data []byte) (src string, err error) {
	var (
		md5Engine = md5.New()
		hashMID   string
		hashRand  string
		fileName  string
		dirPath   string
		dateStr   string
	)
	md5Engine.Write([]byte(strconv.FormatInt(mid, 10)))
	hashMID = hex.EncodeToString(md5Engine.Sum(nil))
	md5Engine.Reset()
	md5Engine.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	md5Engine.Write([]byte(strconv.FormatInt(rand.Int63n(1000000), 10)))
	hashRand = hex.EncodeToString(md5Engine.Sum(nil))
	fileName = fmt.Sprintf("%s_%s.txt", hashMID[:6], hashRand)
	dateStr = time.Now().Format("20060102")
	dirPath = fmt.Sprintf("%s/%s/", s.c.Realname.DataDir, dateStr)

	var (
		dataFile      *os.File
		writeFileSize int
		encrptedData  []byte
	)

	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		mask := syscall.Umask(0)
		defer syscall.Umask(mask)
		if err = os.MkdirAll(dirPath, 0777); err != nil {
			err = errors.WithStack(err)
			return
		}
	}
	if encrptedData, err = s.mainCryptor.IMGEncrypt(data); err != nil {
		err = errors.WithStack(err)
		return
	}
	if dataFile, err = os.Create(dirPath + fileName); err != nil {
		err = errors.Wrapf(err, "create file %s failed", dirPath+fileName)
		return
	}
	defer dataFile.Close()
	if writeFileSize, err = dataFile.Write(encrptedData); err != nil {
		err = errors.Wrapf(err, "write file %s size %d failed", dirPath+fileName, len(encrptedData))
		return
	}
	if writeFileSize != len(encrptedData) {
		err = errors.Errorf("Write file data to %s , expected %d actual %d", dirPath+fileName, len(encrptedData), writeFileSize)
		return
	}
	src = fmt.Sprintf("%s/%s", dateStr, strings.TrimSuffix(fileName, ".txt"))
	return
}

// Preview preview id card
func (s *Service) Preview(c context.Context, mid int64, src string) (img []byte, err error) {
	var (
		filePath string
		file     *os.File
		fileInfo os.FileInfo
	)
	if !s.validateSrc(mid, src) {
		err = errors.Errorf("Preview src %s invalid", src)
		return
	}
	filePath = fmt.Sprintf("%s/%s.txt", s.c.Realname.DataDir, src)
	fileInfo, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		err = errors.WithStack(err)
		return
	}
	if time.Since(fileInfo.ModTime()) > time.Duration(s.c.Realname.ImageExpire) {
		err = errors.Errorf("Realname upload image %s expired %+v", filePath, s.c.Realname.ImageExpire)
		return
	}
	if file, err = os.Open(filePath); err != nil {
		err = errors.WithStack(err)
		return
	}
	defer file.Close()
	if img, err = ioutil.ReadAll(file); err != nil {
		err = errors.WithStack(err)
		return
	}
	return s.mainCryptor.IMGDecrypt(img)
}

func (s *Service) validateSrc(mid int64, src string) (ok bool) {
	var (
		paths     []string
		fileNames []string
	)
	if paths = strings.Split(src, "/"); len(paths) != 2 {
		return
	}
	if fileNames = strings.Split(paths[1], "_"); len(fileNames) != 2 {
		return
	}
	var (
		hash    = md5.New()
		midHash string
	)
	hash.Write([]byte(strconv.FormatInt(mid, 10)))
	midHash = hex.EncodeToString(hash.Sum(nil))[:6]
	if midHash != fileNames[0] {
		return
	}
	ok = true
	return
}
