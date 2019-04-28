package lic

import (
	"encoding/xml"
	"fmt"
	model "go-common/app/job/main/tv/model/pgc"
	"math/rand"
	"net/url"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-"
	_serviceID  = "dataSync"
)

// BuildLic builds the skeleton of a license
func BuildLic(sign string, ps []*model.PS, count int) *model.License {
	var (
		tid = RandStringBytesRmndr(32)
		now = time.Now()
	)
	return &model.License{
		TId:       tid,
		InputTime: now.Format("20060102"),
		Sign:      sign,
		XMLData: &model.XMLData{
			Service: &model.Service{
				ID: _serviceID,
				Head: &model.Head{
					TradeID: tid,
					Date:    now.Format("2006-01-02"),
					Count:   count,
				},
				Body: &model.Body{
					ProgramSetList: &model.PSList{
						ProgramSet: ps,
					},
				},
			},
		},
	}
}

// RandStringBytesRmndr generates an random string
func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// DelLic creates the license message with only the Season ID, for deletion
func DelLic(sign string, prefix string, sid int64) *model.License {
	var (
		ps       []*model.PS
		programS = &model.PS{
			ProgramSetID: fmt.Sprintf("%s%d", prefix, sid),
		}
	)
	ps = append(ps, programS)
	return BuildLic(sign, ps, 1)
}

// DelEpLic creates the license message with only the Ep IDs, for deletion
func DelEpLic(prefix string, sign string, delEps []int) string {
	// message skeleton
	var tid = RandStringBytesRmndr(32)
	type Service struct {
		ID   string `xml:"id,attr"`
		Head *model.Head
		Body *model.DelBody `xml:"Body"`
	}
	Msg := &Service{
		ID: _serviceID,
		Head: &model.Head{
			TradeID: tid,
			Date:    time.Now().Format("2006-01-02"),
			Count:   len(delEps),
		},
		Body: &model.DelBody{
			ProgramList: &model.ProgramList{},
		},
	}
	for _, v := range delEps {
		pm := &model.Program{
			ProgramID: fmt.Sprintf("%s%d", prefix, v),
		}
		Msg.Body.ProgramList.Program = append(Msg.Body.ProgramList.Program, pm)
	}
	// combine the xml message
	xmlRes, _ := xml.MarshalIndent(Msg, " ", " ")
	params := url.Values{}
	params.Set("tId", tid)
	params.Set("inputTime", time.Now().Format("20060102"))
	params.Set("sign", sign)
	body := params.Encode()
	body = body + "&xmlData=<?xml version=\"1.0\" encoding=\"UTF-8\"?> " + string(xmlRes)
	return body
}
