package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go-common/app/admin/main/activity/model"
	tagmdl "go-common/app/interface/main/tag/model"
	artmdl "go-common/app/interface/openplatform/article/model"
	acccli "go-common/app/service/main/account/api"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
)

const (
	_approved   = 1
	_pending    = 0
	_tagArcType = 3
)

// LikesList .
func (s *Service) LikesList(c context.Context, arg *model.LikesParam) (outRes *model.LikesRes, err error) {
	var (
		likeSubject     *model.ActSubject
		list            []*model.Like
		likeList        map[int64]*model.Like
		ids, wids, mids []int64
		count           int64
		offset          int
	)
	if likeSubject, err = s.dao.ActSubject(c, arg.Sid); err != nil {
		return
	}
	db := s.DB
	db = db.Where("sid = ?", likeSubject.ID)
	if len(arg.States) > 0 {
		db = db.Where("state in (?)", arg.States)
	}
	if arg.Mid > 0 {
		db = db.Where("mid = ?", arg.Mid)
	}
	if arg.Wid > 0 {
		db = db.Where("wid = ?", arg.Wid)
	}
	if err = db.Model(model.Like{}).Count(&count).Error; err != nil {
		log.Error("db.Model(model.Like{}).Count() arg(%v) error(%v) ", arg, err)
		return
	}
	offset = (arg.Page - 1) * arg.PageSize
	if err = db.Offset(offset).Limit(arg.PageSize).Order("id asc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("db.Model(model.Like{}).Find() arg(%v) error(%v)", arg, err)
		return
	}
	likeList = make(map[int64]*model.Like, len(list))
	ids = make([]int64, 0, len(list))
	wids = make([]int64, 0, len(list))
	mids = make([]int64, 0, len(list))
	for _, val := range list {
		ids = append(ids, val.ID)
		wids = append(wids, val.Wid)
		mids = append(mids, val.Mid)
		likeList[val.ID] = val
		//like 根据最新的逻辑获取，现在无法获取 todo
		likeList[val.ID].Like = 0
	}
	if len(list) > 0 {
		if err = s.GetContent(c, likeSubject.Type, likeList, ids, wids, mids); err != nil {
			log.Error("s.GetContent(%d,%v,%v,%v,%v) error(%v)", likeSubject.Type, likeList, ids, mids, wids, err)
			return
		}
	}
	outRes = &model.LikesRes{
		Likes: likeList,
	}
	outRes.Size = arg.PageSize
	outRes.Num = arg.Page
	outRes.Total = count
	return
}

// Likes .
func (s *Service) Likes(c context.Context, Sid int64, lids []int64) (likeList map[int64]*model.Like, err error) {
	var (
		likeSubject     *model.ActSubject
		like            []*model.Like
		ids, wids, mids []int64
	)
	if likeSubject, err = s.dao.ActSubject(c, Sid); err != nil {
		return
	}
	if err = s.DB.Where("id in (?)", lids).Find(&like).Error; err != nil {
		log.Error("s.DB.Where(id in (%v)).Find() error(%v)", lids, err)
		return
	}
	likeList = make(map[int64]*model.Like, len(like))
	ids = make([]int64, 0, len(like))
	wids = make([]int64, 0, len(like))
	mids = make([]int64, 0, len(like))
	for _, val := range like {
		ids = append(ids, val.ID)
		wids = append(wids, val.Wid)
		mids = append(mids, val.Mid)
		likeList[val.ID] = val
		//like 根据最新的逻辑获取，现在无法获取 todo
		likeList[val.ID].Like = 0
	}
	if len(like) > 0 {
		if err = s.GetContent(c, likeSubject.Type, likeList, ids, wids, mids); err != nil {
			log.Error("s.GetContent( %d, %v, %v, %v, %v) error(%v)", likeSubject.Type, likeList, ids, mids, wids, err)
		}
	}
	return
}

