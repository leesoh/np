package scan

import (
	"net"
	"strconv"

	"github.com/Masterminds/log-go/impl/cli"
	"github.com/leesoh/np/internal/result"
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

func (s *Scan) IsNaabu() bool {
	if _, err := naabu.Parse(s.Bytes); err != nil {
		s.Logger.Debugf("not a Naabu scan: %v", err)
		return false
	}
	return true
}

func (s *Scan) ParseNaabu() {
	ns, err := naabu.Parse(s.Bytes)
	if err != nil {
		s.Logger.Errorf("error parsing Naabu scan: %v", err)
	}
	for _, hh := range ns.Hosts {
		h := &result.Host{
			Name:     hh.Name,
			IP:       s.stringToIP(hh.IPAddress),
			TCPPorts: s.intToPort(hh.Port),
			UDPPorts: map[int]*result.Port{}, // Naabu doesn't scan UDP ports
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
