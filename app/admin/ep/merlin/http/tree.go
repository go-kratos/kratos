package http

import (
	"net/http"

	"go-common/app/admin/ep/merlin/model"
	bm "go-common/library/net/http/blademaster"
)

func userTree(c *bm.Context) {
	session, err := c.Request.Cookie(_sessIDKey)
	if err != nil {
		return
	}
	c.JSON(svc.UserTreeAsOption(c, session.Value))

}

func userTreeContainer(c *bm.Context) {
	var (
		tnr     = &model.TreeNode{}
		err     error
		session *http.Cookie
	)
	if err = c.Bind(tnr); err != nil {
		return
	}
	if err = tnr.VerifyFieldValue(); err != nil {
		return
	}
	if session, err = c.Request.Cookie(_sessIDKey); err != nil {
		return
	}
	c.JSON(nil, svc.VerifyTreeContainerNode(c, session.Value, tnr))
}

func treeAuditors(c *bm.Context) {
	v := new(struct {
		FirstNode string `form:"first_node"`
	})
	if err := c.Bind(v); err != nil {
		return
	}

	session, err := c.Request.Cookie(_sessIDKey)
	if err != nil {
		return
	}
	c.JSON(svc.TreeRoleAsAuditor(c, session.Value, v.FirstNode))
}
