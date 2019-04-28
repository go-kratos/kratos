package service

import (
	"context"
	"sort"
	"strconv"
	"time"

	"go-common/app/interface/main/esports/model"
	arcmdl "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_gameNoSub = 6
	_gameSub   = 3
	_gameIn    = 5
	_gameLive  = 4
	_gameOver  = 1
	_caleDay   = 3
	_typeMatch = "matchs"
	_typeGame  = "games"
	_typeTeam  = "teams"
	_typeYear  = "years"
	_typeTag   = "tags"
	_downline  = 0
)

var (
	_emptContest       = make([]*model.Contest, 0)
	_emptCalendar      = make([]*model.Calendar, 0)
	_emptFilter        = make([]*model.Filter, 0)
	_emptVideoList     = make([]*arcmdl.Arc, 0)
	_emptSeason        = make([]*model.Season, 0)
	_emptContestDetail = make([]*model.ContestsData, 0)
)

// FilterMatch filter match.
func (s *Service) FilterMatch(c context.Context, p *model.ParamFilter) (rs map[string][]*model.Filter, err error) {
	var (
		tmpRs                map[string][]*model.Filter
		fm                   *model.FilterES
		fMap                 map[string]map[int64]*model.Filter
		matchs, games, teams []*model.Filter
	)
	isAll := p.Tid == 0 && p.Gid == 0 && p.Mid == 0 && p.Stime == ""
	if rs, err = s.dao.FMatCache(c); err != nil {
		err = nil
	}
	if isAll && len(rs) > 0 {
		return
	}
	matchs, games, teams = s.filterLeft(c)
	tmpRs = make(map[string][]*model.Filter, 3)
	tmpRs[_typeMatch] = matchs
	tmpRs[_typeGame] = games
	tmpRs[_typeTeam] = teams
	if fm, err = s.dao.FilterMatch(c, p); err != nil {
		log.Error("s.dao.FilterMatch error(%v)", err)
		return
	}
	fMap = s.filterMap(tmpRs)
	if tmpRs, err = s.fmtES(fm, fMap); err != nil {
		log.Error("FilterMatch s.filterES error(%v)", err)
	}
	rs = make(map[string][]*model.Filter, 3)
	rs[_typeMatch] = tmpRs[_typeMatch]
	rs[_typeGame] = tmpRs[_typeGame]
	rs[_typeTeam] = tmpRs[_typeTeam]
	if isAll {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetFMatCache(c, rs)
		})
	} else {
		if len(rs[_typeMatch]) == 0 && len(rs[_typeGame]) == 0 && len(rs[_typeTeam]) == 0 {
			if tmpRs, err = s.dao.FMatCache(c); err != nil {
				err = nil
			}
			if len(tmpRs) > 0 {
				rs = tmpRs
			}
		}
	}
	return
}

func (s *Service) filterLeft(c context.Context) (matchs, games, teams []*model.Filter) {
	var (
		matchErr, gameErr, teamErr error
	)
	group := &errgroup.Group{}
	group.Go(func() error {
		if matchs, matchErr = s.dao.Matchs(context.Background()); matchErr != nil {
			log.Error("s.dao.Matchs error %v", matchErr)
		}
		return nil
	})
	group.Go(func() error {
		if games, gameErr = s.dao.Games(context.Background()); gameErr != nil {
			log.Error("s.dao.Games error %v", gameErr)
		}
		return nil
	})
	group.Go(func() error {
		if teams, teamErr = s.dao.Teams(context.Background()); teamErr != nil {
			log.Error("s.dao.Teams error %v", teamErr)
		}
		return nil
	})
	group.Wait()
	if len(matchs) == 0 {
		matchs = _emptFilter
	}
	if len(games) == 0 {
		games = _emptFilter
	}
	if len(teams) == 0 {
		teams = _emptFilter
	}
	return
}

// Calendar contest calendar count
func (s *Service) Calendar(c context.Context, p *model.ParamFilter) (rs []*model.Calendar, err error) {
	var fc map[string]int64
	before3 := time.Now().AddDate(0, 0, -_caleDay).Format("2006-01-02")
	after3 := time.Now().AddDate(0, 0, _caleDay).Format("2006-01-02")
	todayAll := p.Mid == 0 && p.Gid == 0 && p.Tid == 0 && p.Stime == before3 && p.Etime == after3
	if todayAll {
		if rs, err = s.dao.CalendarCache(c, p); err != nil {
			err = nil
		}
		if len(rs) > 0 {
			return
		}
	}
	if fc, err = s.dao.FilterCale(c, p); err != nil {
		log.Error("s.dao.FilterCale error(%v)", err)
		return
	}
	if len(fc) == 0 {
		rs = _emptCalendar
		return
	}
	for d, c := range fc {
		rs = append(rs, &model.Calendar{Stime: d, Count: c})
	}
	if todayAll {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetCalendarCache(c, p, rs)
		})
	}
	return
}

