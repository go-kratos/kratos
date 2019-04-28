package service

import (
	"context"
	"strconv"
	"strings"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"
	"go-common/library/log"

	"github.com/xanzy/go-gitlab"
)

// collectprojectproc cron func
func (s *Service) collectprojectproc() {
	/*defer func() {
		if x := recover(); x != nil {
			log.Error("collectprojectproc panic(%v)", errors.WithStack(fmt.Errorf("%v", x)))
			go s.collectprojectproc()
			log.Info("collectprojectproc recover")
		}
	}()*/
	var err error
	if err = s.CollectProject(context.TODO()); err != nil {
		log.Error("s.CollectProject err (%+v)", err)
	}
}

// CollectProject collect project information
func (s *Service) CollectProject(c context.Context) (err error) {
	var (
		projects []*gitlab.Project
		total    = 0
		page     = 1
	)

	log.Info("Collect Project start")
	for page <= 1000 {

		if projects, err = s.gitlab.ListProjects(page); err != nil {
			return
		}

		num := len(projects)
		if num <= 0 {
			break
		}
		total = total + num

		for _, p := range projects {
			if err = s.insertDB(p); err != nil {
				return
			}
		}

		page = page + 1
	}
	log.Info("Collect Project end, find %d projects", total)

	return
}

// insertDB
func (s *Service) insertDB(project *gitlab.Project) (err error) {
	var (
		b           bool
		parseFail   bool
		projectInfo = &model.ProjectInfo{
			ProjectID:     project.ID,
			Name:          project.Name,
			Description:   project.Description,
			WebURL:        project.WebURL,
			Repo:          project.SSHURLToRepo,
			DefaultBranch: project.DefaultBranch,
			//Owner:         project.Owner.Name,
			SpaceName:  project.Namespace.Name,
			SpaceKind:  project.Namespace.Kind,
			Saga:       false,
			Runner:     false,
			Department: "",
			Business:   "",
			Language:   "",
		}
	)

	if project.Namespace.Kind == "user" {
		return
	}

	if b, err = s.dao.HasProjectInfo(project.ID); err != nil {
		return
	}

	if len(project.Description) > 6 {
		projectInfo.Department, projectInfo.Business, projectInfo.Language, parseFail = parseDes(project.Description)
	}

	if parseFail {
		projectInfo.Department, projectInfo.Business = parseGroup(project.Namespace.Name)
	}

	if b {
		/*if err = s.update(project.ID, projectInfo); err != nil {
			return
		}*/
		if err = s.dao.UpdateProjectInfo(project.ID, projectInfo); err != nil {
			log.Warn("UpdateProjectInfo ProjectID(%d), Description: (%s)", projectInfo.ProjectID, projectInfo.Description)
			if strings.Contains(err.Error(), "Incorrect string value") {
				projectInfo.Description = strconv.QuoteToASCII(projectInfo.Description)
			}
			if err = s.dao.UpdateProjectInfo(project.ID, projectInfo); err != nil {
				return
			}
		}
	} else {

		if err = s.dao.AddProjectInfo(projectInfo); err != nil {
			log.Warn("AddProjectInfo ProjectID(%d), Description: (%s)", projectInfo.ProjectID, projectInfo.Description)
			if strings.Contains(err.Error(), "Incorrect string value") {
				projectInfo.Description = strconv.QuoteToASCII(projectInfo.Description)
			}
			if err = s.dao.AddProjectInfo(projectInfo); err != nil {
				return
			}
		}
	}

	return
}

// update database
/*func (s *Service) update(projectID int, projectSrc *model.ProjectInfo) (err error) {
	var (
		projectDes *model.ProjectInfo
	)

	if projectDes, err = s.dao.ProjectInfoByID(projectID); err != nil {
		return
	}

	if *projectSrc == *projectDes {
		return
	}

	s.dao.UpdateProjectInfo(projectID, projectSrc)

	return
}*/

// parseDes get info from project description
func parseDes(s string) (department, business, language string, parseFail bool) {
	//[主站 android java]
	ids := strings.Index(s, "[")
	idx := strings.LastIndex(s, "]")
	if ids == -1 || idx == -1 {
		parseFail = true
		return
	}
	str := s[ids+1 : idx]

	fields := strings.Fields(str)
	if len(fields) < 3 {
		parseFail = true
		return
	}
	department = fields[0]
	business = fields[1]
	language = fields[2]

	for _, de := range conf.Conf.Property.DeInfo {
		if department == de.Label {
			department = de.Value
		}
	}
	for _, bu := range conf.Conf.Property.BuInfo {
		if business == bu.Label {
			business = bu.Value
		}
	}

	return
}

func parseGroup(s string) (department, business string) {

	group := strings.Fields(conf.Conf.Property.Group.Name)
	de := strings.Fields(conf.Conf.Property.Group.Department)
	bu := strings.Fields(conf.Conf.Property.Group.Business)

	for i := 0; i < len(group); i++ {

		if s == group[i] {
			department, business = de[i], bu[i]
		}
	}
	return
}
