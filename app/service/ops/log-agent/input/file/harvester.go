package file

import (
	"fmt"
	"os"
	"errors"
	"io"
	"time"
	"context"
	"bytes"
	"strconv"
	"regexp"
	"go-common/app/service/ops/log-agent/event"
	"go-common/library/log"
	"go-common/app/service/ops/log-agent/pkg/lancerroute"
)

type Source interface {
	io.ReadCloser
	Name() string
	Stat() (os.FileInfo, error)
	Continuable() bool // can we continue processing after EOF?
	HasState() bool    // does this source have a state?
}

type Hfile struct {
	*os.File
}

var lineReadTimeout = errors.New("lineReadTimeout")

func (Hfile) Continuable() bool { return true }
func (Hfile) HasState() bool    { return true }

// Harvester contains all harvester related data
type Harvester struct {
	config          *Config
	source          *os.File
	ctx             context.Context
	cancel          context.CancelFunc
	state           State
	register        *Registrar
	reader          Reader
	output          chan<- *event.ProcessorEvent
	active          time.Time
	lineBuffer      *bytes.Buffer
	multilineBuffer *bytes.Buffer
	firstLine       []byte
	readFunc        func() ([]byte, error)
}

// startHarvester starts a new harvester with the given offset
func (f *File) startHarvester(ctx context.Context, c *Config, register *Registrar, state State, offset int64) (err error) {
	// Set state to "not" finished to indicate that a harvester is running
	state.Finished = false
	state.Offset = offset

	// Create harvester with state
	h, err := NewHarvester(c, register, state, f.output)

	if err != nil {
		return err
	}
	h.ctx, h.cancel = context.WithCancel(ctx)

	err = h.Setup()
	if err != nil {
		return fmt.Errorf("error setting up harvester: %s", err)
	}

	// Update state before staring harvester
	// This makes sure the states is set to Finished: false
	// This is synchronous state update as part of the scan
	h.register.SendStateUpdate(h.state)
	h.active = time.Now()
	h.lineBuffer = bytes.NewBuffer(make([]byte, 0, h.config.MaxLength))
	h.multilineBuffer = bytes.NewBuffer(make([]byte, 0, h.config.MaxLength))

	if h.config.Multiline != nil {
		h.readFunc = h.readMultiLine
	} else {
		h.readFunc = h.readOneLine
	}

	go h.stateUpdatePeriodically()
	go h.Run()
	go h.activeCheck()

	return err
}

// harvestExistingFile continues harvesting a file with a known state if needed
func (f *File) harvestExistingFile(ctx context.Context, c *Config, register *Registrar, newState State, oldState State) {
	//log.Info("Update existing file for harvesting: %s, offset: %v", newState.Source, oldState.Offset)

	// No harvester is running for the file, start a new harvester
	// It is important here that only the size is checked and not modification time, as modification time could be incorrect on windows
	// https://blogs.technet.microsoft.com/asiasupp/2010/12/14/file-date-modified-property-are-not-updating-while-modifying-a-file-without-closing-it/

	if oldState.Finished && newState.Fileinfo.Size() > oldState.Offset {
		// Resume harvesting of an old file we've stopped harvesting from
		// This could also be an issue with force_close_older that a new harvester is started after each scan but not needed?
		// One problem with comparing modTime is that it is in seconds, and scans can happen more then once a second
		log.Info("Resuming harvesting of file: %s, offset: %d, new size: %d", newState.Source, oldState.Offset, newState.Fileinfo.Size())
		err := f.startHarvester(ctx, c, register, newState, oldState.Offset)
		if err != nil {
			log.Error("Harvester could not be started on existing file: %s, Err: %s", newState.Source, err)
		}
		return
	}

	// File size was reduced -> truncated file
	if newState.Fileinfo.Size() < oldState.Offset {
		log.Info("Old file was truncated. Starting from the beginning: %s, old size: %d  new size: %d", newState.Source, oldState.Offset, newState.Fileinfo.Size())
		if oldState.Finished {
			err := f.startHarvester(ctx, c, register, newState, 0)
			if err != nil {
				log.Error("Harvester could not be started on truncated file: %s, Err: %s", newState.Source, err)
			}
			return
		}
		// just stop old harvester
		h := f.register.GetHarvester(oldState.Inode)
		if h != nil {
			h.Stop()
		}
		return
	}

	// Check if file was renamed
	if oldState.Source != "" && oldState.Source != newState.Source {
		// This does not start a new harvester as it is assume that the older harvester is still running
		// or no new lines were detected. It sends only an event status update to make sure the new name is persisted.
		log.Info("File rename was detected: %s -> %s, Current offset: %v", oldState.Source, newState.Source, oldState.Offset)

		oldState.Source = newState.Source
		f.register.SendStateUpdate(oldState)
	}

	if !oldState.Finished {
		// Nothing to do. Harvester is still running and file was not renamed
		log.V(1).Info("Harvester for file is still running: %s, inode %d", newState.Source, newState.Inode)
	} else {
		log.V(1).Info("File didn't change: %s, inode %d", newState.Source, newState.Inode)
	}
}

