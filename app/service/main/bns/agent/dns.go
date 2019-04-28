package agent

import (
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/miekg/dns"

	"go-common/app/service/main/bns/conf"
	"go-common/app/service/main/bns/lib/resolvconf"
	"go-common/app/service/main/bns/lib/shuffle"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

const (
	maxRecurseRecords = 5
)

var grpcPrefixs = []string{"_grpclb._tcp.", "_grpc_config."}

var dnsProm = prom.New().WithTimer("go_bns_server", []string{"time"})

func wrapProm(handler func(dns.ResponseWriter, *dns.Msg)) func(dns.ResponseWriter, *dns.Msg) {
	return func(w dns.ResponseWriter, m *dns.Msg) {
		start := time.Now()
		handler(w, m)
		dt := int64(time.Since(start) / time.Millisecond)
		dnsProm.Timing("bns:dns_query", dt)
	}
}

// DNSServer is used to wrap an Agent and expose various
// service discovery endpoints using a DNS interface.
type DNSServer struct {
	*dns.Server
	cfg       *conf.DNSServer
	agent     *Agent
	domain    string
	recursors []string

	// disableCompression is the config.DisableCompression flag that can
	// be safely changed at runtime. It always contains a bool and is
	// initialized with the value from config.DisableCompression.
	disableCompression atomic.Value

	udpClient *dns.Client
	tcpClient *dns.Client
}

// NewDNSServer new dns server
func NewDNSServer(a *Agent, cfg *conf.DNSServer) (*DNSServer, error) {
	var recursors []string
	var confRecursors = cfg.Config.Recursors
	if len(confRecursors) == 0 {
		resolv, err := resolvconf.ParseResolvConf()
		if err != nil {
			log.Warn("read resolv.conf error: %s", err)
		} else {
			confRecursors = resolv
		}
	}
	for _, r := range confRecursors {
		ra, err := recursorAddr(r)
		if err != nil {
			return nil, fmt.Errorf("Invalid recursor address: %v", err)
		}
		recursors = append(recursors, ra)
	}

	log.Info("recursors %v", recursors)

	// Make sure domain is FQDN, make it case insensitive for ServeMux
	domain := dns.Fqdn(strings.ToLower(cfg.Config.Domain))

	srv := &DNSServer{
		agent:     a,
		domain:    domain,
		recursors: recursors,
		cfg:       cfg,

		udpClient: &dns.Client{Net: "udp", Timeout: time.Duration(cfg.Config.RecursorTimeout)},
		tcpClient: &dns.Client{Net: "tcp", Timeout: time.Duration(cfg.Config.RecursorTimeout)},
	}
	srv.disableCompression.Store(cfg.Config.DisableCompression)

	return srv, nil
}

// ListenAndServe listen and serve dns
func (s *DNSServer) ListenAndServe(network, addr string, notif func()) error {
	mux := dns.NewServeMux()
	mux.HandleFunc("arpa.", wrapProm(s.handlePtr))
	mux.HandleFunc(".", wrapProm(s.handleRecurse))
	mux.HandleFunc(s.domain, wrapProm(s.handleQuery))

	s.Server = &dns.Server{
		Addr:              addr,
		Net:               network,
		Handler:           mux,
		NotifyStartedFunc: notif,
	}
	if network == "udp" {
		s.UDPSize = 65535
	}
	return s.Server.ListenAndServe()
}

// recursorAddr is used to add a port to the recursor if omitted.
func recursorAddr(recursor string) (string, error) {
	// Add the port if none
START:
	_, _, err := net.SplitHostPort(recursor)
	if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
		recursor = fmt.Sprintf("%s:%d", recursor, 53)
		goto START
	}
	if err != nil {
		return "", err
	}

	// Get the address
	addr, err := net.ResolveTCPAddr("tcp", recursor)
	if err != nil {
		return "", err
	}

	// Return string
	return addr.String(), nil
}

// handlePtr is used to handle "reverse" DNS queries
func (s *DNSServer) handlePtr(resp dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]
	defer func(s time.Time) {
		log.V(5).Info("dns: request for %v (%v) from client %s (%s)",
			q, time.Since(s), resp.RemoteAddr().String(),
			resp.RemoteAddr().Network())
	}(time.Now())

	// Setup the message response
	m := new(dns.Msg)
	m.SetReply(req)
	m.Compress = !s.disableCompression.Load().(bool)
	m.Authoritative = true
	m.RecursionAvailable = (len(s.recursors) > 0)

	// Only add the SOA if requested
	if req.Question[0].Qtype == dns.TypeSOA {
		s.addSOA(m)
	}

	// Get the QName without the domain suffix
	qName := strings.ToLower(dns.Fqdn(req.Question[0].Name))

	// FIXME: should return multiple nameservers?
	log.V(5).Info("dns: we said handled ptr with %v", qName)

	// nothing found locally, recurse
	if len(m.Answer) == 0 {
		s.handleRecurse(resp, req)
		return
	}

	// Enable EDNS if enabled
	if edns := req.IsEdns0(); edns != nil {
		m.SetEdns0(edns.UDPSize(), false)
	}

	// Write out the complete response
	if err := resp.WriteMsg(m); err != nil {
		log.Warn("dns: failed to respond: %v", err)
	}
}

