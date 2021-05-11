package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/kenpusney/cra/core/boot"
	"github.com/kenpusney/cra/core/common"
	"os"
)

func main() {
	var opts common.Opts

	p := arg.MustParse(&opts)

	if opts.Endpoint == "" {
		p.WriteHelp(os.Stdout)
		return
	}

	context := boot.NewContext(&opts)

	fmt.Println("Started at port:", opts.Port, "Proxying:", opts.Endpoint)
	_ = context.Serve()
}