// NewHarvester creates a new harvester
func NewHarvester(c *Config, register *Registrar, state State, output chan<- *event.ProcessorEvent) (*Harvester, error) {
	h := &Harvester{
		config:   c,
		state:    state,
		register: register,
		output:   output,
	}

	// Add ttl if cleanInactive is set
	if h.config.CleanInactive > 0 {
		h.state.TTL = time.Duration(h.config.CleanInactive)
	}

	// Add outlet signal so harvester can also stop itself
	return h, nil
}

// Setup opens the file handler and creates the reader for the harvester
func (h *Harvester) Setup() error {
	err := h.openFile()
	if err != nil {
		return fmt.Errorf("Harvester setup failed. Unexpected file opening error: %s", err)
	}

	h.reader, err = h.newLogFileReader()
	if err != nil {
		if h.source != nil {
			h.source.Close()
		}
		return fmt.Errorf("Harvester setup failed. Unexpected encoding line reader error: %s", err)
	}

	return nil
}

// openFile opens a file and checks for the encoding. In case the encoding cannot be detected
// or the file cannot be opened because for example of failing read permissions, an error
// is returned and the harvester is closed. The file will be picked up again the next time
// the file system is scanned
func (h *Harvester) openFile() error {
	f, err := os.OpenFile(h.state.Source, os.O_RDONLY, os.FileMode(0))
	if err != nil {
		return fmt.Errorf("Failed opening %s: %s", h.state.Source, err)
	}

	// Makes sure file handler is also closed on errors
	err = h.validateFile(f)
	if err != nil {
		f.Close()
		return err
	}

	h.source = f
	return nil
}

func (h *Harvester) validateFile(f *os.File) error {
	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("Failed getting stats for file %s: %s", h.state.Source, err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("Tried to open non regular file: %q %s", info.Mode(), info.Name())
	}

	// Compares the stat of the opened file to the state given by the input. Abort if not match.
	if !os.SameFile(h.state.Fileinfo, info) {
		return errors.New("file info is not identical with opened file. Aborting harvesting and retrying file later again")
	}

	// get file offset. Only update offset if no error
	offset, err := h.initFileOffset(f)
	if err != nil {
		return err
	}

	log.V(1).Info("harvester Setting offset for file: %s inode %d. Offset: %d ", h.state.Source, h.state.Inode, offset)
	h.state.Offset = offset

	return nil
}

// initFileOffset set offset for file handler
func (h *Harvester) initFileOffset(file *os.File) (int64, error) {
	// continue from last known offset
	if h.state.Offset > 0 {
		return file.Seek(h.state.Offset, os.SEEK_SET)
	}

	var firstRun = false
	if v := h.ctx.Value("firstRun"); v != nil {
		firstRun = v.(bool)
	}
	if h.config.ReadFrom == "newest" && firstRun {
		return file.Seek(0, os.SEEK_END)
	}
	return file.Seek(0, os.SEEK_CUR)
}

func (h *Harvester) newLogFileReader() (Reader, error) {
	return NewLineReader(h.source, h.config.MaxLength)
}

func (h *Harvester) WriteToProcessor(message []byte) {
	e := event.GetEvent()
	e.Write(message)
	e.AppId = []byte(h.config.AppId)
	e.LogId = h.config.LogId
	e.Source = "file"
	// update fields
	for k, v := range h.config.Fields {
		e.Fields[k] = v
	}
	e.Fields["file"] = h.state.Source
	e.Destination = lancerroute.GetLancerByLogid(e.LogId)
	// time maybe overwrite by processor
	e.Time = time.Now()
	e.TimeRangeKey = strconv.FormatInt(e.Time.Unix()/100*100, 10)

	h.output <- e
}

func (h *Harvester) stateUpdatePeriodically() {
	interval := time.Tick(time.Second * 5)
	var offset int64
	for {
		select {
		case <-interval:
			if h.state.Offset > offset {
				offset = h.state.Offset
				h.register.SendStateUpdate(h.state)
				h.active = time.Now()
			}
		case <-h.ctx.Done():
			h.register.SendStateUpdate(h.state)
			return
		}
	}
}

func (h *Harvester) Stop() {
	h.state.Finished = true
	h.register.SendStateUpdate(h.state)
	h.cancel()
	log.Info("Harvester for File: %s Inode %d Existed", h.state.Source, h.state.Inode)
}

func (h *Harvester) activeCheck() {
	interval := time.Tick(time.Minute * 1)
	for {
		select {
		case <-interval:
			if time.Now().Sub(h.active) > time.Duration(h.config.HarvesterTTL) {
				log.Info("Harvester for file: %s, inode: %d is inactive longer than HarvesterTTL, Ended", h.state.Source, h.state.Inode)
				h.Stop()
				return
			}
		case <-h.ctx.Done():
			return
		}
	}
}