// handleQuery is used to handle DNS queries in the configured domain
func (s *DNSServer) handleQuery(resp dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]
	defer func(s time.Time) {
		log.V(5).Info("dns: request for %v (%v) from client %s (%s)",
			q, time.Since(s), resp.RemoteAddr().String(),
			resp.RemoteAddr().Network())
	}(time.Now())

	// Switch to TCP if the client is
	network := "udp"
	if _, ok := resp.RemoteAddr().(*net.TCPAddr); ok {
		network = "tcp"
	}

	// Setup the message response
	m := new(dns.Msg)
	m.SetReply(req)
	m.Compress = !s.disableCompression.Load().(bool)
	m.Authoritative = true
	m.RecursionAvailable = (len(s.recursors) > 0)

	switch req.Question[0].Qtype {
	case dns.TypeSOA:
		ns, glue := s.nameservers(req.IsEdns0() != nil)
		m.Answer = append(m.Answer, s.soa())
		m.Ns = append(m.Ns, ns...)
		m.Extra = append(m.Extra, glue...)
		m.SetRcode(req, dns.RcodeSuccess)

	case dns.TypeNS:
		ns, glue := s.nameservers(req.IsEdns0() != nil)
		m.Answer = ns
		m.Extra = glue
		m.SetRcode(req, dns.RcodeSuccess)

	case dns.TypeAXFR:
		m.SetRcode(req, dns.RcodeNotImplemented)

	default:
		s.dispatch(network, req, m)
	}

	// Handle EDNS
	if edns := req.IsEdns0(); edns != nil {
		m.SetEdns0(edns.UDPSize(), false)
	}

	// Write out the complete response
	if err := resp.WriteMsg(m); err != nil {
		log.Warn("dns: failed to respond: %v", err)
	}
}

func (s *DNSServer) soa() *dns.SOA {
	return &dns.SOA{
		Hdr: dns.RR_Header{
			Name:   s.domain,
			Rrtype: dns.TypeSOA,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		Ns:     "ns." + s.domain,
		Serial: uint32(time.Now().Unix()),

		// todo(fs): make these configurable
		Mbox:    "hostmaster." + s.domain,
		Refresh: 3600,
		Retry:   600,
		Expire:  86400,
		Minttl:  30,
	}
}

// addSOA is used to add an SOA record to a message for the given domain
func (s *DNSServer) addSOA(msg *dns.Msg) {
	msg.Ns = append(msg.Ns, s.soa())
}

// formatNodeRecord takes an Easyns Agent node and returns an A, AAAA, or CNAME record
func (s *DNSServer) formatNodeRecord(addr, qName string, qType uint16, ttl time.Duration, edns bool) (records []dns.RR) {
	// Parse the IP
	ip := net.ParseIP(addr)
	var ipv4 net.IP
	if ip != nil {
		ipv4 = ip.To4()
	}
	switch {
	case ipv4 != nil && (qType == dns.TypeANY || qType == dns.TypeA):
		return []dns.RR{&dns.A{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(ttl / time.Second),
			},
			A: ip,
		}}

	case ip != nil && ipv4 == nil && (qType == dns.TypeANY || qType == dns.TypeAAAA):
		return []dns.RR{&dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    uint32(ttl / time.Second),
			},
			AAAA: ip,
		}}

	case ip == nil && (qType == dns.TypeANY || qType == dns.TypeCNAME ||
		qType == dns.TypeA || qType == dns.TypeAAAA):
		// Get the CNAME
		cnRec := &dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    uint32(ttl / time.Second),
			},
			Target: dns.Fqdn(addr),
		}
		records = append(records, cnRec)

		// Recurse
		more := s.resolveCNAME(cnRec.Target)
		extra := 0
	MORE_REC:
		for _, rr := range more {
			switch rr.Header().Rrtype {
			case dns.TypeCNAME, dns.TypeA, dns.TypeAAAA:
				records = append(records, rr)
				extra++
				if extra == maxRecurseRecords && !edns {
					break MORE_REC
				}
			}
		}
	}
	return records
}

// nameservers returns the names and ip addresses of up to three random servers
// in the current cluster which serve as authoritative name servers for zone.
func (s *DNSServer) nameservers(edns bool) (ns []dns.RR, extra []dns.RR) {
	// TODO: get list of bns dns nameservers
	// Then, construct them into NS RR.
	// We just hardcode here right now...
	name := "bns"
	addr := s.agent.cfg.DNS.Addr
	fqdn := name + "." + s.domain
	fqdn = dns.Fqdn(strings.ToLower(fqdn))

	// NS record
	nsrr := &dns.NS{
		Hdr: dns.RR_Header{
			Name:   s.domain,
			Rrtype: dns.TypeNS,
			Class:  dns.ClassINET,
			Ttl:    uint32(time.Duration(s.agent.cfg.DNS.Config.TTL) / time.Second),
		},
		Ns: fqdn,
	}
	ns = append(ns, nsrr)
	// A or AAAA glue record
	glue := s.formatNodeRecord(addr, fqdn, dns.TypeANY, time.Duration(s.agent.cfg.DNS.Config.TTL), edns)
	extra = append(extra, glue...)

	return
}

