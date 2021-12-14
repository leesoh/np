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
	for _, h := range r.Hosts {
		r.Logger.Debugf("working on host: %v", h)
		if r.allPortsClosed(h) {
			continue
		}
		if h.Name != "" {
			fmt.Printf("%v (%v)\n", h.Name, h.IP)
		} else {
			fmt.Printf("%v\n", h.IP)
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(writer, "PORT\tSERVICE\tPRODUCT\tVERSION")
		r.portPrinter(writer, "tcp", h.TCPPorts)
		r.portPrinter(writer, "udp", h.UDPPorts)
		writer.Flush()
		fmt.Println()
	}
}

func (r *Result) PrintHost(ip net.IP) {
	for _, h := range r.Hosts {
		if h.IP.Equal(ip) {
			if r.allPortsClosed(h) {
				continue
			}
			if h.Name != "" {
				fmt.Printf("%v (%v)\n", h.Name, h.IP)
			} else {
				fmt.Printf("%v\n", h.IP)
			}
			writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			fmt.Fprintln(writer, "PORT\tSERVICE\tPRODUCT\tVERSION")
			r.portPrinter(writer, "tcp", h.TCPPorts)
			r.portPrinter(writer, "udp", h.UDPPorts)
			writer.Flush()
			fmt.Println()
		}
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
	for _, h := range r.Hosts {
		if r.allPortsClosed(h) {
			continue
		}
		for _, v := range h.TCPPorts {
			if matched, _ := regexp.MatchString(service, v.Name); matched {
				r.Logger.Debugf("matched: %v", v.Name)
				r.PrintHost(h.IP)
			}
		}
	}
}

func (r *Result) PrintServices() {
	for _, h := range r.Hosts {
		for _, v := range h.TCPPorts {
			r.printIfValue(v.Name)
		}
	}
}

func (r *Result) PrintByPort(port int) {
	for _, h := range r.Hosts {
		if r.allPortsClosed(h) {
			continue
		}
		if _, hasPort := h.TCPPorts[port]; hasPort {
			r.Logger.Debugf("%v has port: %v", h.Name, port)
			r.PrintHost(h.IP)
		}
		if _, hasPort := h.UDPPorts[port]; hasPort {
			r.Logger.Debugf("%v has port: %v", h.Name, port)
			r.PrintHost(h.IP)
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
	for _, h := range r.Hosts {
		if r.allPortsClosed(h) {
			continue
		}
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

func (r *Result) allPortsClosed(h *Host) bool {
	if len(h.TCPPorts) == 0 && len(h.UDPPorts) == 0 {
		r.Logger.Debugf("no open ports on host: %v", h)
		return true
	}
	return false
}
