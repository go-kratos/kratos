package service

import (
	"context"
	"sync"

	"go-common/app/job/main/account-summary/model"
	member "go-common/app/service/main/member/model"
	"go-common/app/service/main/member/model/block"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func (s *Service) block(ctx context.Context, mid int64) (*model.BlockSummary, error) {
	bl, err := s.dao.MemberService.BlockInfo(ctx, &block.RPCArgInfo{MID: mid})
	if err != nil {
		return nil, err
	}

	blSummary := &model.BlockSummary{
		EmbedMid:    model.EmbedMid{Mid: bl.MID},
		BlockStatus: int64(bl.BlockStatus),
		StartTime:   model.Datetime(xtime.Time(bl.StartTime)),
		EndTime:     model.Datetime(xtime.Time(bl.EndTime)),
	}
	return blSummary, nil
}

func (s *Service) relationStat(ctx context.Context, mid int64) (*model.RelationStat, error) {
	stat, err := s.dao.RelationService.Stat(ctx, &relation.ArgMid{Mid: mid})
	if err != nil {
		return nil, err
	}

	reStat := &model.RelationStat{
		EmbedMid:  model.EmbedMid{Mid: stat.Mid},
		Follower:  stat.Follower,
		Following: stat.Following,
		Whisper:   stat.Whisper,
		Black:     stat.Black,
	}
	return reStat, nil
}

func (s *Service) passportSummary(ctx context.Context, mid int64) (*model.PassportSummary, error) {
	ps := &model.PassportSummary{
		EmbedMid: model.EmbedMid{Mid: mid},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() error {
		defer wg.Done()
		pp, err := s.dao.PassportProfile(ctx, mid)
		if err != nil {
			log.Error("Failed to fetch passport profile: %+v", err)
			return err
		}
		ps.TelStatus = pp.TelStatus()
		ps.CountryID = pp.CountryCode
		ps.JoinIP = pp.JoinIP
		ps.JoinTime = model.Datetime(pp.JoinTime)
		ps.EmailSuffix = pp.EmailSuffix()
		return nil
	}()
	wg.Add(1)
	go func() error {
		defer wg.Done()
		origin, err := s.dao.AsoAccountRegOrigin(ctx, mid)
		if err != nil {
			log.Error("Failed to fetch passport aso account reg origin: %+v", err)
			return err
		}
		ps.RegType = origin.RegType
		ps.OriginType = origin.OriginType
		return nil
	}()
	wg.Wait()

	return ps, nil
}

func (s *Service) member(ctx context.Context, mid int64) (*model.MemberBase, *model.MemberExp, *model.MemberOfficial, error) {
	mb, err := s.dao.MemberService.Member(ctx, &member.ArgMemberMid{Mid: mid})
	if err != nil {
		return nil, nil, nil, err
	}

	var base *model.MemberBase
	if mb.BaseInfo != nil {
		base = &model.MemberBase{
			EmbedMid: model.EmbedMid{Mid: mb.Mid},
			Name:     mb.Name,
			Face:     mb.Face,
			Rank:     int64(mb.Rank),
			Sex:      mb.Sex,
			Sign:     mb.Sign,
			Birthday: model.Date(mb.Birthday),
		}
	}

	var exp *model.MemberExp
	if mb.LevelInfo != nil {
		exp = &model.MemberExp{
			EmbedMid: model.EmbedMid{Mid: mb.Mid},
			Exp:      int64(mb.NowExp),
		}
	}

	var of *model.MemberOfficial
	if mb.OfficialInfo != nil {
		of = &model.MemberOfficial{
			EmbedMid:    model.EmbedMid{Mid: mb.Mid},
			Role:        int64(mb.Role),
			Title:       mb.Title,
			Description: mb.Desc,
		}
	}

	return base, exp, of, nil
}

func (s *Service) syncMember(ctx context.Context, mid int64) error {
	base, exp, of, err := s.member(ctx, mid)
	if err != nil {
		return err
	}

	syncable := make([]Syncable, 0, 3)
	if base != nil {
		syncable = append(syncable, base)
	}
	if exp != nil {
		syncable = append(syncable, exp)
	}
	if of != nil {
		syncable = append(syncable, of)
	}

	for _, data := range syncable {
		if err := s.SyncToHBase(ctx, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) syncRelationStat(ctx context.Context, mid int64) error {
	reStat, err := s.relationStat(ctx, mid)
	if err != nil {
		return err
	}
	if err := s.SyncToHBase(ctx, reStat); err != nil {
		return err
	}
	return nil
}

func (s *Service) syncBlock(ctx context.Context, mid int64) error {
	blSummary, err := s.block(ctx, mid)
	if err != nil {
		return err
	}
	if err := s.SyncToHBase(ctx, blSummary); err != nil {
		return err
	}
	return nil
}

func (s *Service) syncPassportSummary(ctx context.Context, mid int64) error {
	ps, err := s.passportSummary(ctx, mid)
	if err != nil {
		return err
	}
	if err := s.SyncToHBase(ctx, ps); err != nil {
		return err
	}
	return nil
}

func (s *Service) syncRangeproc(ctx context.Context, start, end int64, worker uint64) {
	log.Info("Sync with range: start: %d, end %d, worker: %d", start, end, worker)

	syncChan := make(chan int64, worker*128)
	defer close(syncChan)

	// initial
	wg := sync.WaitGroup{}
	wg.Add(1)
	for i := uint64(0); i < worker; i++ {
		go func() {
			defer wg.Done()
			for mid := range syncChan {
				if err := s.SyncOne(context.Background(), mid); err != nil {
					log.Error("Failed to sync user with mid: %d: %+v", mid, err)
				}
			}
		}()
	}

	for j := start; j <= end; j++ {
		syncChan <- j
	}

	wg.Wait()
}