// archiveWithTag get archives and tags.
func (s *Service) archiveWithTag(c context.Context, aids []int64, likes map[int64]*model.Like) (err error) {
	var (
		archives       *arcmdl.ArcsReply
		arcErr, tagErr error
		tags           map[int64][]*tagmdl.Tag
		ip             = metadata.String(c, metadata.RemoteIP)
	)
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if archives, arcErr = s.arcClient.Arcs(errCtx, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
			log.Error("s.arcClient.Arcs(%v, %s) error(%v)", aids, ip, err)
			return arcErr
		}
		return nil
	})
	group.Go(func() error {
		arg := &tagmdl.ArgResTags{Oids: aids, Type: _tagArcType, RealIP: ip}
		if tags, tagErr = s.tagRPC.ResTags(errCtx, arg); tagErr != nil {
			log.Error("ResTags接口错误 s.tag.ResTag(%+v) error(%v)", arg, tagErr)
			return tagErr
		}
		return nil
	})
	if err = group.Wait(); err != nil {
		return
	}
	for _, val := range likes {
		if val.Wid != 0 {
			tem := make(map[string]interface{}, 2)
			if arch, ok := archives.Arcs[val.Wid]; ok && arch.IsNormal() {
				tem["archives"] = arch
			}
			if tag, ok := tags[val.Wid]; ok {
				temps := make([]string, 0, len(tag))
				for _, val := range tag {
					temps = append(temps, val.Name)
				}
				tem["tags"] = temps
			}
			val.Object = tem
		}
	}
	return
}

// accountAndContent get likecontent and accountinfo .
func (s *Service) accountAndContent(c context.Context, ids []int64, mids []int64, likes map[int64]*model.Like) (err error) {
	var (
		contents map[int64]*model.LikeContent
		accRly   *acccli.CardsReply
	)
	if contents, err = s.dao.GetLikeContent(c, ids); err != nil {
		log.Error(" s.dao.GetLikeContent(%v) error(%v)", ids, err)
		return
	}
	if accRly, err = s.accClient.Cards3(c, &acccli.MidsReq{Mids: mids}); err != nil {
		log.Error("s.AccountsInfo(%v) error(%v)", mids, err)
		return
	}
	for _, val := range likes {
		temp := make(map[string]interface{}, 2)
		if cont, ok := contents[val.ID]; ok {
			temp["content"] = cont
		}
		if val.Mid != 0 && accRly != nil {
			if acct, ok := accRly.Cards[val.Mid]; ok {
				temp["owner"] = map[string]interface{}{
					"mid":   acct.Mid,
					"name":  acct.Name,
					"face":  acct.Face,
					"sex":   acct.Sex,
					"level": acct.Level,
				}
			}
		}
		val.Object = temp
	}
	return
}

// articles .
func (s *Service) articles(c context.Context, wids []int64, likes map[int64]*model.Like) (err error) {
	var artiRes map[int64]*artmdl.Meta
	if artiRes, err = s.artRPC.ArticleMetas(c, &artmdl.ArgAids{Aids: wids}); err != nil {
		log.Error("s.ArticleMetas(%v) error(%v)", wids, err)
		return
	}
	for _, val := range likes {
		if val.Wid != 0 {
			if v, ok := artiRes[val.Wid]; ok {
				val.Object = map[string]interface{}{
					"article": v,
				}
			}
		}
	}
	return
}

// musicsAndAct .
func (s *Service) musicsAndAct(c context.Context, wids, mids []int64, likes map[int64]*model.Like) (err error) {
	var (
		musics *model.MusicRes
		accRly *acccli.CardsReply
		ip     = metadata.String(c, metadata.RemoteIP)
	)
	if musics, err = s.dao.Musics(c, wids, ip); err != nil {
		log.Error("s.dao.Musics(%v) error(%+v)", wids, err)
		return
	}
	if accRly, err = s.accClient.Cards3(c, &acccli.MidsReq{Mids: mids}); err != nil {
		log.Error("s.AccountsInfo(%v) error(%v)", mids, err)
		return
	}
	for _, val := range likes {
		temp := make(map[string]interface{}, 2)
		if v, ok := musics.Data[val.Wid]; ok {
			temp["music"] = v
		}
		if val.Mid != 0 && accRly != nil {
			if acct, ok := accRly.Cards[val.Mid]; ok {
				temp["owner"] = map[string]interface{}{
					"mid":   acct.Mid,
					"name":  acct.Name,
					"face":  acct.Face,
					"sex":   acct.Sex,
					"level": acct.Level,
				}
			}
		}
		val.Object = temp
	}
	return
}

