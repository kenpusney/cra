package main

import (
	"fmt"
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

	fmt.Println("Started at port:", opts.Port, "Proxying:", opts.Endpoint)
	_ = context.Serve()
}
