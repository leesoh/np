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
	r.Logger.Debugf("processing host: %v", new.IP)
	// We search our list of hosts for a match. If one is found,
	// we add ports.
	for _, hh := range r.Hosts {
		if hh.IP.Equal(new.IP) {
			hh.updateName(new.Name)
			r.Logger.Debugf("adding ports to %v", hh.IP)
			hh.addTCPPorts(new.TCPPorts)
			hh.addUDPPorts(new.UDPPorts)
			return
		}
	}
	// We didn't find an existing host, add this one
	r.Logger.Debugf("found new host: %v", new.IP)
	// good to here
	r.Hosts = append(r.Hosts, new)
}

type Host struct {
	IP       net.IP        `json:"ip"`
	Name     string        `json:"hostname,omitempty"`
	TCPPorts map[int]*Port `json:"tcp_ports,omitempty"`
	UDPPorts map[int]*Port `json:"udp_ports,omitempty"`
}

func (h *Host) updateName(hostname string) {
	// Don't overwrite
	if h.Name == "" && hostname != "" {
		h.Name = hostname
	}
}

func (h *Host) addTCPPorts(p map[int]*Port) {
	for k, v := range p {
		h.addOrUpdatePort(k, v, "tcp")
	}
}

func (h *Host) addUDPPorts(p map[int]*Port) {
	for k, v := range p {
		h.addOrUpdatePort(k, v, "udp")
	}
}

func (h *Host) addOrUpdatePort(k int, p *Port, proto string) {
	if proto == "tcp" {
		if _, hasPort := h.TCPPorts[k]; !hasPort {
			h.TCPPorts[k] = p
			return
		}
		h.TCPPorts[k].Update(p)
	}
	if proto == "udp" {
		if _, hasPort := h.UDPPorts[k]; !hasPort {
			h.UDPPorts[k] = p
			return
		}
		h.UDPPorts[k].Update(p)
	}
}

func (h *Host) allPortsClosed() bool {
	return len(h.TCPPorts) == 0 && len(h.UDPPorts) == 0
}

func (h *Host) GetName() string {
	if h.Name != "" {
		return h.Name
	} else {
		return h.IP.String()
	}
}

type Port struct {
	Name      string `json:"name,omitempty"`
	Product   string `json:"product,omitempty"`
	Version   string `json:"version,omitempty"`
	ExtraInfo string `json:"extra_info,omitempty"`
}

func (p *Port) Update(portData *Port) {
	if portData.Name == "" && p.Name != "" {
		portData.Name = p.Name
	}
	if portData.Product == "" && p.Product != "" {
		portData.Product = p.Product
	}
	if portData.Version == "" && p.Version != "" {
		portData.Version = p.Version
	}
	if portData.ExtraInfo == "" && p.ExtraInfo != "" {
		portData.ExtraInfo = p.ExtraInfo
	}
}
