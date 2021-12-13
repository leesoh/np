package scan

import (
	"github.com/Masterminds/log-go/impl/cli"
	"github.com/leesoh/np/internal/result"
)

type Scan struct {
	Bytes  []byte
	Logger *cli.Logger
	Nmap   NmapScan
	Result *result.Result
}

func New(b []byte, logger *cli.Logger, r *result.Result) *Scan {
	return &Scan{
		Bytes:  b,
		Logger: logger,
		Result: r,
	}
}
