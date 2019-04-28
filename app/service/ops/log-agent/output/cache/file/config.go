package file

import (
	"time"
	"errors"
	"path"

	xtime "go-common/library/time"
)

type Config struct {
	CacheFlushInterval xtime.Duration `tome:"cacheFlushInterval"`
	WriteBuffer        int           `tome:"writeBuffer"`
	Storage            string        `tome:"storage"`
	StorageMaxMB       int           `tome:"storageMaxMB"`
	FileBytes          int           `tome:"fileBytes"`
	Suffix             string        `tome:"suffix"`
	ReadBuffer         int           `tome:"readBuffer"`
	Index              string        `tome:"index"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of fileCache is nil")
	}

	if time.Duration(c.CacheFlushInterval) == 0 {
		c.CacheFlushInterval = xtime.Duration(time.Second * 5)
	}

	if c.WriteBuffer == 0 {
		c.WriteBuffer = 1024 * 1024 * 2 // 2M by default
	}

	if c.Storage == "" {
		return errors.New("storage settings for lancer output can't be nil")
	}

	if c.StorageMaxMB == 0 {
		c.StorageMaxMB = 5120
	}

	if c.FileBytes == 0 {
		c.FileBytes = 1024 * 1024 * 2 // 2M by default
	}

	if c.Suffix == "" {
		c.Suffix = ".log"
	}

	if c.ReadBuffer == 0 {
		c.ReadBuffer = 1024 * 1024 * 2 // 2M by default
	}

	if c.Index == "" {
		c.Index = path.Join(c.Storage, "output.index")
	}

	return nil
}
