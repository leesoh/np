package runner

import (
	"flag"
)

type Options struct {
	JSON    bool
	Host    string
	Hosts   bool
	Path    string
	Verbose bool
}

func ParseOptions() *Options {
	options := &Options{}
	flag.BoolVar(&options.JSON, "json", false, "Display JSON output")
	flag.StringVar(&options.Path, "path", ".", "Path to scan file")
	flag.StringVar(&options.Host, "host", "", "Show results for specified host")
	flag.BoolVar(&options.Hosts, "hosts", false, "Print alive hosts")
	flag.BoolVar(&options.Verbose, "verbose", false, "Display verbose output")
	flag.Parse()
	return options
}
