package service

//var (
//may = time.Now().AddDate(0, -1, 0)
//apr = time.Now().AddDate(0, -2, 0)
//)

/*
func (s *Service) syncToken(tabIdx string, min, max int64) {
	var minID int64
	var t time.Time
	if tabIdx == "201805" {
		t = may
		minID = 167003927
	} else if tabIdx == "201804" {
		t = apr
		minID = min
	}
	for {
		tokens, err := s.dao.TokenMultiQuery(context.Background(), tabIdx, minID, s.c.SyncLines)
		if err != nil {
			log.Error("s.dao.TokenMultiQuery err, sleep 1s %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(tokens) == 0 {
			break
		}
		minID = setToken(tokens, s, minID, t)
		if max > 0 && max <= minID {
			break
		}
	}
}

func setToken(tokens []*model.OldToken, s *Service, minID int64, t time.Time) int64 {
	newTokens := make([]*model.Token, 0, 51)
	newRefreshs := make([]*model.Refresh, 0, 51)
	for _, token := range tokens {
		newToken := &model.Token{
			Mid:     token.Mid,
			AppID:   token.AppID,
			Token:   token.AccessToken,
			Expires: token.Expires,
			Type:    token.Type,
		}
		newTokens = append(newTokens, newToken)
		if len(newTokens) == 50 {
			s.dao.BatchAddToken(context.Background(), newTokens, t)
			newTokens = make([]*model.Token, 0, 51)
		}
		if minID < token.ID {
			minID = token.ID
		}
		if token.RefreshToken == "" || token.Expires < time.Now().Unix() {
			continue
		}
		newRefresh := &model.Refresh{
			Mid:     token.Mid,
			AppID:   token.AppID,
			Refresh: token.RefreshToken,
			Token:   token.AccessToken,
			Expires: token.Expires,
		}
		newRefreshs = append(newRefreshs, newRefresh)
		if len(newRefreshs) == 50 {
			s.dao.BatchAddRefresh(context.Background(), newRefreshs, t)
			newRefreshs = make([]*model.Refresh, 0, 51)
		}
	}
	if len(newTokens) != 0 {
		for _, token := range newTokens {
			tokenByte, err := hex.DecodeString(token.Token)
			if err != nil {
				log.Error("setNewToken error: tokenByte decode hex err: %v", token)
				continue
			}
			s.dao.AddToken(context.Background(), token, tokenByte, t)
		}
	}
	if len(newRefreshs) != 0 {
		for _, refresh := range newRefreshs {
			refreshByte, err := hex.DecodeString(refresh.Refresh)
			tokenByteR, err1 := hex.DecodeString(refresh.Refresh)
			if err != nil || err1 != nil {
				log.Error("setNewToken error: refreshByte decode hex err: %v", refresh)
				continue
			}
			s.dao.AddRefresh(context.Background(), refresh, refreshByte, tokenByteR, t)
		}
	}
	return minID
}

func (s *Service) setNewToken(token *model.OldToken) (err error) {
	newToken := &model.Token{
		Mid:     token.Mid,
		AppID:   token.AppID,
		Token:   token.AccessToken,
		Expires: token.Expires,
		Type:    token.Type,
	}
	var t time.Time
	t, err = time.ParseInLocation("2006-01-02T15:04:05Z07:00", token.CTime, time.Local)
	if err != nil {
		log.Error("setNewToken error: ctime parse err. %v", token)
		return
	}
	tokenByte, err := hex.DecodeString(token.AccessToken)
	if err != nil {
		log.Error("setNewToken error: tokenByte decode hex err: %v", token)
		return
	}
	s.dao.AddToken(context.Background(), newToken, tokenByte, t)
	if token.RefreshToken == "" || t.Unix() < 1525104000 {
		return
	}
	newRefresh := &model.Refresh{
		Mid:     token.Mid,
		AppID:   token.AppID,
		Refresh: token.RefreshToken,
		Token:   token.AccessToken,
		Expires: token.Expires,
	}
	refreshByte, err := hex.DecodeString(token.RefreshToken)
	if err != nil {
		log.Error("setNewToken error: refreshByte decode hex err: %v", token)
		return
	}
	tokenByteR, err := hex.DecodeString(token.AccessToken)
	if err != nil {
		log.Error("setNewToken error: tokenByte decode hex err: %v", token)
		return
	}
	s.dao.AddRefresh(context.Background(), newRefresh, refreshByte, tokenByteR, t)
	return
}

func (s *Service) syncCookie(tabIdx int64) {
	var minID int64
	for {
		now := time.Now().Unix()
		cookies, err := s.dao.CookieMultiQuery(context.Background(), tabIdx, minID, now, s.c.SyncLines)
		if err != nil {
			log.Error("s.dao.CookieMultiQuery err, sleep 1s %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(cookies) == 0 {
			log.Info("syncCookie finished tabIdx : %d", tabIdx)
			break
		}
		minID = s.setCookie(cookies, now, minID)
	}
}

func (s *Service) setCookie(cookies []*model.OldCookie, now int64, minID int64) int64 {
	today := time.Now()
	newCookies := make([]*model.Cookie, 0, 51)
	var count int64
	for _, cookie := range cookies {
		if count%1000 == 0 {
			// for down cpu
			time.Sleep(10 * time.Millisecond)
		}
		if minID < cookie.ID {
			minID = cookie.ID
		}
		if len(newCookies) == 50 {
			s.dao.AddMultiCookie(context.Background(), newCookies, today)
			newCookies = make([]*model.Cookie, 0, 51)
		}
		var t time.Time
		t, err := time.ParseInLocation("2006-01-02T15:04:05Z07:00", cookie.ModifyTime, time.Local)
		if err != nil {
			log.Error("setNewToken error: ctime parse err. %v", cookie)
			continue
		}
		if cookie.Expires > now && t.Unix() > 1527782400 {
			newCookies = append(newCookies, &model.Cookie{
				Mid:     cookie.Mid,
				Session: cookie.Session,
				CSRF:    cookie.CSRFToken,
				Type:    cookie.Type,
				Expires: cookie.Expires,
			})
		}
		count = count + 1
	}
	if len(newCookies) != 0 {
		for _, newCookie := range newCookies {
			csrfByte, err := hex.DecodeString(newCookie.CSRF)
			if err != nil {
				log.Error("hex.DecodeString csrf err: %v, %v", newCookie, err)
				continue
			}
			s.dao.AddCookie(context.Background(), newCookie, []byte(newCookie.Session), csrfByte, today)
		}
	}
	return minID
}
*/
// func (s *Service) setNewCookie(cookie *model.OldCookie) (err error) {
// 	newCookie := &model.Cookie{
// 		Mid:     cookie.Mid,
// 		Session: cookie.Session,
// 		CSRF:    cookie.CSRFToken,
// 		Type:    cookie.Type,
// 		Expires: cookie.Expires,
// 	}
// 	csrfByte, err := hex.DecodeString(newCookie.CSRF)
// 	if err != nil {
// 		log.Error("hex.DecodeString csrf err: %v, %v", cookie, err)
// 		return
// 	}
// 	if _, err = s.dao.AddCookie(context.Background(), newCookie, []byte(newCookie.Session), csrfByte, time.Now()); err != nil {
// 		log.Error("s.dao.AddCookie %v, %v", cookie, err)
// 		return
// 	}
// 	return nil
// }
/*
func (s *Service) compareToken(tabIdx string, min int64, max int64) {
	var t time.Time
	if tabIdx == "201805" {
		t = may
	} else if tabIdx == "201804" {
		t = apr
	} else {
		t = time.Now()
	}
	minID := min
	for {
		tokens, err := s.dao.BatchGetToken(context.Background(), minID, s.c.SyncLines, t)
		if err != nil {
			log.Error("compareToken error: s.dao.BatchGetToken err, sleep 1s %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(tokens) == 0 {
			break
		}
		tokenMap, err := s.dao.TokenBatchQuery(context.Background(), tabIdx, tokens)
		if err != nil {
			log.Error("compareToken error: s.dao.TokenBatchQuery err, sleep 1s %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, newToken := range tokens {
			if tokenMap[newToken.Token] == "" {
				tokenByte, err := hex.DecodeString(newToken.Token)
				if err != nil {
					log.Error("compareToken error: tokenByte decode hex err: %v", err)
					continue
				}
				s.dao.DelToken(context.Background(), tokenByte, t)
			}
		}
		minID = tokens[len(tokens)-1].ID
		if max > 0 && max <= minID {
			break
		}
	}
	log.Info("token compare finished %s, %d - %d", tabIdx, minID, max)
}

func (s *Service) compareCookie(tabIdx string, min int64) {
	var t time.Time
	if tabIdx == "201805" {
		t = may
	} else if tabIdx == "201804" {
		t = apr
	} else {
		t = time.Now()
	}
	minID := min
	for {
		cookies, err := s.dao.BatchGetCookie(context.Background(), minID, s.c.SyncLines, t)
		if err != nil {
			log.Error("compareCookie error: s.dao.BatchGetCookie err, sleep 1s %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(cookies) == 0 {
			break
		}
		tabIdxMap := make(map[int64][]string)
		for _, newCookie := range cookies {
			oldTabIdx := newCookie.Mid % 30
			tabIdxMap[oldTabIdx] = append(tabIdxMap[oldTabIdx], newCookie.Session)
		}

		sessionMap, err := s.dao.CookieBatchQuery(context.Background(), tabIdxMap)
		if err != nil {
			log.Error("compareCookie error: s.dao.CookieBatchQuery err, sleep 1s %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, newCookie := range cookies {
			if sessionMap[newCookie.Session] == "" {
				s.dao.DelCookie(context.Background(), []byte(newCookie.Session), t)
			}
		}
		minID = cookies[len(cookies)-1].ID
	}
}*/
