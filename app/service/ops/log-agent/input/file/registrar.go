package file

import (
	"sync"
	"os"
	"fmt"
	"time"
	"io"
	"path/filepath"
	"encoding/json"
	"context"

	"go-common/library/log"
)

type Registrar struct {
	Channel              chan State
	registryFile         string // Path to the Registry File
	wg                   sync.WaitGroup
	states               *States // Map with all file paths inside and the corresponding state
	bufferedStateUpdates int
	flushInterval        time.Duration
	harvesters           map[uint64]*Harvester
	hLock                sync.RWMutex
	ctx                  context.Context
	cancel               context.CancelFunc
}

// New creates a new Registrar instance, updating the registry file on
// `file.State` updates. New fails if the file can not be opened or created.
func NewRegistry(ctx context.Context, registryFile string) (*Registrar, error) {
	r := &Registrar{
		registryFile:  registryFile,
		states:        NewStates(),
		Channel:       make(chan State, 100),
		wg:            sync.WaitGroup{},
		flushInterval: time.Second * 5,
		harvesters:    make(map[uint64]*Harvester),
	}
	r.ctx, r.cancel = context.WithCancel(ctx)
	err := r.Init()
	if err != nil {
		return nil, err
	}
	go r.Run()

	return r, err
}

func (r *Registrar) RegisterHarvester(h *Harvester) error {
	r.hLock.Lock()
	defer r.hLock.Unlock()

	if _, ok := r.harvesters[h.state.Inode]; !ok {
		r.harvesters[h.state.Inode] = h
		return nil
	}
	return fmt.Errorf("harvestor of inode %s Re registered", h.state.Inode)
}

func (r *Registrar) UnRegisterHarvester(h *Harvester) error {
	r.hLock.Lock()
	defer r.hLock.Unlock()

	if _, ok := r.harvesters[h.state.Inode]; ok {
		delete(r.harvesters, h.state.Inode)
		return nil
	}

	return fmt.Errorf("harvestor of inode %d not found", h.state.Inode)
}

func (r *Registrar) GetHarvester(i uint64) *Harvester {
	r.hLock.RLock()
	defer r.hLock.RUnlock()
	if h, ok := r.harvesters[i]; ok {
		return h
	}
	return nil
}

// Init sets up the Registrar and make sure the registry file is setup correctly
func (r *Registrar) Init() (err error) {
	// Create directory if it does not already exist.
	registryPath := filepath.Dir(r.registryFile)
	err = os.MkdirAll(registryPath, 0750)
	if err != nil {
		return fmt.Errorf("Failed to created registry file dir %s: %v", registryPath, err)
	}

	// Check if files exists
	fileInfo, err := os.Lstat(r.registryFile)
	if os.IsNotExist(err) {
		log.Info("No registry file found under: %s. Creating a new registry file.", r.registryFile)
		// No registry exists yet, write empty state to check if registry can be written
		return r.writeRegistry()
	}
	if err != nil {
		return err
	}

	// Check if regular file, no dir, no symlink
	if !fileInfo.Mode().IsRegular() {
		// Special error message for directory
		if fileInfo.IsDir() {
			return fmt.Errorf("Registry file path must be a file. %s is a directory.", r.registryFile)
		}
		return fmt.Errorf("Registry file path is not a regular file: %s", r.registryFile)
	}

	log.Info("Registry file set to: %s", r.registryFile)

	// load states
	if err = r.loadStates(); err != nil {
		return err
	}

	return nil
}

// writeRegistry writes the new json registry file to disk.
func (r *Registrar) writeRegistry() error {
	// First clean up states
	r.gcStates()

	// TODO lock for reading r.states.states
	tempfile, err := writeTmpFile(r.registryFile, r.states.states)
	if err != nil {
		return err
	}

	err = SafeFileRotate(r.registryFile, tempfile)
	if err != nil {
		return err
	}

	log.V(1).Info("Registry file %s updated. %d states written.", r.registryFile, len(r.states.states))
	return nil
}

