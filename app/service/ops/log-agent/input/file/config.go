package file

import (
	"errors"
	"time"

	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Paths           []string               `toml:"paths"`
	Symlinks        bool                   `toml:"symlinks"`
	AppId           string                 `toml:"appId"`
	LogId           string                 `toml:"logId"`
	ConfigPath      string                 `toml:"-"`
	MetaPath        string                 `toml:"-"`
	ID              string                 `toml:"-"`
	ReadFrom        string                 `toml:"readFrom"`
	MaxLength       int                    `toml:"maxLength"`
	IgnoreOlder     xtime.Duration         `toml:"ignoreOlder"`
	CleanFilesOlder xtime.Duration         `toml:"cleanFilesOlder"`
	ScanFrequency   xtime.Duration         `toml:"scanFrequency"`
	CleanInactive   xtime.Duration         `toml:"cleanInactive"`
	HarvesterTTL    xtime.Duration         `toml:"harvesterTTL"` // harvester will stop itself if inactive longer than HarvesterTTL
	Multiline       *MultilineConf         `toml:"multiline"`
	Timeout         xtime.Duration         `toml:"timeout"`
	Fields          map[string]interface{} `toml:"fields"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of file Input is nil")
	}

	if len(c.Paths) == 0 {
		return errors.New("paths of file Input can't be nil")
	}

	if c.LogId == "" {
		c.LogId = "000161"
	}

	if c.AppId == "" {
		return errors.New("appId of file Input can't be nil")
	}

	if c.IgnoreOlder == 0 {
		c.IgnoreOlder = xtime.Duration(time.Hour * 24)
	}

	if c.ScanFrequency == 0 {
		c.ScanFrequency = xtime.Duration(time.Second * 10)
	}

	// Note: CleanInactive should be greater chan ignore_older + scan_frequency
	if c.CleanInactive == 0 {
		c.CleanInactive = xtime.Duration(time.Hour * 24 * 7)
	}

	if c.CleanInactive < c.IgnoreOlder+c.ScanFrequency {
		return errors.New("CleanInactive must be greater than ScanFrequency + IgnoreOlder")
	}

	if c.HarvesterTTL == 0 {
		c.HarvesterTTL = xtime.Duration(time.Hour * 1)
	}

	if c.Timeout == 0 {
		c.Timeout = xtime.Duration(time.Second * 5)
	}

	if c.ReadFrom != "" && c.ReadFrom != "newest" && c.ReadFrom != "oldest" {
		return errors.New("ReadFrom of file input can only be newest or oldest")
	}

	if c.ReadFrom == "" {
		c.ReadFrom = "newest"
	}

	if c.MaxLength == 0 || c.MaxLength > 1024*10*64 {
		c.MaxLength = 1024 * 10 * 64
	}
	// Symlinks is always disabled
	c.Symlinks = false

	if c.Multiline != nil {
		if err := c.Multiline.ConfigValidate(); err != nil {
			return err
		}
	}

	return nil
}

func DecodeConfig(md toml.MetaData, primValue toml.Primitive) (c interface{}, err error) {
	c = new(Config)
	if err = md.PrimitiveDecode(primValue, c); err != nil {
		return nil, err
	}
	return c, nil
}
