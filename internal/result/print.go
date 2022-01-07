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
		r.PrintHost(hh)
	}
}

// PrintHost prints the full host details for a given host
func (r *Result) PrintHost(h *Host) {
	for _, hh := range r.Hosts {
		if hh.IP.Equal(h.IP) {
			if r.allPortsClosed(hh) {
				continue
			}
			if hh.Name != "" {
				fmt.Printf("%v (%v)\n", hh.Name, hh.IP)
				r.Logger.Debugf("printed host with non-blank name: %+v", hh)
			} else {
				fmt.Printf("%v\n", hh.IP)
				r.Logger.Debugf("printed host with blank name: %+v", hh)
			}
		} else if hh.Name != "" && hh.Name == h.Name {
			if r.allPortsClosed(hh) {
				continue
			}
			fmt.Printf("%v (%v)\n", hh.Name, hh.IP)
		} else {
			continue
		}
		r.portPrinter(hh)
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
	for _, h := range r.Hosts {
		if r.allPortsClosed(h) {
			continue
		}
		if h.Name != "" {
			fmt.Printf("%v (%v)\n", h.IP, h.Name)
		} else {
			fmt.Printf("%v\n", h.IP)
		}
	}
}

func (r *Result) PrintByService(service string) {
	var hosts []string
	for _, hh := range r.Hosts {
		if r.allPortsClosed(hh) {
			continue
		}
		for k, v := range hh.TCPPorts {
			if matched, _ := regexp.MatchString(service, v.Name); matched {
				r.Logger.Debugf("matched: %v", hh.GetName())
				s := fmt.Sprintf("%v:%v", hh.GetName(), k)
				hosts = append(hosts, s)
			}
		}
		for k, v := range hh.UDPPorts {
			if matched, _ := regexp.MatchString(service, v.Name); matched {
				r.Logger.Debugf("matched: %v", hh.GetName())
				s := fmt.Sprintf("%v:%v", hh.GetName(), k)
				hosts = append(hosts, s)
			}
		}
	}
	for _, host := range hosts {
		fmt.Println(host)
	}
}

func (r *Result) PrintServices() {
	s := make(map[string]struct{})
	for _, hh := range r.Hosts {
		for _, v := range hh.TCPPorts {
			if v.Name != "" {
				s[v.Name] = struct{}{}
			}
		}
	}
	sorted := sortStringMap(s)
	for _, ss := range sorted {
		fmt.Println(ss)
	}
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
		for k, _ := range h.TCPPorts {
			p[k] = struct{}{}
		}
		for k, _ := range h.UDPPorts {
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