// SafeFileRotate safely rotates an existing file under path and replaces it with the tempfile
func SafeFileRotate(path, tempfile string) error {
	parent := filepath.Dir(path)

	if e := os.Rename(tempfile, path); e != nil {
		return e
	}

	// best-effort fsync on parent directory. The fsync is required by some
	// filesystems, so to update the parents directory metadata to actually
	// contain the new file being rotated in.
	f, err := os.Open(parent)
	if err != nil {
		return nil // ignore error
	}
	defer f.Close()
	f.Sync()

	return nil
}

// loadStates fetches the previous reading state from the configure RegistryFile file
// The default file is `registry` in the data path.
func (r *Registrar) loadStates() error {
	f, err := os.Open(r.registryFile)
	if err != nil {
		return err
	}

	defer f.Close()

	log.Info("Loading registrar data from %s", r.registryFile)

	states, err := readStatesFrom(f)
	if err != nil {
		return err
	}

	states = r.preProcessStates(states)

	r.states.SetStates(states)
	log.V(1).Info("States Loaded from registrar%s : %+v", r.registryFile, len(states))

	return nil
}

func (r *Registrar) preProcessStates(states map[uint64]State) map[uint64]State {
	for key, state := range states {
		// set all states to finished
		state.Finished = true
		states[key] = state
	}
	return states
}

func readStatesFrom(in io.Reader) (map[uint64]State, error) {
	states := make(map[uint64]State)
	decoder := json.NewDecoder(in)
	if err := decoder.Decode(&states); err != nil {
		return nil, fmt.Errorf("Error decoding states: %s", err)
	}
	return states, nil
}

func writeTmpFile(baseName string, states map[uint64]State) (string, error) {
	tempfile := baseName + ".new"
	f, err := os.OpenFile(tempfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC|os.O_SYNC, 0640)
	if err != nil {
		log.Error("Failed to create tempfile (%s) for writing: %s", tempfile, err)
		return "", err
	}

	defer f.Close()

	encoder := json.NewEncoder(f)

	if err := encoder.Encode(states); err != nil {
		log.Error("Error when encoding the states: %s", err)
		return "", err
	}

	// Commit the changes to storage to avoid corrupt registry files
	if err = f.Sync(); err != nil {
		log.Error("Error when syncing new registry file contents: %s", err)
		return "", err
	}

	return tempfile, nil
}

// gcStates runs a registry Cleanup. The method check if more event in the
// registry can be gc'ed in the future. If no potential removable state is found,
// the gcEnabled flag is set to false, indicating the current registrar state being
// stable. New registry update events can re-enable state gc'ing.
func (r *Registrar) gcStates() {
	//if !r.gcRequired {
	//	return
	//}

	beforeCount := len(r.states.states)
	cleanedStates := r.states.Cleanup()

	log.V(1).Info(
		"Registrar states %s cleaned up. Before: %d, After: %d", r.registryFile,
		beforeCount, beforeCount-cleanedStates)
}

// FindPrevious lookups a registered state, that matching the new state.
// Returns a zero-state if no match is found.
func (r *Registrar) FindPrevious(newState State) State {
	return r.states.FindPrevious(newState)
}

func (r *Registrar) Run() {
	log.Info("Starting Registrar for: %s", r.registryFile)
	flushC := time.Tick(r.flushInterval)
	for {
		select {
		case <-flushC:
			r.flushRegistry()
		case state := <-r.Channel:
			r.processEventStates(state)
		case <-r.ctx.Done():
			r.flushRegistry()
			return
		}
	}
}

func (r *Registrar) flushRegistry() {
	if err := r.writeRegistry(); err != nil {
		log.Error("Writing of registry returned error: %v. Continuing...", err)
	}
}

// processEventStates gets the states from the events and writes them to the registrar state
func (r *Registrar) processEventStates(state State) {
	r.states.UpdateWithTs(state, time.Now())
}

func (r *Registrar) SendStateUpdate(state State) {
	r.Channel <- state
	//select {
	//case r.Channel <- state:
	//default:
	//	log.Warn("state update receiving chan full")
	//}
}
