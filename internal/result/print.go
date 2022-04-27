package result

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

// Print prints all hosts and the details of their open ports
func (r *Result) Print() {
	for _, hh := range r.Hosts {
		r.Logger.Debugf("now printing: %v", hh.Name)
		r.hostPrinter(hh)
		r.portPrinter(hh)
	}
}

// PrintHost prints the full host details for the requested host
func (r *Result) PrintHost(h *Host) {
	// Iterate through all hosts to find a match for the requested host
	for _, hh := range r.Hosts {
		if hh.IP.Equal(h.IP) {
			r.hostPrinter(hh)
		} else if hh.Name != "" && hh.Name == h.Name {
			r.hostPrinter(hh)
		} else {
			continue
		}
		r.portPrinter(hh)
	}
}

func (r *Result) hostPrinter(h *Host) {
	if r.allPortsClosed(h) {
		r.Logger.Debugf("all ports closed: %v", h.IP)
		return
	}
	if h.Name != "" {
		fmt.Printf("%v (%v)\n", h.Name, h.IP)
	} else {
		fmt.Printf("%v\n", h.IP)
	}
}

func (r *Result) portPrinter(h *Host) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(writer, "PORT\tSERVICE\tPRODUCT\tVERSION")
	for k, v := range h.TCPPorts {
		if v.Name == "tcpwrapped" {
			continue
		}
		line := fmt.Sprintf("%v/%v\t%v\t%v\t%v", k, "tcp", v.Name, v.Product, v.Version)
		fmt.Fprintln(writer, line)
	}
	for k, v := range h.UDPPorts {
		if v.Name == "tcpwrapped" {
			continue
		}
		line := fmt.Sprintf("%v/%v\t%v\t%v\t%v", k, "udp", v.Name, v.Product, v.Version)
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
	fmt.Println()
}

func (r *Result) PrintAlive() {
	for _, hh := range r.Hosts {
		if r.allPortsClosed(hh) {
			r.Logger.Debugf("all ports closed: %v", hh.IP)
			continue
		}
		// Prefer hostname over IP
		if hh.Name != "" {
			fmt.Printf("%v (%v)\n", hh.Name, hh.IP)
		} else {
			fmt.Printf("%v\n", hh.IP)
		}
	}
}

// PrintByService prints a specific service
func (r *Result) PrintByService(service string) {
	var hosts []string
	for _, hh := range r.Hosts {
		if r.allPortsClosed(hh) {
			continue
		}
		for k, v := range hh.TCPPorts {
			if matched, _ := regexp.MatchString(service, v.Name); matched {
				r.Logger.Debugf("matched: %v", hh.GetName())
				s := r.formatService(hh.GetName(), k, v.Name)
				hosts = append(hosts, s)
			}
		}
		for k, v := range hh.UDPPorts {
			if matched, _ := regexp.MatchString(service, v.Name); matched {
				r.Logger.Debugf("matched: %v", hh.GetName())
				s := r.formatService(hh.GetName(), k, v.Name)
				hosts = append(hosts, s)
			}
		}
	}
	for _, hh := range hosts {
		fmt.Println(hh)
	}
}

// PrintService prints all services in host:port service format
func (r *Result) PrintServices() {
	var hosts []string
	for _, hh := range r.Hosts {
		for k, v := range hh.TCPPorts {
			r.Logger.Debugf("matched: %v", hh.GetName())
			s := r.formatService(hh.GetName(), k, v.Name)
			hosts = append(hosts, s)
		}
		for k, v := range hh.UDPPorts {
			r.Logger.Debugf("matched: %v", hh.GetName())
			s := r.formatService(hh.GetName(), k, v.Name)
			hosts = append(hosts, s)
		}
	}
	for _, hh := range hosts {
		fmt.Println(hh)
	}
}

func (r *Result) formatService(name string, port int, service string) string {
	return fmt.Sprintf("%v:%v %v", name, port, service)
}

func (r *Result) PrintByPort(port []int) {
	for _, hh := range r.Hosts {
		if r.allPortsClosed(hh) {
			continue
		}
		for _, pp := range port {
			if _, hasPort := hh.TCPPorts[pp]; hasPort {
				r.Logger.Debugf("%v has pp: %v", hh.Name, pp)
				fmt.Printf("%v:%v\n", hh.IP, pp)
			}
			if _, hasPort := hh.UDPPorts[pp]; hasPort {
				r.Logger.Debugf("%v has pp: %v", hh.Name, pp)
				fmt.Printf("%v:%v\n", hh.IP, pp)
			}
		}
	}
}

func (r *Result) PrintPortSummary() {
	p := make(map[int]struct{})
	for _, h := range r.Hosts {
		for k := range h.TCPPorts {
			p[k] = struct{}{}
		}
		for k := range h.UDPPorts {
			p[k] = struct{}{}
		}
	}
	sorted := sortIntMap(p)
	var ps strings.Builder
	for _, k := range sorted {
		s := fmt.Sprintf("%s,", fmt.Sprint(k))
		ps.WriteString(s)
	}
	fmt.Println(strings.TrimSuffix(ps.String(), ","))
}

func (r *Result) printIfValue(s string) {
	if s != "" {
		fmt.Println(s)
	}
}

func (r *Result) PrintJSON() {
	b, err := json.MarshalIndent(r.Hosts, "", "  ")
	if err != nil {
		r.Logger.Error(err)
	}
	fmt.Println(string(b))
}

func (r *Result) PrintableIPList(ips []net.IP) string {
	var ipList strings.Builder
	for _, ii := range ips {
		fmt.Fprintf(&ipList, "%v,", ii)
	}
	return strings.TrimSuffix(ipList.String(), ",")
}

func (r *Result) allPortsClosed(h *Host) bool {
	if len(h.TCPPorts) == 0 && len(h.UDPPorts) == 0 {
		r.Logger.Debugf("no open ports on host: %v", h)
		return true
	}
	return false
}
