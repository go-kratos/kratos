package dao

import (
	"bytes"
	"context"
	"time"

	"go-common/library/log"
)

// Play push databus .
func (d *Dao) Play(c context.Context, plat, aid, cid, part, mid, level, ftime, stime, did, ip, agent, buvid, cookieSid, refer, typeID, subType, sid, epid, playMode, platform, device, mobiAapp, autoPlay, session string) {
	bf := d.bfp.Get().(*bytes.Buffer)
	bf.WriteString(plat)
	bf.Write(d.spliter)
	bf.WriteString(aid)
	bf.Write(d.spliter)
	bf.WriteString(cid)
	bf.Write(d.spliter)
	bf.WriteString(part)
	bf.Write(d.spliter)
	bf.WriteString(mid)
	bf.Write(d.spliter)
	bf.WriteString(level)
	bf.Write(d.spliter)
	bf.WriteString(ftime)
	bf.Write(d.spliter)
	bf.WriteString(stime)
	bf.Write(d.spliter)
	bf.WriteString(did)
	bf.Write(d.spliter)
	bf.WriteString(ip)
	bf.Write(d.spliter)
	bf.WriteString(agent)
	bf.Write(d.spliter)
	bf.WriteString(buvid)
	bf.Write(d.spliter)
	bf.WriteString(cookieSid)
	bf.Write(d.spliter)
	bf.WriteString(refer)
	bf.Write(d.spliter)
	bf.WriteString(typeID)
	bf.Write(d.spliter)
	bf.WriteString(subType)
	bf.Write(d.spliter)
	bf.WriteString(sid)
	bf.Write(d.spliter)
	bf.WriteString(epid)
	bf.Write(d.spliter)
	bf.WriteString(playMode)
	bf.Write(d.spliter)
	bf.WriteString(platform)
	bf.Write(d.spliter)
	bf.WriteString(device)
	bf.Write(d.spliter)
	bf.WriteString(mobiAapp)
	bf.Write(d.spliter)
	bf.WriteString(autoPlay)
	bf.Write(d.spliter)
	bf.WriteString(session)

	buf := make([]byte, len(bf.Bytes()))
	copy(buf, bf.Bytes())
	select {
	case d.msgs <- buf:
	default:
		log.Warn("d.Play() msgs is full !")
	}
	bf.Reset()
	d.bfp.Put(bf)
}

// pubproc send history to databus.
func (d *Dao) pubproc() {
	var (
		msg    []byte
		ms     [][]byte
		ticker = time.NewTicker(time.Second)
	)
	for {
		select {
		case msg = <-d.msgs:
			if len(msg) == 0 {
				continue
			}
			if d.spliter[0] != msg[0] {
				ms = append(ms, msg)
				if len(ms) < 100 {
					continue
				}
			}
		case <-ticker.C:
		}
		if len(ms) == 0 {
			continue
		}
		d.mergePub(ms)
		ms = make([][]byte, 0, 100)
	}
}

func (d *Dao) mergePub(ms [][]byte) {
	key := string(ms[0][:50])
	for j := 0; j < 3; j++ {
		if err := d.merge.Send(context.Background(), key, ms); err == nil {
			break
		}
		log.Error("d.merge.Send(%+v)", ms)
	}
}