func (h *Harvester) readMultiLine() (b []byte, err error) {
	h.multilineBuffer.Reset()
	counter := 0
	if h.firstLine != nil {
		h.multilineBuffer.Write(h.firstLine)
	}

	ctx, _ := context.WithTimeout(h.ctx, time.Duration(h.config.Timeout))
	for {
		select {
		case <-ctx.Done():
			h.firstLine = nil
			return h.multilineBuffer.Bytes(), lineReadTimeout
		default:
		}

		message, err := h.readOneLine()

		if err != nil && err != io.EOF && err != lineReadTimeout {
			h.firstLine = nil
			if h.multilineBuffer.Len() == 0 {
				return message, nil
			}
			if len(message) > 0 {
				h.multilineBuffer.Write([]byte{'\n'})
				h.multilineBuffer.Write(message)
			}
			return h.multilineBuffer.Bytes(), err
		}

		if len(message) == 0 {
			continue
		}

		matched, err := regexp.Match(h.config.Multiline.Pattern, message)

		if matched {
			// old multiline ended
			if h.firstLine != nil {
				h.firstLine = message
				return h.multilineBuffer.Bytes(), nil
			}
			// pure new multiline
			if h.firstLine == nil {
				h.firstLine = message
				h.multilineBuffer.Write(message)
				continue
			}
		}

		if !matched {
			// multiline not begin
			if h.firstLine == nil {
				return message, nil
			}

			if h.firstLine != nil {
				h.multilineBuffer.Write([]byte{'\n'})
				h.multilineBuffer.Write(message)
				counter += 1
			}
		}

		if counter > h.config.Multiline.MaxLines || h.multilineBuffer.Len() > h.config.MaxLength {
			h.firstLine = nil
			return h.multilineBuffer.Bytes(), nil
		}
	}
}

func (h *Harvester) readOneLine() (b []byte, err error) {
	h.lineBuffer.Reset()
	ctx, _ := context.WithTimeout(h.ctx, time.Duration(h.config.Timeout))

	for {
		select {
		case <-ctx.Done():
			return h.lineBuffer.Bytes(), lineReadTimeout
		default:
		}

		message, advance, err := h.reader.Next()
		// update offset
		h.state.Offset += int64(advance)

		if err == nil && h.lineBuffer.Len() == 0 {
			return message, nil
		}

		h.lineBuffer.Write(message)

		if err == nil {
			return h.lineBuffer.Bytes(), nil
		}

		if err == io.EOF && advance == 0 {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		if err == io.EOF && advance > 0 && h.lineBuffer.Len() >= h.config.MaxLength {
			return h.lineBuffer.Bytes(), nil
		}

		if err != nil {
			return h.lineBuffer.Bytes(), err
		}
	}
}

//// Run start the harvester and reads files line by line and sends events to the defined output
//func (h *Harvester) Run() {
//	log.V(1).Info("Harvester started for file: %s, inode %d", h.state.Source, h.state.Inode)
//	h.register.RegisterHarvester(h)
//	defer h.register.UnRegisterHarvester(h)
//
//	var line = make([]byte, 0, h.config.MaxLength)
//	for {
//		select {
//		case <-h.ctx.Done():
//			return
//		default:
//		}
//		//TODO MaxLength check
//		message, advance, err := h.reader.Next()
//		// update offset
//		h.state.Offset += int64(advance)
//
//		if err == nil {
//			if len(line) == 0 {
//				h.WriteToProcessor(message)
//				continue
//			}
//			line = append(line, message...)
//			h.WriteToProcessor(line)
//			line = line[:0]
//			continue
//		}
//
//		if err == io.EOF && advance == 0 {
//			time.Sleep(time.Millisecond * 100)
//			continue
//		}
//
//		if err == io.EOF && advance > 0 {
//			line = append(line, message...)
//			if len(line) >= h.config.MaxLength {
//				h.WriteToProcessor(line)
//				line = line[:0]
//			}
//			continue
//		}
//
//		if err != nil {
//			log.Error("Harvester Read line error: %v; File: %v", err, h.state.Source)
//			h.Stop()
//			return
//		}
//	}
//}

// Run start the harvester and reads files line by line and sends events to the defined output
func (h *Harvester) Run() {
	log.V(1).Info("Harvester started for file: %s, inode %d", h.state.Source, h.state.Inode)
	h.register.RegisterHarvester(h)
	defer h.register.UnRegisterHarvester(h)

	for {
		select {
		case <-h.ctx.Done():
			return
		default:
		}

		message, err := h.readFunc()

		if err == lineReadTimeout {
			log.V(1).Info("lineReadTimeout when harvesting %s", h.state.Source)
		}

		if len(message) > 0 {
			h.WriteToProcessor(message)
		}

		if err != nil && err != lineReadTimeout && err != io.EOF {
			log.Error("Harvester Read line error: %v; File: %v", err, h.state.Source)
			h.Stop()
			return
		}
	}
}
