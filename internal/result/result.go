package result

import (
	"bytes"
	"net"
	"sort"

	"github.com/Masterminds/log-go/impl/cli"
)

type Result struct {
	Logger  *cli.Logger
	Hosts   []*Host
	Exclude []string
}

func (r *Result) SortByIP() {
	sort.Slice(r.Hosts, func(p, q int) bool {
		return bytes.Compare(r.Hosts[p].IP, r.Hosts[q].IP) < 0
	})

}

func New(logger *cli.Logger, exclude []string) *Result {
	return &Result{
		Logger:  logger,
		Exclude: exclude,
	}
}

func (r *Result) AddHost(newHost *Host) {
	r.Logger.Debugf("processing host: %v", newHost.IP)
	// We search our list of hosts for a match. If one is found,
	// we add ports.
	for _, hh := range r.Hosts {
		// Don't process
		if hh.IP.Equal(newHost.IP) {
			r.Logger.Debugf("found existing host: %v", hh.IP)
			hh.updateName(newHost.Name)
			r.Logger.Debugf("adding ports to %v", hh.IP)
			hh.addTCPPorts(newHost.TCPPorts)
			hh.addUDPPorts(newHost.UDPPorts)
			return
		}
	}
	// We didn't find an existing host, add this one
	r.Logger.Debugf("found new host: %v", newHost.IP)
	// good to here
	r.Hosts = append(r.Hosts, newHost)
}

type Host struct {
	IP       net.IP        `json:"ip"`
	Name     string        `json:"hostname,omitempty"`
	TCPPorts map[int]*Port `json:"tcp_ports,omitempty"`
	UDPPorts map[int]*Port `json:"udp_ports,omitempty"`
}

func (h *Host) IsExcluded(el []string) bool {
	for _, ee := range el {
		if h.IP.String() == ee || h.Name == ee {
			return true
		}
	}
	return false
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
