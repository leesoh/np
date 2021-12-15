package runner

import (
	"flag"
	"strconv"
	"strings"
)

type Options struct {
	JSON     bool
	Host     string
	Hosts    bool
	Path     string
	Port     []int
	Ports    bool
	Service  string
	Services bool
	Verbose  bool
}

func ParseOptions() *Options {
	options := &Options{}
	flag.BoolVar(&options.JSON, "json", false, "Display JSON output")
	flag.StringVar(&options.Host, "host", "", "Show results for specified host")
	flag.BoolVar(&options.Hosts, "hosts", false, "Print alive hosts")
	flag.StringVar(&options.Path, "path", ".", "Path to scan file")
	port := flag.String("port", "", "Display hosts with matching port(s)")
	flag.BoolVar(&options.Ports, "ports", false, "Print all ports")
	flag.StringVar(&options.Service, "service", "", "Display hosts with matching service name")
	flag.BoolVar(&options.Services, "services", false, "Print all services")
	flag.BoolVar(&options.Verbose, "verbose", false, "Display verbose output")
	flag.Parse()
	if *port != "" {
		pl := strings.Split(*port, ",")
		for _, pp := range pl {
			pi, err := strconv.Atoi(pp)
			if err != nil {
				continue
			}
			options.Port = append(options.Port, pi)
		}
	}
	return options
}
