package parse

import "github.com/Masterminds/log-go/impl/cli"

type Parser struct {
	Logger *cli.Logger
}

func NewParser(logger *cli.Logger) *Parser {
	return &Parser{
		Logger: logger,
	}
}
