package main

import "github.com/leesoh/np/internal/runner"

func main() {
	options := runner.ParseOptions()
	r := runner.NewRunner(options)
	r.Run()
}
