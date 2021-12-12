package result

import (
	"fmt"
	"net"

	"github.com/Masterminds/log-go/impl/cli"
)

type Result struct {
	Logger *cli.Logger
	Hosts  []Host
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

func New(logger *cli.Logger) {
	return &Result{Logger: logger}
}

func (r *Result) Print() {
	for _, h := range r.Hosts {
		if h.Name {
			fmt.Println(h.Name)
		} else {
			fmt.Println(h.IPs[0])
		}
		for _, s := range h.Services {
			fmt.Printf("\t%v/%v %v", s.Port, s.Protocol, s.Name)
		}
	}
}
