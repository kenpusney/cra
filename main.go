package main

import (
	"github.com/alexflint/go-arg"
	"github.com/kenpusney/cra/core"
	"os"
)

func main() {
	var opts core.Opts

	p := arg.MustParse(&opts)

	if opts.Endpoint == "" {
		p.WriteHelp(os.Stdout)
		return
	}

	context := core.NewContext(&opts)

	_ = context.Serve()
}
