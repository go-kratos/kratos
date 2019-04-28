package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/interface/main/esports/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_matchsSQL       = "SELECT id,title,sub_title,logo,rank FROM es_matchs  WHERE status=0 order by rank DESC , ID ASC"
	_gamesSQL        = "SELECT id,title,sub_title,logo FROM es_games   WHERE status=0 order by id ASC"
	_teamsSQL        = "SELECT id,title,sub_title,logo FROM es_teams   WHERE is_deleted=0 order by id ASC"
	_tagsSQL         = "SELECT id,name FROM es_tags  WHERE status=0 order by id ASC"
	_yearsSQL        = "SELECT distinct year as id, year FROM es_year_map  WHERE is_deleted=0 order by id ASC"
	_dayContestSQL   = "SELECT FROM_UNIXTIME(stime, '%Y-%m-%d') as s,count(1) as c FROM `es_contests` WHERE status=0  AND  stime >= ? and stime <= ? GROUP BY s ORDER BY stime"
	_seasonSQL       = "SELECT id,mid,title,sub_title,stime,etime,sponsor,logo,dic,ctime,mtime,status,rank,is_app,url,data_focus,focus_url FROM  es_seasons  WHERE  status = 0  ORDER BY stime DESC"
	_epSeasonSQL     = "SELECT id,mid,title,sub_title,stime,etime,sponsor,logo,dic,ctime,mtime,status,rank,is_app,url,data_focus,focus_url FROM  es_seasons  WHERE  status = 0 and id in (%s) ORDER BY stime DESC"
	_seasonMSQL      = "SELECT id,mid,title,sub_title,stime,etime,sponsor,logo,dic,ctime,mtime,status,rank,is_app,url,data_focus,focus_url FROM  es_seasons  WHERE  status = 0 AND is_app = 1  ORDER BY rank DESC,stime DESC"
	_seasonsSQL      = "SELECT id,title,sub_title,logo,url,data_focus,focus_url FROM  es_seasons WHERE status = 0  ORDER BY stime DESC"
	_contestSQL      = "SELECT id,game_stage,stime,etime,home_id,away_id,home_score,away_score,live_room,aid,collection,game_state,dic,ctime,mtime,status,sid,mid,special,special_name,special_tips,success_team,special_image,playback,collection_url,live_url,data_type,match_id FROM `es_contests` WHERE id= ?"
	_contestsSQL     = "SELECT id,game_stage,stime,etime,home_id,away_id,home_score,away_score,live_room,aid,collection,game_state,dic,ctime,mtime,status,sid,mid,special,special_name,special_tips,success_team,special_image,playback,collection_url,live_url,data_type,match_id FROM `es_contests` WHERE id in (%s) ORDER BY ID ASC"
	_contestLeidaSQL = "SELECT id,game_stage,stime,etime,home_id,away_id,home_score,away_score,live_room,aid,collection,game_state,dic,ctime,mtime,status,sid,mid,special,special_name,special_tips,success_team,special_image,playback,collection_url,live_url,data_type,match_id FROM `es_contests` WHERE match_id > 0 and status = 0"
	_moduleSQL       = "SELECT id,ma_id,name,oids FROM `es_matchs_module` WHERE id = ? AND status = 0"
	_activeSQL       = "SELECT id,mid,sid,background,live_id,intr,focus,url,back_color,color_step,h5_background,h5_back_color,intr_logo,intr_title,intr_text,h5_focus,h5_url FROM es_matchs_active WHERE id = ? AND `status`= 0"
	_modulesSQL      = "SELECT id,ma_id,name,oids FROM `es_matchs_module` WHERE ma_id = ? AND status = 0 ORDER BY ID ASC"
	_pDetailSQL      = "SELECT ma_id,game_type,stime,etime FROM es_matchs_detail WHERE id = ? AND `status` = 0"
	_actDetail       = "SELECT id,ma_id,game_type,stime,etime,score_id,game_stage,knockout_type,winner_type,online FROM es_matchs_detail WHERE ma_id = ? AND status = 0"
	_treeSQL         = "SELECT id,ma_id,mad_id,pid,root_id,game_rank,mid FROM es_matchs_tree WHERE mad_id = ? AND is_deleted=0 ORDER BY root_id ASC,pid ASC,game_rank ASC"
	_teamsInSQL      = "SELECT id,title,sub_title,logo FROM es_teams WHERE is_deleted=0 AND id in (%s)"
	_kDetailsSQL     = "SELECT id,ma_id,game_type,stime,etime,online FROM es_matchs_detail WHERE `status` = 0 AND game_type = 2"
	_contestDataSQL  = "SELECT id,cid,url,point_data FROM `es_contests_data` WHERE cid = ? AND is_deleted = 0"
	_contestRecent   = "SELECT id,game_stage,stime,etime,home_id,away_id,home_score,away_score,live_room,aid,collection,game_state,dic,ctime,mtime,status,sid,mid,special,special_name,special_tips,success_team,special_image,playback,collection_url,live_url,data_type FROM es_contests WHERE ( `status` = 0 AND home_id = ? AND away_id = ? ) OR ( `status` = 0 AND home_id = ? AND away_id = ? ) ORDER BY stime DESC LIMIT ?"
)

