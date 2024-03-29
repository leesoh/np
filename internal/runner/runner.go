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
	if options.Debug {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.ErrorLevel
	}
	runner := &Runner{
		Logger:  logger,
		Options: options,
	}
	return runner
}

func (r *Runner) Run() {
	r.GetScanFiles()
	res := result.New(r.Logger, r.Options.Exclude, r.Options.CIDR)
	for _, ff := range r.Files {
		r.Logger.Debugf("processing %v", ff)
		// Get last modified time for timeline output
		b, err := ioutil.ReadFile(ff)
		if err != nil {
			r.Logger.Errorf("error reading file: %v", err)
		}
		s := scan.New(b, r.Logger, res)
		// If we don't process this first it will overwrite other scans
		if s.IsNP() {
			r.Logger.Debugf("found np scan: %s", ff)
			// We have to send the exclude list so we don't process previously saved hosts
			s.ParseNP(r.Options.Exclude)
		}
		// This handles both Nmap and Masscan files as they have the same structure
		if s.IsNmap() {
			r.Logger.Debugf("found nmap scan: %s", ff)
			s.ParseNmap()
		}
		if s.IsNaabuV1() {
			r.Logger.Debugf("found naabu v1 scan: %s", ff)
			s.ParseNaabuV1()
		}
		if s.IsNaabuV2() {
			r.Logger.Debugf("found naabu v2 scan: %s", ff)
			s.ParseNaabuV2()
		}
		if s.IsDNSx() {
			r.Logger.Debugf("found DNSx scan: %s", ff)
			s.ParseDNSx()
		}
	}
	// -host
	if r.Options.Host != "" {
		ip := net.ParseIP(r.Options.Host)
		h := &result.Host{}
		if ip == nil {
			h.Name = r.Options.Host
		} else {
			h.IP = ip
		}
		res.PrintHost(h)
		return
	}
	// -hosts
	if r.Options.Hosts {
		res.PrintAlive()
		return
	}
	// -service
	if r.Options.Service != "" {
		res.PrintByService(r.Options.Service)
		return
	}
	// -services
	if r.Options.Services {
		res.PrintServices()
		return
	}
	// -port
	if r.Options.Port != nil {
		res.PrintByPort(r.Options.Port)
		return
	}
	// -ports
	if r.Options.Ports {
		res.PrintPortSummary()
		return
	}
	// -json
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
		r.Logger.Debugf("queued XML file for processing: %v", filepath.Base(path))
	case ".json":
		r.Files = append(r.Files, path)
		r.Logger.Debugf("queued JSON file for processing: %v", filepath.Base(path))
	default:
		r.Logger.Debugf("unsupported file type: %v", filepath.Base(path))
	}
	return nil
}
