package v1

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/gob"
	"encoding/json"
	v1pb "go-common/app/service/live/broadcast-proxy/api/v1"
	"go-common/app/service/live/broadcast-proxy/server"
	"go-common/library/ecode"
	"go-common/library/log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

const (
	kRequestIsNil = "request is nil"
)

// DanmakuService struct
type DanmakuService struct {
	proxy      *server.BroadcastProxy
	dispatcher *server.CometDispatcher
}

//NewDanmakuService init
func NewDanmakuService(p *server.BroadcastProxy, d *server.CometDispatcher) (s *DanmakuService) {
	s = &DanmakuService{
		proxy:      p,
		dispatcher: d,
	}
	return s
}

func (s *DanmakuService) writeLog(method string, begin time.Time, req interface{}, resp interface{}) {
	end := time.Now()
	log.Info("method %s, request:%v, response:%v, time cost:%s", method, req, resp, end.Sub(begin).String())
}

// RoomMessage implementation
func (s *DanmakuService) RoomMessage(ctx context.Context, req *v1pb.RoomMessageRequest) (resp *v1pb.GeneralResponse, err error) {
	resp = &v1pb.GeneralResponse{}
	if req == nil {
		return resp, ecode.Error(-1, kRequestIsNil)
	}
	httpRequest, err := http.NewRequest("POST", "/dm/1/push", strings.NewReader(req.Message))
	if err != nil {
		return resp, ecode.Error(-2, err.Error())
	}
	//httpRequest = httpRequest.WithContext(ctx)
	q := httpRequest.URL.Query()
	q.Add("cid", strconv.FormatInt(int64(req.RoomId), 10))
	q.Add("ensure", strconv.FormatInt(int64(req.Ensure), 10))
	httpRequest.URL.RawQuery = q.Encode()
	httpRecorder := httptest.NewRecorder()
	s.proxy.HandleRequest(httpRecorder, httpRequest)
	httpResponse := httpRecorder.Result()
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return resp, ecode.Errorf(-3, "remote http response code:%s", httpResponse.Status)
	}
	result, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return resp, ecode.Error(-4, err.Error())
	}
	var jsonRespBody struct {
		Code int `json:"ret"`
	}
	if err = json.Unmarshal(result, &jsonRespBody); err != nil {
		return resp, ecode.Error(-5, err.Error())
	}
	if jsonRespBody.Code != 1 {
		return resp, ecode.Error(-6, "internal server error")
	}
	return resp, nil
}

// BroadcastMessage implementation
func (s *DanmakuService) BroadcastMessage(ctx context.Context, req *v1pb.BroadcastMessageRequest) (resp *v1pb.GeneralResponse, err error) {
	resp = &v1pb.GeneralResponse{}
	if req == nil {
		return resp, ecode.Error(-1, kRequestIsNil)
	}
	httpRequest, err := http.NewRequest("POST", "/dm/1/push/all", strings.NewReader(req.Message))
	if err != nil {
		return resp, ecode.Error(-2, err.Error())
	}
	var exclude strings.Builder
	if len(req.ExcludeRoomId) > 0 {
		exclude.WriteString(strconv.FormatInt(int64(req.ExcludeRoomId[0]), 10))
	}
	if len(req.ExcludeRoomId) > 1 {
		for _, excludeRoom := range req.ExcludeRoomId[1:] {
			exclude.WriteByte(',')
			exclude.WriteString(strconv.FormatInt(int64(excludeRoom), 10))
		}
	}
	//httpRequest = httpRequest.WithContext(ctx)
	q := httpRequest.URL.Query()
	q.Add("exclude_room", exclude.String())
	httpRequest.URL.RawQuery = q.Encode()
	httpRecorder := httptest.NewRecorder()
	s.proxy.HandleRequest(httpRecorder, httpRequest)
	httpResponse := httpRecorder.Result()
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return resp, ecode.Errorf(-5, "remote http response code:%s", httpResponse.Status)
	}
	result, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return resp, ecode.Error(-3, err.Error())
	}
	var jsonRespBody struct {
		Code int `json:"ret"`
	}
	if err = json.Unmarshal(result, &jsonRespBody); err != nil {
		return resp, ecode.Error(-5, err.Error())
	}
	if jsonRespBody.Code != 1 {
		return resp, ecode.Error(-4, "internal server error")
	}
	return resp, nil
}

