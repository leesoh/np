package runner

import (
	"flag"
)

type Options struct {
	JSON     bool
	Host     string
	Hosts    bool
	Path     string
	Port     int
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
	flag.IntVar(&options.Port, "port", 0, "Display hosts with matching port")
	flag.BoolVar(&options.Ports, "ports", false, "Print all ports")
	flag.StringVar(&options.Service, "service", "", "Display hosts with matching service name")
	flag.BoolVar(&options.Services, "services", false, "Print all services")
	flag.BoolVar(&options.Verbose, "verbose", false, "Display verbose output")
	flag.Parse()
	return options
}
