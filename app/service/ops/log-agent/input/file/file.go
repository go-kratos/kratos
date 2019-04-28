package file

import (
	"context"
	"fmt"
	"time"
	"path"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"os"

	"go-common/app/service/ops/log-agent/event"
	"go-common/app/service/ops/log-agent/input"
	"go-common/library/log"
)

type File struct {
	c        *Config
	output   chan<- *event.ProcessorEvent
	ctx      context.Context
	cancel   context.CancelFunc
	register *Registrar
}

func init() {
	err := input.Register("file", NewFile)
	if err != nil {
		panic(err)
	}
}

func NewFile(ctx context.Context, config interface{}, output chan<- *event.ProcessorEvent) (input.Input, error) {
	f := new(File)
	if c, ok := config.(*Config); !ok {
		return nil, fmt.Errorf("Error config for File Input")
	} else {
		if err := c.ConfigValidate(); err != nil {
			return nil, err
		}
		f.c = c
	}
	f.output = output
	f.ctx, f.cancel = context.WithCancel(ctx)

	// set config by ctx
	if f.c.ConfigPath == "" {
		configPath := ctx.Value("configPath")
		if configPath == nil {
			return nil, errors.New("can't get configPath from context")
		}
		f.c.ConfigPath = configPath.(string)
	}

	if f.c.ID == "" {
		hasher := sha1.New()
		hasher.Write([]byte(f.c.ConfigPath))
		f.c.ID = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	}

	if f.c.MetaPath == "" {
		f.c.MetaPath = ctx.Value("MetaPath").(string)
	}

	// init register
	if err := f.initRegister(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *File) Run() (err error) {
	log.Info("start collect log configured in %s", f.c.ConfigPath)

	f.scan()

	ticker := time.Tick(time.Duration(f.c.ScanFrequency))
	go func() {
		for {
			select {
			case <-f.ctx.Done():
				return
			case <-ticker:
				f.scan()
			}
		}
	}()
	return nil
}

func (f *File) Stop() {
	f.cancel()
}

func (f *File) Ctx() (context.Context) {
	return f.ctx
}

func (f *File) initRegister() error {
	path := path.Join(f.c.MetaPath, f.c.ID)

	register, err := NewRegistry(f.ctx, path)
	if err != nil {
		return err
	}
	f.register = register
	return nil
}

// Scan starts a scanGlob for each provided path/glob
func (f *File) scan() {
	paths := f.getFiles()

	// clean files older than
	if time.Duration(f.c.CleanFilesOlder) != 0 {
		f.cleanOldFiles(paths)
	}

	for path, info := range paths {
		select {
		case <-f.ctx.Done():
			return
		default:
		}

		newState, err := getFileState(path, info)
		if err != nil {
			log.Error("Skipping file %s due to error %s", path, err)
			continue
		}

		// Load last state
		lastState := f.register.FindPrevious(newState)

		// Ignores all files which fall under ignore_older
		if f.isIgnoreOlder(newState) {
			continue
		}

		// Decides if previous state exists
		if lastState.IsEmpty() {
			log.Info("Start harvester for new file: %s, inode: %d", newState.Source, newState.Inode)
			ctx := context.WithValue(f.ctx, "firstRun", true)
			err := f.startHarvester(ctx, f.c, f.register, newState, 0)
			if err != nil {
				log.Error("Harvester could not be started on new file: %s, Err: %s", newState.Source, err)
			}
		} else {
			ctx := context.WithValue(f.ctx, "firstRun", false)
			f.harvestExistingFile(ctx, f.c, f.register, newState, lastState)
		}
	}
}

func (f *File) cleanOldFiles(paths map[string]os.FileInfo) {
	if time.Duration(f.c.CleanFilesOlder) == 0 {
		return
	}

	var latestFile *State
	for path, info := range paths {
		newState, err := getFileState(path, info)

		if err != nil {
			log.Error("Skipping file %s due to error %s", path, err)
			continue
		}

		if latestFile == nil {
			latestFile = &newState
			continue
		}

		if newState.Fileinfo.ModTime().After(latestFile.Fileinfo.ModTime()) {
			// delete latestFile if newer file existing and modtime of latestFile is older than f.c.CleanFilesOlder
			if time.Since(latestFile.Fileinfo.ModTime()) > time.Duration(f.c.CleanFilesOlder) {
				if err := os.Remove(latestFile.Source); err != nil {
					log.Error("Failed to delete file %s", latestFile.Source)
				} else {
					log.Info("Delete file %s older than %s", latestFile.Source, time.Duration(f.c.CleanFilesOlder).String())
				}
			}
			latestFile = &newState
			continue
		}

		if newState.Fileinfo.ModTime().Before(latestFile.Fileinfo.ModTime()) {
			if time.Since(newState.Fileinfo.ModTime()) > time.Duration(f.c.CleanFilesOlder) {
				if err := os.Remove(newState.Source); err != nil {
					log.Error("Failed to delete file %s", newState.Source)
				} else {
					log.Info("Delete file %s older than %s", newState.Source, time.Duration(f.c.CleanFilesOlder))
				}
			}
		}
	}
}

// isIgnoreOlder checks if the given state reached ignore_older
func (f *File) isIgnoreOlder(state State) bool {
	// ignore_older is disable
	if f.c.IgnoreOlder == 0 {
		return false
	}

	modTime := state.Fileinfo.ModTime()
	if time.Since(modTime) > time.Duration(f.c.IgnoreOlder) {
		return true
	}

	return false
}
