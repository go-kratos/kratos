package service

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go-common/app/admin/main/upload/model"
	"go-common/library/ecode"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const (
	_downloadURL            = "http://%s/bfs/%s/%s"
	_uploadURL              = "http://%s/bfs/%s/%s"
	_deleteURL              = "http://%s/bfs/%s/%s"
	_privateBucket          = "facepri" // bucket to save yellow pic
	_privateBucketAppKey    = "8923aff2e1124bb2"
	_privateBucketAppSecret = "b237e8927823cc2984aee980123cb0"
)

// Add will add a record into bfs_upload_admin table
func (s *Service) Add(c context.Context, ap *model.AddParam) (err error) {
	record := new(model.Record)
	record.Bucket = ap.Bucket
	record.FileName = ap.FileName
	record.URL = ap.URL
	record.Sex = ap.Sex
	record.Politics = ap.Politics

	err = s.orm.Create(&record).Error
	return
}

// List lists records
func (s *Service) List(c context.Context, lp *model.ListParam) (listResSlice []*model.Record, err error) {
	listResSlice = make([]*model.Record, 0)
	err = s.orm.Limit(10).Order("id desc").Where("state=?", lp.State).Where("bucket=?", lp.Bucket).Find(&listResSlice).Error
	return
}

// MultiList lists records from multi bucket
func (s *Service) MultiList(c context.Context, lp *model.MultiListParam) (result []*model.MultiListResult, err error) {
	result = make([]*model.MultiListResult, 0)
	if len(lp.Bucket) == 0 {
		var buckets []*model.Bucket
		if err = s.orm.Table("bucket").Order("id desc").Find(&buckets).Error; err != nil {
			log.Error("read bucket error(%v)", err)
			return
		}
		for _, v := range buckets {
			lp.Bucket = append(lp.Bucket, v.BucketName)
		}
	}
	for _, bucket := range lp.Bucket {
		tmpResult := &model.MultiListResult{}
		tmpResult.Bucket = bucket
		tmpRecord := make([]*model.Record, 0)
		if err = s.orm.Limit(10).Order("id desc").Where("state=?", lp.State).Where("bucket=?", bucket).Find(&tmpRecord).Error; err != nil {
			return
		}
		tmpResult.Imgs = tmpRecord
		result = append(result, tmpResult)
	}
	return
}

// Delete deletes a record and delete file in bfs
func (s *Service) Delete(c context.Context, dp *model.DeleteParam) (err error) {
	var (
		downloadBytes []byte
		contentType   string
	)
	record := new(model.Record)
	if err = s.orm.Where("id=?", dp.Rid).Find(&record).Error; err != nil {
		err = errors.Wrapf(err, "Query(%d)", dp.Rid)
		return
	}
	dp.Bucket = record.Bucket
	dp.FileName = record.FileName
	if downloadBytes, contentType, err = s.download(dp); err != nil {
		return
	}
	if err = s.upload(dp, contentType, downloadBytes); err != nil {
		return
	}
	if err = s.delete(dp); err != nil {
		return
	}
	if err = s.orm.Where("id=?", dp.Rid).Update("state", 1).Update("adminid", dp.AdminID).Error; err != nil {
		err = errors.Wrapf(err, "Update(%d,%d,%d)", dp.Rid, 1, dp.AdminID)
		return
	}
	return
}

// DeleteRaw delete file in bfs
func (s *Service) DeleteRaw(c context.Context, dp *model.DeleteRawParam) (err error) {
	d := &model.DeleteParam{
		Bucket:   dp.Bucket,
		FileName: dp.FileName,
	}
	return s.delete(d)
}

// download from bfs
func (s *Service) download(dp *model.DeleteParam) (downloadBytes []byte, contentType string, err error) {
	var (
		downloadReq    *http.Request
		resp           *http.Response
		bfsDownloadURL string
	)
	client := &http.Client{
		Timeout: time.Duration(s.c.HTTPClient.Read.Timeout),
	}
	bfsDownloadURL = fmt.Sprintf(_downloadURL, s.c.BfsDownloadHost, dp.Bucket, dp.FileName)
	if downloadReq, err = http.NewRequest(http.MethodGet, bfsDownloadURL, nil); err != nil {
		log.Error("client.NewRequest(%s) error(%v)", bfsDownloadURL, err)
		return
	}
	if resp, err = client.Do(downloadReq); err != nil {
		log.Error("client.Do(%v) error(%v)", downloadReq, err)
		return
	}
	contentType = resp.Header.Get("Content-Type")
	if downloadBytes, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Error("ioutil.ReadAll(%v) error(%v)", resp.Body, err)
		return
	}
	return
}