func (s *Service) fmtContest(c context.Context, contests []*model.Contest, mid int64) {
	cids := s.contestIDs(contests)
	favContest, _ := s.isFavs(c, mid, cids)
	for _, contest := range contests {
		if contest.Etime > 0 && time.Now().Unix() > contest.Etime {
			contest.GameState = _gameOver
		} else if contest.Stime <= time.Now().Unix() && (contest.Etime >= time.Now().Unix() || contest.Etime == 0) {
			if contest.LiveRoom == 0 {
				contest.GameState = _gameIn
			} else {
				contest.GameState = _gameLive
			}
		} else if contest.LiveRoom > 0 {
			if v, ok := favContest[contest.ID]; ok && v && mid > 0 {
				contest.GameState = _gameSub
			} else {
				contest.GameState = _gameNoSub
			}
		}
	}
}

// ListContest contest list.
func (s *Service) ListContest(c context.Context, mid int64, p *model.ParamContest) (rs []*model.Contest, total int, err error) {
	var (
		teams, seasons              []*model.Filter
		cData, tmpRs                []*model.Contest
		dbContests                  map[int64]*model.Contest
		group                       *errgroup.Group
		cids                        []int64
		contErr, teamErr, seasonErr error
	)
	// get from cache.
	isFirst := p.Mid == 0 && p.Gid == 0 && p.Tid == 0 && p.Stime == "" && p.GState == "" && p.Pn == 1 && len(p.Sids) == 0 && p.Sort == 0
	if isFirst {
		if rs, total, err = s.dao.ContestCache(c, p.Ps); err != nil {
			err = nil
		} else if len(rs) > 0 {
			s.fmtContest(c, rs, mid)
			return
		}
	}
	group, errCtx := errgroup.WithContext(c)
	group.Go(func() error {
		if cData, total, contErr = s.dao.SearchContest(errCtx, p); contErr != nil {
			log.Error("s.dao.SearchContest error(%v)", contErr)
		}
		return contErr
	})
	group.Go(func() error {
		if teams, teamErr = s.dao.Teams(errCtx); teamErr != nil {
			log.Error("s.dao.Teams error(%v)", teamErr)
		}
		return nil
	})
	group.Go(func() error {
		if seasons, seasonErr = s.dao.SeasonAll(errCtx); seasonErr != nil {
			log.Error("s.dao.SeasonAll error %v", seasonErr)
		}
		return nil
	})
	err = group.Wait()
	if err != nil {
		return
	}
	if total == 0 || len(cData) == 0 {
		rs = _emptContest
		return
	}
	cids = s.contestIDs(cData)
	if len(cids) > 0 {
		if dbContests, err = s.dao.EpContests(c, cids); err != nil {
			log.Error("s.dao.EpContests error(%v)", err)
			return
		}
	}
	for _, c := range cData {
		if contest, ok := dbContests[c.ID]; ok {
			tmpRs = append(tmpRs, contest)
		}
	}
	rs = s.ContestInfo(c, cids, tmpRs, teams, seasons, mid)
	if isFirst {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetContestCache(c, p.Ps, rs, total)
		})
	}
	return
}

func (s *Service) contestIDs(cData []*model.Contest) (rs []int64) {
	for _, contest := range cData {
		rs = append(rs, contest.ID)
	}
	return
}

