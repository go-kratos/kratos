package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"sync"
	"time"

	"go-common/app/admin/main/usersuit/model"
	accmdl "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"

	"go-common/library/sync/errgroup"
)

const (
	_geneMaxLimit     = int64(1000)
	_geneSubCount     = 200
	_batch            = 50
	_fetchInfoTimeout = time.Second * 5
)

var (
	_emptyRichInvites = make([]*model.RichInvite, 0)
	_emptyInfoMap     = make(map[int64]*accmdl.Info)
)

// Generate generate invite codes in batch.
func (s *Service) Generate(c context.Context, mid, num, expireDay int64) (res []*model.RichInvite, err error) {
	if num > _geneMaxLimit {
		err = ecode.UsersuitInviteReachMaxGeneLimit
		return
	}
	expireSeconds := expireDay * 86400
	nowTs := time.Now().Unix()
	cm, err1 := concurrentGenerateCode(mid, nowTs, int(num), _geneSubCount)
	if err1 != nil {
		log.Error("concurrentGenerateCode(%d, %d, %d, %d) error(%v)", mid, nowTs, num, _geneSubCount, err)
	}
	ginvs := make([]*model.Invite, 0, num)
	buyIP := net.ParseIP(metadata.String(c, metadata.RemoteIP))
	for code := range cm {
		ginvs = append(ginvs, &model.Invite{
			Mid:     mid,
			Code:    code,
			IP:      IPv4toN(buyIP),
			IPng:    buyIP,
			Expires: nowTs + expireSeconds,
			Ctime:   xtime.Time(nowTs),
		})
	}
	invs := make([]*model.Invite, 0)
	var rc int64
	for _, inv := range ginvs {
		if rc, err = s.d.AddIgnoreInvite(c, inv); err != nil {
			err = nil
			break
		}
		if rc == 0 {
			log.Error("service.dao.AddIgnoreInvite(%s), duplicate entry for invite_code %s", inv.Code, inv.Code)
			continue
		}
		invs = append(invs, inv)
	}
	res = s.fillStatusAndInviteeInfo(c, invs)
	return
}

func concurrentGenerateCode(mid, ts int64, num, subCount int) (res map[string]int, err error) {
	batches := num / subCount
	eg, _ := errgroup.WithContext(context.TODO())
	ims := make([]map[string]int, batches)
	mu := sync.Mutex{}
	for i := 0; i < batches; i++ {
		idx := i
		eg.Go(func() error {
			im := make(map[string]int)
			for len(im) < subCount {
				im[geneInviteCode(mid, ts)] = 1
			}
			mu.Lock()
			ims[idx] = im
			mu.Unlock()
			return nil
		})
	}
	err = eg.Wait()
	m := make(map[string]int)
	for _, im := range ims {
		for code := range im {
			m[code] = 1
		}
	}
	for len(m) < num {
		m[geneInviteCode(mid, ts)] = 1
	}
	res = m
	return
}

func geneInviteCode(mid int64, ts int64) string {
	data := md5.Sum([]byte(fmt.Sprintf("%d,%d,%d", ts, mid, rand.Int63())))
	h := hex.EncodeToString(data[:])
	return h[8:24]
}

// List list one's invite codes range time start and end.
func (s *Service) List(c context.Context, mid, start, end int64) (res []*model.RichInvite, err error) {
	if start > end {
		res = _emptyRichInvites
		return
	}
	var invs []*model.Invite
	if invs, err = s.d.RangeInvites(c, mid, time.Unix(start, 0), time.Unix(end, 0)); err != nil {
		return
	}
	sort.Slice(invs, func(i, j int) bool {
		return int64(invs[i].Ctime) > int64(invs[j].Ctime)
	})
	res = s.fillStatusAndInviteeInfo(c, invs)
	return
}

func (s *Service) fillStatusAndInviteeInfo(c context.Context, invs []*model.Invite) []*model.RichInvite {
	if len(invs) == 0 {
		return _emptyRichInvites
	}
	imidm := make(map[int64]struct{})
	now := time.Now().Unix()
	for _, inv := range invs {
		inv.FillStatus(now)
		if inv.Status == model.StatusUsed {
			imidm[inv.Imid] = struct{}{}
		}
	}
	infom := _emptyInfoMap
	if len(imidm) > 0 {
		imids := make([]int64, 0, len(imidm))
		for imid := range imidm {
			imids = append(imids, imid)
		}
		var err1 error
		if infom, err1 = s.fetchInfos(c, imids, _fetchInfoTimeout); err1 != nil {
			log.Error("service.fetchInfos(%v, %s) error(%v)", imids, _fetchInfoTimeout, err1)
		}
	}
	rinvs := make([]*model.RichInvite, 0)
	for _, inv := range invs {
		rinvs = append(rinvs, model.NewRichInvite(inv, infom[inv.Imid]))
	}
	return rinvs
}

func (s *Service) fetchInfos(c context.Context, mids []int64, timeout time.Duration) (res map[int64]*accmdl.Info, err error) {
	if len(mids) == 0 {
		res = _emptyInfoMap
		return
	}
	batches := len(mids)/_batch + 1
	tc, cancel := context.WithTimeout(c, timeout)
	defer cancel()
	eg, errCtx := errgroup.WithContext(tc)
	bms := make([]map[int64]*accmdl.Info, batches)
	mu := sync.Mutex{}
	for i := 0; i < batches; i++ {
		idx := i
		end := (idx + 1) * _batch
		if idx == batches-1 {
			end = len(mids)
		}
		ids := mids[idx*_batch : end]
		eg.Go(func() error {
			m, err1 := s.accountClient.Infos3(errCtx, &accmdl.MidsReq{Mids: ids})
			mu.Lock()
			bms[idx] = m.Infos
			mu.Unlock()
			return err1
		})
	}
	err = eg.Wait()
	res = make(map[int64]*accmdl.Info)
	for _, bm := range bms {
		for mid, info := range bm {
			res[mid] = info
		}
	}
	return
}

// IPv4toN is
func IPv4toN(ip net.IP) (sum uint32) {
	v4 := ip.To4()
	if v4 == nil {
		return
	}
	sum += uint32(v4[0]) << 24
	sum += uint32(v4[1]) << 16
	sum += uint32(v4[2]) << 8
	sum += uint32(v4[3])
	return sum
}
