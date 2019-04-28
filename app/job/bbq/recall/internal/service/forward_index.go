package service

import (
	"context"
	"encoding/hex"
	"os"
	"os/exec"

	"go-common/app/job/bbq/recall/proto"
	"go-common/library/log"

	"github.com/golang/snappy"
)

// GenForwardIndex 生产正排索引
func (s *Service) GenForwardIndex() {
	log.Info("run [%s]", "GenForwardIndex")
	c := context.Background()
	vInfo, err := s.videoBasicInfo(c)
	if err != nil {
		log.Error("video info: %v", err)
		return
	}

	outputFile, err := os.OpenFile(s.c.Job.ForwardIndex.Output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Error("open file: %v", err)
		return
	}
	shadowFile, _ := os.OpenFile(s.c.Job.ForwardIndex.Output+".bak", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer outputFile.Close()
	defer shadowFile.Close()

	for _, v := range vInfo {
		qu, _ := s.dao.FetchVideoQuality(c, v.SVID)
		tmp := &proto.ForwardIndex{
			SVID:         v.SVID,
			BasicInfo:    v,
			VideoQuality: qu,
		}
		raw, err := tmp.Marshal()
		if err != nil {
			log.Error("json marshal: %v", err)
			continue
		}
		_, err = outputFile.WriteString(hex.EncodeToString(snappy.Encode(nil, raw)))
		if err != nil {
			log.Error("output: %v", err)
		}
		outputFile.Write([]byte("\n"))
		if err != nil {
			log.Error("output endline: %v", err)
		}
		shadowFile.WriteString(tmp.String())
		if err != nil {
			log.Error("shadow: %v", err)
		}
		shadowFile.Write([]byte("\n"))
		if err != nil {
			log.Error("shadow endline: %v", err)
		}
	}
	exec.Command(s.c.Job.ForwardIndex.Output + ".sh").Run()

	log.Info("finish [GenForwardIndex]")

	s.GenRealTimeInvertedIndex()
}

func (s *Service) videoBasicInfo(c context.Context) (result []*proto.VideoInfo, err error) {
	// fetch tag info from db
	tags, err := s.dao.FetchVideoTagAll(c)
	if err != nil {
		return
	}
	tagIDMap := make(map[int32]*proto.Tag)
	tagNameMap := make(map[string]*proto.Tag)
	for _, v := range tags {
		tagIDMap[v.TagID] = v
		tagNameMap[v.TagName] = v
	}

	// fetch video info from db
	offset := 0
	size := 1000
	basic, err := s.dao.FetchVideoInfo(c, offset, size)
	if err != nil {
		log.Error("FetchVideoInfo: %v", err)
		return
	}
	for len(basic) > 0 && err == nil {
		log.Info("FetchVideoInfo: %v", len(result))
		for _, v := range basic {
			vInfo := &proto.VideoInfo{
				SVID:     uint64(v.SVID),
				Title:    v.Title,
				Content:  v.Content,
				MID:      uint64(v.MID),
				AVID:     uint64(v.AVID),
				CID:      uint64(v.CID),
				PubTime:  v.PubTime.Time().Unix(),
				CTime:    v.CTime.Time().Unix(),
				MTime:    v.MTime.Time().Unix(),
				Duration: uint32(v.Duration),
				State:    int32(v.State),
			}
			vTags := make([]*proto.Tag, 0)
			// 一级标签
			if tag, ok := tagIDMap[v.TID]; ok {
				vTags = append(vTags, tag)
			}
			// 二级标签
			if subTag, ok := tagIDMap[v.SubTID]; ok {
				vTags = append(vTags, subTag)
			}
			// 三级标签
			if textTags, e := s.dao.FetchVideoTextTag(c, v.SVID); e == nil {
				for _, v := range textTags {
					if tmp, ok := tagNameMap[v]; ok {
						vTags = append(vTags, tmp)
					}
				}
			}

			vInfo.Tags = vTags
			result = append(result, vInfo)
		}
		offset += size
		basic, err = s.dao.FetchVideoInfo(c, offset, size)
		if err != nil {
			log.Error("FetchVideoInfo: %v", err)
		}
	}
	return
}
