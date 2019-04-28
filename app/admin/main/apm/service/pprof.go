package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/pprof"
	"go-common/library/log"
	xtime "go-common/library/time"
)

var (
	port   = "2333"
	cpuURL = "%s/x/admin/apm/pprof/svg?name=%s&uri=%s&hostname=%s"
	msg    = "<a href=\"%s\">点击查看详情</a>(注：生成时间有延迟，如未显示，稍后重试)"
)

// kind .
const (
	CPUPerformace   = 1 // CPU性能图
	CPUFrame        = 2 // CPU火焰图
	HeapPerformance = 3 // 内存性能图
	HeapFrame       = 4 // 内存火焰图
)

// Pprof ...
func (s *Service) Pprof(url, uri, svgName, hostName string, time int64, sType int8) (err error) {
	var (
		out    bytes.Buffer
		errOut bytes.Buffer
	)
	goPath := "go"
	if len(conf.Conf.Pprof.GoPath) > 0 {
		goPath = conf.Conf.Pprof.GoPath
	}
	f, err := exec.LookPath(goPath)
	if err != nil {
		log.Error("pprof go error(%v) goPath=(%v)", err, goPath)
		fmt.Printf("pprof=(%v)", err)
		return
	}
	cmd := exec.Command(f, "tool", "pprof", "--seconds="+strconv.FormatInt(time, 10), "--svg", "--output="+conf.Conf.Pprof.Dir+"/"+svgName+"_"+hostName+"_"+uri+".svg", url+"/debug/pprof/"+uri)
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	// 执行命令
	if sType == 1 { //串行
		if err = cmd.Run(); err != nil {
			log.Error("pprof Run stdout=(%s) stderr=(%s) error(%v)", out.String(), errOut.String(), err)
		}
	} else { //阻塞
		if err = cmd.Start(); err != nil {
			log.Error("pprof Start stdout=(%s) stderr=(%s) error(%v)", out.String(), errOut.String(), err)
		}
		if err = cmd.Wait(); err != nil {
			log.Error("s.Pprof cmd.Wait() error(%v)", err)
		}
	}
	return
}

// Torch ...
func (s *Service) Torch(c context.Context, url, uri, svgName, hostName string, time int64, sType int8) (err error) {
	goPath := "go-torch"
	// if len(conf.Conf.Pprof.GoPath) > 0 {
	// 	goPath = conf.Conf.Pprof.GoPath
	// }
	f, err := exec.LookPath(goPath)
	if err != nil {
		log.Error("go-torch error(%v) goPath=(%v)", err, goPath)
		fmt.Printf("go-torch=(%v)", err)
		return
	}
	cmd := exec.Command(f, "--url="+url, "--suffix=/debug/pprof/"+uri, "--seconds="+strconv.FormatInt(time, 10), "-f="+conf.Conf.Pprof.Dir+"/"+svgName+"_"+hostName+"_"+uri+"_flame.svg")
	var (
		out    bytes.Buffer
		cmdErr bytes.Buffer
	)
	cmd.Stdout = &out
	cmd.Stderr = &cmdErr
	// 执行命令
	if sType == 1 { //串行
		if err = cmd.Run(); err != nil {
			log.Error("go-torch Run stdout=(%s) stderr=(%s) error(%v)", out.String(), cmdErr.String(), err)
		}
	} else { //阻塞
		if err = cmd.Start(); err != nil {
			log.Error("go-torch Start stdout=(%s) stderr=(%s) error(%v)", out.String(), cmdErr.String(), err)
		}
		if err = cmd.Wait(); err != nil {
			log.Error("cmd.Wait() error(%v)", err)
		}
	}
	return
}

