package diskqueue

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	// max memroy use equal to BucketByte * (MemBucket + DynamicMemBucket)
	_defaultBucketByte       = _blockByte * 2 * 1024 * 16 // 16MB
	_defaultMemBucket        = 1
	_defaultDynamicMemBucket = 3
	_filePrefix              = "disk_queue_"
)

// node status
const (
	_inmem int8 = iota
	_indisk
	_freed
)

var _globalID int64

// ErrQueueFull .
var ErrQueueFull = errors.New("error queue is full, can't create new membucket")

func nextNodeID() int64 {
	return atomic.AddInt64(&_globalID, 1)
}

// DiskQueue disk queue
type DiskQueue interface {
	Push(p []byte) error
	Pop() ([]byte, error)
	Close() error
}

// Option Ringbuffer option
type Option func(opt *option)

// SetBucketByte bucketbyte
func SetBucketByte(n int) Option {
	return func(opt *option) {
		opt.bucketByte = (int32(n) / _blockByte) * _blockByte
	}
}

// SetMemBucket set the number of mem bucket
func SetMemBucket(n int) Option {
	return func(opt *option) {
		opt.memBucket = int32(n)
	}
}

// SetDynamicMemBucket set the number of dynamic mem bucket
func SetDynamicMemBucket(n int) Option {
	return func(opt *option) {
		opt.dynamicMemBucket = int32(n)
	}
}

// SetMaxBucket set the number of max bucket 0 represent unlimit
func SetMaxBucket(n int) Option {
	return func(opt *option) {
		opt.maxBucket = int32(n)
	}
}

type option struct {
	bucketByte       int32
	memBucket        int32
	maxBucket        int32
	dynamicMemBucket int32
	fpath            string
}

func (o option) validate() error {
	if o.bucketByte <= 0 {
		return fmt.Errorf("bucket byte must > 0")
	}
	if o.memBucket <= 0 {
		return fmt.Errorf("mem bucket must > 0")
	}
	if o.dynamicMemBucket <= 0 {
		return fmt.Errorf("dynamic mem bucket must > 0")
	}
	return nil
}

var _defaultOpt = option{
	bucketByte:       _defaultBucketByte,
	memBucket:        _defaultMemBucket,
	dynamicMemBucket: _defaultDynamicMemBucket,
}

// New Ringbuffer
func New(fpath string, options ...Option) (DiskQueue, error) {
	info, err := os.Stat(fpath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("stat %s error: %s", fpath, err)
		}
		if err = os.MkdirAll(fpath, 0755); err != nil {
			return nil, fmt.Errorf("fpath %s not exists try create directry error: %s", fpath, err)
		}
	} else if !info.IsDir() {
		return nil, fmt.Errorf("fpath: %s already exists and not a directory", fpath)
	}
	// TODO: check permission
	opt := _defaultOpt
	opt.fpath = fpath
	for _, fn := range options {
		fn(&opt)
	}
	if err = opt.validate(); err != nil {
		return nil, err
	}
	b := &base{
		opt: opt,
	}
	if opt.maxBucket == 0 {
		return &queue{base: b}, b.init()
	}
	return nil, nil
}

type node struct {
	id       int64
	mx       sync.Mutex
	flushing bool
	bucket   *memBucket
	next     *node
	fpath    string
	fbucket  *fileBucket
	kind     int8
}

func (n *node) setFlushing(flushing bool) {
	n.mx.Lock()
	n.flushing = flushing
	n.mx.Unlock()
}

func (n *node) pop() ([]byte, error) {
	n.mx.Lock()
	defer n.mx.Unlock()
	if n.bucket != nil {
		return n.bucket.pop()
	}
	var err error
	if n.fbucket == nil {
		if n.fbucket, err = newFileBucket(n.fpath); err != nil {
			return nil, err
		}
	}
	return n.fbucket.pop()
}

type base struct {
	opt       option
	head      *node
	tail      *node
	pool      *memBucketPool
	length    int32
	memBucket int32
}

func (b *base) init() error {
	b.pool = newMemBucketPool(b.opt.bucketByte)
	if loaded, err := b.loadFromFile(); err != nil || loaded {
		return err
	}
	current := &node{
		id:     nextNodeID(),
		bucket: b.pool.new(),
	}
	b.head = current
	b.tail = current
	return nil
}

func (b *base) loadFromFile() (bool, error) {
	infos, err := ioutil.ReadDir(b.opt.fpath)
	if err != nil {
		return false, fmt.Errorf("readdir %s error: %s", b.opt.fpath, err)
	}
	var files []string
	for _, info := range infos {
		if info.IsDir() || !strings.HasPrefix(info.Name(), _filePrefix) {
			continue
		}
		files = append(files, path.Join(b.opt.fpath, info.Name()))
	}
	if len(files) == 0 {
		return false, nil
	}
	nodeID := func(name string) int64 {
		id, err := strconv.ParseInt(path.Base(name)[len(_filePrefix):], 10, 64)
		if err != nil {
			panic(fmt.Errorf("invalid file name: %s error: %s", name, err))
		}
		return id
	}
	sort.Slice(files, func(i int, j int) bool {
		return nodeID(files[i]) < nodeID(files[j])
	})
	_globalID = nodeID(files[len(files)-1])
	current := &node{
		id:    nodeID(files[0]),
		fpath: files[0],
		kind:  _indisk,
	}
	b.head = current
	for _, file := range files[1:] {
		next := &node{
			id:    nodeID(file),
			fpath: file,
			kind:  _indisk,
		}
		current.next = next
		current = next
	}
	b.memBucket = 1
	next := &node{
		id:     nextNodeID(),
		bucket: b.pool.new(),
	}
	current.next = next
	current = next
	b.tail = current
	return true, nil
}

