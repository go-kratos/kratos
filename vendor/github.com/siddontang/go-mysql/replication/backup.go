package replication

import (
	"context"
	"io"
	"os"
	"path"
	"time"

	"github.com/juju/errors"
	. "github.com/siddontang/go-mysql/mysql"
)

// Like mysqlbinlog remote raw backup
// Backup remote binlog from position (filename, offset) and write in backupDir
func (b *BinlogSyncer) StartBackup(backupDir string, p Position, timeout time.Duration) error {
	if timeout == 0 {
		// a very long timeout here
		timeout = 30 * 3600 * 24 * time.Second
	}

	// Force use raw mode
	b.parser.SetRawMode(true)

	os.MkdirAll(backupDir, 0755)

	s, err := b.StartSync(p)
	if err != nil {
		return errors.Trace(err)
	}

	var filename string
	var offset uint32

	var f *os.File
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		e, err := s.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			return nil
		}

		if err != nil {
			return errors.Trace(err)
		}

		offset = e.Header.LogPos

		if e.Header.EventType == ROTATE_EVENT {
			rotateEvent := e.Event.(*RotateEvent)
			filename = string(rotateEvent.NextLogName)

			if e.Header.Timestamp == 0 || offset == 0 {
				// fake rotate event
				continue
			}
		} else if e.Header.EventType == FORMAT_DESCRIPTION_EVENT {
			// FormateDescriptionEvent is the first event in binlog, we will close old one and create a new

			if f != nil {
				f.Close()
			}

			if len(filename) == 0 {
				return errors.Errorf("empty binlog filename for FormateDescriptionEvent")
			}

			f, err = os.OpenFile(path.Join(backupDir, filename), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return errors.Trace(err)
			}

			// write binlog header fe'bin'
			if _, err = f.Write(BinLogFileHeader); err != nil {
				return errors.Trace(err)
			}

		}

		if n, err := f.Write(e.RawData); err != nil {
			return errors.Trace(err)
		} else if n != len(e.RawData) {
			return errors.Trace(io.ErrShortWrite)
		}
	}

	return nil
}
