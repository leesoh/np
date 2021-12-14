package result

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
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

func (r *Result) PrintHost(host string) {
	for _, h := range r.Hosts {
		if h.Name == host || h.IP.String() == host {
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