// ContestInfo contest add  team season.
func (s *Service) ContestInfo(c context.Context, cids []int64, cData []*model.Contest, teams, seasons []*model.Filter, mid int64) (rs []*model.Contest) {
	var (
		mapTeam, mapSeason map[int64]*model.Filter
	)
	mapTeam = make(map[int64]*model.Filter, len(teams))
	for _, team := range teams {
		mapTeam[team.ID] = team
	}
	mapSeason = make(map[int64]*model.Filter, len(seasons))
	for _, season := range seasons {
		mapSeason[season.ID] = season
	}
	favContest, _ := s.isFavs(c, mid, cids)
	for _, contest := range cData {
		if contest == nil {
			continue
		}
		if v, ok := mapTeam[contest.HomeID]; ok && v != nil {
			contest.HomeTeam = v
		} else {
			contest.HomeTeam = struct{}{}
		}
		if v, ok := mapTeam[contest.AwayID]; ok && v != nil {
			contest.AwayTeam = v
		} else {
			contest.AwayTeam = struct{}{}
		}
		if v, ok := mapTeam[contest.SuccessTeam]; ok && v != nil {
			contest.SuccessTeaminfo = v
		} else {
			contest.SuccessTeaminfo = struct{}{}
		}
		if v, ok := mapSeason[contest.Sid]; ok && v != nil {
			contest.Season = v
		} else {
			contest.Season = struct{}{}
		}
		if contest.Etime > 0 && time.Now().Unix() > contest.Etime {
			contest.GameState = _gameOver
		} else if contest.Stime <= time.Now().Unix() && (contest.Etime >= time.Now().Unix() || contest.Etime == 0) {
			if contest.LiveRoom == 0 {
				contest.GameState = _gameIn
			} else {
				contest.GameState = _gameLive
			}
		} else if contest.LiveRoom > 0 {
			if v, ok := favContest[contest.ID]; ok && v && mid > 0 {
				contest.GameState = _gameSub
			} else {
				contest.GameState = _gameNoSub
			}
		}
		rs = append(rs, contest)
	}
	return
}

// ListVideo video list.
func (s *Service) ListVideo(c context.Context, p *model.ParamVideo) (rs []*arcmdl.Arc, total int, err error) {
	var (
		vData     []*model.SearchVideo
		aids      []int64
		arcsReply *arcmdl.ArcsReply
	)
	isFirst := p.Mid == 0 && p.Gid == 0 && p.Tid == 0 && p.Year == 0 && p.Tag == 0 && p.Sort == 0 && p.Pn == 1
	if isFirst {
		// get from cache.
		if rs, total, err = s.dao.VideoCache(c, p.Ps); err != nil {
			err = nil
		} else if len(rs) > 0 {
			return
		}
	}
	if vData, total, err = s.dao.SearchVideo(c, p); err != nil {
		log.Error("s.dao.SearchVideo(%v) error(%v)", p, err)
		return
	}
	if total == 0 {
		rs = _emptVideoList
		return
	}
	for _, arc := range vData {
		aids = append(aids, arc.AID)
	}
	if arcsReply, err = s.arcClient.Arcs(c, &arcmdl.ArcsRequest{Aids: aids}); err != nil {
		log.Error("ListVideo s.arc.Archives3 error(%v)", err)
		return
	}
	for _, aid := range aids {
		if arc, ok := arcsReply.Arcs[aid]; ok && arc.IsNormal() {
			rs = append(rs, arc)
		}
	}
	if isFirst {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetVideoCache(c, p.Ps, rs, total)
		})
	}
	return
}

// FilterVideo filter video.
func (s *Service) FilterVideo(c context.Context, p *model.ParamFilter) (rs map[string][]*model.Filter, err error) {
	var (
		tmpRs                             map[string][]*model.Filter
		fv                                *model.FilterES
		fMap                              map[string]map[int64]*model.Filter
		matchs, games, teams, tags, years []*model.Filter
	)
	isAll := p.Year == 0 && p.Tag == 0 && p.Tid == 0 && p.Gid == 0 && p.Mid == 0
	if rs, err = s.dao.FVideoCache(c); err != nil {
		err = nil
	}
	if isAll && len(rs) > 0 {
		return
	}
	matchs, games, teams, tags, years = s.filterTop(c)
	tmpRs = make(map[string][]*model.Filter, 3)
	tmpRs[_typeMatch] = matchs
	tmpRs[_typeGame] = games
	tmpRs[_typeTeam] = teams
	tmpRs[_typeYear] = years
	tmpRs[_typeTag] = tags
	if fv, err = s.dao.FilterVideo(c, p); err != nil {
		log.Error("s.dao.FilterVideo error(%v)", err)
		return
	}
	fMap = s.filterMap(tmpRs)
	if rs, err = s.fmtES(fv, fMap); err != nil {
		log.Error("FilterVideo s.filterES error(%v)", err)
	}
	if isAll {
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetFVideoCache(c, rs)
		})
	} else {
		if len(rs[_typeMatch]) == 0 && len(rs[_typeGame]) == 0 && len(rs[_typeTeam]) == 0 && len(rs[_typeYear]) == 0 && len(rs[_typeTag]) == 0 {
			if tmpRs, err = s.dao.FVideoCache(c); err != nil {
				err = nil
			}
			if len(tmpRs) > 0 {
				rs = tmpRs
			}
		}
	}
	return
}

