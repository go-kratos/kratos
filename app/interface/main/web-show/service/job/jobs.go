package job

import (
	"context"

	jobmdl "go-common/app/interface/main/web-show/model/job"
)

// Jobs get job infos
func (s *Service) Jobs(c context.Context) (js []*jobmdl.Job) {
	js = s.cache
	return
}
