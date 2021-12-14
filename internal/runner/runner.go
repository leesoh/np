package runner

import (
	"io/fs"
	"io/ioutil"
	"net"
	"path/filepath"

	"github.com/Masterminds/log-go"
	"github.com/Masterminds/log-go/impl/cli"
	"github.com/leesoh/np/internal/parse"
	"github.com/leesoh/np/internal/result"
	"github.com/leesoh/np/internal/scan"
)

type Runner struct {
	Files   []string
	Logger  *cli.Logger
	Options *Options
	Parser  *parse.Parser
}

func New(options *Options) *Runner {
	logger := cli.NewStandard()
	if options.Verbose {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.FatalLevel
	}
	parser := parse.NewParser(logger)
	runner := &Runner{
		Logger:  logger,
		Options: options,
		Parser:  parser,
	}
	return runner
}

func (r *Runner) Run() {
	r.GetScanFiles()
	res := result.New(r.Logger)
	for _, ff := range r.Files {
		b, err := ioutil.ReadFile(ff)
		if err != nil {
			r.Logger.Errorf("error reading file: %v", err)
		}
		s := scan.New(b, r.Logger, res)
		if s.IsNmap() {
			s.ParseNmap()
		}
	}
	if r.Options.Host != "" {
		ip := net.ParseIP(r.Options.Host)
		h := &result.Host{}
		r.Logger.Debugf("host option with ip: %+v", ip)
		if ip == nil {
			h.Name = r.Options.Host
		} else {
			h.IP = ip
		}
		r.Logger.Debugf("printing host: %v", h)
		res.PrintHost(h)
		return
	}
	if r.Options.Hosts {
		res.PrintAlive()
		return
	}
	if r.Options.Service != "" {
		res.PrintByService(r.Options.Service)
		return
	}
	if r.Options.Services {
		res.PrintServices()
		return
	}
	if r.Options.Port != 0 {
		res.PrintByPort(r.Options.Port)
		return
	}
	if r.Options.Ports {
		res.PrintPorts()
		return
	}
	if r.Options.JSON {
		res.PrintJSON()
	} else {
		res.Print()
	}
}

// GetScanFiles gathers scan files from the provided path
func (r *Runner) GetScanFiles() {
	r.Logger.Debugf("searching path: %v", r.Options.Path)
	err := filepath.WalkDir(r.Options.Path, r.walkScans)
	if err != nil {
		r.Logger.Errorf("error parsing path: %v", err)
	}
}

// walkScans walks the provided path and queues likely scan files for parsing
func (r *Runner) walkScans(path string, d fs.DirEntry, err error) error {
	// TODO: I'm not sure why this is here
	if err != nil {
		return err
	}
	r.Logger.Debugf("found file: %v", filepath.Base(path))
	if filepath.Ext(path) == ".xml" {
		r.Files = append(r.Files, path)
		r.Logger.Debugf("added file: %v", filepath.Base(path))
	}
	return nil
}
