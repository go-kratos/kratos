package permit

import (
	"net/url"

	mng "go-common/app/admin/main/manager/api"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/net/rpc/warden"

	"github.com/pkg/errors"
)

const (
	_verifyURI             = "/api/session/verify"
	_permissionURI         = "/x/admin/manager/permission"
	_sessIDKey             = "_AJSESSIONID"
	_sessUIDKey            = "uid"      // manager user_id
	_sessUnKey             = "username" // LDAP username
	_defaultDomain         = ".bilibili.co"
	_defaultCookieName     = "mng-go"
	_defaultCookieLifeTime = 2592000
	// CtxPermissions will be set into ctx.
	CtxPermissions = "permissions"
)

// permissions .
type permissions struct {
	UID   int64    `json:"uid"`
	Perms []string `json:"perms"`
}

// Permit is an auth middleware.
type Permit struct {
	verifyURI       string
	permissionURI   string
	dashboardCaller string
	dsClient        *bm.Client // dashboard client
	maClient        *bm.Client // manager-admin client

	sm *SessionManager // user Session

	mng.PermitClient // mng grpc client
}

//Verify only export Verify function because of less configure
type Verify interface {
	Verify() bm.HandlerFunc
}

// Config identify config.
type Config struct {
	DsHTTPClient    *bm.ClientConfig // dashboard client config. appkey can not reuse.
	MaHTTPClient    *bm.ClientConfig // manager-admin client config
	Session         *SessionConfig
	ManagerHost     string
	DashboardHost   string
	DashboardCaller string
}

// Config2 .
type Config2 struct {
	MngClient *warden.ClientConfig
	Session   *SessionConfig
}

// New new an auth service.
func New(c *Config) *Permit {
	a := &Permit{
		dashboardCaller: c.DashboardCaller,
		verifyURI:       c.DashboardHost + _verifyURI,
		permissionURI:   c.ManagerHost + _permissionURI,
		dsClient:        bm.NewClient(c.DsHTTPClient),
		maClient:        bm.NewClient(c.MaHTTPClient),
		sm:              newSessionManager(c.Session),
	}
	return a
}

// New2 .
func New2(c *warden.ClientConfig) *Permit {
	permitClient, err := mng.NewClient(c)
	if err != nil {
		panic(errors.WithMessage(err, "Failed to dial mng rpc server"))
	}
	return &Permit{
		PermitClient: permitClient,
		sm:           &SessionManager{},
	}
}

// NewVerify new a verify service.
func NewVerify(c *Config) Verify {
	a := &Permit{
		verifyURI:       c.DashboardHost + _verifyURI,
		dsClient:        bm.NewClient(c.DsHTTPClient),
		dashboardCaller: c.DashboardCaller,
		sm:              newSessionManager(c.Session),
	}
	return a
}

// Verify2 check whether the user has logged in.
func (p *Permit) Verify2() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		sid, username, err := p.login2(ctx)
		if err != nil {
			ctx.JSON(nil, ecode.Unauthorized)
			ctx.Abort()
			return
		}
		ctx.Set(_sessUnKey, username)
		p.sm.setHTTPCookie(ctx, _defaultCookieName, sid)
	}
}

// Verify return bm HandlerFunc which check whether the user has logged in.
func (p *Permit) Verify() bm.HandlerFunc {
	return func(ctx *bm.Context) {
		si, err := p.login(ctx)
		if err != nil {
			ctx.JSON(nil, ecode.Unauthorized)
			ctx.Abort()
			return
		}
		p.sm.SessionRelease(ctx, si)
	}
}

// Permit return bm HandlerFunc which check whether the user has logged in and has the access permission of the location.
// If `permit` is empty,it will allow any logged in request.
func (p *Permit) Permit(permit string) bm.HandlerFunc {
	return func(ctx *bm.Context) {
		var (
			si    *Session
			uid   int64
			perms []string
			err   error
		)
		si, err = p.login(ctx)
		if err != nil {
			ctx.JSON(nil, ecode.Unauthorized)
			ctx.Abort()
			return
		}
		defer p.sm.SessionRelease(ctx, si)
		uid, perms, err = p.permissions(ctx, si.Get(_sessUnKey).(string))
		if err == nil {
			si.Set(_sessUIDKey, uid)
			ctx.Set(_sessUIDKey, uid)
			if md, ok := metadata.FromContext(ctx); ok {
				md[metadata.Uid] = uid
			}
		}
		if len(perms) > 0 {
			ctx.Set(CtxPermissions, perms)
		}
		if !p.permit(permit, perms) {
			ctx.JSON(nil, ecode.AccessDenied)
			ctx.Abort()
			return
		}
	}
}