type queue struct {
	*base
	mx     sync.Mutex
	closed bool
	lastID int64
	wg     sync.WaitGroup
}

func (q *queue) Push(p []byte) (err error) {
	if len(p) >= int(q.opt.bucketByte) {
		return fmt.Errorf("data too large")
	}
	if q.closed {
		return fmt.Errorf("queue already closed")
	}
	for {
		err = q.tail.bucket.push(p)
		if err == nil {
			atomic.AddInt32(&q.length, 1)
			return
		}
		if err == errBucketFull {
			if err = q.moveTail(); err != nil {
				return err
			}
			continue
		}
		return
	}
}

func (q *queue) moveTail() error {
	bucket := atomic.LoadInt32(&q.memBucket)
	if bucket >= q.opt.memBucket+q.opt.dynamicMemBucket {
		return fmt.Errorf("can't assign memory bucket any more")
	}
	if bucket >= q.opt.maxBucket {
		q.notifyStore()
	}
	// take tail snapshot
	p := q.tail
	// lock queue
	q.mx.Lock()
	defer q.mx.Unlock()
	// tail alreay changed
	if p != q.tail {
		return nil
	}
	atomic.AddInt32(&q.memBucket, 1)
	n := &node{
		id:     nextNodeID(),
		bucket: q.pool.new(),
		kind:   _inmem,
	}
	// move to new tail
	q.tail.next = n
	q.tail = n
	return nil
}

func (q *queue) notifyStore() {
	n := q.head
	for n.next != nil {
		read := q.head
		if n.id > q.lastID && n.kind != _indisk && n != read {
			q.lastID = n.id
			go q.storeNode(n)
			return
		}
		n = n.next
	}
}

func (q *queue) Pop() (data []byte, err error) {
	defer func() {
		if err != nil {
			atomic.AddInt32(&q.length, -1)
		}
	}()
	if q.closed {
		return nil, fmt.Errorf("queue already closed")
	}
	data, err = q.head.pop()
	if err != nil {
		if err == io.EOF {
			if err = q.moveHead(); err != nil {
				return nil, err
			}
			return q.head.pop()
		}
		return nil, err
	}
	return data, nil
}

func (q *queue) moveHead() error {
	tail := q.tail
	if q.head == tail {
		return io.EOF
	}
	// move head to next
	q.mx.Lock()
	head := q.head
	q.head = q.head.next
	q.mx.Unlock()
	// reset head to new read node
	q.freeNode(head)
	return nil
}

func (q *queue) freeNode(n *node) {
	n.mx.Lock()
	defer n.mx.Unlock()
	if n.flushing {
		n.kind = _freed
		return
	}
	if n.bucket != nil {
		q.freeBucket(n.bucket)
		n.bucket = nil
	}
	if n.fbucket != nil {
		n.fbucket.close()
	}
	if n.fpath != "" {
		if err := os.Remove(n.fpath); err != nil {
			//fmt.Fprintf(os.Stderr, "[ERROR] diskqueue: remove file %s error: %s", n.fpath, err)
		}
	}
}

func (q *queue) storeNode(n *node) (err error) {
	fpath := storePath(q.opt.fpath, n)
	q.wg.Add(1)
	defer q.wg.Done()
	n.setFlushing(true)
	// if node already free return direct
	if n.bucket == nil {
		return
	}
	// if node be freed just release membucket
	if n.kind == _freed {
		q.freeBucket(n.bucket)
		return
	}
	// store bucket to disk
	if err = store(fpath, n); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] diskqueue: store node error: %s", err)
	}
	n.fpath = fpath
	n.setFlushing(false)
	if n.kind == _freed {
		q.freeBucket(n.bucket)
		n.bucket = nil
		if err := os.Remove(fpath); err != nil {
			//fmt.Fprintf(os.Stderr, "[ERROR] diskqueue: remove file %s error: %s", n.fpath, err)
		}
		return
	}
	n.kind = _indisk
	q.mx.Lock()
	if q.head != n {
		q.freeBucket(n.bucket)
		n.bucket = nil
	}
	q.mx.Unlock()
	return
}

func (q *queue) freeBucket(bucket *memBucket) {
	q.pool.free(bucket)
	atomic.AddInt32(&q.memBucket, -1)
}

func (q *queue) Close() error {
	// set closed
	q.closed = true
	// wait all store goroutines finish
	q.wg.Wait()
	var messages []string
	// store all leave node
	current := q.head
	for current != nil {
		if current.kind == _inmem && current.bucket != nil {
			fpath := storePath(q.opt.fpath, current)
			if err := store(fpath, current); err != nil {
				messages = append(messages, err.Error())
			}
		}
		current = current.next
	}
	if len(messages) == 0 {
		return nil
	}
	return fmt.Errorf("close queue error: %s", strings.Join(messages, "; "))
}

func store(fpath string, n *node) (err error) {
	// ignore empty bucket
	if n.bucket.writeAt == n.bucket.readAt {
		return nil
	}
	var fp *os.File
	fp, err = os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file %s error: %s", fpath, err)
	}
	_, err = n.bucket.dump(fp)
	if err != nil {
		return fmt.Errorf("dump data to file %s error: %s", fpath, err)
	}
	return
}

func storePath(base string, n *node) string {
	return path.Join(base, _filePrefix+strconv.FormatInt(n.id, 10))
}
