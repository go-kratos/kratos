package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"go-common/app/job/main/tag/model"
	"go-common/library/log"
	"go-common/library/xstr"

	ftp "github.com/ftp-master"
)

const (
	fileName    = "tag"
	fileMD5Name = "tag.md5"
)

func (s *Service) channelRule() (res map[int64][]string, err error) {
	var (
		lastID       int64
		n            = 50
		tids         = make([]int64, 0)
		tidMap       = make(map[int64]struct{})
		tagMap       = make(map[int64]*model.Tag)
		channelRules = make([]*model.ChannelRule, 0)
	)
	for {
		var rules []*model.ChannelRule
		if rules, err = s.dao.ChannelRules(context.TODO(), lastID); err != nil {
			return
		}
		if len(rules) == 0 {
			break
		}
		for _, rule := range rules {
			if lastID < rule.ID {
				lastID = rule.ID
			}
			if rule.InRule == "" || rule.InRule == "0" {
				continue
			}
			if rule.NotInRule != "" && rule.NotInRule != "0" {
				continue
			}
			var ruleTids []int64
			if ruleTids, err = xstr.SplitInts(rule.InRule); err != nil {
				continue
			}
			usable := true
			for index, ruleTid := range ruleTids {
				if ruleTid <= 0 {
					usable = false
					break
				}
				switch index {
				case 0:
					rule.ATid = ruleTid
				case 1:
					rule.BTid = ruleTid
				default:
					usable = false
				}
				if !usable {
					break
				}
				if _, ok := tidMap[ruleTid]; ok {
					continue
				}
				tidMap[ruleTid] = struct{}{}
				tids = append(tids, ruleTid)
			}
			if usable {
				channelRules = append(channelRules, rule)
			}
		}
	}
	for len(tids) > 0 {
		if n > len(tids) {
			n = len(tids)
		}
		var tags []*model.Tag
		if tags, err = s.dao.Tags(context.TODO(), tids[:n]); err != nil {
			return
		}
		for _, tag := range tags {
			if tag.State != model.TagStateNormal {
				continue
			}
			tagMap[tag.ID] = tag
		}
		tids = tids[n:]
	}
	res = make(map[int64][]string)
	for _, r := range channelRules {
		var rule string
		tagA, ok := tagMap[r.ATid]
		if !ok || tagA == nil {
			continue
		}
		rule = tagA.Name
		tagB, b := tagMap[r.BTid]
		if b && tagB != nil {
			rule = fmt.Sprintf("%s+%s", rule, tagB.Name)
		}
		res[r.Tid] = append(res[r.Tid], rule)
	}
	return
}

// WriteTagInfo .
func (s *Service) writeTagInfo(path string) (err error) {
	var (
		lastTid int64
	)
	tagFileName := path + fileName
	file, err := os.OpenFile(tagFileName, os.O_CREATE|syscall.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("WriteTagInfo os.OpenFile(%s) error(%v)", tagFileName, err)
		return
	}
	defer file.Close()
	channelTidMap, err := s.dao.ChannelMap(context.TODO())
	if err != nil {
		return
	}
	ruleMap, err := s.channelRule()
	if err != nil {
		return
	}
	for {
		var (
			tags []*model.PlatformTagInfo
			bs   = make([]byte, 0, 1000*800)
		)
		if tags, err = s.dao.TagInfo(context.TODO(), lastTid); err != nil {
			return
		}
		if len(tags) == 0 {
			break
		}
		for _, tag := range tags {
			if _, ok := channelTidMap[tag.ID]; ok {
				tag.Channel = model.IsChannel
				if rules, b := ruleMap[tag.ID]; b {
					tag.CommonList = strings.Join(rules, string(byte(1)))
				}
			}
			var b []byte
			if b, err = json.Marshal(tag); err != nil {
				log.Warn("json.Marshal(%v) error(%v)", tag, err)
				continue
			}
			bs = append(bs, b...)
			bs = append(bs, '\n')
			lastTid = tag.ID
		}
		if _, err = file.Write(bs); err != nil {
			log.Error("writeTagInfo file.Write error(%v)", err)
			return
		}
		time.Sleep(time.Millisecond * 300)
	}
	if err = s.writeTagInfoMD5(path); err != nil {
		return
	}
	s.uploadFile(path)
	return
}

// WriteTagInfoMD5 .
func (s *Service) writeTagInfoMD5(path string) (err error) {
	var out bytes.Buffer
	f, err := exec.LookPath("md5sum")
	if err != nil {
		log.Error("writeTagInfoMD5 exec.LookPath(md5sum) error(%v)", err)
		f = "md5sum"
	}
	cmd := exec.Command(f, path+fileName)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Error("writeTagInfoMD5 cmd.Run() error(%v)", err)
		return
	}
	outStrs := strings.Split(out.String(), " ")
	if len(outStrs) == 0 {
		log.Error("writeTagInfoMD5 len(outStrs) == 0")
		return errors.New("NONE MD5")
	}
	md5FileName := path + fileMD5Name
	md5File, err := os.OpenFile(md5FileName, os.O_CREATE|syscall.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		log.Error("WriteTagInfoMD5 os.OpenFile(md5 file) error(%v)", err)
		return
	}
	defer md5File.Close()
	_, err = md5File.WriteString(outStrs[0])
	return
}

// WriteTagInfoproc .
func (s *Service) writeTagInfoproc() {
	for {
		s.writeTagInfo(s.conf.Tag.TagInfoPath)
		time.Sleep(time.Minute)
	}
}

func (s *Service) uploadFile(path string) (err error) {
	ftp, err := ftp.Connect(s.conf.FTP.Addr)
	if err != nil {
		log.Error("connect to ftp(%s) error(%v)", s.conf.FTP.Addr, err)
		return
	}
	defer ftp.Quit()
	err = ftp.Login(s.conf.FTP.User, s.conf.FTP.Password)
	if err != nil {
		log.Error("ftp login(user:%s) error(%v)", s.conf.FTP.User, err)
		return
	}
	defer ftp.Logout()
	ftp.ChangeDir(s.conf.FTP.HomeDir)
	// delete tag file, and upload file to ftp.
	if err = ftp.Delete(fileName); err != nil {
		log.Error("ftp.Delete(%s) error(%v)", fileName, err)
	}
	tagFileName := path + fileName
	file, err := os.Open(tagFileName)
	if err != nil {
		log.Error("os.Open(%s) error(%v)", tagFileName, err)
		return
	}
	defer file.Close()
	if err = ftp.Stor(fileName, file); err != nil {
		log.Error("ftp.Stor(%s) error(%v)", fileName, err)
		return
	}
	// delete tag.md5 file, and upload file to ftp.
	if err = ftp.Delete(fileMD5Name); err != nil {
		log.Error("ftp.Delete(%s) error(%v)", fileMD5Name, err)
	}
	md5FileName := path + fileMD5Name
	md5File, err := os.Open(md5FileName)
	if err != nil {
		log.Error("os.Open(%s) error(%v)", md5FileName, err)
		return
	}
	defer md5File.Close()
	if err = ftp.Stor(fileMD5Name, md5File); err != nil {
		log.Error("ftp.Stor(%s) error(%v)", fileMD5Name, err)
	}
	return
}