func (s *Service) filterTop(c context.Context) (matchs, games, teams, tags, years []*model.Filter) {
	var (
		matchErr, gameErr, teamErr, tagErr, yearErr error
	)
	group := &errgroup.Group{}
	group.Go(func() error {
		if matchs, matchErr = s.dao.Matchs(context.Background()); matchErr != nil {
			log.Error("s.dao.Matchs error %v", matchErr)
		}
		return nil
	})
	group.Go(func() error {
		if games, gameErr = s.dao.Games(context.Background()); gameErr != nil {
			log.Error("s.dao.Games error %v", gameErr)
		}
		return nil
	})
	group.Go(func() error {
		if teams, teamErr = s.dao.Teams(context.Background()); teamErr != nil {
			log.Error("s.dao.Teams error %v", teamErr)
		}
		return nil
	})
	group.Go(func() error {
		if tags, tagErr = s.dao.Tags(context.Background()); tagErr != nil {
			log.Error("s.dao.tags error %v", tagErr)
		}
		return nil
	})
	group.Go(func() error {
		if years, yearErr = s.dao.Years(context.Background()); yearErr != nil {
			log.Error("s.dao.Years error %v", yearErr)
		}
		return nil
	})
	group.Wait()
	if len(matchs) == 0 {
		matchs = _emptFilter
	}
	if len(games) == 0 {
		games = _emptFilter
	}
	if len(teams) == 0 {
		teams = _emptFilter
	}
	if len(years) == 0 {
		years = _emptFilter
	}
	if len(tags) == 0 {
		tags = _emptFilter
	}
	return
}

func (s *Service) fmtES(fv *model.FilterES, fMap map[string]map[int64]*model.Filter) (rs map[string][]*model.Filter, err error) {
	var (
		intMid, intGid, intTeam, intTag, intYear int64
		matchs, games, teams, tags, years        []*model.Filter
	)
	group := &errgroup.Group{}
	group.Go(func() error {
		for _, midGroup := range fv.GroupByMatch {
			if intMid, err = strconv.ParseInt(midGroup.Key, 10, 64); err != nil {
				err = nil
				continue
			}
			if match, ok := fMap[_typeMatch][intMid]; ok {
				matchs = append(matchs, match)
			}
		}
		return nil
	})
	group.Go(func() error {
		for _, gidGroup := range fv.GroupByGid {
			if intGid, err = strconv.ParseInt(gidGroup.Key, 10, 64); err != nil {
				err = nil
				continue
			}
			if game, ok := fMap[_typeGame][intGid]; ok {
				games = append(games, game)
			}
		}
		return nil
	})
	group.Go(func() error {
		for _, teamGroup := range fv.GroupByTeam {
			if intTeam, err = strconv.ParseInt(teamGroup.Key, 10, 64); err != nil {
				err = nil
				continue
			}
			if team, ok := fMap[_typeTeam][intTeam]; ok {
				teams = append(teams, team)
			}
		}
		return nil
	})
	group.Go(func() error {
		for _, tagGroup := range fv.GroupByTag {
			if intTag, err = strconv.ParseInt(tagGroup.Key, 10, 64); err != nil {
				err = nil
				continue
			}
			if tag, ok := fMap[_typeTag][intTag]; ok {
				tags = append(tags, tag)
			}
		}
		return nil
	})
	group.Go(func() error {
		for _, yearGroup := range fv.GroupByYear {
			if intYear, err = strconv.ParseInt(yearGroup.Key, 10, 64); err != nil {
				err = nil
				continue
			}
			if year, ok := fMap[_typeYear][intYear]; ok {
				years = append(years, year)
			}
		}
		return nil
	})
	group.Wait()
	rs = make(map[string][]*model.Filter, 5)
	if len(matchs) == 0 {
		matchs = _emptFilter
	} else {
		sort.Slice(matchs, func(i, j int) bool {
			return matchs[i].Rank > matchs[j].Rank || (matchs[i].Rank == matchs[j].Rank && matchs[i].ID < matchs[j].ID)
		})
	}
	if len(games) == 0 {
		games = _emptFilter
	} else {
		sort.Slice(games, func(i, j int) bool { return games[i].ID < games[j].ID })
	}
	if len(teams) == 0 {
		teams = _emptFilter
	} else {
		sort.Slice(teams, func(i, j int) bool { return teams[i].ID < teams[j].ID })
	}
	if len(years) == 0 {
		years = _emptFilter
	} else {
		sort.Slice(years, func(i, j int) bool { return years[i].ID < years[j].ID })
	}
	if len(tags) == 0 {
		tags = _emptFilter
	} else {
		sort.Slice(tags, func(i, j int) bool { return tags[i].ID < tags[j].ID })
	}
	rs[_typeMatch] = matchs
	rs[_typeGame] = games
	rs[_typeTeam] = teams
	rs[_typeTag] = tags
	rs[_typeYear] = years
	return
}

