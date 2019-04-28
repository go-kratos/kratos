package service

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"go-common/app/admin/ep/tapd/model"
	"go-common/library/ecode"
)

// TapdCallBack Tapd Call Back.
func (s *Service) TapdCallBack(c context.Context, body io.ReadCloser) (err error) {

	var (
		URLs           []string
		jsonByte       []byte
		eventRequest   *model.EventRequest
		eventInterface = make(map[string]interface{})
		createdTime    time.Time
		workspaceID    int
		eventID        int
	)

	if jsonByte, err = ioutil.ReadAll(body); err != nil {
		return
	}

	//get event
	if err = json.Unmarshal(jsonByte, &eventRequest); err != nil {
		return
	}

	if eventRequest.Secret != s.c.Tapd.CallbackToken {
		err = ecode.Unauthorized
		return
	}

	if err = json.Unmarshal(jsonByte, &eventInterface); err != nil {
		return
	}

	// add log
	workspaceID, _ = strconv.Atoi(eventRequest.WorkspaceID)
	eventID, _ = strconv.Atoi(eventRequest.EventID)

	eventLog := &model.EventLog{
		Event:       string(eventRequest.Event),
		WorkspaceID: workspaceID,
		EventID:     eventID,
	}
	if err = s.dao.AddEventLog(eventLog); err != nil {
		return
	}

	//handle special param
	if createdTime, err = time.Parse("2006-01-02 15:04:05", eventRequest.Created); err != nil {
		return
	}
	eventInterface["id"] = eventRequest.EventID
	eventInterface["created"] = createdTime.Unix()

	if URLs, err = s.GetEnableHookURL(c, eventRequest.Event, workspaceID); err != nil {
		return
	}

	for _, URL := range URLs {
		s.transferChan.Do(context.Background(), func(c context.Context) {
			s.dao.CallHookUrlAsForm(context.Background(), URL, eventInterface)
		})
	}
	return
}

// GetEnableHookURL Get Enable Hook URL.
func (s *Service) GetEnableHookURL(c context.Context, eventType model.Event, workspaceID int) (URLs []string, err error) {
	var hookURLs []*model.HookUrl

	if s.c.Tapd.UseCache {
		if hookURLs, err = s.dao.GetEnableHookURLFromCacheAndSaveIfNot(eventType); err != nil {
			return
		}
	} else {
		if hookURLs, err = s.dao.GetEnableHookURLFromDB(eventType); err != nil {
			return
		}
	}

	for _, hookURL := range hookURLs {
		if hookURL.ID > 0 && hookURL.Status == model.HookURLStatusEnable {
			if hookURL.WorkspaceID == 0 {
				URLs = append(URLs, hookURL.URL)
			} else {
				if workspaceID == hookURL.WorkspaceID {
					URLs = append(URLs, hookURL.URL)
				}
			}
		}
	}
	return
}

// SaveEnableHookURL Save Enable Hook URL.
func (s *Service) SaveEnableHookURL(c context.Context) (eventMap map[string][]*model.HookUrl, err error) {
	return s.dao.SaveEnableHookURLToCache()
}

// QueryEnableHookURLInCache query Enable Hook URL In Cache.
func (s *Service) QueryEnableHookURLInCache(c context.Context) (hookURLs []*model.HookUrl, err error) {
	return s.dao.GetEnableHookURLFromCache(model.StoryCreate)
}

