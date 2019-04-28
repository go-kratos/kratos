package service

import (
	"os"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
)

const pathPerm = 0775

// uniqueFolderPath Unique Folder Path
func (s *Service) uniqueFolderPath(path string) (uniquePath string, err error) {
	uniquePath = path + strconv.Itoa(time.Now().Nanosecond()) + "/"
	for {
		var isExists bool
		if isExists, err = exists(uniquePath); err != nil {
			return
		}
		if !isExists {
			if err = os.MkdirAll(uniquePath, pathPerm); err != nil {
				uniquePath = ""
				log.Error("Create err ... (%v)", err)
				return
			}
			break
		} else {
			uniquePath = path + strconv.Itoa(time.Now().Nanosecond()) + "/"
		}
	}
	return
}

// getSessionInCookie get session
func (s *Service) getSessionInCookie(cookie string) (session string) {
	cookieStr := strings.Split(cookie, ";")
	for _, value := range cookieStr {
		strt := strings.TrimSpace(value)
		strs := strings.Split(strt, "=")
		if strs[0] == "_AJSESSIONID" {
			session = strs[1]
			return
		}
	}
	return
}
