package http

import (
	"bytes"
	"net/http"

	modsvg "go-common/app/admin/main/aegis/model/svg"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// HTMLContentType
var (
	HTMLContentType        = []string{"text/html"}
	_               Render = HTML{}
)

// HTML str.
type HTML struct {
	Content []byte
	Title   string
}

// WriteContentType fn
func (j HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, HTMLContentType, j.Title, "html")
}

// Render (JSON) writes data with json ContentType.
func (j HTML) Render(w http.ResponseWriter) (err error) {
	if _, err = w.Write(j.Content); err != nil {
		err = errors.WithStack(err)
	}
	return
}

func svg(c *bm.Context) {
	opt := new(struct {
		NetID int64 `form:"net_id" validate:"required"`
		Debug int8  `form:"debug"`
	})
	if err := c.Bind(opt); err != nil {
		return
	}

	var (
		nv  *modsvg.NetView
		err error
	)
	if opt.Debug > 0 {
		nv = modsvg.DebugSVG()
	} else {
		if nv, err = srv.NetSVG(c, opt.NetID); err != nil {
			c.JSON(nil, err)
			return
		}
	}

	bs := bytes.NewBufferString("")
	nv.Execute(c.Writer, nv.Data)

	c.Render(http.StatusOK, CSV{
		Content: bs.Bytes(),
		Title:   "流程网图",
	})
}
