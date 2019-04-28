package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"go-common/app/job/main/passport-game-data/model"
	"go-common/library/log"
)

const (
	_cloudJobGoroutineNum = 32
)

// ParseDiffLog parse diff log printed by compare proc.
func ParseDiffLog(src, dst string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		log.Error("failed to open file %s, error(%v)", src, err)
		return
	}
	defer f.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		log.Error("failed to open file %s, error(%v)", dst, err)
		return
	}
	defer dstFile.Close()

	var (
		line         string
		skippedCount = 0
		res          = make([]*model.CompareRes, 0)
		rd           = bufio.NewReader(f)
	)
	for {
		line, err = rd.ReadString('\n')

		if err != nil || io.EOF == err {
			break
		}

		idx := strings.LastIndex(line, "]")
		if idx == -1 {
			log.Error("failed to parse log, expected have ] in string but not")
			skippedCount++
			continue
		}

		logJSON := line[idx+1:]
		l := new(model.Log)
		if err = json.Unmarshal([]byte(logJSON), &l); err != nil {
			log.Error("failed to parse log, json.Unmarshal(%s) error(%v), skip", logJSON, err)
			skippedCount++
			continue
		}

		var cRes *model.CompareRes
		if cRes, err = diffLog2CompareRes(l.Log); err != nil {
			log.Error("diffLog2CompareRes(%s) error(%v), skip", l.Log, err)
			skippedCount++
			continue
		}

		// compare local encrypted and cloud, parse diff flags
		flags := diff(cRes.Cloud, cRes.LocalEncrypted)

		if flags == _diffTypeNon {
			continue
		}

		cRes.Flags = flags
		cRes.FlagsDesc = formatFlags(flags)
		cRes.Seq = cRes.Local.Mid % _cloudJobGoroutineNum
		res = append(res, cRes)
	}

	percentMap := make(map[uint8]*model.CountAndPercent)
	seqMap := make(map[int64]*model.SeqCountAndPercent)

	for _, v := range res {
		percent, ok := percentMap[v.Flags]
		if !ok {
			percent = &model.CountAndPercent{
				DiffType: v.FlagsDesc,
			}
			percentMap[v.Flags] = percent
		}
		percent.Count++

		seq, ok := seqMap[v.Seq]
		if !ok {
			seq = &model.SeqCountAndPercent{
				Seq: v.Seq,
			}
			seqMap[v.Seq] = seq
		}
		seq.Count++
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Cloud.Mtime.After(res[j].Cloud.Mtime)
	})

	percentList := make([]*model.CountAndPercent, 0)
	for _, v := range percentMap {
		v.Percent = fmt.Sprintf("%0.2f", 100*float64(v.Count)/float64(len(res))) + "%"
		percentList = append(percentList, v)
	}
	sort.Slice(percentList, func(i, j int) bool {
		return percentList[i].Count > percentList[j].Count
	})

	seqList := make([]*model.SeqCountAndPercent, 0)
	for _, v := range seqMap {
		v.Percent = fmt.Sprintf("%0.2f", 100*float64(v.Count)/float64(_cloudJobGoroutineNum)) + "%"
		seqList = append(seqList, v)
	}
	sort.Slice(seqList, func(i, j int) bool {
		return seqList[i].Count > seqList[j].Count
	})

	stat := &model.DiffParseResp{
		Total:            len(res),
		SeqAndPercents:   seqList,
		CompareResList:   res,
		CountAndPercents: percentList,
	}

	str, _ := json.Marshal(stat)
	_, err = dstFile.WriteString(string(str))
	if err != nil {
		log.Info("failed to write parse diff log result to file %s, error(%v)", dst, err)
	}

	log.Info("len res: %d, write ok", len(res))
	return
}

func diffLog2CompareRes(str string) (*model.CompareRes, error) {
	idx := strings.Index(str, "local")
	if idx == -1 {
		return nil, fmt.Errorf("failed to parse diff log, expected have local in string but not")
	}
	res := replace(str[idx:])
	cRes := new(model.CompareRes)
	err := json.Unmarshal([]byte(res), &cRes)
	return cRes, err
}

// parse string like "local({\"mid\":1}) local_encrypted({\"mid\":1}) cloud({\"mid\":1})" to json string {"local":{},"local_encrypted":{},"cloud":{}}
func replace(str string) string {
	res := strings.Replace(str, "local(", `{"local":`, -1)

	res = strings.Replace(res, "local_encrypted(", `"local_encrypted":`, -1)

	res = strings.Replace(res, "cloud(", `"cloud":`, -1)

	res = strings.Replace(res, ")", ",", -1)

	res = strings.Replace(res, "\\", "", -1)

	if strings.HasSuffix(res, ",") {
		res = res[:len(res)-1]
	}

	res = res + "}"
	return res
}

const (
	_diffTypeNon            = uint8(0)  // 0x00000000
	_diffTypePwd            = uint8(1)  // 0x00000001
	_diffTypeEmail          = uint8(2)  // 0x00000010
	_diffTypeTel            = uint8(4)  // 0x00000100
	_diffTypeCountryID      = uint8(16) // 0x00001000
	_diffTypeMobileVerified = uint8(32) // 0x00010000
	_diffTypeIsLeak         = uint8(64) // 0x00100000
)

func formatFlags(flags uint8) string {
	fs := make([]string, 0)

	if flags&_diffTypePwd > 0 {
		fs = append(fs, "pwd")
	}

	if flags&_diffTypeEmail > 0 {
		fs = append(fs, "email")
	}

	if flags&_diffTypeTel > 0 {
		fs = append(fs, "tel")
	}

	if flags&_diffTypeCountryID > 0 {
		fs = append(fs, "country_id")
	}

	if flags&_diffTypeMobileVerified > 0 {
		fs = append(fs, "mobile_verified")
	}

	if flags&_diffTypeIsLeak > 0 {
		fs = append(fs, "is_leak")
	}

	if len(fs) == 0 {
		return "non"
	}
	return strings.Join(fs, ",")
}

func diff(cloud, localEncrypted *model.AsoAccount) uint8 {
	if localEncrypted == cloud {
		return _diffTypeNon
	}
	if localEncrypted == nil || cloud == nil {
		return _diffTypePwd | _diffTypeEmail | _diffTypeTel
	}

	res := _diffTypeNon

	if cloud.Salt != localEncrypted.Salt || cloud.Pwd != localEncrypted.Pwd {
		res = res | _diffTypePwd
	}

	if cloud.Email != localEncrypted.Email {
		res = res | _diffTypeEmail
	}

	if cloud.Tel != localEncrypted.Tel {
		res = res | _diffTypeTel
	}

	if cloud.CountryID != localEncrypted.CountryID {
		res = res | _diffTypeCountryID
	}

	if cloud.MobileVerified != localEncrypted.MobileVerified {
		res = res | _diffTypeMobileVerified
	}

	if cloud.Isleak != localEncrypted.Isleak {
		res = res | _diffTypeIsLeak
	}

	return res
}
