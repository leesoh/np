package result

import (
	"net"

	"github.com/Masterminds/log-go/impl/cli"
)

type Result struct {
	Logger *cli.Logger
	Hosts  []*Host
}

func New(logger *cli.Logger) *Result {
	return &Result{
		Logger: logger,
	}
}

// TODO: Refactor this mess
func (r *Result) AddHost(new *Host) {
	r.Logger.Debugf("processing host: %+v", new)
	for _, hh := range r.Hosts {
		// host is a match
		if hh.IP.Equal(new.IP) {
			hh.addTCPPorts(new.TCPPorts)
			//hh.addUDPPorts(new.UDPPorts)
			return
		}
	}
	// we didn't find an existing host, add this one
	r.Hosts = append(r.Hosts, new)
}

type Host struct {
	IP       net.IP        `json:"ip"`
	Name     string        `json:"hostname"`
	TCPPorts map[int]*Port `json:"tcp_ports,omitempty"`
	UDPPorts map[int]*Port `json:"udp_ports,omitempty"`
}

func (h *Host) addTCPPorts(p map[int]*Port) {
	for k, v := range p {
		h.addOrUpdatePort(k, v)
	}
}

func (h *Host) addOrUpdatePort(k int, p *Port) {
	// port doesn't exist
	if _, hasPort := h.TCPPorts[k]; !hasPort {
		h.TCPPorts[k] = p
		return
	}
	h.mergePorts(k, p)
}

func (h *Host) mergePorts(k int, p *Port) {
	if h.TCPPorts[k].Name == "" && p.Name != "" {
		h.TCPPorts[k].Name = p.Name
	}
	if h.TCPPorts[k].Product == "" && p.Product != "" {
		h.TCPPorts[k].Product = p.Product
	}
	if h.TCPPorts[k].Version == "" && p.Version != "" {
		h.TCPPorts[k].Version = p.Version
	}
	if h.TCPPorts[k].ExtraInfo == "" && p.ExtraInfo != "" {
		h.TCPPorts[k].ExtraInfo = p.ExtraInfo
	}
}

type Port struct {
	Name      string `json:"name,omitempty"`
	Product   string `json:"product,omitempty"`
	Version   string `json:"version,omitempty"`
	ExtraInfo string `json:"extra_info,omitempty"`
}