// ActiveWarning active
func (s *Service) ActiveWarning(c context.Context, text string) (err error) {
	var (
		ins     *pprof.Ins
		title   = "【%s】性能告警抓取通知"
		warn    = &pprof.Warning{}
		pws     = make([]*pprof.Warn, 0)
		times   = time.Now().Unix()
		curTime = xtime.Time(times)
		reg     = regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
	)
	if err = json.Unmarshal([]byte(text), warn); err != nil {
		log.Error("s.ActiveWarning json.Unmarshal data error(%v)", err)
		return
	}
	if warn.Tags.App == "" {
		log.Info("s.ActiveWarning skipped(AppName is empty)")
		return
	}
	title = fmt.Sprintf(title, warn.Tags.App)
	if ins, err = s.dao.Instances(c, warn.Tags.App); err != nil {
		log.Error("s.ActiveWaring get instances error(%v)", err)
		return
	}
	for _, instance := range ins.Instances {
		var (
			ip        string
			addrs     = instance.Addrs
			strs      = strings.Split(instance.Hostname, "-")
			hostName  = fmt.Sprintf("%s-%d", strs[len(strs)-1], times)
			pprofWarn = &pprof.Warn{}
		)
		if len(addrs) < 1 {
			log.Info("s.ActiveWarning not found adds")
			continue
		}
		ip = reg.FindString(addrs[0])
		host := fmt.Sprintf("http://%s:%s", ip, port)
		pprofWarn.IP = ip
		pprofWarn.AppID = warn.Tags.App
		pprofWarn.SvgName = hostName
		pprofWarn.Ctime = curTime
		pprofWarn.Mtime = curTime
		if err = s.Torch(c, host, "profile", warn.Tags.App, hostName, 30, 2); err == nil {
			pprofWarn.Kind = CPUPerformace
			pws = append(pws, packing(pprofWarn))
		}
		if err = s.Pprof(host, "profile", warn.Tags.App, hostName, 30, 2); err == nil {
			pprofWarn.Kind = CPUFrame
			pws = append(pws, packing(pprofWarn))
		}
		if err = s.Torch(c, host, "heap", warn.Tags.App, hostName, 30, 2); err == nil {
			pprofWarn.Kind = HeapPerformance
			pws = append(pws, packing(pprofWarn))
		}
		if err = s.Pprof(host, "heap", warn.Tags.App, hostName, 30, 2); err == nil {
			pprofWarn.Kind = HeapFrame
			pws = append(pws, packing(pprofWarn))
		}
	}
	if len(pws) == 0 {
		return
	}
	if err = s.AddPprofWarn(c, pws); err != nil {
		return
	}
	return s.SendWeChat(c, title, fmt.Sprintf(msg, s.c.Host.SVENCo), warn.Tags.App, strings.Join(s.c.WeChat.Users, ","))
}

// packing .
func packing(pw *pprof.Warn) (pprofWarn *pprof.Warn) {
	pprofWarn = &pprof.Warn{
		AppID:   pw.AppID,
		SvgName: pw.SvgName,
		IP:      pw.IP,
		Kind:    pw.Kind,
		Mtime:   pw.Mtime,
		Ctime:   pw.Ctime,
	}
	return
}

// AddPprofWarn .
func (s *Service) AddPprofWarn(c context.Context, pws []*pprof.Warn) (err error) {
	var (
		sql   = "INSERT INTO `pprof_warn`(`app_id`, `svg_name`, `ip`, `kind`, `ctime`, `mtime`) VALUES %s"
		key   = make([]string, 0)
		value = make([]interface{}, 0)
	)
	for _, pw := range pws {
		key = append(key, "(?,?,?,?,?,?)")
		value = append(value, pw.AppID, pw.SvgName, pw.IP, pw.Kind, pw.Ctime, pw.Mtime)
	}
	if err = s.DB.Exec(fmt.Sprintf(sql, strings.Join(key, ",")), value...).Error; err != nil {
		log.Error("s.AddPprofWarn error(%v)", err)
	}
	return
}

// PprofWarn .
func (s *Service) PprofWarn(c context.Context, req *pprof.Params) (pws []*pprof.Warn, err error) {
	var (
		query = s.DB.Where("1=1")
	)
	pws = make([]*pprof.Warn, 0)
	if req.AppID == "" && req.IP == "" && req.SvgName == "" && req.Kind == 0 && req.StartTime == 0 && req.EndTime == 0 {
		return
	}
	if req.AppID != "" {
		query = query.Where("app_id=?", req.AppID)
	}
	if req.SvgName != "" {
		query = query.Where("svg_name=?", req.SvgName)
	}
	if req.IP != "" {
		query = query.Where("ip=?", req.IP)
	}
	if req.Kind != 0 {
		query = query.Where("kind=?", req.Kind)
	}
	if req.StartTime != 0 && req.EndTime != 0 {
		query = query.Where("mtime between ? and ?", req.StartTime, req.EndTime)
	}
	if err = query.Order("mtime desc").Find(&pws).Error; err != nil {
		log.Error("s.PprofWarn query error(%v)", err)
	}
	s.setSvgURL(pws)
	return
}

// setSvgURL .
func (s *Service) setSvgURL(pws []*pprof.Warn) {
	for _, pw := range pws {
		switch {
		case pw.Kind == CPUPerformace:
			pw.URL = fmt.Sprintf(cpuURL, s.c.Host.SVENCo, pw.AppID, "profile", pw.SvgName)
		case pw.Kind == CPUFrame:
			pw.URL = fmt.Sprintf(cpuURL, s.c.Host.SVENCo, pw.AppID, "profile_flame", pw.SvgName)
		case pw.Kind == HeapPerformance:
			pw.URL = fmt.Sprintf(cpuURL, s.c.Host.SVENCo, pw.AppID, "heap", pw.SvgName)
		case pw.Kind == HeapFrame:
			pw.URL = fmt.Sprintf(cpuURL, s.c.Host.SVENCo, pw.AppID, "heap_flame", pw.SvgName)
		default:
		}
	}
}
