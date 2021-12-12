package result

import (
	"fmt"
	"net"
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
	Services []Service
}

type Service struct {
	Port     int
	Protocol string
	Name     string
}

func New(logger *cli.Logger) *Result {
	return &Result{Logger: logger}
}

func (r *Result) AddNmapHost(ns *scan.NmapScan) {
	h := &Host{
		Name: r.GetHostname(ns),
		IPs:  r.GetIPs(ns),
	}
	r.Hosts = append(r.Hosts, h)
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

func (r *Result) GetServices(ns *scan.NmapScan) []Service {
	var services []Service
	for _, p := range ns.Host.Ports {
		//TBD
	}
}

func (r *Result) Print() {
	for _, h := range r.Hosts {
		r.Logger.Debugf("working on host: %v", h)
		if h.Name != "" {
			fmt.Printf("%v (%v)\n", h.Name, r.PrintableIPList(h.IPs))
		}

		//} else {
		//	fmt.Println(h.IPs[0])
		//}
		//for _, s := range h.Services {
		//	fmt.Printf("\t%v/%v %v", s.Port, s.Protocol, s.Name)
		//}
	}
}

func (r *Result) PrintableIPList(ips []net.IP) string {
	var ipList strings.Builder
	for _, i := range ips {
		fmt.Fprintf(&ipList, "%v,", i)
	}
	return strings.TrimSuffix(ipList.String(), ",")
}
