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
		if h.Name != "" {
			fmt.Printf("%v (%v)\n", h.Name, h.IP)
		}
		r.portPrinter(h.TCPPorts)
		r.portPrinter(h.UDPPorts)
		fmt.Println()
	}
}

func (r *Result) portPrinter(p map[int]*Port) {
	//	PORT      STATE    SERVICE
	//19/tcp    filtered chargen
	//22/tcp    open     ssh
	//25/tcp    filtered smtp
	//80/tcp    open     http
	//135/tcp   filtered msrpc
	//139/tcp   filtered netbios-ssn
	//445/tcp   filtered microsoft-ds
	//5631/tcp  filtered pcanywheredata
	//9929/tcp  open     nping-echo
	//31337/tcp open     Elite

	writer := tabwriter.NewWriter(os.Stdout, 0, 2, 2, '\t', 0)
	fmt.Fprintln(writer, "PORT\tSERVICE\tPRODUCT\tVERSION")
	for k, v := range p {
		line := fmt.Sprintf("%v\t%v\t%v\t%v", k, v.Name, v.Product, v.Version)
		fmt.Fprintln(writer, line)
		writer.Flush()
	}
	//for k, v := range p {
	//	tp := fmt.Sprintf("%v/%v %v", k, "tcp", v.Name)
	//	if v.Product != "" {
	//		tp = fmt.Sprintf("%v %v", tp, v.Product)
	//		if v.Version != "" {
	//			tp = fmt.Sprintf("%v %v", tp, v.Version)
	//			if v.ExtraInfo != "" {
	//				tp = fmt.Sprintf("%v (%v)", tp, v.ExtraInfo)
	//			}
	//		}
	//	}
	//	fmt.Println(tp)
	//}
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
