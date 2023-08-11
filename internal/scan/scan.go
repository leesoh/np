package scan

import (
	"net"
	"strconv"

	"github.com/Masterminds/log-go/impl/cli"
	"github.com/leesoh/np/internal/result"
	"github.com/leesoh/np/internal/scan/dnsx"
	"github.com/leesoh/np/internal/scan/naabu"
	"github.com/leesoh/np/internal/scan/nmap"
)

type Scan struct {
	Bytes  []byte
	Logger *cli.Logger
	Result *result.Result
}

func New(b []byte, logger *cli.Logger, r *result.Result) *Scan {
	return &Scan{
		Bytes:  b,
		Logger: logger,
		Result: r,
	}
}

func (s *Scan) IsNmap() bool {
	if _, err := nmap.Parse(s.Bytes); err != nil {
		s.Logger.Debug("not an Nmap scan")
		return false
	}
	return true
}

func (s *Scan) ParseNmap() {
	ns, err := nmap.Parse(s.Bytes)
	if err != nil {
		s.Logger.Errorf("error parsing Nmap scan: %v", err)
	}
	for _, hh := range ns.Hosts {
		h := &result.Host{
			Name:     s.getNmapHostname(hh),
			IP:       s.stringToIP(hh.Address.Addr),
			TCPPorts: s.getNmapPorts(hh, "tcp"),
			UDPPorts: s.getNmapPorts(hh, "udp"),
		}
		s.Result.AddHost(h)
		s.Logger.Debugf("added host: %v", h.IP)
	}
}

func (s *Scan) getNmapHostname(h nmap.Host) string {
	for _, hh := range h.Hostnames {
		if hh.Type == "user" {
			s.Logger.Debugf("found hostname: %v", hh.Name)
			return hh.Name
		}
	}
	return ""
}

func (s *Scan) stringToIP(ipString string) net.IP {
	ip := net.ParseIP(ipString)
	if ip != nil {
		s.Logger.Debugf("added IP: %v", ip)
		return ip
	}
	return nil
}

func (s *Scan) getNmapPorts(h nmap.Host, protocol string) map[int]*result.Port {
	ports := make(map[int]*result.Port)
	for _, pp := range h.Ports {
		if pp.State.State == "open" && pp.Protocol == protocol {
			number, err := strconv.Atoi(pp.PortID)
			if err != nil {
				s.Logger.Errorf("error casting port: %v", pp.PortID)
			}
			port := &result.Port{
				Name:      pp.Service.Name,
				Product:   pp.Service.Product,
				Version:   pp.Service.Version,
				ExtraInfo: pp.Service.ExtraInfo,
			}
			ports[number] = port
			s.Logger.Debugf("found port: %v/%v", number, protocol)
		}
	}
	return ports
}

// IsNaabuV1 checks whether the v2.1.11 or earlier template is used.
// I'm calling it IsNaabuV1 since it's the first type of output they used.
func (s *Scan) IsNaabuV1() bool {
	if _, err := naabu.ParseV1(s.Bytes); err != nil {
		s.Logger.Debugf("not a Naabu V1 scan: %v", err)
		return false
	}
	return true
}

// IsNaabuV2 checks whether the v2.1.11 or earlier template is used.
// I'm calling it IsNaabuV2 since it's the first type of output they used.
func (s *Scan) IsNaabuV2() bool {
	if _, err := naabu.ParseV2(s.Bytes); err != nil {
		s.Logger.Debugf("not a Naabu V2 scan: %v", err)
		return false
	}
	return true
}

func (s *Scan) ParseNaabuV1() {
	ns, err := naabu.ParseV1(s.Bytes)
	if err != nil {
		s.Logger.Errorf("error parsing Naabu scan: %v", err)
	}
	for _, hh := range ns.V1Hosts {
		h := &result.Host{
			Name:     hh.Name,
			IP:       s.stringToIP(hh.IPAddress),
			TCPPorts: s.intToPort(hh.Port),
			UDPPorts: map[int]*result.Port{}, // Naabu V1 format
		}
		s.Result.AddHost(h)
		s.Logger.Debugf("added host: %v", h.IP)
	}
}

func (s *Scan) ParseNaabuV2() {
	ns, err := naabu.ParseV2(s.Bytes)
	if err != nil {
		s.Logger.Errorf("error parsing Naabu scan: %v", err)
	}
	for _, hh := range ns.V2Hosts {
		h := &result.Host{
			Name: hh.Name,
			IP:   s.stringToIP(hh.IPAddress),
		}
		// https://github.com/projectdiscovery/naabu/blob/46dc6d250523f0047532d6009017863a6276040b/v2/pkg/protocol/protocol.go#L6
		switch hh.Port.Protocol {
		case 0:
			h.TCPPorts = s.intToPort(hh.Port.Port)
		case 1:
			h.UDPPorts = s.intToPort(hh.Port.Port)
		default:
			s.Logger.Debugf("unsupported protocol")
		}
		s.Result.AddHost(h)
		s.Logger.Debugf("added host: %v", h.IP)
	}
}

func (s *Scan) intToPort(portInt int) map[int]*result.Port {
	ports := make(map[int]*result.Port)
	ports[portInt] = &result.Port{}
	return ports
}

func (s *Scan) IsDNSx() bool {
	if _, err := dnsx.Parse(s.Bytes); err != nil {
		s.Logger.Debugf("not a DNSx scan: %v", err)
		return false
	}
	return true
}

func (s *Scan) ParseDNSx() {
	ds, err := dnsx.Parse(s.Bytes)
	if err != nil {
		s.Logger.Errorf("error parsing DNSx scan: %v", err)
	}
	for _, rr := range ds.Records {
		// We will add each IP:host mapping as a discrete host
		for i := range rr.IPAddresses {
			h := &result.Host{
				Name:     rr.Hostname,
				IP:       s.stringToIP(rr.IPAddresses[i]),
				TCPPorts: map[int]*result.Port{},
				UDPPorts: map[int]*result.Port{},
			}
			s.Result.AddHost(h)
			s.Logger.Debugf("added host: %v", h.Name)
		}
	}
}
