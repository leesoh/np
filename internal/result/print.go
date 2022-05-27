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
	r.SortByIP()
	for _, hh := range r.Hosts {
		// We have to do it here since some scans will add by IP.
		// So if we exclude example.com, but then massscan imports
		// the IP for example.com, we want to wait until the last minute
		// To ignore
		if hh.IsExcluded(r.Exclude) {
			continue
		}
		r.hostPrinter(hh)
		r.portPrinter(hh)
	}
}

// PrintHost prints the full host details for the requested host (-host)
func (r *Result) PrintHost(h *Host) {
	// Iterate through all hosts to find a match for the requested host
	for _, hh := range r.Hosts {
		if hh.IsExcluded(r.Exclude) {
			continue
		}
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
		return
	}
	if h.Name != "" {
		fmt.Printf("%v (%v)\n", h.Name, h.IP)
	} else {
		fmt.Printf("%v\n", h.IP)
	}
}

// portPrinter prints a nice port table
func (r *Result) portPrinter(h *Host) {
	// We don't print ports if we don't have ports
	if h.allPortsClosed() {
		return
	}
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(writer, "PORT\tSERVICE\tPRODUCT\tVERSION")
	// tcpwrapped are boring, do not print
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

// PrintAlive prints all hosts with at least one open port (-hosts)
func (r *Result) PrintAlive() {
	r.SortByIP()
	for _, hh := range r.Hosts {
		if hh.IsExcluded(r.Exclude) {
			continue
		}
		if r.allPortsClosed(hh) {
			r.Logger.Debugf("all ports closed: %v", hh.IP)
			continue
		}
		// Prefer hostname over IP
		if hh.Name != "" {
			r.Logger.Debugf("%v", hh.Name)
			fmt.Printf("%v (%v)\n", hh.Name, hh.IP)
		} else {
			fmt.Printf("%v\n", hh.IP)
		}
	}
}

// PrintByService prints a specific service (-service)
func (r *Result) PrintByService(service string) {
	r.SortByIP()
	var hosts []string
	for _, hh := range r.Hosts {
		if hh.IsExcluded(r.Exclude) {
			continue
		}
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

// PrintService prints all services in host:port service format (-services)
func (r *Result) PrintServices() {
	r.SortByIP()
	var hosts []string
	for _, hh := range r.Hosts {
		if hh.IsExcluded(r.Exclude) {
			continue
		}
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

// PrintByPort prints all hosts with a specific port open (-port)
func (r *Result) PrintByPort(port []int) {
	r.SortByIP()
	for _, hh := range r.Hosts {
		if hh.IsExcluded(r.Exclude) {
			continue
		}
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

// PrintPortSummary prints a comma-delimited list of ports (-ports)
func (r *Result) PrintPortSummary() {
	p := make(map[int]struct{})
	for _, hh := range r.Hosts {
		if hh.IsExcluded(r.Exclude) {
			continue
		}
		for k := range hh.TCPPorts {
			p[k] = struct{}{}
		}
		for k := range hh.UDPPorts {
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
	r.SortByIP()
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
		r.Logger.Debugf("no open ports on host: %v", h.IP)
		return true
	}
	return false
}
