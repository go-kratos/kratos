package service

import (
	"testing"

	"go-common/app/job/main/passport-auth/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_HandleToken(t *testing.T) {
	once.Do(startService)
	Convey("Test del cookie by cookie", t, func() {
		token := &tokenBMsg{
			Action: "insert",
			Table:  "asp_app_perm_201701",
			New: &model.OldToken{
				Mid:          4186264,
				AppID:        1,
				AccessToken:  "test token",
				RefreshToken: "test refresh",
				AppSubID:     0,
				CreateAt:     123456789000,
			},
		}
		err := s.handleToken(token)
		So(err, ShouldBeNil)
	})
}

func Test_HandleCookie(t *testing.T) {
	once.Do(startService)
	Convey("Test del cookie by cookie", t, func() {
		cookie := &cookieBMsg{
			Action: "insert",
			Table:  "asp_app_perm_201701",
			New: &model.OldCookie{
				Mid:       4186264,
				Session:   "test session",
				CSRFToken: "test csrf",
				Expires:   1234567890000,
			},
		}
		err := s.handleCookie(cookie)
		So(err, ShouldBeNil)
	})
}