func trimDomainSuffix(name string, domain string) (reversed string) {
	reversed = strings.TrimSuffix(name, "."+domain)
	return strings.Trim(reversed, ".")
}

// Answers answers
type Answers []dns.RR

// Len the number of answer
func (as Answers) Len() int {
	return len(as)
}

// Swap order
func (as Answers) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

// dispatch is used to parse a request and invoke the correct handler
func (s *DNSServer) dispatch(network string, req, resp *dns.Msg) {
	var answers Answers
	// Get the QName
	qName := strings.ToLower(dns.Fqdn(req.Question[0].Name))
	name := trimDomainSuffix(qName, s.agent.cfg.DNS.Config.Domain)
	for _, prefix := range grpcPrefixs {
		name = strings.TrimPrefix(name, prefix)
	}
	inss, err := s.agent.Query(name)
	if err != nil {
		log.Error("dns: query %s failed to resolve from bns server, err: %s", name, err)
		goto INVALID
	}

	if len(inss) == 0 {
		log.Error("dns: QName %s has no upstreams found!", qName)
		goto INVALID
	}

	for _, ins := range inss {
		answers = append(answers, &dns.A{
			Hdr: dns.RR_Header{
				Name:   qName,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(time.Duration(s.agent.cfg.DNS.Config.TTL) / time.Second),
			},
			A: ins.IPAddr,
		})
		log.V(5).Info("dns: QName resolved ipAddress: %s - %s", qName, ins.IPAddr)
	}
	shuffle.Shuffle(answers)
	resp.Answer = []dns.RR(answers)
	return
INVALID:
	s.addSOA(resp)
	resp.SetRcode(req, dns.RcodeNameError)
}

// handleRecurse is used to handle recursive DNS queries
func (s *DNSServer) handleRecurse(resp dns.ResponseWriter, req *dns.Msg) {
	q := req.Question[0]
	network := "udp"
	client := s.udpClient

	defer func(s time.Time) {
		log.V(5).Info("dns: request for %v (%s) (%v) from client %s (%s)",
			q, network, time.Since(s), resp.RemoteAddr().String(),
			resp.RemoteAddr().Network())
	}(time.Now())

	// Switch to TCP if the client is
	if _, ok := resp.RemoteAddr().(*net.TCPAddr); ok {
		network = "tcp"
		client = s.tcpClient
	}

	for _, recursor := range s.recursors {
		r, rtt, err := client.Exchange(req, recursor)
		if err == nil || err == dns.ErrTruncated {
			// Compress the response; we don't know if the incoming
			// response was compressed or not, so by not compressing
			// we might generate an invalid packet on the way out.
			r.Compress = !s.disableCompression.Load().(bool)

			// Forward the response
			log.V(5).Info("dns: recurse RTT for %v (%v)", q, rtt)
			if err = resp.WriteMsg(r); err != nil {
				log.Warn("dns: failed to respond: %v", err)
			}
			return
		}
		log.Error("dns: recurse failed: %v", err)
	}

	// If all resolvers fail, return a SERVFAIL message
	log.Error("dns: all resolvers failed for %v from client %s (%s)",
		q, resp.RemoteAddr().String(), resp.RemoteAddr().Network())
	m := &dns.Msg{}
	m.SetReply(req)
	m.Compress = !s.disableCompression.Load().(bool)
	m.RecursionAvailable = true
	m.SetRcode(req, dns.RcodeServerFailure)
	if edns := req.IsEdns0(); edns != nil {
		m.SetEdns0(edns.UDPSize(), false)
	}
	resp.WriteMsg(m)
}

// resolveCNAME is used to recursively resolve CNAME records
func (s *DNSServer) resolveCNAME(name string) []dns.RR {
	// If the CNAME record points to a Easyns Name address, resolve it internally
	// Convert query to lowercase because DNS is case insensitive; d.domain is
	// already converted
	if strings.HasSuffix(strings.ToLower(name), "."+s.domain) {
		req := &dns.Msg{}
		resp := &dns.Msg{}

		req.SetQuestion(name, dns.TypeANY)
		s.dispatch("udp", req, resp)

		return resp.Answer
	}

	// Do nothing if we don't have a recursor
	if len(s.recursors) == 0 {
		return nil
	}

	// Ask for any A records
	m := new(dns.Msg)
	m.SetQuestion(name, dns.TypeA)

	// Make a DNS lookup request
	c := &dns.Client{Net: "udp", Timeout: time.Duration(s.agent.cfg.DNS.Config.RecursorTimeout)}
	var r *dns.Msg
	var rtt time.Duration
	var err error
	for _, recursor := range s.recursors {
		r, rtt, err = c.Exchange(m, recursor)
		if err == nil {
			log.V(5).Info("dns: cname recurse RTT for %v (%v)", name, rtt)
			return r.Answer
		}
		log.Error("dns: cname recurse failed for %v: %v", name, err)
	}
	log.Error("dns: all resolvers failed for %v", name)
	return nil
}
