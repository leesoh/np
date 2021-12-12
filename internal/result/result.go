package result

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/Masterminds/log-go/impl/cli"
	"github.com/leesoh/np/internal/scan"
)

type Result struct {
	Logger *cli.Logger
	Hosts  []*Host
}

type Host struct {
	IPs      []net.IP
	Name     string
	Services []*Service
}

type Service struct {
	Port     int
	Protocol string
	Name     string
}

func New(logger *cli.Logger) *Result {
	return &Result{Logger: logger}
}

func (r *Result) AddNmapHosts(ns *scan.NmapScan) {
	for _, hh := range ns.Host.Hostnames.Hostname {
		h := &Host{
			Name:     r.GetHostname(ns),
			IPs:      r.GetIPs(ns),
			Services: r.GetServices(ns),
		}
		r.addHost(h)
	}
}

func (r *Result) GetHostname(ns *scan.NmapScan) string {
	for _, h := range ns.Host.Hostnames.Hostname {
		if h.Type == "user" {
			return h.Name
		}
	}
	return ""
}

func (r *Result) GetIPs(ns *scan.NmapScan) []net.IP {
	var ips []net.IP
	//jfor _, a := range ns.Host.Address {
	//j	ip := net.ParseIP(a)
	//j	if a != "" {
	//j		ips = append(ips, ip)
	//j	}
	//j}
	r.Logger.Debugf("checking IP: %v", ns.Host.Address.Addr)
	ip := net.ParseIP(ns.Host.Address.Addr)
	if ip != nil {
		ips = append(ips, ip)
		r.Logger.Debugf("added IP: %v", ip)
	}
	return ips
}

func (r *Result) GetServices(ns *scan.NmapScan) []*Service {
	var services []*Service
	for _, p := range ns.Host.Ports.Port {
		if p.State.State == "open" {
			port, _ := strconv.Atoi(p.Portid)
			s := &Service{
				Port:     port,
				Protocol: p.Protocol,
				Name:     p.Service.Name,
			}
			services = append(services, s)
			r.Logger.Debugf("found services: %v", s)
		}
	}
	return services
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

func (r *Result) Print() {
	for _, h := range r.Hosts {
		r.Logger.Debugf("working on host: %v", h)
		if h.Name != "" {
			fmt.Printf("%v (%v)\n", h.Name, r.PrintableIPList(h.IPs))
		}
		for _, s := range h.Services {
			fmt.Printf("\t%v/%v %v\n", s.Port, s.Protocol, s.Name)
		}
	}
}

func (r *Result) PrintJSON() {
	for _, h := range r.Hosts {
		r.Logger.Debugf("working on host: %v", h)
		b, err := json.Marshal(h)
		if err != nil {
			r.Logger.Error(err)
		}
		fmt.Println(string(b))
	}
}

func (r *Result) PrintableIPList(ips []net.IP) string {
	var ipList strings.Builder
	for _, i := range ips {
		fmt.Fprintf(&ipList, "%v,", i)
	}
	return strings.TrimSuffix(ipList.String(), ",")
}
