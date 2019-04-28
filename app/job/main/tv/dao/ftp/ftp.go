package ftp

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"go-common/library/log"

	"github.com/ftp-master"
)

const (
	_ftpRetry = 3
	errFormat = "Func:[%s] - Step:[%s] - Error:[%v]"
	_sleep    = 100 * time.Millisecond
)

// Retry . retry one function until no error
func Retry(callback func() error, retry int, sleep time.Duration) (err error) {
	for i := 0; i < retry; i++ {
		if err = callback(); err == nil {
			return
		}
		time.Sleep(sleep)
	}
	return
}

// FileMd5 calculates the local file's md5 and store it in a file
func (d *Dao) FileMd5(path string, md5Path string) (err error) {
	var (
		content []byte
	)
	if content, err = ioutil.ReadFile(path); err != nil {
		log.Error(errFormat+" FilePath: %s", "fileMd5", "ReadFile", err, path)
		return
	}
	md5hash := md5.New()
	if _, err = io.Copy(md5hash, bytes.NewReader(content)); err != nil {
		log.Error(errFormat, "fileMd5", "CopyContent", err)
		return
	}
	md5 := md5hash.Sum(nil)
	fMd5 := hex.EncodeToString(md5[:])
	file, error := os.OpenFile(md5Path, os.O_RDWR|os.O_CREATE, 0766)
	if error != nil {
		log.Error(errFormat, "fileMd5", "OpenFile", err)
		return
	}
	file.WriteString(fMd5)
	file.Close()
	return
}

// UploadFile the file to remote frp server and update the md5 file
func (d *Dao) UploadFile(localPath string, remotePath string, url string) (err error) {
	var (
		ftpInfo  = d.conf.Search.FTP
		c        *ftp.ServerConn
		content  []byte // file's content
		fileSize int64
	)
	// Dial
	if c, err = ftp.DialTimeout(ftpInfo.Host, time.Duration(ftpInfo.Timeout)); err != nil {
		log.Error(errFormat, "uploadFile", "DialTimeout", err)
		return
	}
	// use EPSV or not
	if !ftpInfo.UseEPSV {
		c.DisableEPSV = true
	}
	// Login
	if err = c.Login(ftpInfo.User, ftpInfo.Pass); err != nil {
		log.Error(errFormat, "uploadFile", "Login", err)
		return
	}
	// Change dir
	if err = c.ChangeDir(url); err != nil {
		log.Error(errFormat, "uploadFile", "ChangeDir", err)
		return
	}
	// Upload the file
	if content, err = ioutil.ReadFile(localPath); err != nil {
		log.Error(errFormat, "uploadFile", "ReadFile", err)
		return
	}
	data := bytes.NewBuffer(content)
	if err = Retry(func() (err error) {
		return c.Stor(remotePath, data)
	}, _ftpRetry, _sleep); err != nil {
		log.Error("upArchives Error %+v", err)
		return
	}
	// Calculate the file size to check it's ok
	if fileSize, err = c.FileSize(remotePath); err != nil {
		log.Error(errFormat, "uploadFile", "FileSize", err)
		return
	}
	if localSize := int64(len(content)); localSize != fileSize {
		err = fmt.Errorf("LocalSize is %d, RemoteSize is %d", localSize, fileSize)
		log.Error(errFormat, "uploadFile", "FileSize", err)
		return
	}
	log.Info("File %s is uploaded successfully, size: %d", remotePath, fileSize)
	return
}