// upload file to facepri bucket
func (s *Service) upload(dp *model.DeleteParam, contentType string, body []byte) (err error) {
	var (
		uploadReq    *http.Request
		resp         *http.Response
		bfsUploadURL string
	)
	client := &http.Client{
		Timeout: time.Duration(s.c.HTTPClient.Read.Timeout),
	}
	bfsUploadURL = fmt.Sprintf(_uploadURL, s.c.BfsUpdateHost, dp.Bucket, dp.FileName)
	if uploadReq, err = http.NewRequest(http.MethodPut, bfsUploadURL, bytes.NewReader(body)); err != nil {
		return
	}
	auth := s.authorize(_privateBucketAppKey, _privateBucketAppSecret, http.MethodPut, _privateBucket, dp.FileName, time.Now().Unix())
	uploadReq.Header.Add("Host", "bfs.bilibili.co")
	uploadReq.Header.Add("Date", time.Now().String())
	uploadReq.Header.Add("Authorization", auth)
	uploadReq.Header.Add("Content-Type", contentType)
	uploadReq.Header.Add("Date", fmt.Sprint(time.Now().Unix()))
	if resp, err = client.Do(uploadReq); err != nil {
		log.Error("client.Do(%v) error(%v)", uploadReq, err)
		return
	}
	// judge response code
	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusBadRequest:
		err = ecode.RequestErr
		return
	case http.StatusUnauthorized:
		// 验证不通过
		err = ecode.BfsUploadAuthErr
		return
	case http.StatusRequestEntityTooLarge:
		err = ecode.FileTooLarge
		return
	case http.StatusNotFound:
		err = ecode.NothingFound
		return
	case http.StatusMethodNotAllowed:
		err = ecode.MethodNotAllowed
		return
	case http.StatusServiceUnavailable:
		err = ecode.BfsUploadServiceUnavailable
		return
	case http.StatusInternalServerError:
		err = ecode.ServerErr
		return
	default:
		err = ecode.BfsUploadStatusErr
		return
	}
	code, err := strconv.Atoi(resp.Header.Get("code"))
	if err != nil || code != 200 {
		err = ecode.BfsUploadCodeErr
		return
	}
	return
}

// delete file in old bucket
func (s *Service) delete(dp *model.DeleteParam) (err error) {
	var (
		deleteReq    *http.Request
		resp         *http.Response
		bfsDeleteURL string
	)
	client := &http.Client{
		Timeout: time.Duration(s.c.HTTPClient.Read.Timeout),
	}
	bfsDeleteURL = fmt.Sprintf(_deleteURL, s.c.BfsDeleteHost, dp.Bucket, dp.FileName)
	if deleteReq, err = http.NewRequest("DELETE", bfsDeleteURL, nil); err != nil {
		log.Error("client.NewRequest(%s) error(%v)", bfsDeleteURL, err)
		return
	}
	item, ok := s.bucketCache[dp.Bucket]
	if !ok {
		err = errors.Wrapf(ecode.NothingFound, "bucket not exist: %s", dp.Bucket)
		log.Error("bucket not exist: %s", dp.Bucket)
		return
	}
	deleteReq.Header.Add("Host", "bfs.bilibili.co")
	deleteReq.Header.Add("Date", fmt.Sprint(time.Now().Unix()))
	deleteReq.Header.Add("Authorization", s.authorize(item.KeyID, item.KeySecret, http.MethodDelete, dp.Bucket, dp.FileName, time.Now().Unix()))
	if resp, err = client.Do(deleteReq); err != nil {
		log.Error("client.Do(%v) error(%v)", deleteReq, err)
		return
	}
	if resp.StatusCode != 200 {
		log.Error("bfs delete error code: %d", resp.StatusCode)
		return
	}
	return
}

// authorize return token
func (s *Service) authorize(key, secret, method, bucket, fileName string, expire int64) (authorization string) {
	var (
		content   string
		mac       hash.Hash
		signature string
	)
	content = fmt.Sprintf("%s\n%s\n%s\n%d\n", method, bucket, fileName, expire)
	mac = hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	signature = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	authorization = fmt.Sprintf("%s:%s:%d", key, signature, expire)
	return
}

// DeleteV2 deletes a record and delete file in bfs
func (s *Service) DeleteV2(c context.Context, dp *model.DeleteV2Param) (err error) {
	switch dp.Status {
	case model.PassStatus:
		if err = s.orm.Table("upload_yellowing").Where("id=?", dp.Rid).Update("state", model.PassStatus).Update("adminid", dp.AdminID).Error; err != nil {
			err = errors.Wrapf(err, "Update(%d,%d,%d)", dp.Rid, model.PassStatus, dp.AdminID)
			return
		}
	case model.DeleteStatus:
		record := new(model.Record)
		if err = s.orm.Table("upload_yellowing").Where("id=?", dp.Rid).Find(&record).Error; err != nil {
			err = errors.Wrapf(err, "Query(%d)", dp.Rid)
			return
		}
		dp.Bucket = record.Bucket
		dp.FileName = record.FileName
		if err = s.delete(&model.DeleteParam{
			Bucket:   dp.Bucket,
			FileName: dp.FileName,
		}); err != nil {
			return
		}
		if err = s.orm.Table("upload_yellowing").Where("id=?", dp.Rid).Update("state", model.DeleteStatus).Update("adminid", dp.AdminID).Error; err != nil {
			err = errors.Wrapf(err, "Update(%d,%d,%d)", dp.Rid, model.DeleteStatus, dp.AdminID)
			return
		}
	default:
		err = errors.Wrapf(err, "illegal Status(%d,%d,%d)", dp.Rid, model.DeleteStatus, dp.AdminID)
		return
	}
	return
}