// logoURL convert logo url to full url.
func logoURL(uri string) (logo string) {
	if uri == "" {
		return
	}
	logo = uri
	if strings.Index(uri, "http://") == 0 || strings.Index(uri, "//") == 0 {
		return
	}
	if len(uri) >= 10 && uri[:10] == "/templets/" {
		return
	}
	if strings.HasPrefix(uri, "group1") {
		logo = "//i0.hdslb.com/" + uri
		return
	}
	if pos := strings.Index(uri, "/uploads/"); pos != -1 && (pos == 0 || pos == 3) {
		logo = uri[pos+8:]
	}
	logo = strings.Replace(logo, "{IMG}", "", -1)
	logo = "//i0.hdslb.com" + logo
	return
}

// Matchs filter matchs.
func (d *Dao) Matchs(c context.Context) (res []*model.Filter, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _matchsSQL); err != nil {
		log.Error("Match:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Filter)
		if err = rows.Scan(&r.ID, &r.Title, &r.SubTitle, &r.Logo, &r.Rank); err != nil {
			log.Error("Match:row.Scan() error(%v)", err)
			return
		}
		r.Logo = logoURL(r.Logo)
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Games filter games.
func (d *Dao) Games(c context.Context) (res []*model.Filter, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _gamesSQL); err != nil {
		log.Error("Games:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Filter)
		if err = rows.Scan(&r.ID, &r.Title, &r.SubTitle, &r.Logo); err != nil {
			log.Error("Games:row.Scan() error(%v)", err)
			return
		}
		r.Logo = logoURL(r.Logo)
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Teams filter teams.
func (d *Dao) Teams(c context.Context) (res []*model.Filter, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _teamsSQL); err != nil {
		log.Error("Teams:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Filter)
		if err = rows.Scan(&r.ID, &r.Title, &r.SubTitle, &r.Logo); err != nil {
			log.Error("Teams:row.Scan() error(%v)", err)
			return
		}
		r.Logo = logoURL(r.Logo)
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Tags filter Tags.
func (d *Dao) Tags(c context.Context) (res []*model.Filter, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _tagsSQL); err != nil {
		log.Error("Tags:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Filter)
		if err = rows.Scan(&r.ID, &r.Title); err != nil {
			log.Error("Tags:row.Scan() error(%v)", err)
			return
		}
		r.Logo = logoURL(r.Logo)
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Years filter years.
func (d *Dao) Years(c context.Context) (res []*model.Filter, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _yearsSQL); err != nil {
		log.Error("Years:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Year)
		if err = rows.Scan(&r.ID, &r.Year); err != nil {
			log.Error("Years:row.Scan() error(%v)", err)
			return
		}
		res = append(res, &model.Filter{ID: r.ID, Title: strconv.FormatInt(r.Year, 10)})
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Calendar  calendar count.
func (d *Dao) Calendar(c context.Context, stime, etime int64) (res []*model.Calendar, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _dayContestSQL, stime, etime); err != nil {
		log.Error("Calendar:d.db.Query(%d,%d) error(%v)", stime, etime, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Calendar)
		if err = rows.Scan(&r.Stime, &r.Count); err != nil {
			log.Error("Calendar:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	return
}

// Season season list.
func (d *Dao) Season(c context.Context) (res []*model.Season, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _seasonSQL); err != nil {
		log.Error("Contest:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Season)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Title, &r.SubTitle, &r.Stime, &r.Etime, &r.Sponsor, &r.Logo, &r.Dic, &r.Ctime,
			&r.Mtime, &r.Status, &r.Rank, &r.IsApp, &r.URL, &r.DataFocus, &r.FocusURL); err != nil {
			log.Error("Contest:row.Scan() error(%v)", err)
			return
		}
		r.Logo = logoURL(r.Logo)
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Module active module
func (d *Dao) Module(c context.Context, mmid int64) (mod *model.Module, err error) {
	mod = &model.Module{}
	row := d.db.QueryRow(c, _moduleSQL, mmid)
	if err = row.Scan(&mod.ID, &mod.MAid, &mod.Name, &mod.Oids); err != nil {
		if err == sql.ErrNoRows {
			mod = nil
			err = nil
		} else {
			log.Error("Esport dao Module:row.Scan error(%v)", err)
		}
	}
	return
}

// Modules active module
func (d *Dao) Modules(c context.Context, aid int64) (mods []*model.Module, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _modulesSQL, aid); err != nil {
		log.Error("Esport dao Modules:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Module)
		if err = rows.Scan(&r.ID, &r.MAid, &r.Name, &r.Oids); err != nil {
			log.Error("Esport dao Modules:row.Scan() error(%v)", err)
			return
		}
		mods = append(mods, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("Esport dao Modules.Err() error(%v)", err)
	}
	return
}

// Trees match tree
func (d *Dao) Trees(c context.Context, madID int64) (mods []*model.Tree, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _treeSQL, madID); err != nil {
		log.Error("Esport dao Trees:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Tree)
		if err = rows.Scan(&r.ID, &r.MaID, &r.MadID, &r.Pid, &r.RootID, &r.GameRank, &r.Mid); err != nil {
			log.Error("Esport dao Trees:row.Scan() error(%v)", err)
			return
		}
		mods = append(mods, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("Esport dao Trees.Err() error(%v)", err)
	}
	return
}

// Active matchs active
func (d *Dao) Active(c context.Context, aid int64) (mod *model.Active, err error) {
	mod = &model.Active{}
	row := d.db.QueryRow(c, _activeSQL, aid)
	if err = row.Scan(&mod.ID, &mod.Mid, &mod.Sid, &mod.Background, &mod.Liveid, &mod.Intr, &mod.Focus, &mod.URL, &mod.BackColor, &mod.ColorStep, &mod.H5Background, &mod.H5BackColor, &mod.IntrLogo, &mod.IntrTitle, &mod.IntrText, &mod.H5Focus, &mod.H5Url); err != nil {
		if err == sql.ErrNoRows {
			mod = nil
			err = nil
		} else {
			log.Error("Esport dao Active:row.Scan error(%v)", err)
		}
	}
	return
}

// PActDetail poin match detail
func (d *Dao) PActDetail(c context.Context, id int64) (mod *model.ActiveDetail, err error) {
	mod = &model.ActiveDetail{}
	row := d.db.QueryRow(c, _pDetailSQL, id)
	if err = row.Scan(&mod.Maid, &mod.GameType, &mod.STime, &mod.ETime); err != nil {
		if err == sql.ErrNoRows {
			mod = nil
			err = nil
		} else {
			log.Error("Esport dao Contest:row.Scan error(%v)", err)
		}
	}
	return
}

// ActDetail data module
func (d *Dao) ActDetail(c context.Context, aid int64) (actDetail []*model.ActiveDetail, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _actDetail, aid); err != nil {
		log.Error("Esport dao Modules:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.ActiveDetail)
		if err = rows.Scan(&r.ID, &r.Maid, &r.GameType, &r.STime, &r.ETime, &r.ScoreID, &r.GameStage, &r.KnockoutType, &r.WinnerType, &r.Online); err != nil {
			log.Error("Esport dao ActDetail:row.Scan() error(%v)", err)
			return
		}
		actDetail = append(actDetail, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("Esport dao ActDetail.Err() error(%v)", err)
	}
	return
}

// AppSeason season match list.
func (d *Dao) AppSeason(c context.Context) (res []*model.Season, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _seasonMSQL); err != nil {
		log.Error("Contest:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Season)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Title, &r.SubTitle, &r.Stime, &r.Etime, &r.Sponsor, &r.Logo,
			&r.Dic, &r.Ctime, &r.Mtime, &r.Status, &r.Rank, &r.IsApp, &r.URL, &r.DataFocus, &r.FocusURL); err != nil {
			log.Error("Contest:row.Scan() error(%v)", err)
			return
		}
		r.Logo = logoURL(r.Logo)
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// SeasonAll all season list.
func (d *Dao) SeasonAll(c context.Context) (res []*model.Filter, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _seasonsSQL); err != nil {
		log.Error("SeasonAll:d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.Filter)
		if err = rows.Scan(&r.ID, &r.Title, &r.SubTitle, &r.Logo, &r.URL, &r.DataFocus, &r.FocusURL); err != nil {
			log.Error("SeasonAll:row.Scan() error(%v)", err)
			return
		}
		r.Logo = logoURL(r.Logo)
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// Contest get contest by id.
func (d *Dao) Contest(c context.Context, cid int64) (res *model.Contest, err error) {
	res = &model.Contest{}
	row := d.db.QueryRow(c, _contestSQL, cid)
	if err = row.Scan(&res.ID, &res.GameStage, &res.Stime, &res.Etime, &res.HomeID, &res.AwayID, &res.HomeScore, &res.AwayScore,
		&res.LiveRoom, &res.Aid, &res.Collection, &res.GameState, &res.Dic, &res.Ctime, &res.Mtime, &res.Status, &res.Sid, &res.Mid,
		&res.Special, &res.SpecialName, &res.SpecialTips, &res.SuccessTeam, &res.SpecialImage, &res.Playback, &res.CollectionURL,
		&res.LiveURL, &res.DataType, &res.MatchID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("Contest:row.Scan error(%v)", err)
		}
	}
	return
}

// ContestRecent get recent contest
func (d *Dao) ContestRecent(c context.Context, homeid, awayid, contestid, ps int64) (res []*model.Contest, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _contestRecent, homeid, awayid, awayid, homeid, ps+1); err != nil {
		log.Error("ContestRecent: db.Exec(%s) error(%v)", _contestRecent, err)
		return
	}
	defer rows.Close()
	res = make([]*model.Contest, 0)
	for rows.Next() {
		r := new(model.Contest)
		if err = rows.Scan(&r.ID, &r.GameStage, &r.Stime, &r.Etime, &r.HomeID, &r.AwayID, &r.HomeScore, &r.AwayScore,
			&r.LiveRoom, &r.Aid, &r.Collection, &r.GameState, &r.Dic, &r.Ctime, &r.Mtime, &r.Status, &r.Sid, &r.Mid,
			&r.Special, &r.SpecialName, &r.SpecialTips, &r.SuccessTeam, &r.SpecialImage, &r.Playback, &r.CollectionURL, &r.LiveURL, &r.DataType); err != nil {
			log.Error("Contests:row.Scan() error(%v)", err)
			return
		}
		if r.ID != contestid && len(res) != int(ps) {
			res = append(res, r)
		}
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// ContestData get contest by id.
func (d *Dao) ContestData(c context.Context, cid int64) (res []*model.ContestsData, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _contestDataSQL, cid); err != nil {
		log.Error("ContestsData: db.Exec(%s) error(%v)", _contestDataSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := new(model.ContestsData)
		if err = rows.Scan(&data.ID, &data.Cid, &data.URL, &data.PointData); err != nil {
			log.Error("ContestsData:row.Scan() error(%v)", err)
			return
		}
		res = append(res, data)
	}
	if err = rows.Err(); err != nil {
		log.Error("ContestssData rows.Err() error(%v)", err)
	}
	return
}

// ContestDatas contest datas.
func (d *Dao) ContestDatas(c context.Context) (res []*model.Contest, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _contestLeidaSQL); err != nil {
		log.Error("Contests: db.Exec(%s) error(%v)", _contestLeidaSQL, err)
		return
	}
	defer rows.Close()
	res = make([]*model.Contest, 0)
	for rows.Next() {
		r := new(model.Contest)
		if err = rows.Scan(&r.ID, &r.GameStage, &r.Stime, &r.Etime, &r.HomeID, &r.AwayID, &r.HomeScore, &r.AwayScore,
			&r.LiveRoom, &r.Aid, &r.Collection, &r.GameState, &r.Dic, &r.Ctime, &r.Mtime, &r.Status, &r.Sid, &r.Mid,
			&r.Special, &r.SpecialName, &r.SpecialTips, &r.SuccessTeam, &r.SpecialImage, &r.Playback, &r.CollectionURL, &r.LiveURL, &r.DataType, &r.MatchID); err != nil {
			log.Error("Contests:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// RawEpContests get contests by ids.
func (d *Dao) RawEpContests(c context.Context, cids []int64) (res map[int64]*model.Contest, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, fmt.Sprintf(_contestsSQL, xstr.JoinInts(cids))); err != nil {
		log.Error("Contests: db.Exec(%s) error(%v)", xstr.JoinInts(cids), err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Contest, len(cids))
	for rows.Next() {
		r := new(model.Contest)
		if err = rows.Scan(&r.ID, &r.GameStage, &r.Stime, &r.Etime, &r.HomeID, &r.AwayID, &r.HomeScore, &r.AwayScore,
			&r.LiveRoom, &r.Aid, &r.Collection, &r.GameState, &r.Dic, &r.Ctime, &r.Mtime, &r.Status, &r.Sid, &r.Mid,
			&r.Special, &r.SpecialName, &r.SpecialTips, &r.SuccessTeam, &r.SpecialImage, &r.Playback, &r.CollectionURL,
			&r.LiveURL, &r.DataType, &r.MatchID); err != nil {
			log.Error("Contests:row.Scan() error(%v)", err)
			return
		}
		res[r.ID] = r
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// RawEpSeasons get seasons by ids.
func (d *Dao) RawEpSeasons(c context.Context, sids []int64) (res map[int64]*model.Season, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, fmt.Sprintf(_epSeasonSQL, xstr.JoinInts(sids))); err != nil {
		log.Error("Contests: db.Exec(%s) error(%v)", xstr.JoinInts(sids), err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Season, len(sids))
	for rows.Next() {
		r := new(model.Season)
		if err = rows.Scan(&r.ID, &r.Mid, &r.Title, &r.SubTitle, &r.Stime, &r.Etime, &r.Sponsor, &r.Logo, &r.Dic, &r.Ctime,
			&r.Mtime, &r.Status, &r.Rank, &r.IsApp, &r.URL, &r.DataFocus, &r.FocusURL); err != nil {
			log.Error("Contest:row.Scan() error(%v)", err)
			return
		}
		res[r.ID] = r
	}
	if err = rows.Err(); err != nil {
		log.Error("rows.Err() error(%v)", err)
	}
	return
}

// ActTeams get teams by ids in
func (d *Dao) ActTeams(c context.Context, tids []int64) (res []*model.Team, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, fmt.Sprintf(_teamsInSQL, xstr.JoinInts(tids))); err != nil {
		log.Error("ActTeams: db.Exec(%s) error(%v)", xstr.JoinInts(tids), err)
		return
	}
	defer rows.Close()
	res = make([]*model.Team, len(tids))
	for rows.Next() {
		r := new(model.Team)
		if err = rows.Scan(&r.ID, &r.Title, &r.SubTitle, &r.Logo); err != nil {
			log.Error("ActTeams:row.Scan() error(%v)", err)
			return
		}
		res = append(res, r)
	}
	if err = rows.Err(); err != nil {
		log.Error("ActTeams rows.Err() error(%v)", err)
	}
	return
}

// RawEpTeams get seasons by ids.
func (d *Dao) RawEpTeams(c context.Context, tids []int64) (res map[int64]*model.Team, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, fmt.Sprintf(_teamsInSQL, xstr.JoinInts(tids))); err != nil {
		log.Error("RawEpTeams: db.Exec(%s) error(%v)", xstr.JoinInts(tids), err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Team, len(tids))
	for rows.Next() {
		r := new(model.Team)
		if err = rows.Scan(&r.ID, &r.Title, &r.SubTitle, &r.Logo); err != nil {
			log.Error("RawEpTeams:row.Scan() error(%v)", err)
			return
		}
		res[r.ID] = r
	}
	if err = rows.Err(); err != nil {
		log.Error("RawEpTeams.Err() error(%v)", err)
	}
	return
}

// KDetails knockout detail
func (d *Dao) KDetails(c context.Context) (res []*model.ActiveDetail, err error) {
	var (
		rows *xsql.Rows
	)
	if rows, err = d.db.Query(c, _kDetailsSQL); err != nil {
		log.Error("ActPDetails: db.Exec(%s) error(%v)", _kDetailsSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		detail := new(model.ActiveDetail)
		if err = rows.Scan(&detail.ID, &detail.Maid, &detail.GameType, &detail.STime, &detail.ETime, &detail.Online); err != nil {
			log.Error("KDetails:row.Scan() error(%v)", err)
			return
		}
		res = append(res, detail)
	}
	if err = rows.Err(); err != nil {
		log.Error("KDetails rows.Err() error(%v)", err)
	}
	return
}
