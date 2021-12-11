package runner

import (
	"github.com/Masterminds/log-go"
	"github.com/Masterminds/log-go/impl/cli"
)

type Runner struct {
	Options *Options
	Logger  *cli.Logger
}

func NewRunner(options *Options) *Runner {
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
	r.ReadScans()
}
