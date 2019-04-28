package dao

import (
	"context"
	"go-common/app/service/live/rtc/internal/model"
	"go-common/library/cache/redis"
	"time"

	"go-common/app/service/live/rtc/internal/conf"
	xsql "go-common/library/database/sql"
)

// Dao dao
type Dao struct {
	c *conf.Config
	//mc    *memcache.Pool
	redis *redis.Pool
	db    *xsql.DB
}

// New init mysql db
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{
		c:     c,
		redis: redis.NewPool(c.Redis),
		db:    xsql.NewMySQL(c.MySQL),
	}
	return
}

// Close close the resource.
func (d *Dao) Close() {
	//d.mc.Close()
	d.redis.Close()
	d.db.Close()
}

// Ping dao ping
func (d *Dao) Ping(ctx context.Context) error {
	// TODO: add mc,redis... if you use
	return d.db.Ping(ctx)
}

func (d *Dao) GetMediaSource(ctx context.Context, channelID uint64) ([]*model.RtcMediaSource, error) {
	sql := "SELECT `id`,`channel_id`,`user_id`,`type`,`codec`,`media_specific` FROM `rtc_media_source` WHERE `channel_id` = ? AND `status` = 0"
	stmt := d.db.Prepared(sql)
	defer stmt.Close()
	rows, err := stmt.Query(ctx, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	source := make([]*model.RtcMediaSource, 0)
	for rows.Next() {
		s := &model.RtcMediaSource{}
		if err = rows.Scan(&s.SourceID, &s.ChannelID, &s.UserID, &s.Type, &s.Codec, &s.MediaSpecific); err != nil {
			return nil, err
		}
		source = append(source, s)
	}
	return source, nil
}

func (d *Dao) CreateCall(ctx context.Context, call *model.RtcCall) (uint32, error) {
	sql := "INSERT INTO `rtc_call`(`user_id`,`channel_id`,`version`,`token`,`join_time`,`leave_time`,`status`) VALUES(?,?,?,?,?,?,?)"
	stmt := d.db.Prepared(sql)
	defer stmt.Close()
	r, err := stmt.Exec(ctx, call.UserID, call.ChannelID, call.Version, call.Token, call.JoinTime, call.LeaveTime, call.Status)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}
	call.CallID = uint32(id)
	return call.CallID, nil
}