func (s *Service) filterMap(f map[string][]*model.Filter) (rs map[string]map[int64]*model.Filter) {
	var (
		match, game, team, tag, year                *model.Filter
		mapMatch, mapGame, mapTeam, mapTag, mapYear map[int64]*model.Filter
	)
	group := &errgroup.Group{}
	group.Go(func() error {
		mapMatch = make(map[int64]*model.Filter, len(f[_typeMatch]))
		for _, match = range f[_typeMatch] {
			if match != nil {
				mapMatch[match.ID] = match
			}
		}
		return nil
	})
	group.Go(func() error {
		mapGame = make(map[int64]*model.Filter, len(f[_typeGame]))
		for _, game = range f[_typeGame] {
			if game != nil {
				mapGame[game.ID] = game
			}
		}
		return nil
	})
	group.Go(func() error {
		mapTeam = make(map[int64]*model.Filter, len(f[_typeTeam]))
		for _, team = range f[_typeTeam] {
			if team != nil {
				mapTeam[team.ID] = team
			}
		}
		return nil
	})
	group.Go(func() error {
		mapTag = make(map[int64]*model.Filter, len(f[_typeTag]))
		for _, tag = range f[_typeTag] {
			if tag != nil {
				mapTag[tag.ID] = tag
			}
		}
		return nil
	})
	group.Go(func() error {
		mapYear = make(map[int64]*model.Filter, len(f[_typeYear]))
		for _, year = range f[_typeYear] {
			if year != nil {
				mapYear[year.ID] = year
			}
		}
		return nil
	})
	group.Wait()
	rs = make(map[string]map[int64]*model.Filter, 5)
	rs[_typeMatch] = mapMatch
	rs[_typeGame] = mapGame
	rs[_typeTeam] = mapTeam
	rs[_typeTag] = mapTag
	rs[_typeYear] = mapYear
	return
}

// Season season list.
func (s *Service) Season(c context.Context, p *model.ParamSeason) (rs []*model.Season, count int, err error) {
	var (
		seasons []*model.Season
		start   = (p.Pn - 1) * p.Ps
		end     = start + p.Ps - 1
	)
	if rs, count, err = s.dao.SeasonCache(c, start, end); err != nil || len(rs) == 0 {
		err = nil
		if seasons, err = s.dao.Season(c); err != nil {
			log.Error("s.dao.Season error(%v)", err)
			return
		}
		count = len(seasons)
		if count == 0 || count < start {
			rs = _emptSeason
			return
		}
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetSeasonCache(c, seasons, count)
		})
		if count > end+1 {
			rs = seasons[start : end+1]
		} else {
			rs = seasons[start:]
		}
	}
	return
}

// AppSeason  app season list.
func (s *Service) AppSeason(c context.Context, p *model.ParamSeason) (rs []*model.Season, count int, err error) {
	var (
		seasons []*model.Season
		start   = (p.Pn - 1) * p.Ps
		end     = start + p.Ps - 1
	)
	if rs, count, err = s.dao.SeasonMCache(c, start, end); err != nil || len(rs) == 0 {
		err = nil
		if seasons, err = s.dao.AppSeason(c); err != nil {
			log.Error("s.dao.AppSeason error(%v)", err)
			return
		}
		count = len(seasons)
		if count == 0 || count < start {
			rs = _emptSeason
			return
		}
		s.cache.Do(c, func(c context.Context) {
			s.dao.SetSeasonMCache(c, seasons, count)
		})
		if count > end+1 {
			rs = seasons[start : end+1]
		} else {
			rs = seasons[start:]
		}
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Rank > rs[j].Rank || (rs[i].Rank == rs[j].Rank && rs[i].Stime > rs[j].Stime)
	})
	return
}

