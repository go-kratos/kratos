package http

import (
	"net/http"

	"go-common/app/admin/ep/merlin/model"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
)

func projects(c *bm.Context) {
	var (
		username string
		err      error
	)
	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(svc.Projects(c, username))
}

func accessPullProjects(c *bm.Context) {
	var (
		username string
		err      error
	)
	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(svc.AccessPullProjects(c, username))
}

func authHub(c *bm.Context) {
	var (
		err     error
		session *http.Cookie
	)
	if session, err = c.Request.Cookie(_sessIDKey); err != nil {
		return
	}

	c.JSON(nil, svc.AuthHub(c, session.Value))
}

func accessAuthHub(c *bm.Context) {
	var (
		username string
		err      error
	)
	if username, err = getUsername(c); err != nil {
		return
	}

	c.JSON(svc.AccessAuthHub(c, username))
}

func repos(c *bm.Context) {
	v := new(struct {
		model.Pagination
		ProjectID int    `form:"project_id"`
		KeyWord   string `form:"key_word"`
	})

	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.ProjectRepositories(c, v.ProjectID, v.PageNum, v.PageSize, v.KeyWord))
}

func tags(c *bm.Context) {
	v := new(struct {
		RepoName string `form:"repository_name"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.RepositoryTags(c, v.RepoName))
}

func deleteRepoTag(c *bm.Context) {
	var (
		v = new(struct {
			RepoName string `form:"repository_name"`
			TagName  string `form:"tag_name"`
		})
		username string
		err      error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.DeleteRepositoryTag(c, username, v.RepoName, v.TagName))
}

func deleteRepo(c *bm.Context) {
	var (
		v = new(struct {
			RepoName string `form:"repository_name"`
		})
		username string
		err      error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.DeleteRepository(c, username, v.RepoName))
}

func allImage(c *bm.Context) {
	c.JSON(svc.GetAllImagesInDocker())
}

func addTag(c *bm.Context) {
	var (
		v = new(struct {
			RepoName    string `json:"repository_name"`
			TagName     string `json:"tag_name"`
			NewRepoName string `json:"new_repository"`
			NewTagName  string `json:"new_tag"`
		})
		username string
		err      error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}
	c.JSON(svc.AddRepositoryTag(c, username, v.RepoName, v.TagName, v.NewRepoName, v.NewTagName))
}

func push(c *bm.Context) {
	var (
		v = new(struct {
			RepoName string `json:"repository_name"`
			TagName  string `json:"tag_name"`
		})
		username string
		err      error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	c.JSON(svc.Push(c, username, v.RepoName, v.TagName, 0))
}

func reTag(c *bm.Context) {
	var (
		v = new(struct {
			RepoName    string `json:"repository_name"`
			TagName     string `json:"tag_name"`
			NewRepoName string `json:"new_repository"`
			NewTagName  string `json:"new_tag"`
		})
		username string
		err      error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}
	c.JSON(svc.ReTag(c, username, v.RepoName, v.TagName, v.NewRepoName, v.NewTagName, 0))
}

func pull(c *bm.Context) {
	var (
		v = new(struct {
			RepoName string `json:"repository_name"`
			TagName  string `json:"tag_name"`
		})
		username string
		err      error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	c.JSON(svc.Pull(c, username, v.RepoName, v.TagName, 0))
}

func snapshot(c *bm.Context) {
	var (
		v = new(struct {
			MachineID int64 `form:"machine_id"`
		})
		username string
		err      error
	)

	if username, err = getUsername(c); err != nil {
		return
	}

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.CreateSnapShot(c, username, v.MachineID))
}

func querySnapshot(c *bm.Context) {
	var (
		v = new(struct {
			MachineID int64 `form:"machine_id"`
		})
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.QuerySnapShot(c, v.MachineID))
}

func callbackSnapshot(c *bm.Context) {
	var (
		v = new(struct {
			MachineName  string `json:"name"`
			ImageName    string `json:"image_name"`
			ResultStatus bool   `json:"status"`
			Message      string `json:"msg"`
		})
		err error
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}
	c.JSON(nil, svc.CallBackSnapShot(c, v.MachineName, v.ImageName, v.Message, v.ResultStatus))
}

func machine2image(c *bm.Context) {
	var (
		username string
		err      error
		v        = new(struct {
			MachineID    int64  `json:"machine_id"`
			ImageName    string `json:"image_name"`
			NewImageName string `json:"new_image_name"`
		})
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}
	c.JSON(nil, svc.Machine2Image(c, username, v.ImageName, v.NewImageName, v.MachineID))
}

func queryMachine2ImageLog(c *bm.Context) {
	var (
		v   = &model.QueryMachine2ImageLogRequest{}
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.QueryMachine2ImageLog(c, v))
}

func machine2imageForceFailed(c *bm.Context) {
	var (
		v = new(struct {
			MachineID int64 `form:"machine_id"`
		})
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.Machine2ImageForceFailed(c, v.MachineID))
}

func updateImageConf(c *bm.Context) {
	var (
		username string
		err      error
		v        = &model.ImageConfiguration{}
	)

	if err = c.BindWith(v, binding.JSON); err != nil {
		return
	}

	if username, err = getUsername(c); err != nil {
		return
	}
	c.JSON(svc.UpdateImageConf(c, username, v))
}

func queryImageConf(c *bm.Context) {
	var (
		v = new(struct {
			ImageName string `form:"image_full_name"`
		})
		err error
	)

	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(svc.QueryImageConf(c, v.ImageName))
}