func (d *Dao) UpdateCallStatus(ctx context.Context, channelID uint64, callID uint32, userID uint64, leave time.Time, status uint8) error {
	sql := "UPDATE `rtc_call` SET `leave_time` = ?,`status` = ? WHERE `id` = ? AND `user_id` = ? LIMIT 1"
	stmt := d.db.Prepared(sql)
	defer stmt.Close()
	_, err := stmt.Exec(ctx, leave, status, callID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) UpdateMediaSourceStatus(ctx context.Context, channelID uint64, callID uint32, userID uint64, status uint8) error {
	sql := "UPDATE `rtc_media_source` SET `status` = ? WHERE `call_id` = ? AND `channel_id` = ? AND `user_id` = ?"
	stmt := d.db.Prepared(sql)
	defer stmt.Close()
	_, err := stmt.Exec(ctx, status, callID, channelID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) CreateMediaSource(ctx context.Context, source *model.RtcMediaSource) (uint32, error) {
	sql := "INSERT INTO `rtc_media_source`(`channel_id`,`user_id`,`type`,`codec`,`media_specific`,`status`) VALUES(?,?,?,?,?,?)"
	stmt := d.db.Prepared(sql)
	defer stmt.Close()
	r, err := stmt.Exec(ctx, source.ChannelID, source.UserID, source.Type, source.Codec, source.MediaSpecific, source.Status)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint32(id), nil
}

func (d *Dao) CreateMediaPublish(ctx context.Context, publish *model.RtcMediaPublish) error {
	mixConfigSql := "REPLACE INTO `rtc_mix_config`(`call_id`,`config`) VALUES(?,?)"
	mixConfigStmt := d.db.Prepared(mixConfigSql)
	defer mixConfigStmt.Close()
	var err error
	mixConfigResult, err := mixConfigStmt.Exec(ctx, publish.CallID, publish.MixConfig)
	if err != nil {
		return err
	}
	_, err = mixConfigResult.LastInsertId()
	if err != nil {
		return err
	}
	publishSql := "REPLACE INTO `rtc_media_publish`(`call_id`,`channel_id`,`user_id`,`switch`,`width`,`height`,`frame_rate`,`video_codec`,`video_profile`,`channel`,`sample_rate`,`audio_codec`,`bitrate`) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)"
	publishStmt := d.db.Prepared(publishSql)
	defer publishStmt.Close()
	_, err = publishStmt.Exec(ctx, publish.CallID, publish.ChannelID, publish.UserID, publish.Switch,
		publish.Width, publish.Height, publish.FrameRate, publish.VideoCodec, publish.VideoProfile,
		publish.Channel, publish.SampleRate, publish.AudioCodec, publish.Bitrate)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) GetMediaPublishConfig(ctx context.Context, channelID uint64, callID uint32) (*model.RtcMediaPublish, error) {
	publishSql := "SELECT `user_id`,`switch`,`width`,`height`,`frame_rate`,`video_codec`,`video_profile`,`channel`,`sample_rate`,`audio_codec`,`bitrate`,`mix_config_id` FROM `rtc_media_publish` WHERE `call_id` = ? AND `channel_id` = ? LIMIT 1"
	publishStmt := d.db.Prepared(publishSql)
	defer publishStmt.Close()
	publishRow := publishStmt.QueryRow(ctx, callID, channelID)
	var publish model.RtcMediaPublish
	var mixConfigID uint32

	if err := publishRow.Scan(&publish.UserID, &publish.Switch, &publish.Width, &publish.Height, &publish.FrameRate,
		&publish.VideoCodec, &publish.VideoProfile, &publish.Channel, &publish.SampleRate,
		&publish.AudioCodec, &publish.Bitrate, &mixConfigID); err != nil {
		return nil, err
	}

	mixConfigSql := "SELECT `config` FROM `rtc_mix_config` WHERE `id` = ? "
	mixConfigStmt := d.db.Prepared(mixConfigSql)
	defer mixConfigStmt.Close()
	mixConfigRow := mixConfigStmt.QueryRow(ctx, mixConfigID)
	if err := mixConfigRow.Scan(&publish.MixConfig); err != nil {
		return nil, err
	}
	return &publish, nil
}

func (d *Dao) UpdateMediaPublishConfig(ctx context.Context, channelID uint64, callID uint32, config string) error {
	sql := "UPDATE `rtc_mix_config` SET `config` = ? WHERE  `call_id` = ? LIMIT 1"
	stmt := d.db.Prepared(sql)
	defer stmt.Close()
	_, err := stmt.Exec(ctx, config, callID)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) TerminateStream(ctx context.Context, channelID uint64, callID uint32) error {
	sql := "UPDATE `rtc_media_publish` SET `switch` = 0 WHERE `call_id` = ? LIMIT 1"
	stmt := d.db.Prepared(sql)
	defer stmt.Close()
	_, err := stmt.Exec(ctx, callID)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) GetChannelIP(ctx context.Context, channelID uint64) ([]string, error) {
	sql := "SELECT `ip` FROM `rtc_call` WHERE `channel_id` = ? AND `status` = 0"
	stmt := d.db.Prepared(sql)
	defer stmt.Close()
	rows, err := stmt.Query(ctx, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]string, 0)
	for rows.Next() {
		var ip string
		if err = rows.Scan(&ip); err != nil {
			return nil, err
		}
		result = append(result, ip)
	}
	return result, nil
}

func (d *Dao) GetToken(ctx context.Context, channelID uint64, callID uint32) (string, error) {
	sql := "SELECT `token` FROM `rtc_call` WHERE `id` = ? AND `channel_id` = ?"
	stmt := d.db.Prepared(sql)
	defer stmt.Close()
	row := stmt.QueryRow(ctx, callID, channelID)
	var token string
	err := row.Scan(&token)
	if err == xsql.ErrNoRows {
		err = nil
	}
	return token, err
}