// Contest contest data.
func (s *Service) Contest(c context.Context, mid, cid int64) (res *model.ContestDataPage, err error) {
	var (
		contest             *model.Contest
		contestData         []*model.ContestsData
		teams               map[int64]*model.Team
		season              map[int64]*model.Season
		teamErr, contestErr error
		games               []*model.Game
		gameMap             map[int64]*model.Game
	)
	if res, err = s.dao.GetCSingleData(c, cid); err != nil || res == nil {
		err = nil
		res = &model.ContestDataPage{}
		group, errCtx := errgroup.WithContext(c)
		group.Go(func() error {
			if contest, contestErr = s.dao.Contest(errCtx, cid); contestErr != nil {
				log.Error("SingleData.dao.Contest error(%v)", teamErr)
			}
			return contestErr
		})
		group.Go(func() error {
			if contestData, _ = s.dao.ContestData(errCtx, cid); err != nil {
				log.Error("SingleData.dao.ContestData error(%v)", teamErr)
			}
			return nil
		})
		err = group.Wait()
		if err != nil {
			return
		}
		if contest.ID == 0 {
			err = ecode.NothingFound
			return
		}
		if len(contestData) == 0 {
			contestData = _emptContestDetail
		}
		if teams, err = s.dao.EpTeams(c, []int64{contest.HomeID, contest.AwayID}); err != nil {
			log.Error("SingleData.dao.Teams error(%v)", err)
			err = nil
		}
		if season, err = s.dao.EpSeasons(c, []int64{contest.Sid}); err != nil {
			log.Error("SingleData.dao.EpSeasons error(%v)", err)
			err = nil
		}
		s.ContestInfos(contest, teams, season)
		res.Contest = contest
		res.Detail = contestData
		s.cache.Do(c, func(c context.Context) {
			s.dao.AddCSingleData(c, cid, res)
		})
	}
	if res.Contest.DataType == _lolType {
		games = s.lolGameMap.Data[res.Contest.MatchID]
	} else if res.Contest.DataType == _dotaType {
		games = s.dotaGameMap.Data[res.Contest.MatchID]
	}
	if len(games) > 0 {
		gameMap = make(map[int64]*model.Game, len(games))
		for _, game := range games {
			gameMap[game.ID] = game
		}
		for _, data := range res.Detail {
			if g, ok := gameMap[data.PointData]; ok {
				if g.Finished == true || (res.Contest.Etime > 0 && time.Now().Unix() > res.Contest.Etime) {
					data.GameStatus = 1
				} else if g.Finished == false {
					data.GameStatus = 2
				}
			}
		}
	}
	tmp := []*model.Contest{res.Contest}
	s.fmtContest(c, tmp, mid)
	return
}

// ContestInfos contest infos.
func (s *Service) ContestInfos(contest *model.Contest, teams map[int64]*model.Team, season map[int64]*model.Season) {
	if homeTeam, ok := teams[contest.HomeID]; ok {
		contest.HomeTeam = homeTeam
	} else {
		contest.HomeTeam = struct{}{}
	}
	if awayTeam, ok := teams[contest.AwayID]; ok {
		contest.AwayTeam = awayTeam
	} else {
		contest.AwayTeam = struct{}{}
	}
	if sea, ok := season[contest.Sid]; ok {
		contest.Season = sea
	} else {
		contest.Season = struct{}{}
	}
}

// Recent contest recents.
func (s *Service) Recent(c context.Context, mid int64, param *model.ParamCDRecent) (res []*model.Contest, err error) {
	var (
		teams  map[int64]*model.Team
		season map[int64]*model.Season
	)
	if res, err = s.dao.GetCRecent(c, param); err != nil || len(res) == 0 {
		err = nil
		if res, err = s.dao.ContestRecent(c, param.HomeID, param.AwayID, param.CID, param.Ps); err != nil {
			log.Error("ContestRecent.dao.ContestRecent error(%v)", err)
			return
		}
		if len(res) == 0 {
			res = _emptContest
			return
		}
		for _, contest := range res {
			if teams, err = s.dao.EpTeams(c, []int64{contest.HomeID, contest.AwayID}); err != nil {
				log.Error("SingleData.dao.Teams error(%v)", err)
				err = nil
			}
			if season, err = s.dao.EpSeasons(c, []int64{contest.Sid}); err != nil {
				log.Error("SingleData.dao.EpSeasons error(%v)", err)
				err = nil
			}
			s.ContestInfos(contest, teams, season)
			contest.SuccessTeaminfo = struct{}{}
		}
		s.cache.Do(c, func(c context.Context) {
			s.dao.AddCRecent(c, param, res)
		})
	}
	s.fmtContest(c, res, mid)
	return
}
