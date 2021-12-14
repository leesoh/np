package result

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

func (r *Result) Print() {
	for _, h := range r.Hosts {
		r.Logger.Debugf("working on host: %v", h)
		if h.Name != "" {
			fmt.Printf("%v (%v)\n", h.Name, h.IP)
		}
		for _, t := range h.TCPPorts {
			fmt.Printf("%v/%v %v\n", t.Number, "tcp", t.Name)
		}
		for _, u := range h.UDPPorts {
			fmt.Printf("%v/%v %v\n", u.Number, "udp", u.Name)
		}
		fmt.Println()
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