// GetContent get act_subjet extensions .
func (s *Service) GetContent(c context.Context, likeSubType int, likes map[int64]*model.Like, ids []int64, wids []int64, mids []int64) (err error) {
	switch likeSubType {
	case model.PICTURE, model.PICTURELIKE, model.DRAWYOO, model.DRAWYOOLIKE, model.TEXT, model.TEXTLIKE, model.QUESTION:
		err = s.accountAndContent(c, ids, mids, likes)
	case model.VIDEO, model.VIDEOLIKE, model.ONLINEVOTE, model.VIDEO2, model.PHONEVIDEO, model.SMALLVIDEO:
		err = s.archiveWithTag(c, wids, likes)
	case model.ARTICLE:
		err = s.articles(c, wids, likes)
	case model.MUSIC:
		err = s.musicsAndAct(c, wids, mids, likes)
	default:
		err = ecode.RequestErr
	}
	return
}

// ItemUp .
func (s *Service) ItemUp(c context.Context, sid, wid int64, state int) (likeList *model.Like, err error) {
	var likeSubject *model.ActSubject
	if likeSubject, err = s.dao.ActSubject(c, sid); err != nil {
		log.Error("s.GetLikesSubjectByID(%d) error(%v)", sid, err)
		return
	}
	likeList = new(model.Like)
	if err = s.DB.Where("sid =?", likeSubject.ID).Where("wid = ?", wid).Last(likeList).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("s.DB.Model(model.Like{}).Where(sid =? %d).Where(wid = ?, %d).Last(),error(%v)", likeSubject.ID, wid, err)
		return
	}
	if likeList.ID == 0 {
		err = nil
		return
	}
	likeList.State = state
	likeList.Mtime = xtime.Time(time.Now().Unix())
	if err = s.DB.Model(&model.Like{}).Where("id =?", likeList.ID).Update(map[string]interface{}{"state": state, "mtime": likeList.Mtime.Time().Format("2006-01-02 15:04:05")}).Error; err != nil {
		log.Error("s.DB.Model(&model.Like{}).Where(id =?, %d).Update() error(%v) ", likeList.ID, err)
	}
	return
}

// ItemAdd .
func (s *Service) ItemAdd(c context.Context, args *model.AddLikes) (likeList *model.Like, err error) {
	var (
		likeSubject *model.ActSubject
		ip          = metadata.String(c, metadata.RemoteIP)
	)
	if likeSubject, err = s.dao.ActSubject(c, args.Sid); err != nil {
		log.Error(" s.GetLikesSubjectByID(%d) error(%v) ", args.Sid, err)
		return
	}
	//怀疑原代码这里是个bug？？ todo
	if args.Type != likeSubject.Type {
		err = errors.New("ItemAdd params liketype error")
		return
	}
	likeList = new(model.Like)
	if err = s.DB.Where("sid =?", likeSubject.ID).Where("wid = ?", args.Wid).Last(likeList).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("s.DB.Model(model.Like{}).Where(sid =?, %d).Where(wid = ?, %d).Last() error(%v)", likeSubject.ID, args.Wid, err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		addLike := &model.Like{
			Type:  args.Type,
			Wid:   args.Wid,
			Sid:   likeSubject.ID,
			Mid:   args.Mid,
			State: args.State,
		}
		ipv6 := net.ParseIP(ip)
		if ipv6 == nil {
			ipv6 = []byte{}
		}
		addLikeContent := &model.LikeContent{
			Plat:    args.Plat,
			Device:  args.Device,
			Message: "no",
			Image:   "no",
			IPv6:    ipv6,
		}
		tx := s.DB.Begin()
		if err = tx.Create(addLike).Error; err != nil {
			log.Error("db.Model(&model.Like{}).Create(%v) error(%v)", addLike, err)
			tx.Rollback()
			return
		}
		addLikeContent.ID = addLike.ID
		if err = tx.Create(addLikeContent).Error; err != nil {
			log.Error("db.Model(&model.LikeContent{}).Create(%v) error(%v)\n", addLikeContent, err)
			tx.Rollback()
			return
		}
		tx.Commit()
		likeList = addLike
	} else {
		likeList.State = args.State
		likeList.Mtime = xtime.Time(time.Now().Unix())
		if err = s.DB.Model(&model.Like{}).Where("id =?", likeList.ID).Update(map[string]interface{}{"state": args.State, "mtime": likeList.Mtime.Time().Format("2006-01-02 15:04:05")}).Error; err != nil {
			log.Error("s.DB.Model(&model.Like{}).Where(id =?, %d).Update() error(%v) ", likeList.ID, err)
		}
	}
	return

}

