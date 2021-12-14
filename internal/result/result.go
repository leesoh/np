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
	Number    int    `json:"port,omitempty"`
	Name      string `json:"name,omitempty"`
	Product   string `json:"product,omitempty"`
	Version   string `json:"version,omitempty"`
	ExtraInfo string `json:"extra_info,omitempty"`
}

func (p *Port) merge(np Port) {
	if p.Name == "" && np.Name != "" {
		p.Name = np.Name
	}
	if p.Product == "" && np.Product != "" {
		p.Product = np.Product
	}
	if p.ExtraInfo == "" && np.ExtraInfo != "" {
		p.ExtraInfo = np.ExtraInfo
	}
	if p.Version == "" && np.Version != "" {
		p.Version = np.Version
	}
}

func New(logger *cli.Logger) *Result {
	return &Result{Logger: logger}
}

// TODO: Refactor this mess
func (r *Result) AddHost(new *Host) {
	r.Logger.Debugf("processing host: %+v", new)
	for hostIndex, hh := range r.Hosts {
		// host is a match
		if hh.IP.Equal(new.IP) {
			// check all the new ports for a match
			for _, tt := range new.TCPPorts {
				// against all the old ports
				for _, ht := range hh.TCPPorts {
					// port is a match
					if ht.Number == tt.Number {
						r.Logger.Debugf("current: %v new: %v", ht, tt)
						ht.merge(tt)
						return
					}
				}
				// new port, add it
				r.Hosts[hostIndex].TCPPorts = append(r.Hosts[hostIndex].TCPPorts, tt)
			}
			return
		}
	}
	// we didn't find an existing host, add this one
	r.Hosts = append(r.Hosts, new)
}
