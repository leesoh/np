package runner

import (
	"flag"
	"strconv"
	"strings"

	"github.com/Masterminds/log-go"
	"github.com/Masterminds/log-go/impl/cli"
)

type Options struct {
	Exclude  []string
	JSON     bool
	Host     string
	Hosts    bool
	CIDR     string
	Path     string
	Port     []int
	Ports    bool
	Service  string
	Services bool
	Debug    bool
}

func ParseOptions() *Options {
	options := &Options{}
	logger := cli.NewStandard()
	logger.Level = log.ErrorLevel
	exclude := flag.String("exclude", "", "Exclude these hosts from output")
	flag.BoolVar(&options.JSON, "json", false, "Display JSON output")
	flag.StringVar(&options.Host, "host", "", "Show results for specified host")
	flag.BoolVar(&options.Hosts, "hosts", false, "Print alive hosts")
	flag.StringVar(&options.CIDR, "cidr", "0.0.0.0/0", "CIDR range to output")
	flag.StringVar(&options.Path, "path", ".", "Path to scan file")
	port := flag.String("port", "", "Display hosts with matching port(s)")
	flag.BoolVar(&options.Ports, "ports", false, "Print all ports")
	flag.StringVar(&options.Service, "service", "", "Display hosts with matching service name")
	flag.BoolVar(&options.Services, "services", false, "Print all services")
	flag.BoolVar(&options.Debug, "debug", false, "Display debug output")
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
	if *exclude != "" {
		el := strings.Split(*exclude, ",")
		options.Exclude = append(options.Exclude, el...)
	}
	return options
}
