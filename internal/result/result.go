package result

import (
	"net"

	"github.com/Masterminds/log-go/impl/cli"
)

type Result struct {
	Logger *cli.Logger
	Hosts  []*Host
}

type Host struct {
	IP       net.IP `json:"ip"`
	Name     string `json:"hostname"`
	TCPPorts []Port `json:"tcp_ports,omitempty"`
	UDPPorts []Port `json:"udp_ports,omitempty"`
}

type Port struct {
	Number int    `json:"port,omitempty"`
	Name   string `json:"name,omitempty"`
}

func (p *Port) addNameIfBlank(name string) {
	if p.Name == "" && name != "" {
		p.Name = name
	}
}

func (h *Host) hasIP(ip net.IP) bool {
	return h.IP.Equal(ip)
}

func New(logger *cli.Logger) *Result {
	return &Result{Logger: logger}
}

//If hostname or IP is a match, merge. If hostname is blank, update If IP is
//missing, add it For each service, match on port and protocol If name is
//blank, add it.
func (r *Result) AddHost(new *Host) {
	for hostIndex, hh := range r.Hosts {
		if hh.IP.Equal(new.IP) {
			for _, tt := range new.TCPPorts {
				for _, ht := range hh.TCPPorts {
					if ht.Number == tt.Number {
						ht.addNameIfBlank(tt.Name)
						return
					}
				}
				r.Hosts[hostIndex].TCPPorts = append(r.Hosts[hostIndex].TCPPorts, tt)
			}
			return
		}
	}
	r.Hosts = append(r.Hosts, new)
	//if hh.Name == "" && new.Name != "" {
	//	hh.Name = new.Name
	//	r.Logger.Debugf("added hostname %v to ip %v", new.Name, hh.IP.String())
	//}
	//for newPort, newName := range new.TCPPorts {
	//	if r.hasPort(newPort, hh.TCPPorts) {
	//		if hh.TCPPorts[newPort] == nil && new.TCPPorts[newPort] != nil {
	//			hh.TCPPorts[newPort] = new.TCPPorts[newPort]
	//		}
	//	} else {
	//		hh.TCPPorts[newPort] = newName
	//	}
	//}
	//for newPort, newName := range new.UDPPorts {
	//	if r.hasPort(newPort, hh.UDPPorts) {
	//		if hh.UDPPorts[newPort] == nil && new.UDPPorts[newPort] != nil {
	//			hh.UDPPorts[newPort] = new.UDPPorts[newPort]
	//		}
	//	} else {
	//		hh.UDPPorts[newPort] = newName
	//	}
	//}
}

func (r *Result) hostExists(h *Host) bool {
	for _, hh := range r.Hosts {
		if hh.IP.Equal(h.IP) {
			r.Logger.Debugf("host exists: %v", h.IP)
			return true
		}
	}
	return false
}

//func (r *Result) mergeHost(new *Host) {
//	for k, existingHost := range r.Hosts {
//		if existingHost.IP.Equal(new.IP) {
//			// update blank name
//			if existingHost.Name == "" && new.Name != "" {
//				r.Hosts[k].Name = new.Name
//			}
//			//update missing port
//			for _, newTCPPort := range new.TCPPorts {
//				if existingHost.hasTCPPort(newTCPPort) {
//					existingHost.updatePort(newTCPPort)
//				}
//			}
//		}
//	}
//}