// MultiRoomMessage implementation
func (s *DanmakuService) MultiRoomMessage(ctx context.Context, req *v1pb.MultiRoomMessageRequest) (resp *v1pb.GeneralResponse, err error) {
	resp = &v1pb.GeneralResponse{}
	if req == nil {
		return resp, ecode.Error(-1, kRequestIsNil)
	}
	jsonRequestData, err := json.Marshal(req)
	if err != nil {
		return resp, ecode.Error(-2, err.Error())
	}
	var httpRequestBuffer bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&httpRequestBuffer, gzip.BestCompression)
	if err != nil {
		return resp, ecode.Error(-3, err.Error())
	}
	gzipWriter.Write(jsonRequestData)
	gzipWriter.Close()
	httpRequest, err := http.NewRequest("POST", "/dm/v1/push/multi_room", &httpRequestBuffer)
	if err != nil {
		return resp, ecode.Error(-4, err.Error())
	}
	httpRecorder := httptest.NewRecorder()
	s.proxy.HandleRequest(httpRecorder, httpRequest)
	httpResponse := httpRecorder.Result()
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return resp, ecode.Errorf(-7, "remote http response code:%s", httpResponse.Status)
	}
	result, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return resp, ecode.Error(-5, err.Error())
	}
	var jsonRespBody struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
	if err = json.Unmarshal(result, &jsonRespBody); err != nil {
		return resp, ecode.Error(-6, err.Error())
	}
	if jsonRespBody.Code != 0 {
		return resp, ecode.Error(ecode.Code(jsonRespBody.Code), jsonRespBody.Message)
	}
	return resp, nil
}

// BatchRoomMessage implementation
func (s *DanmakuService) BatchRoomMessage(ctx context.Context, req *v1pb.BatchRoomMessageRequest) (resp *v1pb.GeneralResponse, err error) {
	resp = &v1pb.GeneralResponse{}
	if req == nil {
		return resp, ecode.Error(-1, kRequestIsNil)
	}
	jsonRequestData, err := json.Marshal(req)
	if err != nil {
		return resp, ecode.Error(-2, err.Error())
	}
	var httpRequestBuffer bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&httpRequestBuffer, gzip.BestCompression)
	if err != nil {
		return resp, ecode.Error(-3, err.Error())
	}
	gzipWriter.Write(jsonRequestData)
	gzipWriter.Close()
	httpRequest, err := http.NewRequest("POST", "/dm/v1/push/multi_msg", &httpRequestBuffer)
	if err != nil {
		return resp, ecode.Error(-4, err.Error())
	}
	httpRecorder := httptest.NewRecorder()
	s.proxy.HandleRequest(httpRecorder, httpRequest)
	httpResponse := httpRecorder.Result()
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return resp, ecode.Errorf(-7, "remote http response code:%s", httpResponse.Status)
	}
	result, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return resp, ecode.Error(-5, err.Error())
	}
	var jsonRespBody struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
	if err = json.Unmarshal(result, &jsonRespBody); err != nil {
		return resp, ecode.Error(-6, err.Error())
	}
	if jsonRespBody.Code != 0 {
		return resp, ecode.Error(ecode.Code(jsonRespBody.Code), jsonRespBody.Message)
	}
	return resp, nil
}

// UserMessage implementation
func (s *DanmakuService) UserMessage(ctx context.Context, req *v1pb.UserMessageRequest) (resp *v1pb.GeneralResponse, err error) {
	resp = &v1pb.GeneralResponse{}
	if req == nil {
		return resp, ecode.Error(-1, kRequestIsNil)
	}
	httpRequest, err := http.NewRequest("POST", "/dm/v1/push/user_msg", strings.NewReader(req.Message))
	if err != nil {
		return resp, ecode.Error(-2, err.Error())
	}
	//httpRequest = httpRequest.WithContext(ctx)
	q := httpRequest.URL.Query()
	q.Add("uid", strconv.FormatInt(int64(req.UserId), 10))
	httpRequest.URL.RawQuery = q.Encode()
	httpRecorder := httptest.NewRecorder()
	s.proxy.HandleRequest(httpRecorder, httpRequest)
	httpResponse := httpRecorder.Result()
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return resp, ecode.Errorf(-7, "remote http response code:%s", httpResponse.Status)
	}
	result, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return resp, ecode.Error(-3, err.Error())
	}
	var jsonRespBody struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
	if err = json.Unmarshal(result, &jsonRespBody); err != nil {
		return resp, ecode.Error(-4, err.Error())
	}
	if jsonRespBody.Code != 0 {
		return resp, ecode.Error(ecode.Code(jsonRespBody.Code), jsonRespBody.Message)
	}
	return resp, nil
}

// BatchUserMessage implementation
func (s *DanmakuService) BatchUserMessage(ctx context.Context, req *v1pb.BatchUserMessageRequest) (resp *v1pb.GeneralResponse, err error) {
	resp = &v1pb.GeneralResponse{}
	if req == nil {
		return resp, ecode.Error(-1, kRequestIsNil)
	}
	jsonRequestData, err := json.Marshal(req)
	if err != nil {
		return resp, ecode.Error(-2, err.Error())
	}
	var httpRequestBuffer bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&httpRequestBuffer, gzip.BestCompression)
	if err != nil {
		return resp, ecode.Error(-3, err.Error())
	}
	gzipWriter.Write(jsonRequestData)
	gzipWriter.Close()
	httpRequest, err := http.NewRequest("POST", "/dm/v1/push/multi_user_msg", &httpRequestBuffer)
	if err != nil {
		return resp, ecode.Error(-4, err.Error())
	}
	httpRecorder := httptest.NewRecorder()
	s.proxy.HandleRequest(httpRecorder, httpRequest)
	httpResponse := httpRecorder.Result()
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return resp, ecode.Errorf(-7, "remote http response code:%s", httpResponse.Status)
	}
	result, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return resp, ecode.Error(-5, err.Error())
	}
	var jsonRespBody struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
	if err = json.Unmarshal(result, &jsonRespBody); err != nil {
		return resp, ecode.Error(-6, err.Error())
	}
	if jsonRespBody.Code != 0 {
		return resp, ecode.Error(ecode.Code(jsonRespBody.Code), jsonRespBody.Message)
	}
	return resp, nil
}