// UpdateHookURL Update Hook URL
func (s *Service) UpdateHookURL(c context.Context, username string, req *model.HookURLUpdateReq) (rep map[string]interface{}, err error) {
	var repID int64

	if req.ID > 0 {
		// update
		var (
			urlEvents     []*model.UrlEvent
			urlEventsInDB []*model.UrlEvent
			hookUrlInDB   *model.HookUrl
		)

		if hookUrlInDB, err = s.dao.QueryHookURLByID(req.ID); err != nil {
			return
		}

		if hookUrlInDB.ID == 0 {
			err = ecode.NothingFound
			return
		}

		hookUrl := &model.HookUrl{
			ID:          req.ID,
			URL:         req.URL,
			WorkspaceID: req.WorkspaceID,
			Status:      req.Status,
			UpdateBy:    username,
		}

		if urlEventsInDB, err = s.dao.QueryURLEventByUrl(req.ID); err != nil {
			return
		}

		//新增 或修改为可用
		for _, eventStr := range req.Events {
			urlEvent := &model.UrlEvent{
				Event:  eventStr,
				Status: model.HookEventStatusEnable,
				UrlID:  req.ID,
			}

			for _, urlEventInDB := range urlEventsInDB {
				if urlEventInDB.Event == eventStr {
					urlEvent.ID = urlEventInDB.ID
					break
				}
			}
			urlEvents = append(urlEvents, urlEvent)
		}

		// 未传值，视为disable
		for _, urlEventInDB := range urlEventsInDB {
			if urlEventInDB.Status == model.HookEventStatusEnable {
				var isEnable bool
				for _, eventStr := range req.Events {
					if eventStr == urlEventInDB.Event {
						isEnable = true
						break
					}
				}

				if !isEnable {
					urlEvent := &model.UrlEvent{
						Event:  urlEventInDB.Event,
						Status: model.HookEventStatusDisable,
						ID:     urlEventInDB.ID,
						UrlID:  urlEventInDB.UrlID,
					}
					urlEvents = append(urlEvents, urlEvent)
				}
			}
		}

		if err = s.dao.UpdateHookURLandEvent(hookUrl, urlEvents); err != nil {
			return
		}

		repID = req.ID

	} else {
		// add
		var urlEvents []*model.UrlEvent

		hookUrl := &model.HookUrl{
			URL:         req.URL,
			WorkspaceID: req.WorkspaceID,
			Status:      req.Status,
			UpdateBy:    username,
		}

		for _, eventStr := range req.Events {
			urlEvent := &model.UrlEvent{
				Event:  eventStr,
				Status: model.HookEventStatusEnable,
			}
			urlEvents = append(urlEvents, urlEvent)

		}
		if err = s.dao.AddHookURLandEvent(hookUrl, urlEvents); err != nil {
			return
		}
		repID = hookUrl.ID
	}

	rep = make(map[string]interface{})
	rep["url_id"] = repID
	return
}

// QueryHookURL  Query Hook URL.
func (s *Service) QueryHookURL(c context.Context, req *model.QueryHookURLReq) (rep *model.QueryHookURLRep, err error) {
	var (
		total    int64
		hookUrls []*model.HookUrl
	)
	if total, hookUrls, err = s.dao.FindHookURLs(req); err != nil {
		return
	}
	rep = &model.QueryHookURLRep{
		Pagination: model.Pagination{
			PageNum:  req.PageNum,
			PageSize: req.PageSize,
		},
		Total:    total,
		HookUrls: hookUrls,
	}
	return
}

// QueryURLEvent  Query URL Event.
func (s *Service) QueryURLEvent(c context.Context, urlID int64) (rep []*model.UrlEvent, err error) {
	var tmpUrlEvents []*model.UrlEvent

	if tmpUrlEvents, err = s.dao.QueryURLEventByUrl(urlID); err != nil {
		return
	}

	for _, tmpUrlEvent := range tmpUrlEvents {
		if tmpUrlEvent.Status == model.HookEventStatusEnable {
			rep = append(rep, tmpUrlEvent)
		}
	}

	return
}

// QueryEventLog Query Event Log.
func (s *Service) QueryEventLog(c context.Context, req *model.QueryEventLogReq) (rep *model.QueryEventLogRep, err error) {
	var (
		total     int64
		eventLogs []*model.EventLog
	)
	if total, eventLogs, err = s.dao.FindEventLogs(req); err != nil {
		return
	}
	rep = &model.QueryEventLogRep{
		Pagination: model.Pagination{
			PageNum:  req.PageNum,
			PageSize: req.PageSize,
		},
		Total:     total,
		EventLogs: eventLogs,
	}
	return
}
