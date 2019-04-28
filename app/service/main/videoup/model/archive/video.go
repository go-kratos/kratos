package archive

import (
	"sync"
	xtime "time"

	"go-common/library/log"
	"go-common/library/time"
)

// const Video Status
const (
	// video xcode and dispatch state.
	VideoUploadInfo      = int8(0)
	VideoXcodeSDFail     = int8(1)
	VideoXcodeSDFinish   = int8(2)
	VideoXcodeHDFail     = int8(3)
	VideoXcodeHDFinish   = int8(4)
	VideoDispatchRunning = int8(5)
	VideoDispatchFinish  = int8(6)
	// video status.
	VideoStatusOpen         = int16(0)
	VideoStatusAccess       = int16(10000)
	VideoStatusWait         = int16(-1)
	VideoStatusRecicle      = int16(-2)
	VideoStatusLock         = int16(-4)
	VideoStatusXcodeFail    = int16(-16)
	VideoStatusSubmit       = int16(-30)
	VideoStatusUploadSubmit = int16(-50)
	VideoStatusDelete       = int16(-100)
	// xcode fail
	XcodeFailZero = 0
)

// Video is archive_video model.
type Video struct {
	ID           int64      `json:"-"`
	Aid          int64      `json:"aid"`
	Title        string     `json:"title"`
	Desc         string     `json:"desc"`
	Filename     string     `json:"filename"`
	SrcType      string     `json:"src_type"`
	Cid          int64      `json:"cid"`
	Sid          int64      `json:"-"`
	Duration     int64      `json:"duration"`
	Filesize     int64      `json:"-"`
	Resolutions  string     `json:"-"`
	Index        int        `json:"index"`
	Playurl      string     `json:"-"`
	Status       int16      `json:"status"`
	FailCode     int8       `json:"fail_code"`
	XcodeState   int8       `json:"xcode_state"`
	Attribute    int32      `json:"-"`
	RejectReason string     `json:"reject_reason"`
	CTime        time.Time  `json:"ctime"`
	MTime        time.Time  `json:"-"`
	Dimension    *Dimension `json:"dimension"`
}

// Dimension Archive video dimension
type Dimension struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
	Rotate int64 `json:"rotate"`
}

// SimpleVideo for Archive History
type SimpleVideo struct {
	Cid    int64     `json:"cid"`
	Index  int       `json:"part_id"`
	Title  string    `json:"part_name"`
	Status int16     `json:"status"`
	MTime  time.Time `json:"dm_modified"`
}

