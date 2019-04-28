package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go-common/app/admin/main/apm/model/ut"
	"go-common/library/log"

	"github.com/smartystreets/goconvey/web/server/contract"
	"github.com/smartystreets/goconvey/web/server/parser"
)

// Upload upload file to bfs with no filename
func (s *Service) Upload(c context.Context, fileType string, expire int64, body []byte) (url string, err error) {
	if len(body) == 0 {
		err = fmt.Errorf("s.Upload error(empty file)")
		return
	}
	if url, err = s.dao.UploadProxy(c, fileType, expire, body); err != nil {
		log.Error("s.Upload error(%v)", err)
	}
	return
}

// ParseContent parse go test output to go convey result json.
// If result not contains "=== RUN", that means execute ut err.
func (s *Service) ParseContent(c context.Context, body []byte) (content []byte, err error) {
	var (
		text  = string(body)
		index = strings.Index(text, "=== RUN")
	)
	if index == -1 {
		err = fmt.Errorf(text)
		return
	}
	var res = new(contract.PackageResult)
	parser.ParsePackageResults(res, text[index:])
	if content, err = json.Marshal(res); err != nil {
		log.Error("service.Upload json.Marshal err (%v)", err)
	}
	return
}

// CalcCount calc count
func (s *Service) CalcCount(c context.Context, body []byte) (pkg *ut.PkgAnls, err error) {
	res := new(contract.PackageResult)
	pkg = new(ut.PkgAnls)
	if err = json.Unmarshal(body, res); err != nil {
		log.Error("service.CalcCount json.Unmarshal err (%v)", err)
		return
	}
	pkg.Coverage = res.Coverage * 100
	for _, v := range res.TestResults {
		if len(v.Stories) == 0 {
			pkg.Assertions++
			if v.Error != "" {
				pkg.Panics++
			} else if !v.Passed {
				pkg.Failures++
			} else if v.Skipped {
				pkg.Skipped++
			} else {
				pkg.Passed++
			}
		}
		for _, story := range v.Stories {
			for _, ass := range story.Assertions {
				pkg.Assertions++
				if ass.Skipped {
					pkg.Skipped++
				} else if ass.Failure != "" {
					pkg.Failures++
				} else if ass.Error != nil {
					pkg.Panics++
				} else {
					pkg.Passed++
				}
			}
		}
	}
	return
}

// CalcCountFiles calculating lines and statements in files
func (s *Service) CalcCountFiles(c context.Context, res *ut.UploadRes, body []byte) (utFiles []*ut.File, err error) {
	var (
		fblocks = make(map[string][]ut.Block)
		data    = strings.Split(string(body[:]), "\n")
		reg     = regexp.MustCompile(`^(.+):([0-9]+).([0-9]+),([0-9]+).([0-9]+) ([0-9]+) ([0-9]+)$`)
	)
	if !strings.HasPrefix(data[0], "mode:") {
		return nil, fmt.Errorf("Wrong cover.dat/cover.out file format")
	}
	for i := 1; i < len(data); i++ {
		if data[i] == "" {
			continue
		}
		m := reg.FindStringSubmatch(data[i])
		if m == nil {
			return nil, fmt.Errorf("line %s doesn't match expected format: %v", data[i], reg)
		}
		b := ut.Block{
			Start:      toInt(m[2]),
			End:        toInt(m[4]),
			Statements: toInt(m[6]),
			Count:      toInt(m[7]),
		}
		fblocks[m[1]] = append(fblocks[m[1]], b)
	}
	for name, blocks := range fblocks {
		utFile := &ut.File{
			Name:     name,
			CommitID: res.CommitID,
			PKG:      res.PKG,
		}
		for i := 0; i < len(blocks); i++ {
			utFile.Statements += int64(blocks[i].Statements)
			if blocks[i].Count > 0 {
				utFile.CoveredStatements += int64(blocks[i].Statements)
			}
		}
		utFiles = append(utFiles, utFile)
	}
	return
}

// Assist functions
func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Error("strconv.Atoi(%s) error(%v)", s, err)
	}
	return i
}
