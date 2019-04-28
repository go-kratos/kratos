package service

import (
	"context"

	"go-common/app/admin/ep/saga/model"
	"go-common/library/log"

	"github.com/xanzy/go-gitlab"
)

// QueryProjectMembers ...
func (s *Service) QueryProjectMembers(c context.Context, req *model.ProjectDataReq) (resp []*gitlab.ProjectMember, err error) {

	var (
		members    []*gitlab.ProjectMember
		response   *gitlab.Response
		memberItem []*gitlab.ProjectMember
	)
	if members, response, err = s.gitlab.ListProjectMembers(req.ProjectID, 1); err != nil {
		return
	}

	page := 2
	for page <= response.TotalPages {
		memberItem, _, err = s.gitlab.ListProjectMembers(req.ProjectID, page)
		members = append(members, memberItem...)
		page++
	}
	resp = members
	return
}

/*-------------------------------------- sync member ----------------------------------------*/

// SyncProjectMember ...
func (s *Service) SyncProjectMember(c context.Context, projectID int) (totalPage, totalNum int, err error) {
	var (
		resp        *gitlab.Response
		members     []*gitlab.ProjectMember
		projectInfo *model.ProjectInfo
	)

	if projectInfo, err = s.dao.ProjectInfoByID(projectID); err != nil {
		return
	}

	for page := 1; ; page++ {
		totalPage++
		if members, resp, err = s.gitlab.ListProjectMembers(projectID, page); err != nil {
			return
		}
		for _, member := range members {
			memberDB := &model.StatisticsMembers{
				ProjectID:   projectID,
				ProjectName: projectInfo.Name,
				MemberID:    member.ID,
				Username:    member.Username,
				Email:       member.Email,
				Name:        member.Name,
				State:       member.State,
				CreatedAt:   member.CreatedAt,
				AccessLevel: int(member.AccessLevel),
			}
			if err = s.SaveDatabaseMember(c, memberDB); err != nil {
				log.Error("member Save Database err: projectID(%d), MemberID(%d)", projectID, member.ID)
				err = nil
				continue
			}
			totalNum++
		}

		if resp.NextPage == 0 {
			break
		}
	}

	return
}

// SaveDatabaseMember ...
func (s *Service) SaveDatabaseMember(c context.Context, memberDB *model.StatisticsMembers) (err error) {
	var total int

	if total, err = s.dao.HasMember(c, memberDB.ProjectID, memberDB.MemberID); err != nil {
		log.Error("SaveDatabaseMember HasMember(%+v)", err)
		return
	}

	// update found row
	if total == 1 {
		if err = s.dao.UpdateMember(c, memberDB.ProjectID, memberDB.MemberID, memberDB); err != nil {
			log.Error("SaveDatabaseMember UpdateMember(%+v)", err)
		}
		return
	} else if total > 1 {
		// found repeated row, this situation will not exist under normal
		log.Warn("SaveDatabaseMember Note has more rows(%d)", total)
		return
	}

	// insert row now
	if err = s.dao.CreateMember(c, memberDB); err != nil {
		log.Error("SaveDatabaseMember CreateMember(%+v)", err)
		return

	}

	return
}