// VideoFn for Archive Video table filename check
type VideoFn struct {
	Cid      int64     `json:"cid"`
	Filename string    `json:"filename"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

type param struct {
	undoneCount int // 待完成数总和
	failCount   int // 错误数总和
	timeout     *xtime.Timer
}

/*VideosEditor 处理超多分p的稿件编辑
* 0.对该稿件的编辑操作加锁
* 1.将所有视频分组分事务更新
* 2.全部更新成功后发送成功信号量
* 3.更新错误的分组多次尝试，超时或超过次数后失败
* 4.收到更新成功信号量进行回调-（1.同步信息给视频云 2.解锁该稿件编辑）
* 5.收到更新失败信号量进行回调- (1.记录错误日志信息 2.推送错误消息 3.解锁该稿件编辑)
 */
type VideosEditor struct {
	sync.Mutex

	failTH             int
	closeFlag          bool
	wg                 sync.WaitGroup
	params             map[int64]*param
	cbSuccess          map[int64]func()
	cbFail             map[int64]func()
	sigSuccess         chan int64
	sigFail            chan int64
	chanRetry          chan func() (int64, int, int, error)
	closechan, cbclose chan struct{}
}

// NewEditor new VideosEditor
func NewEditor(failTH int) *VideosEditor {
	editor := &VideosEditor{
		failTH:     failTH,
		wg:         sync.WaitGroup{},
		params:     make(map[int64]*param),
		cbSuccess:  make(map[int64]func()),
		cbFail:     make(map[int64]func()),
		sigSuccess: make(chan int64, 10),
		sigFail:    make(chan int64, 10),
		chanRetry:  make(chan func() (int64, int, int, error), 100),
		closechan:  make(chan struct{}, 1),
		cbclose:    make(chan struct{}),
	}
	editor.wg.Add(1)
	go editor.consumerRetry(&editor.wg)
	editor.wg.Add(1)
	go editor.consumercb(&editor.wg)

	return editor
}

// Close 等待所有消息消费完才退出
func (m *VideosEditor) Close() {
	m.closeFlag = true
	m.closechan <- struct{}{}
	m.wg.Wait()
}

// Add add to editor
func (m *VideosEditor) Add(aid int64, cbSuccess, cbFail func(), timeout xtime.Duration, retrys ...func() (int64, int, int, error)) {
	if m.closeFlag {
		log.Warn("VideosEditor closed")
		return
	}
	log.Info("VideosEditor Add(%d) len(%d)", aid, len(retrys))

	timer := xtime.AfterFunc(timeout, func() { m.notifyTimeout(aid) })
	m.params[aid] = &param{
		undoneCount: len(retrys),
		failCount:   0,
		timeout:     timer,
	}
	m.cbSuccess[aid] = cbSuccess
	m.cbFail[aid] = cbFail
	for _, ry := range retrys {
		m.addRetry(ry, 0)
	}
}

// NotifySuccess notify success
func (m *VideosEditor) notifySuccess(aid int64) {
	m.Lock()
	defer m.Unlock()
	param, ok := m.params[aid]
	if ok {
		param.undoneCount--
		if param.undoneCount <= 0 {
			log.Info("notifySuccess(%d) undoneCount(%d)", aid, param.undoneCount)
			m.sigSuccess <- aid
		}
	}
}

// NotifyFail 返回值表示达到触发阈值
func (m *VideosEditor) notifyFail(aid int64) (retry bool) {
	m.Lock()
	defer m.Unlock()
	param, ok := m.params[aid]
	if ok {
		retry = true
		param.failCount++
		if param.failCount >= m.failTH {
			log.Warn("notifyFail(%d) failCount(%d)", aid, param.failCount)
			retry = false
			if _, ok := m.cbFail[aid]; ok {
				m.sigFail <- aid
			}
		}
	}
	return
}

func (m *VideosEditor) notifyTimeout(aid int64) {
	log.Warn("notifyTimeout(%d)", aid)
	m.Lock()
	defer m.Unlock()
	_, ok := m.params[aid]
	if ok {
		if _, ok := m.cbFail[aid]; ok {
			m.sigFail <- aid
		}
	}
}

func (m *VideosEditor) consumercb(g *sync.WaitGroup) {
	defer g.Done()

	for {
		select {
		case aid, ok := <-m.sigSuccess:
			if !ok {
				log.Info("consumercb close")
				return
			}
			if f, ok := m.cbSuccess[aid]; ok {
				f()
				m.release(aid)
			}
		case aid, ok := <-m.sigFail:
			if !ok {
				log.Info("consumercb close")
				return
			}
			if f, ok := m.cbFail[aid]; ok {
				f()
				m.release(aid)
			}
		case <-m.closechan:
			if len(m.cbFail) == 0 && len(m.cbSuccess) == 0 {
				m.cbclose <- struct{}{}
				log.Info("consumercb closechan")
				return
			}
			xtime.Sleep(50 * xtime.Millisecond)
			m.closechan <- struct{}{}
		}
	}
}

func (m *VideosEditor) consumerRetry(g *sync.WaitGroup) {
	defer g.Done()

	for {
		log.Info("consumerRetry")
		select {
		case <-m.cbclose:
			log.Info("consumerRetry closed")
			return
		case f, ok := <-m.chanRetry:
			if !ok {
				log.Info("m.chanRetry closed")
				break
			}
			id, head, tail, err := f()
			log.Info("consumerRetry(%d) head(%d) tail(%d) err(%v) ", id, head, tail, err)
			if err != nil {
				if m.notifyFail(id) {
					go func() {
						xtime.Sleep(3 * xtime.Second)
						m.addRetry(f, 3)
					}()
				}
			} else {
				m.notifySuccess(id)
			}
		}
	}
}

func (m *VideosEditor) addRetry(f func() (int64, int, int, error), asynctime int) {
	if asynctime == 0 {
		m.chanRetry <- f
		return
	}
	go func() {
		xtime.Sleep(3 * xtime.Second)
		m.chanRetry <- f
	}()
}

func (m *VideosEditor) release(aid int64) {
	if p, ok := m.params[aid]; ok {
		p.timeout.Stop()
	}
	delete(m.cbFail, aid)
	delete(m.cbSuccess, aid)
}
