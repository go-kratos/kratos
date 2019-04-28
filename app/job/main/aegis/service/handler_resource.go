package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go-common/app/job/main/aegis/model"
)

//ERROR
var (
	ErrMangaNoIndex = errors.New("漫画无图")
)

type baseResourceAddHandler struct {
	*Service
}

type mangaResourceAddHandler struct {
	baseResourceAddHandler
}

type baseResourceUpdateHandler struct {
	*Service
}

type baseResourceCancelHandler struct {
	*Service
}

func (h baseResourceAddHandler) CheckMessage(msg json.RawMessage) (addObj interface{}, err error) {
	return h.checkRscAddMsg(msg)
}

func (h baseResourceAddHandler) HandleMessage(c context.Context, addObj interface{}) error {
	return h.writeRscAdd(c, addObj.(*model.AddOption))
}

//漫画的，校验是否有图
func (h mangaResourceAddHandler) CheckMessage(msg json.RawMessage) (addObj interface{}, err error) {
	if addObj, err = h.baseResourceAddHandler.CheckMessage(msg); err != nil {
		return
	}
	addopt := addObj.(*model.AddOption)
	metas := make(map[string]interface{})
	if err = json.Unmarshal([]byte(addopt.MetaData), &metas); err != nil {
		return
	}
	if index, ok := metas["index"]; !ok || len(fmt.Sprint(index)) == 0 {
		return nil, ErrMangaNoIndex
	}

	return
}

func (h mangaResourceAddHandler) HandleMessage(c context.Context, addObj interface{}) error {
	return h.baseResourceAddHandler.HandleMessage(c, addObj.(*model.AddOption))
}

func (h baseResourceUpdateHandler) CheckMessage(msg json.RawMessage) (updateObj interface{}, err error) {
	return h.checkRscUpdateMsg(msg)
}

func (h baseResourceUpdateHandler) HandleMessage(c context.Context, updateObj interface{}) error {
	return h.writeRscUpdate(c, updateObj.(*model.UpdateOption))
}

func (h baseResourceCancelHandler) CheckMessage(msg json.RawMessage) (cancelObj interface{}, err error) {
	return h.checkRscCancelMsg(msg)
}

func (h baseResourceCancelHandler) HandleMessage(c context.Context, cancelObj interface{}) error {
	return h.writeRscCancel(c, cancelObj.(*model.CancelOption))
}
