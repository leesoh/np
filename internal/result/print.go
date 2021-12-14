package result

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
)

func (r *Result) Print() {
	for _, hh := range r.Hosts {
		r.Logger.Debugf("working on host: %v", hh)
		if r.allPortsClosed(hh) {
			continue
		}
		if hh.Name != "" {
			fmt.Printf("%v (%v)\n", hh.Name, hh.IP)
		} else {
			fmt.Printf("%v\n", hh.IP)
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(writer, "PORT\tSERVICE\tPRODUCT\tVERSION")
		r.portPrinter(writer, "tcp", hh.TCPPorts)
		r.portPrinter(writer, "udp", hh.UDPPorts)
		writer.Flush()
		fmt.Println()
	}
}

func (r *Result) PrintHost(h *Host) {
	for _, hh := range r.Hosts {
		if hh.IP.Equal(h.IP) {
			if r.allPortsClosed(hh) {
				continue
			}
			if hh.Name != "" {
				fmt.Printf("%v (%v)\n", hh.Name, hh.IP)
			} else {
				fmt.Printf("%v\n", hh.IP)
			}
		} else if hh.Name == h.Name {
			if r.allPortsClosed(hh) {
				continue
			}
			fmt.Printf("%v (%v)\n", hh.Name, hh.IP)
		} else {
			continue
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(writer, "PORT\tSERVICE\tPRODUCT\tVERSION")
		r.portPrinter(writer, "tcp", hh.TCPPorts)
		r.portPrinter(writer, "udp", hh.UDPPorts)
		writer.Flush()
		fmt.Println()
	}
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
	hm := make(map[*Host]struct{})
	for _, hh := range r.Hosts {
		if r.allPortsClosed(hh) {
			continue
		}
		for _, v := range hh.TCPPorts {
			if matched, _ := regexp.MatchString(service, v.Name); matched {
				r.Logger.Debugf("matched: %v", v.Name)
				hm[hh] = struct{}{}
			}
		}
		for _, v := range hh.UDPPorts {
			if matched, _ := regexp.MatchString(service, v.Name); matched {
				r.Logger.Debugf("matched: %v", v.Name)
				hm[hh] = struct{}{}
			}
		}
	}
	for k := range hm {
		r.PrintHost(k)
	}
}

func (r *Result) PrintServices() {
	for _, hh := range r.Hosts {
		for _, v := range hh.TCPPorts {
			r.printIfValue(v.Name)
		}
	}
}

func (r *Result) PrintByPort(port int) {
	for _, hh := range r.Hosts {
		if r.allPortsClosed(hh) {
			continue
		}
		if _, hasPort := hh.TCPPorts[port]; hasPort {
			r.Logger.Debugf("%v has port: %v", hh.Name, port)
			r.PrintHost(hh)
		}
		if _, hasPort := hh.UDPPorts[port]; hasPort {
			r.Logger.Debugf("%v has port: %v", hh.Name, port)
			r.PrintHost(hh)
		}
	}
}

func (r *Result) PrintPorts() {
	p := make(map[int]struct{})
	for _, h := range r.Hosts {
		for k, _ := range h.TCPPorts {
			p[k] = struct{}{}
		}
		for k, _ := range h.UDPPorts {
			p[k] = struct{}{}
		}
	}
	var keys []int
	for k := range p {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	var ps strings.Builder
	for _, k := range keys {
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
func (r *Result) portPrinter(writer *tabwriter.Writer, protocol string, p map[int]*Port) {
	for k, v := range p {
		line := fmt.Sprintf("%v/%v\t%v\t%v\t%v", k, protocol, v.Name, v.Product, v.Version)
		fmt.Fprintln(writer, line)
	}
}

func (r *Result) PrintJSON() {
	for _, hh := range r.Hosts {
		if r.allPortsClosed(hh) {
			continue
		}
		r.Logger.Debugf("working on host: %v", hh)
		b, err := json.Marshal(hh)
		if err != nil {
			r.Logger.Error(err)
		}
		fmt.Println(string(b))
	}
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