// AddPicContent .
func (s *Service) AddPicContent(c context.Context, args *model.AddPic) (likeID int64, err error) {
	var (
		likeSubject *model.ActSubject
		ip          = metadata.String(c, metadata.RemoteIP)
	)
	if likeSubject, err = s.dao.ActSubject(c, args.Sid); err != nil {
		log.Error(" s.GetLikesSubjectByID(%d) error(%v) ", args.Sid, err)
		return
	}
	likes := &model.Like{
		Sid:   likeSubject.ID,
		Wid:   args.Wid,
		Type:  args.Type,
		Mid:   args.Mid,
		State: _pending,
	}
	ipv6 := net.ParseIP(ip)
	if ipv6 == nil {
		ipv6 = []byte{}
	}
	likeContent := &model.LikeContent{
		IPv6:    ipv6,
		Plat:    args.Plat,
		Device:  args.Device,
		Message: args.Message,
		Link:    args.Link,
		Image:   args.Image, //只支持传image串（需要支持传文件）todo
	}
	tx := s.DB.Begin()
	if err = tx.Create(likes).Error; err != nil {
		log.Error("db.Model(&model.Like{}).Create(%v) error(%v)", likes, err)
		tx.Rollback()
		return
	}
	likeContent.ID = likes.ID
	if err = tx.Create(likeContent).Error; err != nil {
		log.Error("db.Model(&model.LikeContent{}).Create(%v) error(%v)\n", likeContent, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	likeID = likes.ID // need to update todo
	return
}

// BatchLikes .
func (s *Service) BatchLikes(c context.Context, args *model.BatchLike) (err error) {
	var (
		likes   []*model.Like
		likeMap map[int64]*model.Like
		addWid  []int64
		ipv6    = net.ParseIP(metadata.String(c, metadata.RemoteIP))
	)
	if ipv6 == nil {
		ipv6 = []byte{}
	}
	if err = s.DB.Where(fmt.Sprintf("sid = ? and wid in (%s)", xstr.JoinInts(args.Wid)), args.Sid).Find(&likes).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("s.DB.Where(%d,%v) error(%v)", args.Sid, args.Wid, err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		addWid = args.Wid
	} else {
		likeMap = make(map[int64]*model.Like, len(likes))
		for _, v := range likes {
			likeMap[v.Wid] = v
		}
		addWid = make([]int64, 0, len(args.Wid))
		for _, v := range args.Wid {
			if _, ok := likeMap[v]; !ok {
				addWid = append(addWid, v)
			}
		}
	}
	if len(addWid) == 0 {
		return
	}
	item := &model.Like{
		Sid:   args.Sid,
		Mid:   args.Mid,
		Type:  args.Type,
		State: _pending,
	}
	if err = s.dao.BatchLike(c, item, addWid, ipv6); err != nil {
		log.Error("s.dao.BatchLike(%v) error(%+v)", addWid, err)
	}
	return
}

// VideoAdd .
func (s *Service) VideoAdd(c context.Context, args *model.AddLikes) (likeList *model.Like, err error) {
	var (
		likeSubject *model.ActSubject
		ip          = metadata.String(c, metadata.RemoteIP)
	)
	if likeSubject, err = s.dao.ActSubject(c, args.Sid); err != nil {
		log.Error("s.dao.ActSubject(%d) error(%v) ", args.Sid, err)
		return
	}
	if likeSubject.Type != model.VIDEO && likeSubject.Type != model.VIDEOLIKE && likeSubject.Type != model.VIDEO2 {
		err = errors.New("VideoAdd params liketype error")
		return
	}
	likeList = new(model.Like)
	if err = s.DB.Where("sid =?", likeSubject.ID).Where("wid = ?", args.Wid).Last(likeList).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("s.DB.Model(&model.Like{}).Where(sid =?, %d).Where(wid = ?, %d).Last() error(%v)", likeSubject.ID, args.Wid, err)
		return
	}
	if err == gorm.ErrRecordNotFound {
		addLike := &model.Like{
			Type:  args.Type,
			Wid:   args.Wid,
			Sid:   likeSubject.ID,
			Mid:   args.Mid,
			State: _approved,
		}
		ipv6 := net.ParseIP(ip)
		if ipv6 == nil {
			ipv6 = []byte{}
		}
		addLikeContent := &model.LikeContent{
			IPv6:    ipv6,
			Plat:    args.Plat,
			Device:  args.Device,
			Message: "no",
			Image:   "no",
		}
		tx := s.DB.Begin()
		if err = tx.Create(addLike).Error; err != nil {
			log.Error("db.Model(&model.Like{}).Create(%v) error(%v)", addLike, err)
			tx.Rollback()
			return
		}
		addLikeContent.ID = addLike.ID
		if err = tx.Create(addLikeContent).Error; err != nil {
			log.Error("tx.Model(&model.LikeContent{}).Create(%v) error(%v)", addLikeContent, err)
			tx.Rollback()
			return
		}
		tx.Commit()
		likeList = addLike
	} else {
		likeList.State = _approved
		likeList.Mtime = xtime.Time(time.Now().Unix())
		if err = s.DB.Model(&model.Like{}).Where("id =?", likeList.ID).Update(map[string]interface{}{"state": _approved, "mtime": likeList.Mtime.Time().Format("2006-01-02 15:04:05")}).Error; err != nil {
			log.Error("ls.DB.Model(&model.Like{}).Where(id =?, %d).Update() error(%v) ", likeList.ID, err)
		}
	}
	return
}

// UpLikesState .
func (s *Service) UpLikesState(c context.Context, IDs []int64, state int, reply string, username string) (err error) {
	if err = s.DB.Model(&model.Like{}).Where("id in (?)", xstr.JoinInts(IDs)).Update(map[string]interface{}{"state": state}).Error; err != nil {
		log.Info("s.DB.Model(&model.Like{}).Where(id = ?,%d).Update() error(%v)", IDs, err)
		return
	}
	if reply != "" {
		s.DB.Model(&model.LikeContent{}).Where("id in (?)", xstr.JoinInts(IDs)).Update(map[string]interface{}{"reply": reply})
	}
	if state != 0 {
		for _, v := range IDs {
			likeLog := &model.ActLikeLog{
				Lid:   v,
				User:  username,
				State: int64(state),
			}
			s.DB.Create(likeLog)
		}
	}
	return
}

// UpLike .
func (s *Service) UpLike(c context.Context, args *model.UpLike, username string) (res int64, err error) {
	likes := map[string]interface{}{
		"Type":     args.Type,
		"Mid":      args.Mid,
		"Wid":      args.Wid,
		"State":    args.State,
		"StickTop": args.StickTop,
	}
	if err = s.DB.Model(model.Like{}).Where("id = ?", args.Lid).Update(likes).Error; err != nil {
		log.Error("s.DB.Model(model.Like{}).Where(id = ?, %d).Update() error(%v)", args.Lid, err)
		return
	}
	likeContent := map[string]interface{}{
		"Message": args.Message,
		"Reply":   args.Reply,
		"Link":    args.Link,
		"Image":   args.Image,
	}
	if err = s.DB.Model(model.LikeContent{}).Where("id = ?", args.Lid).Update(likeContent).Error; err != nil {
		log.Error("s.DB.Model(model.LikeContent{}).Where(id = ?,%d).Update() error(%v)", args.Lid, err)
		return
	}
	if args.State != 0 {
		likeLog := &model.ActLikeLog{
			Lid:   args.Lid,
			User:  username,
			State: int64(args.State),
		}
		s.DB.Create(likeLog)
	}
	res = args.Lid
	return
}

// AddLike .
func (s *Service) AddLike(c context.Context, args *model.AddLikes) (likesRes *model.Like, err error) {
	switch args.DealType {
	case "itemUp":
		likesRes, err = s.ItemUp(c, args.Sid, args.Wid, args.State)
	case "itemAdd":
		likesRes, err = s.ItemAdd(c, args)
	case "videoAdd":
		likesRes, err = s.VideoAdd(c, args)
	default:
		err = errors.New("type error")
	}
	return
}

// UpWid .
func (s *Service) UpWid(c context.Context, args *model.UpWid) (err error) {
	var (
		subject *model.ActSubject
	)
	if subject, err = s.dao.ActSubject(c, args.Sid); err != nil {
		return
	}
	like := map[string]interface{}{
		"state": args.State,
		"mtime": xtime.Time(time.Now().Unix()),
	}
	if err := s.DB.Model(&model.Like{}).Where("sid = ? and wid = ?", subject.ID, args.Wid).Update(like).Limit(1).Error; err != nil {
		log.Error("actSrv.DB.Where(sid = %d and wid = %d).Update().Limit(1) error(%v)", subject.ID, args.Wid, err)
	}
	return
}
