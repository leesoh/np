package runner

import (
	"flag"
)

type Options struct {
	JSON    bool
	Verbose bool
	Path    string
}

func ParseOptions() *Options {
	options := &Options{}
	flag.BoolVar(&options.JSON, "json", false, "Display JSON output")
	flag.BoolVar(&options.Verbose, "verbose", false, "Display verbose output")
	flag.StringVar(&options.Path, "path", "", "Path to scan file")
	flag.Parse()
	return options
}
