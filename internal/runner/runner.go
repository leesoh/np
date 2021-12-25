package runner

import (
	"io/fs"
	"io/ioutil"
	"net"
	"path/filepath"

	"github.com/Masterminds/log-go"
	"github.com/Masterminds/log-go/impl/cli"
	"github.com/leesoh/np/internal/result"
	"github.com/leesoh/np/internal/scan"
)

type Runner struct {
	Files   []string
	Logger  *cli.Logger
	Options *Options
}

func New(options *Options) *Runner {
	logger := cli.NewStandard()
	if options.Verbose {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.FatalLevel
	}
	runner := &Runner{
		Logger:  logger,
		Options: options,
	}
	return runner
}

func (r *Runner) Run() {
	r.GetScanFiles()
	res := result.New(r.Logger)
	for _, ff := range r.Files {
		r.Logger.Debugf("processing %v", ff)
		b, err := ioutil.ReadFile(ff)
		if err != nil {
			r.Logger.Errorf("error reading file: %v", err)
		}
		s := scan.New(b, r.Logger, res)
		if s.IsNmap() {
			r.Logger.Debugf("found nmap scan: %s", ff)
			s.ParseNmap()
			continue
		}
		if s.IsNP() {
			r.Logger.Debugf("found np scan: %s", ff)
			s.ParseNP()
			continue
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
		res.PrintHost(h)
		r.Logger.Debugf("printed host: %v", h)
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
	if r.Options.Port != nil {
		res.PrintByPort(r.Options.Port)
		r.Logger.Debugf("printed by port")
		return
	}
	if r.Options.Ports {
		res.PrintPortSummary()
		return
	}
	if r.Options.JSON {
		res.PrintJSON()
		return
	}
	r.Logger.Debugf("no options selected, printing default")
	res.Print()
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
	switch filepath.Ext(path) {
	case ".xml":
		r.Files = append(r.Files, path)
		r.Logger.Debugf("added XML file: %v", filepath.Base(path))
	case ".json":
		r.Files = append(r.Files, path)
		r.Logger.Debugf("added JSON file: %v", filepath.Base(path))
	}
	return nil
}