// login check whether the user has logged in.
func (p *Permit) login(ctx *bm.Context) (si *Session, err error) {
	si = p.sm.SessionStart(ctx)
	if si.Get(_sessUnKey) == nil {
		var username string
		if username, err = p.verify(ctx); err != nil {
			return
		}
		si.Set(_sessUnKey, username)
	}
	ctx.Set(_sessUnKey, si.Get(_sessUnKey))
	if md, ok := metadata.FromContext(ctx); ok {
		md[metadata.Username] = si.Get(_sessUnKey)
	}
	return
}

// Permit2 same function as permit function but reply on grpc.
func (p *Permit) Permit2(permit string) bm.HandlerFunc {
	return func(ctx *bm.Context) {
		sid, username, err := p.login2(ctx)
		if err != nil {
			ctx.JSON(nil, ecode.Unauthorized)
			ctx.Abort()
			return
		}
		p.sm.setHTTPCookie(ctx, _defaultCookieName, sid)
		ctx.Set(_sessUnKey, username)
		if md, ok := metadata.FromContext(ctx); ok {
			md[metadata.Username] = username
		}
		reply, err := p.Permissions(ctx, &mng.PermissionReq{Username: username})
		if err != nil {
			if ecode.NothingFound.Equal(err) && permit != "" {
				ctx.JSON(nil, ecode.AccessDenied)
				ctx.Abort()
			}
			return
		}
		ctx.Set(_sessUIDKey, reply.Uid)
		if md, ok := metadata.FromContext(ctx); ok {
			md[metadata.Uid] = reply.Uid
		}
		if len(reply.Perms) > 0 {
			ctx.Set(CtxPermissions, reply.Perms)
		}
		if !p.permit(permit, reply.Perms) {
			ctx.JSON(nil, ecode.AccessDenied)
			ctx.Abort()
			return
		}
	}
}

// login2 .
func (p *Permit) login2(ctx *bm.Context) (sid, uname string, err error) {
	var dsbsid, mngsid string
	dsbck, err := ctx.Request.Cookie(_sessIDKey)
	if err == nil {
		dsbsid = dsbck.Value
	}
	if dsbsid == "" {
		err = ecode.Unauthorized
		return
	}
	mngck, err := ctx.Request.Cookie(_defaultCookieName)
	if err == nil {
		mngsid = mngck.Value
	}
	reply, err := p.Login(ctx, &mng.LoginReq{Mngsid: mngsid, Dsbsid: dsbsid})
	if err != nil {
		log.Error("mng rpc Login error(%v)", err)
		return
	}
	sid = reply.Sid
	uname = reply.Username
	return
}

func (p *Permit) verify(ctx *bm.Context) (username string, err error) {
	var (
		sid string
		r   = ctx.Request
	)
	session, err := r.Cookie(_sessIDKey)
	if err == nil {
		sid = session.Value
	}
	if sid == "" {
		err = ecode.Unauthorized
		return
	}
	username, err = p.verifyDashboard(ctx, sid)
	return
}

// permit check whether user has the access permission of the location.
func (p *Permit) permit(permit string, permissions []string) bool {
	if permit == "" {
		return true
	}
	for _, p := range permissions {
		if p == permit {
			// access the permit
			return true
		}
	}
	return false
}

// verifyDashboard check whether the user is valid from Dashboard.
func (p *Permit) verifyDashboard(ctx *bm.Context, sid string) (username string, err error) {
	params := url.Values{}
	params.Set("session_id", sid)
	params.Set("encrypt", "md5")
	params.Set("caller", p.dashboardCaller)
	var res struct {
		Code     int    `json:"code"`
		UserName string `json:"username"`
	}
	if err = p.dsClient.Get(ctx, p.verifyURI, metadata.String(ctx, metadata.RemoteIP), params, &res); err != nil {
		log.Error("dashboard get verify Session url(%s) error(%v)", p.verifyURI+"?"+params.Encode(), err)
		return
	}
	if ecode.Int(res.Code) != ecode.OK {
		log.Error("dashboard get verify Session url(%s) error(%v)", p.verifyURI+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	username = res.UserName
	return
}

// permissions get user's permisssions from manager-admin.
func (p *Permit) permissions(ctx *bm.Context, username string) (uid int64, perms []string, err error) {
	params := url.Values{}
	params.Set(_sessUnKey, username)
	var res struct {
		Code int         `json:"code"`
		Data permissions `json:"data"`
	}
	if err = p.maClient.Get(ctx, p.permissionURI, metadata.String(ctx, metadata.RemoteIP), params, &res); err != nil {
		log.Error("dashboard get permissions url(%s) error(%v)", p.permissionURI+"?"+params.Encode(), err)
		return
	}
	if ecode.Int(res.Code) != ecode.OK {
		log.Error("dashboard get permissions url(%s) error(%v)", p.permissionURI+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
		return
	}
	perms = res.Data.Perms
	uid = res.Data.UID
	return
}
