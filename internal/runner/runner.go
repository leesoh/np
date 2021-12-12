package runner

import (
	"io/fs"
	"path/filepath"

	"github.com/Masterminds/log-go"
	"github.com/Masterminds/log-go/impl/cli"
	"github.com/leesoh/np/internal/scan"
)

type Runner struct {
	Options *Options
	Logger  *cli.Logger
	Files   []string
}

func New(options *Options) *Runner {
	logger := cli.NewStandard()
	if options.Verbose {
		logger.Level = log.DebugLevel
	} else {
		logger.Level = log.FatalLevel
	}
	runner := &Runner{
		Options: options,
		Logger:  logger,
	}
	return runner
}

func (r *Runner) Run() {
	r.GetScanFiles()
	r.LoadScans()
}

func (r *Runner) GetScanFiles() {
	r.Logger.Debugf("searching path: %v", r.Options.Path)
	err := filepath.WalkDir(r.Options.Path, func(path string, d fs.DirEntry, err error) error {
		r.Logger.Debugf("found file: %v", filepath.Base(path))
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".xml" {
			r.Files = append(r.Files, path)
			r.Logger.Debugf("added file: %v", filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		r.Logger.Errorf("error parsing path: %v", err)
	}
}

func (r *Runner) LoadScans() {
	for _, f := range r.Files {
		sc := scan.New(r.Logger)
		sc.Unmarshal(f)
		sc.Print()
	}
}
