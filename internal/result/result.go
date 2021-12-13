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
	IP       net.IP           `json:"ip"`
	Name     string           `json:"hostname"`
	TCPPorts []map[int]string `json:"tcp_ports,omitempty"`
	UDPPorts []map[int]string `json:"udp_ports,omitempty"`
}

func New(logger *cli.Logger) *Result {
	return &Result{Logger: logger}
}

func (r *Result) AddHost(h *Host) {
	r.addHost(h)
}

func (r *Result) addHost(h *Host) {
	for _, hh := range r.Hosts {
		// Names are the same, we should merge
		if hh.Name == h.Name {
			//TBD
		}
	}
	r.Hosts = append(r.Hosts, h)
}
