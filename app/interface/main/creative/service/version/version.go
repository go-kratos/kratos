package version

import (
	"context"
	"go-common/app/interface/main/creative/model/version"
	"go-common/library/ecode"
)

// Versions fn
func (s *Service) versionMap(c context.Context) (versions map[string][]*version.Version, err error) {
	if s.VersionCache == nil {
		err = ecode.NothingFound
		return
	}
	versions = make(map[string][]*version.Version)
	for _, v := range s.VersionCache {
		vs := &version.Version{
			ID:       v.ID,
			Ty:       v.Ty,
			Title:    v.Title,
			Content:  v.Content,
			Link:     v.Link,
			Ctime:    v.Ctime,
			Dateline: v.Dateline,
		}
		versions[vs.Ty] = append(versions[vs.Ty], vs)
	}
	return
}