// Dispatch implementation
func (s *DanmakuService) Dispatch(ctx context.Context, req *v1pb.DispatchRequest) (resp *v1pb.DispatchResponse, err error) {
	resp = &v1pb.DispatchResponse{}
	if req == nil {
		return resp, ecode.Error(-1, kRequestIsNil)
	}
	defer s.writeLog("Dispatch", time.Now(), req, resp)
	ip, host := s.dispatcher.Dispatch(req.UserIp, req.UserId)
	resp.Host = host
	resp.Ip = ip
	return resp, nil
}

func (s *DanmakuService) SetAngryValue(ctx context.Context, req *v1pb.SetAngryValueRequest) (resp *v1pb.SetAngryValueResponse, err error) {
	resp = &v1pb.SetAngryValueResponse{}
	if req == nil {
		return resp, ecode.Error(-1, kRequestIsNil)
	}
	defer s.writeLog("SetAngryValue", time.Now(), req, resp)
	if len(req.AngryValue) == 0 {
		return resp, ecode.Error(-2, "empty angry value")
	}
	var rawAngryValueBuffer bytes.Buffer
	enc := gob.NewEncoder(&rawAngryValueBuffer)
	if err = enc.Encode(req.AngryValue); err != nil {
		return resp, ecode.Error(-3, err.Error())
	}
	var httpRequestBuffer bytes.Buffer
	gzipWriter, _ := gzip.NewWriterLevel(&httpRequestBuffer, gzip.BestCompression)
	gzipWriter.Write(rawAngryValueBuffer.Bytes())
	gzipWriter.Close()
	httpRequest, err := http.NewRequest("POST", "/dm/x/internal/v2/set_angry_value", &httpRequestBuffer)
	if err != nil {
		return resp, ecode.Error(-4, err.Error())
	}
	httpRecorder := httptest.NewRecorder()
	s.proxy.HandleRequest(httpRecorder, httpRequest)
	httpResponse := httpRecorder.Result()
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return resp, ecode.Errorf(-5, "remote http response code:%s", httpResponse.Status)
	}
	result, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return resp, ecode.Error(-6, err.Error())
	}
	var jsonRespBody struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
	if err = json.Unmarshal(result, &jsonRespBody); err != nil {
		return resp, ecode.Error(-7, err.Error())
	}
	if jsonRespBody.Code != 0 {
		return resp, ecode.Error(ecode.Code(jsonRespBody.Code), jsonRespBody.Message)
	}
	return resp, nil
}

func (s *DanmakuService) GetRoomOnlineCount(ctx context.Context, req *v1pb.GetRoomOnlineCountRequest) (resp *v1pb.GetRoomOnlineCountResponse, err error) {
	resp = &v1pb.GetRoomOnlineCountResponse{}
	if req == nil {
		return resp, ecode.Error(-1, kRequestIsNil)
	}
	defer s.writeLog("GetRoomOnlineCount", time.Now(), req, resp)
	if len(req.RoomId) == 0 {
		return resp, ecode.Error(-2, "empty angry value")
	}
	var rawRoomIdBuffer bytes.Buffer
	enc := gob.NewEncoder(&rawRoomIdBuffer)
	if err = enc.Encode(req.RoomId); err != nil {
		return resp, ecode.Error(-3, err.Error())
	}
	var httpRequestBuffer bytes.Buffer
	gzipWriter, _ := gzip.NewWriterLevel(&httpRequestBuffer, gzip.BestCompression)
	gzipWriter.Write(rawRoomIdBuffer.Bytes())
	gzipWriter.Close()
	httpRequest, err := http.NewRequest("POST", "/dm/x/internal/v3/get_room_online_count", &httpRequestBuffer)
	if err != nil {
		return resp, ecode.Error(-2, err.Error())
	}
	httpRecorder := httptest.NewRecorder()
	s.proxy.HandleRequest(httpRecorder, httpRequest)
	httpResponse := httpRecorder.Result()
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != http.StatusOK {
		return resp, ecode.Errorf(-3, "remote http response code:%s", httpResponse.Status)
	}
	reader, err := gzip.NewReader(httpResponse.Body)
	if err != nil {
		return resp, ecode.Error(-4, err.Error())
	}
	defer reader.Close()
	dec := gob.NewDecoder(reader)
	var resultData struct {
		Code    int
		Message string
		Data    map[uint64]uint64
	}
	if err = dec.Decode(&resultData); err != nil {
		return resp, ecode.Error(-5, err.Error())
	}
	if resultData.Code != 0 {
		return resp, ecode.Error(ecode.Code(resultData.Code), resultData.Message)
	}
	resp.RoomOnlineCount = resultData.Data
	return resp, nil
}
